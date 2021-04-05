import logging

import pytest

import burn_lock_functions
import test_utilities
from pytest_utilities import generate_test_account, generate_minimal_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def test_eth_to_ceth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    logging.info(f"transfer_request: {basic_transfer_request}")
    return generate_minimal_test_account(
        base_transfer_request=basic_transfer_request,
        target_ceth_balance=100
    )


def test_eth_to_ceth_and_back_to_eth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgetoken_address,
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
    # send some test account ceth back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = small_amount
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    test_utilities.get_eth_balance(request)

    logging.info("send eth back to ethereum chain")
    request.sifchain_symbol = "ceth"
    request.ethereum_symbol = "eth"
    request.amount = small_amount
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)


def test_transfer_eth_to_ceth_over_limit(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
):
    basic_transfer_request.ethereum_symbol = "eth"
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.amount = 35 * 10 ** 18
    with pytest.raises(Exception):
        basic_transfer_request.ethereum_address = source_ethereum_address
        generate_test_account(
            basic_transfer_request,
            rowan_source_integrationtest_env_transfer_request,
            rowan_source_integrationtest_env_credentials,
            target_ceth_balance=50 * 10 ** 19,
        )
