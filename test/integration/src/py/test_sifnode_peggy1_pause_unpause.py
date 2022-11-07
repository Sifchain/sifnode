import pytest

import siftool_path
from siftool import eth, test_utils, sifchain
from siftool.inflate_tokens import InflateTokens
from siftool.common import *
from siftool.test_utils import EnvCtx

def test_pause_unpause_no_error(ctx: EnvCtx):
    res = ctx.sifnode.pause_peggy_bridge(ctx.sifchain_ethbridge_admin_account)
    assert res[0]['code'] == 0
    res = ctx.sifnode.unpause_peggy_bridge(ctx.sifchain_ethbridge_admin_account)
    assert res[0]['code'] == 0

# We assert a tx is successful before pausing because we test the pause
# functionality by 1. An error response and 2. Balance unchanged within timeout.
# We want to make sure #2 is not a false positive due to lock function not
# working in the first place
def test_pause_lock_valid(ctx: EnvCtx):
    # Test a working flow:
    fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1
    fund_amount_eth = 10 * eth.ETH

    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    eth_balance_before = ctx.eth.get_eth_balance(test_eth_account)
    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)

    print(sif_balance_before)

    send_amount = 10000
    # Submit lock
    # TODO: Fix Hardcoded denom "rowan"
    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, send_amount, "rowan")

    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)

    # Assert tx go through, balance updated correctly.
    balance_diff = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == send_amount

    # Pause the bridge
    print("Using admin account to pause bridge:", ctx.sifchain_ethbridge_admin_account)
    res = ctx.sifnode.pause_peggy_bridge(ctx.sifchain_ethbridge_admin_account)
    print(res)

    # Submit lock
    # eth_balance_before = ctx.eth.get_eth_balance(test_eth_account)
    # sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
    # res = ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, send_amount, ctx.ceth_symbol)
    # # TODO: Assert on RES getting ERROR

    # balance_change_exception = None
    # try:
    #     sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)
    # except Exception as e:
    #     balance_change_exception = e

    # # TODO: Add more precise assertion, e.g. exception type
    # assert balance_change_exception is not None

    # # Unpause the bridge
    # # TODO: Implement this method
    print("Using admin account to unpause bridge:", ctx.sifchain_ethbridge_admin_account)
    res = ctx.sifnode.unpause_peggy_bridge(ctx.sifchain_ethbridge_admin_account)
    print(res)
    # # Submit lock
    # # Assert tx go through, balance updated correctly.

def test_pause_burn_valid(ctx):



    res = ctx.sifnode.pause_peggy_bridge(ctx.sifchain_ethbridge_admin_account)



    pass

def test_non_admin_cant_pause_bridge(ctx: EnvCtx):
    non_admin_test_sif_acct = ctx.create_sifchain_addr()
    res = ctx.sifnode.pause_peggy_bridge(non_admin_test_sif_acct)
    assert res[0]['code'] != 0
    # Assert res gets error,
    # Assert error code is what's expected