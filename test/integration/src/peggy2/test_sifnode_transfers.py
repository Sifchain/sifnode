import time

from siftool import eth, test_utils, sifchain
from siftool.common import *
import web3

fund_amount_eth = eth.ETH
rowan_unit = test_utils.sifnode_funds_for_transfer_peggy1
fund_amount_sif = 2 * rowan_unit
rowan_contract_address = web3.Web3.toChecksumAddress('0x5fbdb2315678afecb367f032d93f642f64180aa3')

def test_rowan_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    # get rowan contract
    rowan_sc = ctx.get_generic_erc20_sc(rowan_contract_address)

    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"], [fund_amount_eth, ctx.ceth_symbol]])

    test_sif_dummy_account = ctx.create_sifchain_addr()

    # init balance for erc20 rowan
    test_eth_account_initial_balance = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("test_eth_account_initial_balance", test_eth_account_initial_balance)

    # we need mint some rowan in contract then lock into bridge
    mint_amount = 5 * rowan_unit
    ctx.mint_generic_erc20_token(rowan_contract_address, mint_amount, test_eth_account)
    ctx.advance_blocks()
    ctx.wait_for_eth_balance_change(test_eth_account, test_eth_account_initial_balance, token_addr=rowan_contract_address, timeout=90)
    eth_account_balance_after_mint = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("rowan balance after mint", eth_account_balance_after_mint)

    # send rowan to dummy account
    ctx.send_from_ethereum_to_sifchain(test_eth_account, test_sif_dummy_account, mint_amount, token_sc=rowan_sc, isLock=False)
    ctx.advance_blocks()
    ctx.wait_for_eth_balance_change(test_eth_account, eth_account_balance_after_mint, token_addr=rowan_contract_address, timeout=90)
    eth_account_balance_after_lock = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("rowan balance after lock", eth_account_balance_after_lock)

    # sif account initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)
    print("test_sif_account_initial_balance", test_eth_account_initial_balance)

    # lock rowan in sifnode
    amount_to_lock = 1 * rowan_unit
    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, amount_to_lock, "rowan")
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, [[amount_to_lock, "rowan"]])

    # wait rowan change happened in sifnode account
    test_sif_account_after_lock_balance = ctx.get_sifchain_balance(test_sif_account)
    print("after lock balance is ", test_sif_account_after_lock_balance)

    # we need take the transaction fee into account
    rowan_balance = test_sif_account_initial_balance["rowan"] - amount_to_lock
    assert rowan_balance >= test_sif_account_after_lock_balance["rowan"]

    # wait the ethereum reciever's rowan balance change
    ctx.wait_for_eth_balance_change(test_eth_account, eth_account_balance_after_lock, token_addr=rowan_contract_address, timeout=90)
    eth_account_final_balance = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("eth_account_final_balance is ", eth_account_final_balance)

    # check the rowan balance as expected
    assert eth_account_final_balance == eth_account_balance_after_lock + amount_to_lock
