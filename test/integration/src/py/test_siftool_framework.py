import logging
import pytest

import siftool_path
from siftool import eth
from siftool.common import *


# Note: these tests burn a lot of ether very inefficiently. If you care about that make sure to recover
# any funds from test accounts.
# max_gas_price has to be sufficiently high so that test_account will be able to send the remaining funds back.
# Especially with ganache the account balance has to be much higher than the actual fee, otherwise ganache will
# reject the transaction with "sender doesn't have enough funds to send tx. The upfront cost is: 10000123000000000".
# (This upfront fee corresponds to approx. gas price of 500 gwei).

min_tx_gas = eth.MIN_TX_GAS
max_tx_gas = 5000000
min_gas_price = 1 * eth.GWEI
max_gas_price = 500 * eth.GWEI
min_tx_cost = min_tx_gas * min_gas_price
max_tx_cost = max_tx_gas * max_gas_price


def test_sanity_checks(ctx):
    ctx.sanity_check()

def test_eth_fee_functions_1(ctx):
    e = ctx.eth.w3_conn.eth
    null_txn = {"to": eth.NULL_ADDRESS}
    gas_price = None
    try:
        gas_price = e.gas_price
    except Exception as ex:
        assert ctx.eth.is_contract_logic_error_method_not_found(ex)
    if gas_price is not None:
        assert gas_price > 1 * eth.GWEI
    fee_history = None
    try:
        fee_history = e.fee_history(1, "latest", [25])
    except Exception as ex:
        assert ctx.eth.is_contract_logic_error_method_not_found(ex, "eth_maxPriorityFeePerGas")
    if fee_history is not None:
        assert "baseFeePerGas" in fee_history
        assert "oldestBlock" in fee_history
        assert "reward" in fee_history
        assert "baseFeePerGas" in fee_history
        assert len(fee_history.reward) == 1
        assert len(fee_history.reward[0]) == 1
    is_hardhat = ctx.eth.is_local_node and on_peggy2_branch
    if is_hardhat:
        assert e.estimate_gas(null_txn) == eth.MIN_TX_GAS + 1
    else:
        assert e.estimate_gas(null_txn) == eth.MIN_TX_GAS
    assert gas_price is not None
    if ctx.eth.is_local_node:
        assert (fee_history is None) == (not is_hardhat)
        assert ctx.eth.gas_estimate_fn is None
        assert bool(ctx.eth.fixed_gas_args)
    else:
        assert fee_history is not None
        assert ctx.eth.gas_estimate_fn is not None
        gas, max_fee_per_gas, max_priority_fee_per_gas, gas_price = ctx.eth.gas_estimate_fn(null_txn)
    return

@pytest.mark.filterwarnings("ignore::UserWarning")
def test_eth_fee_functions_2(ctx):
    e = ctx.eth.w3_conn.eth
    max_priority_fee = None
    try:
        max_priority_fee = e.max_priority_fee
    except Exception as ex:
        assert ctx.eth.is_contract_logic_error_method_not_found(ex, "eth_maxPriorityFeePerGas")
    if ctx.eth.is_local_node:
        assert max_priority_fee is None
    else:
        assert max_priority_fee is not None
        assert max_priority_fee >= 1 * eth.GWEI


def test_send_ether(ctx):
    operator = ctx.operator

    test_account_addr = ctx.create_and_fund_eth_account()

    operator_balance_before = ctx.eth.get_eth_balance(operator)
    tmp_balance_before = ctx.eth.get_eth_balance(test_account_addr)

    logging.info(f"Operator balance before: {operator_balance_before}")
    logging.info(f"Test account balance before: {tmp_balance_before}")

    amount_to_transfer = 123456 * eth.GWEI

    assert ctx.eth.get_eth_balance(operator) > amount_to_transfer + 2 * max_tx_cost

    def send(src_addr, dst_addr, amount):
        src_balance_before = ctx.eth.get_eth_balance(src_addr)
        dst_balance_before = ctx.eth.get_eth_balance(dst_addr)
        assert src_balance_before >= amount + max_tx_cost
        txrcpt = ctx.eth.send_eth(src_addr, dst_addr, amount)
        src_balance_after = ctx.eth.get_eth_balance(src_addr)
        dst_balance_after = ctx.eth.get_eth_balance(dst_addr)
        src_balance_change = src_balance_after - src_balance_before
        dst_balance_change = dst_balance_after - dst_balance_before
        effective_fee = - src_balance_change - amount
        assert src_balance_change < - amount
        assert dst_balance_change == amount
        assert txrcpt.gasUsed >= min_tx_gas
        assert txrcpt.gasUsed <= max_tx_gas
        # Hardhat and Alchemy support effectiveGasPrice, ganache does not
        assert ("effectiveGasPrice" in txrcpt) == (not ctx.eth.is_local_node) or on_peggy2_branch
        if ctx.eth.is_local_node:
            effective_gas_price = ctx.eth.w3_conn.eth.gas_price
        else:
            effective_gas_price = txrcpt.effectiveGasPrice
        assert effective_fee == txrcpt.gasUsed * effective_gas_price
        assert effective_fee >= min_tx_cost
        assert effective_fee <= max_tx_cost
        assert effective_gas_price >= min_gas_price
        assert effective_gas_price <= max_gas_price
        logging.info(f"Effective transaction fee: {effective_fee}")

    send(operator, test_account_addr, amount_to_transfer + max_tx_cost)
    send(test_account_addr, operator, amount_to_transfer)


