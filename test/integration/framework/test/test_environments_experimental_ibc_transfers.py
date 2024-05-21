from siftool.common import *
from siftool import command, environments, project


# The goal of this test is to setup an local environment for testing IBC transfers.
# We are trying to spin up single-node sifnoded validator + single-node axelard validator, and then setup and run
# and Osmosis relayer ("rly") to relay IBC tokens between the chains.
# In addition to sifnoded, axelard and rly must be in PATH.
# TODO Not working yet
class TestSifnodedEnvironmentWithIBCTransfers:
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

    def test_ibc_setup_and_transfers(self):
        env = environments.SifnodedEnvironment(self.cmd, sifnoded_home_root=self.sifnoded_home_root)
        env.add_validator()
        env.start()

        axelard = Axelard(home=os.path.join(env.sifnoded_home_root + "axelard"))


class Axelard:
    def __init__(self, /, chain_id: Optional[str] = None, binary: Optional[str] = None, home: Optional[str] = None):
        self.chain_id = chain_id or "axelar-localnet"
        self.binary = binary or "axelard"
        self.home = home
        self.ports = None

    def start(self):
        pass
        # For generic Cosmos-based chain:
        # - <binary> keys add <admin-name>
        # - <binary> init <moniker>
        # - <binary> add-genesis-account ...
        # - <binary> add-genesis-validator ... <--- no such command in axelard
        # - <binary> gentx <admin-name> <stake-amount> --home <home> --keyring-backend test --chain-id <chain-id>
        # - <binary> collect-gentxs --home <home>
        # - <binary> validate-genesis
