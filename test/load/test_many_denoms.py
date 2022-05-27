import sys
import time
import logging
import argparse
from typing import Callable, Tuple, Iterable, List

import siftool_path
from siftool import cosmos, eth, sifchain, test_utils, command
from siftool.cosmos import balance_add, balance_sub, balance_equal, balance_mul
from siftool.common import *
from load_testing import *


log = logging.getLogger(__name__)


def batch_deploy_erc20_tokens(ctx: test_utils.EnvCtx, count: int, deployer_addr: eth.Address,
    token_data_provider: Callable[[int], Tuple[str, str, int]]
) -> List[eth.Address]:
    abi, bytecode, _ = ctx.abi_provider.get_descriptor(ctx.generic_erc20_contract)
    token_sc = ctx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
    txhashes = []
    for i in range(count):
        name, symbol, decimals = token_data_provider(i)
        constructor_args = [name, symbol, decimals, "dummy_value_for_cosmos_denom"]  # Dummy value for cosmos_denom, actual denom will be sifBridgeDDDD0xXXX...X
        txhash = ctx.eth.transact(token_sc.constructor, deployer_addr)(*constructor_args)
        txhashes.append(txhash)
    txrcpts = ctx.eth.wait_for_all_transaction_receipts(txhashes)
    contract_addresses = [txrcpt.contractAddress for txrcpt in txrcpts]
    return contract_addresses


def batch_mint_erc20_tokens(ctx: test_utils.EnvCtx, minter_account: eth.Address,
    minted_tokens_recipient: eth.Address, amount: int, contracts: Iterable[eth.Address]
):
    txhashes = []
    abi, bytecode, _ = ctx.abi_provider.get_descriptor(ctx.generic_erc20_contract)
    for contract_address in contracts:
        token_sc = ctx.w3_conn.eth.contract(abi=abi, address=contract_address)
        txhash = ctx.eth.transact(token_sc.functions.mint, minter_account)(minted_tokens_recipient, amount)
        txhashes.append(txhash)
    ctx.eth.wait_for_all_transaction_receipts(txhashes)


def batch_approve_and_lock_erc20_tokens(ctx: test_utils.EnvCtx, from_eth_acct: eth.Address,
    to_sif_acct: cosmos.Address, token_addrs: Iterable[eth.Address], amount: int
):
    abi, bytecode, _ = ctx.abi_provider.get_descriptor(ctx.generic_erc20_contract)
    bridge_bank_sc = ctx.get_bridge_bank_sc()
    to_sif_acct_encoded = test_utils.sif_addr_to_evm_arg(to_sif_acct)
    txhashes = []
    for token_addr in token_addrs:
        token_sc = ctx.w3_conn.eth.contract(abi=abi, address=token_addr)
        txhash = ctx.eth.transact(token_sc.functions.approve, from_eth_acct)(bridge_bank_sc.address, amount)
        txhashes.append(txhash)
    ctx.eth.wait_for_all_transaction_receipts(txhashes)
    txhashes = []
    for token_addr in token_addrs:
        tx_opts = {"value": 0}
        txhash = ctx.eth.transact(bridge_bank_sc.functions.lock, from_eth_acct, tx_opts=tx_opts)(to_sif_acct_encoded, token_addr, amount)
        txhashes.append(txhash)
    ctx.eth.wait_for_all_transaction_receipts(txhashes)


def timed_tx_send_loop(ctx: test_utils.EnvCtx, src_sif_address: cosmos.Address, dst_sif_address: cosmos.Address,
    tx_amount: cosmos.Balance, loop_count: int
) -> float:
    sif_tx_fee = get_sif_tx_fees(ctx)
    send_from_sifchain_to_sifchain(ctx, ctx.rowan_source, src_sif_address, balance_mul(sif_tx_fee, loop_count))
    src_initial_balance = ctx.get_sifchain_balance(src_sif_address)
    dst_initial_balance = ctx.get_sifchain_balance(dst_sif_address)
    time_before = time.time()
    for i in range(loop_count):
        send_from_sifchain_to_sifchain(ctx, src_sif_address, dst_sif_address, tx_amount)
        src_expected_balance = balance_sub(src_initial_balance, balance_mul(balance_add(tx_amount, sif_tx_fee), i + 1))
        dst_expected_balance = balance_add(dst_initial_balance, balance_mul(tx_amount, i + 1))
        src_actual_balance = ctx.get_sifchain_balance(src_sif_address)
        dst_actual_balance = ctx.get_sifchain_balance(dst_sif_address)
        assert balance_equal(src_actual_balance, src_expected_balance)
        assert balance_equal(dst_actual_balance, dst_expected_balance)
    return time.time() - time_before


