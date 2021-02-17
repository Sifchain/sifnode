import logging
import os
import time

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account, \
    get_shell_output_json, get_shell_output, detect_errors_in_sifnodecli_output, get_transaction_result


def get_faucet_balance(sifnodecli_node):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli q faucet balance {node} -o json"
    result = get_shell_output_json(command_line)
    return result


def get_pools(sifnodecli_node):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli q clp pools {node} -o json"
    # returns error when empty
    try:
        json_str = get_shell_output_json(command_line)
        return json_str
    except Exception as e:
        logging.debug(f"get_pools is empty.")


# sifnodecli tx clp create-pool --from sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --symbol ceth --nativeAmount 100000 --externalAmount 1000000
def create_pool(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"create_pool")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp create-pool",
        f"--from {transfer_request.sifchain_address}",
        f"--symbol {transfer_request.sifchain_symbol}",
        f"--nativeAmount {transfer_request.amount}",
        f"--externalAmount {transfer_request.amount}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn

# sifnodecli tx clp swap --from sif --sentSymbol ceth --receivedSymbol cdash --sentAmount 20
def swap_pool(
        transfer_request: EthereumToSifchainTransferRequest,
        sent_symbol, received_symbol,
        credentials: SifchaincliCredentials
):
    logging.debug(f"swap_pool")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp swap",
        f"--from {transfer_request.sifchain_address}",
        f"--sentSymbol {sent_symbol}",
        f"--receivedSymbol {received_symbol}",
        f"--sentAmount {transfer_request.amount}",
        f"--minReceivingAmount {int(transfer_request.amount * 0.99)}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn


# sifnodecli tx clp remove-liquidity --from sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5 --symbol cacoin --wBasis 5001 --asymmetry -1 --yes --node tcp://54.218.170.168:26657 --trust-node --chain-id sandpit --gas-prices "0.5rowan"
def remove_pool_liquidity(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"remove_pool_liquidity")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp remove-liquidity",
        f"--from {transfer_request.sifchain_address}",
        f"--symbol {transfer_request.sifchain_symbol}",
        f"--wBasis 10000",
        f"--asymmetry 0",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn

# sifnodecli tx clp add-liquidity --from sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5 --symbol cacoin --nativeAmount 10000000000000000000000 --externalAmount 1000000000000000000000000 --node tcp://54.218.170.168:26657 --trust-node --chain-id sandpit --gas-prices "0.5rowan"
def add_pool_liquidity(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"add_pool_liquidity")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp add-liquidity",
        f"--from {transfer_request.sifchain_address}",
        f"--symbol {transfer_request.sifchain_symbol}",
        f"--nativeAmount {transfer_request.amount}",
        f"--externalAmount {transfer_request.amount}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn


def add_faucet_coins(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"add_faucet_coins {transfer_request}")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx faucet add-coins",
        f"{transfer_request.amount}{transfer_request.sifchain_symbol}",
        f"--from {transfer_request.sifchain_address}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y"
    ])
    result = get_shell_output(cmd)
    detect_errors_in_sifnodecli_output(result)
    return result


def request_faucet_coins(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"request_faucet_coins {transfer_request}")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx faucet request-coins",
        f"{transfer_request.amount}{transfer_request.sifchain_symbol}",
        f"--from {transfer_request.sifchain_address}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y"
    ])
    result = get_shell_output(cmd)
    detect_errors_in_sifnodecli_output(result)
    return result


@pytest.mark.skip(reason="only test removal")
def test_add_faucet_coins(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=0,
        target_rowan_balance=target_rowan_balance
    )
    time.sleep(10)

    # sifaddress = rowan_source_integrationtest_env_transfer_request.sifchain_address
    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")

    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10 ** 17
    basic_transfer_request.sifchain_symbol = "rowan"
    basic_transfer_request.sifchain_address = sifaddress
    add_faucet_coins(basic_transfer_request, credentials)
    time.sleep(10)
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")


@pytest.mark.skip(reason="not working yet")
def test_request_faucet_coins(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=0,
        target_rowan_balance=target_rowan_balance
    )
    time.sleep(10)

    # sifaddress = rowan_source_integrationtest_env_transfer_request.sifchain_address
    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")

    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10 ** 17
    basic_transfer_request.sifchain_symbol = "rowan"
    basic_transfer_request.sifchain_address = sifaddress
    request_faucet_coins(basic_transfer_request, credentials)
    time.sleep(10)
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")

