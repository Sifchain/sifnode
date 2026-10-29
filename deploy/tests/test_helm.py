"""Render both cloud configurations with real Helm and inspect Kubernetes objects."""

from pathlib import Path
import subprocess
import unittest

import yaml


CHART = Path(__file__).resolve().parents[1] / "helm/sifnode"


class HelmTest(unittest.TestCase):
    def render(self, provider="azure", *extra, valid=True):
        command = ["helm", "template", "test", str(CHART), "--namespace", "testnet",
                   "--set-string", f"cloudProvider={provider}",
                   "--set-string", "image.repository=sifnode-test",
                   "--set-string", "image.tag=local",
                   "--set-string", "chain.id=localnet",
                   "--set-string", "genesis.existingConfigMap=reviewed-genesis", *extra]
        result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        if not valid:
            self.assertNotEqual(result.returncode, 0)
            return result
        self.assertEqual(result.returncode, 0, result.stderr)
        return list(yaml.safe_load_all(result.stdout))

    def statefulset(self, objects):
        return next(obj for obj in objects if obj["kind"] == "StatefulSet")

    def test_azure_persistent_disk_and_single_signer(self):
        objects = self.render()
        sts = self.statefulset(objects)
        self.assertEqual(sts["spec"]["replicas"], 1)
        claim = sts["spec"]["volumeClaimTemplates"][0]["spec"]
        self.assertEqual(claim["storageClassName"], "managed-csi-premium")
        self.assertEqual(claim["accessModes"], ["ReadWriteOnce"])
        self.assertEqual(claim["resources"]["requests"]["storage"], "500Gi")
        self.assertFalse(any(obj["kind"] == "StorageClass" for obj in objects))

    def test_aws_ebs_csi_storage(self):
        objects = self.render("aws")
        storage = next(obj for obj in objects if obj["kind"] == "StorageClass")
        self.assertEqual(storage["provisioner"], "ebs.csi.aws.com")
        self.assertEqual(storage["parameters"], {"type": "gp3", "encrypted": "true"})
        self.assertEqual(storage["volumeBindingMode"], "WaitForFirstConsumer")
        self.assertEqual(storage["reclaimPolicy"], "Retain")
        claim = self.statefulset(objects)["spec"]["volumeClaimTemplates"][0]
        self.assertEqual(claim["spec"]["storageClassName"], storage["metadata"]["name"])

    def test_only_p2p_is_public(self):
        for provider in ("azure", "aws"):
            with self.subTest(provider=provider):
                services = [obj for obj in self.render(provider) if obj["kind"] == "Service"]
                external = [obj for obj in services if obj["spec"].get("type") == "LoadBalancer"]
                self.assertEqual(len(external), 1)
                self.assertEqual([p["port"] for p in external[0]["spec"]["ports"]], [26656])
                rpc = next(obj for obj in services if obj["metadata"]["name"].endswith("-rpc"))
                self.assertEqual(rpc["spec"]["type"], "ClusterIP")
                self.assertNotIn("annotations", external[0]["metadata"], "no hard-coded AWS annotations on Azure")

    def test_bootstrap_script_matches_tested_source(self):
        objects = self.render()
        config = next(obj for obj in objects if obj["kind"] == "ConfigMap")
        self.assertEqual(config["data"]["init.sh"].rstrip(), (CHART / "scripts/init.sh").read_text().rstrip())
        pod = self.statefulset(objects)["spec"]["template"]["spec"]
        self.assertFalse(pod["automountServiceAccountToken"])
        self.assertEqual(pod["securityContext"]["runAsUser"], 1000)
        for container in pod["containers"] + pod["initContainers"]:
            self.assertTrue(container["securityContext"]["readOnlyRootFilesystem"])
            self.assertFalse(container["securityContext"]["allowPrivilegeEscalation"])
        self.assertFalse(any(obj["kind"] in ("Secret", "Role", "RoleBinding") for obj in objects))

    def test_existing_volume_and_local_storage_overrides(self):
        objects = self.render("aws", "--set-string", "persistence.existingClaim=restored-data")
        sts = self.statefulset(objects)
        self.assertNotIn("volumeClaimTemplates", sts["spec"])
        volumes = sts["spec"]["template"]["spec"]["volumes"]
        self.assertIn({"name": "data", "persistentVolumeClaim": {"claimName": "restored-data"}}, volumes)
        self.assertFalse(any(obj["kind"] == "StorageClass" for obj in objects))
        local = self.render("azure", "--set-string", "persistence.storageClassName=standard")
        claim = self.statefulset(local)["spec"]["volumeClaimTemplates"][0]
        self.assertEqual(claim["spec"]["storageClassName"], "standard")

    def test_network_options_are_arguments_not_shell_commands(self):
        objects = self.render("azure", "--set-string", "chain.persistentPeers=id@example:26656",
                              "--set-string", "chain.externalAddress=203.0.113.1:26656")
        container = self.statefulset(objects)["spec"]["template"]["spec"]["containers"][0]
        self.assertEqual(container["command"], ["sifnoded"])
        self.assertIn("--p2p.persistent_peers=id@example:26656", container["args"])
        self.assertIn("--p2p.external-address=203.0.113.1:26656", container["args"])

    def test_large_genesis_can_use_a_readonly_pvc(self):
        objects = self.render("azure", "--set-string", "genesis.existingConfigMap=",
                              "--set-string", "genesis.existingClaim=large-genesis")
        pod = self.statefulset(objects)["spec"]["template"]["spec"]
        volume = next(item for item in pod["volumes"] if item["name"] == "genesis")
        self.assertEqual(volume["persistentVolumeClaim"], {"claimName": "large-genesis", "readOnly": True})

    def test_invalid_values_are_rejected_before_deployment(self):
        for key, value in [("cloudProvider", "gcp"), ("image.tag", ""), ("image.repository", ""),
                           ("chain.id", ""), ("genesis.existingConfigMap", ""),
                           ("persistence.size", "0Gi"), ("service.type", "ExternalName")]:
            with self.subTest(key=key):
                self.render("azure", "--set-string", f"{key}={value}", valid=False)


if __name__ == "__main__":
    unittest.main()
