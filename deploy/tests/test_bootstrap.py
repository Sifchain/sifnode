"""Run the chart's exact bootstrap script against an operator-supplied real binary."""

import os
from pathlib import Path
import subprocess
import tempfile
import unittest


SCRIPT = Path(__file__).resolve().parents[1] / "helm/sifnode/scripts/init.sh"
BINARY = os.environ.get("SIFNODE_TEST_BINARY")


@unittest.skipUnless(BINARY, "set SIFNODE_TEST_BINARY to a compatible real sifnoded binary")
class BootstrapTest(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="sifnode-bootstrap-")
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        self.node = self.root / "node"
        self.node.mkdir()
        self.bin = self.root / "bin"
        self.bin.mkdir()
        (self.bin / "sifnoded").symlink_to(Path(BINARY).resolve())
        source = self.root / "fixture"
        result = subprocess.run([BINARY, "init", "fixture", "--chain-id", "localnet", "--home", str(source)],
                                capture_output=True, text=True)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.genesis = self.root / "genesis.json"
        self.genesis.write_bytes((source / "config/genesis.json").read_bytes())
        self.env = dict(os.environ, PATH=f"{self.bin}{os.pathsep}{os.environ['PATH']}",
                        SIFNODE_HOME=str(self.node), GENESIS_FILE=str(self.genesis),
                        CHAIN_ID="localnet", MONIKER="test-node")

    def bootstrap(self):
        return subprocess.run(["sh", str(SCRIPT)], env=self.env, capture_output=True, text=True)

    def test_initializes_and_validates_real_genesis(self):
        result = self.bootstrap()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue((self.node / ".initialized").is_file())
        self.assertEqual((self.node / "config/genesis.json").read_bytes(), self.genesis.read_bytes())
        self.assertTrue((self.node / "config/priv_validator_key.json").is_file())

    def test_restart_preserves_keys_genesis_and_consensus_state(self):
        self.assertEqual(self.bootstrap().returncode, 0)
        paths = [self.node / item for item in ("config/priv_validator_key.json", "config/node_key.json",
                                                "config/genesis.json", "data/priv_validator_state.json")]
        before = [path.read_bytes() for path in paths]
        self.genesis.write_text("invalid replacement genesis")
        self.assertEqual(self.bootstrap().returncode, 0)
        self.assertEqual([path.read_bytes() for path in paths], before)

    def test_invalid_genesis_never_marks_initialized_and_can_retry(self):
        valid = self.genesis.read_bytes()
        self.genesis.write_text("invalid json")
        self.assertNotEqual(self.bootstrap().returncode, 0)
        self.assertFalse((self.node / ".initialized").exists())
        key = (self.node / "config/priv_validator_key.json").read_bytes()
        self.genesis.write_bytes(valid)
        self.assertEqual(self.bootstrap().returncode, 0)
        self.assertEqual((self.node / "config/priv_validator_key.json").read_bytes(), key)

    def test_refuses_unmarked_existing_application_database(self):
        (self.node / "data/application.db").mkdir(parents=True)
        result = self.bootstrap()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Refusing to replace genesis", result.stderr)
        self.assertFalse((self.node / "config/priv_validator_key.json").exists())


if __name__ == "__main__":
    unittest.main()
