# to test against ropsten, define:
# ETHEREUM_ADDRESS
# ETHEREUM_PRIVATE_KEY
# ETHEREUM_NETWORK = ropsten
#
import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from pytest_utilities import generate_test_account
import test_utilities
from test_utilities import get_required_env_var, SifchaincliCredentials


def test_fee_charged_to_transfer_rowan_to_erowan(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgetoken_address,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    logging.info(f"credentials: {rowan_source_integrationtest_env_credentials}")
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
    logging.info(f"sending rowan to erowan and checking that a ceth fee was charged")
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = 31500

    # get the starting ceth balance, transfer some rowan to erowan, get the ending ceth
    # balance.  The difference is the fee charged and should be equal to request.ceth_amount

    starting_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, "ceth")
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, "ceth")
    fee = starting_ceth_balance - ending_ceth_balance
    assert fee == request.ceth_amount


def test_do_not_transfer_if_fee_allowed_is_too_low(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgetoken_address,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    target_ceth_balance = 10 ** 18
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=target_ceth_balance,
        target_rowan_balance=target_rowan_balance
    )
    # send some test account ceth back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = 31500

    logging.info("try to transfer rowan to erowan with a ceth_amount that's too low")
    with pytest.raises(Exception):
        request.ceth_amount = test_utilities.lock_gas_cost - 1
        burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, "ceth")
    assert ending_ceth_balance == target_ceth_balance

    logging.info("try with not owning enough ceth to cover the offer")
    with pytest.raises(Exception):
        request.ceth_amount = target_ceth_balance + 1
        burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, "ceth")
    assert ending_ceth_balance == target_ceth_balance