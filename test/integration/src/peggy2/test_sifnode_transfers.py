import time

from integration_framework import main, common, eth, test_utils, inflate_tokens, sifchain
import eth
import test_utils
import sifchain
from common import *
import web3

fund_amount_eth = 10 * eth.ETH
rowan_unit = test_utils.sifnode_funds_for_transfer_peggy1
fund_amount_sif = 10 * rowan_unit
rowan_contract_address = web3.Web3.toChecksumAddress('0x5fbdb2315678afecb367f032d93f642f64180aa3')
mint_rowan_amount = 10 ** 28

def test_rowan_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])

    # mint and lock rowan in ethereum, make sure enough rowan in bridge contract for unlock
    rowan_contract = ctx.get_generic_erc20_sc(rowan_contract_address)
    print("", rowan_contract.functions)

    # init balance for erc20 rowan
    test_eth_account_initial_balance = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)

    # ctx.mint_generic_erc20_token(rowan_contract_address, mint_rowan_amount, test_eth_account)
    # ctx.approve_erc20_token(rowan_contract, test_eth_account, mint_rowan_amount)
    # ctx.bridge_bank_lock_erc20(rowan_contract, )

    # sif account initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    print("", test_eth_account_initial_balance)
    # exit()
    # Send from ethereum to sifchain by locking
    amount_to_send = 1 * eth.ETH
    assert amount_to_send < fund_amount_eth

    ctx.bridge_bank_lock_eth(test_eth_account, test_sif_account, amount_to_send)
    ctx.advance_blocks()

    # Verify final balance
    test_sif_account_final_balance = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == amount_to_send

    # Verify initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)
    print("+++++++++++++ init balance is ", test_sif_account_initial_balance)

    # Send from ethereum to sifchain by locking
    amount_to_lock = 1 * rowan_unit
    print("+++++++++++++ amount_to_lock is ", amount_to_lock)

    # test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)
    ctx.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, amount_to_lock, "rowan",)
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, [[amount_to_lock, "rowan"]])

    test_sif_account_after_lock_balance = ctx.get_sifchain_balance(test_sif_account)
    print("+++++++++++++ after lock balance is ", test_sif_account_after_lock_balance)

    # we need take the transaction fee into account
    rowan_balance = test_sif_account_initial_balance["rowan"] - amount_to_lock
    assert rowan_balance >= test_sif_account_after_lock_balance["rowan"]

    # Verify final balance
    ctx.wait_for_eth_balance_change(rowan_contract_address, test_eth_account_initial_balance, amount_to_lock)
    test_eth_account_balance_after_lock = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("+++++++++++++ after lock eth balance is ", test_eth_account_balance_after_lock)
    assert test_eth_account_balance_after_lock - amount_to_lock == test_eth_account_initial_balance
    # test_utilities.wait_for_sifchain_addr_balance(test_sif_account, "rowan", amount_to_lock,
    #                                               basic_transfer_request.sifnoded_node, 180)

    # test_sif_account_final_balance = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    # balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    # assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    # assert balance_diff[ctx.ceth_symbol] == amount_to_send
    #
    # # Send from sifchain to ethereum by burning on sifchain side,
    # # > sifnoded tx ethbridge burn
    # # Reduce amount for cross-chain fee. The same formula is used inside this function.
    # eth_balance_before = ctx.eth.get_eth_balance(test_eth_account)
    # amount_to_send = amount_to_send - ctx.cross_chain_fee_base * ctx.cross_chain_burn_fee
    # ctx.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, amount_to_send, ctx.ceth_symbol)
    #
    # # Verify final balance
    # ctx.wait_for_eth_balance_change(test_eth_account, eth_balance_before)


    