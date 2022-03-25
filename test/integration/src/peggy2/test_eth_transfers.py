import pytest
import siftool_path

from siftool import eth, test_utils, sifchain
from siftool.common import *
from siftool.test_utils import EnvCtx
from typing import Iterable

# TODO for PR: Remove all print outs

fund_amount_eth = 10 * eth.ETH
fund_amount_sif = 10 * test_utils.sifnode_funds_for_transfer_peggy1  # TODO How much rowan do we need? (this is 10**18)
fund_amount_ceth_cross_chain_fee = 10 * eth.GWEI

def test_eth_to_ceth_and_back_to_eth_transfer_valid(ctx):
    # Create/retrieve a test ethereum account
    test_eth_account = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)

    # create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])

    # Verify initial balance
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)

    # Send from ethereum to sifchain by locking
    amount_to_send = 123456 * eth.GWEI
    assert amount_to_send < fund_amount_eth

    ctx.bridge_bank_lock_eth(test_eth_account, test_sif_account, amount_to_send)
    ctx.advance_blocks()

    # Verify final balance
    test_sif_account_final_balance = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    balance_diff = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    assert exactly_one(list(balance_diff.keys())) == ctx.ceth_symbol
    assert balance_diff[ctx.ceth_symbol] == amount_to_send

    # Send from sifchain to ethereum by burning on sifchain side,
    # > sifnoded tx ethbridge burn
    # Reduce amount for cross-chain fee. The same formula is used inside this function.
    eth_balance_before = ctx.eth.get_eth_balance(test_eth_account)
    amount_to_send = amount_to_send - ctx.eth.cross_chain_fee_base * ctx.eth.cross_chain_burn_fee
    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_account, amount_to_send, ctx.ceth_symbol)

    # Verify final balance
    ctx.wait_for_eth_balance_change(test_eth_account, eth_balance_before)


def test_erc20_to_sifnode_and_back_first_time(ctx):
    token_decimals = 18
    token_sc = deploy_erc20_token_for_test(ctx, token_decimals)
    transfer_erc20_to_sifnode_and_back(ctx, token_sc, token_decimals, 1, False)


# def test_erc20_to_sifnode_and_back_multiple_times(ctx):
#     token_decimals = 18
#     token_sc = deploy_erc20_token_for_test(ctx, token_decimals)
#     transfer_erc20_to_sifnode_and_back(ctx, token_sc, token_decimals, 5, False)


# def test_trolltoken_to_sifnode_and_back(ctx):
#     token_decimals = 0
#     token_sc = deploy_trolltoken_for_test(ctx)
    # TODO Depending on timeout is not reliable confirmation of test passing. Timeout can happen for many other reasons
    #      such as slow machine, Ethereum chain congestion, relayer crashing etc. It would be much better if we could
    #      use some specific information to confirm that we did in fact handle the token correctly (i.e. seeing a
    #      particular message in logs, or a specific event being emitted by BridgeBank etc.)
    # transfer_erc20_to_sifnode_and_back(ctx, token_sc, token_decimals, 1, True)


