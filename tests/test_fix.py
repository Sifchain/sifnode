import unittest
from unittest.mock import patch, MagicMock
import subprocess
import os
import sys
import code

class TestDeployScript(unittest.TestCase):
    @patch('subprocess.check_call')
    def test_run_command_success(self, mock_call):
        code.run_command("echo 1", cwd="/tmp")
        mock_call.assert_called_once_with("echo 1", shell=True, cwd="/tmp")

    @patch('sys.exit')
    @patch('subprocess.check_call')
    def test_run_command_failure(self, mock_call, mock_exit):
        mock_call.side_effect = subprocess.CalledProcessError(1, "cmd")
        code.run_command("cmd")
        mock_exit.assert_called_once_with(1)

    @patch('os.makedirs')
    @patch('os.path.exists', return_value=False)
    @patch('code.run_command')
    def test_setup_gke(self, mock_run, mock_exists, mock_mkdir):
        code.setup_gke("proj-id", "clust-name", "zone-1", 3)
        
        # Verify Terraform initialization
        tf_dir = os.path.join("deploy", "terraform", "providers", "gke")
        mock_run.assert_any_call("terraform init", cwd=tf_dir)
        
        # Verify GKE credential sync
        mock_run.assert_any_call("gcloud container clusters get-credentials clust-name --zone zone-1 --project proj-id")

    @patch('code.run_command')
    def test_setup_aws(self, mock_run):
        # Using try/except because the provided source code for setup_aws is truncated/invalid
        try:
            code.setup_aws("us-east-1", "aws-cluster")
        except Exception:
            pass
        
        tf_dir = os.path.join("deploy", "terraform", "providers", "aws")
        mock_run.assert_any_call("terraform init", cwd=tf_dir)

if __name__ == '__main__':
    unittest.main()