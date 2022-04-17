import time
import threading
import copy 
from typing import List

import siftool_path
from siftool import eth, test_utils, sifchain
from siftool.eth import NULL_ADDRESS
from siftool.cosmos import balance_equal, balance_add
from siftool.common import *
from web3.eth import Contract
from web3.eth import Account
import web3

# In separate window: test/integration/framework/siftool run-env
# test/integration/framework/venv/bin/python3 test/load/test_load_tx_ethereum.py

fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1
rowan_contract_address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

threads_num = 3
amount_to_send = 10000

def bridge_bank_lock_burn(eth_tx_wrapper: eth.EthereumTxWrapper, bridge_bank_sc: Contract, test_eth_account: str, 
    recipient: str, token_amount: int, nonce: int, token_sc: Union[Contract,  None], isLock: bool):

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


def batch_create_eth_account(ctx: test_utils.EnvCtx, sc_address: Union[str, List[str], None]) -> List[Account]:
    test_eth_accounts: List[Account] = []
    for i in range(threads_num):
        test_eth_account: Account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
        test_eth_accounts.append(test_eth_account)
        if sc_address is not None:
            if isinstance(sc_address, list):
                ctx.mint_generic_erc20_token(sc_address[i], amount_to_send, test_eth_account)
            else:
                ctx.mint_generic_erc20_token(sc_address, amount_to_send, test_eth_account)
        
    return test_eth_accounts

def run_multi_thread(ctx: test_utils.EnvCtx, test_sif_account: str, test_eth_accounts: List[Account], sc_address: Union[str, List[str], None], isLock: bool):
    w3_url = ctx.w3_conn.provider.endpoint_uri
    threads: List[threading.Thread] = []
    conn_list: List[web3.Web3] = []
    eth_tx_wrappers: List[eth.EthereumTxWrapper] = []

    for i in range(threads_num):
        w3_conn = eth.web3_connect(w3_url, websocket_timeout=90)
        conn_list.append(w3_conn)
        bridge_bank_abi, _, bridge_bank_address = ctx.abi_provider.get_descriptor("BridgeBank")
        bridge_bank_sc = w3_conn.eth.contract(abi=bridge_bank_abi, address=bridge_bank_address)
        if sc_address is not None:
            sc_abi, _, _ = ctx.abi_provider.get_descriptor("BridgeToken")
            if isinstance(sc_address, list):
                erc20_contract = w3_conn.eth.contract(abi=sc_abi, address=sc_address[i])
            else:
                erc20_contract = w3_conn.eth.contract(abi=sc_abi, address=sc_address)
        else:
            erc20_contract = None
        test_eth_account = test_eth_accounts[i]
        nonce = w3_conn.eth.get_transaction_count(test_eth_account)
        recipient = test_utils.sif_addr_to_evm_arg(test_sif_account)
        eth_tx_wrapper = eth.EthereumTxWrapper(w3_conn, ctx.eth.is_local_node)
        eth_tx_wrapper.set_private_key(test_eth_account, ctx.eth._get_private_key(test_eth_account))
        eth_tx_wrappers.append(eth_tx_wrapper)
        threads.append(threading.Thread(target=bridge_bank_lock_burn, args=(
            eth_tx_wrapper, bridge_bank_sc, test_eth_account, recipient, amount_to_send, nonce, erc20_contract, isLock)))

    for t in threads:
        t.start()

    for t in threads:
        t.join()

    ctx.advance_blocks()

def test_load_burn_rowan(ctx: test_utils.EnvCtx):
    rowan_cosmos_denom = "rowan"
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)
    
    test_eth_accounts = batch_create_eth_account(ctx, rowan_contract_address)    

    run_multi_thread(ctx, test_sif_account, test_eth_accounts, rowan_contract_address, False)
    expected_change = {rowan_cosmos_denom: amount_to_send * threads_num}
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_change)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    assert balance_equal(test_sif_account_final_balance, balance_add(test_sif_account_initial_balance, expected_change))


def test_load_erc20_to_sifnode(ctx: test_utils.EnvCtx):
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    token_decimals = 18
    token_data: test_utils.ERC20TokenData = ctx.generate_random_erc20_token_data()
    erc20_sc_address = ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, token_decimals).address

    erc20_cosmos_denom = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, erc20_sc_address)

    test_eth_accounts = batch_create_eth_account(ctx, erc20_sc_address)    
    run_multi_thread(ctx, test_sif_account, test_eth_accounts, erc20_sc_address, True)

    expected_change = {erc20_cosmos_denom: amount_to_send * threads_num}
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_change)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    assert balance_equal(test_sif_account_final_balance, balance_add(test_sif_account_initial_balance, expected_change))

def test_load_multiple_erc20_to_sifnode(ctx: test_utils.EnvCtx):
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

    test_eth_accounts = batch_create_eth_account(ctx, erc20_contract_addresses)    
    run_multi_thread(ctx, test_sif_account, test_eth_accounts, erc20_contract_addresses, True)

    expected_changes = {}
    for i in range(threads_num):
        expected_changes = balance_add(expected_changes, {erc20_cosmos_denoms[i]: amount_to_send})

    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_changes)
    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    assert balance_equal(test_sif_account_final_balance, balance_add(test_sif_account_initial_balance, expected_changes))

def test_load_tx_eth(ctx: test_utils.EnvCtx):
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    test_eth_accounts = batch_create_eth_account(ctx, None)   
    run_multi_thread(ctx, test_sif_account, test_eth_accounts, None, True)

    assert amount_to_send < fund_amount_eth
    
    # Verify final balance
    expected_change = {ctx.ceth_symbol: amount_to_send * threads_num}
    ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance, expected_change)

    test_sif_account_final_balance = ctx.get_sifchain_balance(test_sif_account)
    assert balance_equal(balance_add(expected_change, test_sif_account_initial_balance), test_sif_account_final_balance)


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
