import logging

import pytest

import burn_lock_functions
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials

# FEEDFACE is hardcoded in genesis.rake
feedface_token = "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE"


# Only one of lock or burn will be valid depending on how we implement IBC

@pytest.mark.skip(reason="not the way we do IBC yet")
def test_burn_ibc_coins(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        sifchain_fees_int,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    small_amount = 100

    logging.info("the test account needs enough rowan and ceth for one burn and one lock, make sure it has that")
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=test_utilities.burn_gas_cost + test_utilities.lock_gas_cost + small_amount,
        target_rowan_balance=sifchain_fees_int * 2 + small_amount
    )

    logging.info("create an ERC20 token for use by FEEDFACE")
    new_currency = test_utilities.create_new_currency(
        10 ** 20,
        feedface_token,
        feedface_token,
        18,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path
    )

    # send some test account FEEDFACE back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.sifchain_symbol = feedface_token
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.amount = small_amount
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    feedface_ethereum_balance = test_utilities.get_eth_balance(request)
    assert feedface_ethereum_balance == small_amount


def test_lock_ibc_coins(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        sifchain_fees_int,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    small_amount = 100

    logging.info("the test account needs enough rowan and ceth for one burn and one lock, make sure it has that")
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=test_utilities.burn_gas_cost + test_utilities.lock_gas_cost + small_amount,
        target_rowan_balance=sifchain_fees_int * 2 + small_amount
    )

    logging.info(
        "this test of locking an ibc token relies on the smart contracts deploying a destination ERC20 contract")

    # send some test account FEEDFACE back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.sifchain_symbol = feedface_token
    logging.warning("not clear what symbol to use for ethereum_symbol")
    request.amount = small_amount
    burn_lock_functions.send_from_sifchain_to_ethereum(request, credentials)
    feedface_token_data = test_utilities.wait_for_ethereum_token(request, feedface_token)
    request.ethereum_symbol = feedface_token_data["token"]
    feedface_ethereum_balance = test_utilities.get_eth_balance(request)
    assert feedface_ethereum_balance == small_amount