# Shared code used by different test scenarios
def transfer_erc20_to_sifnode_and_back(ctx: EnvCtx, token_sc, token_decimals, number_of_times, is_troll_token):
    # Create/retrieve 2 test ethereum accounts
    test_eth_acct_0, test_eth_acct_1 = [ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth) for _ in range(2)]

    # Create/retrieve a test sifchain account
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"]])

    # We must fund test_sif_acct with some ceth to pay for transaction
    amount_to_send = 5000000 * eth.GWEI * number_of_times # TODO How to set properly?
    assert amount_to_send < fund_amount_eth
    test_sif_account_initial_balance = ctx.get_sifchain_balance(test_sif_account)
    ctx.bridge_bank_lock_eth(test_eth_acct_0, test_sif_account, amount_to_send)
    ctx.advance_blocks()
    test_sif_account_final_balance = ctx.wait_for_sif_balance_change(test_sif_account, test_sif_account_initial_balance)
    sif_balance_delta = sifchain.balance_delta(test_sif_account_initial_balance, test_sif_account_final_balance)
    assert len(sif_balance_delta) == 1
    assert sif_balance_delta[ctx.ceth_symbol] == amount_to_send

    token_addr = token_sc.address
    sif_denom_hash = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, token_addr)

    send_amount_0 = 10 * 10**token_decimals
    send_amount_1 = 3 * 10**token_decimals
    assert send_amount_1 < send_amount_0
    total_amount = send_amount_0 * number_of_times

    # We do minting and approving just once for all iterations, but we could also do it each time separately.
    ctx.mint_generic_erc20_token(token_addr, total_amount, test_eth_acct_0)

    for i in range(number_of_times):
        # Send from Ethereum account 1 to Sifchain
        eth_balance_before_0 = ctx.get_erc20_token_balance(token_addr, test_eth_acct_0)
        sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
        ctx.send_from_ethereum_to_sifchain(test_eth_acct_0, test_sif_account, send_amount_0, token_sc=token_sc)
        ctx.advance_blocks()
        sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)
        eth_balance_after_0 = ctx.get_erc20_token_balance(token_addr, test_eth_acct_0)
        sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)

        assert len(sif_balance_delta) == 1
        assert sif_balance_delta[sif_denom_hash] == send_amount_0
        assert eth_balance_before_0 == total_amount - send_amount_0 * i
        assert eth_balance_after_0 == eth_balance_before_0 - send_amount_0

        # test_eth_account2 = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
        eth_balance_before_1 = ctx.get_erc20_token_balance(token_addr, test_eth_acct_1)
        sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
        ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_acct_1, send_amount_1, sif_denom_hash)

        # TrollToken should time out, any legit ERC20 should pass.
        # Timeout needs to be long enough for any legit token (90s works for Hardhat, but might not work for Ropsten).
        try:
            eth_balance_after_1 = ctx.wait_for_eth_balance_change(test_eth_acct_1, eth_balance_before_1,
                token_addr=token_addr, timeout=90)
            assert not is_troll_token
        except Exception as e:
            assert is_troll_token
            assert i + 1 == number_of_times == 1
            assert e.__class__ == Exception
            assert len(e.args) == 1
            assert e.args[0] == "Timeout waiting for Ethereum balance to change"
            return

        sif_balance_after = ctx.get_sifchain_balance(test_sif_account)
        sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)

        assert sif_balance_delta[sif_denom_hash] == - send_amount_1
        assert sif_balance_delta["rowan"] == -100000  # TODO Where is this value from?
        assert sif_balance_delta[ctx.ceth_symbol] == -1  # TODO Where is this value from?
        assert eth_balance_after_1 == eth_balance_before_1 + send_amount_1
        assert eth_balance_after_1 == send_amount_1 * (i + 1)


