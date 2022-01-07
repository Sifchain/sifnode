import pytest

from integration_framework import main, common, eth, test_utils, inflate_tokens
from common import *


max_gas_required = 200000
max_gas_price = 500 * eth.GWEI


def deploy_new_erc20_token_for_testing(ctx):
    # Symbol must be unique on the blocklist
    token = ctx.generate_random_erc20_token_data()
    token.decimals = 0
    mint_amount = 1000 * 10**token.decimals
    token_sc = ctx.deploy_new_generic_erc20_token(token.name, token.symbol, token.decimals, mint_amount=mint_amount)
    return token_sc

def bridge_bank_lock_eth(ctx, from_eth_acct, to_sif_acct, amount):
    assert ctx.eth.get_eth_balance(from_eth_acct) > max_gas_required * max_gas_price, "Not enough gas for test"
    return ctx.bridge_bank_lock_eth(from_eth_acct, to_sif_acct, amount)

def bridge_bank_lock_erc20(ctx, bridge_token, from_eth_acct, to_sif_acct, amount):
    assert ctx.eth.get_eth_balance(from_eth_acct) > max_gas_required * max_gas_price, "Not enough gas for test"
    assert ctx.get_erc20_token_balance(bridge_token.address, from_eth_acct) >= amount, "Not enough tokens for test"
    return ctx.bridge_bank_lock_erc20(bridge_token.address, from_eth_acct, to_sif_acct, amount)

def is_blocklisted_exception(ctx, exception):
    return ctx.eth.is_contract_logic_error(exception, "Address is blocklisted")

@pytest.mark.skipif("on_peggy2_branch")
def test_blocklist_eth():
    with test_utils.get_test_env_ctx() as ctx:
        _test_blocklist_eth(ctx)

def _test_blocklist_eth(ctx):
    w3 = ctx.eth.w3_conn

    amount_to_fund = 1 * eth.ETH
    amount_to_send = eth.ETH // 100

    acct1, acct2 = [ctx.create_and_fund_eth_account(fund_amount=amount_to_fund) for _ in range(2)]

    to_sif_acct = ctx.create_sifchain_addr()
    sif_symbol = test_utils.CETH

    bridge_bank = ctx.get_bridge_bank_sc()

    filter = bridge_bank.events.LogLock.createFilter(fromBlock="latest")

    # Valid negative test outcome: transaction rejected with the string "Address is blocklisted"
    def assert_blocked(addr):
        assert len(filter.get_new_entries()) == 0

        try:
            bridge_bank_lock_eth(ctx, addr, to_sif_acct, amount_to_send)
            assert False
        except Exception as e:
            assert is_blocklisted_exception(ctx, e)

        assert len(filter.get_new_entries()) == 0

    # Valid positive test outcome: event emitted, optionally: funds appear on sifchain
    def assert_not_blocked(addr):
        assert len(filter.get_new_entries()) == 0

        balances_before = ctx.get_sifchain_balance(to_sif_acct)
        txrcpt = bridge_bank_lock_eth(ctx, addr, to_sif_acct, amount_to_send)
        ctx.advance_blocks()
        balances_after = ctx.wait_for_sif_balance_change(to_sif_acct, balances_before)

        assert balances_after.get(sif_symbol, 0) == balances_before.get(sif_symbol, 0) + amount_to_send

        entries = filter.get_new_entries()
        assert len(entries) == 1
        assert entries[0].event == "LogLock"
        assert entries[0].transactionHash == txrcpt.transactionHash
        assert entries[0].address == bridge_bank.address
        assert entries[0].args["_from"] == addr
        assert entries[0].args["_to"] == test_utils.sif_addr_to_evm_arg(to_sif_acct)
        assert entries[0].args["_value"] == amount_to_send

    try:
        assert_not_blocked(acct1)
        assert_not_blocked(acct2)
        ctx.set_ofac_blocklist_to([acct2])
        assert_not_blocked(acct1)
        assert_blocked(acct2)
        ctx.set_ofac_blocklist_to([])
        assert_not_blocked(acct1)
        assert_not_blocked(acct2)
    finally:
        w3.eth.uninstall_filter(filter.filter_id)


