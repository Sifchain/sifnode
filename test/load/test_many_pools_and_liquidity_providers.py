# Scenarion description: https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5
# Ticket: https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifnode/3020

import argparse
import random
import sys
import time
from typing import Tuple, Iterable
import siftool_path
from siftool.common import *
from siftool import test_utils, command, sifchain, project, cosmos
from siftool.sifchain import ROWAN, STAKE


log = siftool_logger(__name__)


def create_rewards_descriptor(start_block: int, end_block: int, multipliers: Iterable[Tuple[str, int]], allocation: int) -> sifchain.RewardsParams:
    return {
        "reward_period_id": "RP_1",
        "reward_period_start_block": start_block,
        "reward_period_end_block": end_block,
        "reward_period_allocation": str(allocation),
        "reward_period_pool_multipliers": [{
            "pool_multiplier_asset": denom,
            "multiplier": str(multiplier)
        } for denom, multiplier in multipliers],
        "reward_period_default_multiplier": "0.0",
        "reward_period_distribute": False,
        "reward_period_mod": 1
    }

def create_lppd_params(start_block, end_block, rate) -> sifchain.LPPDParams:
    return {
        "distribution_period_block_rate": str(rate),
        "distribution_period_start_block": start_block,
        "distribution_period_end_block": end_block,
        "distribution_period_mod": 1
    }

