import siftool_path
from siftool import eth, test_utils, sifchain
from siftool.common import *

fund_amount_eth = eth.ETH
rowan_unit = 10 ** 18
fund_amount_sif = 2 * rowan_unit
ibc_token_symbol = 'ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE'
double_peggy_symbol = 'sifBridge00030x1111111111111111111111111111111111111111'
fund_amount = 10000

def bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, cosmos_denom: str, isFirstTime: bool):
    # get cosmos denom's bridge token contract
    cosmos_denom_bridge_token_address = ctx.get_destination_contract_address(cosmos_denom)
    if isFirstTime:
        assert cosmos_denom_bridge_token_address == eth.NULL_ADDRESS
    else:
        assert cosmos_denom_bridge_token_address != eth.NULL_ADDRESS
        cosmos_denom_sc = ctx.get_generic_erc20_sc(cosmos_denom_bridge_token_address)

    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    if cosmos_denom == "rowan":
        test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"], [fund_amount_eth, ctx.ceth_symbol]])
    else:
        test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"], [fund_amount_eth, ctx.ceth_symbol], [fund_amount, cosmos_denom]])

    # init balance for erc20 should be 0
    if isFirstTime:
        test_eth_account_initial_balance = 0
    else:
        test_eth_account_initial_balance = ctx.get_erc20_token_balance(cosmos_denom_sc.address, test_eth_account)

    # send bridge token to ethereum
    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, fund_amount, cosmos_denom)
    ctx.advance_blocks()

    # wait the bridge token created
    if isFirstTime:
        cosmos_denom_bridge_token_address = ctx.wait_for_new_bridge_token_created(cosmos_denom)
        cosmos_denom_sc = ctx.get_generic_erc20_sc(cosmos_denom_bridge_token_address)

    # wait the ethereum reciever's balance change
    eth_account_final_balance = ctx.wait_for_eth_balance_change(test_eth_account, test_eth_account_initial_balance, token_addr=cosmos_denom_bridge_token_address, timeout=90)

    # check the bridge token balance as expected
    assert eth_account_final_balance == test_eth_account_initial_balance + fund_amount

    test_sif_account_before_receive = ctx.get_sifchain_balance(test_sif_account)

    ctx.send_from_ethereum_to_sifchain(test_eth_account, test_sif_account, fund_amount, token_sc=cosmos_denom_sc, isLock=False)
    ctx.advance_blocks()

    test_sif_account_after_receive = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_before_receive, [[fund_amount, cosmos_denom]])
    # check the bridge token back to sifnode side
    assert sifchain.balance_delta(test_sif_account_before_receive, test_sif_account_after_receive)[cosmos_denom] == fund_amount


def test_rowan_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, "rowan", isFirstTime=True)
    bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, "rowan", isFirstTime=False)


def test_ibc_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, ibc_token_symbol, isFirstTime=True)
    bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, ibc_token_symbol, isFirstTime=False)


def test_double_peg_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, double_peggy_symbol, isFirstTime=True)
    bridge_token_to_eth_and_back_to_sifnode_transfer_valid(ctx, double_peggy_symbol, isFirstTime=False)
