import time

from siftool import eth, test_utils, sifchain
from siftool.common import *
import web3
from json_file_writer import write_registry_json, file_path

fund_amount_eth = eth.ETH
rowan_unit = test_utils.sifnode_funds_for_transfer_peggy1
fund_amount_sif = 2 * rowan_unit
rowan_contract_address = web3.Web3.toChecksumAddress('0x5fbdb2315678afecb367f032d93f642f64180aa3')
ibc_token_symbol = 'ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE'

fund_amount_ibc = 10000
admin_account = "sifnodeadmin"

def test_ibc_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    # events = ctx.smart_contract_get_past_events(ctx.get_cosmos_bridge_sc(), "LogProphecyCompleted")
    # print("all events, ", events)

    # events = ctx.smart_contract_get_past_events(ctx.get_bridge_bank_sc(), "LogBridgeTokenMint")
    # print("all events, ", events)

    # exit(0)

    # deploy erc20 for ibc token
    token_data: test_utils.ERC20TokenData = ctx.generate_random_erc20_token_data()
    ibc_token_sc = ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, 18, cosmosDenom=ibc_token_symbol)
    print("ibc token address", ibc_token_sc.address)

    # add bridge token into whitelist
    ctx.bridge_bank_add_existing_bridge_token(ibc_token_sc.address)

    ctx.tx_grant_minter_role(ibc_token_sc, ctx.get_cosmos_bridge_sc().address)
    ctx.tx_grant_minter_role(ibc_token_sc, ctx.get_bridge_bank_sc().address)
    # exit(0)

    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(
        fund_amounts=[[fund_amount_sif, "rowan"], [fund_amount_eth, ctx.ceth_symbol], [fund_amount_ibc, ibc_token_symbol]])

    # write the ibc token entry to a json file
    write_registry_json(ibc_token_symbol, ibc_token_sc.address)

    # register the ibc token
    ctx.sifnode.peggy2_token_registry_register_all(file_path, [0.5, "rowan"], 1.5, admin_account, ctx.sifnode_chain_id)

    test_eth_account_init_balance = ctx.get_erc20_token_balance(ibc_token_sc.address, test_eth_account)
    print("test_eth_account_init_balance", test_eth_account_init_balance)

    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, fund_amount_ibc, ibc_token_symbol)
    ctx.wait_for_eth_balance_change(test_eth_account, test_eth_account_init_balance, token_addr=ibc_token_sc.address, timeout=90)
    # time.sleep(120)

    test_eth_account_after_lock = ctx.get_erc20_token_balance(ibc_token_sc.address, test_eth_account)
    print("test_eth_account_after_lock", test_eth_account_after_lock)
    assert test_eth_account_after_lock == test_eth_account_init_balance + fund_amount_ibc
    # exit()

    test_sif_account_init_balance = ctx.get_sifchain_balance(test_sif_account)

    ctx.send_from_ethereum_to_sifchain(test_eth_account, test_sif_account, fund_amount_ibc, token_sc=ibc_token_sc,
                                       isLock=False)
    ctx.advance_blocks()

    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_init_balance, [[fund_amount_ibc, ibc_token_symbol]])
    test_sif_account_after_receive = ctx.get_sifchain_balance(test_sif_account)
    print("test_sif_account_after_receive is ", test_sif_account_after_receive)

    assert sifchain.balance_delta(test_sif_account_init_balance, test_sif_account_after_receive)[ibc_token_symbol] == fund_amount_ibc

    # assert test_sif_account_after_receive[ibc_token_symbol] == test_sif_account_init_balance[ibc_token_symbol] + fund_amount_ibc

    # exit()


    # lock ibc token
    # wait for ibc token balance change in eth account
    # assert balance change
    # burn ibc toke in eth side
    # wait for ibc token balance change in sifnode side
    # assert balance change

def not_test_rowan_to_eth_and_back_to_sifnode_transfer_valid(ctx):
    # get rowan contract
    rowan_sc = ctx.get_generic_erc20_sc(rowan_contract_address)

    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"], [fund_amount_eth, ctx.ceth_symbol]])

    # init balance for erc20 rowan
    test_eth_account_initial_balance = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("test_eth_account_initial_balance", test_eth_account_initial_balance)

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
    ctx.wait_for_eth_balance_change(test_eth_account, test_eth_account_initial_balance, token_addr=rowan_contract_address, timeout=90)
    eth_account_final_balance = ctx.get_erc20_token_balance(rowan_contract_address, test_eth_account)
    print("eth_account_final_balance is ", eth_account_final_balance)

    # check the rowan balance as expected
    assert eth_account_final_balance == test_eth_account_initial_balance + amount_to_lock

    test_sif_account_before_receive = ctx.get_sifchain_balance(test_sif_account)
    print("test_sif_account_before_receive is ", test_sif_account_before_receive)

    ctx.send_from_ethereum_to_sifchain(test_eth_account, test_sif_account, amount_to_lock, token_sc=rowan_sc, isLock=False)
    ctx.advance_blocks()

    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_before_receive, [[amount_to_lock, "rowan"]])
    test_sif_account_after_receive = ctx.get_sifchain_balance(test_sif_account)
    print("test_sif_account_after_receive is ", test_sif_account_after_receive)

    assert test_sif_account_after_receive["rowan"] == amount_to_lock + test_sif_account_before_receive["rowan"]



