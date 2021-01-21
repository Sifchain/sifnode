import copy
import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output, amount_in_wei, \
    run_yarn_command, ganache_accounts, get_sifchain_addr_balance, get_eth_balance

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
bridgetoken_address = get_required_env_var("BRIDGE_TOKEN_ADDRESS")


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
        sifchain_address=new_addr["address"],
        smart_contracts_dir=smart_contracts_dir,
        ethereum_address=accounts["accounts"][0],
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=bridgebank_address,
        bridgetoken_address=bridgetoken_address,
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=amount,
        ceth_amount=2 * (10 ** 16)
    )

    return (request, credentials)


def test_transfer_tokens_with_a_capital_letter_in_the_name():
    new_account_key = "Foo"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key)
    logging.info(f"new_currency is {new_currency}")
    (foo_request, credentials) = build_request(new_currency, amount)

    starting_balance = get_eth_balance(foo_request)
    assert(starting_balance == amount)
    logging.info(f"starting balance for {foo_request} is {starting_balance}")
    burn_lock_functions.transfer_ethereum_to_sifchain(foo_request, 10)
    ending_balance = get_eth_balance(foo_request)
    assert(starting_balance == amount)
    assert(ending_balance == 0)

    logging.info(f"ending balance for {foo_request} is {ending_balance}")

    logging.info("get ceth to pay lock/burn fees")
    accounts = ganache_accounts(smart_contracts_dir=smart_contracts_dir)
    eth_request: EthereumToSifchainTransferRequest = copy.deepcopy(foo_request)
    eth_request.sifchain_symbol = "ceth"
    eth_request.ethereum_symbol = "eth"
    eth_request.ethereum_address=accounts["accounts"][0]
    eth_request.bridgebank_address=bridgebank_address
    eth_request.bridgetoken_address=bridgetoken_address
    burn_lock_functions.transfer_ethereum_to_sifchain(eth_request, 5)

    logging.info(f"request balances using {foo_request}")
    get_eth_balance(foo_request)
    get_sifchain_addr_balance(foo_request.sifchain_address, foo_request.sifnodecli_node, "ceth")

    logging.info("sending cFoo back to ethereum")
    return_request: EthereumToSifchainTransferRequest = copy.deepcopy(foo_request)
    return_request.sifchain_symbol = return_request.sifchain_symbol.lower()
    burn_lock_functions.transfer_sifchain_to_ethereum(return_request, credentials)
