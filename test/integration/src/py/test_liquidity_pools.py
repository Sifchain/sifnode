import logging
import os
import time
import json
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

#LOCAL UTILITIES
def generate_new_currency(symbol, amount, solidity_json_path, operator_address, ethereum_network):
    logging.info(f"create new currency: "+symbol)
    new_currency = test_utilities.create_new_currency(
        amount,
        symbol,
        token_name=symbol,
        decimals=18,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path,
        operator_address=operator_address,
        ethereum_network=ethereum_network
    )
    return new_currency


def transfer_new_currency(request, symbol, currency, amount):
    logging.info(f"transfer some of the new currency {symbol} to the test sifchain address")
    request.ethereum_symbol = currency["newtoken_address"]
    request.sifchain_symbol = symbol
    request.amount = amount
    burn_lock_functions.transfer_ethereum_to_sifchain(request)


def get_currency_faucet_balance(currency, json_input):
    amount = 0;
    for obj in json_input:
        for key, value in obj.items():
            if value == currency:
                amount = int(obj['amount'])
    return amount


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
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn


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
        "-y -o json"
    ])
    json_str = get_shell_output_json(cmd)
    assert(json_str.get("code", 0) == 0)
    txn = get_transaction_result(json_str["txhash"], transfer_request.sifnodecli_node, transfer_request.chain_id)
    tx = txn["tx"]
    logging.debug(f"resulting tx: {tx}")
    return txn


