# This is for load testing of LPD/rewardss (and in future, margin)
#
# Scenario description: https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5
# Ticket: https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifnode/3020
# How to run a validator in multi-node setup:
# - https://docs.sifchain.finance/network-security/validators/running-sifnode-and-becoming-a-validator
# - https://docs.sifchain.finance/developers/tutorials/setup-standalone-validator-node-manually
#
# Requirements: Python (3.9 is best, 3.8 and 3.10 should also work, support for other versions is currently unknownn)
#
# Example usage:
# cd test/load
# ../integration/framework/siftool venv
# ../integration/framework/venv/bin/python3 test_many_pools_and_liquidity_providers.py \
#     --number-of-liquidity-pools 5 \
#     --number-of-wallets 10 \
#     --liquidity-providers-per-wallet 3 \
#     --reward-period-default-multiplier 0.0 \
#     --reward-period-distribute \
#     --reward-period-mod 1 \
#     --reward-period-pool-count 5 \
#     --test-duration-blocks 5 \
#     --number-of-nodes 4
#
# The test runs, in sequence:
# 1. test_duration_blocks of neither rewards nor lppd
# 2. test_duration_blocks of rewards without lppd
# 3. test_duration_blocks of rewards and lppd
# 4. test_duration_blocks of lppd without rewards
#
# |<---test time------------------------------------------------------------------->|
#
#                          |<--------------rewards-------------->|
#                                             ]<---------------lppd---------------->|
#
#       |<----neither----->|<--rewards only-->|<--rewards+lppd-->|<----lppd only--->|
#       ^-- time0          ^-- time1          ^-- time2          ^-- time3          ^-- time4
#
#
#
# TODO - improvements
#
# (1) Log HTTP response for block_height, maybe we can get a bit more information like this:
# $ curl -i "http://rpc-archive.sifchain.finance/block_results?height=9000000"
# HTTP/1.1 500 Internal Server Error
# ...
# {
#   "jsonrpc": "2.0",
#   "id": -1,
#   "error": {
#     "code": -32603,
#     "message": "Internal error",
#     "data": "height 9000000 must be less than or equal to the current blockchain height 7856892"
#   }
# }
#
# (2) use rocksb, see https://raw.githubusercontent.com/Sifchain/sifchain-devops/1218ff79b22ab2a6bd22b81d6aa4385a247cafc9/scripts/sifnode/testing/sifnode_n_node_network_simulator.py?token=GHSAT0AAAAAABLH7LII5AAWD6YDWG7THBHGYW7DSVA
#
# (3) Exceptions / printing of _debug...
#
# (4) Use parameter rpc.laddr for sifnoded start instead of self.node


import argparse
import json
import random
import sys
import time
import siftool_path
from siftool.common import *
from siftool import command, sifchain, project, cosmos, environments, test_utils2
from siftool.sifchain import ROWAN, STAKE, ROWAN_DECIMALS


log = siftool_logger()


