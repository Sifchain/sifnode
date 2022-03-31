import time
import threading

import siftool_path
from siftool import eth, test_utils, sifchain
import siftool.cosmos  # gPRC generated stubs use "cosmos" namespace
from siftool.common import *

import cosmos.tx.v1beta1.service_pb2 as cosmos_tx
import cosmos.tx.v1beta1.service_pb2_grpc as cosmos_tx_grpc


# How to use:
# test/integration/framework/siftool generate-python-protobuf-stubs
# ... (be patient) ...
# In separate window: test/integration/framework/siftool run-env
# test/integration/framework/venv/bin/python3 test/integration/src/peggy2/test_eth_transfers_grpc.py

# Current issues:
# python (probably from grpc):
# E0326 11:41:27.125006742  636480 fork_posix.cc:70]           Fork support is only compatible with the epoll1 and poll polling strategies
# Maybe: https://github.com/grpc/grpc/issues/29044


# Fees for sifchain -> sifchain transactions, paid by the sender.
sif_tx_fee_in_rowan = 1 * 10**17

# Fees for "ethbridge burn" transactions. Determined experimentally
sif_tx_burn_fee_in_rowan = 100000
sif_tx_burn_fee_in_ceth = 1

# There seems to be a minimum amount of rowan that a sif account needs to own in order for the bridge to do an
# "ethbridge burn". This amount does not seem to be actually used. For example, if you fund the account just with
# sif_tx_burn_fee_in_rowan, We observed that if you try to fund sif accounts with just the exact amount of rowan
# needed to pay fees (sif_tx_burn_fee_in_rowan * number_of_transactions), the bridge would stop forwarding after
# approx. 200 transactions, and you would see in sifnoded logs this message:
# {"level":"debug","module":"mempool","err":null,"peerID":"","res":{"check_tx":{"code":5,"data":null,"log":"0rowan is smaller than 500000000000000000rowan: insufficient funds: insufficient funds","info":"","gas_wanted":"1000000000000000000","gas_used":"19773","events":[],"codespace":"sdk"}},"tx":"H\ufffd\ufffdx\ufffd,4\u0004\ufffd\u001fWSnn\ufffd\ufffd\ufffdp\ufffd\ufffdg\ufffdGں^\ufffd\ufffd*i\ufffdX","time":"2022-03-26T10:09:26+01:00","message":"rejected bad transaction"}
sif_tx_burn_fee_buffer_in_rowan = 5 * sif_tx_fee_in_rowan

rowan = "rowan"


def test_load_tx_ethbridge_burn(ctx):
    # Matrix of transactions that we want to send. A row (list) in the table corresponds to a sif account sending
    # transactions to eth accounts. The numbers are transaction counts, where each transaction is for amount_per_tx.
    # Each sif account uses a dedicated send thread.
    transfer_table = [
        [100, 100, 100],
        [100, 100, 100],
        [100, 100, 100],
        [10, 20, 30],
    ]

    # transfer_table = [[1] * 3] * 4

    amount_per_tx = 123456 * eth.GWEI

    _test_load_tx_ethbridge_burn(ctx, amount_per_tx, transfer_table)


