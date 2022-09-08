from siftool.common import *
from siftool import command, cosmos, sifchain, project, environments


class TestSifnodedCLIWrapper:
    def setup_method(self):
        self.cmd = command.Command()
        self.sifnoded_home_root = self.cmd.tmpdir("siftool.tmp")
        self.cmd.rmdir(self.sifnoded_home_root)
        self.cmd.mkdir(self.sifnoded_home_root)
        prj = project.Project(self.cmd, project_dir())
        prj.pkill()

    def teardown_method(self):
        prj = project.Project(self.cmd, project_dir())
        prj.pkill()

    # We do two different reads - "query bank balances" and "query clp pools" since they use slightly ddi
    def test_batching_and_paged_reads(self):
        tmpdir = self.cmd.mktempdir()

        # Note: since all the denoms are passed by "initial_balance" they appear on the command line of "sifnoded gentx".
        # For more than 1000 denoms we might get an "OSError of "too many parameters".
        # TODO Apparently we need more than 5000 denoms to actually trigger the paging in "query bank balances"
        denoms = ["test-{}".format(i) for i in range(1000)]
        try:
            sifnoded = sifchain.Sifnoded(self.cmd, home=tmpdir)
            test_addr = sifnoded.create_addr()
            test_coins_balance = {denom: 10**18 for denom in denoms}
            test_addr_balance = cosmos.balance_add({sifchain.ROWAN: 10**30}, test_coins_balance)

            env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
            env.add_validator()
            env.init(extra_accounts={test_addr: test_addr_balance})
            env.start()

            validator0_admin = env.node_info[0]["admin_addr"]
            clp_admin = validator0_admin

            sifnoded = sifchain.Sifnoded(self.cmd, home=tmpdir, chain_id=env.chain_id,
                node=sifchain.format_node_url(env.node_info[0]["host"], env.node_info[0]["ports"]["rpc"]))
            test_addr_actual_balance = sifnoded.get_balance(test_addr)
            assert test_addr_actual_balance == test_addr_balance

            # Send from addr to clp_admin
            sifnoded.send_batch(test_addr, clp_admin, test_coins_balance)

            # Create pools
            sifnoded = env._sifnoded_for(env.node_info[0])
            sifnoded.token_registry_register_batch(clp_admin,
                tuple(sifnoded.create_tokenregistry_entry(denom, denom, 18, ["CLP"]) for denom in denoms))

            sifnoded.create_liquidity_pools_batch(clp_admin,
                tuple((denom, 10**18, 10**18) for denom in denoms))

            assert set(p["external_asset"]["symbol"] for p in sifnoded.query_pools()) == set(denoms)
        finally:
            self.cmd.rmdir(tmpdir)
