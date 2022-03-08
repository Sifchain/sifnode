import os
from command import Command
from common import project_dir

class Localnet(Command):
    def __init__(self, script_dir=None, config_dir=None, bin_dir=None):
        self.script_dir = script_dir if script_dir else project_dir("test/localnet")
        self.config_dir = config_dir if config_dir else os.path.join("/tmp/localnet", "./config")
        self.bin_dir = bin_dir if bin_dir else os.path.join("/tmp/localnet", "./bin")
        self.node_module_dir = os.path.join(self.script_dir, "./node_modules")

    @staticmethod
    def is_enabled():
        return 'LOCALNET' in os.environ and os.environ['LOCALNET'] == "true"

    def install_deps(self):
        self.execst(["yarn"], cwd=self.script_dir)

    def download_binaries(self):
        self.execst(["yarn", "downloadBinaries"], cwd=self.script_dir)

    def init_all_chains(self):
        self.execst(["yarn", "initAllChains"], cwd=self.script_dir)

    def start_all_chains(self):
        self.execst(['yarn', 'startAllChains'], cmd=self.script_dir)

    def init_all_relayers(self):
        self.execst(['yarn', 'initAllRelayers'], cmd=self.script_dir)

    def start_all_relayers(self):
        self.execst(['yarn', 'startAllRelayers'], cmd=self.script_dir)

    def build_local_net(self):
        self.execst(['yarn', 'buildLocalNet'], cmd=self.script_dir)

    def loal_local_net(self):
        self.execst(['yarn', 'loadLocalNet'], cmd=self.script_dir)

    def take_snapshot(self):
        self.execst(['yarn', 'takeSnapshot'], cmd=self.script_dir)

    def create_snapshot(self):
        self.execst(['yarn', 'createSnapshot'], cmd=self.script_dir)
        
    def test(self):
        self.execst(['yarn', 'test'], cmd=self.script_dir)
        