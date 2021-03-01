import copy
import json
import logging
import subprocess

import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def test_get_all_accounts_with_ceth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
):
    cmd="sifnodecli keys list --keyring-backend test -o json | jq -r '.[].address' | parallel sifnodecli q auth account --node tcp://44.241.55.154:26657 -o json {} | grep coins"
    sub = subprocess.run(cmd, shell=True, capture_output=True)
    output = sub.stdout.decode("utf-8").rstrip()
    result = []
    for line in output.split('\n'):
        json_output = json.loads(line)
        result.append(json_output)
    logging.info(f"srx: {json.dumps(result)}")