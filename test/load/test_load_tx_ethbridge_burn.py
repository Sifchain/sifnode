import time
import threading
from typing import List, Any, Iterable
from web3.eth import Contract

import siftool_path
from siftool.eth import NULL_ADDRESS
from load_testing import *
from siftool import eth, test_utils, cosmos, sifchain
from siftool.common import *

import cosmos.tx.v1beta1.service_pb2 as cosmos_tx
import cosmos.tx.v1beta1.service_pb2_grpc as cosmos_tx_grpc

rowan_contract_address = "0x000000000000000"

# How to use:
# test/integration/framework/siftool generate-python-protobuf-stubs
# ... (be patient) ...
# In separate window: test/integration/framework/siftool run-env
# test/integration/framework/venv/bin/python3 test/integration/src/peggy2/test_eth_transfers_grpc.py

# Current issues:
# python (probably from grpc):
# E0326 11:41:27.125006742  636480 fork_posix.cc:70]           Fork support is only compatible with the epoll1 and poll polling strategies
# Maybe: https://github.com/grpc/grpc/issues/29044

eth_account_number = 2
sif_account_number = 1
transaction_number = 2

def build_transfer_table() -> [[int]]:
    transfer_table = [[int]]
    for i in range(sif_account_number):
        column = []
        for i in range(eth_account_number):
            column.append(transaction_number)
        transfer_table.append(column)
    return transfer_table

def burn_rowan_get_erc20_address(ctx: test_utils.EnvCtx):
    token_address = ctx.get_destination_contract_address(rowan)
    if token_address != NULL_ADDRESS:
        return token_address
    fund_amount_sif = 10 ** 20
    fund_amount_eth = 10 ** 20
    rowan_transfer_amount = 10 ** 18
    test_sif_account = ctx.create_sifchain_addr(fund_amounts={rowan:fund_amount_sif, ctx.ceth_symbol:fund_amount_eth})
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, rowan_transfer_amount, rowan)
    ctx.advance_blocks()
    token_address = ctx.wait_for_new_bridge_token_created(rowan)
    return token_address

# test single sif account burn erc20 to multiple ethereum accounts
def test_single_sif_to_multiple_eth_account_lock_rowan(ctx: test_utils.EnvCtx):
    # get rowan contract address
    rowan_token_address = burn_rowan_get_erc20_address(ctx)
    rowan_sc = ctx.get_generic_erc20_sc(rowan_token_address)

    transfer_table = build_transfer_table()
    amount_per_tx = 1000100101

    _test_load_tx_ethbridge_lock_burn(ctx, amount_per_tx, transfer_table, rowan_sc, isRowan=True)

# test single sif account burn erc20 to multiple ethereum accounts
def test_single_sif_to_multiple_eth_account_burn_erc20(ctx: test_utils.EnvCtx):
    # deploy an erc20 contract
    token_decimals = 18
    total_amount = 10 ** 28
    token_data: test_utils.ERC20TokenData = ctx.generate_random_erc20_token_data()
    erc20_sc = ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, token_decimals)
    # mint token to operator account
    ctx.mint_generic_erc20_token(erc20_sc.address, total_amount, ctx.operator) 
    ctx.advance_blocks()

    transfer_table = build_transfer_table()
    amount_per_tx = 1000100101

    _test_load_tx_ethbridge_lock_burn(ctx, amount_per_tx, transfer_table, erc20_sc)

# test single sif account burn ceth to multiple ethereum accounts
def test_single_sif_to_multiple_eth_account_burn_eth(ctx: test_utils.EnvCtx):
    transfer_table = build_transfer_table()
    amount_per_tx = 1000100101
    eth_sc = ctx.get_generic_erc20_sc(NULL_ADDRESS)

    _test_load_tx_ethbridge_lock_burn(ctx, amount_per_tx, transfer_table, eth_sc)

def test_load_tx_ethbridge_burn_eth_short(ctx: test_utils.EnvCtx):
    transfer_table = [[2, 2], [2, 2]]
    amount_per_tx = 1000100101
    eth_sc = ctx.get_generic_erc20_sc(NULL_ADDRESS)
    _test_load_tx_ethbridge_lock_burn(ctx, amount_per_tx, transfer_table, eth_sc)

