# Scenarion description: https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5
# Ticket: https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifnode/3020

import argparse
import random
import sys
import time
import siftool_path
from siftool.common import *
from siftool import command, sifchain, project, cosmos
from siftool.sifchain import ROWAN, STAKE


log = siftool_logger()


class Test:
    def __init__(self, cmd: command.Command, prj: project.Project, sifnoded_home: str = None):
        self.cmd = cmd
        self.prj = prj
        self.rnd = random.Random(5646067977921730044)  # Use a fixed random seed for repeatable results
        self.sifnoded = sifchain.Sifnoded(cmd, home=sifnoded_home, chain_id="localnet")

        # Number of pools == number of tokens.
        # We create this many tokens, and create a pool for each token.
        self.number_of_liquidity_pools = 10

        # Define one token per liquidity pool.
        self.token_decimals = 18
        self.token_unit = 10**self.token_decimals
        self.tokens = ["test{}".format(i) for i in range(self.number_of_liquidity_pools)]

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
        # The timing starts with the next block after setup. The accuracty of the test is limited by polling for the
        # current block number (1s). The total time will be 4 * test_duration_blocks * block_time, i.e.
        # 4 * 6s = 24s for one unit of test_duration_blocks.
        self.test_duration_blocks = 5  # 2 min

        # Number of wallets. Since each wallet provides liquidity to 1 or more liquidity pool, this is also the number
        # of unique liquidity providers.
        self.number_of_wallets = 10

        # The number of liquidity pools to which each wallet provides liquidity. The pools are chosen randomly from
        # all `number_of_liquidity_pools`. This is also the same of number of different tokens per wallet (not counting
        # rowan).
        self.liquidity_providers_per_wallet = 5

        self.reward_period_default_multiplier = 1.0
        self.reward_period_distribute = True
        self.reward_period_mod = 1
        self.reward_period_pool_count = 10

        # We are only dealing with symmetrical liquidity pools here.
        # This means that each liquidity provider uses `native_amount == external_amount`.
        # The ratio `native_amount/external_amount` has to be the same as defined per liquidity pool
        # (withing a certain threshold).
        # We use 1000/1000 for pools and 500/500 for liquidity providers.
        self.amount_of_denom_per_wallet = 1000 * self.token_unit
        self.amount_of_liquidity_added_by_wallet = 500 * self.token_unit

        self.amount_of_rowan_per_wallet = 10000 * 10**18  # TODO How much?

    def run_test(self):
        assert self.liquidity_providers_per_wallet > 0
        assert self.liquidity_providers_per_wallet <= self.number_of_liquidity_pools
        assert self.reward_period_pool_count <= self.number_of_liquidity_pools
        assert self.test_duration_blocks > 0

        sifnoded = self.sifnoded

        self.prj.pkill()
        time.sleep(1)

        # Admin balances per denom. This has to cover any rewards and lppd distribution.
        denom_total_supply = 10000 * self.number_of_wallets * self.amount_of_denom_per_wallet

        sifnode_moniker = "test"
        sifnoded.init(sifnode_moniker)
        sif_name = "sif"
        sif_mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow".split()
        sif = sifnoded.keys_add(sif_name, sif_mnemonic)["address"]
        sif_bech = sifnoded.get_val_address(sif)
        tokens_sif = {
            ROWAN: 999 * 10**30,
            STAKE: 999 * 10**30,
        } | {denom: denom_total_supply for denom in self.tokens}

        # Set up test wallets with test tokens. We do this in genesis for performance reasons. For each wallet we choose
        # number_of_denoms_per_wallet` random denoms.
        wallets = {}
        fund_amounts = {}
        for i in range(self.number_of_wallets):
            chosen_tokens = [self.tokens[i] for i in random_choice(self.liquidity_providers_per_wallet, len(self.tokens), rnd=self.rnd)]
            balances = {denom: self.amount_of_denom_per_wallet for denom in chosen_tokens}
            addr = sifnoded.create_addr()
            wallets[addr] = chosen_tokens
            fund_amounts[addr] = cosmos.balance_add(balances, {ROWAN: self.amount_of_rowan_per_wallet})
        sifnoded.add_genesis_account_directly_to_existing_genesis_json(fund_amounts)

        self._debug_fund_amounts = fund_amounts

        sifnoded.add_genesis_account(sif, tokens_sif)
        sifnoded.add_genesis_clp_admin(sif)
        sifnoded.set_genesis_oracle_admin(sif_name)
        sifnoded.add_genesis_validators(sif_bech)
        sifnoded.set_genesis_whitelister_admin(sif_name)
        sifnoded.gentx(sif_name, {STAKE: 10**24})
        sifnoded.collect_gentx()
        sifnoded.validate_genesis()

        # Start process
        sifnoded_log_file = open(self.cmd.tmpdir("sifnoded.log"), "w")
        sifnoded_proc = sifnoded.sifnoded_start(log_file=sifnoded_log_file, log_level="debug", trace=True)
        # Currently returning 404
        # sifnoded.wait_up("localhost", 26657)
        time.sleep(5)
        sifnoded.wait_for_last_transaction_to_be_mined()
        self.wait_for_block(sifnoded.get_current_block() + 3)

        # Check balances
        assert all(cosmos.balance_equal(sifnoded.get_balance(addr), fund_amounts[addr]) for addr in fund_amounts)
        assert all(sifnoded.get_acct_seq(addr)[1] == 0 for addr in wallets)

        # Add tokens to token registry. The minimum required permissions are  CLP.
        # TODO Might want to use `tx tokenregistry set-registry` to do it in one step (faster)
        #      According to @sheokapr `tx tokenregistry set-registry` also works for only one entry
        #      But`tx tokenregistry register-all` also works only for one entry.
        for denom in self.tokens:
            entry = sifnoded.create_tokenregistry_entry(denom, denom, 18, ["CLP"])
            sifnoded.token_registry_register(entry, sif)
            sifnoded.wait_for_last_transaction_to_be_mined()  # Must be run synchronously! (if not, only the first will work)
        tokenregistry_entries = sifnoded.query_tokenregistry_entries()
        assert set(x["denom"] for x in tokenregistry_entries) == set(self.tokens)

        # Set up liquidity pools. We create them symmetrically (`native_amount == externam_amount`).
        for denom in self.tokens:
            sifnoded.tx_clp_create_pool(sif, denom, self.amount_of_denom_per_wallet, self.amount_of_denom_per_wallet)
            sifnoded.wait_for_last_transaction_to_be_mined()
        pools = sifnoded.query_pools()
        assert set(x["external_asset"]["symbol"] for x in pools) == set(self.tokens)

        # Set up liquidity providers. We create them symmetrically (`native_amount == externam_amount`). The ratio of
        # native vs. external amount has to be the same as for corresponding pool (within certain rounding tolerance).
        # Calling `tx_clp_add_liquidity` to add multiple liquidity providers within the same block does not work (only
        # the first call gets through). To avoid `--broadcast-mode block` or waiting for new block, we need to use
        # account sequence numbers.
        for addr, denoms in wallets.items():
            account_number, account_sequence = sifnoded.get_acct_seq(addr)
            for denom in denoms:
                sifnoded.tx_clp_add_liquidity(addr, denom, self.amount_of_liquidity_added_by_wallet,
                    self.amount_of_liquidity_added_by_wallet, account_seq=(account_number, account_sequence))
                account_sequence += 1
        sifnoded.wait_for_last_transaction_to_be_mined()
        actual_lp_providers = {}
        for denom in self.tokens:
            for lp in sifnoded.query_clp_liquidity_providers(denom):
                addr = lp["liquidity_provider_address"]
                symbol = lp["asset"]["symbol"]
                if addr not in actual_lp_providers:
                    actual_lp_providers[addr] = set()
                actual_lp_providers[addr].add(symbol)
        # Note: "sif" address will automatically be a liquidity provider for all => remove "sif" before asserting
        actual_lp_providers.pop(sif)
        assert set(actual_lp_providers) == set(wallets)  # Keys
        assert all(set(actual_lp_providers[addr]) == set(wallets[addr]) for addr in wallets)  # Values

        # Determine start and end blocks for rewards and LPPD
        current_block = sifnoded.get_current_block()
        start_block = current_block + 5
        rewards_start_block = start_block + self.test_duration_blocks
        rewards_end_block = rewards_start_block + 2 * self.test_duration_blocks - 1
        lppd_start_block = start_block + 2 * self.test_duration_blocks
        lppd_end_block = lppd_start_block + 2 * self.test_duration_blocks - 1

        # Set up rewards
        reward_params = sifchain.create_rewards_descriptor("RP_1", rewards_start_block, rewards_end_block,
            [(token, 1) for token in self.tokens][:self.reward_period_pool_count], 100000 * self.token_unit,
            self.reward_period_default_multiplier, self.reward_period_distribute, self.reward_period_mod)
        sifnoded.clp_reward_period(sif, reward_params)
        sifnoded.wait_for_last_transaction_to_be_mined()

        # Set up LPPD policies
        lppd_params = sifchain.create_lppd_params(lppd_start_block, lppd_end_block, 0.00045)
        sifnoded.clp_set_lppd_params(sif, lppd_params)
        sifnoded.wait_for_last_transaction_to_be_mined()

        time0 = self.wait_for_block(start_block)
        log.info("In phase 'neither'")
        time1 = self.wait_for_block(rewards_start_block)
        log.info("In phase 'rewards only'")
        time2 = self.wait_for_block(lppd_start_block)
        log.info("In phase 'rewards + LPPD'")
        time3 = self.wait_for_block(rewards_end_block)
        log.info("In phase 'LPPD only'")
        time4 = self.wait_for_block(lppd_end_block)

        accuracy = 1.0 / self.test_duration_blocks

        log.info("Neither: {:.2f} +/- {:.2f} s/block".format((time1 - time0) / self.test_duration_blocks, accuracy))
        log.info("Rewards only: {:.2f} +/- {:.2f} s/block".format((time2 - time1) / self.test_duration_blocks, accuracy))
        log.info("Rewards + LPPD: {:.2f} +/- {:.2f} s/block".format((time3 - time2) / self.test_duration_blocks, accuracy))
        log.info("LPPD only: {:.2f} +/- {:.2f} s/block".format((time4 - time3) / self.test_duration_blocks, accuracy))

        # TODO LPPD and rewards assertions
        # See https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5#7392be2c1a034d2db83b9b38ab89ff9e

        return sifnoded_proc, sifnoded_log_file

    # TODO Refactor - move to Sifnoded
    def wait_for_block(self, block_number: int) -> float:
        current_block = self.sifnoded.get_current_block()
        prev_block = None
        assert current_block < block_number
        while current_block < block_number:
            if (prev_block is None) or (current_block != prev_block):
                # This is just for collecting statistics while we wait, the test result does not depend on it.
                # Check also https://github.com/cosmos/cosmos-sdk/issues/6105
                try:
                    blk = self.sifnoded.get_block_results()
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
            current_block = self.sifnoded.get_current_block()
        return time.time()