# Lock an eth to
# We expect the transfer from evm to sifchain to be successful,
# We expect the transfer from sifchain to evm to fail, w/
#   - Tokens burned on sifchain side
#   - Tokens not depsoited on evm side
#   - It SHOULD NOT halt the bridge
def test_failhard_token_to_sifnode_and_back_does_not_halt_bridge(ctx: EnvCtx):
    test_eth_acct = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"],
                                                              [fund_amount_ceth_cross_chain_fee, ctx.ceth_symbol]])
    test_account_token_balance = 30000

    token_sc = deploy_failhard_for_test(ctx, test_eth_acct, test_account_token_balance)
    token_addr = token_sc.address
    sif_denom_hash = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, token_sc.address)

    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
    eth_token_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)
    # Locking erc20 token to sifchain
    ctx.send_from_ethereum_to_sifchain(test_eth_acct, test_sif_account, test_account_token_balance, token_sc)
    ctx.advance_blocks()

    # Group these into 1 func
    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 1, "User should only have changes in token balance. Received {}".format(sif_balance_delta)
    assert sif_denom_hash in sif_balance_delta, "User should see changes in the bridged token"
    assert sif_balance_delta[sif_denom_hash] == test_account_token_balance

    eth_token_balance_after = ctx.wait_for_eth_balance_change(test_eth_acct, eth_token_balance_before, token_addr=token_addr)
    eth_token_balance_delta = eth_token_balance_after - eth_token_balance_before
    assert eth_token_balance_delta == (test_account_token_balance * -1), "User's token on evm should have decreased by sent amount"

    # The user has successfully locked token on evm, and got balance on sifchain

    # Attempting to send arbitrary amount back to evm
    test_send_amount_back = test_account_token_balance - 15

    eth_token_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)
    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)

    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_acct, test_send_amount_back, sif_denom_hash)

    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before, min_changes=[[1, "rowan"], [1, ctx.ceth_symbol], [1, sif_denom_hash]])
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 3, "User should only have changes in token balance. Delta: {}".format(sif_balance_delta)
    assert "rowan" in sif_balance_delta, "User should see change in rowan balance"
    assert sif_balance_delta["rowan"] < 0, "User should see rowan balance decreased to pay for mining"
    assert ctx.ceth_symbol in sif_balance_delta, "User should see changes in the bridged token"
    assert sif_balance_delta[ctx.ceth_symbol] < 0, "User should see ceth balance decreased to pay for cross chain"
    assert sif_denom_hash in sif_balance_delta, "User should see change in token balance"
    assert sif_balance_delta[sif_denom_hash] == -1 * test_send_amount_back, "User's token should've been burned regardless of evm tx status"

    with pytest.raises(Exception) as exception:
        ctx.wait_for_eth_balance_change(test_eth_acct, eth_token_balance_before, token_addr=token_addr, timeout=90)
        assert exception.args[0] == "Timeout waiting for Ethereum balance to change", "We shouldn't see any changes on evm side, this was supposed to fail"

    test_erc20_to_sifnode_and_back_first_time(ctx)

# Token with unicode in symbol should not have abnormal behavior in import & export flow
# We expect the transfer from evm to sifchain to be successful,
# We expect the transfer from sifchain to evm to be successful
def test_unicodeToken_token_to_sifnode_and_back_succeed_and_does_not_halt_bridge(ctx: EnvCtx):
    test_eth_acct = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"],
                                                              [fund_amount_ceth_cross_chain_fee, ctx.ceth_symbol]])
    test_account_token_balance = 30000

    token_sc = deploy_unicodeToken_for_test(ctx, test_eth_acct, test_account_token_balance)
    token_addr = token_sc.address
    sif_denom_hash = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, token_sc.address)

    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
    eth_token_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)

    ctx.send_from_ethereum_to_sifchain(test_eth_acct, test_sif_account, test_account_token_balance, token_sc)
    ctx.advance_blocks()

    # Group these into 1 func
    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 1, "User should only have changes in token balance. Received {}".format(sif_balance_delta)
    assert sif_denom_hash in sif_balance_delta, "User should see changes in the bridged token"
    assert sif_balance_delta[sif_denom_hash] == test_account_token_balance

    eth_token_balance_after = ctx.wait_for_eth_balance_change(test_eth_acct, eth_token_balance_before, token_addr=token_addr)
    eth_token_balance_delta = eth_token_balance_after - eth_token_balance_before
    assert eth_token_balance_delta == (test_account_token_balance * -1), "User's token on evm should have decreased by sent amount"

    # The user has successfully locked token on evm, and got balance on sifchain

    # Attempting to send arbitrary amount back to evm
    test_send_amount_back = test_account_token_balance - 15

    eth_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)
    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)

    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_acct, test_send_amount_back, sif_denom_hash)
    ctx.advance_blocks()

    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before, min_changes=[[1, "rowan"], [1, ctx.ceth_symbol], [1, sif_denom_hash]])
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 3, "User should only have changes in token balance. Delta: {}".format(sif_balance_delta)
    assert "rowan" in sif_balance_delta, "User should see change in rowan balance"
    assert sif_balance_delta["rowan"] < 0, "User should see rowan balance decreased to pay for mining"
    assert ctx.ceth_symbol in sif_balance_delta, "User should see changes in the bridged token"
    assert sif_balance_delta[ctx.ceth_symbol] < 0, "User should see ceth balance decreased to pay for cross chain"
    assert sif_denom_hash in sif_balance_delta, "User should see change in token balance"
    assert sif_balance_delta[sif_denom_hash] == -1 * test_send_amount_back, "User's token should've been burned regardless of evm tx status"

    eth_balance_after = ctx.wait_for_eth_balance_change(test_eth_acct, eth_balance_before,
        token_addr=token_addr, timeout=90)
    eth_balance_delta = eth_balance_after - eth_balance_before
    assert eth_balance_delta == test_send_amount_back

    test_erc20_to_sifnode_and_back_first_time(ctx)

