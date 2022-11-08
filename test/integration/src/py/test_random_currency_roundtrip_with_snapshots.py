import logging

import pytest
import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from pytest_utilities import generate_test_account
from test_utilities import get_shell_output, amount_in_wei, SifchaincliCredentials


def do_currency_test(
    ctx,
    new_currency_symbol,
    basic_transfer_request: EthereumToSifchainTransferRequest,
    source_ethereum_address: str,
    rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
    rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
):
    amount = amount_in_wei(9)
    logging.info(f"create new currency")
    new_currency = test_utilities.create_new_currency(
        amount,
        new_currency_symbol,
        new_currency_symbol,
        18,
        smart_contracts_dir=ctx.smart_contracts_dir,
        bridgebank_address=ctx.bridgebank_address,
        solidity_json_path=ctx.solidity_json_path)

    logging.info(f"create test account to use with new currency {new_currency_symbol}")
    basic_transfer_request.ethereum_address = source_ethereum_address
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=10**17,
        target_rowan_balance=10**18)
    test_amount = 39000
    logging.info(f"transfer some of the new currency {new_currency_symbol} to the test sifchain address")
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.sifchain_symbol = ("c" + new_currency["newtoken_symbol"]).lower()
    request.amount = test_amount
    burn_lock_functions.transfer_ethereum_to_sifchain(request)

    logging.info("send some new currency to ethereum")
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        ctx.smart_contracts_dir, ctx.ethereum_network)
    request.amount = test_amount - 1
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)


@pytest.mark.usefixtures("with_snapshot")
@pytest.mark.parametrize("snapshot_name", ["s1"])
def test_transfer_tokens_with_some_currency(ctx_legacy):
    ctx = ctx_legacy  # TODO Refactoring in progress, pytest fixture name "ctx" is now used for new code
    ctx.set_operator_private_key_env_var()  # TODO Instead of @pytest.mark.usefixtures("operator_private_key")
    basic_transfer_request = ctx.basic_transfer_request
    source_ethereum_address = ctx.source_ethereum_address
    rowan_source_integrationtest_env_credentials = ctx.rowan_source_integrationtest_env_credentials
    rowan_source_integrationtest_env_transfer_request = ctx.rowan_source_integrationtest_env_transfer_request(basic_transfer_request)
    new_currency_symbol = ("a" + get_shell_output("uuidgen").replace("-", ""))[:4]
    do_currency_test(
        ctx,
        new_currency_symbol,
        basic_transfer_request,
        source_ethereum_address,
        rowan_source_integrationtest_env_credentials,
        rowan_source_integrationtest_env_transfer_request)


@pytest.mark.usefixtures("with_snapshot")
@pytest.mark.parametrize("snapshot_name", ["s1"])
def test_three_letter_currency_with_capitals_in_name(ctx_legacy):
    ctx = ctx_legacy  # TODO Refactoring in progress, pytest fixture name "ctx" is now used for new code
    ctx.set_operator_private_key_env_var()  # TODO Instead of @pytest.mark.usefixtures("operator_private_key")
    basic_transfer_request = ctx.basic_transfer_request
    source_ethereum_address = ctx.source_ethereum_address
    rowan_source_integrationtest_env_credentials = ctx.rowan_source_integrationtest_env_credentials
    rowan_source_integrationtest_env_transfer_request = ctx.rowan_source_integrationtest_env_transfer_request(basic_transfer_request)
    new_currency_symbol = ("F" + get_shell_output("uuidgen").replace("-", ""))[:3]
    do_currency_test(
        ctx,
        new_currency_symbol,
        basic_transfer_request,
        source_ethereum_address,
        rowan_source_integrationtest_env_credentials,
        rowan_source_integrationtest_env_transfer_request)