class Test:
    def __init__(self, cmd: command.Command, prj: project.Project, sifnoded_home: str = None):
        self.cmd = cmd
        self.prj = prj
        self.rnd = random.Random(5646067977921730044)  # Use a fixed random seed for repeatable results
        self.sifnoded = sifchain.Sifnoded(cmd, home=sifnoded_home, chain_id="localnet")

        # Note: pools need to be symmetrical (down to certain decimals)!

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
        #       |<----neither----->|<--rewards only-->|<--rewards+lppd-->|<--rewards only-->|
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
        assert self.liquidity_providers_per_wallet > 0
        assert self.liquidity_providers_per_wallet <= self.number_of_liquidity_pools

        # We are only dealing with symmetrical liquidity pools here.
        # This means that each liquidity provider uses `native_amount == external_amount`.
        # The ratio `native_amount/external_amount` has to be the same as defined per liquidity pool
        # (withing a certain threshold).
        # We use 1000/1000 for pool and 500/500 for liquidity provider.
        self.amount_of_denom_per_wallet = 1000 * self.token_unit
        self.amount_of_liquidity_added_by_wallet = 500 * self.token_unit

        self.amount_of_rowan_per_wallet = 10000 * 10**18  # TODO How much?

    def run_test(self):
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
        for i in range(self.number_of_wallets):
            chosen_tokens = [self.tokens[i] for i in random_choice(self.liquidity_providers_per_wallet, len(self.tokens), rnd=self.rnd)]
            balances = {denom: self.amount_of_denom_per_wallet for denom in chosen_tokens}
            fund_amounts = cosmos.balance_add(balances, {ROWAN: self.amount_of_rowan_per_wallet})
            addr = sifnoded.create_addr()
            wallets[addr] = chosen_tokens
            sifnoded.add_genesis_account(addr, fund_amounts)

        sifnoded.add_genesis_account(sif, tokens_sif)
        sifnoded.add_genesis_clp_admin(sif)
        sifnoded.set_genesis_oracle_admin(sif_name)
        sifnoded.add_genesis_validators(sif_bech)
        sifnoded.set_genesis_whitelister_admin(sif_name)
        sifnoded.gentx(sif_name, {STAKE: 10**24})
        sifnoded.collect_gentx()
        sifnoded.validate_genesis()

        # Start process
        sifnoded_log_file = open("/tmp/sifnoded.log", "w")
        sifnoded_proc = sifnoded.sifnoded_start(log_file=sifnoded_log_file, trace=True)
        # Currently returning 404
        # sifnoded.wait_up("localhost", 26657)
        time.sleep(5)
        sifnoded.wait_for_last_transaction_to_be_mined()
        self.wait_for_block(sifnoded.get_current_block() + 3)

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
        # TODO This is too slow (1 operation per block, we need at least 10000)
        for addr, denoms in wallets.items():
            for denom in denoms:
                sifnoded.tx_clp_add_liquidity(addr, denom, self.amount_of_liquidity_added_by_wallet, self.amount_of_liquidity_added_by_wallet)
                sifnoded.wait_for_last_transaction_to_be_mined()
        lp_providers_check = {}
        for denom in self.tokens:
            for lp in sifnoded.query_clp_liquidity_providers(denom):
                addr = lp["liquidity_provider_address"]
                symbol = lp["asset"]["symbol"]
                if addr not in lp_providers_check:
                    lp_providers_check[addr] = set()
                lp_providers_check[addr].add(symbol)
        # Note: "sif" address will have all of them
        assert lp_providers_check

        # Determine start and end blocks for rewards and LPPD
        current_block = sifnoded.get_current_block()
        start_block = current_block + 5
        rewards_start_block = start_block + self.test_duration_blocks
        rewards_end_block = rewards_start_block + 2 * self.test_duration_blocks
        lppd_start_block = start_block + 2 * self.test_duration_blocks
        lppd_end_block = lppd_start_block + 2 * self.test_duration_blocks

        # Set up rewards
        reward_params = create_rewards_descriptor(rewards_start_block, rewards_end_block,
            [(token, 1) for token in self.tokens], 100000 * self.token_unit)
        sifnoded.clp_reward_period(sif, reward_params)
        sifnoded.wait_for_last_transaction_to_be_mined()

        # Set up LPPD policies
        lppd_params = create_lppd_params(lppd_start_block, lppd_end_block, 0.00045)
        sifnoded.clp_set_lppd_params(sif, lppd_params)
        sifnoded.wait_for_last_transaction_to_be_mined()

        time0 = self.wait_for_block(start_block)
        time1 = self.wait_for_block(rewards_start_block)
        time2 = self.wait_for_block(lppd_start_block)
        time3 = self.wait_for_block(rewards_end_block)
        time4 = self.wait_for_block(lppd_end_block)

        print("Neither:       {:.2f} s/block".format((time1 - time0) / self.test_duration_blocks))
        print("Rewards only:  {:.2f} s/block".format((time2 - time1) / self.test_duration_blocks))
        print("Reards + LPPD: {:.2f} s/block".format((time3 - time2) / self.test_duration_blocks))
        print("LPPD only:     {:.2f} s/block".format((time4 - time3) / self.test_duration_blocks))

        # TODO LPPD and rewards assertions
        # See https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5#7392be2c1a034d2db83b9b38ab89ff9e

        return sifnoded_proc, sifnoded_log_file

    # TODO Refactor - move to Sifnoded
    def wait_for_block(self, block_number):
        current_block = self.sifnoded.get_current_block()
        assert current_block < block_number
        while current_block < block_number:
            blk = self.sifnoded.get_block_results()
            begin_block_events = blk["begin_block_events"]
            histogram = {}
            for e in begin_block_events:
                type = e["type"]
                if type not in histogram:
                    histogram[type] = 0
                histogram[type] += 1
            log.debug("Block events: {}".format(repr(histogram)))
            time.sleep(1)
            current_block = self.sifnoded.get_current_block()
        return time.time()


def run(number_of_liquidity_pools: int, number_of_wallets: int, liquidity_providers_per_wallet: int,
    test_duration_blocks
):
    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    sifnoded_home = cmd.tmpdir("siftool-test.tmp")
    cmd.rmdir(sifnoded_home)
    test = Test(cmd, prj, sifnoded_home=sifnoded_home)
    test.number_of_liquidity_pools = number_of_liquidity_pools
    test.number_of_wallets = number_of_wallets
    test.liquidity_providers_per_wallet = liquidity_providers_per_wallet
    test.test_duration_blocks = test_duration_blocks
    test.run_test()
    log.info("Finished successfully")

def main(argv: List[str]):
    basic_logging_setup()
    parser = argparse.ArgumentParser()
    parser.add_argument("--number-of-liquidity-pools", type=int, default=10)
    parser.add_argument("--number-of-wallets", type=int, default=10)
    parser.add_argument("--liquidity-providers-per-wallet", type=int, default=5)
    parser.add_argument("--test-duration-blocks", type=int, default=5)
    args = parser.parse_args(argv[1:])
    run(args.number_of_liquidity_pools, args.number_of_wallets, args.liquidity_providers_per_wallet, args.test_duration_blocks)

if __name__ == "__main__":
    main(sys.argv)