# Token with built-in fee deduction in _transfer() should not have abnormal behavior in import & export flow
# We expect the transfer from evm to sifchain to be successful,
#   - User should receive net balance reflecting fee deduction in contract's _transfer function
# We expect the transfer from sifchain to evm to be successful
def test_commission_token_to_sifnode_and_back_succeed_and_does_not_halt_bridge(ctx: EnvCtx):
    test_eth_acct, commission_dev_acct = [ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth) for _ in range(2)]

    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"],
                                                              [fund_amount_ceth_cross_chain_fee, ctx.ceth_symbol]])
    test_account_token_balance = 30000

    token_dev_fee = 10
    token_sc = deploy_commissiontoken_for_test(ctx, commission_dev_acct, token_dev_fee, test_eth_acct, test_account_token_balance)
    token_addr = token_sc.address
    sif_denom_hash = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, token_sc.address)

    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
    eth_token_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)
    assert eth_token_balance_before == test_account_token_balance

    # Locking erc20 token to sifchain
    ctx.send_from_ethereum_to_sifchain(test_eth_acct, test_sif_account, test_account_token_balance, token_sc)
    ctx.advance_blocks()

    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 1, "User should only have changes in token balance. Received {}".format(sif_balance_delta)
    assert sif_denom_hash in sif_balance_delta, "User should see changes in the bridged token"
    assert sif_balance_delta[sif_denom_hash] < test_account_token_balance, "User sifchain should have received less than transfered amount, due to commission on erc20 transfer"

    eth_token_balance_after = ctx.wait_for_eth_balance_change(test_eth_acct, eth_token_balance_before, token_addr=token_addr)
    eth_token_balance_delta = eth_token_balance_after - eth_token_balance_before
    assert eth_token_balance_delta == (test_account_token_balance * -1), "User's token on evm should have decreased by sent amount"

    # Arbitrary amount to be sent back
    test_send_amount_back = 20000

    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
    eth_token_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)
    ctx.sifnode_client.send_from_sifchain_to_ethereum(test_sif_account, test_eth_acct, test_send_amount_back, sif_denom_hash)

    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before, min_changes=[[1, "rowan"], [1, ctx.ceth_symbol], [1, sif_denom_hash]])
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 3, "User should only have changes in token balance. Delta: {}".format(sif_balance_delta)
    assert "rowan" in sif_balance_delta, "User should see rowan decreased for cross chain fee"
    assert sif_balance_delta["rowan"] < 0, "User should have rowan decreased to pay for fee"
    assert ctx.ceth_symbol in sif_balance_delta, "User should see changes in the bridged token"
    assert sif_balance_delta[ctx.ceth_symbol] < 0, "User should have ceth decreased to pay for fee"
    assert sif_denom_hash in sif_balance_delta
    assert sif_balance_delta[sif_denom_hash] == -1 * test_send_amount_back, "User's token should've been burned regardless of evm tx status"

    eth_token_balance_after = ctx.wait_for_eth_balance_change(test_eth_acct, old_balance=eth_token_balance_before, token_addr=token_addr)
    eth_token_balance_delta = eth_token_balance_after - eth_token_balance_before
    assert eth_token_balance_delta > 0, "Should've received non-zero amount"
    assert eth_token_balance_delta < test_send_amount_back, "Should've received less amount than indicated due to token's fee"

    test_erc20_to_sifnode_and_back_first_time(ctx)

