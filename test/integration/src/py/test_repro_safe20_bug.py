import copy
import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output, amount_in_wei, \
    run_yarn_command, ganache_accounts, get_sifchain_addr_balance

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")


def create_new_currency(amount, symbol):
    return run_yarn_command(
        f"yarn --cwd {smart_contracts_dir} "
        f"integrationtest:enableNewToken "
        f"--bridgebank_address {bridgebank_address} "
        f"--symbol {symbol} "
        f"--amount {amount} "
        f"--limit_amount {amount}"
    )


def build_request(new_currency, amount):
    accounts = ganache_accounts(smart_contracts_dir=smart_contracts_dir)
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    request = EthereumToSifchainTransferRequest(
        ethereum_symbol=new_currency["newtoken_address"],
        sifchain_symbol="c" + new_currency["newtoken_symbol"],
        sifchain_address=get_required_env_var("OWNER_ADDR"),
        smart_contracts_dir=smart_contracts_dir,
        ethereum_address=accounts["accounts"][0],
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=amount,
        ceth_amount=2 * (10 ** 16)
    )

    return request


def test_reproduce_safe_erc20_bug():
    new_account_key = "Foo"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key)
    request1 = build_request(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)

    logging.info("get ceth to pay lock/burn fees")
    eth_request: EthereumToSifchainTransferRequest = copy.deepcopy(request1)
    eth_request.sifchain_symbol = "ceth"
    eth_request.ethereum_symbol = "eth"
    eth_request.ethereum_address = bridgebank_address
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 5)

    return_request: EthereumToSifchainTransferRequest = copy.deepcopy(request1)
    burn_lock_functions.transfer_sifchain_to_ethereum(return_request)
