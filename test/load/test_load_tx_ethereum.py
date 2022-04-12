import time
import threading
import copy 

import siftool_path
from siftool import eth, test_utils, sifchain
from siftool.eth import NULL_ADDRESS
from siftool.common import *
from web3.eth import Contract
import web3
# In separate window: test/integration/framework/siftool run-env
# test/integration/framework/venv/bin/python3 test/load/test_load_tx_ethereum.py

fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1
rowan_contract_address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

threads_num = 3
amount_to_send = 10000

def bridge_bank_lock_burn(eth_tx_wrapper: eth.EthereumTxWrapper, bridge_bank_sc: Contract, test_eth_account: str, 
    recipient: str, token_amount: int, nonce: int, token_sc: Contract = None, isLock: bool = True):

    if token_sc is None:
        tx_opts = {"value": token_amount, "nonce": nonce}
        token_addr = NULL_ADDRESS
    else:
        tx_opts = {"value": 0, "nonce": nonce}
        token_addr = token_sc.address 
        eth_tx_wrapper.transact_sync(token_sc.functions.approve, test_eth_account, tx_opts=tx_opts)(bridge_bank_sc.address, token_amount)
        tx_opts = {"value": 0, "nonce": nonce+1}

    if isLock:
        function = bridge_bank_sc.functions.lock
    else:
        function = bridge_bank_sc.functions.burn
    eth_tx_wrapper.transact(function, test_eth_account, tx_opts=tx_opts)(recipient, token_addr, token_amount)


def test_load_burn_rowan(ctx):
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    rowan_cosmos_denom = "rowan"
    w3_url = ctx.w3_conn.provider.endpoint_uri
    test_eth_accounts = []
    threads = []
    conn_list = []
    eth_tx_wrappers = []

    # Create test ethereum accounts and mint token to each account
    for i in range(threads_num):
        test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
        ctx.mint_generic_erc20_token(rowan_contract_address, amount_to_send, test_eth_account)
        test_eth_accounts.append(test_eth_account)

    for i in range(threads_num):
        w3_conn = eth.web3_connect(w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        bridge_bank_abi, _, bridge_bank_address = ctx.abi_provider.get_descriptor("BridgeBank")
        bridge_bank_sc = w3_conn.eth.contract(abi=bridge_bank_abi, address=bridge_bank_address)
        rowan_abi, _, _ = ctx.abi_provider.get_descriptor("BridgeToken")
        rowan_sc = w3_conn.eth.contract(abi=rowan_abi, address=rowan_contract_address)
        test_eth_account = test_eth_accounts[i]
        nonce = w3_conn.eth.get_transaction_count(test_eth_account)
        recipient = test_utils.sif_addr_to_evm_arg(test_sif_account)
        eth_tx_wrapper = eth.EthereumTxWrapper(w3_conn, ctx.eth.is_local_node)
        eth_tx_wrapper.set_private_key(test_eth_account, ctx.eth._get_private_key(test_eth_account))
        eth_tx_wrappers.append(eth_tx_wrapper)
        threads.append(threading.Thread(target=bridge_bank_lock_burn, args=(
            eth_tx_wrapper, bridge_bank_sc, test_eth_account, recipient, amount_to_send, nonce, rowan_sc, False)))

    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()

    # Verify final balance
    expected_change = [[amount_to_send * threads_num, rowan_cosmos_denom]]
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_change)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)

    assert exactly_one(list(balance_diff.keys())) == rowan_cosmos_denom
    assert balance_diff[rowan_cosmos_denom] == amount_to_send * threads_num


def test_load_erc20_to_sifnode(ctx):
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    token_decimals = 18
    token_data: test_utils.ERC20TokenData = ctx.generate_random_erc20_token_data()
    erc20_sc_address = ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, token_decimals).address

    erc20_cosmos_denom = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, erc20_sc_address)

    w3_url = ctx.w3_conn.provider.endpoint_uri
    test_eth_accounts = []
    threads = []
    conn_list = []
    eth_tx_wrappers = []

    # Create test ethereum accounts and mint token to each account
    for i in range(threads_num):
        test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
        ctx.mint_generic_erc20_token(erc20_sc_address, amount_to_send, test_eth_account)
        test_eth_accounts.append(test_eth_account)

    for i in range(threads_num):
        w3_conn = eth.web3_connect(w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        bridge_bank_abi, _, bridge_bank_address = ctx.abi_provider.get_descriptor("BridgeBank")
        bridge_bank_sc = w3_conn.eth.contract(abi=bridge_bank_abi, address=bridge_bank_address)
        erc20_abi, _, _ = ctx.abi_provider.get_descriptor("BridgeToken")
        erc20_sc = w3_conn.eth.contract(abi=erc20_abi, address=erc20_sc_address)
        test_eth_account = test_eth_accounts[i]
        nonce = w3_conn.eth.get_transaction_count(test_eth_account)
        recipient = test_utils.sif_addr_to_evm_arg(test_sif_account)
        eth_tx_wrapper = eth.EthereumTxWrapper(w3_conn, ctx.eth.is_local_node)
        eth_tx_wrapper.set_private_key(test_eth_account, ctx.eth._get_private_key(test_eth_account))
        eth_tx_wrappers.append(eth_tx_wrapper)
        threads.append(threading.Thread(target=bridge_bank_lock_burn, args=(
            eth_tx_wrapper, bridge_bank_sc, test_eth_account, recipient, amount_to_send, nonce, erc20_sc)))

    start_time = time.time()
    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()

    expected_change = [[amount_to_send * threads_num, erc20_cosmos_denom]]
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_change)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)

    assert exactly_one(list(balance_diff.keys())) == erc20_cosmos_denom
    assert balance_diff[erc20_cosmos_denom] == amount_to_send * threads_num

