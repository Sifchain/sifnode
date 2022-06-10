import random
from typing import Sequence, Any
from siftool import eth, test_utils, cosmos, sifchain
from siftool.common import *

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
# {"level":"debug","module":"mempool","err":null,"peerID":"","res":{"check_tx":{"code":5,"data":null,"log":"0rowan is smaller than 500000000000000000rowan: insufficient funds: insufficient funds","info":"","gas_wanted":"1000000000000000000","gas_used":"19773","events":[],"codespace":"sdk"}},"tx":"...","time":"2022-03-26T10:09:26+01:00","message":"rejected bad transaction"}
sif_tx_burn_fee_buffer_in_rowan = 5 * sif_tx_fee_in_rowan

ROWAN = sifchain.ROWAN

# Fee for transfering ERC20 tokens from an ethereum account to sif account (approve + lock). This is the maximum cost
# for a single transfer (regardless of amount) that the sender needs to have in his account in order for transaction to
# be processed. This value was determined experimentally with hardhat. Typical effective fee is 210542 GWEI per
# transaction, but for some reason the logic requires sender to have more funds in his account.
max_eth_transfer_fee = 10000000 * eth.GWEI


def get_sif_tx_fees(ctx):
    return {rowan: sif_tx_fee_in_rowan}


def get_sif_burn_fees(ctx):
    return {rowan: sif_tx_burn_fee_in_rowan, ctx.ceth_symbol: sif_tx_burn_fee_in_ceth}


def send_from_sifchain_to_sifchain(ctx: test_utils.EnvCtx, from_addr: cosmos.Address, to_addr: cosmos.Address,
    amounts: cosmos.Balance
) -> cosmos.Balance:
    from_balance_before = ctx.get_sifchain_balance(from_addr)
    to_balance_before = ctx.get_sifchain_balance(to_addr)
    ctx.send_from_sifchain_to_sifchain(from_addr, to_addr, amounts)
    from_expected_balance = cosmos.balance_sub(from_balance_before, amounts)
    to_expected_balance = cosmos.balance_add(to_balance_before, amounts)
    to_balance_after = ctx.wait_for_sif_balance_change(to_addr, to_balance_before, expected_balance=to_expected_balance)
    from_balance_after = cosmos.balance_sub(from_balance_before, amounts)
    assert to_balance_after == ctx.get_sifchain_balance(to_addr)
    assert cosmos.balance_equal(from_balance_after, from_expected_balance)
    assert cosmos.balance_equal(to_balance_after, to_expected_balance)
    return to_balance_after


def send_erc20_from_sifchain_to_ethereum(ctx: test_utils.EnvCtx, from_addr: cosmos.Address, to_addr: eth.Address,
    amount: int, denom: str
):
    token_address = get_erc20_token_address(ctx, denom)
    sif_balance_before = ctx.get_sifchain_balance(from_addr)
    eth_balance_before = ctx.get_erc20_token_balance(token_address, to_addr)
    ctx.sifnode_client.send_from_sifchain_to_ethereum(from_addr, to_addr, amount, denom)
    ctx.wait_for_eth_balance_change(to_addr, eth_balance_before, token_addr=token_address)
    sif_balance_after = ctx.get_sifchain_balance(from_addr)
    eth_balance_after = ctx.get_erc20_token_balance(token_address, to_addr)
    sif_burn_fees = get_sif_burn_fees(ctx)
    assert cosmos.balance_equal(sif_balance_after, cosmos.balance_sub(sif_balance_before, {denom: amount},  sif_burn_fees))
    assert eth_balance_after == eth_balance_before + amount


def get_erc20_token_address(ctx: test_utils.EnvCtx, sif_denom_hash: str) -> eth.Address:
    assert on_peggy2_branch
    ethereum_network_descriptor, token_address = sifchain.sifchain_denom_hash_to_token_contract_address(sif_denom_hash)
    assert ethereum_network_descriptor == ctx.eth.ethereum_network_descriptor  # Note: peggy2 only
    return token_address


def choose_from(distr: Sequence[Any], rnd: Optional[random.Random] = None) -> int:
    r = (rnd.randrange(sum(distr))) if rnd else 0
    s = 0
    for i, di in enumerate(distr):
        s += di
        if r < s:
            distr[i] -= 1
            return i
    assert False
