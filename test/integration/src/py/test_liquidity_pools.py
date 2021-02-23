import logging
import os
import time

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account, \
    get_shell_output_json, get_shell_output, detect_errors_in_sifnodecli_output, get_transaction_result, amount_in_wei


smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
bridgetoken_address = get_required_env_var("BRIDGE_TOKEN_ADDRESS")


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
    assert(json_str.get("code", 0) == 0)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn

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
    assert(json_str.get("code", 0) == 0)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn


def remove_pool_liquidity(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials,
        wBasis
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
        f"--wBasis {wBasis}",
        f"--asymmetry 0",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn

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
    assert(json_str.get("code", 0) == 0)
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

    sifaddress = request.sifchain_address
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

    sifaddress = request.sifchain_address
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

@pytest.mark.skip(reason="not now")
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

    sifaddress = request.sifchain_address
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
@pytest.mark.usefixtures("operator_private_key")
def test_pools(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        solidity_json_path
):
    # max symbol length in clp is 10
    new_currency_symbol = ("a" + get_shell_output("uuidgen").replace("-", ""))[:8]
    target_new_currency_balance = 5 * 10 ** 18
    logging.info(f"create new currency")
    new_currency = test_utilities.create_new_currency(
        1000 * target_new_currency_balance,
        new_currency_symbol,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path
    )
    sifchain_symbol = ("c" + new_currency["newtoken_symbol"]).lower()

    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 5 * 10 ** 18
    target_ceth_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=target_ceth_balance,
        target_rowan_balance=target_rowan_balance
    )

    logging.info(f"transfer some of the new currency {new_currency_symbol} to the test sifchain address")
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.sifchain_symbol = sifchain_symbol
    request.amount = target_new_currency_balance
    burn_lock_functions.transfer_ethereum_to_sifchain(request)

    sifaddress = request.sifchain_address
    sifchain_fees = 100000  # Should probably make this a constant
    #from_key = credentials.from_key
    # wait for balance
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, "rowan", target_rowan_balance, basic_transfer_request.sifnodecli_node)
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, sifchain_symbol, target_new_currency_balance, basic_transfer_request.sifnodecli_node)

    pools = get_pools(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.sifchain_symbol = sifchain_symbol
    basic_transfer_request.sifchain_address = sifaddress
    current_coin_balance = target_new_currency_balance
    current_rowan_balance = target_rowan_balance

    change_amount = 10 ** 19
    basic_transfer_request.amount = change_amount
    # Fail if amount is greater than user has
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 12)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)

    change_amount = 10 ** 17
    basic_transfer_request.amount = change_amount
    # Fail if amount is less than or equal to minimum
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 7)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)

    change_amount = 10 ** 18
    basic_transfer_request.amount = change_amount
    # Only works the first time, fails later.
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_coin_balance = current_coin_balance - change_amount
    current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to create a pool twice
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 14)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # ensure we can add liquidity, money gets transferred
    txn = add_pool_liquidity(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_coin_balance = current_coin_balance - change_amount
    current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to remove all liquidity
    txn = remove_pool_liquidity(basic_transfer_request, credentials, 10000)
    assert(txn["code"] == 26)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # ensure we can remove liquidity, money gets transferred
    txn = remove_pool_liquidity(basic_transfer_request, credentials, 5000)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_coin_balance = current_coin_balance + change_amount
    current_rowan_balance = current_rowan_balance + change_amount - sifchain_fees
    #assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    #assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_ceth_balance)
    # no slippage if pool is perfectly balanced.

    # TODO: compute this precisely?
    slip_pct = 0.01
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    slip_cost = (slip_pct * current_rowan_balance)
    assert(balance >= current_rowan_balance - slip_cost and balance <= current_rowan_balance + slip_cost )
    current_rowan_balance = balance
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol)
    slip_cost = (slip_pct * current_coin_balance)
    assert(balance >= current_coin_balance - slip_cost and balance <= current_coin_balance + slip_cost)
    current_coin_balance = balance

    # check for failure if we try to add too much liquidity
    basic_transfer_request.amount = 10 ** 19
    txn = add_pool_liquidity(basic_transfer_request, credentials)
    assert(txn["code"] == 25)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to remove too much liquidity, only occurs now if not 100% liquidity provider
    """
    txn = remove_pool_liquidity(basic_transfer_request, credentials, 10000)
    assert(txn["code"] == 3)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)
    """

    # check for failure if we try to swap too much for user
    txn = swap_pool(basic_transfer_request, "rowan", sifchain_symbol, credentials)
    assert(txn["code"] == 27)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to swap too much for pool
    basic_transfer_request.amount = 5 * 10 ** 17
    txn = swap_pool(basic_transfer_request, "rowan", sifchain_symbol, credentials)
    assert(txn["code"] == 31)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # now try to do a swap that works
    change_amount = 10 ** 15
    basic_transfer_request.amount = change_amount
    txn = swap_pool(basic_transfer_request, "rowan", sifchain_symbol, credentials)
    # TODO: compute this precisely?
    slip_pct = 0.01
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    assert(balance < current_rowan_balance)
    current_rowan_balance = balance
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol)
    assert(balance > current_coin_balance)
    current_coin_balance = balance
