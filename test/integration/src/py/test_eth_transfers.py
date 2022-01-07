import logging

import pytest

import burn_lock_functions
import test_utilities
from pytest_utilities import generate_test_account, generate_minimal_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def disabled_test_eth_to_ceth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    logging.info(f"transfer_request: {basic_transfer_request}")
    return generate_minimal_test_account(
        base_transfer_request=basic_transfer_request,
        target_ceth_balance=100
    )


def disabled_test_eth_to_ceth_and_back_to_eth(
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
    request.sifchain_symbol = "evmnative"
    request.ethereum_symbol = "eth"
    request.amount = small_amount
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)


def disabled_test_transfer_eth_to_ceth_over_limit(
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


# Group 0
# Sending eth from ethereum to sifchain

from integration_framework import main, common, eth, test_utils, inflate_tokens
import eth
import test_utils
from common import *


fund_amount_eth = 10 * eth.ETH
fund_amount_sif = [1 * 10 ** 18, "rowan"]  # TODO How much do we need?


def disabled_test_eth_to_ceth_valid():
    # create/retrieve a test ethereum account
    # create/retrieve a test sifchain account

    # Verify initial balance

    # Send from ethereum to sifchain by locking

    # Verify final balance

    pass


def test_eth_to_ceth_and_back_to_eth_transfer_valid():
    with test_utils.get_test_env_ctx() as ctx:
        _test_eth_to_ceth_and_back_to_eth_transfer_valid(ctx)

def _test_eth_to_ceth_and_back_to_eth_transfer_valid(ctx):
    # create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[fund_amount_sif])

    # Verify initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    # Send from ethereum to sifchain by locking
    amount_to_send = 3 * eth.ETH
    ctx.bridge_bank_lock_eth(test_eth_account, test_sif_account, amount_to_send)
    ctx.advance_blocks()

    # Verify final balance
    test_sif_account_final_balance = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    balance_diff = ctx.sif_balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == amount_to_send

    # Send from sifchain to ethereum by burning on sifchain side,
    # > sifnoded tx ethbridge burn
    # Reduce amount for cross-chain fee. The same formula is used inside this function.
    eth_balance_before = ctx.eth.get_eth_balance(test_eth_account)
    amount_to_send = amount_to_send - ctx.cross_chain_fee_base * ctx.cross_chain_burn_fee
    ctx.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, amount_to_send, ctx.ceth_symbol)

    # Verify final balance
    ctx.wait_for_eth_balance_change(test_eth_account, eth_balance_before)


def disabled_test_erc20_to_sifnode_and_back():
    with test_utils.get_test_env_ctx() as ctx:
        _test_erc20_to_sifnode_and_back(ctx)

def _test_erc20_to_sifnode_and_back(ctx):
    token_data = ctx.generate_random_erc20_token_data()
    token_sc = ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, token_data.decimals)

    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    fund_amount_sif = [1 * 10**18, "rowan"]  # TODO How much do we need?
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[fund_amount_sif])

    # TODO

    return
