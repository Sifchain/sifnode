from contextlib import contextmanager
from siftool.common import *
from siftool.sifchain import ROWAN, ROWAN_DECIMALS
from siftool import command, environments, project, sifchain, cosmos


def get_validators(env):
    sifnoded = env._sifnoded_for(env.node_info[0])
    return {v["description"]["moniker"]: v for v in sifnoded.query_staking_validators()}


def test_transfer(env):
    sifnoded = env.sifnoded
    sifnoded.send_and_check(env.faucet, sifnoded.create_addr(), {sifchain.ROWAN: 10 ** sifchain.ROWAN_DECIMALS})


def assert_validators_working(env, expected_monikers):
    assert set(get_validators(env)) == expected_monikers
    for i in range(len(env.node_info)):
        test_transfer(env)


def swap_pricing_formula(x, X, y, Y):
    # Pricing for swaps:
    # y = x*Y / (x + X)
    # x = y*X / (y + Y)
    # x = swap input or output (non-rowan)
    # y = swap input or output (rowan)
    # X = non-rowan token pool balance
    # Y = rowan pool balance
    # Note: a "swap" is subjected to price calculation of all swaps that happen in that block, so the final prices
    # are not known until after swaps have been processed.
    raise NotImplemented()


class TestMargin:
    def setup_method(self):
        self.default_pool_setup = {
            # denom,  decimals, faucet balance, pool native amount, pool external amount
            "cusdc": (6,        10**30,         10**25,             10**25),
            "cdash": (18,       10**30,         10**25,             10**25),
            "ceth":  (18,       10**30,         10**25,             10**25),
        }

        self.cmd = command.Command()
        self.sifnoded_home_root = self.cmd.tmpdir("siftool.tmp")
        self.cmd.rmdir(self.sifnoded_home_root)
        self.cmd.mkdir(self.sifnoded_home_root)
        prj = project.Project(self.cmd, project_dir())
        prj.pkill()

    def teardown_method(self):
        prj = project.Project(self.cmd, project_dir())
        prj.pkill()

    @contextmanager
    def with_test_env_with_tokens_and_pools(self, pool_definitions):
        faucet_balance = cosmos.balance_add({
            sifchain.ROWAN: 10**30,
            sifchain.STAKE: 10**30,
        }, {
            denom: vals[1] for denom, vals in pool_definitions.items()
        })

        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        env.add_validator()
        env.init(faucet_balance=faucet_balance)
        env.start()

        sifnoded = env.sifnoded
        sifnoded.send_batch(env.faucet, env.clp_admin, cosmos.balance_add(
            {denom: vals[3] for denom, vals in pool_definitions.items()},
            {sifchain.ROWAN: sum(vals[2] for _, vals in pool_definitions.items())}))

        # We need to register rowan in token registry, otherwise swaps from/to rowan will error out with
        # "token not supported by sifchain"
        # Note: original_entry = {
        #     "decimals": str(ROWAN_DECIMALS),
        #     "denom": ROWAN,
        #     "base_denom": ROWAN,
        #     "permissions": [1]
        # }
        sifnoded.token_registry_register(sifnoded.create_tokenregistry_entry(ROWAN, ROWAN, ROWAN_DECIMALS, ["CLP"]),
            env.clp_admin, broadcast_mode="block")

        sifnoded.token_registry_register_batch(env.clp_admin,
            [sifnoded.create_tokenregistry_entry(denom, denom, vals[0], ["CLP"]) for denom, vals in pool_definitions.items()])
        sifnoded.create_liquidity_pools_batch(env.clp_admin,
            [(denom, vals[2], vals[3]) for denom, vals in pool_definitions.items()])

        pools = env.sifnoded.query_pools()
        assert len(pools) == len(pool_definitions)

        yield env

        env.close()

    @contextmanager
    def with_default_env(self):
        with self.with_test_env_with_tokens_and_pools(self.default_pool_setup) as env:
            yield env, env.sifnoded, self.default_pool_setup

    def test_swap(self):

        def test_rowan_external():
            with self.with_default_env() as (env, sifnoded, pools_definitions):
                src_denom = "rowan"
                dst_denom = "ceth"
                amount = 10**18

                account = sifnoded.create_addr()
                env.fund(account, {ROWAN: 10**25})
                balance_before = sifnoded.get_balance(account)
                pools_before = sifnoded.query_pools_sorted()
                res = sifnoded.tx_clp_swap(account, src_denom, amount, dst_denom, 0, broadcast_mode="block")
                balance_after = sifnoded.get_balance(account)
                pools_after = sifnoded.query_pools_sorted()

                src_balance_delta = balance_after.get(src_denom, 0) - balance_before.get(src_denom, 0)
                dst_balance_delta = balance_after.get(dst_denom, 0) - balance_before.get(dst_denom, 0)

                assert balance_before.get(dst_denom, 0) == 0
                assert balance_after.get(dst_denom, 0) > 0
                
                # Account rowan balance should change by 1 tx fee + amount swapped
                assert src_balance_delta == - (sifchain.sif_tx_fee_in_rowan + amount)

                # Pool native asset balance should increase by swapped amount
                assert int(pools_after[dst_denom]["native_asset_balance"]) - int(pools_before[dst_denom]["native_asset_balance"]) == amount

                # Pool external asset balance should decrease by amount received
                assert int(pools_after[dst_denom]["external_asset_balance"]) - int(pools_before[dst_denom]["external_asset_balance"]) == -dst_balance_delta

                # y = x*Y / (x + X)
                # x = y*X / (y + Y)
                # x = swap input or output (non-rowan)
                # y = swap input or output (rowan)
                # X = non-rowan token pool balance
                # Y = rowan pool balance

                y = amount
                X = int(pools_before[dst_denom]["external_asset_balance"])
                Y = int(pools_before[dst_denom]["native_asset_balance"])
                expected_dst_delta = amount * int(pools_before[dst_denom]["external_asset_balance"]) / (amount + int(pools_before[dst_denom]["native_asset_balance"]))
                expected_dst_delta2 = amount * int(pools_after[dst_denom]["external_asset_balance"]) / (amount + int(pools_after[dst_denom]["native_asset_balance"]))

                x = y*X / (y + Y)

                assert True  # no-op line just for setting a breakpoint
        test_rowan_external()


        def test_external_external():
            with self.with_default_env() as (env, sifnoded, pools_definitions):
                account = sifnoded.create_addr()
                balances = []
                pools = []

                from_denom = "cusdc"
                to_denom = "ceth"
                amount1 = 10**23

                balances.append(sifnoded.get_balance(account))
                pools.append(sifnoded.query_pools_sorted())
                env.fund(account, {ROWAN: 10**25, "cusdc": amount1})

                balance_before = sifnoded.get_balance(account)
                pools_before = sifnoded.query_pools_sorted()
                res = sifnoded.tx_clp_swap(account, from_denom, amount1, to_denom, 0, broadcast_mode="block")
                balance_after = sifnoded.get_balance(account)
                pools_after = sifnoded.query_pools_sorted()

                rowan_delta = balance_after.get(ROWAN, 0) - balance_before.get(ROWAN, 0)
                from_delta = balance_after.get(from_denom, 0) - balance_before.get(from_denom, 0)
                to_delta = balance_after.get(to_denom, 0) - balance_before.get(to_denom, 0)

                # Before swap, the account should have expected amount of from_denom and zero to_denom
                assert balance_before.get(from_denom, 0) == amount1
                assert balance_before.get(to_denom, 0) == 0

                # The account should pay a transaction fee in rowan of 10**17
                assert rowan_delta == -sifchain.sif_tx_fee_in_rowan

                # Source pool's external asset should increase by swapped amount
                assert int(pools_after[from_denom]["external_asset_balance"]) - int(pools_before[from_denom]["external_asset_balance"]) == amount1

                # Destination pool's external asset should decrease by amount received
                assert int(pools_after[to_denom]["external_asset_balance"]) - int(pools_before[to_denom]["external_asset_balance"]) == -to_delta

                from_pool_native_balance_delta = int(pools_after[from_denom]["native_asset_balance"]) - int(pools_before[from_denom]["native_asset_balance"])
                to_pool_native_balance_delta = int(pools_after[to_denom]["native_asset_balance"]) - int(pools_before[to_denom]["native_asset_balance"])

                # Source pool native balance should decrease
                assert from_pool_native_balance_delta < 0

                # Destination pool  native balance should increase by the same amount
                assert to_pool_native_balance_delta == -from_pool_native_balance_delta

                pools.append(pools_before)
                pools.append(pools_after)
                balances.append(balance_before)
                balances.append(balance_after)