@pytest.mark.skipif("on_peggy2_branch")
def test_blocklist_erc20():
    with test_utils.get_test_env_ctx() as ctx:
        _test_blocklist_erc20(ctx)

# For ERC20 tokens, we need to create a new instance of Blocklist smart contract, deploy it and whitelist it with
# BridgeBank. In peggy1, the token matching in BridgeBank is done by symbol, so we need to give our token a unique
# symbol such as TEST or MOCK + random suffix + call updateEthWtiteList() + mint() + approve().
# See smart-contracts/test/test_bridgeBank.js:131-160 for example.
def _test_blocklist_erc20(ctx):
    w3 = ctx.eth.w3_conn

    test_token = deploy_new_erc20_token_for_testing(ctx)

    eth_token_symbol = test_token.functions.symbol().call()
    sif_token_symbol = ctx.eth_symbol_to_sif_symbol(eth_token_symbol)

    bridge_bank = ctx.get_bridge_bank_sc()

    amount_to_fund = 1 * eth.ETH
    amount_to_send = 1

    acct1, acct2 = [ctx.create_and_fund_eth_account(fund_amount=amount_to_fund) for _ in range(2)]

    for acct in [acct1, acct2]:
        # Transfer 10 tokens from operator to acct
        # TODO This would be better done as ctx.send_erc20_tokens(), but we're currently using BridgeToken
        ctx.eth.transact_sync(test_token.functions.transfer, ctx.operator)(acct, 10)

        # Set allowance for BridgeBank to spend 10 tokens on behalf of acct1 and acct2
        # Without this we get transaction rejected with "SafeERC20: low-level call failed"
        # TODO Move to Peggy1EnvCtx.bridge_bank_lock_erc20() as they should always be together
        ctx.eth.transact_sync(test_token.functions.approve, acct)(bridge_bank.address, 10)

    to_sif_acct = ctx.create_sifchain_addr()

    filter = bridge_bank.events.LogLock.createFilter(fromBlock="latest")

    def assert_blocked(addr):
        assert len(filter.get_new_entries()) == 0

        try:
            bridge_bank_lock_erc20(ctx, test_token, addr, to_sif_acct, amount_to_send)
            assert False
        except Exception as e:
            assert is_blocklisted_exception(ctx, e)

        assert len(filter.get_new_entries()) == 0

    def assert_not_blocked(addr):
        assert len(filter.get_new_entries()) == 0

        balances_before = ctx.get_sifchain_balance(to_sif_acct)
        txrcpt = bridge_bank_lock_erc20(ctx, test_token, addr, to_sif_acct, amount_to_send)
        ctx.advance_blocks()
        balances_after = ctx.wait_for_sif_balance_change(to_sif_acct, balances_before)

        assert balances_after.get(sif_token_symbol, 0) == balances_before.get(sif_token_symbol, 0) + amount_to_send

        entries = filter.get_new_entries()
        assert len(entries) == 1
        assert entries[0].event == "LogLock"
        assert entries[0].transactionHash == txrcpt.transactionHash
        assert entries[0].address == bridge_bank.address
        assert entries[0].args["_from"] == addr
        assert entries[0].args["_to"] == test_utils.sif_addr_to_evm_arg(to_sif_acct)
        assert entries[0].args["_value"] == amount_to_send

    try:
        assert_not_blocked(acct1)
        assert_not_blocked(acct2)
        ctx.set_ofac_blocklist_to([acct2])
        assert_not_blocked(acct1)
        assert_blocked(acct2)
        ctx.set_ofac_blocklist_to([])
        assert_not_blocked(acct1)
        assert_not_blocked(acct2)
    finally:
        w3.eth.uninstall_filter(filter.filter_id)
