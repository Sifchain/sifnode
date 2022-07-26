# Scenarion description: https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5
# Ticket: https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/sifnode/3020
#
# Example usage:
# cd test/load
# ../integration/framework/siftool venv-init
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
# (2) Exceptions / printing of _debug...
#
# (3) Use parameter rpc.laddr for sifnoded start instead of self.node


import argparse
import random
import sys
import time
import siftool_path
from siftool.common import *
from siftool import command, sifchain, project, cosmos
from siftool.sifchain import ROWAN, STAKE, ROWAN_DECIMALS


log = siftool_logger()


class Test:
    def __init__(self, cmd: command.Command, prj: project.Project, sifnoded_home_root: str):
        self.cmd = cmd
        self.prj = prj
        self.rnd = random.Random(5646067977921730044)  # Use a fixed random seed for repeatable results

        self.number_of_nodes = 1

        # Number of pools == number of tokens.
        # We create this many tokens, and create a pool for each token.
        self.number_of_liquidity_pools = 10

        # Number of wallets. Since each wallet provides liquidity to 1 or more liquidity pool, this is also the number
        # of unique liquidity providers.
        self.number_of_wallets = 10


        # The number of liquidity pools to which each wallet provides liquidity. The pools are chosen randomly from
        # all `number_of_liquidity_pools`. This is also the same of number of different tokens per wallet (not counting
        # rowan).
        self.liquidity_providers_per_wallet = 5

        self.reward_period_default_multiplier = 0.0
        self.reward_period_distribute = False
        self.reward_period_mod = 1
        self.reward_period_pool_count = 10

        # The timing starts with the next block after setup. The accuracty of the test is limited by polling for the
        # current block number (1s). The total time will be 4 * test_duration_blocks * block_time, i.e.
        # 4 * 6s = 24s for one unit of test_duration_blocks.
        self.test_duration_blocks = 5

        self.sifnoded_home_root = sifnoded_home_root

        self.chain_id = "localnet"
        self.sifnoded: Optional[List[sifchain.Sifnoded]] = None
        self.sifnoded = []
        self.node_info = None

    def setup(self):
        assert self.number_of_nodes >= 1
        assert self.number_of_nodes < 10  # Change ports to ensure they do not overlap
        assert self.liquidity_providers_per_wallet > 0
        assert self.liquidity_providers_per_wallet <= self.number_of_liquidity_pools
        assert self.reward_period_pool_count <= self.number_of_liquidity_pools
        assert self.test_duration_blocks > 0

        log.info("number_of_nodes: {}".format(self.number_of_nodes))
        log.info("number_of_liquidity_pools: {}".format(self.number_of_liquidity_pools))
        log.info("number_of_wallets: {}".format(self.number_of_wallets))
        log.info("liquidity_providers_per_wallet: {}".format(self.liquidity_providers_per_wallet))
        log.info("reward_period_default_multiplier: {}".format(self.reward_period_default_multiplier))
        log.info("reward_period_distribute: {}".format(self.reward_period_distribute))
        log.info("reward_period_mod: {}".format(self.reward_period_mod))
        log.info("reward_period_pool_count: {}".format(self.reward_period_pool_count))
        log.info("test_duration_blocks: {}".format(self.test_duration_blocks))

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

        self.log_level = "debug"
        self.validator0_mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow".split()

    def run(self):
        self.prj.pkill()
        time.sleep(1)

        import socket
        hostname = socket.gethostname()
        ip_address = socket.gethostbyname(hostname)

        self.sifnoded = []
        self.node_info = []
        for i in range(self.number_of_nodes):
            ports = self.ports_for_node(i)
            home = os.path.join(self.sifnoded_home_root, "sifnoded-{}".format(i))
            sifnoded_i = sifchain.Sifnoded(self.cmd, node=sifchain.format_node_url(ANY_ADDR, ports["rpc"]),
                home=home, chain_id=self.chain_id)
            moniker = "sifnode-{}".format(i)
            acct_name = "sif-{}".format(i)
            acct_addr = sifnoded_i.create_addr(acct_name, mnemonic=self.validator0_mnemonic if i == 0 else None)
            sifnoded_i.init(moniker)
            node_id = sifnoded_i.tendermint_show_node_id()  # Taken from ${sifnoded_home}/config/node_key.json
            pubkey = sifnoded_i.tendermint_show_validator()  # Taken from ${sifnoded_home}/config/priv_validator_key.json
            node_info = {
                "moniker": moniker,
                "home": home,
                "node_id": node_id,
                "pubkey": pubkey,
                "acct_name": acct_name,
                "acct_addr": acct_addr,
                "ports": ports,
                "external_address": sifchain.format_node_url(LOCALHOST, ports["rpc"])  # For --node
            }
            self.sifnoded.append(sifnoded_i)
            self.node_info.append(node_info)

        # Set up admin account balances. We add these with "add-genesis-account"
        # TODO It is not clear if we really need to fund all of them (and how much).
        # TODO Does this have to cover for rewards and lppd distribution? If rewards are minted, then no.
        # For rewards, the funds are minted and in case we opted for a distribution of the rewards to the LP wallet the
        # minted rowans are transferred there, you can see the minting process here: https://github.com/Sifchain/sifnode/blob/master/x/clp/keeper/rewards.go#L54
        # For LPD, we only transfer the existing funds in CLP to the LP's wallet, you can see here: https://github.com/Sifchain/sifnode/blob/8b2f9c45130c79e07555735185fbe1d00279fab0/x/clp/keeper/pool.go#L128
        denom_total_supply = 10000 * self.number_of_wallets * self.amount_of_denom_per_wallet
        validator_account_balance = cosmos.balance_add({
            ROWAN: 999 * 10**30,
            STAKE: 999 * 10**30,
        }, {denom: denom_total_supply for denom in self.tokens})

        sifnoded0 = self.sifnoded[0]

        for i in range(self.number_of_nodes):
            acct_addr = self.node_info[i]["acct_addr"]
            acct_bech = self.sifnoded[i].get_val_address(acct_addr)
            sifnoded0.add_genesis_account(acct_addr, validator_account_balance)
            sifnoded0.add_genesis_validators(acct_bech)
        admin0_addr = self.node_info[0]["acct_addr"]
        admin0_name = self.node_info[0]["acct_name"]
        sifnoded0.add_genesis_clp_admin(admin0_addr)
        sifnoded0.set_genesis_oracle_admin(admin0_name)
        sifnoded0.set_genesis_whitelister_admin(admin0_name)

        genesis_balances_to_add = {}

        # Set up test wallets with test tokens. We do this in genesis for performance reasons. For each wallet we choose
        # number_of_denoms_per_wallet` random denoms.
        client_home = os.path.join(self.sifnoded_home_root, "sifnoded-client")
        sifnoded_client = sifchain.Sifnoded(self.cmd, home=client_home, node=self.node_info[0]["external_address"],
            chain_id=self.chain_id)

        wallets = {}
        for i in range(self.number_of_wallets):
            chosen_tokens = [self.tokens[i] for i in random_choice(self.liquidity_providers_per_wallet, len(self.tokens), rnd=self.rnd)]
            balances = {denom: self.amount_of_denom_per_wallet for denom in chosen_tokens}
            addr = sifnoded_client.create_addr()
            wallets[addr] = chosen_tokens
            genesis_balances_to_add[addr] = cosmos.balance_add(balances, {ROWAN: self.amount_of_rowan_per_wallet})

        genesis = sifnoded0.load_genesis_json()
        sifnoded0.add_accounts_to_existing_genesis(genesis, genesis_balances_to_add)

        app_state = genesis["app_state"]
        app_state["gov"]["voting_params"] = {"voting_period": "120s"}
        app_state["gov"]["deposit_params"]["min_deposit"] = [{"denom": ROWAN, "amount": "10000000"}]
        app_state["crisis"]["constant_fee"] = {"denom": ROWAN, "amount": "1000"}
        app_state["staking"]["params"]["bond_denom"] = ROWAN
        app_state["mint"]["params"]["mint_denom"] = ROWAN
        sifnoded0.save_genesis_json(genesis)

        # sifnoded0.gentx(admin0_name, {STAKE: 10**24})
        sifnoded0.gentx(admin0_name, {ROWAN: 10**24})
        sifnoded0.collect_gentx()
        sifnoded0.validate_genesis()

        # According to gzukel, nodes need just one peer to make sync work.
        peers = [sifchain.format_peer_address(node_info["node_id"], LOCALHOST, node_info["ports"]["p2p"])
            for node_info in [self.node_info[0]]]
        genesis = sifnoded0.load_genesis_json()
        for i in range(self.number_of_nodes):
            sifnoded_i = self.sifnoded[i]
            if i != 0:
                sifnoded_i.save_genesis_json(genesis)  # Copy genesis from validator 0 to all other
            info = self.node_info[i]
            app_toml = sifnoded_i.load_app_toml()
            app_toml["minimum-gas-prices"] = sif_format_amount(0.5, ROWAN)
            app_toml['api']['enable'] = True
            app_toml["api"]["address"] = sifchain.format_node_url(ANY_ADDR, info["ports"]["api"])
            sifnoded_i.save_app_toml(app_toml)
            config_toml = sifnoded_i.load_config_toml()
            config_toml["log_level"] = self.log_level  # TODO Probably redundant
            config_toml['p2p']["external_address"] = "{}:{}".format(ip_address, info["ports"]["p2p"])
            if i != 0:
                config_toml["p2p"]["persistent_peers"] = ",".join(peers)
            config_toml['p2p']['max_num_inbound_peers'] = 50
            config_toml['p2p']['max_num_outbound_peers'] = 50
            config_toml['p2p']['allow_duplicate_ip'] = True
            config_toml["rpc"]["pprof_laddr"] = "{}:{}".format(LOCALHOST, info["ports"]["pprof"])
            config_toml['moniker'] = info["moniker"]
            sifnoded_i.save_config_toml(config_toml)

        # Start processes
        processes = []
        log_files = []
        for i, sifnoded_i in enumerate(self.sifnoded):
            node_info = self.node_info[i]
            ports = node_info["ports"]
            log_file_path = os.path.join(sifnoded_i.home, "sifnoded.log")
            log_file = open(log_file_path, "w")
            log_files.append(log_file)
            process = sifnoded_i.sifnoded_start(log_file=log_file, log_level="debug", trace=True,
                tcp_url="tcp://{}:{}".format(ANY_ADDR, ports["rpc"]), p2p_laddr="{}:{}".format(ANY_ADDR, ports["p2p"]),
                grpc_address="{}:{}".format(ANY_ADDR, ports["grpc"]),
                grpc_web_address="{}:{}".format(ANY_ADDR, ports["grpc_web"]),
                address="tcp://{}:{}".format(ANY_ADDR, ports["address"])
            )
            sifnoded_i._wait_up()
            processes.append(process)

        # Wait for some time so that nodes are fully booted
        sifnoded0.wait_for_last_transaction_to_be_mined()

        # Create a validator for all non-0 nodes. Node 0 needs to be up, but node i may or may not be up.
        for i in [x for x in range(self.number_of_nodes) if x != 0]:
            node_info = self.node_info[i]
            # This needs to have the private key ("home") of i-th validator but "node" of the 0-th.
            # TODO We need to use "rpc" for --node, not p2p / external_address!
            sifnoded_tmp = sifchain.Sifnoded(self.cmd, home=node_info["home"], chain_id=self.chain_id,
                node=self.node_info[0]["external_address"])
            sifnoded_tmp.staking_create_validator((10 ** 24, ROWAN), node_info["pubkey"], node_info["moniker"],
                0.10, 0.20, 0.01, 1000000, node_info["acct_addr"])

        sifnoded0.wait_for_last_transaction_to_be_mined()

        # On each node, do a sample transfer of one rowan from admin to a new wallet and check that the change of
        # balances is visible on all nodes
        test_transfer_amount = {ROWAN: 10**ROWAN_DECIMALS}
        for i in range(self.number_of_nodes):
            test_addr = self.sifnoded[i].create_addr()
            self.sifnoded[i].send(self.node_info[i]["acct_addr"], test_addr, test_transfer_amount)
            for j in range(self.number_of_nodes):
                self.sifnoded[j].wait_for_balance_change(test_addr, {}, test_transfer_amount)

        assert all(len(self.sifnoded[i].query_staking_validators()) == self.number_of_nodes
            for i in range(self.number_of_nodes))

        # Check balances
        assert all(cosmos.balance_equal(sifnoded0.get_balance(addr), genesis_balances_to_add[addr])
            for addr in genesis_balances_to_add)
        assert all(sifnoded0.get_acct_seq(addr)[1] == 0 for addr in wallets)

        sifnoded = self.sifnoded[0]
        sif = self.node_info[0]["acct_addr"]

        # Add tokens to token registry. The minimum required permissions are  CLP.
        # TODO Might want to use `tx tokenregistry set-registry` to do it in one step (faster)
        #      According to @sheokapr `tx tokenregistry set-registry` also works for only one entry
        #      But`tx tokenregistry register-all` also works only for one entry.
        for denom in self.tokens:
            entry = sifnoded.create_tokenregistry_entry(denom, denom, 18, ["CLP"])
            sifnoded.token_registry_register(entry, sif)
            sifnoded.wait_for_last_transaction_to_be_mined()  # Must be run synchronously! (if not, only the first will work)
        assert set(e["denom"] for e in sifnoded.query_tokenregistry_entries()) == set(self.tokens)

        # Set up liquidity pools. We create them symmetrically (`native_amount == externam_amount`).
        for denom in self.tokens:
            sifnoded.tx_clp_create_pool(sif, denom, self.amount_of_denom_per_wallet, self.amount_of_denom_per_wallet)
            sifnoded.wait_for_last_transaction_to_be_mined()
        assert set(p["external_asset"]["symbol"] for p in sifnoded.query_pools()) == set(self.tokens)

        # Set up liquidity providers. We create them symmetrically (`native_amount == externam_amount`). The ratio of
        # native vs. external amount has to be the same as for corresponding pool (within certain rounding tolerance).
        # Calling `tx_clp_add_liquidity` to add multiple liquidity providers within the same block does not work (only
        # the first call gets through). To avoid `--broadcast-mode block` or waiting for new block, we need to use
        # account sequence numbers.
        for addr, denoms in wallets.items():
            account_number, account_sequence = sifnoded.get_acct_seq(addr)
            for denom in denoms:
                sifnoded_client.tx_clp_add_liquidity(addr, denom, self.amount_of_liquidity_added_by_wallet,
                    self.amount_of_liquidity_added_by_wallet, account_seq=(account_number, account_sequence))
                account_sequence += 1
        sifnoded_client.wait_for_last_transaction_to_be_mined()
        actual_lp_providers = {}
        for denom in self.tokens:
            for lp in sifnoded_client.query_clp_liquidity_providers(denom):
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
        # TODO start and end blocks are both inclusive, adjust
        current_block = sifnoded_client.get_current_block()
        start_block = current_block + 5
        rewards_start_block = start_block + self.test_duration_blocks
        rewards_end_block = rewards_start_block + 2 * self.test_duration_blocks
        lppd_start_block = start_block + 2 * self.test_duration_blocks
        lppd_end_block = lppd_start_block + 2 * self.test_duration_blocks

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
        log.info("In phase 'neither' (blocks {} - {})".format(start_block, rewards_start_block))
        time1 = self.wait_for_block(rewards_start_block)
        log.info("In phase 'rewards only' (blocks {} - {})".format(rewards_start_block, lppd_start_block))
        time2 = self.wait_for_block(lppd_start_block)
        log.info("In phase 'rewards + LPPD' (blocks {} - {})".format(lppd_start_block, rewards_end_block))
        time3 = self.wait_for_block(rewards_end_block)
        log.info("In phase 'LPPD only' (blocks {} - {})".format(rewards_end_block, lppd_end_block))
        time4 = self.wait_for_block(lppd_end_block)

        accuracy = 1.0 / self.test_duration_blocks

        log.info("Neither: {:.2f} +/- {:.2f} s/block".format((time1 - time0) / self.test_duration_blocks, accuracy))
        log.info("Rewards only: {:.2f} +/- {:.2f} s/block".format((time2 - time1) / self.test_duration_blocks, accuracy))
        log.info("Rewards + LPPD: {:.2f} +/- {:.2f} s/block".format((time3 - time2) / self.test_duration_blocks, accuracy))
        log.info("LPPD only: {:.2f} +/- {:.2f} s/block".format((time4 - time3) / self.test_duration_blocks, accuracy))

        # TODO LPPD and rewards assertions
        # See https://www.notion.so/sifchain/Rewards-2-0-Load-Testing-972fbe73b04440cd87232aa60a3146c5#7392be2c1a034d2db83b9b38ab89ff9e

        return processes, log_files

    # TODO Refactor - move to Sifnoded
    def wait_for_block(self, block_number: int) -> float:
        sifnoded = self.sifnoded[0]
        current_block = sifnoded.get_current_block()
        prev_block = None
        assert current_block < block_number
        while current_block < block_number:
            if (prev_block is None) or (current_block != prev_block):
                # This is just for collecting statistics while we wait, the test result does not depend on it.
                # Check also https://github.com/cosmos/cosmos-sdk/issues/6105
                try:
                    blk = sifnoded.get_block_results()
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

    def ports_for_node(self, i: int) -> JsonDict:
        return {
            "p2p": 10276 + i,
            "grpc": 10909 + i,
            "grpc_web": 10919 + i,
            "address": 10276 + i,
            "rpc": 10286 + i,
            "api": 10131 + i,
            "pprof": 10606 + i,
        }


