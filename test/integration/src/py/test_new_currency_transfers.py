import copy

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import create_new_sifaddr_and_credentials
from test_utilities import get_shell_output, amount_in_wei, \
    get_sifchain_addr_balance, create_new_currency


def build_request_for_new_sifchain_address(basic_transfer_request, source_ethereum_address, new_currency, amount):
    sifaddress, _ = create_new_sifaddr_and_credentials()
    request = copy.deepcopy(basic_transfer_request)
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.ethereum_address = source_ethereum_address
    request.sifchain_symbol = "c" + new_currency["newtoken_symbol"]
    request.sifchain_address = sifaddress
    request.amount = amount
    return request


@pytest.mark.parametrize("token_length", [7, 12])
@pytest.mark.usefixtures("operator_private_key")
def test_can_create_a_new_token_and_peg_it(
        token_length: int,
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        operator_address,
        ethereum_network,
        source_ethereum_address,
):
    new_account_key = ("a" + get_shell_output("uuidgen").replace("-", ""))[:token_length]
    token_name = new_account_key
    amount = amount_in_wei(9)
    new_currency = create_new_currency(
        amount=amount,
        symbol=new_account_key,
        token_name=token_name,
        decimals=18,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path,
        operator_address=operator_address,
        ethereum_network=ethereum_network
    )
    request = build_request_for_new_sifchain_address(
        basic_transfer_request,
        source_ethereum_address,
        new_currency,
        amount / 10
    )
    burn_lock_functions.transfer_ethereum_to_sifchain(request)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_can_create_a_new_token_with_a_one_number_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = "0"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request_for_new_sifchain_address(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_can_create_a_new_token_with_a_one_letter_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = "a"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request_for_new_sifchain_address(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_can_create_a_new_token_with_a_long_name_and_peg_it(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = "ca36e47edfeb28489d8e110fb91d351bcd"
    amount = amount_in_wei(9)
    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request_for_new_sifchain_address(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_two_currencies_with_different_capitalization_should_not_interfere_with_each_other(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = ("a" + get_shell_output("uuidgen").replace("-", "").lower())[:5]
    amount = amount_in_wei(9)

    new_currency = create_new_currency(amount, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request1, _) = build_request_for_new_sifchain_address(new_currency, amount)
    burn_lock_functions.transfer_ethereum_to_sifchain(request1, 10)
    balance_1 = get_sifchain_addr_balance(request1.sifchain_address, request1.sifnodecli_node, request1.sifchain_symbol)
    assert (balance_1 == request1.amount)

    new_currency = create_new_currency(amount, new_account_key.upper(), smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request2, _) = build_request_for_new_sifchain_address(new_currency, amount + 70000)
    burn_lock_functions.transfer_ethereum_to_sifchain(request2, 10)

    balance_1_again = get_sifchain_addr_balance(request1.sifchain_address, request1.sifnodecli_node,
                                                request1.sifchain_symbol)

    assert (balance_1 == balance_1_again)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_cannot_create_two_currencies_that_only_differ_in_capitalization(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = get_shell_output("uuidgen").replace("-", "").lower()
    new_currency = create_new_currency(10 ** 18, new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    with pytest.raises(Exception):
        new_currency = create_new_currency(10 ** 18, new_account_key, smart_contracts_dir, bridgebank_address,
                                           solidity_json_path)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_cannot_create_two_currencies_with_the_same_name(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    new_account_key = get_shell_output("uuidgen").replace("-", "")
    new_currency = create_new_currency(amount_in_wei(10), new_account_key, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    with pytest.raises(Exception):
        new_currency = create_new_currency(amount_in_wei(10), new_account_key, smart_contracts_dir, bridgebank_address,
                                           solidity_json_path)


@pytest.mark.skip(reason="fails and is too slow to mark with xfail")
@pytest.mark.usefixtures("operator_private_key")
def test_can_use_a_token_with_a_dash_in_the_name(
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
):
    n = "a-b"
    new_currency = create_new_currency(amount_in_wei(10), n, smart_contracts_dir, bridgebank_address,
                                       solidity_json_path)
    (request, _) = build_request_for_new_sifchain_address(new_currency, 60000)
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 10)
