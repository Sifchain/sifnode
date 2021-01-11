import logging
import os

# to test against ropsten, define:
# ETHEREUM_ADDRESS
# ETHEREUM_PRIVATE_KEY
# ETHEREUM_NETWORK = ropsten
#
from copy import copy
from json import JSONDecodeError

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output, SifchaincliCredentials, amount_in_wei

request: EthereumToSifchainTransferRequest
credentials: SifchaincliCredentials


def test_transfer_eth_to_ceth():
    global request
    global credentials
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)

    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=get_required_env_var("SMART_CONTRACTS_DIR"),
        ethereum_address=get_required_env_var("ETHEREUM_ADDRESS"),
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=90000
    )

    result = burn_lock_functions.transfer_ethereum_to_sifchain(request, 50)
    logging.info(f"transfer_ethereum_to_sifchain_result_json: {result}")


def test_transfer_eth_to_ceth_over_limit():
    global request
    global credentials
    assert credentials is not None, "run test_transfer_eth_to_ceth first"
    assert request is not None
    invalid_request = copy(request)
    invalid_request.amount = amount_in_wei(35)
    with pytest.raises(JSONDecodeError):
        burn_lock_functions.transfer_ethereum_to_sifchain(invalid_request, credentials)


def test_transfer_ceth_to_eth():
    global request
    global credentials
    assert credentials is not None, "run test_transfer_eth_to_ceth first"
    assert request is not None
    return_result = burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    logging.info(f"transfer_sifchain_to_ethereum__result_json: {return_result}")