# Token with random return values should not halt the bridge
def test_randomtroll_token_to_sifnode_does_not_halt_bridge(ctx: EnvCtx):
    test_eth_acct = ctx.create_and_fund_eth_account(fund_amount=fund_amount_eth)
    test_sif_account = ctx.create_sifchain_addr(fund_amounts=[[fund_amount_sif, "rowan"],
                                                              [fund_amount_ceth_cross_chain_fee, ctx.ceth_symbol]])
    test_account_token_balance = 30000

    token_sc = deploy_randomtrolltoken_for_test(ctx, [test_eth_acct], [test_account_token_balance])
    token_addr = token_sc.address
    sif_denom_hash = sifchain.sifchain_denom_hash(ctx.eth.ethereum_network_descriptor, token_sc.address)

    sif_balance_before = ctx.get_sifchain_balance(test_sif_account)
    eth_token_balance_before = ctx.get_erc20_token_balance(token_addr, test_eth_acct)

    ctx.send_from_ethereum_to_sifchain(test_eth_acct, test_sif_account, test_account_token_balance, token_sc)
    ctx.advance_blocks()

    sif_balance_after = ctx.wait_for_sif_balance_change(test_sif_account, sif_balance_before)
    sif_balance_delta = sifchain.balance_delta(sif_balance_before, sif_balance_after)
    assert len(sif_balance_delta) == 1, "User should only have changes in token balance. Received {}".format(sif_balance_delta)
    assert sif_denom_hash in sif_balance_delta, "User should see changes in the bridged token"

    test_erc20_to_sifnode_and_back_first_time(ctx)









# TODO: Token_data also has token decimals field, why are we taking as param
def deploy_erc20_token_for_test(ctx, token_decimals):
    token_data: test_utils.ERC20TokenData = ctx.generate_random_erc20_token_data()
    return ctx.deploy_new_generic_erc20_token(token_data.name, token_data.symbol, token_decimals)

def deploy_trolltoken_for_test(ctx: EnvCtx):
    token = ctx.generate_random_erc20_token_data()
    abi, bytecode, _ = ctx.abi_provider.get_descriptor("TrollToken")
    token_sc = ctx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
    txrcpt = ctx.eth.transact_sync(token_sc.constructor, ctx.operator)(token.name, token.symbol)
    return ctx.w3_conn.eth.contract(abi=abi, address=txrcpt.contractAddress)

def deploy_failhard_for_test(ctx: EnvCtx, account: str, amount: int):
    token = ctx.generate_random_erc20_token_data()
    abi, bytecode, _ = ctx.abi_provider.get_descriptor("FailHardToken")
    token_sc = ctx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
    txrcpt = ctx.eth.transact_sync(token_sc.constructor, ctx.operator)("fail hard token", "FHT", account, amount)
    return ctx.w3_conn.eth.contract(abi=abi, address=txrcpt.contractAddress)

def deploy_unicodeToken_for_test(ctx: EnvCtx, account: str, amount: int):
    token = ctx.generate_random_erc20_token_data()
    abi, bytecode, _ = ctx.abi_provider.get_descriptor("UnicodeToken")
    token_sc = ctx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
    txrcpt = ctx.eth.transact_sync(token_sc.constructor, ctx.operator)()
    token_sc = ctx.w3_conn.eth.contract(abi=abi, address=txrcpt.contractAddress)

    ctx.mint_generic_erc20_token(token_sc.address, amount, account)
    return token_sc

def deploy_commissiontoken_for_test(ctx, dev: str, devFee: int, user: str, quantity: int):
    token = ctx.generate_random_erc20_token_data()
    abi, bytecode, _ = ctx.abi_provider.get_descriptor("CommissionToken")
    token_sc = ctx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
    txrcpt = ctx.eth.transact_sync(token_sc.constructor, ctx.operator)(dev, devFee, user, quantity)
    return ctx.w3_conn.eth.contract(abi=abi, address=txrcpt.contractAddress)

def deploy_randomtrolltoken_for_test(ctx: EnvCtx, initialAccounts: Iterable[str], quantity: Iterable[int]):
    token = ctx.generate_random_erc20_token_data()
    abi, bytecode, _ = ctx.abi_provider.get_descriptor("RandomTrollToken")
    token_sc = ctx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
    txrcpt = ctx.eth.transact_sync(token_sc.constructor, ctx.operator)(initialAccounts, quantity)
    return ctx.w3_conn.eth.contract(abi=abi, address=txrcpt.contractAddress)
