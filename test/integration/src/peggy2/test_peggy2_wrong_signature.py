import siftool_path
from siftool import eth, test_utils, sifchain
from siftool.common import *


fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1  # TODO How much rowan do we need? (this is 10**18)


def test_sign_prophecy_with_wrong_signature_grpc(ctx):
    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    moniker = "temp-moniker"
    test_sif_account = ctx.create_sifchain_addr(moniker=moniker, fund_amounts=[[fund_amount_sif, "rowan"]])
    val_address = ctx.sifnode.get_val_address(moniker)

    # create other one for wrong cosmos sender
    moniker = "temp-moniker-2"
    val_address_2 = ctx.sifnode.get_val_address(moniker)

    # parameter for sign prophecy tx
    prophecy_id = "1"
    signature_for_sign_prophecy = "1"
    result = ctx.sifnode_client.send_sign_prophecy_with_wrong_signature_grpc(
        test_sif_account, val_address, val_address_2, test_eth_account, prophecy_id, signature_for_sign_prophecy)

    # Verify failed tx
    assert result == False