#@pytest.mark.skip(reason="not now")
def test_create_pools(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 19
    target_ceth_balance = 10 ** 19
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=target_ceth_balance,
        target_rowan_balance=target_rowan_balance
    )
    #time.sleep(10)

    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    # wait for balance
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, "rowan", target_rowan_balance, basic_transfer_request.sifnodecli_node)
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, "ceth", target_ceth_balance, basic_transfer_request.sifnodecli_node)

    pools = get_pools(basic_transfer_request.sifnodecli_node)
    change_amount = 10 ** 18
    sifchain_fees = 100000  # Should probably make this a constant
    basic_transfer_request.amount = change_amount
    basic_transfer_request.sifchain_symbol = "ceth"
    basic_transfer_request.sifchain_address = sifaddress
    current_ceth_balance = target_ceth_balance
    current_rowan_balance = target_rowan_balance

    # Only works the first time, fails later.  Make this flexible for manual and private net testing for now.
    if pools is None:
        create_pool(basic_transfer_request, credentials)
        get_pools(basic_transfer_request.sifnodecli_node)
        current_ceth_balance = current_ceth_balance - change_amount
        current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees
        assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
        assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # check for failure if we try to create a pool twice
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 14)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)


# @pytest.mark.skip(reason="not now")
def test_pools(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 18
    target_ceth_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=target_ceth_balance,
        target_rowan_balance=target_rowan_balance
    )
    #time.sleep(10)

    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    # wait for balance
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, "rowan", target_rowan_balance, basic_transfer_request.sifnodecli_node)
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, "ceth", target_ceth_balance, basic_transfer_request.sifnodecli_node)

    pools = get_pools(basic_transfer_request.sifnodecli_node)
    change_amount = 10 ** 17
    sifchain_fees = 100000  # Should probably make this a constant
    basic_transfer_request.amount = change_amount
    basic_transfer_request.sifchain_symbol = "ceth"
    basic_transfer_request.sifchain_address = sifaddress
    current_ceth_balance = target_ceth_balance
    current_rowan_balance = target_rowan_balance

    # ensure we can add liquidity, money gets transferred
    txn = add_pool_liquidity(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_ceth_balance = current_ceth_balance - change_amount
    current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # ensure we can remove liquidity, money gets transferred
    txn = remove_pool_liquidity(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_ceth_balance = current_ceth_balance + change_amount
    current_rowan_balance = current_rowan_balance + change_amount - sifchain_fees
    #assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    #assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)
    # no slippage if pool is perfectly balanced.

    #"""
    # TODO: compute this precisely?
    slip_pct = 0.01
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    slip_cost = (slip_pct * current_rowan_balance)
    assert(balance >= current_rowan_balance - slip_cost and balance <= current_rowan_balance + slip_cost )
    current_rowan_balance = balance
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
    slip_cost = (slip_pct * current_ceth_balance)
    assert(balance >= current_ceth_balance - slip_cost and balance <= current_ceth_balance + slip_cost)
    current_ceth_balance = balance
    #"""

    # check for failure if we try to add too much liquidity
    basic_transfer_request.amount = 10 ** 19
    txn = add_pool_liquidity(basic_transfer_request, credentials)
    assert(txn["code"] == 25)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # check for failure if we try to remove too much liquidity
    txn = remove_pool_liquidity(basic_transfer_request, credentials)
    assert(txn["code"] == 3)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # check for failure if we try to swap too much for user
    txn = swap_pool(basic_transfer_request, "rowan", "ceth", credentials)
    assert(txn["code"] == 27)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # check for failure if we try to swap too much for pool
    basic_transfer_request.amount = 5 * 10 ** 17
    txn = swap_pool(basic_transfer_request, "rowan", "ceth", credentials)
    assert(txn["code"] == 31)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # now try to do a swap that works
    change_amount = 10 ** 15
    basic_transfer_request.amount = change_amount
    txn = swap_pool(basic_transfer_request, "rowan", "ceth", credentials)
    # TODO: compute this precisely?
    slip_pct = 0.01
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    assert(balance < current_rowan_balance)
    current_rowan_balance = balance
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
    assert(balance > current_ceth_balance)
    current_ceth_balance = balance

    #current_ceth_balance = current_ceth_balance + change_amount
    #current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees
    #assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    #assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)
