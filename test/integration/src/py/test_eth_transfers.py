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
        target_ceth_balance=10 ** 15
    )


def test_eth_to_ceth_and_back_to_eth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=10 ** 18,
        target_rowan_balance=10 ** 18
    )
    # send some test account ceth back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.amount = int(request.amount / 2)
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
