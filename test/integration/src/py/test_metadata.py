import logging

import pytest
import hashlib

import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from pytest_utilities import generate_test_account
from test_utilities import get_required_env_var, get_shell_output, amount_in_wei, \
    SifchaincliCredentials, get_token_metadata, calculate_denom_hash

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
bridgetoken_address = get_required_env_var("BRIDGE_TOKEN_ADDRESS")


def do_currency_test(
        new_currency_symbol,
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        solidity_json_path,
):
    amount = amount_in_wei(9)
    logging.info(f"create new currency")
    new_currency = test_utilities.create_new_currency(
        amount,
        new_currency_symbol,
        new_currency_symbol,
        18,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path
    )

    logging.info(
        f"create test account to use with new currency {new_currency_symbol}")
    basic_transfer_request.ethereum_address = source_ethereum_address
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=10 ** 17,
        target_rowan_balance=10 ** 18
    )
    test_amount = 39000
    logging.info(
        f"transfer some of the new currency {new_currency_symbol} to the test sifchain address")
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.sifchain_symbol = ("c" + new_currency["newtoken_symbol"]).lower()
    request.amount = test_amount
    burn_lock_functions.transfer_ethereum_to_sifchain(request)

    logging.info("send some new currency to ethereum")
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.amount = test_amount - 1
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)

    # Validate that the new token metadata is available to the metadata module
    network_descriptor = 1
    token_contract_address = new_currency["newtoken_address"]
    decimals = int(new_currency["newtoken_decimals"])
    token_name = new_currency["newtoken_name"]
    token_symbol = new_currency["newtoken_symbol"]
    denom_hash = calculate_denom_hash(
        network_descriptor, token_contract_address, decimals, token_name, token_symbol)
    token_metadata = get_token_metadata(denom_hash)
    assert int(token_metadata["decimals"]) == decimals
    assert token_metadata["name"] == token_name
    assert token_metadata["symbol"] == token_symbol
    assert token_metadata["token_address"] == token_contract_address
    assert int(token_metadata["network_descriptor"]) == network_descriptor


@pytest.mark.usefixtures("operator_private_key")
def test_transfer_tokens_with_some_currency(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        solidity_json_path,
):
    new_currency_symbol = (
        "a" + get_shell_output("uuidgen").replace("-", ""))[:4]
    do_currency_test(
        new_currency_symbol,
        basic_transfer_request,
        source_ethereum_address,
        rowan_source_integrationtest_env_credentials,
        rowan_source_integrationtest_env_transfer_request,
        ethereum_network,
        solidity_json_path=solidity_json_path
    )


@pytest.mark.usefixtures("operator_private_key")
def test_three_letter_currency_with_capitals_in_name(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        solidity_json_path,
):
    new_currency_symbol = (
        "F" + get_shell_output("uuidgen").replace("-", ""))[:3]
    do_currency_test(
        new_currency_symbol,
        basic_transfer_request,
        source_ethereum_address,
        rowan_source_integrationtest_env_credentials,
        rowan_source_integrationtest_env_transfer_request,
        ethereum_network,
        solidity_json_path=solidity_json_path
    )
