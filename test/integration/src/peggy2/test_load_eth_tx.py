import threading
import time

import pytest
import siftool_path

from siftool import eth, test_utils, sifchain
from siftool.common import *
from siftool.test_utils import EnvCtx
from typing import Iterable

# TODO for PR: Remove all print outs

fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1  # TODO How much rowan do we need? (this is 10**18)
fund_amount_ceth_cross_chain_fee = 10 * eth.GWEI
double_peggy_symbol = 'sifBridge99990x0000000000000000000000000000000000000000'

def bridge_bank_lock_eth(ctx, test_eth_account, test_sif_account, amount_to_send, nonce):
    ctx.bridge_bank_lock_eth(test_eth_account, test_sif_account, amount_to_send, nonce)

def test_eth_to_ceth_and_back_to_eth_transfer_valid(ctx):
    threads_num = 2
    ctx.w3_url = "ws://localhost:8545"
    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
    nonce = 0
    conn_list = []
    eth_list = []
    for i in range(threads_num):
        w3_conn = eth.web3_connect(ctx.w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        eth_list.append(eth.EthereumTxWrapper(w3_conn, True))

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])

    # Verify initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    print("++++++ final balance is ", test_sif_account_initial_balance)

    # Send from ethereum to sifchain by locking
    amount_to_send = 123456 * eth.GWEI
    assert amount_to_send < fund_amount_eth
    threads = []
    for i in range(threads_num):
        ctx.w3_conn = conn_list[i]
        ctx.eth = eth_list[i]
        threads.append(threading.Thread(target=bridge_bank_lock_eth, args=(ctx, test_eth_account, test_sif_account, amount_to_send, nonce)))
        nonce = nonce + 1

    start_time = time.time()
    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()

    # Verify final balance
    time.sleep(90)
    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
        # ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    print("++++++ final balance is ", test_sif_account_final_balance)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == amount_to_send