def main(argv: List[str]):
    basic_logging_setup()
    parser = argparse.ArgumentParser()
    parser.add_argument("--number-of-nodes", type=int, default=1)
    parser.add_argument("--number-of-liquidity-pools", type=int, default=10)
    parser.add_argument("--number-of-wallets", type=int, default=10)
    parser.add_argument("--liquidity-providers-per-wallet", type=int, default=5)
    parser.add_argument("--reward-period-default-multiplier", type=float, default=0.0)
    parser.add_argument("--reward-period-distribute", action="store_true")
    parser.add_argument("--reward-period-mod", type=int, default=1)
    parser.add_argument("--reward-period-pool-count", type=int, default=10)
    parser.add_argument("--test-duration-blocks", type=int, default=5)
    args = parser.parse_args(argv[1:])

    cmd = command.Command()
    prj = project.Project(cmd, project_dir())
    sifnoded_home_root = cmd.tmpdir("siftool-test.tmp")
    cmd.rmdir(sifnoded_home_root)
    test = Test(cmd, prj, sifnoded_home_root=sifnoded_home_root)

    test.number_of_nodes = args.number_of_nodes
    test.number_of_liquidity_pools = args.number_of_liquidity_pools
    test.number_of_wallets = args.number_of_wallets
    test.liquidity_providers_per_wallet = args.liquidity_providers_per_wallet
    test.reward_period_default_multiplier = args.reward_period_default_multiplier
    test.reward_period_distribute = args.reward_period_distribute
    test.reward_period_mod = args.reward_period_mod
    test.reward_period_pool_count = args.reward_period_pool_count
    test.test_duration_blocks = args.test_duration_blocks

    test.setup()

    try:
        test_start_time = time.time()
        test.run()
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


if __name__ == "__main__":
    main(sys.argv)