def _test_load_tx_ethbridge_burn(ctx, amount_per_tx, transfer_table, randomize=None):
    n_sif = len(transfer_table)
    assert n_sif > 0
    n_eth = len(transfer_table[0])
    assert all([len(row) == n_eth for row in transfer_table]), "transfer_table has to be rectangular"
    sum_sif = [sum(x) for x in transfer_table]
    sum_eth = [sum([x[i] for x in transfer_table]) for i in range(n_eth)]
    sum_all = sum([sum(x) for x in transfer_table])

    # Create n_sif test sif accounts.
    # Each sif account needs sif_tx_burn_fee_in_rowan * rowan + sif_tx_burn_fee_in_ceth ceth for every transaction.
    # Theoretically, we could fund accounts with ceth here, but in a strict sense this would violate the balance sheets
    # (i.e. there might be an attempt to unlock ETH in the bridge bank without enough locked ETH available).
    sif_acct_funds = [{
        rowan: sif_tx_burn_fee_in_rowan * n + sif_tx_burn_fee_buffer_in_rowan,
        # ctx.ceth_symbol: sif_tx_burn_fee_in_ceth * n
    } for n in sum_sif]
    sif_accts = [ctx.create_sifchain_addr(fund_amounts=f) for f in sif_acct_funds]

    # Create a test ethereum accounts. They are just receiving ETH, so we don't need to fund them.
    eth_accts = [ctx.create_and_fund_eth_account() for _ in range(n_eth)]

    # Get initial balances
    sif_balances_initial = [ctx.get_sifchain_balance(sif_acct) for sif_acct in sif_accts]
    eth_balances_initial = [ctx.eth.get_eth_balance(eth_acct) for eth_acct in eth_accts]
    assert all([b == 0 for b in eth_balances_initial])  # Might be non-zero if we're recycling accounts

    # Create a dispensation sif account that will receive all locked ETH and dispense it to each sif account
    # (we do this in one transaction because lock transactions take a lot of time).
    # Dispensation account needs rowan for distributing ceth to sif_accts.
    dispensation_sif_acct = ctx.create_sifchain_addr(fund_amounts={rowan: n_sif * sif_tx_fee_in_rowan})

    # Transfer ETH from operator to dispensation_sif_acct (aka lock)
    old_balances = ctx.get_sifchain_balance(dispensation_sif_acct)
    ctx.bridge_bank_lock_eth(ctx.operator, dispensation_sif_acct, sum_all * (amount_per_tx + sif_tx_burn_fee_in_ceth))
    ctx.advance_blocks()
    new_balances = ctx.wait_for_sif_balance_change(dispensation_sif_acct, old_balances)

    # Dispense from sif_dispensation_acct to sif_accts
    for i, sif_acct in enumerate(sif_accts):
        b_sif_acct_before = ctx.get_sifchain_balance(sif_acct)
        b_disp_acct_before = ctx.get_sifchain_balance(dispensation_sif_acct)
        amount_ceth = sum_sif[i] * (amount_per_tx + sif_tx_burn_fee_in_ceth)
        ctx.send_from_sifchain_to_sifchain(dispensation_sif_acct, sif_acct, {ctx.ceth_symbol: amount_ceth})
        b_sif_acct_after = ctx.wait_for_sif_balance_change(sif_acct, b_sif_acct_before)
        b_disp_acct_after = ctx.get_sifchain_balance(dispensation_sif_acct)

    # Get sif account info (for account_number and sequence)
    sif_acct_infos = [ctx.sifnode_client.query_account(sif_acct) for sif_acct in sif_accts]

    # Generate transactions.
    # If randomize is given, then we pick target ethereum accounts in random order, otherwise we go from first to last.
    start_time = time.time()
    signed_encoded_txns = []
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
                ctx.ceth_symbol, generate_only=True)
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
    def sif_acct_sender_fn(sif_acct, tx_stub, reqs):
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
    assert all([ctx.eth.get_eth_balance(eth_accts[i]) == eth_balances_initial[i] for i in range(n_eth)])  # Assert == eth_balances_initial
    assert all([ctx.get_sifchain_balance(sif_accts[i])[rowan] == sif_balances_before[i][rowan] for i in range(n_sif)])

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
    last_change_timeout = 90
    cumulative_timeout = sum_all * 10  # Equivalent to min rate of 0.1 tps
    while True:
        eth_balances = [ctx.eth.get_eth_balance(eth_acct) for eth_acct in eth_accts]
        balance_delta = sum([eth_balances[i] - eth_balances_initial[i] for i in range(n_eth)])
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
        assert actual_balance.get(rowan, 0) == sif_tx_burn_fee_buffer_in_rowan

    # Verify final eth balances
    for i, eth_acct in enumerate(eth_accts):
        expected_balance = sum_eth[i] * amount_per_tx
        actual_balance = ctx.eth.get_eth_balance(eth_acct)
        assert expected_balance == actual_balance


def choose_from(distr, rnd=None):
    r = (rnd.randrange(sum(distr))) if rnd else 0
    s = 0
    for i, di in enumerate(distr):
        s += di
        if r < s:
            distr[i] -= 1
            return i
    assert False


# Enable running directly, i.e. without pytest
if __name__ == "__main__":
    basic_logging_setup()
    from siftool import test_utils
    ctx = test_utils.get_env_ctx()
    test_load_tx_ethbridge_burn(ctx)
