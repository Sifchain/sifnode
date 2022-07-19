# Scenarion description: https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5
# Ticket: https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifnode/3020

import argparse
import random
import sys
import time
import siftool_path
from siftool.common import *
from siftool import test_utils, command, sifchain, project, cosmos
from siftool.sifchain import ROWAN, STAKE


log = siftool_logger(__name__)


def create_rewards_descriptor(start_block, end_block, allocation) -> sifchain.RewardsParams:
    return {
        "reward_period_id": "RP_1",
        "reward_period_start_block": start_block,
        "reward_period_end_block": end_block,
        "reward_period_allocation": str(allocation),
        "reward_period_pool_multipliers": [{
            "pool_multiplier_asset": "cusdt",
            "multiplier": "2"
        }, {
            "pool_multiplier_asset": "clink",
            "multiplier": "1"
        }],
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
        self.sifnode = sifchain.Sifnoded(cmd, home=sifnoded_home, chain_id="localnet")

        # Note: pools need to be symmetrical (down to certain decimals)!

        # Number of pools == number of tokens.
        # We create this many tokens, and create a pool for each token.
        self.number_of_tokens = 100

        self.token_decimals = 18
        self.token_unit = 10**self.token_decimals
        self.tokens = ["test{}".format(i) for i in range(self.number_of_tokens)]

        # The test runs, in sequence
        # test_duration_blocks of neither rewards nor lppd
        # test_duration_blocks of rewards without lppd
        # test_duration_blocks of rewards and lppd
        # test_duration_blocks of lppd without rewards
        # The timing starts with the next block after setup. The accuracty of the test is limited by polling for the
        # current block number (1s). The total time will be 4 * test_duration_blocks * block_time, i.e.
        # 4 * 6s for one unit of test_duration_blocks.
        self.test_duration_blocks = 1

        self.number_of_wallets = 10

        # Number of different tokens in a wallet. Tokens are randomly chosen from all (number_of_tokens).
        self.number_of_denoms_per_wallet = 5  # Also how much of liquidity every wallet provides
        assert self.number_of_tokens >= self.number_of_denoms_per_wallet

        self.amount_of_denom_per_wallet = 1000 * self.token_unit
        self.amount_of_liquidity_added_by_wallet = 500 * self.token_unit

        # self.number_of_pools = 80  # choice from [80, 100, 200]


    def initial_setup(self):
        sifnoded = self.sifnode

        self.prj.pkill()
        time.sleep(1)

        # Total minted supply per denom. This has to cover the sum of all balances, rewards and lppd distribution.
        denom_total_supply = 1000 * self.number_of_wallets * self.amount_of_denom_per_wallet

        sifnode_moniker = "test"
        self.sifnode.init(sifnode_moniker)
        sif_name = "sif"
        akasha_name = "akasha"
        sif_mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow".split()
        akasha_mnemonic = "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard".split(" ")
        sif = self.sifnode.keys_add(sif_name, sif_mnemonic)["address"]
        akasha = self.sifnode.keys_add(akasha_name, akasha_mnemonic)["address"]
        sif_bech = self.sifnode.get_val_address(sif)
        tokens_sif = {
            ROWAN: 999 * 10**30,
            STAKE: 999 * 10**30,
            "ceth": 999 * 10**30,
            "cusdc": 999 * 10**30,
            "cusdt": 999 * 10**30,
            "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": 999 * 10**30,
            "ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA": 999 * 10**30,
            "ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D": 999 * 10**30,
        } | {denom: denom_total_supply for denom in self.tokens}
        tokens_akasha = {
            ROWAN: 5 * 10**23,
            STAKE: 99 * 10**25,
            "catk": 99 * 10**25,
            "cbtk": 99 * 10**25,
            "ceth": 99 * 10**25,
            "cdash": 99 * 10**25,
            "clink": 99 * 10**25,
        }
        mkey = self.sifnode.keys_add_multisig("mkey", [sif_name, akasha_name], 2)["address"]
        sifnoded.add_genesis_account(sif, tokens_sif)
        sifnoded.add_genesis_account(akasha, tokens_akasha)
        sifnoded.add_genesis_clp_admin(sif)
        sifnoded.add_genesis_clp_admin(akasha)
        sifnoded.set_genesis_oracle_admin(sif_name)
        sifnoded.add_genesis_validators(sif_bech)
        sifnoded.set_genesis_whitelister_admin(sif_name)
        # sifnoded get-gen-denom-whitelist scripts/denoms.json
        sifnoded.gentx(sif_name, {STAKE: 10**24})
        sifnoded.collect_gentx()
        sifnoded.validate_genesis()

        # Start process

        sifnoded_log_file = open("/tmp/sifnoded.log", "w")
        sifnoded_proc = sifnoded.sifnoded_start(log_file=sifnoded_log_file, trace=True)
        # Currently returning 404
        # self.sifnode.wait_up("localhost", 26657)
        time.sleep(5)

        self.wait_for_block(sifnoded.get_current_block() + 3)

        # Set up pools

        for denom in self.tokens:
            entry = sifnoded.create_tokenregistry_entry(denom, denom, 18, ["CLP"])
            sifnoded.token_registry_register(entry, sif)
        sifnoded.wait_for_last_transaction_to_be_mined()

        # Set up each wallets with test tokens

        # TODO We want to use EnvCtx.create_sifchain_addr
        def create_sifchain_addr(moniker: str, fund_amounts: Optional[cosmos.Balance]):
            new_addr = sifnoded.keys_add(moniker)["address"]
            if fund_amounts:
                sifnoded.send_batch(sif, new_addr, fund_amounts)
                assert cosmos.balance_equal(sifnoded.get_balance(new_addr), fund_amounts)
            return new_addr

        # Set up wallets.
        # For each wallet we choose `number_of_denoms_per_wallet` random denoms and top it up with 500 `token_units`
        # from `sif` account.
        wallets = {}
        for i in range(self.number_of_wallets):
            chosen_tokens = [self.tokens[i] for i in random_choice(self.number_of_denoms_per_wallet, len(self.tokens), rnd=self.rnd)]
            balances = {denom: self.amount_of_denom_per_wallet for denom in chosen_tokens}
            fund_amounts = cosmos.balance_add(balances, {ROWAN: 10000 * 10**18})  # Add some rowan TODO How much?
            addr = create_sifchain_addr("test-wallet-{}".format(i), fund_amounts=fund_amounts)
            wallets[addr] = chosen_tokens

        sifnoded.wait_for_last_transaction_to_be_mined()

        # Add liquidity
        for addr, denoms in wallets.items():
            for denom in denoms:
                sifnoded.tx_clp_add_liquidity(addr, denom, self.amount_of_liquidity_added_by_wallet, self.amount_of_denom_per_wallet)

        current_block = self.sifnode.get_current_block()

        # |<---test time------------------------------------------------------------------->|
        #
        #                          |<--------------rewards-------------->|
        #                                             ]<---------------lppd---------------->|
        #
        #       |<----neither----->|<--rewards only-->|<--rewards+lppd-->|<--rewards only-->|
        #       ^-- time0          ^-- time1          ^-- time2          ^-- time3          ^-- time4

        start_block = current_block + 10
        rewards_start_block = start_block + self.test_duration_blocks
        rewards_end_block = rewards_start_block + 2 * self.test_duration_blocks
        lppd_start_block = start_block + 2 * self.test_duration_blocks
        lppd_end_block = lppd_start_block + 2 * self.test_duration_blocks

        # Set up liquidity providers
        reward_params = create_rewards_descriptor(rewards_start_block, rewards_end_block, 100000 * self.token_unit)
        sifnoded.clp_reward_period(sif, reward_params)

        sifnoded.wait_for_last_transaction_to_be_mined()

        lppd_params = create_lppd_params(lppd_start_block, lppd_end_block, 0.00045)
        sifnoded.clp_set_lppd_params(sif, lppd_params)

        time0 = self.wait_for_block(start_block)
        time1 = self.wait_for_block(rewards_start_block)
        time2 = self.wait_for_block(lppd_start_block)
        time3 = self.wait_for_block(rewards_end_block)
        time4 = self.wait_for_block(lppd_end_block)

        print("Nothing:       {:.2f} s/block".format((time1 - time0) / self.test_duration_blocks))
        print("Rewards only:  {:.2f} s/block".format((time2 - time1) / self.test_duration_blocks))
        print("Reards + LPPD: {:.2f} s/block".format((time3 - time2) / self.test_duration_blocks))
        print("LPPD only:     {:.2f} s/block".format((time4 - time3) / self.test_duration_blocks))

        return sifnoded_proc, sifnoded_log_file

    # TODO Refactor - move to Sifnoded
    def wait_for_block(self, block_number):
        current_block = self.sifnode.get_current_block()
        assert current_block < block_number
        while current_block < block_number:
            time.sleep(1)
            current_block = self.sifnode.get_current_block()
        return time.time()

    def run(self, number_of_pools: int, number_of_liqudity_providers: int):
        self.initial_setup()
        return


def main(argv: List[str]):
    basic_logging_setup()
    # ctx = test_utils.get_env_ctx()
    parser = argparse.ArgumentParser()
    parser.add_argument("--number-of-pools", type=int, default=10)
    parser.add_argument("--number-of-liquidity-providers", type=int, default=10)
    args = parser.parse_args(argv[1:])
    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    sifnoded_home = cmd.tmpdir("siftool-test.tmp")
    cmd.rmdir(sifnoded_home)
    test = Test(cmd, prj, sifnoded_home=sifnoded_home)
    test.run(args.number_of_pools, args.number_of_liquidity_providers)
    log.info("Finished successfully")


if __name__ == "__main__":
    main(sys.argv)