def test_load_multiple_erc20_to_sifnode(ctx):
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    erc20_contract_addresses = []
    erc20_cosmos_denoms = []
    token_decimals = 18
    for i in range(threads_num):
        token_data: test_utils.ERC20TokenData = ctx.generate_random_erc20_token_data()
        erc20_sc = ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, token_decimals)
        erc20_contract_addresses.append(erc20_sc.address)
        erc20_cosmos_denoms.append(sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, erc20_sc.address))

    w3_url = ctx.w3_conn.provider.endpoint_uri
    test_eth_accounts = []
    threads = []
    conn_list = []
    eth_tx_wrappers = []

    # Create test ethereum accounts and mint token to each account
    for i in range(threads_num):
        test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
        ctx.mint_generic_erc20_token(erc20_contract_addresses[i], amount_to_send, test_eth_account)
        test_eth_accounts.append(test_eth_account)

    for i in range(threads_num):
        w3_conn = eth.web3_connect(w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        
        bridge_bank_abi, _, bridge_bank_address = ctx.abi_provider.get_descriptor("BridgeBank")
        bridge_bank_sc = w3_conn.eth.contract(abi=bridge_bank_abi, address=bridge_bank_address)
        erc20_abi, _, _ = ctx.abi_provider.get_descriptor("BridgeToken")
        erc20_sc = w3_conn.eth.contract(abi=erc20_abi, address=erc20_contract_addresses[i])
        test_eth_account = test_eth_accounts[i]
        nonce = w3_conn.eth.get_transaction_count(test_eth_account)
        recipient = test_utils.sif_addr_to_evm_arg(test_sif_account)
        eth_tx_wrapper = eth.EthereumTxWrapper(w3_conn, ctx.eth.is_local_node)
        eth_tx_wrapper.set_private_key(test_eth_account, ctx.eth._get_private_key(test_eth_account))
        eth_tx_wrappers.append(eth_tx_wrapper)
        threads.append(threading.Thread(target=bridge_bank_lock_burn, args=(
            eth_tx_wrapper, bridge_bank_sc, test_eth_account, recipient, amount_to_send, nonce, erc20_sc)))

    start_time = time.time()
    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()

    expected_changes = []
    for i in range(threads_num):
        expected_changes.append([amount_to_send, erc20_cosmos_denoms[i]])

    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_changes)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    for i in range(threads_num):
        assert balance_diff[erc20_cosmos_denoms[i]] == amount_to_send

def test_load_tx_eth(ctx):
    w3_url = ctx.w3_conn.provider.endpoint_uri
    eth_tx_wrappers = []
    threads = []
    test_eth_accounts = []
    conn_list = []

    for i in range(threads_num):
        test_eth_accounts.append(ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth))

    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    assert amount_to_send < fund_amount_eth
    
    for i in range(threads_num):
        w3_conn = eth.web3_connect(w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        bridge_bank_abi, _, bridge_bank_address = ctx.abi_provider.get_descriptor("BridgeBank")
        bridge_bank_sc = w3_conn.eth.contract(abi=bridge_bank_abi, address=bridge_bank_address)
        test_eth_account = test_eth_accounts[i]
        nonce = w3_conn.eth.get_transaction_count(test_eth_account)
        recipient = test_utils.sif_addr_to_evm_arg(test_sif_account)
        eth_tx_wrapper = eth.EthereumTxWrapper(w3_conn, ctx.eth.is_local_node)
        eth_tx_wrapper.set_private_key(test_eth_account, ctx.eth._get_private_key(test_eth_account))
        eth_tx_wrappers.append(eth_tx_wrapper)
        threads.append(threading.Thread(target=bridge_bank_lock_burn, args=(
            eth_tx_wrapper, bridge_bank_sc, test_eth_account, recipient, amount_to_send, nonce, None)))

    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()

    # Verify final balance
    expected_change = [[amount_to_send * threads_num, ctx.ceth_symbol]]
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_change)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)

    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == amount_to_send * threads_num


# Enable running directly, i.e. without pytest
if __name__ == "__main__":
    basic_logging_setup()
    from siftool import test_utils
    ctx = test_utils.get_env_ctx()
    test_load_burn_rowan(ctx)
    test_load_tx_eth(ctx)
    test_load_erc20_to_sifnode(ctx)
    test_load_multiple_erc20_to_sifnode(ctx)
    print("load test from ethereum to sifnode done.")