# test multiple sif accounts burn ceth to multiple ethereum accounts
def test_load_tx_ethbridge_burn_eth(ctx: test_utils.EnvCtx):   
    # Matrix of transactions that we want to send. A row (list) in the table corresponds to a sif account sending
    # transactions to eth accounts. The numbers are transaction counts, where each transaction is for amount_per_tx.
    # Each sif account uses a dedicated send thread.
    transfer_table = [
            [100, 100, 100],
            [100, 100, 100],
            [100, 100, 100],
            [10, 20, 30],
        ]
    amount_per_tx = 1000100101
    eth_sc = ctx.get_generic_erc20_sc(NULL_ADDRESS)
    _test_load_tx_ethbridge_lock_burn(ctx, amount_per_tx, transfer_table, eth_sc)

def _test_load_tx_ethbridge_lock_burn(ctx: test_utils.EnvCtx, amount_per_tx: int, 
    transfer_table: List[List[int]], token_sc: Contract, isRowan: bool = False, randomize: bool = None):
    # rowan is natvie token, denom not from contract in Ethereum
    if isRowan:
        token_denom = rowan
    else:
        token_denom = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, token_sc.address)

    n_sif: int = len(transfer_table)
    assert n_sif > 0
    n_eth: int = len(transfer_table[0])
    assert all([len(row) == n_eth for row in transfer_table]), "transfer_table has to be rectangular"
    sum_sif: List[int] = [sum(x) for x in transfer_table]
    sum_eth: List[int] = [sum([x[i] for x in transfer_table]) for i in range(n_eth)]
    sum_all: int = sum([sum(x) for x in transfer_table])

    # Create n_sif test sif accounts.
    # Each sif account needs sif_tx_burn_fee_in_rowan * rowan + sif_tx_burn_fee_in_ceth ceth for every transaction.
    # Theoretically, we could fund accounts with ceth here, but in a strict sense this would violate the balance sheets
    # (i.e. there might be an attempt to unlock ETH in the bridge bank without enough locked ETH available).
    sif_acct_funds: List[cosmos.Balance] = [{
        rowan: sif_tx_burn_fee_in_rowan * n + sif_tx_burn_fee_buffer_in_rowan,
        # ctx.ceth_symbol: sif_tx_burn_fee_in_ceth * n
    } for n in sum_sif]
    sif_accts: List[cosmos.Address] = [ctx.create_sifchain_addr(fund_amounts=f) for f in sif_acct_funds]

    # Create a test ethereum accounts. They are just receiving ETH, so we don't need to fund them.
    eth_accts: List[str] = [ctx.create_and_fund_eth_account() for _ in range(n_eth)]

    # Get initial balances
    sif_balances_initial: List[cosmos.Balance] = [ctx.get_sifchain_balance(sif_acct) for sif_acct in sif_accts]
    eth_balances_initial: List[str] = [ctx.eth.get_eth_balance(eth_acct) for eth_acct in eth_accts]
    assert all([b == 0 for b in eth_balances_initial])  # Might be non-zero if we're recycling accounts
    
    if token_denom != ctx.ceth_symbol:
        erc20_balances_initial: List[str] = [ctx.get_erc20_token_balance(token_sc.address, eth_acct) for eth_acct in eth_accts]
        assert all([b == 0 for b in eth_balances_initial])  # Might be non-zero if we're recycling accounts


    # Create a dispensation sif account that will receive all locked ETH and dispense it to each sif account
    # (we do this in one transaction because lock transactions take a lot of time).
    # Dispensation account needs rowan for distributing ceth to sif_accts.
    if token_denom == rowan:
        amount = sum_all * amount_per_tx + n_sif * 2 * sif_tx_fee_in_rowan
        dispensation_sif_acct: cosmos.Address = ctx.create_sifchain_addr(fund_amounts={rowan: amount})
    elif token_denom == ctx.ceth_symbol:
        dispensation_sif_acct: cosmos.Address = ctx.create_sifchain_addr(fund_amounts={rowan: n_sif * sif_tx_fee_in_rowan})
    else:
        # for erc20 token, dispensation_sif_acct need distribute both ceth for cross chain fee and erc20 to burn
        dispensation_sif_acct: cosmos.Address = ctx.create_sifchain_addr(fund_amounts={rowan: n_sif * 2 * sif_tx_fee_in_rowan})

    # Transfer ETH from operator to dispensation_sif_acct (aka lock)
    old_balances = ctx.get_sifchain_balance(dispensation_sif_acct)
    lock_eth_amount = sif_tx_burn_fee_in_ceth
    if token_denom == ctx.ceth_symbol:
        lock_eth_amount += amount_per_tx
    ctx.bridge_bank_lock_eth(ctx.operator, dispensation_sif_acct, sum_all * lock_eth_amount)
    ctx.advance_blocks()
    old_balances = ctx.wait_for_sif_balance_change(dispensation_sif_acct, old_balances)

    # just for erc20 token
    if str.startswith(token_denom, "sifBridge") and token_denom != ctx.ceth_symbol:
        ctx.send_from_ethereum_to_sifchain(ctx.operator, dispensation_sif_acct, sum_all * amount_per_tx, token_sc=token_sc, isLock=True)
        _ = ctx.wait_for_sif_balance_change(dispensation_sif_acct, old_balances)

    # Dispense from sif_dispensation_acct to sif_accts
    for i, sif_acct in enumerate(sif_accts):
        b_sif_acct_before = ctx.get_sifchain_balance(sif_acct)
        b_disp_acct_before = ctx.get_sifchain_balance(dispensation_sif_acct)

        # if token is ceth, combine and fund
        if ctx.ceth_symbol == token_denom:
            amount = sum_sif[i] * (amount_per_tx + sif_tx_burn_fee_in_ceth)
            ctx.send_from_sifchain_to_sifchain(dispensation_sif_acct, sif_acct, {ctx.ceth_symbol: amount})
            b_sif_acct_after = ctx.wait_for_sif_balance_change(sif_acct, b_sif_acct_before)

        else:
            ctx.send_from_sifchain_to_sifchain(dispensation_sif_acct, sif_acct, {ctx.ceth_symbol: sum_sif[i] * sif_tx_burn_fee_in_ceth})
            # rowan got when the sif account init
            # if token_denom != rowan:
            b_sif_acct_before = ctx.wait_for_sif_balance_change(sif_acct, b_sif_acct_before)
            ctx.send_from_sifchain_to_sifchain(dispensation_sif_acct, sif_acct, {token_denom: sum_sif[i] * amount_per_tx})

            b_sif_acct_after = ctx.wait_for_sif_balance_change(sif_acct, b_sif_acct_before)
        b_disp_acct_after = ctx.get_sifchain_balance(dispensation_sif_acct)

    # Get sif account info (for account_number and sequence)
    sif_acct_infos = [ctx.sifnode_client.query_account(sif_acct) for sif_acct in sif_accts]

    # Generate transactions.
    # If randomize is given, then we pick target ethereum accounts in random order, otherwise we go from first to last.
    start_time: float = time.time()
    signed_encoded_txns: List[List[bytes]] = []
    for i in range(n_sif):
        sif_acct = sif_accts[i]
        log.debug("Generating {} txns from {}...".format(sum_sif[i], sif_acct))
        account_number = int(sif_acct_infos[i]["account_number"])
        sequence = int(sif_acct_infos[i]["sequence"])
        choice_histogram = transfer_table[i].copy()
        ordereed_eth_accounts = [eth_accts[choose_from(choice_histogram, rnd=randomize)] for _ in range(sum_sif[i])]
        txn_list = []
        rowan_before = None
        fee_histogram = {}
        for eth_acct in ordereed_eth_accounts:
            tx = ctx.sifnode_client.send_from_sifchain_to_ethereum(sif_acct, eth_acct, amount_per_tx,
                token_denom, generate_only=True)
            signed_tx = ctx.sifnode_client.sign_transaction(tx, sif_acct, sequence=sequence,
                account_number=account_number)
            encoded_tx = ctx.sifnode_client.encode_transaction(signed_tx)
            rowan_after = ctx.get_sifchain_balance(sif_acct)[rowan]
            if rowan_before is not None:
                rowan_fee = rowan_after - rowan_before
                fee_histogram[rowan_fee] = fee_histogram.get(rowan_fee, 0) + 1
            rowan_before = rowan_after
            txn_list.append(encoded_tx)
            sequence += 1
        signed_encoded_txns.append(txn_list)
        log.debug("Fee histogram: {}".format(repr(fee_histogram)))
    log.debug("Transaction generation speed: {:.2f}/s".format(sum_all / (time.time() - start_time)))

    # Per-thread function for broadcasting transactions
    def sif_acct_sender_fn(sif_acct: cosmos.Address, tx_stub: cosmos_tx_grpc.ServiceStub, reqs: Sequence[cosmos_tx.BroadcastTxRequest]):
        log.debug("Broadcasting {} txns from {}...".format(len(reqs), sif_acct))
        for req in reqs:
            tx_stub.BroadcastTx(req)

    # Prepare gRPC messages and sending threads. We use one thread for each sif_accts since the messages for account
    # have to be broadcast by ascending sequnce number.
    threads = []
    channels = []
    broadcast_mode = cosmos_tx.BROADCAST_MODE_ASYNC
    for i in range(n_sif):
        sif_acct = sif_accts[i]
        channel = ctx.sifnode_client.open_grpc_channel()
        channels.append(channel)
        tx_stub = cosmos_tx_grpc.ServiceStub(channel)
        reqs = [cosmos_tx.BroadcastTxRequest(tx_bytes=tx_bytes, mode=broadcast_mode) for tx_bytes in signed_encoded_txns[i]]
        threads.append(threading.Thread(target=sif_acct_sender_fn, args=(sif_acct, tx_stub, reqs)))

    # Get initial balances. The balances should not have been changed by now.
    sif_balances_before = [ctx.get_sifchain_balance(x) for x in sif_accts]  # Assert == sif_balances_initial (for rowan)
    assert all(ctx.eth.get_eth_balance(eth_accts[i]) == eth_balances_initial[i] for i in range(n_eth))
    assert all(ctx.get_sifchain_balance(sif_accts[i])[rowan] == sif_balances_before[i][rowan] for i in range(n_sif))

    # Broadcast transactions
    start_time = time.time()
    for t in threads:
        t.start()

    for t in threads:
        t.join()

    log.debug("Transaction broadcast speed: {:.2f}/s".format(sum_all / (time.time() - start_time)))

    for c in channels:
        c.close()

    avg_tx_fees = [(sif_balances_before[i].get(rowan, 0) - ctx.get_sifchain_balance(sif_accts[i]).get(rowan, 0)) / sum_sif[i] for i in range(n_sif)]
    log.debug("Average used fee per transaction: {}".format(repr(avg_tx_fees)))

    # Wait for eth balances
    start_time = time.time()
    last_change_time = None
    last_change = None
    last_change_timeout = 180
    cumulative_timeout = 90  # Equivalent to min rate of 0.1 tps
    while True:
        if token_denom == ctx.ceth_symbol:
            token_balances = [ctx.eth.get_eth_balance(eth_acct) for eth_acct in eth_accts]
            balance_delta = sum([token_balances[i] - eth_balances_initial[i] for i in range(n_eth)])
        else:
            token_balances = [ctx.get_erc20_token_balance(token_sc.address, eth_acct) for eth_acct in eth_accts]
            balance_delta = sum([token_balances[i] - eth_balances_initial[i] for i in range(n_eth)])

        now = time.time()
        total = sum_all * amount_per_tx

        still_to_go = total - balance_delta
        pct_done = balance_delta / total * 100
        txns_done = balance_delta / amount_per_tx
        time_elapsed = time.time() - start_time
        log.debug("Test progress: {} / {} ({:.9f} txns done, {:.2f}%, {:.4f} avg tps)".format(balance_delta, total,
            txns_done, pct_done, (txns_done / time_elapsed if time_elapsed > 0 else 0)))
        if (last_change is None) or (balance_delta != last_change):
            last_change_time = now
            last_change = balance_delta
            log.debug("sif_accts balances: {}".format(repr([ctx.get_sifchain_balance(x) for x in sif_accts])))
        if still_to_go == 0:
            break
        if now - last_change_time > last_change_timeout:
            raise Exception("Last change timeout exceeded")
        elif now - start_time > cumulative_timeout:
            raise Exception("Cumulative timeout exceeded")
        time.sleep(3)

    # Verify final sif balances. There should be no ceth left. There is some rowan left since we oversupplied it.
    for sif_acct in sif_accts:
        actual_balance = ctx.get_sifchain_balance(sif_acct)
        assert actual_balance.get(ctx.ceth_symbol, 0) == 0
        print("========= actual_balance.get(rowan", actual_balance.get(rowan, 0))
        assert actual_balance.get(rowan, 0) == sif_tx_burn_fee_buffer_in_rowan

    # Verify final eth balances
    for i, eth_acct in enumerate(eth_accts):
        expected_balance = sum_eth[i] * amount_per_tx
        if token_denom == ctx.ceth_symbol:
            actual_balance = ctx.eth.get_eth_balance(eth_acct)
        else:
            actual_balance = ctx.get_erc20_token_balance(token_sc.address, eth_acct)
        assert expected_balance == actual_balance


# Enable running directly, i.e. without pytest
if __name__ == "__main__":
    basic_logging_setup()
    from siftool import test_utils
    ctx = test_utils.get_env_ctx()
    test_single_sif_to_multiple_eth_account_lock_rowan(ctx)
    test_single_sif_to_multiple_eth_account_burn_erc20(ctx)
    test_single_sif_to_multiple_eth_account_burn_eth(ctx)
    test_load_tx_ethbridge_burn_eth_short(ctx)
    test_load_tx_ethbridge_burn_eth(ctx)



