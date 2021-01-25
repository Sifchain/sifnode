import logging
import os
from copy import copy, deepcopy
from functools import lru_cache
from json import JSONDecodeError

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
import test_utilities
from test_utilities import get_required_env_var, get_shell_output, SifchaincliCredentials, amount_in_wei, \
    get_optional_env_var, ganache_owner_account

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")

ethereum_address = get_optional_env_var(
    "ETHEREUM_ADDRESS",
    ganache_owner_account(smart_contracts_dir)
)


def build_request() -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=smart_contracts_dir,
        ethereum_address=ethereum_address,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=9 * 10 ** 18,
        ceth_amount=2 * (10 ** 16)
    )
    return request, credentials


def test_transfer_eth_to_ceth_and_back():
    # def whitelist_token(token: str, smart_contracts_dir: str, setting:bool = True):
    request, credentials = build_request()
    test_utilities.whitelist_token(request.ethereum_symbol, request.smart_contracts_dir, True)
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 15)
    logging.info(f"send ceth back to {request.ethereum_address}")
    return_request = deepcopy(request)
    return_request.amount = 20000
    burn_lock_functions.transfer_sifchain_to_ethereum(return_request, credentials)


def test_transfer_eth_to_ceth_over_limit():
    request, credentials = build_request()
    invalid_request = copy(request)
    invalid_request.amount = amount_in_wei(35)
    with pytest.raises(JSONDecodeError):
        burn_lock_functions.transfer_ethereum_to_sifchain(invalid_request, credentials)