class Test:
    def __init__(self, cmd: command.Command, prj: project.Project, sifnoded_home_root: str):
        self.cmd = cmd
        self.prj = prj
        self.rnd = random.Random(5646067977921730044)  # Use a fixed random seed for repeatable results

        self.number_of_nodes = None
        # The number of liquidity pools to which each wallet provides liquidity. The pools are chosen randomly from
        # all `number_of_liquidity_pools`. This is also the same of number of different tokens per wallet (not counting
        # rowan). Number of pools == number of tokens. We create this many tokens, and create a pool for each token.
        self.number_of_liquidity_pools = None
        # Number of wallets. Since each wallet provides liquidity to 1 or more liquidity pool, this is also the number
        # of unique liquidity providers.
        self.number_of_wallets = None
        self.liquidity_providers_per_wallet = None
        self.reward_period_default_multiplier = None
        self.reward_period_distribute = None
        self.reward_period_mod = None
        self.reward_period_pool_count = None
        self.rewards_offset_blocks = None
        self.rewards_duration_blocks = None
        self.lpd_period_mod = None
        self.lpd_offset_blocks = None
        self.lpd_duration_blocks = None
        self.block_results_offset = None
        self.run_forever = False
        self.disable_assertions = False
        self.wallets_dir = None

        # The timing starts with the next block after setup. The accuracty of the test is limited by polling for the
        # current block number (1s). The total time will be 4 * test_duration_blocks * block_time, i.e.
        # 4 * 6s = 24s for one unit of test_duration_blocks.

        self.sifnoded_home_root = sifnoded_home_root

        self.env = None
        self.custom_wallet_mnemonics = [
            # sif1rruvw03utshn7ry3emeqf2gzkg6eap6hu5shun
            "zebra sentence tape you spawn forget catalog veteran rocket steel ticket slender follow rubber spoil thing into liar twin document ring clock shell skirt",
        ]

        self.sifnoded = None
        self.sifnoded_client = None

    def setup(self):
        # Define one token per liquidity pool.
        self.token_decimals = 18
        self.token_unit = 10**self.token_decimals
        self.tokens = ["test{}".format(i) for i in range(self.number_of_liquidity_pools)]

        # We are only dealing with symmetrical liquidity pools here.
        # This means that each liquidity provider uses `native_amount == external_amount`.
        # The ratio `native_amount/external_amount` has to be the same as defined per liquidity pool
        # (withing a certain threshold).
        # We use 1000/1000 for pools and 500/500 for liquidity providers.
        self.amount_of_denom_per_wallet = 1000 * self.token_unit
        self.amount_of_liquidity_added_by_wallet = 500 * self.token_unit

        self.amount_of_rowan_per_wallet = 10000 * 10**18  # TODO How much?

        self.validator0_mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow".split()

        self.prj.pkill()
        time.sleep(1)

        if self.wallets_dir is not None:
            predefined_wallets = test_utils2.PredefinedWallets(self.cmd, self.wallets_dir)
            client_home = self.wallets_dir
        else:
            predefined_wallets = None
            client_home = os.path.join(self.sifnoded_home_root, "sifnoded-client")

        # Set up test wallets with test tokens. We do this in genesis for performance reasons. For each wallet we choose
        # number_of_denoms_per_wallet` random denoms.
        # Set up admin account balances. We add these with "add-genesis-account"
        # TODO It is not clear if we really need to fund all of them (and how much).
        # TODO Does this have to cover for rewards and lppd distribution? If rewards are minted, then no.
        # For rewards, the funds are minted and in case we opted for a distribution of the rewards to the LP wallet the
        # minted rowans are transferred there, you can see the minting process here: https://github.com/Sifchain/sifnode/blob/master/x/clp/keeper/rewards.go#L54
        # For LPD, we only transfer the existing funds in CLP to the LP's wallet, you can see here: https://github.com/Sifchain/sifnode/blob/8b2f9c45130c79e07555735185fbe1d00279fab0/x/clp/keeper/pool.go#L128
        sifnoded_client = sifchain.Sifnoded(self.cmd, home=client_home)
        denom_total_supply = 10000 * self.number_of_wallets * self.amount_of_denom_per_wallet
        wallets = {}
        extra_accounts = {}
        for i in range(self.number_of_wallets):
            chosen_tokens = [self.tokens[i] for i in random_choice(self.liquidity_providers_per_wallet, len(self.tokens), rnd=self.rnd)]
            balances = {denom: self.amount_of_denom_per_wallet for denom in chosen_tokens}
            if predefined_wallets:
                addr = predefined_wallets.create_acct()
            else:
                mnemonic = None if ((self.custom_wallet_mnemonics is None) or (i >= len(self.custom_wallet_mnemonics))) else self.custom_wallet_mnemonics[i].split(" ")
                addr = sifnoded_client.create_addr(mnemonic=mnemonic)
            wallets[addr] = chosen_tokens
            extra_accounts[addr] = cosmos.balance_add(balances, {ROWAN: self.amount_of_rowan_per_wallet})

        # Create validators.
        # To create liquidity pools, faucet needs to have enough balances for all denoms. During
        # setup_liquidity_pools_simple() tokens will be transferred from faucet to clp_admin automatically.
        # (Currently, clp_admin == admin account of validator 0).
        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        validator_admin_initial_balance = {ROWAN: 10**25}
        for i in range(self.number_of_nodes):
            mnemonic = self.validator0_mnemonic if i == 0 else None
            env.add_validator(admin_mnemonic=mnemonic, initial_balance=validator_admin_initial_balance)

        faucet_balance = cosmos.balance_add({denom: denom_total_supply for denom in self.tokens}, {ROWAN: 10**25})
        env.init(faucet_balance=faucet_balance, extra_accounts=extra_accounts)
        env.start()

        self.env = env
        sifnoded = env._sifnoded_for(env.node_info[0])
        sifnoded_client = sifchain.Sifnoded(self.cmd, home=client_home, node=sifchain.format_node_url(
            self.env.node_info[0]["host"], self.env.node_info[0]["ports"]["rpc"]), chain_id=env.chain_id)
        sifnoded_client.get_balance_default_retries = 5

        # Set up liquidity pools. We create them symmetrically (`native_amount == external_amount`).
        native_amount = external_amount = self.amount_of_denom_per_wallet
        pools_definitions = {denom: (18, native_amount, external_amount) for denom in self.tokens}
        env.setup_liquidity_pools_simple(pools_definitions)

        # Set up liquidity providers. We create them symmetrically (`native_amount == externam_amount`). The ratio of
        # native vs. external amount has to be the same as for corresponding pool (within certain rounding tolerance).
        # Calling `tx_clp_add_liquidity` to add multiple liquidity providers within the same block does not work (only
        # the first call gets through). To avoid `--broadcast-mode block` or waiting for new block, we need to use
        # account sequence numbers.
        for addr, denoms in wallets.items():
            account_number, account_sequence = sifnoded.get_acct_seq(addr)
            for denom in denoms:
                res = sifnoded_client.tx_clp_add_liquidity(addr, denom, self.amount_of_liquidity_added_by_wallet,
                    self.amount_of_liquidity_added_by_wallet, account_seq=(account_number, account_sequence))
                sifchain.check_raw_log(res)
                account_sequence += 1
        sifnoded_client.wait_for_last_transaction_to_be_mined()
        self.check_actual_liquidity_providers(sifnoded_client, env.clp_admin, wallets)

        self.sifnoded = sifnoded
        self.sifnoded_client = sifnoded_client

    def run(self):
        sifnoded = self.sifnoded
        sifnoded_client = self.sifnoded_client
        admin_addr = self.env.node_info[0]["admin_addr"]

        # Determine start and end blocks for rewards and LPPD
        # TODO start and end blocks are both inclusive, adjust
        current_block = sifnoded_client.get_current_block()
        start_block = current_block + 5
        rewards_start_block = start_block + self.rewards_offset_blocks
        rewards_end_block = rewards_start_block + self.rewards_duration_blocks
        lppd_start_block = start_block + self.lpd_offset_blocks
        lppd_end_block = lppd_start_block + self.lpd_duration_blocks

        wait_boundaries = set()
        wait_boundaries.add(start_block)

        # Set up rewards
        if self.rewards_duration_blocks > 0:
            reward_params = sifchain.create_rewards_descriptor("RP_1", rewards_start_block, rewards_end_block,
                [(token, 1) for token in self.tokens][:self.reward_period_pool_count], 100000 * self.token_unit,
                self.reward_period_default_multiplier, self.reward_period_distribute, self.reward_period_mod)
            sifnoded.clp_reward_period(admin_addr, reward_params)
            sifnoded.wait_for_last_transaction_to_be_mined()
            wait_boundaries.add(rewards_start_block)
            wait_boundaries.add(rewards_end_block)
        # TODO sifnoded query reward params --node --chain-id (check if/when implemented)

        # Set up LPD policies
        if self.lpd_duration_blocks > 0:
            lppd_params = sifchain.create_lppd_params(lppd_start_block, lppd_end_block, 0.00045, self.lpd_period_mod)
            sifnoded.clp_set_lppd_params(admin_addr, lppd_params)
            sifnoded.wait_for_last_transaction_to_be_mined()
            wait_boundaries.add(lppd_start_block)
            wait_boundaries.add(lppd_end_block)

        wait_boundaries = sorted(list(wait_boundaries))
        cnt = len(wait_boundaries)
        if cnt < 2:
            log.info("Not measuring block time - nothing to wait for")
        else:
            block_time_per_phase = []
            log.info("Waiting for phase 0...")
            prev_time = self.wait_for_block(wait_boundaries[0])
            for i in range(1, cnt):
                log.info("Waiting for phase {}...".format(i))
                next_time = self.wait_for_block(wait_boundaries[i])
                block_time_per_phase.append((next_time - prev_time) / (wait_boundaries[i] - wait_boundaries[i - 1]))
                prev_time = next_time
            for i in range(cnt - 1):
                self.report("Block time for phase {} (blocks {} - {}): {}".format(i + 1,
                    wait_boundaries[i], wait_boundaries[i + 1], block_time_per_phase[i]))

        # TODO LPD and rewards assertions
        # See https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5#7392be2c1a034d2db83b9b38ab89ff9e

        # run_forever means we're not interested in average block times but want to run this
        # as an environment
        if self.run_forever:
            wait_for_enter_key_pressed()

    def report(self, message):
        log.info(message)

    # TODO Refactor - move to Sifnoded
    def wait_for_block(self, block_number: int) -> float:
        sifnoded = self.sifnoded
        current_block = sifnoded.get_current_block()
        prev_block = None
        assert current_block < block_number
        while current_block < block_number:
            if self.block_results_offset is not None:
                if (prev_block is None) or (current_block != prev_block):
                    # This is just for collecting statistics while we wait, the test result does not depend on it.
                    # Check also https://github.com/cosmos/cosmos-sdk/issues/6105
                    try:
                        height = current_block - self.block_results_offset if self.block_results_offset is not None else None
                        blk = sifnoded.get_block_results(height=height)
                        histogram = {}
                        for key in ["begin_block_events", "end_block_events"]:
                            items = blk[key]
                            if items is not None:
                                for evt in items:
                                    evt_type = evt["type"]
                                    if evt_type not in histogram:
                                        histogram[evt_type] = 0
                                    histogram[evt_type] += 1
                        log.debug("Block events for block {}: {}".format(current_block, repr(histogram)))
                        prev_block = current_block
                    except Exception as e:
                        log.error("HTTP request for block_results failed: {}".format(repr(e)))
            time.sleep(1)
            current_block = sifnoded.get_current_block()
        return time.time()

    def assert_set_equal(self, message: str, actual: set, expected: set):
        if actual != expected:
            actual_only = actual.difference(expected)
            expected_only = expected.difference(actual)
            log.error("Assertion failed: {}: actual={}".format(message, repr(actual)))
            log.error("Assertion failed: {}: expected={}".format(message, repr(expected)))
            log.error("Assertion failed: {}: actual_only={}".format(message, repr(actual_only)))
            log.error("Assertion failed: {}: expected_only={}".format(message, repr(expected_only)))

    def check_actual_liquidity_providers(self, sifnoded, clp_admin, wallets):
        actual_lp_providers = {}
        for denom in self.tokens:
            for lp in sifnoded.query_clp_liquidity_providers(denom):
                addr = lp["liquidity_provider_address"]
                symbol = lp["asset"]["symbol"]
                if addr not in actual_lp_providers:
                    actual_lp_providers[addr] = set()
                actual_lp_providers[addr].add(symbol)
        # Note: "clp_admin" will automatically be a liquidity provider for all since it had created the pools
        expected_lp_providers = dict_merge({clp_admin: set(self.tokens)}, {addr: wallets[addr] for addr in wallets})
        if self.disable_assertions:
            act = set(actual_lp_providers)
            exp = set(expected_lp_providers)
            self.assert_set_equal("LP providers mismatch", act, exp)
            for addr in wallets:
                act = set(actual_lp_providers[addr])
                exp = set(expected_lp_providers[addr])
                self.assert_set_equal("LP mismatch for wallet {}".format(addr), act, exp)
        else:
            assert set(actual_lp_providers) == set(expected_lp_providers)  # Keys
            assert all(set(actual_lp_providers[addr]) == set(expected_lp_providers[addr]) for addr in actual_lp_providers)  # Values


