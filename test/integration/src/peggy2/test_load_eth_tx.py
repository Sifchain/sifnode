import threading
import time
import copy

import pytest
import siftool_path
from siftool.eth import NULL_ADDRESS
from siftool import eth, test_utils, sifchain
from siftool.common import *
from siftool.test_utils import EnvCtx
from typing import Iterable

fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1  # TODO How much rowan do we need? (this is 10**18)
fund_amount_ceth_cross_chain_fee = 10 * eth.GWEI
double_peggy_symbol = 'sifBridge99990x0000000000000000000000000000000000000000'

def bridge_bank_lock_eth(eth, bridge_bank, test_eth_account, recipient, amount_to_send, nonce):

    token_addr = NULL_ADDRESS  # For "eth", otherwise use coin's address
    # Mandatory tx_opts: {"from": from_eth_acct, "gas": max_gas_required, "value": amount}
    # If "value" is missing, we get "call to non-contract"
    tx_opts = {"value": amount_to_send, "nonce": nonce}

    txhash = eth.transact(bridge_bank.functions.lock, test_eth_account, tx_opts=tx_opts)(recipient, token_addr, amount_to_send)
    # eth.wait_for_transaction_receipt(txhash)

def test_eth_to_ceth_and_back_to_eth_transfer_valid(ctx):
    threads_num = 3
    # ctx.w3_url = "ws://localhost:8545"
    w3_url = ctx.w3_conn.provider.endpoint_uri
    test_eth_accounts = []
    # Create/retrieve a test ethereum account
    for i in range(threads_num):
        test_eth_accounts.append(ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth))

    bridge_bank = ctx.get_bridge_bank_sc()

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])

    # Verify initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    print("++++++ final balance is ", test_sif_account_initial_balance)

    nonce = 0
    conn_list = []
    eth_tx_wrappers = []

    # Send from ethereum to sifchain by locking
    amount_to_send = 123456 * eth.GWEI
    assert amount_to_send < fund_amount_eth
    threads = []
    for i in range(threads_num):
        w3_conn = eth.web3_connect(w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        bridge_bank_abi, _, bridge_bank_address = ctx.abi_provider.get_descriptor("BridgeBank")
        bridge_bank = w3_conn.eth.contract(abi=bridge_bank_abi, address=bridge_bank_address)
        test_eth_account = test_eth_accounts[i]
        nonce = w3_conn.eth.get_transaction_count(test_eth_account)
        recipient = test_utils.sif_addr_to_evm_arg(test_sif_account)
        eth_tx_wrapper = eth.EthereumTxWrapper(w3_conn, ctx.eth.is_local_node)
        eth_tx_wrapper.set_private_key(test_eth_account, ctx.eth._get_private_key(test_eth_account))
        eth_tx_wrappers.append(eth_tx_wrapper)
        threads.append(threading.Thread(target=bridge_bank_lock_eth, args=(
            eth_tx_wrapper, bridge_bank, test_eth_account, recipient, amount_to_send, nonce)))

    start_time = time.time()
    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()
    time.sleep(90)
    # Verify final balance
    # ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)

    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    print("++++++ final balance is ", test_sif_account_final_balance)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == amount_to_send * threads_num
