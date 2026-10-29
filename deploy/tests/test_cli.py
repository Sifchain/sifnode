"""Exercise the CLI without cloud credentials or provisioning resources."""

import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest


DEPLOY = Path(os.environ.get("SIFNODE_TEST_DEPLOY", Path(__file__).resolve().parents[1]))


class ClusterCLITest(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="sifnode-cli-")
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        self.directory = self.root / "environment with spaces"
        self.bin = self.root / "bin"
        self.bin.mkdir()
        self.record = self.root / "commands.jsonl"
        # Transport doubles only: Terraform/Helm semantics are checked separately.
        script = '''#!/usr/bin/env python3
import json, os, sys
with open(os.environ["COMMAND_RECORD"], "a") as stream:
    stream.write(json.dumps(sys.argv) + "\\n")
if "output" in sys.argv:
    print("apiVersion: v1\\nkind: Config\\n# synthetic-test-config")
sys.exit(int(os.environ.get("TOOL_EXIT_CODE", "0")))
'''
        for name in ("terraform", "helm"):
            executable = self.bin / name
            executable.write_text(script)
            executable.chmod(0o755)
        self.env = dict(os.environ, PATH=f"{self.bin}{os.pathsep}{os.environ['PATH']}",
                        COMMAND_RECORD=str(self.record))

    def cli(self, command, *extra, provider="azure"):
        return subprocess.run(
            ["python3", str(DEPLOY / "cluster.py"), command, "--provider", provider,
             "--name", "testnet", "--directory", str(self.directory), *map(str, extra)],
            env=self.env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)

    def scaffold(self, provider="azure"):
        result = self.cli("scaffold", "--admin-cidr", "192.0.2.1/32", provider=provider)
        self.assertEqual(result.returncode, 0, result.stderr)

    def commands(self):
        return [json.loads(line) for line in self.record.read_text().splitlines()]

    def test_scaffold_azure(self):
        self.scaffold()
        main = (self.directory / "main.tf").read_text()
        self.assertIn(str(DEPLOY / "terraform/providers/azure"), main)
        self.assertNotIn("__MODULE_SOURCE__", main)
        variables = json.loads((self.directory / "terraform.tfvars.json").read_text())
        self.assertEqual(variables["location"], "eastus")
        self.assertEqual(variables["cluster_name"], "sifchain-azure-testnet")
        self.assertEqual(variables["api_server_authorized_ip_ranges"], ["192.0.2.1/32"])
        self.assertFalse(self.record.exists(), "scaffold must not provision anything")

    def test_scaffold_aws(self):
        self.scaffold("aws")
        self.assertIn(str(DEPLOY / "terraform/providers/aws"), (self.directory / "main.tf").read_text())
        variables = json.loads((self.directory / "terraform.tfvars.json").read_text())
        self.assertEqual(variables["region"], "us-west-2")
        self.assertNotIn("location", variables)

    def test_custom_region_version_and_multiple_admin_ranges(self):
        result = self.cli("scaffold", "--location", "westeurope", "--kubernetes-version", "1.33",
                          "--admin-cidr", "192.0.2.1/32", "--admin-cidr", "198.51.100.0/24")
        self.assertEqual(result.returncode, 0, result.stderr)
        variables = json.loads((self.directory / "terraform.tfvars.json").read_text())
        self.assertEqual(variables["location"], "westeurope")
        self.assertEqual(variables["kubernetes_version"], "1.33")
        self.assertEqual(len(variables["api_server_authorized_ip_ranges"]), 2)

    def test_scaffold_never_overwrites_existing_environment(self):
        self.scaffold()
        marker = self.directory / "terraform.tfstate"
        marker.write_text("existing state")
        self.assertNotEqual(self.cli("scaffold", "--admin-cidr", "192.0.2.1/32").returncode, 0)
        self.assertEqual(marker.read_text(), "existing state")

    def test_requires_admin_cidr_before_writing(self):
        self.assertNotEqual(self.cli("scaffold").returncode, 0)
        self.assertFalse(self.directory.exists())

    def test_rejects_unsafe_or_invalid_ranges(self):
        for value in ("0.0.0.0/0", "::/0", "garbage", "192.0.2.1/24"):
            with self.subTest(value=value):
                self.assertNotEqual(self.cli("scaffold", "--admin-cidr", value).returncode, 0)
        self.assertFalse(self.directory.exists())

    def test_rejects_unknown_provider_and_path_names(self):
        self.assertNotEqual(self.cli("scaffold", provider="gcp").returncode, 0)
        self.assertNotEqual(self.cli("scaffold", "--name", "../other").returncode, 0)

    def test_refuses_wrong_provider_for_existing_state(self):
        self.scaffold()
        self.assertNotEqual(self.cli("destroy", provider="aws").returncode, 0)
        self.assertFalse(self.record.exists())

    def test_apply_and_destroy_keep_interactive_confirmation(self):
        self.scaffold()
        for action in ("apply", "destroy"):
            self.assertEqual(self.cli(action).returncode, 0)
            self.assertEqual(self.commands()[-1][1:], [f"-chdir={self.directory}", action])

    def test_propagates_tool_failure(self):
        self.scaffold()
        self.env["TOOL_EXIT_CODE"] = "7"
        self.assertNotEqual(self.cli("apply").returncode, 0)

    def test_private_kubeconfig_without_stdout_disclosure(self):
        self.scaffold()
        result = self.cli("kubeconfig")
        self.assertEqual(result.returncode, 0, result.stderr)
        config = self.directory / "kubeconfig"
        self.assertEqual(config.stat().st_mode & 0o777, 0o600)
        self.assertIn("synthetic-test-config", config.read_text())
        self.assertNotIn("synthetic-test-config", result.stdout)
        self.assertEqual(self.cli("kubeconfig").returncode, 0, "regeneration should replace atomically")

    def test_kubeconfig_does_not_follow_temporary_symlink(self):
        self.scaffold()
        target = self.root / "unrelated"
        target.write_text("keep me")
        (self.directory / "kubeconfig.tmp").symlink_to(target)
        self.assertNotEqual(self.cli("kubeconfig").returncode, 0)
        self.assertEqual(target.read_text(), "keep me")

    def test_node_requires_explicit_environment_kubeconfig(self):
        self.scaffold()
        self.assertNotEqual(self.cli("node").returncode, 0)
        self.assertFalse(self.record.exists())

    def test_node_uses_explicit_cluster_and_selected_provider(self):
        self.scaffold()
        self.assertEqual(self.cli("kubeconfig").returncode, 0)
        values = self.root / "node values.yaml"
        values.write_text("cloudProvider: aws\n")
        result = self.cli("node", "--values", values)
        self.assertEqual(result.returncode, 0, result.stderr)
        command = self.commands()[-1]
        self.assertEqual(command[1:3], ["--kubeconfig", str(self.directory / "kubeconfig")])
        self.assertEqual(command[command.index("--values") + 1], str(values))
        self.assertGreater(command.index("cloudProvider=azure"), command.index("--values"))
        self.assertIn("--wait", command)

    def test_node_requires_values(self):
        self.scaffold()
        self.assertEqual(self.cli("kubeconfig").returncode, 0)
        self.assertNotEqual(self.cli("node").returncode, 0)
        self.assertEqual(len(self.commands()), 1)


if __name__ == "__main__":
    unittest.main()