def timed_tx_burn_loop(ctx: test_utils.EnvCtx, src_sif_address: cosmos.Address, dst_eth_address: eth.Address,
    denom: str, amount: int, loop_count: int
) -> float:
    sif_burn_fee = get_sif_burn_fees(ctx)
    send_from_sifchain_to_sifchain(ctx, ctx.rowan_source, src_sif_address, balance_mul(sif_burn_fee, loop_count))
    tx_denom_token_address = get_erc20_token_address(ctx, denom)
    src_initial_balance = ctx.get_sifchain_balance(src_sif_address)
    dst_initial_balance = ctx.get_erc20_token_balance(tx_denom_token_address, dst_eth_address)
    time_before = time.time()
    for i in range(loop_count):
        send_erc20_from_sifchain_to_ethereum(ctx, src_sif_address, dst_eth_address, amount, denom)
        src_expected_balance = balance_sub(src_initial_balance, balance_mul(balance_add({denom: amount}, sif_burn_fee), i + 1))
        dst_expected_balance = dst_initial_balance + amount * (i + 1)
        src_actual_balance = ctx.get_sifchain_balance(src_sif_address)
        dst_actual_balance = ctx.get_erc20_token_balance(tx_denom_token_address, dst_eth_address)
        assert balance_equal(src_actual_balance, src_expected_balance)
        assert dst_actual_balance == dst_expected_balance
    return time.time() - time_before


def timed_query_bank_balances_loop(ctx: test_utils.EnvCtx, sif_address: cosmos.Address, paginate=None):
    ctx.get_sifchain_balance(sif_address)


def report(report_lines: List[str], message: str):
    log.info(message)
    report_lines.append(message)


def test(ctx: test_utils.EnvCtx):
    _parametric_test(ctx, 2)


