import logging.handlers
import time
import logging
from typing import Callable, Tuple, Iterable, List

import siftool_path

from siftool import cosmos, eth, sifchain, test_utils
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
        constructor_args = [name, symbol, decimals, "cosmos_denom"]  # Dummy value for cosmos_denom, actual denom will be sifBridgeDDDD0xXXX...X
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
    txhashes = []
    for token_addr in token_addrs:
        token_sc = ctx.w3_conn.eth.contract(abi=abi, address=token_addr)
        txhash = ctx.eth.transact(token_sc.functions.approve, from_eth_acct)(bridge_bank_sc.address, amount)
        txhashes.append(txhash)
    ctx.eth.wait_for_all_transaction_receipts(txhashes)

    to_sif_acct_encoded = test_utils.sif_addr_to_evm_arg(to_sif_acct)
    txrcpts = []
    for token_addr in token_addrs:
        tx_opts = {"value": 0}
        txrcpt = ctx.eth.transact(bridge_bank_sc.functions.lock, from_eth_acct, tx_opts=tx_opts)(to_sif_acct_encoded, token_addr, amount)
        txrcpts.append(txrcpt)
    ctx.eth.wait_for_all_transaction_receipts(txrcpts)


def test(ctx):
    _test(ctx, 5000)


def _test(ctx, number_of_erc20_tokens):
    token_name_base = random_string(4)
    owner = ctx.operator

    eth_sender = ctx.create_and_fund_eth_account(fund_amount=max_eth_transfer_fee * number_of_erc20_tokens)
    sif_recipient = ctx.create_sifchain_addr()

    def token_data_provider(i: int) -> Tuple[str, str, int]:
        token_name = "{}{}".format(token_name_base, i)
        token_symbol = "eth-symbol-{}".format(i)
        token_decimals = 6
        return token_name, token_symbol, token_decimals

    start_time = time.time()
    time_before = start_time
    contract_addresses = batch_deploy_erc20_tokens(ctx, number_of_erc20_tokens, ctx.operator, token_data_provider)
    deploy_time = time.time() - time_before

    log.debug("batch_deploy_erc20_tokens(): {:.2f} s, {:.2f} items/s".format(deploy_time,
        number_of_erc20_tokens / deploy_time if deploy_time > 0 else 0))

    sif_denoms = [sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, addr) for addr in contract_addresses]

    amount = 123456

    time_before = time.time()
    batch_mint_erc20_tokens(ctx, owner, eth_sender, amount, contract_addresses)
    mint_time = time.time() - time_before

    log.debug("batch_mint_erc20_tokens(): {:.2f} s, {:.2f} items/s".format(mint_time,
        number_of_erc20_tokens / mint_time if mint_time > 0 else 0))

    eth_balance_before = ctx.eth.get_eth_balance(eth_sender)
    sif_balance_before = ctx.get_sifchain_balance(sif_recipient)

    time_before = time.time()
    batch_approve_and_lock_erc20_tokens(ctx, eth_sender, sif_recipient, contract_addresses, amount)
    approve_and_lock_time = time.time() - time_before
    log.debug("batch_approve_and_lock_erc20_tokens(): {:.2f} s, {:.2f} items/s".format(approve_and_lock_time,
        number_of_erc20_tokens / approve_and_lock_time if approve_and_lock_time > 0 else 0))

    expected_balance = {denom: amount for denom in sif_denoms}
    time_before = time.time()
    sif_balance_after = ctx.wait_for_sif_balance_change(sif_recipient, sif_balance_before, expected_balance=expected_balance)
    balance_change_time = time.time() - time_before

    log.debug("wait_for_sif_balance_change(): {:.2f} s, {:.2f} items/s".format(balance_change_time,
        number_of_erc20_tokens / balance_change_time if balance_change_time > 0 else 0))

    total_time = time.time() - start_time

    eth_balance_after = ctx.eth.get_eth_balance(eth_sender)

    log.debug("Total: {:.2f} s, {:.2f} items/s".format(total_time,
        number_of_erc20_tokens / total_time if total_time > 0 else 0))
