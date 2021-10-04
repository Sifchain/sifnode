import json
from common import *


def force_kill_processes(cmd):
    cmd.execst(["pkill", "node"], check_exit=False)
    cmd.execst(["pkill", "ebrelayer"], check_exit=False)
    cmd.execst(["pkill", "sifnoded"], check_exit=False)

def killall(processes):
    # TODO Order - ebrelayer, sifnoded, ganache
    for p in processes:
        if p is not None:
            p.kill()
            p.wait()


class Project:
    """Represents a checked out copy of a project in a particular directory."""

    def __init__(self, cmd, base_dir):
        self.cmd = cmd
        self.base_dir = base_dir
        self.smart_contracts_dir = project_dir("smart-contracts")
        self.go_bin_dir = os.environ["GOBIN"]

    def project_dir(self, *paths):
        return os.path.abspath(os.path.join(self.base_dir, *paths))

    def __rmdir(self, path):
        if self.cmd.exists(path):
            log.info("Removing '{}'...".format(path))
            self.cmd.rmdir(path)
        else:
            log.info("Nothing to delete for '{}'".format(path))

    def cleanup_and_reset_state(self):
        force_kill_processes(self.cmd)

        # rm -rvf /tmp/tmp.xxxx (ganache DB, unique for every run)
        self.__rmdir(self.project_dir("test", "integration", "sifchainrelayerdb"))  # TODO move to /tmp
        self.__rmdir(self.project_dir("smart-contracts", "build"))  # truffle deploy
        self.__rmdir(self.project_dir("test", "integration", "vagrant", "data"))
        self.__rmdir(self.cmd.get_user_home(".sifnoded"))  # Probably needed for "--keyring-backend test"

        self.__rmdir(self.project_dir("deploy", "networks"))  # from running integration tests

        # Peggy/devenv/hardhat cleanup
        # For full clean, also: cd smart-contracts && rm -rf node_modules && npm install
        # TODO Difference between yarn vs. npm install?
        # (1) = cd smart-contracts; npx hardhat run scripts/deploy_contracts.ts --network localhost
        # (2) = cd smart-contracts; GOBIN=/home/anderson/go/bin npx hardhat run scripts/devenv.ts
        self.__rmdir(self.project_dir("smart-contracts", "build"))  # (1)
        self.__rmdir(self.project_dir("smart-contracts", "artifacts"))  # (1)
        self.__rmdir(self.project_dir("smart-contracts", "cache"))  # (1)
        self.__rmdir(self.project_dir("smart-contracts", ".openzeppelin"))  # (1)
        self.__rmdir(self.project_dir("smart-contracts", "relayerdb"))  # (2)
        self.__rmdir(self.project_dir("smart-contracts", "venv"))

        # Additional cleanup (not neccessary to make it work)
        # self.cmd.rm(self.project_dir("smart-contracts/combined.log"))
        # self.cmd.rmdir(self.project_dir("test/integration/.pytest_cache"))
        # self.cmd,rm(self.project_dir("smart-contracts/.env"))
        # self.cmd.rmdir(self.project_dir("deploy/networks"))
        # self.cmd.rmdir(self.project_dir("smart-contracts/.openzeppelin"))

        # rmdir ~/.cache/yarn

    def fullclean(self):
        go_home = self.cmd.get_user_home("go")
        self.cmd.execst(["chmod", "-R", "+w", go_home])
        self.__rmdir(go_home)
        self.__rmdir(self.cmd.get_user_home("go"))
        self.__rmdir(self.cmd.get_user_home(".npm"))
        self.__rmdir(self.cmd.get_user_home(".npm-global"))
        self.__rmdir(self.cmd.get_user_home(".npm-global"))
        self.__rmdir(self.cmd.get_user_home(".cache/yarn"))
        self.__rmdir(self.cmd.get_user_home(".sifnoded"))
        self.__rmdir(self.cmd.get_user_home(".sifnode-integration"))
        self.__rmdir(project_dir("smart-contracts/node_modules"))
        self.cmd.execst(["npm", "install", "-g", "ganache-cli", "dotenv", "yarn"], cwd=self.smart_contracts_dir)
        self.install_smart_contracts_dependencies()

    def yarn(self, args, cwd=None, env=None):
        return self.cmd.execst(["yarn"] + args, cwd=cwd, env=env, pipe=False)

    def npx(self, args, env=None, cwd=None, pipe=True):
        # Typically we want any npx commands to inherit stdout and strerr
        return self.cmd.execst(["npx"] + args, env=env, cwd=cwd, pipe=pipe)

    # IntegrationEnvironment
    # TODO Merge
    def make_go_binaries(self):
        # make go binaries (TODO Makefile needs to be trimmed down, especially "find")
        # cd test/integration; BASEDIR=... make
        # (checks all *.go files and, runs make in $BASEDIR, touches sifnoded, removes ~/.sifnoded/localnet
        self.cmd.execst(["make"], cwd=project_dir("test", "integration"), env={"BASEDIR": project_dir()}, pipe=False)

    # From PeggyEnvironment
    # TODO Merge
    # Main Makefile requires GOBIN to be set to an absolute path. Compiled executables ebrelayer, sifgen and
    # sifnoded will be written there. The directory will be created if it doesn't exist yet.
    #
    def make_go_binaries_2(self):
        # Original: cd smart-contracts; make -C .. install
        self.cmd.execst(["make", "install"], cwd=project_dir(), pipe=False)

    def install_smart_contracts_dependencies(self):
        self.cmd.execst(["make", "clean-smartcontracts"], cwd=self.smart_contracts_dir)  # = rm -rf build .openzeppelin
        # According to peggy2, the plan is to move from npm install to yarn, but there are some issues with yarn atm.
        # self.yarn(["install"], cwd=self.smart_contracts_dir)
        self.cmd.execst(["npm", "install"], cwd=self.smart_contracts_dir, pipe=False)

    def write_vagrantenv_sh(self, state_vars, data_dir, ethereum_websocket_address, chainnet):
        # Trace of test_utilities.py get_required_env_var/get_optional_env_var:
        #
        # BASEDIR (required), value=/home/jurez/work/projects/sif/sifnode/local
        # BRIDGE_BANK_ADDRESS (optional), value=0x30753E4A8aad7F8597332E813735Def5dD395028
        # BRIDGE_BANK_ADDRESS (required), value=0x30753E4A8aad7F8597332E813735Def5dD395028
        # BRIDGE_REGISTRY_ADDRESS (required), value=0xf204a4Ef082f5c04bB89F7D5E6568B796096735a
        # BRIDGE_TOKEN_ADDRESS (optional), value=0x82D50AD3C1091866E258Fd0f1a7cC9674609D254
        # BRIDGE_TOKEN_ADDRESS (required), value=0x82D50AD3C1091866E258Fd0f1a7cC9674609D254
        # CHAINDIR (required), 3x value
        # CHAINNET (required), value=localnet
        # DEPLOYMENT_NAME (optional), value=None
        # ETHEREUM_ADDRESS (optional), value=None
        # ETHEREUM_NETWORK (optional), value=None
        # ETHEREUM_NETWORK_ID (optional), value=None
        # GANACHE_KEYS_FILE (optional), value=None
        # HOME (required), value=/home/jurez
        # MNEMONIC (required), value=future tattoo gesture artist tomato accuse chuckle polar ivory strategy rail flower apart virus burger rhythm either describe habit attend absurd aspect predict parent
        # MONIKER (required), value=wandering-flower
        # OPERATOR_ADDRESS (optional), value=None
        # OPERATOR_PRIVATE_KEY (optional), value=None
        # OPERATOR_PRIVATE_KEY (optional), value=c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
        # ROWAN_SOURCE (optional), value=None
        # ROWAN_SOURCE_KEY (optional), value=None
        # SIFCHAIN_ADMIN_ACCOUNT (required), value=sif1896ner48vrg8m05k48ykc6yydlxc4yvm23hp5m
        # SIFNODE (optional), value=None
        # SMART_CONTRACTS_DIR (required), 2x value
        # SMART_CONTRACT_ARTIFACT_DIR (optional), value=None
        # SOLIDITY_JSON_PATH (optional), value=None
        # TEST_INTEGRATION_DIR (required), value=/home/jurez/work/projects/sif/sifnode/local/test/integration
        # VALIDATOR1_ADDR (optional), 3x value
        # VALIDATOR1_PASSWORD (optional), 3x value
        env = dict_merge(state_vars, {
            # For running test/integration/execute_integration_tests_against_*.sh
            "TEST_INTEGRATION_DIR": project_dir("test/integration"),
            "TEST_INTEGRATION_PY_DIR": project_dir("test/integration/src/py"),
            "SMART_CONTRACTS_DIR": self.cmd.smart_contracts_dir,
            "datadir": data_dir,  # Needed by test_rollback_chain.py that calls ganache_start.sh
            "GANACHE_KEYS_JSON": os.path.join(data_dir, "ganachekeys.json"),  # Needed by test_rollback_chain.py that calls ganache_start.sh
            "ETHEREUM_WEBSOCKET_ADDRESS": ethereum_websocket_address,   # Needed by test_ebrelayer_replay.py (and possibly others)
            "CHAINNET": chainnet,   # Needed by test_ebrelayer_replay.py (and possibly others)
        })
        vagrantenv_path = project_dir("test/integration/vagrantenv.sh")
        self.cmd.write_text_file(vagrantenv_path, joinlines(format_as_shell_env_vars(env)))
        self.cmd.write_text_file(project_dir("test/integration/vagrantenv.json"), json.dumps(env))

    def init(self):
        self.cleanup_and_reset_state()
        self.cmd.rmdir(project_dir("smart-contracts/node_modules"))
        self.make_go_binaries_2()
        self.install_smart_contracts_dependencies()