def test_add_faucet_coins(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        solidity_json_path,
        operator_address,
        ethereum_network
):
    #Generate New Random Currency
    new_currency_symbol = ("a" + get_shell_output("uuidgen").replace("-", ""))[:8]
    target_new_currency_balance = 300000
    new_currency = generate_new_currency(new_currency_symbol, target_new_currency_balance, solidity_json_path, operator_address, ethereum_network)
    sifchain_symbol = ("c" + new_currency["newtoken_symbol"]).lower()

    logging.info(f"request: {basic_transfer_request}")
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=100,
        target_rowan_balance=target_rowan_balance
    )
    #Transfer Newly Generated Currency
    transfer_new_currency(request, sifchain_symbol, new_currency, target_new_currency_balance)

    sifaddress = request.sifchain_address
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    balance_ceth = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
    balance_new_currency = test_utilities.get_sifchain_addr_balance(sifaddress, request.sifnodecli_node, sifchain_symbol)
    assert balance_new_currency == 300000, "new currency balance is wrong in the test account"
    assert balance_ceth == 100, "new ceth currency balance is wrong in the test account"
    assert balance == 10 ** 18, "new rowan currency balance is wrong in the test account"
    logging.info(f"{sifchain_symbol} balance = {balance_new_currency}")
    logging.info(f"ceth  balance = {balance_ceth}")
    logging.info(f"rowan  balance = {balance}")

    # add rowan
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 50000
    basic_transfer_request.sifchain_symbol = "rowan"
    basic_transfer_request.sifchain_address = sifaddress
    txn = add_faucet_coins(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    # add ceth
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10
    basic_transfer_request.sifchain_symbol = "ceth"
    basic_transfer_request.sifchain_address = sifaddress
    txn = add_faucet_coins(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    balance_ceth = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
    # add randome currency
    get_faucet_balance(request.sifnodecli_node)
    request.amount = 200000
    request.sifchain_symbol = sifchain_symbol
    request.sifchain_address = sifaddress
    txn = add_faucet_coins(request, credentials)
    assert(txn.get("code", 0) == 0)
    balance_new_currency = test_utilities.get_sifchain_addr_balance(sifaddress, request.sifnodecli_node, sifchain_symbol)
  
    #get account remaining balances
    logging.info("get ending balances just to have those commands in the history")
    assert balance_new_currency == 100000, "new currency balance is wrong in the test account"
    assert balance_ceth == 90, "new ceth currency balance is wrong in the test account"
    assert balance == 999999999999750000, "new rowan currency balance is wrong in the test account"

    logging.info(f"{sifchain_symbol} balance = {balance_new_currency}")
    logging.info(f"ceth  balance = {balance_ceth}")
    logging.info(f"rowan  balance = {balance}")
    get_faucet_balance(request.sifnodecli_node)

def test_request_faucet_coins(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        solidity_json_path,
        operator_address,
        ethereum_network
):

    #Generate New Random Currency
    new_currency_symbol = ("a" + get_shell_output("uuidgen").replace("-", ""))[:8]
    target_new_currency_balance = 300000
    new_currency = generate_new_currency(new_currency_symbol, target_new_currency_balance, solidity_json_path, operator_address, ethereum_network)
    sifchain_symbol = ("c" + new_currency["newtoken_symbol"]).lower()

    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=100,
        target_rowan_balance=target_rowan_balance
    )
    #Transfer Newly Generated Currency
    transfer_new_currency(request, sifchain_symbol, new_currency, target_new_currency_balance)

    sifaddress = request.sifchain_address
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    balance_ceth = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
    balance_new_currency = test_utilities.get_sifchain_addr_balance(sifaddress, request.sifnodecli_node, sifchain_symbol)
    
    assert balance_new_currency == 300000, "new currency balance is wrong in the test account"
    assert balance_ceth == 100, "new ceth currency balance is wrong in the test account"
    assert balance == 10 ** 18, "new rowan currency balance is wrong in the test account"

    logging.info(f"{sifchain_symbol} balance = {balance_new_currency}")
    logging.info(f"ceth  balance = {balance_ceth}")
    logging.info(f"rowan  balance = {balance}")

    #rowan
    start_faucet_rowan_balance = get_currency_faucet_balance("rowan", get_faucet_balance(basic_transfer_request.sifnodecli_node))
    logging.info(get_faucet_balance(basic_transfer_request.sifnodecli_node))
    basic_transfer_request.amount = 50000
    basic_transfer_request.sifchain_symbol = "rowan"
    basic_transfer_request.sifchain_address = sifaddress
    # add
    txn = add_faucet_coins(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    after_add_faucet_rowan_balance = get_currency_faucet_balance("rowan", get_faucet_balance(basic_transfer_request.sifnodecli_node))
    assert (start_faucet_rowan_balance + 50000 == after_add_faucet_rowan_balance), "faucet balance in rowan should be higher after adding coins to faucet"
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    # request
    basic_transfer_request.amount = 25000
    txn = request_faucet_coins(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    after_request_faucet_rowan_balance = get_currency_faucet_balance("rowan", get_faucet_balance(basic_transfer_request.sifnodecli_node))
    assert (start_faucet_rowan_balance + 25000 == after_request_faucet_rowan_balance), "current faucet balance in rowan should be equal with the initial faucet balance + 25000"
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
     
    #ceth
    start_faucet_ceth_balance = get_currency_faucet_balance("ceth", get_faucet_balance(basic_transfer_request.sifnodecli_node))
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10
    basic_transfer_request.sifchain_symbol = "ceth"
    basic_transfer_request.sifchain_address = sifaddress
    # add
    txn = add_faucet_coins(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    after_add_faucet_ceth_balance = get_currency_faucet_balance("ceth", get_faucet_balance(basic_transfer_request.sifnodecli_node))
    assert (start_faucet_ceth_balance + 10 == after_add_faucet_ceth_balance), "faucet balance in ceth should be higher after adding coins to faucet"
    get_faucet_balance(request.sifnodecli_node)
    # request
    basic_transfer_request.amount = 5
    txn = request_faucet_coins(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance_ceth = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
    
    #random currency
    start_faucet_random_balance = get_currency_faucet_balance(sifchain_symbol, get_faucet_balance(request.sifnodecli_node))
    get_faucet_balance(request.sifnodecli_node)
    request.amount = 20000
    request.sifchain_symbol = sifchain_symbol
    request.sifchain_address = sifaddress
    # add
    txn = add_faucet_coins(request, credentials)
    assert(txn.get("code", 0) == 0)
    after_add_faucet_random_balance = get_currency_faucet_balance(sifchain_symbol, get_faucet_balance(request.sifnodecli_node))
    assert (start_faucet_random_balance + 20000 == after_add_faucet_random_balance), "faucet balance in random currency should be higher after adding coins to faucet"
    get_faucet_balance(request.sifnodecli_node)
    balance_new_currency = test_utilities.get_sifchain_addr_balance(sifaddress, request.sifnodecli_node, sifchain_symbol)
   
    # request
    request.amount = 1000
    txn = request_faucet_coins(request, credentials)
    assert(txn.get("code", 0) == 0)
    get_faucet_balance(request.sifnodecli_node)
    balance_new_currency = test_utilities.get_sifchain_addr_balance(sifaddress, request.sifnodecli_node, sifchain_symbol)
    
    #get faucet balances
    logging.info("get remaining balances just to have those commands in the history")
    assert balance_new_currency == 281000, "new currency balance is wrong in the test account"
    assert balance_ceth == 95, "new ceth currency balance is wrong in the test account"
    assert balance == 999999999999575000, "new rowan currency balance is wrong in the test account"
    logging.info(f"{sifchain_symbol} balance = {balance_new_currency}")
    logging.info(f"ceth  balance = {balance_ceth}")
    logging.info(f"rowan  balance = {balance}")
    logging.info(get_faucet_balance(request.sifnodecli_node))



@pytest.mark.skip(reason="not now")
def test_create_pools(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        sifchain_fees_int
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
        current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees_int
        assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
        assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)

    # check for failure if we try to create a pool twice
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 14)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth") == current_ceth_balance)


# @pytest.mark.skip(reason="not now")
@pytest.mark.usefixtures("operator_private_key")
def test_pools(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        operator_address,
        ethereum_network,
        source_ethereum_address,
        sifchain_fees_int
):
    # max symbol length in clp is 10
    new_currency_symbol = ("a" + get_shell_output("uuidgen").replace("-", ""))[:8]
    target_new_currency_balance = 5 * 10 ** 18
    logging.info(f"create new currency")
    new_currency = test_utilities.create_new_currency(
        1000 * target_new_currency_balance,
        new_currency_symbol,
        token_name=new_currency_symbol,
        decimals=18,
        smart_contracts_dir=smart_contracts_dir,
        bridgebank_address=bridgebank_address,
        solidity_json_path=solidity_json_path,
        operator_address=operator_address,
        ethereum_network=ethereum_network
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
    # wait for balance
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, "rowan", target_rowan_balance, basic_transfer_request.sifnodecli_node)
    test_utilities.wait_for_sifchain_addr_balance(sifaddress, sifchain_symbol, target_new_currency_balance, basic_transfer_request.sifnodecli_node)

    request2, credentials2 = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=target_ceth_balance,
        target_rowan_balance=target_rowan_balance
    )

    logging.info(f"transfer some of the new currency {new_currency_symbol} to the test sifchain address")
    request2.ethereum_symbol = new_currency["newtoken_address"]
    request2.sifchain_symbol = sifchain_symbol
    request2.amount = target_new_currency_balance
    burn_lock_functions.transfer_ethereum_to_sifchain(request2)

    sifaddress2 = request2.sifchain_address
    # wait for balance
    test_utilities.wait_for_sifchain_addr_balance(sifaddress2, "rowan", target_rowan_balance, basic_transfer_request.sifnodecli_node)
    test_utilities.wait_for_sifchain_addr_balance(sifaddress2, sifchain_symbol, target_new_currency_balance, basic_transfer_request.sifnodecli_node)

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
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)

    change_amount = 10 ** 17
    basic_transfer_request.amount = change_amount
    # Fail if amount is less than or equal to minimum
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 7)
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)

    change_amount = 10 ** 18
    basic_transfer_request.amount = change_amount
    # Only works the first time, fails later.
    txn = create_pool(basic_transfer_request, credentials2)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_coin_balance = current_coin_balance - change_amount
    current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to create a pool twice
    txn = create_pool(basic_transfer_request, credentials)
    assert(txn["code"] == 14)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # ensure we can add liquidity, money gets transferred
    txn = add_pool_liquidity(basic_transfer_request, credentials)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_coin_balance = current_coin_balance - change_amount
    current_rowan_balance = current_rowan_balance - change_amount - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # ensure we can remove liquidity, money gets transferred
    txn = remove_pool_liquidity(basic_transfer_request, credentials, 5000)
    assert(txn.get("code", 0) == 0)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_coin_balance = current_coin_balance + change_amount
    current_rowan_balance = current_rowan_balance + change_amount - sifchain_fees_int

    # check for failure if we try to remove more
    txn = remove_pool_liquidity(basic_transfer_request, credentials, 10000)
    assert(txn["code"] == 26)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

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
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to swap too much for user
    basic_transfer_request.amount = 10 ** 19
    txn = swap_pool(basic_transfer_request, "rowan", sifchain_symbol, credentials)
    assert(txn["code"] == 27)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan") == current_rowan_balance)
    assert(test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, sifchain_symbol) == current_coin_balance)

    # check for failure if we try to swap too much for pool
    basic_transfer_request.amount = 5 * 10 ** 17
    txn = swap_pool(basic_transfer_request, "rowan", sifchain_symbol, credentials)
    assert(txn["code"] == 31)
    get_pools(basic_transfer_request.sifnodecli_node)
    current_rowan_balance = current_rowan_balance - sifchain_fees_int
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