def run_test_case(args):
    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    sifnoded_home_root = cmd.tmpdir("siftool-test.tmp")
    cmd.rmdir(sifnoded_home_root)
    test = Test(cmd, prj, sifnoded_home_root=sifnoded_home_root)

    scenario_vars = json.loads(cmd.read_text_file(args.scenario_file))

    def get_arg(name, default):
        result = args.__dict__.get(name, None)
        if result is not None:
            return result
        if name in scenario_vars:
            return scenario_vars[name]
        else:
            return default

    test.number_of_nodes = get_arg("number_of_nodes", 1)
    test.number_of_liquidity_pools = get_arg("number_of_liquidity_pools", 10)
    test.number_of_wallets = get_arg("number_of_wallets", 10)
    test.liquidity_providers_per_wallet = get_arg("liquidity_providers_per_wallet", 3)
    test.reward_period_default_multiplier = get_arg("reward_period_default_multiplier", 0.0)
    test.reward_period_distribute = get_arg("reward_period_distribute", False)
    test.reward_period_mod = get_arg("reward_period_mod", 1)
    test.reward_period_pool_count = get_arg("reward_period_pool_count", test.number_of_liquidity_pools)
    test.rewards_offset_blocks = get_arg("rewards_offset_blocks", 100)
    test.rewards_duration_blocks = get_arg("rewards_duration_blocks", 200)
    test.lpd_period_mod = get_arg("lpd_period_mod", 1)
    test.lpd_offset_blocks = get_arg("lpd_offset_blocks", 200)
    test.lpd_duration_blocks = get_arg("lpd_duration_blocks", 200)
    test.block_results_offset = get_arg("block_results_offset", None)
    test.run_forever = get_arg("run_forever", False)
    test.disable_assertions = get_arg("disable_assertions", False)
    test.wallets_dir = get_arg("wallets_dir", None)

    assert test.number_of_nodes >= 1
    assert test.number_of_nodes < 10  # Change ports to ensure they do not overlap
    assert test.liquidity_providers_per_wallet > 0
    assert test.liquidity_providers_per_wallet <= test.number_of_liquidity_pools
    assert test.reward_period_pool_count <= test.number_of_liquidity_pools

    test.report("number_of_nodes: {}".format(test.number_of_nodes))
    test.report("number_of_liquidity_pools: {}".format(test.number_of_liquidity_pools))
    test.report("number_of_wallets: {}".format(test.number_of_wallets))
    test.report("liquidity_providers_per_wallet: {}".format(test.liquidity_providers_per_wallet))
    test.report("reward_period_default_multiplier: {}".format(test.reward_period_default_multiplier))
    test.report("reward_period_distribute: {}".format(test.reward_period_distribute))
    test.report("reward_period_mod: {}".format(test.reward_period_mod))
    test.report("reward_period_pool_count: {}".format(test.reward_period_pool_count))
    test.report("rewards_offset_blocks: {}".format(test.rewards_offset_blocks))
    test.report("rewards_duration_blocks: {}".format(test.rewards_duration_blocks))
    test.report("lpd_period_mod: {}".format(test.lpd_period_mod))
    test.report("lpd_offset_blocks: {}".format(test.lpd_offset_blocks))
    test.report("lpd_duration_blocks: {}".format(test.lpd_duration_blocks))
    test.report("block_results_offset: {}".format(test.block_results_offset))
    test.report("run_forever: {}".format(test.run_forever))
    test.report("wallets_dir: {}".format(test.wallets_dir))

    test_start_time = time.time()
    test.setup()
    run_start_time = time.time()

    try:
        test.run()
        test_finish_time = time.time()
        log.info("Finished successfully, setup: {:.2f}s, total {:.2f}s".format(run_start_time - test_start_time,
            test_finish_time - test_start_time))
    except Exception as e:
        log.error("Test failed", exc_info=True)
        try:
            log.error("Checking some balances to see if the thing is dead or alive...")
            addr = test.node_info[0]["admin_addr"]
            balance = cosmos.balance_format(test.sifnoded[0].get_balance(addr))
            log.debug("Balance of {}: {}".format(addr, balance))
        except Exception as e:
            log.error("Balance check failed", exc_info=True)


def main(argv: List[str]):
    basic_logging_setup()
    parser = argparse.ArgumentParser()
    parser.add_argument("scenario_file", type=str)
    parser.add_argument("--block-results-offset", type=int, default=None)
    parser.add_argument("--run-forever", action="store_true")
    parser.add_argument("--disable-assertions", action="store_true")
    parser.add_argument("--wallets-dir", type=str)
    args = parser.parse_args(argv[1:])

    run_test_case(args)


if __name__ == "__main__":
    main(sys.argv)
