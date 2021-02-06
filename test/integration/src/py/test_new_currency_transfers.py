import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output, amount_in_wei, \
    ganache_accounts, get_sifchain_addr_balance, create_new_currency

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
bridgetoken_address = get_required_env_var("BRIDGE_TOKEN_ADDRESS")


def build_request(new_currency, amount, solidity_json_path = None):
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
        bridgebank_address=bridgebank_address,
        bridgetoken_address=bridgetoken_address,
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=amount,
        solidity_json_path=solidity_json_path,
        ceth_amount=2 * (10 ** 16)
    )

    return (request, credentials)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_can_create_a_new_token_with_a_one_number_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = "0"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_can_create_a_new_token_with_a_one_letter_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = "a"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_can_create_a_new_token_with_a_long_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = "ca36e47edfeb28489d8e110fb91d351bcd"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


def test_can_create_a_new_token_with_a_7_char_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = ("a" + get_shell_output("uuidgen").replace("-", ""))[:7]
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request(new_currency, amount, solidity_json_path)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_two_currencies_with_different_capitalization_should_not_interfere_with_each_other(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = ("a" + get_shell_output("uuidgen").replace("-", "").lower())[:5]
    amount = amount_in_wei(9)

    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)
    balance_1 = get_sifchain_addr_balance(request1.sifchain_address, request1.sifnodecli_node, request1.sifchain_symbol)
    assert (balance_1 == request1.amount)

    new_currency = create_new_currency(amount, new_account_key.upper(), smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request2, _) = build_request(new_currency, amount + 70000)
    burn_lock_functions.transfer_ethereum_to_sifchain(request2, 10)

    balance_1_again = get_sifchain_addr_balance(request1.sifchain_address, request1.sifnodecli_node,
                                                request1.sifchain_symbol)

    assert (balance_1 == balance_1_again)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_cannot_create_two_currencies_that_only_differ_in_capitalization(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = get_shell_output("uuidgen").replace("-", "").lower()
    new_currency = create_new_currency(10 ** 18, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    with pytest.raises(Exception):
        create_new_currency(amount_in_wei(10), new_account_key.upper())


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_cannot_create_two_currencies_with_the_same_name(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = get_shell_output("uuidgen").replace("-", "")
    create_new_currency(amount_in_wei(10), new_account_key)
    with pytest.raises(Exception):
        create_new_currency(amount_in_wei(10), new_account_key)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
def test_can_use_a_token_with_a_dash_in_the_name(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    n = "a-b"
    new_currency = create_new_currency(amount_in_wei(10), n)
    (request, _) = build_request(new_currency, 60000)
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 10)
