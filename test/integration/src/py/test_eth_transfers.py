import logging

import pytest

import burn_lock_functions
import test_utilities
from pytest_utilities import generate_test_account, generate_minimal_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def test_eth_to_ceth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        ceth_amount,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    logging.info(f"transfer_request: {basic_transfer_request}")
    return generate_minimal_test_account(
        base_transfer_request=basic_transfer_request,
        target_ceth_balance=ceth_amount
    )


def test_eth_to_ceth_and_back_to_eth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        rowan_amount,
        ceth_amount,
        bridgetoken_address,
):
    # logging.info("shark")
    # x = 4.9908339925e+18
    # s = "{:d}".format(int(x))
    # logging.info(f"xiss: {s}")
    # logging.info(f"xiss: {x}")
    # assert False
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    logging.info(f"rq: {basic_transfer_request}")
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=ceth_amount,
        target_rowan_balance=rowan_amount
    )
    # send some test account ceth back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    logging.critical(f"get balance of test account")
    test_utilities.get_sifchain_addr_balance(
        request.sifchain_address,
        sifnodecli_node=request.sifnodecli_node, denom="ceth"
    )
    logging.info("send erowan back to ethereum chain, saving 100k for ceth transfer fees")
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = rowan_amount - 200000
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    test_utilities.get_eth_balance(request)

    logging.critical("send eth back to ethereum chain")
    logging.critical("get ceth balance to decide how much to return")
    request.sifchain_symbol = "ceth"
    request.ethereum_symbol = "eth"
    ceth_balance = test_utilities.get_sifchain_addr_balance(
        request.sifchain_address,
        sifnodecli_node=request.sifnodecli_node,
        denom="ceth"
    )
    request.amount = ceth_balance - request.ceth_amount
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    logging.critical("get final eth balance")
    test_utilities.get_eth_balance(request)
    test_utilities.get_sifchain_addr_balance(
        request.sifchain_address,
        sifnodecli_node=request.sifnodecli_node, denom="ceth"
    )


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
