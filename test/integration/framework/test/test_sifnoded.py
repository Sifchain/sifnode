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
        denoms = ["test-{}".format(i) for i in range(1000)]
        try:
            sifnoded = sifchain.Sifnoded(self.cmd, home=tmpdir)
            addr = sifnoded.create_addr()

            env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
            balances = {denom: 10**18 for denom in denoms}
            env.add_validator(initial_balance=cosmos.balance_add({sifchain.ROWAN: 10**30}, balances))
            env.init(extra_accounts={addr: balances})
            env.start()

            sifnoded = sifchain.Sifnoded(self.cmd, home=tmpdir, chain_id=env.chain_id,
                node=sifchain.format_node_url(env.node_info[0]["host"], env.node_info[0]["ports"]["rpc"]))
            actual_balances = sifnoded.get_balance(addr)
            assert actual_balances == balances

            # Create pools
            sifnoded = env._sifnoded_for(env.node_info[0])
            sifnoded.disable_assertions = True  # Bug: "sifnoded query tokenregistry entries" does not return more than 103 entries
            admin = env.node_info[0]["admin_addr"]
            sifnoded.token_registry_register_batch(admin,
                tuple(sifnoded.create_tokenregistry_entry(denom, denom, 18, ["CLP"]) for denom in denoms))

            sifnoded.create_liquidity_pools_batch(admin,
                tuple((denom, 10**18, 10**18) for denom in denoms))

            assert set(p["external_asset"]["symbol"] for p in sifnoded.query_pools()) == set(denoms)
        finally:
            self.cmd.rmdir(tmpdir)