def _parametric_test(ctx: test_utils.EnvCtx, number_of_erc20_tokens: int, sample_loop_size: int = 20,
    report_lines: Optional[List[str]] = None,
):
    report_lines = [] if report_lines is None else report_lines
    assert number_of_erc20_tokens > 1
    assert sample_loop_size > 0

    token_name_base = random_string(4)
    owner = ctx.operator
    test_start_time = time.time()

    eth_sender = ctx.create_and_fund_eth_account(fund_amount=max_eth_transfer_fee * number_of_erc20_tokens)
    fat_sif_wallet = ctx.create_sifchain_addr()
    slim_sif_wallet = ctx.create_sifchain_addr()

    def token_data_provider(i: int) -> Tuple[str, str, int]:
        token_name = "{}{}".format(token_name_base, i)
        token_symbol = "eth-symbol-{}".format(i)
        token_decimals = 6
        return token_name, token_symbol, token_decimals

    report(report_lines, "Number of tokens: {}".format(number_of_erc20_tokens))

    setup_start_time = time.time()
    time_before = time.time()
    contract_addresses = batch_deploy_erc20_tokens(ctx, number_of_erc20_tokens, ctx.operator, token_data_provider)
    deploy_time = time.time() - time_before

    report(report_lines, "batch_deploy_erc20_tokens(): {:.2f} s, {:.2f} items/s".format(deploy_time,
        number_of_erc20_tokens / deploy_time if deploy_time > 0 else 0))

    sif_denoms = [sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, addr) for addr in contract_addresses]

    amount = 123456  # Must be > test_loop_count

    time_before = time.time()
    batch_mint_erc20_tokens(ctx, owner, eth_sender, amount, contract_addresses)
    mint_time = time.time() - time_before

    report(report_lines, "batch_mint_erc20_tokens(): {:.2f} s, {:.2f} items/s".format(mint_time,
        number_of_erc20_tokens / mint_time if mint_time > 0 else 0))

    eth_balance_before = ctx.eth.get_eth_balance(eth_sender)
    sif_balance_before = ctx.get_sifchain_balance(fat_sif_wallet)

    time_before = time.time()
    batch_approve_and_lock_erc20_tokens(ctx, eth_sender, fat_sif_wallet, contract_addresses, amount)
    approve_and_lock_time = time.time() - time_before
    report(report_lines, "batch_approve_and_lock_erc20_tokens(): {:.2f} s, {:.2f} items/s".format(approve_and_lock_time,
        number_of_erc20_tokens / approve_and_lock_time if approve_and_lock_time > 0 else 0))

    fat_balance = {denom: amount for denom in sif_denoms}
    time_before = time.time()
    ctx.wait_for_sif_balance_change(fat_sif_wallet, sif_balance_before, expected_balance=fat_balance)
    balance_change_time = time.time() - time_before
    assert balance_equal(ctx.get_sifchain_balance(fat_sif_wallet), fat_balance)

    report(report_lines, "wait_for_sif_balance_change(): {:.2f} s, {:.2f} items/s".format(balance_change_time,
        number_of_erc20_tokens / balance_change_time if balance_change_time > 0 else 0))

    setup_time = time.time() - setup_start_time

    eth_balance_after = ctx.eth.get_eth_balance(eth_sender)
    report(report_lines, "Cost of approve+lock: {:.2f} gwei".format(
        (eth_balance_before - eth_balance_after) / number_of_erc20_tokens / eth.GWEI))

    report(report_lines, "Total for setup (mint+approve+lock): {:.2f} s, {:.2f} items/s".format(setup_time,
        number_of_erc20_tokens / setup_time if setup_time > 0 else 0))

    # Populate slim wallet with only one denom. Slim wallet is used to compare relative performance of a sif account
    # with just one denom to performance of a sif account with many denoms, all other things being equal.
    token0 = contract_addresses[0]
    denom0 = sif_denoms[0]
    expected_slim_balance = {denom0: amount}
    slim_balance_before = ctx.get_sifchain_balance(slim_sif_wallet)
    batch_mint_erc20_tokens(ctx, owner, eth_sender, amount, [token0])
    batch_approve_and_lock_erc20_tokens(ctx, eth_sender, slim_sif_wallet, [token0], amount)
    ctx.wait_for_sif_balance_change(slim_sif_wallet, sif_balance_before, expected_balance=expected_slim_balance)
    assert cosmos.balance_equal(ctx.get_sifchain_balance(slim_sif_wallet), expected_slim_balance)

    # We do few transfers and burns to get the average time per transaction for fat wallet.
    # We assert that the fees are equal to the fee of a single transaction as given by get_sif_tx_fees() and
    # get_sif_burn_fees().

    tmp_sif_account = ctx.create_sifchain_addr()
    tmp_eth_account = ctx.create_and_fund_eth_account()

    sif_send_time_fat = timed_tx_send_loop(ctx, fat_sif_wallet, tmp_sif_account, {denom0: 1}, sample_loop_size)
    report(report_lines, "Average tx_send time for fat_sif_wallet: {:.2f} s".format(sif_send_time_fat / sample_loop_size))

    sif_burn_time_fat = timed_tx_burn_loop(ctx, fat_sif_wallet, tmp_eth_account, denom0, 1, sample_loop_size)
    report(report_lines, "Average tx_burn time for fat_sif_wallet: {:.2f} s".format(sif_burn_time_fat / sample_loop_size))

    sif_send_time_slim = timed_tx_send_loop(ctx, slim_sif_wallet, tmp_sif_account, {denom0: 1}, sample_loop_size)
    report(report_lines, "Average tx_send time for slim_sif_wallet: {:.2f} s".format(sif_send_time_slim / sample_loop_size))

    sif_burn_time_slim = timed_tx_burn_loop(ctx, slim_sif_wallet, tmp_eth_account, denom0, 1, sample_loop_size)
    report(report_lines, "Average tx_burn time for slim_sif_wallet: {:.2f} s".format(sif_burn_time_slim / sample_loop_size))

    report(report_lines, "Relative fat/slim speed for tx_send: {:.2f}".format(sif_send_time_fat / sif_send_time_slim))
    report(report_lines, "Relative fat/slim speed for tx_burn: {:.2f}".format(sif_burn_time_fat / sif_burn_time_slim))

    test_total_time = time.time() - test_start_time
    report(report_lines, "Total test time: {:.2f} s".format(test_total_time))


# Enable running directly, i.e. without pytest
if __name__ == "__main__":
    basic_logging_setup()
    ctx = test_utils.get_env_ctx()
    parser = argparse.ArgumentParser()
    parser.add_argument("--count", type=int, default=2)
    parser.add_argument("--sample-loop-size", type=int, default=10)
    parser.add_argument("--report")
    args = parser.parse_args(sys.argv[1:])
    report_lines = []
    _parametric_test(ctx, args.count, sample_loop_size=args.sample_loop_size, report_lines=report_lines)
    if args.report:
        command.Command().write_text_file(args.report, joinlines(report_lines))
    log.info("Finished successfully")