def test_deploy_erc20_token(ctx):
    operator = ctx.operator

    # TODO Cannot assume that initial balances will be 0 unless we're running in a snapshot

    assert ctx.eth.get_eth_balance(operator) > max_tx_cost

    # Deploy a sample token
    token = ctx.generate_random_erc20_token_data()
    token_sc = ctx.deploy_new_generic_erc20_token(token.name, token.symbol, token.decimals)
    token_addr = token_sc.address

    # Create test account and fund it if necessary
    test_account_addr = ctx.create_and_fund_eth_account(fund_amount=max_tx_cost)

    # As token owner, mint some tokens to test_account_addr and verify the balance
    assert ctx.get_erc20_token_balance(token_addr, test_account_addr) == 0
    assert ctx.get_erc20_token_balance(token_addr, operator) == 0
    mint_amount = 123
    ctx.eth.transact_sync(token_sc.functions.mint, operator)(test_account_addr, mint_amount)
    assert ctx.get_erc20_token_balance(token_addr, test_account_addr) == mint_amount
    assert ctx.get_erc20_token_balance(token_addr, operator) == 0

    # mint() is "owner only" and can only be done by account that created the smart contract.
    # Calling it with another account should fail with a specific exception and message.
    try:
        ctx.eth.transact_sync(token_sc.functions.mint, test_account_addr)(test_account_addr, mint_amount)
        assert False
    except Exception as e:
        assert ctx.eth.is_contract_logic_error_not_in_minter_role(e)

    # Create another account
    another_account_addr = ctx.create_and_fund_eth_account()

    # Send some tokens to another account
    send_amount = 3
    assert (send_amount > 0) and (send_amount <= mint_amount)
    assert ctx.get_erc20_token_balance(token_addr, test_account_addr) == mint_amount
    assert ctx.get_erc20_token_balance(token_addr, another_account_addr) == 0
    assert ctx.get_erc20_token_balance(token_addr, operator) == 0
    ctx.send_erc20_tokens(token_sc.address, test_account_addr, another_account_addr, send_amount)
    assert ctx.get_erc20_token_balance(token_addr, test_account_addr) == mint_amount - send_amount
    assert ctx.get_erc20_token_balance(token_addr, another_account_addr) == send_amount
    assert ctx.get_erc20_token_balance(token_addr, operator) == 0

    # Check that our balance function returns same values as balanceOf() from smart contract
    for addr in [test_account_addr, another_account_addr, operator]:
        assert ctx.get_erc20_token_balance(token_addr, addr) == token_sc.functions.balanceOf(addr).call()

    # Check past events
    past_transfers = ctx.smart_contract_get_past_events(token_sc, "Transfer")
    assert len(past_transfers) == 2
    assert past_transfers[0].address == token_addr
    assert past_transfers[0].args["from"] == eth.NULL_ADDRESS
    assert past_transfers[0].args["to"] == test_account_addr
    assert past_transfers[0].args["value"] == mint_amount
    assert past_transfers[1].address == token_addr
    assert past_transfers[1].args["from"] == test_account_addr
    assert past_transfers[1].args["to"] == another_account_addr
    assert past_transfers[1].args["value"] == send_amount

    # Try to transfer one token too many (over balance). This should fail and the balance should not change.
    try:
        ctx.send_erc20_tokens(token_addr, test_account_addr, another_account_addr, mint_amount - send_amount + 1)
        assert False
    except Exception as e:
        assert ctx.eth.is_contract_logic_error_amount_exceeds_balance(e)

    assert ctx.get_erc20_token_balance(token_sc.address, test_account_addr) == mint_amount - send_amount
    assert ctx.get_erc20_token_balance(token_sc.address, another_account_addr) == send_amount
    assert ctx.get_erc20_token_balance(token_sc.address, operator) == 0
