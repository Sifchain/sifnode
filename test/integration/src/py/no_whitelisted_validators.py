# This test needs an environment where you have no whitelisted validators.
#
# the ./execute_integration_tests_whitelisted_validators.sh script sets that up and runs this test.

import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output


def test_transfer_eth_to_ceth_without_a_validator_should_throw_exception():
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=get_required_env_var("SMART_CONTRACTS_DIR"),
        ethereum_address=get_required_env_var("ETHEREUM_ADDRESS"),
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=90000
    )

    logging.info("try to transfer, but expect a failure since there are no whitelisted validators")
    with pytest.raises(Exception):
        # use a small number for max_retries - on a local system, it shouldn't
        # take more than a second or two for ebrelayer to act
        burn_lock_functions.transfer_ethereum_to_sifchain(request, 3)