def run(number_of_liquidity_pools: int, number_of_wallets: int, liquidity_providers_per_wallet: int,
    reward_period_default_multiplier: float, reward_period_distribute: bool, reward_period_mod: int,
    reward_period_pool_count: int, test_duration_blocks: int
):
    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    sifnoded_home = cmd.tmpdir("siftool-test.tmp")
    cmd.rmdir(sifnoded_home)
    test = Test(cmd, prj, sifnoded_home=sifnoded_home)

    log.info("sifnoded_home: {}".format(sifnoded_home))
    log.info("number_of_liquidity_pools: {}".format(number_of_liquidity_pools))
    log.info("number_of_wallets: {}".format(number_of_wallets))
    log.info("liquidity_providers_per_wallet: {}".format(liquidity_providers_per_wallet))
    log.info("reward_period_default_multiplier: {}".format(reward_period_default_multiplier))
    log.info("reward_period_distribute: {}".format(reward_period_distribute))
    log.info("reward_period_mod: {}".format(reward_period_mod))
    log.info("reward_period_pool_count: {}".format(reward_period_pool_count))
    log.info("test_duration_blocks: {}".format(test_duration_blocks))

    test.number_of_liquidity_pools = number_of_liquidity_pools
    test.number_of_wallets = number_of_wallets
    test.liquidity_providers_per_wallet = liquidity_providers_per_wallet
    test.reward_period_default_multiplier = reward_period_default_multiplier
    test.reward_period_distribute = reward_period_distribute
    test.reward_period_mod = reward_period_mod
    test.reward_period_pool_count = reward_period_pool_count
    test.test_duration_blocks = test_duration_blocks

    try:
        test_start_time = time.time()
        test.run_test()
        test_finish_time = time.time()
        log.info("Finished successfully io {:.2f}s".format(test_finish_time - test_start_time))
    except Exception as e:
        log.error("Test failed", e)
        try:
            log.error("Checking some balances to see if the thing is dead or alive...")
            for addr in list(test._debug_fund_amounts)[5:]:
                log.debug("Balance of {}: {}".format(addr, cosmos.balance_format(test.sifnoded.get_balance(addr))))
        except Exception as e:
            log.error("Balance check failed", e)
        wait_for_enter_key_pressed()


def main(argv: List[str]):
    basic_logging_setup()
    parser = argparse.ArgumentParser()
    parser.add_argument("--number-of-liquidity-pools", type=int, default=10)
    parser.add_argument("--number-of-wallets", type=int, default=10)
    parser.add_argument("--liquidity-providers-per-wallet", type=int, default=5)
    parser.add_argument("--reward-period-default-multiplier", type=float, default=0.0)
    parser.add_argument("--reward-period-distribute", action="store_true")
    parser.add_argument("--reward-period-mod", type=int, default=1)
    parser.add_argument("--reward-period-pool-count", type=int, default=10)
    parser.add_argument("--test-duration-blocks", type=int, default=5)
    args = parser.parse_args(argv[1:])
    run(args.number_of_liquidity_pools, args.number_of_wallets, args.liquidity_providers_per_wallet,
        args.reward_period_default_multiplier, args.reward_period_distribute, args.reward_period_mod,
        args.reward_period_pool_count, args.test_duration_blocks)

if __name__ == "__main__":
    main(sys.argv)
