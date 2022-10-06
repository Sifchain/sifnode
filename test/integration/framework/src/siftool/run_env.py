import json
import os
import re
import time
from typing import List, Tuple, TextIO, Any

import siftool.eth
from siftool import eth, cosmos, command
from siftool.hardhat import Hardhat, default_accounts as hardhat_default_accounts
from siftool.geth import Geth
from siftool.truffle import Ganache
from siftool.localnet import Localnet
from siftool.command import Command
from siftool.sifchain import Sifgen, Sifnoded, Ebrelayer, sifchain_denom_hash, ROWAN
from siftool.project import Project, killall, force_kill_processes
from siftool.common import *


log = siftool_logger(__name__)


# Extension of Command class
# TODO This is now obsolete and should be refactored out
class Integrator(Ganache, Command):
    def __init__(self):
        super().__init__()  # TODO Which super is this? All of them?
        self.project = Project(self, project_dir())

    def primitive_parse_env_file(self, path):
        def split(lines):
            result = dict()
            for line in lines:
                m = patt.match(line)
                result[m[1]] = m[2]
            return result

        tmp1 = self.mktempfile()
        tmp2 = self.mktempfile()
        try:
            self.execst(["env", "-i", "bash", "-c", "set -o posix; IFS=' '; set > {}; source {}; set > {}".format(tmp1, path, tmp2)])
            t1 = self.read_text_file(tmp1).splitlines()
            t2 = self.read_text_file(tmp2).splitlines()
        finally:
            self.rm(tmp1)
            self.rm(tmp2)
        patt = re.compile("^(.*?)=(.*)$")
        d1 = split(t1)
        d2 = split(t2)
        result = dict()
        for k, v in d2.items():
            if (k in d1) and (d1[k] == d2[k]):
                continue
            if k in ["_", "BASH_ARGC"]:
                continue
            result[k] = v
        return result

    def _check_env_vs_file(self, env, env_path):
        if (not self.exists(env_path)) or (env is None):
            return
        fenv = self.primitive_parse_env_file(env_path)
        for k, v in env.items():
            if k in fenv:
                if env[k] == fenv[k]:
                    log.warning(f"Variable '{k}' specified both as a parameter and in '{env_path}'")
                else:
                    log.warning(f"Variable '{k}' has different values, parameter: {env[k]}, in '{env_path}': "
                        f"{fenv[k]}. According to observation, value from parameter will be used.")

    def deploy_smart_contracts_for_integration_tests(self, network_name, consensus_threshold=None, operator=None,
        owner=None, initial_validator_addresses=None, initial_validator_powers=None, pauser=None,
        mainnet_gas_price=None, env_file=None
    ):
        env = {}
        if consensus_threshold is not None:  # required
            env["CONSENSUS_THRESHOLD"] = str(consensus_threshold)
        if operator is not None:  # required
            env["OPERATOR"] = operator
        if owner is not None:  # required
            env["OWNER"] = owner
        if initial_validator_addresses is not None:
            env["INITIAL_VALIDATOR_ADDRESSES"] = ",".join(initial_validator_addresses)
        if initial_validator_powers is not None:
            env["INITIAL_VALIDATOR_POWERS"] = ",".join([str(x) for x in initial_validator_powers])
        if pauser is not None:
            env["PAUSER"] = pauser
        if mainnet_gas_price is not None:
            env["MAINNET_GAS_PRICE"] = mainnet_gas_price

        env_path = os.path.join(self.project.smart_contracts_dir, ".env")
        if env_file is not None:
            self.copy_file(env_file, env_path)

        self._check_env_vs_file(env, env_path)

        # TODO ui scripts use just "yarn; yarn migrate" alias "npx truffle migrate --reset",
        self.project.npx(["truffle", "deploy", "--network", network_name, "--reset"], env=env,
            cwd=self.project.smart_contracts_dir, pipe=False)

    def deploy_smart_contracts_for_ui_stack(self):
        self.copy_file(os.path.join(self.project.smart_contracts_dir, ".env.ui.example"),
            os.path.join(self.project.smart_contracts_dir, ".env"))
        # TODO Might not be neccessary
        self.project.yarn([], cwd=self.project.smart_contracts_dir)
        self.project.yarn(["migrate"], cwd=self.project.smart_contracts_dir)

    # truffle
    def get_smart_contract_address(self, compiled_json_path, network_id):
        return json.loads(self.read_text_file(compiled_json_path))["networks"][str(network_id)]["address"]

    # truffle
    def get_bridge_smart_contract_addresses(self, network_id):
        return [self.get_smart_contract_address(os.path.join(
            self.project.smart_contracts_dir, f"build/contracts/{x}.json"), network_id)
            for x in ["BridgeToken", "BridgeRegistry", "BridgeBank"]]

    def truffle_exec(self, script_name, *script_args, env=None):
        self._check_env_vs_file(env, os.path.join(self.project.smart_contracts_dir, ".env"))
        script_path = os.path.join(self.project.smart_contracts_dir, f"scripts/{script_name}.js")
        # Hint: call web3 directly, avoid npx + truffle + script
        # Maybe: self.cmd.yarn(["integrationtest:setTokenLockBurnLimit", str(amount)])
        self.project.npx(["truffle", "exec", script_path] + list(script_args), env=env, cwd=self.project.smart_contracts_dir, pipe=False)

    # TODO setTokenLockBurnLimit is gone, possibly replaced by bulkSetTokenLockBurnLimit
    def set_token_lock_burn_limit(self, update_address, amount, ethereum_private_key, infura_project_id, local_provider):
        env = {
            "ETHEREUM_PRIVATE_KEY": ethereum_private_key,
            "UPDATE_ADDRESS": update_address,
            "INFURA_PROJECT_ID": infura_project_id,
            "LOCAL_PROVIDER": local_provider,
        }
        # Needs: ETHEREUM_PRIVATE_KEY, INFURA_PROJECT_ID, LOCAL_PROVIDER, UPDATE_ADDRESS
        # TODO script is no longer there!
        self.truffle_exec("setTokenLockBurnLimit", str(amount), env=env)

    # Peggy1 only
    def sifchain_init_integration(self, sifnode, validator_moniker, validator_mnemonic, denom_whitelist_file):
        # now we have to add the validator key to the test keyring so the tests can send rowan from validator1
        sifnode0 = Sifnoded(self)
        sifnode0.keys_add(validator_moniker, validator_mnemonic)
        valoper = sifnode.keys_show(validator_moniker, bech="val")[0]["address"]
        assert valoper == sifnode0.get_val_address(validator_moniker)  # This does not use "home"; if it the assertion holds it could be grouped with sifchain_init_peggy

        # This was deleted in commit f00242302dd226bc9c3060fb78b3de771e3ff429 from sifchain_start_daemon.sh because
        # it was not working. But we assume that we want to keep it.
        sifnode.sifnoded_exec(["add-genesis-validators", valoper] + sifnode._home_args())

        # Add sifnodeadmin to ~/.sifnoded
        sifnode0 = Sifnoded(self)
        adminuser_addr = sifnode0.keys_add("sifnodeadmin")["address"]
        tokens = {ROWAN: 10 ** 28}
        # Original from peggy:
        # self.cmd.execst(["sifnoded", "add-genesis-account", sifnoded_admin_address, "100000000000000000000rowan", "--home", sifnoded_home])
        sifnode.add_genesis_account(adminuser_addr, tokens)
        sifnode.set_genesis_oracle_admin(adminuser_addr)
        sifnode.set_gen_denom_whitelist(denom_whitelist_file)

        return adminuser_addr

    def sifnoded_peggy2_init_validator(self, sifnode, validator_moniker, validator_mnemonic, evm_network_descriptor, validator_power):
        # Add validator key to test keyring
        # This effectively copies key for validator_moniker from what sifgen creates in /tmp/sifnodedNetwork/validators
        # to ~/.sifnoded (note absence of explicit sifnoded_home, therefore it's ~/.sifnoded)
        sifnode0 = Sifnoded(self)
        # TODO don't add keys to the default ~/ keyring
        sifnode0.keys_add(validator_moniker, validator_mnemonic)

        # Read valoper key
        # (Since we now copied the key to main keyring we could also read it from there)
        valoper = sifnode.get_val_address(validator_moniker)

        # Add genesis validator
        sifnode.add_genesis_validators_peggy(evm_network_descriptor, valoper, validator_power)

    # @TODO Move to Sifgen class
    def sifgen_create_network(self, chain_id: str, validator_count: int, networks_dir: str, network_definition_file: str,
        seed_ip_address: str, mint_amount: Optional[cosmos.Balance] = None
    ):
        # Old call (no longer works either):
        # sifgen network create localnet 1 /mnt/shared/work/projects/sif/sifnode/local-tmp/my/deploy/rake/../networks \
        #     192.168.1.2 /mnt/shared/work/projects/sif/sifnode/local-tmp/my/deploy/rake/../networks/network-definition.yml \
        #     --keyring-backend file
        # self.cmd.execst(["sifgen", "network", "create", "localnet", str(validator_count), networks_dir, seed_ip_address,
        #     os.path.join(networks_dir, "network-definition.yml"), "--keyring-backend", "file"])
        # TODO Most likely, this should be "--keyring-backend file"
        args = ["sifgen", "network", "create", chain_id, str(validator_count), networks_dir, seed_ip_address,
            network_definition_file, "--keyring-backend", "test"] + \
            (["--mint-amount", cosmos.balance_format(mint_amount)] if mint_amount else [])
        self.execst(args)

    # def wait_for_sif_account(self, netdef_json, validator1_address):
    #     # TODO Replace with test_utilities.wait_for_sif_account / wait_for_sif_account_up
    #     return self.execst(["python3", os.path.join(self.project.test_integration_dir, "src/py/wait_for_sif_account.py"),
    #         netdef_json, validator1_address], env={"USER1ADDR": "nothing"})

    def wait_for_sif_account_up(self, address: cosmos.Address, tcp_url: str = None, timeout: int = 90):
        # TODO Move to Sifnoded class
        # TODO Deduplicate: this is also in run_ebrelayer()
        # netdef_json is path to file containing json_dump(netdef)
        # while not self.cmd.tcp_probe_connect("localhost", tendermint_port):
        #     time.sleep(1)
        # self.wait_for_sif_account(netdef_json, validator1_address)

        # Peggy2
        # How this works: by default, the command below will try to do a POST to http://localhost:26657.
        # So the port has to be up first, but this query will fail anyway if it is not.
        args = ["sifnoded", "query", "account", address] + \
            (["--node", tcp_url] if tcp_url else [])
        last_exception = None
        start_time = time.time()
        while True:
            try:
                self.execst(args, disable_log=True)
                break
            except Exception as e:
                if last_exception is None:
                    log.debug(f"Waiting for sif account {address}...")
                if time.time() - start_time > timeout:
                    message = "Timeout expired while waiting for sif account {} to become available".format(address)
                    raise Exception(message) from e
                last_exception = e
                time.sleep(1)

    def _npm_install(self):
        self.project.npm_install(self.project.project_dir("smart-contracts"))


class UIStackEnvironment:
    def __init__(self, cmd):
        self.cmd = cmd
        self.project = cmd.project
        self.chain_id = "sifchain-local"
        self.network_name = "develop"
        self.network_id = 5777
        self.keyring_backend = "test"
        self.ganache_db_path = self.cmd.get_user_home(".ganachedb")
        self.sifnoded_path = self.cmd.get_user_home(".sifnoded")
        self.sifnode = Sifnoded(cmd)

        # From ui/chains/credentials.sh
        self.shadowfiend_name = "shadowfiend"
        self.shadowfiend_mnemonic = ["race", "draft", "rival", "universe", "maid", "cheese", "steel", "logic", "crowd",
            "fork", "comic", "easy", "truth", "drift", "tomorrow", "eye", "buddy", "head", "time", "cash", "swing",
            "swift", "midnight", "borrow"]
        self.kasha_name = "akasha"
        self.akasha_mnemonic = ["hand", "inmate", "canvas", "head", "lunar", "naive", "increase", "recycle", "dog",
            "ecology", "inhale", "december", "wide", "bubble", "hockey", "dice", "worth", "gravity", "ketchup", "feed",
            "balance", "parent", "secret", "orchard"]
        self.juniper_name = "juniper"
        self.juniper_mnemonic = ["clump", "genre", "baby", "drum", "canvas", "uncover", "firm", "liberty", "verb",
            "moment", "access", "draft", "erupt", "fog", "alter", "gadget", "elder", "elephant", "divide", "biology",
            "choice", "sentence", "oppose", "avoid"]
        self.ethereum_root_mnemonic = ["candy", "maple", "cake", "sugar", "pudding", "cream", "honey", "rich", "smooth",
            "crumble", "sweet", "treat"]

    def stack_save_snapshot(self):
        # ui-stack.yml
        # cd .; go get -v -t -d ./...
        # cd ui; yarn install --frozen-lockfile --silent
        # Compile smart contracts:
        # cd ui; yarn build

        # yarn stack --save-snapshot -> ui/scripts/stack.sh -> ui/scripts/stack-save-snapshot.sh
        # rm ui/node_modules/.migrate-complete

        # yarn stack --save-snapshot -> ui/scripts/stack.sh -> ui/scripts/stack-save-snapshot.sh => ui/scripts/stack-launch.sh
        # ui/scripts/stack-launch.sh -> ui/scripts/_sif-build.sh -> ui/chains/sif/build.sh
        # killall sifnoded
        # rm $(which sifnoded)
        self.cmd.rmdir(self.sifnoded_path)
        self.project.make_go_binaries_2()

        # ui/scripts/stack-launch.sh -> ui/scripts/_eth.sh -> ui/chains/etc/launch.sh
        self.cmd.rmdir(self.ganache_db_path)
        self.project.yarn([], cwd=project_dir("ui/chains/eth"))  # Installs ui/chains/eth/node_modules
        # Note that this runs ganache-cli from $PATH whereas scripts start it with yarn in ui/chains/eth
        ganache_proc = Ganache.start_ganache_cli(self.cmd, mnemonic=self.ethereum_root_mnemonic, db=self.ganache_db_path,
            port=7545, network_id=self.network_id, gas_price=20000000000, gas_limit=6721975, host=ANY_ADDR)

        sifnode = Sifnoded(self.cmd)
        # ui/scripts/stack-launch.sh -> ui/scripts/_sif.sh -> ui/chains/sif/launch.sh
        sifnode.sifnoded_init("test", self.chain_id)
        self.cmd.copy_file(project_dir("ui/chains/sif/app.toml"), os.path.join(self.sifnoded_path, "config/app.toml"))
        log.info(f"Generating deterministic account - {self.shadowfiend_name}...")
        shadowfiend_account = self.cmd.sifnoded_keys_add(self.shadowfiend_name, self.shadowfiend_mnemonic)
        log.info(f"Generating deterministic account - {self.akasha_name}...")
        akasha_account = self.sifnode.keys_add(self.akasha_name, self.akasha_mnemonic)
        log.info(f"Generating deterministic account - {self.juniper_name}...")
        juniper_account = self.cmd.sifnoded_keys_add(self.juniper_name, self.juniper_mnemonic)
        shadowfiend_address = shadowfiend_account["address"]
        akasha_address = akasha_account["address"]
        juniper_address = juniper_account["address"]
        assert shadowfiend_address == self.sifnode.keys_show(self.shadowfiend_name)[0]["address"]
        assert akasha_address == self.sifnode.keys_show(self.akasha_name)[0]["address"]
        assert juniper_address == self.sifnode.keys_show(self.juniper_name)[0]["address"]

        tokens_shadowfiend = [[10**29, "rowan"], [10**29, "catk"], [10**29, "cbtk"], [10**29, "ceth"], [10**29, "cusdc"], [10**29, "clink"], [10**26, "stake"]]
        tokens_akasha = [[10**29, "rowan"], [10**29, "catk"], [10**29, "cbtk"], [10**29, "ceth"], [10**29, "cusdc"], [10**29, "clink"], [10**26, "stake"]]
        tokens_juniper = [[10**22, "rowan"], [10**22, "cusdc"], [10**20, "clink"], [10**20, "ceth"]]
        sifnode.add_genesis_account(shadowfiend_address, tokens_shadowfiend)
        sifnode.add_genesis_account(akasha_address, tokens_akasha)
        sifnode.add_genesis_account(juniper_address, tokens_juniper)

        shadowfiend_address_bech_val = sifnode.keys_show(self.shadowfiend_name, bech="val")[0]["address"]
        self.cmd.sifnoded_add_genesis_validators(shadowfiend_address_bech_val)

        amount = sif_format_amount(10**24, "stake")
        self.cmd.execst(["sifnoded", "gentx", self.shadowfiend_name, amount, f"--chain-id={self.chain_id}",
            f"--keyring-backend={self.keyring_backend}"])

        log.info("Collecting genesis txs...")
        self.cmd.execst(["sifnoded", "collect-gentxs"])
        log.info("Validating genesis file...")
        self.cmd.execst(["sifnoded", "validate-genesis"])

        log.info("Starting test chain...")
        sifnoded_proc = self.cmd.sifnoded_start(minimum_gas_prices=(0.5, ROWAN))  # TODO sifnoded_home=???

        # sifnoded must be up before continuing
        self.cmd.sif_wait_up("localhost", 1317)

        # ui/scripts/_migrate.sh -> ui/chains/peggy/migrate.sh
        self.cmd.deploy_smart_contracts_for_ui_stack()

        # ui/scripts/_migrate.sh -> ui/chains/eth/migrate.sh
        # send through atk and btk tokens to eth chain
        self.project.yarn(["migrate"], cwd=project_dir("ui/chains/eth"))

        # ui/scripts/_migrate.sh -> ui/chains/sif/migrate.sh
        # Original scripts say "if we don't sleep there are issues"
        time.sleep(10)
        log.info("Creating liquidity pool from catk:rowan...")
        sifnode.tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "catk", [10**5, ROWAN], 10**25, 10**25)
        time.sleep(5)
        log.info("Creating liquidity pool from cbtk:rowan...")
        sifnode.tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "cbtk", [10**5, ROWAN], 10**25, 10**25)
        # should now be able to swap from catk:cbtk
        time.sleep(5)
        log.info("Creating liquidity pool from ceth:rowan...")
        sifnode.tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "ceth", [10**5, ROWAN], 10**25, 83*10**20)
        # should now be able to swap from x:ceth
        time.sleep(5)
        log.info("Creating liquidity pool from cusdc:rowan...")
        sifnode.tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "cusdc", [10**5, ROWAN], 10**25, 10**25)
        time.sleep(5)
        log.info("Creating liquidity pool from clink:rowan...")
        sifnode.tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "clink", [10**5, ROWAN], 10**25, 588235*10**18)
        time.sleep(5)
        log.info("Creating liquidity pool from ctest:rowan...")
        sifnode.tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "ctest", [10**5, ROWAN], 10**25, 10**13)

        # ui/scripts/_migrate.sh -> ui/chains/post_migrate.sh

        atk_address, btk_address, usdc_address, link_address = [
            self.cmd.get_smart_contract_address(project_dir(f"ui/chains/eth/build/contracts/{x}.json"), self.network_id)
            for x in ["AliceToken", "BobToken", "UsdCoin", "LinkCoin"]
        ]

        bridge_token_address, bridge_registry_address, bridge_bank = self.cmd.get_bridge_smart_contract_addresses(self.network_id)

        # From smart-contracts/.env.ui.example
        smart_contracts_env_ui_example_vars = {
            "ETHEREUM_PRIVATE_KEY": "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3",
            "INFURA_PROJECT_ID": "JFSH7439sjsdtqTM23Dz",
            "LOCAL_PROVIDER": "http://localhost:7545",
        }

        # NOTE: this probably doesn't work anymore since setTokenLockBurnLimit.js was replaced
        burn_limits = [
            [eth.NULL_ADDRESS, 31 * 10 ** 18],
            [bridge_token_address, 10 ** 25],
            [atk_address, 10 ** 25],
            [btk_address, 10 ** 25],
            [usdc_address, 10 ** 25],
            [link_address, 10 ** 25],
        ]
        for address, amount in burn_limits:
            self.cmd.set_token_lock_burn_limit(
                address,
                amount,
                smart_contracts_env_ui_example_vars["ETHEREUM_PRIVATE_KEY"],
                smart_contracts_env_ui_example_vars["INFURA_PROJECT_ID"],
                smart_contracts_env_ui_example_vars["LOCAL_PROVIDER"]
            )

        # signal migrate-complete

        # Whitelist test tokens
        for addr in [atk_address, btk_address, usdc_address, link_address]:
            self.project.yarn(["peggy:whiteList", addr, "true"], cwd=self.project.smart_contracts_dir)

        # ui/scripts/stack-launch.sh -> ui/scripts/_peggy.sh -> ui/chains/peggy/launch.sh
        # rm -rf ui/chains/peggy/relayerdb
        # ebrelayer is in $GOBIN, gets installed by "make install"
        ethereum_private_key = smart_contracts_env_ui_example_vars["ETHEREUM_PRIVATE_KEY"]
        ebrelayer_proc = Ebrelayer(self.cmd).init("tcp://localhost:26657", "ws://localhost:7545/",
            bridge_registry_address, self.shadowfiend_name, self.shadowfiend_mnemonic, self.chain_id,
            ethereum_private_key=ethereum_private_key, gas=5*10**12, gas_prices=[0.5, ROWAN])

        # At this point we have 3 running processes - ganache_proc, sifnoded_proc and ebrelayer_proc
        # await sif-node-up and migrate-complete

        time.sleep(30)
        # ui/scripts/_snapshot.sh

        # ui/scripts/stack-pause.sh:
        # killall sifnoded sifnoded ebrelayer ganache-cli
        sifnoded_proc.kill()
        ebrelayer_proc.kill()
        ganache_proc.kill()
        time.sleep(10)

        snapshots_dir = project_dir("ui/chains/snapshots")
        self.cmd.mkdir(snapshots_dir)  # TODO self.cmd.rmdir(snapshots_dir)
        # ui/chains/peggy/snapshot.sh:
        # mkdir -p ui/chains/peggy/relayerdb
        self.cmd.tar_create(project_dir("ui/chains/peggy/relayerdb"), os.path.join(snapshots_dir, "peggy.tar.gz"))
        # mkdir -p smart-contracts/build
        self.cmd.tar_create(project_dir("smart-contracts/build"), os.path.join(snapshots_dir, "peggy_build.tar.gz"))

        # ui/chains/sif/snapshot.sh:
        self.cmd.tar_create(self.sifnoded_path, os.path.join(snapshots_dir, "sif.tar.gz"))

        # ui/chains/etc/snapshot.sh:
        self.cmd.tar_create(self.ganache_db_path, os.path.join(snapshots_dir, "eth.tar.gz"))

    def stack_push(self):
        # ui/scripts/stack-push.sh
        # $PWD=ui

        # User must be logged in to Docker hub:
        # ~/.docker/config.json must exist and .auths['ghcr.io'].auth != null
        log.info("Github Registry Login found.")

        commit = exactly_one(stdout_lines(self.cmd.execst(["git", "rev-parse", "HEAD"], cwd=project_dir())))
        branch = exactly_one(stdout_lines(self.cmd.execst(["git", "rev-parse", "--abbrev-ref", "HEAD"], cwd=project_dir())))

        image_root = "ghcr.io/sifchain/sifnode/ui-stack"
        image_name = "{}:{}".format(image_root, commit)
        stable_tag = "{}:{}".format(image_root, branch.replace("/", "__"))

        running_in_ci = bool(os.environ.get("CI"))

        if running_in_ci:
            res = self.cmd.execst(["git", "status", "--porcelain", "--untracked-files=no"], cwd=project_dir())
            # # reverse grep for go.mod because on CI this can be altered by installing go dependencies
            # if [[ -z "$CI" && ! -z "$(git status --porcelain --untracked-files=no)" ]]; then
            #   echo "Git workspace must be clean to save git commit hash"
            #   exit 1
            # fi
            pass

        log.info("Building new container...")
        log.info(f"New image name: {image_name}")

        self.cmd.execst(["docker", "build", "-f", project_dir("ui/scripts/stack.Dockerfile"), "-t", image_name, "."],
            cwd=project_dir(), env={"DOCKER_BUILDKIT": "1"}, pipe=False)

        if running_in_ci:
            log.info(f"Tagging image as {stable_tag}...")
            self.cmd.execst(["docker", "tag", image_name, stable_tag])
            self.cmd.execst(["docker", "push", image_name])
            self.cmd.execst(["docker", "push", stable_tag])


class IntegrationTestsEnvironment:
    def __init__(self, cmd):
        self.cmd = cmd
        self.project = cmd.project
        # Fixed, set in start-integration-env.sh
        self.ethereum_private_key = "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3"
        self.owner = "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
        # we may eventually switch things so PAUSER and OWNER aren't the same account, but for now they're the same
        self.pauser = self.owner
        # set_persistant_env_var BASEDIR $(fullpath $BASEDIR) $envexportfile
        # set_persistant_env_var SIFCHAIN_BIN $BASEDIR/cmd $envexportfile
        # set_persistant_env_var envexportfile $(fullpath $envexportfile) $envexportfile
        # set_persistant_env_var TEST_INTEGRATION_DIR ${BASEDIR}/test/integration $envexportfile
        # set_persistant_env_var TEST_INTEGRATION_PY_DIR ${BASEDIR}/test/integration/src/py $envexportfile
        # set_persistant_env_var SMART_CONTRACTS_DIR ${BASEDIR}/smart-contracts $envexportfile
        # set_persistant_env_var datadir ${TEST_INTEGRATION_DIR}/vagrant/data $envexportfile
        # set_persistant_env_var CONTAINER_NAME integration_sifnode1_1 $envexportfile
        # set_persistant_env_var NETWORKDIR $BASEDIR/deploy/networks $envexportfile
        # set_persistant_env_var GANACHE_DB_DIR $(mktemp -d /tmp/ganachedb.XXXX) $envexportfile
        # set_persistant_env_var ETHEREUM_WEBSOCKET_ADDRESS ws://localhost:7545/ $envexportfile
        # set_persistant_env_var CHAINNET localnet $envexportfile
        self.network_name = "develop"
        self.network_id = 5777
        self.peruser_storage_dir = self.cmd.get_user_home(".sifnode-integration")
        self.state_vars = {}
        self.test_integration_dir = project_dir("test/integration")
        self.data_dir = project_dir("test/integration/vagrant/data")
        self.chainnet = "localnet"
        self.tcp_url = f"tcp://{ANY_ADDR}:26657"
        self.ethereum_websocket_address = "ws://localhost:7545/"
        self.ganache_mnemonic = ["candy", "maple", "cake", "sugar", "pudding", "cream", "honey", "rich", "smooth",
                "crumble", "sweet", "treat"]
        self.mint_amount = {
            ROWAN: 999999 * 10**21,
            "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE": 137 * 10**16
        }

        self.ganache_executable = self.project.project_dir("smart-contracts", "node_modules", ".bin", "ganache-cli")

    def prepare(self):
        self.project.make_go_binaries()
        self.project.install_smart_contracts_dependencies()

    def run(self):
        self.cmd.mkdir(self.data_dir)

        self.prepare()

        log_dir = self.cmd.tmpdir("sifnode")
        self.cmd.mkdir(log_dir)
        ganache_log_file = open(os.path.join(log_dir, "ganache.log"), "w")  # TODO close
        sifnoded_log_file = open(os.path.join(log_dir, "sifnoded.log"), "w")  # TODO close
        ebrelayer_log_file = open(os.path.join(log_dir, "ebrelayer.log"), "w")  # TODO close

        # test/integration/ganache-start.sh:
        # 1. pkill -9 -f ganache-cli || true
        # 2. while nc -z localhost 7545; do sleep 1; done
        # 3. nohup tmux new-session -d -s my_session "ganache-cli ${block_delay} -h 0.0.0.0 --mnemonic \
        #     'candy maple cake sugar pudding cream honey rich smooth crumble sweet treat' \
        #     --networkId '5777' --port '7545' --db ${GANACHE_DB_DIR} --account_keys_path $GANACHE_KEYS_JSON \
        #     > $GANACHE_LOG 2>&1"
        # 4. sleep 5
        # 5. while ! nc -z localhost 4545; do sleep 5; done
        # GANACHE_LOG=ui/test/integration/vagrant/data/logs/ganache.$(filenamedate).txt
        block_time = None  # TODO
        account_keys_path = os.path.join(self.data_dir, "ganachekeys.json")
        ganache_db_path = self.cmd.mktempdir()
        ganache_proc = Ganache.start_ganache_cli(self.cmd, executable=self.ganache_executable, block_time=block_time,
            host=ANY_ADDR, mnemonic=self.ganache_mnemonic, network_id=self.network_id, port=7545, db=ganache_db_path,
            account_keys_path=account_keys_path, log_file=ganache_log_file)

        self.cmd.wait_for_file(account_keys_path)  # Created by ganache-cli
        time.sleep(2)

        # Pick an account for ebrelayer from 10 hardcoded ganache_keys. In theory, it shouldn't matter which one we pick
        # but since other parts of the code might have some hidden assumptions we just pick a fixed one for now.
        ganache_keys = json.loads(self.cmd.read_text_file(account_keys_path))
        ebrelayer_ethereum_addr = "0x5aeda56215b167893e80b4fe645ba6d5bab767de"
        assert ebrelayer_ethereum_addr in ganache_keys["private_keys"]
        ebrelayer_ethereum_private_key = ganache_keys["private_keys"][ebrelayer_ethereum_addr]

        env_file = project_dir("test/integration/.env.ciExample")
        env_vars = self.cmd.primitive_parse_env_file(env_file)
        self.cmd.deploy_smart_contracts_for_integration_tests(self.network_name, owner=self.owner, pauser=self.pauser,
            operator=env_vars["OPERATOR"], consensus_threshold=int(env_vars["CONSENSUS_THRESHOLD"]),
            initial_validator_powers=[int(x) for x in env_vars["INITIAL_VALIDATOR_POWERS"].split(",")],
            initial_validator_addresses=[ebrelayer_ethereum_addr], env_file=env_file)

        bridge_token_sc_addr, bridge_registry_sc_addr, bridge_bank_sc_addr = \
            self.cmd.get_bridge_smart_contract_addresses(self.network_id)

        # # TODO This should be last (after return from setup_sifchain.sh)
        # burn_limits = [
        #     [eth.NULL_ADDRESS, 31*10**18],
        #     [bridge_token_sc_addr, 10**25],
        # ]
        # env_file_vars = self.cmd.primitive_parse_env_file(env_file)
        # for address, amount in burn_limits:
        #     self.cmd.set_token_lock_burn_limit(
        #         address,
        #         amount,
        #         env_file_vars["ETHEREUM_PRIVATE_KEY"],  # != ebrelayer_ethereum_private_key
        #         env_file_vars["INFURA_PROJECT_ID"],
        #         env_file_vars["LOCAL_PROVIDER"],  # for web3.js to connect to ganache
        #     )

        # test/integration/setup_sifchain.sh:
        networks_dir = project_dir("deploy/networks")
        self.cmd.rmdir(networks_dir)  # networks_dir has many directories without write permission, so change those before deleting it
        self.cmd.mkdir(networks_dir)
        # Old:
        # self.cmd.execst(["rake", f"genesis:network:scaffold[{self.chainnet}]"], env={"BASEDIR": project_dir()}, pipe=False)
        # New:
        # sifgen network create localnet 1 $NETWORKDIR 192.168.1.2 $NETWORKDIR/network-definition.yml --keyring-backend test \
        #     --mint-amount 999999000000000000000000000rowan,1370000000000000000ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE
        chain_id = "localnet"
        validator_count = 1
        network_definition_file = os.path.join(networks_dir, "network-definition.yml")
        seed_ip_address = "192.168.1.2"

        self.cmd.sifgen_create_network(chain_id, validator_count, networks_dir, network_definition_file, seed_ip_address, mint_amount=self.mint_amount)

        netdef, netdef_json = self.process_netdef(network_definition_file)

        validator1_moniker = netdef["moniker"]
        validator1_address = netdef["address"]
        validator1_password = netdef["password"]
        validator1_mnemonic = netdef["mnemonic"].split(" ")
        chaindir = os.path.join(networks_dir, f"validators/{self.chainnet}/{validator1_moniker}")
        sifnoded_home = os.path.join(chaindir, ".sifnoded")
        denom_whitelist_file = os.path.join(self.test_integration_dir, "whitelisted-denoms.json")
        # SIFNODED_LOG=$datadir/logs/sifnoded.log

        sifnode = Sifnoded(self.cmd, home=sifnoded_home)

        adminuser_addr = self.cmd.sifchain_init_integration(sifnode, validator1_moniker, validator1_mnemonic,
            denom_whitelist_file)

        # Start sifnoded
        sifnoded_proc = sifnode.sifnoded_start(tcp_url=self.tcp_url, minimum_gas_prices=(0.5, ROWAN),
            log_file=sifnoded_log_file, trace=True)

        # TODO: wait for sifnoded to come up before continuing
        # in sifchain_start_daemon.sh: "sleep 10"
        # in sifchain_run_ebrelayer.sh (also run_ebrelayer here) we already wait for connection to port 26657 and sif account validator1_addr

        # Removed
        # # TODO Process exits immediately with returncode 1
        # # TODO Why does it not stop start-integration-env.sh?
        # # rest_server_proc = self.cmd.popen(["sifnoded", "rest-server", "--laddr", "tcp://0.0.0.0:1317"])  # TODO cwd

        # test/integration/sifchain_start_ebrelayer.sh -> test/integration/sifchain_run_ebrelayer.sh
        # This script is also called from tests

        relayer_db_path = os.path.join(self.test_integration_dir, "sifchainrelayerdb")

        # Prevent starting over dirty/existing relayer_db_path
        if self.cmd.exists(relayer_db_path):
            assert not self.cmd.ls(relayer_db_path), "relayer_db_path {} not empty".format(relayer_db_path)

        ebrelayer_proc = self.run_ebrelayer(netdef_json, validator1_address, validator1_moniker, validator1_mnemonic,
            ebrelayer_ethereum_private_key, bridge_registry_sc_addr, relayer_db_path, log_file=ebrelayer_log_file)

        vagrantenv_path = project_dir("test/integration/vagrantenv.sh")
        self.state_vars = {
            "ETHEREUM_PRIVATE_KEY": self.ethereum_private_key,
            "OWNER": self.owner,
            "PAUSER": self.pauser,
            "BASEDIR": project_dir(),
            "envexportfile": vagrantenv_path,
            "SMART_CONTRACTS_DIR": self.project.smart_contracts_dir,
            "NETWORKDIR": networks_dir,
            "GANACHE_DB_DIR": ganache_db_path,
            "EBRELAYER_ETHEREUM_ADDR": ebrelayer_ethereum_addr,
            "EBRELAYER_ETHEREUM_PRIVATE_KEY": ebrelayer_ethereum_private_key,  # Needed by sifchain_run_ebrelayer.sh
            "BRIDGE_REGISTRY_ADDRESS": bridge_registry_sc_addr,
            "BRIDGE_TOKEN_ADDRESS": bridge_token_sc_addr,
            "BRIDGE_BANK_ADDRESS": bridge_bank_sc_addr,
            "NETDEF": os.path.join(networks_dir, "network-definition.yml"),
            "NETDEF_JSON": project_dir("test/integration/vagrant/data/netdef.json"),
            "MONIKER": validator1_moniker,
            "VALIDATOR1_PASSWORD": validator1_password,
            "VALIDATOR1_ADDR": validator1_address,
            "MNEMONIC": " ".join(validator1_mnemonic),
            "CHAINDIR": os.path.join(networks_dir, "validators", self.chainnet, validator1_moniker),
            "SIFCHAIN_ADMIN_ACCOUNT": adminuser_addr,  # Needed by test_peggy_fees.py (via conftest.py)
            "EBRELAYER_DB": relayer_db_path,  # Created by sifchain_run_ebrelayer.sh, does not appear to be used anywhere at the moment
        }
        self.project.write_vagrantenv_sh(self.state_vars, self.data_dir, self.ethereum_websocket_address, self.chainnet)

        log.debug("Admin account address (rowan source): {}".format(adminuser_addr))

        from siftool import localnet
        localnet.run_localnet_hook()

        return ganache_proc, sifnoded_proc, ebrelayer_proc

    def remove_and_add_sifnoded_keys(self, moniker, mnemonic):
        # Error: The specified item could not be found in the keyring
        # This is not neccessary during start-integration-env.sh (as the key does not exist yet), but is neccessary
        # during tests that restart ebrelayer
        # res = self.cmd.execst(["sifnoded", "keys", "delete", moniker, "--keyring-backend", "test"], stdin=["y"])
        sifnode = Sifnoded(self.cmd)
        sifnode.keys_delete(moniker)
        sifnode.keys_add(moniker, mnemonic)

    def process_netdef(self, network_definition_file):
        # networks_dir = deploy/networks
        # File deploy/networks/network-definition.yml is created by "rake genesis:network:scaffold", specifically by
        # "sifgen network create"
        # We read it and convert to test/integration/vagrant/data/netdef.json
        netdef = exactly_one(yaml_load(self.cmd.read_text_file(network_definition_file)))
        netdef_json = os.path.join(self.data_dir, "netdef.json")
        self.cmd.write_text_file(netdef_json, json.dumps(netdef, indent=4))
        return netdef, netdef_json

    def run_ebrelayer(self, netdef_json, validator1_address, validator_moniker, validator_mnemonic,
        ebrelayer_ethereum_private_key, bridge_registry_sc_addr, relayer_db_path, log_file=None
    ):
        # TODO Deduplicate
        while not self.cmd.tcp_probe_connect("localhost", 26657):
            time.sleep(1)
        # self.cmd.wait_for_sif_account(netdef_json, validator1_address)
        self.cmd.wait_for_sif_account_up(validator1_address)
        time.sleep(10)
        self.remove_and_add_sifnoded_keys(validator_moniker, validator_mnemonic)  # Creates ~/.sifnoded/keyring-tests/xxxx.address
        ebrelayer_proc = Ebrelayer(self.cmd).init(self.tcp_url, self.ethereum_websocket_address, bridge_registry_sc_addr,
            validator_moniker, validator_mnemonic, self.chainnet, ethereum_private_key=ebrelayer_ethereum_private_key,
            node=self.tcp_url, keyring_backend="test", sign_with=validator_moniker,
            symbol_translator_file=os.path.join(self.test_integration_dir, "config/symbol_translator.json"),
            relayerdb_path=relayer_db_path, cwd=self.test_integration_dir, log_file=log_file)
        return ebrelayer_proc

    def create_own_dirs(self):
        self.cmd.mkdir(self.peruser_storage_dir)
        self.cmd.mkdir(os.path.join(self.peruser_storage_dir, "snapshots"))

    def create_snapshot(self, snapshot_name):
        self.create_own_dirs()
        named_snapshot_dir = os.path.join(self.peruser_storage_dir, "snapshots", snapshot_name)
        if self.cmd.exists(named_snapshot_dir):
            raise Exception(f"Directory '{named_snapshot_dir}' already exists")
        self.cmd.mkdir(named_snapshot_dir)
        self.cmd.tar_create(self.state_vars["GANACHE_DB_DIR"], os.path.join(named_snapshot_dir, "ganache.tar.gz"))
        self.cmd.tar_create(self.state_vars["EBRELAYER_DB"], os.path.join(named_snapshot_dir, "sifchainrelayerdb.tar.gz"))
        self.cmd.tar_create(project_dir("deploy/networks"), os.path.join(named_snapshot_dir, "networks.tar.gz"))
        self.cmd.tar_create(project_dir("smart-contracts/build"), os.path.join(named_snapshot_dir, "smart-contracts.tar.gz"))
        self.cmd.tar_create(self.cmd.get_user_home(".sifnoded"), os.path.join(named_snapshot_dir, "sifnoded.tar.gz"))
        self.cmd.write_text_file(os.path.join(named_snapshot_dir, "vagrantenv.json"), json.dumps(self.state_vars, indent=4))

    def restore_snapshot(self, snapshot_name):
        named_snapshot_dir = os.path.join(self.peruser_storage_dir, "snapshots", snapshot_name)
        state_vars = json.loads(self.cmd.read_text_file(os.path.join(named_snapshot_dir, "vagrantenv.json")))

        def extract(tarfile, path):
            self.cmd.rmdir(path)
            self.cmd.mkdir(path)
            self.cmd.tar_extract(os.path.join(named_snapshot_dir, tarfile), path)

        ganache_db_dir = self.cmd.mktempdir()
        relayer_db_path = state_vars["EBRELAYER_DB"]  # TODO use /tmp
        assert os.path.realpath(relayer_db_path) == os.path.realpath(os.path.join(self.test_integration_dir, "sifchainrelayerdb"))
        extract("ganache.tar.gz", ganache_db_dir)
        extract("sifchainrelayerdb.tar.gz", relayer_db_path)
        deploy_networks_dir = project_dir("deploy/networks")
        extract("networks.tar.gz", deploy_networks_dir)
        smart_contracts_build_dir = project_dir("smart-contracts/build")
        extract("smart-contracts.tar.gz", smart_contracts_build_dir)
        extract("sifnoded.tar.gz", self.cmd.get_user_home(".sifnoded"))  # Needed for "--keyring-backend test"

        state_vars["GANACHE_DB_DIR"] = ganache_db_dir
        state_vars["EBRELAYER_DB"] = relayer_db_path
        self.state_vars = state_vars
        self.project.write_vagrantenv_sh(state_vars, self.data_dir, self.ethereum_websocket_address, self.chainnet)
        self.cmd.mkdir(self.data_dir)

    def restart_processes(self):
        block_time = None
        ganache_db_path = self.state_vars["GANACHE_DB_DIR"]
        account_keys_path = os.path.join(self.data_dir, "ganachekeys.json")  # TODO this is in test/integration/vagrant/data, which is supposed to be cleared

        ganache_proc = Ganache.start_ganache_cli(self.cmd, executable=self.ganache_executable, block_time=block_time,
            host=ANY_ADDR, mnemonic=self.ganache_mnemonic, network_id=self.network_id, port=7545, db=ganache_db_path,
            account_keys_path=account_keys_path)  # TODO log_file

        self.cmd.wait_for_file(account_keys_path)  # Created by ganache-cli
        time.sleep(2)

        validator_moniker = self.state_vars["MONIKER"]
        networks_dir = project_dir("deploy/networks")
        chaindir = os.path.join(networks_dir, f"validators/{self.chainnet}/{validator_moniker}")
        sifnoded_home = os.path.join(chaindir, ".sifnoded")
        sifnoded_proc = self.cmd.sifnoded_start(tcp_url=self.tcp_url, minimum_gas_prices=(0.5, ROWAN),
            sifnoded_home=sifnoded_home, trace=True)

        bridge_token_sc_addr, bridge_registry_sc_addr, bridge_bank_sc_addr = \
            self.cmd.get_bridge_smart_contract_addresses(self.network_id)

        validator_mnemonic = self.state_vars["MNEMONIC"].split(" ")
        account_keys_path = os.path.join(self.data_dir, "ganachekeys.json")
        ganache_keys = json.loads(self.cmd.read_text_file(account_keys_path))
        ebrelayer_ethereum_addr = list(ganache_keys["private_keys"].keys())[9]
        ebrelayer_ethereum_private_key = ganache_keys["private_keys"][ebrelayer_ethereum_addr]
        network_definition_file = project_dir(networks_dir, "network-definition.yml")

        netdef, netdef_json = self.process_netdef(network_definition_file)
        validator1_address = netdef["address"]
        assert validator1_address == self.state_vars["VALIDATOR1_ADDR"]
        relayer_db_path = self.state_vars["EBRELAYER_DB"]
        ebrelayer_proc = self.run_ebrelayer(netdef_json, validator1_address, validator_moniker, validator_mnemonic,
            ebrelayer_ethereum_private_key, bridge_registry_sc_addr, relayer_db_path)

        return ganache_proc, sifnoded_proc, ebrelayer_proc


# Peggy2 environment
# Originally derived from devenv which ran hardhat, sifnoded, ebrelayer like this: cd smart-contracts; GOBIN=~/go/bin npx hardhat run scripts/devenv.ts
class Peggy2Environment(IntegrationTestsEnvironment):
    def __init__(self, cmd: Command, witness_count: int = 1, witness_power: int = 100, consensus_threshold: int = 49,
        ethereum_ws_port: int = 8545, ethereum_chain_id: int = 9999, sif_chain_id = "localnet", use_geth: bool = False
    ):
        super().__init__(cmd)
        self.extra_balances_for_admin_account = None
        self.witness_count = witness_count
        self.witness_power = witness_power
        self.consensus_threshold = consensus_threshold

        self.ethereum_ws_port = ethereum_ws_port
        self.ethereum_chain_id = ethereum_chain_id
        self.chain_id = sif_chain_id
        self.ceth_symbol = sifchain_denom_hash(self.ethereum_chain_id, eth.NULL_ADDRESS)
        assert self.ceth_symbol == "sifBridge99990x0000000000000000000000000000000000000000"

        # This goes to validator0, i.e. sifnode_validators[0]["address"]
        self.validator_mint_amounts: cosmos.Balance = {
            ROWAN: 999999 * 10**27,
            self.ceth_symbol: 999999 * 10**21,
            "sifBridge00030x1111111111111111111111111111111111111111": 137 * 10**16,
            "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE": 137 * 10**16,
        }
        # These go to admin account, relayers and witnesses
        self.admin_account_mint_amounts: cosmos.Balance = {
            ROWAN: 10**27,
            self.ceth_symbol: 2 * 10**22,
            "sifBridge00030x1111111111111111111111111111111111111111": 10 ** 16,
            "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE": 10 ** 16,
        }

        # Log levels for "ebrelayer --log_level". Use `None` to skip this parameter and let ebrelayer use its default
        self.log_level_relayer = None
        self.log_level_witness = None

        # The following are only meaningful if use_geth_instead_of_hardhat is enabled
        self.use_geth_instead_of_hardhat = use_geth
        # From https://stackoverflow.com/a/50972677/8338100:
        # > The gasLimit in the genesis block is only used as a starting point. As blocks are added to the chain, the
        # > block gas limit will change over time based on the miners processing the transactions on the network. To
        # > keep the block gas limit higher, you need to override the default configuration in your node client.
        # In dev_mode this is passed via "--dev.gaslimit", and in regular mode via genesis + "--miner.gaslimit".
        # "None" means to use some arbitrary default (typically 8000000).
        # Ethereum mainnet currently uses gas price of ~30 * 10**6.
        self.geth_gas_limit = 30 * 10**6
        self.geth_gas_limit = 100 * 10**6
        self.geth_dev_mode = False
        self.geth_block_mining_period = 1  # Only applicable when not geth_dev_mode

    # Destuctures a linear array of EVM accounts into:
    # [operator, owner, pauser, [validator-0, validator-1, ...], [...available...]]
    # proxy_admin is the same as operator.
    # Keep this synched with deploy_contracts_dev.ts
    def signer_array_to_ethereum_accounts(self, accounts, n_validators):
        assert len(accounts) >= n_validators + 3
        operator, owner, pauser, *rest = accounts  # Take 3 and store remaining in rest
        validators, available = rest[:n_validators], rest[n_validators:]  # Take n_validators for validators the remaining for available
        return {
            "proxy_admin": operator,
            "operator": operator,
            "owner": owner,
            "pauser": pauser,
            "validators": validators,
            "available": available,
        }

    # Override
    def run(self) -> Tuple[subprocess.Popen, subprocess.Popen, subprocess.Popen, List[subprocess.Popen]]:
        # self.project._make_go_binaries()

        # Ordering (for possible parallelisation):
        # - build_golang_binaries before start_sifchain
        # - start_hardhat before deploy_smart_contracts
        # - start_sifchain before start_witnesses_and_relayers
        # - deploy_smart_contracts before start_witnesses_and_relayers
        # - start_witnesses_and_relayers before return
        # - write_env_file before return

        # TODO: where is log watcher?

        log_dir = self.cmd.tmpdir("sifnode")
        self.cmd.mkdir(log_dir)
        hardhat_log_file = open(os.path.join(log_dir, "hardhat.log"), "w")  # TODO close + use a different name
        sifnoded_log_file = open(os.path.join(log_dir, "sifnoded.log"), "w")  # TODO close
        relayer_log_file = open(os.path.join(log_dir, "relayer.log"), "w")  # TODO close

        self.cmd.rmdir(self.cmd.get_user_home(".sifnoded"))  # Purge test keyring backend

        hardhat = Hardhat(self.cmd)

        hardhat.compile_smart_contracts()

        # This determines how many EVM accounts we want to allocate for validators.
        # Since every validator needs on EVM account, this should be equal to the number of validators (possibly more).
        hardhat_validator_count = self.witness_count
        hardhat_port = 8545  # The port on which to listen for new connections (default: 8545)

        # We need at least 4 accounts for operator, owner, pauser, validator1. They can be any accounts, but if we're
        # running hardhat node (not use_geth_instead_of_hardhat), hardhat supplies its own hardcoded accounts and in
        # that case we need to use the same ones here. If we're running more validators we need more than 4. The first
        # account is also used for geth (if use_geth_instead_of_hardhat). They are used like this:
        # [operatorAccount, ownerAccount, pauserAccount, validator1Account, ...extraAccounts].
        sample_eth_accounts = hardhat_default_accounts()
        assert len(sample_eth_accounts) >= 4

        hardhat_accounts = self.signer_array_to_ethereum_accounts(sample_eth_accounts, hardhat_validator_count)

        # Initialization of smart contracts (technically this is part of deployment)
        operator_acct = hardhat_accounts["operator"]
        operator_addr, operator_private_key = operator_acct

        w3_url = eth.web3_host_port_url("localhost", self.ethereum_ws_port)

        files_to_delete = []
        manual_funds_alloc = None
        try:
            if self.use_geth_instead_of_hardhat:
                # EnvCtx takes environment "ETH_ACCOUNT_OWNER_ADDRESS" for funding ETH accounts which is set from
                # hardhat_accounts["owner"] Smart contract delpoyment however needs funds on 0'th account which is the
                # operator. For now we fund both. TODO cleanup, allocate all ether to make EnvCtx fund accounts from ETH_ACCOUNT_OPERATOR_ADDRESS
                # validators need funds to run, too.
                accounts_to_fund = \
                    [hardhat_accounts[who][0] for who in ["owner", "operator"]] + \
                    [acct[0] for acct in hardhat_accounts["validators"]]
                funds_alloc = {addr: 1000000 * eth.ETH for addr in accounts_to_fund}
                # Default http port for geth is 8545, but we're using it for ws:// for compatibility with hardhat.
                # Since we're using hardhat for deployment, and deployment script scripts/deploy_contracts_dev.ts uses
                # hardhat, and since hardhat can deploy only over http, we also need geth to listen on http port.
                geth_http_port = self.ethereum_ws_port + 1
                geth_ws_port = 8545  # We're reversing default values for http and ws ports to keep ws on the same port as hardhat
                geth_datadir = self.cmd.tmpdir("geth")  # TODO self.cmd.mktempdir()
                self.cmd.rmdir(geth_datadir)
                self.cmd.mkdir(geth_datadir)
                geth = Geth(self.cmd, datadir=geth_datadir)
                if self.geth_dev_mode:
                    # "geth --dev" runs a proof-of-authority chain with on-demand mining and zero gas price.
                    # See https://geth.ethereum.org/docs/getting-started/dev-mode for more information about it.
                    # Unfortunately, there is no standard for connecting to such PoA chains, so we need to inject a
                    # custom middleware into every web3py connection. (It seems that this is only needed for using
                    # eth.sendTransaction(), but not for eth.sendRawTransaction(), so we only use it below, while
                    # connections from test_utils and TypeScript still work without this.)
                    # See https://web3py.readthedocs.io/en/stable/middleware.html#geth-style-proof-of-authority
                    # Also unfortunately, dev mode seems to be incompatible with ebrelayer which waits for 50 blocks
                    # to process transaction, while "geth --dev" mode will only mine blocks on demand => deadlock.
                    # Could be solved by doing 50 no-op transactions or by making ebrelayer not wait.
                    geth_run_args = geth.buid_run_args(self.ethereum_chain_id, dev=True, http_port=geth_http_port,
                        ws_port=self.ethereum_ws_port, mine=True, allow_insecure_unlock=True,
                        rpc_allow_unprotected_txs=True, gas_limit=self.geth_gas_limit)
                    manual_funds_alloc = funds_alloc
                else:
                    # geth "regular" private mode. We need to create a genesis and assign funds there.
                    # geth needs at least one account to run, it has to be unlocked, but doesn't have to be funded.
                    geth_runner_acct = operator_acct
                    geth_runner_password = ""
                    geth.create_account(geth_runner_acct[1], password=geth_runner_password)
                    geth.init(self.ethereum_chain_id, [geth_runner_acct[0]], funds_alloc=funds_alloc,
                        block_mining_period=self.geth_block_mining_period, gas_limit=self.geth_gas_limit)
                    tmp_password_file = self.cmd.mktempfile()
                    files_to_delete.append(tmp_password_file)
                    self.cmd.write_text_file(tmp_password_file, geth_runner_password)
                    # TODO Disable rpc_allow_unprotected_txs and fix the cause of "only replay-protected (EIP-155)
                    #      transactions allowed over RPC" when sending transactions
                    geth_run_args = geth.buid_run_args(self.ethereum_chain_id, http_port=geth_http_port,
                        ws_port=self.ethereum_ws_port, mine=True, unlock=[geth_runner_acct[0]],
                        password=tmp_password_file, allow_insecure_unlock=True, rpc_allow_unprotected_txs=True,
                        gas_limit=self.geth_gas_limit)

                geth_proc = self.cmd.spawn_asynchronous_process(geth_run_args, log_file=hardhat_log_file)
                hardhat_proc = geth_proc
                hardhat_config_section = "geth"
                hardhat_bind_hostname = "localhost"
                hardhat_deploy_url = "http://localhost:{}/".format(geth_http_port)
                # Accounts for deployments of smart contracts that are passed to smart-contracts/hardhat.config.ts for
                # deployment of smart contracts. The actual script, smart-contracts/scripts/deploy_contracts_dev.ts,
                # needs at least 4 accounts, they are destructured like this:
                # const [operatorAccount, ownerAccount, pauserAccount, validator1Account, ...extraAccounts]
                smart_contract_accounts = [private_key for _, private_key in sample_eth_accounts]
                relayer_extra_args = {"max_fee_per_gas": 300 * eth.GWEI, "max_priority_fee_per_gas": 100 * eth.GWEI}
            else:
                hardhat_bind_hostname = "localhost"  # The host to which to bind to for new connections (Defaults to 127.0.0.1 running locally, and 0.0.0.0 in Docker)
                hardhat_exec_args = hardhat.build_start_args(hostname=hardhat_bind_hostname, port=self.ethereum_ws_port)
                hardhat_proc = self.cmd.spawn_asynchronous_process(hardhat_exec_args, log_file=hardhat_log_file)
                hardhat_config_section = "localhost"
                # hardhat has a blockchain node that will support web socket communication but the node it communicates
                # with internally must be HTTP.
                hardhat_deploy_url = None
                smart_contract_accounts = None  # Provided by hardhat (hardcoded)
                relayer_extra_args = {}

            w3_conn = eth.web3_connect(w3_url)
            eth.web3_wait_for_connection_up(w3_conn)

            # In "geth --dev" mode, funds are not allocated in genesis file, but they need to be manually funded from
            # eth.coinbase. Since we don't have coinbase's private key, we have to use the "local account" API for
            # signing transaction. Special PoA middleware is required for implictly signed transaction (see above).
            if manual_funds_alloc:
                eth.web3_inject_geth_poa_middleware(w3_conn)
                coinbase = w3_conn.eth.coinbase
                log.debug("Ethereum coinbase address: {}".format(coinbase))
                for addr, amount in manual_funds_alloc.items():
                    txhash = w3_conn.eth.send_transaction({"from": coinbase, "to": addr, "value": amount})
                    w3_conn.eth.wait_for_transaction_receipt(txhash)

            balances_check = {a[0]: w3_conn.eth.get_balance(a[0]) for a in sample_eth_accounts}
            assert balances_check[hardhat_accounts["owner"][0]] >= 1 * eth.ETH
            assert balances_check[hardhat_accounts["operator"][0]] >= 1 * eth.ETH

            gas_price = w3_conn.eth.gas_price
            log.debug("Gas price: {}".format(gas_price))
        finally:
            for f in files_to_delete:
                self.cmd.rm(f)

        hardhat_scripts = hardhat.script_runner(url=hardhat_deploy_url, network=hardhat_config_section,
            accounts=smart_contract_accounts)

        # Deploy smart ocntracts using hardhat.
        # Note: if the contracts were compiled previously for hardhat, or if previous deployment failed, you might
        # have to remove smart-contracts/{build,cache,artifacts,.openzeppelin}
        peggy_sc_addrs = hardhat_scripts.deploy_smart_contracts()

        admin_account_name = "sifnodeadmin"
        validator_power = self.witness_power
        seed_ip_address = "10.10.1.1"
        tendermint_port = 26657
        denom_whitelist_file = project_dir("test", "integration", "whitelisted-denoms.json")
        registry_json = project_dir("smart-contracts", "src", "devenv", "registry.json")
        sifnoded_network_dir = self.cmd.tmpdir("sifnodedNetwork")  # Gets written to .vscode/launch.json
        self.cmd.rmdir(sifnoded_network_dir)
        self.cmd.mkdir(sifnoded_network_dir)
        network_config_file, sifnoded_exec_args, sifnoded_proc, tcp_url, admin_account_address, sifnode_validators, \
            sifnode_relayers, sifnode_witnesses, sifnode_validator0_home, chain_dir = \
                self.init_sifchain(sifnoded_network_dir, sifnoded_log_file, self.validator_mint_amounts,
                    validator_power, seed_ip_address, tendermint_port, denom_whitelist_file,
                    self.admin_account_mint_amounts, registry_json, admin_account_name, self.ceth_symbol)

        log.debug("Smart contracts operator: {}".format(operator_addr))
        log.debug("ceth symbol: {}".format(self.ceth_symbol))
        log.debug("Admin account address (rowan source): {}".format(admin_account_address))  # tokens
        log.debug("Witness count: {}".format(self.witness_count))
        log.debug("Consensus thresholds {}".format(self.consensus_threshold))
        log.debug("Validator 0 address: {}".format(sifnode_validators[0]["address"]))  # mint
        for sc_name, sc_address in peggy_sc_addrs.items():
            log.debug("{} smart contract address: {}".format(sc_name, sc_address))

        evm_validator_accounts = hardhat_accounts["validators"]
        evm_validator_addresses = [address[0] for address in evm_validator_accounts]

        symbol_translator_file = os.path.join(self.test_integration_dir, "config", "symbol_translator.json")
        [relayer0_exec_args], [witness_exec_args] = self.start_witnesses_and_relayers(
            w3_url, self.ethereum_chain_id, tcp_url, self.chain_id, peggy_sc_addrs, evm_validator_accounts,
            sifnode_validators, sifnode_relayers, sifnode_witnesses, symbol_translator_file, self.log_level_relayer,
            self.log_level_witness, relayer_extra_args)

        hardhat_scripts.update_validator_power(peggy_sc_addrs["CosmosBridge"], evm_validator_addresses, sifnode_witnesses)

        relayer0_proc = self.cmd.spawn_asynchronous_process(relayer0_exec_args, log_file=relayer_log_file)
        witness_procs = []
        for i, w in enumerate(witness_exec_args):
            witness_log_file = open(os.path.join(log_dir, f"witness{i}.log"), "w")  # TODO close; will be empty on non-peggy2 branch
            witness_proc = self.cmd.spawn_asynchronous_process(w, log_file=witness_log_file)
            witness_procs.append(witness_proc)

        # In the future, we want to have one descriptor for entire environment.
        # It should be able to support multiple EVM and multiple Cosmos chains, including all neccessary bridges and
        # relayers. For now this is just a prototype which is not used yet.
        _unused_peggy2_environment = {
            "admin": {
                "name": admin_account_name,
                "address": admin_account_address,
                "home": sifnode_validator0_home,
            },
            "symbol_translator_file": symbol_translator_file,
            "w3_url": w3_url,
            "evm_chain_id": self.ethereum_chain_id,
            "chain_id": self.chain_id,
            "validators": sifnode_validators,  # From yaml file generated by sifgen
            "relayers": sifnode_relayers,
            "smart_contracts": peggy_sc_addrs
        }

        self.write_env_files(self.project.project_dir(), self.project.go_bin_dir, peggy_sc_addrs, hardhat_accounts,
            admin_account_name, admin_account_address, sifnode_validator0_home, sifnode_validators, sifnode_relayers,
            sifnode_witnesses, tcp_url, hardhat_bind_hostname, hardhat_port, self.ethereum_chain_id, chain_dir,
            sifnoded_exec_args, relayer0_exec_args, witness_exec_args
        )

        return hardhat_proc, sifnoded_proc, relayer0_proc, witness_procs

    def init_sifchain(self, sifnoded_network_dir: str, sifnoded_log_file: TextIO,
        validator_mint_amounts: cosmos.Balance, validator_power: int, seed_ip_address: str, tendermint_port: int,
        denom_whitelist_file: str, admin_account_mint_amounts: cosmos.Balance, registry_json: str,
        admin_account_name: str, ceth_symbol: str
    ) -> Tuple[str, command.ExecArgs, subprocess.Popen, str, cosmos.Address, List, List, List, str, str]:
        validator_count = 1
        relayer_count = 1
        witness_count = self.witness_count

        network_config_file_path = self.cmd.mktempfile()
        try:
            self.cmd.sifgen_create_network(self.chain_id, validator_count, sifnoded_network_dir, network_config_file_path,
                seed_ip_address, mint_amount=validator_mint_amounts)
            network_config_file = self.cmd.read_text_file(network_config_file_path)
        finally:
            self.cmd.rm(network_config_file_path)
        validators = yaml_load(network_config_file)

        # netdef_yml is a list of generated validators like below.
        # Each one has its unique IP (starting with base IP + 1), the first one also has is_seed=True.
        #
        # class ValidatorValues:
        #     chain_id: str
        #     node_id: str
        #     ipv4_address: str
        #     moniker: str
        #     password: str
        #     address: str
        #     pub_key: str
        #     mnemonic: str
        #     validator_address: str
        #     validator_consensus_address: str
        #     is_seed: bool
        assert len(validators) == validator_count

        chain_dir_base = os.path.join(sifnoded_network_dir, "validators", self.chain_id)

        for validator in validators:
            validator_moniker = validator["moniker"]
            validator_mnemonic = validator["mnemonic"].split(" ")
            # TODO Not used
            # validator_password = validator["password"]
            evm_network_descriptor = 1  # TODO Why not hardhat_chain_id?
            sifnoded_home = os.path.join(chain_dir_base, validator_moniker, ".sifnoded")
            sifnode = Sifnoded(self.cmd, home=sifnoded_home)
            self.cmd.sifnoded_peggy2_init_validator(sifnode, validator_moniker, validator_mnemonic, evm_network_descriptor, validator_power)

        # TODO Needs to be fixed when we support more than 1 validator
        validator0 = exactly_one(validators)
        validator0_home = os.path.join(chain_dir_base, validator0["moniker"], ".sifnoded")
        validator0_address = validator0["address"]
        chain_dir = os.path.join(chain_dir_base, validator0["moniker"])

        sifnode = Sifnoded(self.cmd, home=validator0_home, chain_id=self.chain_id)

        # Create an ADMIN account on sifnode with name admin_account_name (e.g. "sifnodeadmin")
        admin_account_address = sifnode.peggy2_add_account(admin_account_name, admin_account_mint_amounts, is_admin=True)

        if self.extra_balances_for_admin_account:
            sifnode.add_genesis_account_directly_to_existing_genesis_json({admin_account_address: self.extra_balances_for_admin_account})

        # TODO Check if sifnoded_peggy2_add_relayer_witness_account can be executed offline (without sifnoded running)
        # TODO Check if sifnoded_peggy2_set_cross_chain_fee can be executed offline (without sifnoded running)

        # Create an account for each relayer
        # Note: "--home" is shared with sifnoded's "--home"
        relayers = [{
            "name": name,
            "address": sifnode.peggy2_add_relayer_witness_account(name, admin_account_mint_amounts,
                self.ethereum_chain_id, validator_power, denom_whitelist_file),
            "home": validator0_home,
        } for name in [f"relayer-{i}" for i in range(relayer_count)]]

        # Create an account for each witness
        # Note: "--home" is shared with sifnoded's "--home"
        witnesses = [{
            "name": name,
            "address": sifnode.peggy2_add_relayer_witness_account(name, admin_account_mint_amounts,
                self.ethereum_chain_id, validator_power, denom_whitelist_file),
            "home": validator0_home,
            "power": validator_power
        } for name in [f"witness-{i}" for i in range(witness_count)]]

        tcp_url = "tcp://{}:{}".format(ANY_ADDR, tendermint_port)
        gas_prices = (0.5, ROWAN)
        # @TODO Detect if sifnoded is already running, for now it fails silently and we wait forever in wait_for_sif_account_up
        sifnoded_exec_args = sifnode.build_start_cmd(tcp_url=tcp_url, minimum_gas_prices=gas_prices,
            log_format_json=True, log_level="debug")
        sifnoded_proc = self.cmd.spawn_asynchronous_process(sifnoded_exec_args, log_file=sifnoded_log_file)

        time_before = time.time()
        self.cmd.wait_for_sif_account_up(validator0_address, tcp_url)
        log.debug("Time for sifnoded to come up: {:.2f}s".format(time.time() - time_before))

        # TODO This command exits with status 0, but looks like there are some errros.
        # The same happens also in devenv.
        # TODO Try whitelister account instead of admin
        res = sifnode.peggy2_token_registry_set_registry(registry_json, gas_prices, 1.5, admin_account_address,
            self.chain_id)
        log.debug("Result from token registry: {}".format(repr(res)))
        assert len(res) == 1
        assert res[0]["code"] == 0

        # We need wait for last tx wrapped up in block, otherwise we could get a wrong sequence, resulting in invalid
        # signatures. This delay waits for block production. (See commit 5854d8b6f3970c1254cac0eca0e3817354151853)
        sifnode.wait_for_last_transaction_to_be_mined()
        cross_chain_fee_base = 1
        cross_chain_lock_fee = 1
        cross_chain_burn_fee = 1
        ethereum_cross_chain_fee_token = ceth_symbol
        assert self.ethereum_chain_id == int(ethereum_cross_chain_fee_token[9:13])  # Assume they should match
        gas_prices = [0.5, ROWAN]
        gas_adjustment = 1.5
        sifnode.peggy2_set_cross_chain_fee(admin_account_address, self.ethereum_chain_id,
            ethereum_cross_chain_fee_token, cross_chain_fee_base, cross_chain_lock_fee, cross_chain_burn_fee,
            admin_account_name, gas_prices, gas_adjustment)

        # We need wait for last tx wrapped up in block, otherwise we could get a wrong sequence, resulting in invalid
        # signatures. This delay waits for block production. (See commit 5854d8b6f3970c1254cac0eca0e3817354151853)
        sifnode.wait_for_last_transaction_to_be_mined()
        sifnode.peggy2_update_consensus_needed(admin_account_address, self.ethereum_chain_id, self.consensus_threshold)

        return network_config_file, sifnoded_exec_args, sifnoded_proc, tcp_url, admin_account_address, validators, \
            relayers, witnesses, validator0_home, chain_dir

    def start_witnesses_and_relayers(self, web3_websocket_address: str, hardhat_chain_id: int, tcp_url: str,
        chain_id: str, peggy_sc_addrs: Mapping[str, eth.Address], evm_validator_accounts: List,
        sifnode_validators: Sequence[Mapping[str, Any]], sifnode_relayers: Sequence[Mapping[str, Any]],
        sifnode_witnesses: Sequence[Mapping[str, Any]], symbol_translator_file: str, log_level_relayer: Optional[str],
        log_level_witness: Optional[str], relayer_extra_args: Mapping[str, Any]
    ) -> Tuple[Sequence[command.ExecArgs], Sequence[command.ExecArgs]]:
        # For now we assume a single validator
        evm_validator0_addr, evm_validator0_key = evm_validator_accounts[0]

        sifnode_validator0 = exactly_one(sifnode_validators)
        sifnode_validator0_address = sifnode_validator0["address"]

        sifnode_relayer0 = exactly_one(sifnode_relayers)
        sifnode_relayer0_mnemonic = sifnode_relayer0["name"]
        sifnode_relayer0_address = sifnode_relayer0["address"]
        sifnode_relayer0_home = sifnode_relayer0["home"]

        bridge_registry_contract_addr = peggy_sc_addrs["BridgeRegistry"]

        self.cmd.wait_for_sif_account_up(sifnode_validator0_address, tcp_url=tcp_url)  # Required for both relayer and witness

        ebrelayer = Ebrelayer(self.cmd)

        relayer0_exec_args = ebrelayer.peggy2_build_ebrelayer_cmd(
            "init-relayer",
            hardhat_chain_id,
            tcp_url,
            web3_websocket_address,
            bridge_registry_contract_addr,
            sifnode_relayer0_mnemonic,
            chain_id=chain_id,
            node=tcp_url,
            sign_with=sifnode_relayer0_address,
            symbol_translator_file=symbol_translator_file,
            ethereum_address=evm_validator0_addr,
            ethereum_private_key=evm_validator0_key,
            keyring_backend="test",
            keyring_dir=sifnode_relayer0_home,
            home=sifnode_relayer0_home,
            log_level=log_level_relayer,
            **relayer_extra_args,
        )

        witness_exec_args = []
        for i, w in enumerate(sifnode_witnesses):
            item = ebrelayer.peggy2_build_ebrelayer_cmd(
                "init-witness",
                hardhat_chain_id,
                tcp_url,
                web3_websocket_address,
                bridge_registry_contract_addr,
                w["name"],
                chain_id=chain_id,
                node=tcp_url,
                sign_with=w["address"],
                symbol_translator_file=symbol_translator_file,
                ethereum_address=evm_validator_accounts[i][0],
                ethereum_private_key=evm_validator_accounts[i][1],
                keyring_backend="test",
                keyring_dir=w["home"],
                log_format="json",
                home=w["home"],
                log_level=log_level_witness,
            )
            witness_exec_args.append(item)

        return [relayer0_exec_args], [witness_exec_args]

    def write_env_files(self, project_dir, go_bin_dir, evm_smart_contract_addrs, eth_accounts, admin_account_name,
        admin_account_address, sifnode_validator0_home, sifnode_validators, sifnode_relayers, sifnode_witnesses,
        tcp_url, hardhat_bind_hostname, hardhat_port, hardhat_chain_id, chain_dir, sifnoded_exec_args,
        relayer0_exec_args, witness0_exec_args
    ):
        w3_url = eth.web3_host_port_url(hardhat_bind_hostname, hardhat_port)

        # @TODO At the moment, values are fed from one rendered template into the next.
        #       We should use values directly from parameters instead.

        def format_eth_account(eth_account):
            assert eth_account[0].startswith("0x") and not eth_account[1].startswith("0x")
            return {"address": eth_account[0], "privateKey": "0x" + eth_account[1]}

        def format_sif_account(sif_account):
            return {"account": sif_account["address"], "name": sif_account["name"], "homeDir": sif_account["home"]}

        environment_json = {
            "contractResults": {
                # "completed": True,
                # "output": "...",
                "contractAddresses": {
                    "cosmosBridge": evm_smart_contract_addrs["CosmosBridge"],
                    "bridgeBank": evm_smart_contract_addrs["BridgeBank"],
                    "bridgeRegistry": evm_smart_contract_addrs["BridgeRegistry"],
                    "rowanContract": evm_smart_contract_addrs["Rowan"],
                    "blocklist": evm_smart_contract_addrs["Blocklist"],
                }
            },
            "ethResults": {
                # "process": { ... },
                "accounts": {
                    "proxyAdmin": format_eth_account(eth_accounts["proxy_admin"]),
                    "operator": format_eth_account(eth_accounts["operator"]),
                    "owner": format_eth_account(eth_accounts["owner"]),
                    "pauser": format_eth_account(eth_accounts["pauser"]),
                    "validators": [format_eth_account(a) for a in eth_accounts["validators"]],
                    "available": [format_eth_account(a) for a in eth_accounts["available"]],
                },
                "httpHost": hardhat_bind_hostname,
                "httpPort": hardhat_port,
                "chainId": self.ethereum_chain_id,
            },
            "goResults": {
                # "completed": True,
                # "output": "...",
                "goBin": go_bin_dir
            },
            "sifResults": {
                "validatorValues": sifnode_validators,
                "adminAddress": format_sif_account({"address": admin_account_address, "name": admin_account_name, "home": sifnode_validator0_home}),
                "relayerAddresses": [format_sif_account(a) for a in sifnode_relayers],
                "witnessAddresses": [format_sif_account(a) for a in sifnode_witnesses],
                # "process": { ... },
                "tcpurl": tcp_url,
            }
        }

        # TODO Inconsistent format of deployed smart contract addresses (this was intentionally carried over from
        #      devenv to preserve compatibility with devenv users)
        # TODO Convert to our "unified" json file format

        # TODO Do we want "0x" prefixes here for private keys?
        dot_env = dict_merge({
            "BASEDIR": project_dir,
            "ETHEREUM_ADDRESS": eth_accounts["available"][0][0],
            "ETHEREUM_PRIVATE_KEY": "0x" + eth_accounts["available"][0][1],
            "ETH_ACCOUNT_OPERATOR_ADDRESS": eth_accounts["operator"][0],
            "ETH_ACCOUNT_OPERATOR_PRIVATEKEY": "0x" + eth_accounts["operator"][1],
            "ETH_ACCOUNT_OWNER_ADDRESS": eth_accounts["owner"][0],
            "ETH_ACCOUNT_OWNER_PRIVATEKEY": "0x" + eth_accounts["owner"][1],
            "ETH_ACCOUNT_PAUSER_ADDRESS": eth_accounts["pauser"][0],
            "ETH_ACCOUNT_PAUSER_PRIVATEKEY": "0x" + eth_accounts["pauser"][1],
            "ETH_ACCOUNT_PROXYADMIN_ADDRESS": eth_accounts["proxy_admin"][0],
            "ETH_ACCOUNT_PROXYADMIN_PRIVATEKEY": "0x" + eth_accounts["proxy_admin"][1],
            "ETH_ACCOUNT_VALIDATOR_ADDRESS": eth_accounts["validators"][0][0],
            "ETH_ACCOUNT_VALIDATOR_PRIVATEKEY": "0x" + eth_accounts["validators"][0][1],
            "ETH_CHAIN_ID": str(self.ethereum_chain_id),
            "ETH_HOST": hardhat_bind_hostname,
            "ETH_PORT": str(hardhat_port),
            "ROWAN_SOURCE": admin_account_address,
            "BRIDGE_BANK_ADDRESS": evm_smart_contract_addrs["BridgeBank"],
            # "BRIDGE_REGISTRY_ADDRESS": evm_smart_contract_addrs["BridgeRegistry"],
            "BRIDGE_REGISTERY_ADDRESS": evm_smart_contract_addrs["BridgeRegistry"], # TODO Typo, remove, keeping it for now for compatibility
            "COSMOS_BRIDGE_ADDRESS": evm_smart_contract_addrs["CosmosBridge"],
            "ROWANTOKEN_ADDRESS": evm_smart_contract_addrs["Rowan"],
            "BRIDGE_TOKEN_ADDRESS": evm_smart_contract_addrs["Rowan"],
            "GOBIN": go_bin_dir,
            "TCP_URL": tcp_url,
            "VALIDATOR_ADDRESS": sifnode_validators[0]["address"],
            "VALIDATOR_CONSENSUS_ADDRESS": sifnode_validators[0]["validator_consensus_address"],
            "VALIDATOR_MENOMONIC": sifnode_validators[0]["mnemonic"],
            "VALIDATOR_MONIKER": sifnode_validators[0]["moniker"],
            "VALIDATOR_PASSWORD": sifnode_validators[0]["password"],
            "VALIDATOR_PUB_KEY": sifnode_validators[0]["pub_key"],
            "ADMIN_ADDRESS": admin_account_address,
            "ADMIN_NAME": admin_account_name,
            "CHAINDIR": chain_dir,
            "HOME": chain_dir,
        }, *[{
            f"ETHEREUM_ADDRESS_{i}": account[0],
            f"ETHEREUM_PRIVATE_KEY_{i}": "0x" + account[1],
        } for i, account in enumerate(eth_accounts["available"])], *[{
            f"ETH_ACCOUNT_VALIDATOR_{i}_ADDRESS": account[0],
            f"ETH_ACCOUNT_VALIDATOR_{i}_PRIVATEKEY": "0x" + account[1],
        } for i, account in enumerate(eth_accounts["validators"])], *[{
            f"VALIDATOR_ADDRESS_{i}": validator["address"],
            f"VALIDATOR_CONSENSUS_ADDRESS_{i}": validator["validator_consensus_address"],
            f"VALIDATOR_MENOMONIC_{i}": validator["mnemonic"],
            f"VALIDATOR_MONIKER_{i}": validator["moniker"],
            f"VALIDATOR_PASSWORD_{i}": validator["password"],
            f"VALIDATOR_PUB_KEY_{i}": validator["pub_key"],
        } for i, validator in enumerate(sifnode_validators)], *[{
            f"RELAYER_ADDRESS_{i}": relayer["address"],
            f"RELAYER_NAME_{i}": relayer["name"],
        } for i, relayer in enumerate(sifnode_relayers)], *[{
            f"WITNESS_ADDRESS_{i}": witness["address"],
            f"WITNESS_NAME_{i}": witness["name"],
        } for i, witness in enumerate(sifnode_witnesses)])

        # launch.json for VS Code
        launch_json = {
            "version": "0.2.0",
            "configurations": [
                {
                    "name": "Python: siftool",
                    "type": "python",
                    "request": "launch",
                    "program": "${workspaceFolder}/test/integration/framework/siftool",
                    "args": ["run-env"],
                    "console": "integratedTerminal",
                    "justMyCode": True,
                    "env": {
                        "WITNESS_COUNT": "1",
                        "CONSENSUS_THRESHOLD": "49"
                    }
                }, {
                    "runtimeArgs": ["node_modules/.bin/hardhat", "run"],
                    "cwd": "${workspaceFolder}/smart-contracts",
                    "type": "node",
                    "request": "launch",
                    "name": "Debug DevENV scripts",
                    "skipFiles": ["<node_internals>/**"],
                    "program": "${workspaceFolder}/smart-contracts/scripts/devenv.ts",
                }, {
                    "runtimeArgs": ["node_modules/.bin/hardhat", "test"],
                    "args": ["--network", "localhost"],
                    "cwd": "${workspaceFolder}/smart-contracts",
                    "type": "node",
                    "request": "launch",
                    "name": "Typescript Tests",
                    "skipFiles": ["<node_internals>/**"],
                    "program": "${workspaceFolder}/smart-contracts/test/watcher/watcher.ts",
                }, *[{
                    "name": f"Debug Relayer-{i}",
                    "type": "go",
                    "request": "launch",
                    "mode": "debug",
                    "program": "cmd/ebrelayer",
                    "envFile": "${workspaceFolder}/smart-contracts/.env",
                    # Generally we want to use relayer0_exec_args, except for:
                    # - here we don't have the initial "ebrelayer"
                    # - here we are using "${workspaceFolder} for" "--symbol-translator-file"
                    # - here we don't have ETHEREUM_ADDRESS env
                    # TODO Probable bug, should be "eth_accounts["validators"][0][1]}"
                    "env": {"ETHEREUM_PRIVATE_KEY": eth_accounts["available"][i][1]},
                    # "env": {"ETHEREUM_PRIVATE_KEY": eth_accounts["validators"][0][1]},
                    "args": [
                        "init-relayer",
                        "--network-descriptor", str(self.ethereum_chain_id),
                        "--tendermint-node", tcp_url,
                        "--web3-provider", w3_url,
                        "--bridge-registry-contract-address", evm_smart_contract_addrs["BridgeRegistry"],
                        "--validator-mnemonic", relayer["name"],
                        "--chain-id", "localnet",
                        "--node", tcp_url,
                        "--keyring-backend", "test",
                        "--from", relayer["address"],
                        "--symbol-translator-file", "${workspaceFolder}/test/integration/config/symbol_translator.json",
                        "--keyring-dir", relayer["home"],
                    ]
                } for i, relayer in enumerate(sifnode_relayers)], *[{
                    "name": f"Debug Witness-{i}",
                    "type": "go",
                    "request": "launch",
                    "mode": "debug",
                    "program": "cmd/ebrelayer",
                    "envFile": "${workspaceFolder}/smart-contracts/.env",
                    # Generally we want to use witness0_exec_args, except for:
                    # - here we don't have the initial "ebrelayer"
                    # - here we are using "${workspaceFolder} for" "--symbol-translator-file"
                    # - here we don't have ETHEREUM_ADDRESS env
                    # TODO Probable bug, should be "eth_accounts["validators"][0][1]}"
                    "env": {"ETHEREUM_PRIVATE_KEY": eth_accounts["available"][i][1]},
                    # "env": {"ETHEREUM_PRIVATE_KEY": eth_accounts["validators"][0][1]},
                    "args": [
                        "init-witness",
                        "--network-descriptor", str(self.ethereum_chain_id),
                        "--tendermint-node", tcp_url,
                        "--web3-provider", w3_url,
                        "--bridge-registry-contract-address", evm_smart_contract_addrs["BridgeRegistry"],
                        "--validator-mnemonic", witness["name"],
                        "--chain-id", "localnet",
                        "--node", tcp_url,
                        "--keyring-backend", "test",
                        "--from", witness["address"],
                        # TODO: This shouldnt be needed, it defaults to --home value
                        "--keyring-dir", witness["home"],
                        "--symbol-translator-file", "${workspaceFolder}/test/integration/config/symbol_translator.json",
                        "--home", witness["home"]
                    ]
                } for i, witness in enumerate(sifnode_witnesses)], {
                    "name": "Debug Sifnoded",
                    "type": "go",
                    "request": "launch",
                    "mode": "debug",
                    "program": "cmd/sifnoded",
                    # Generally we want to use sifnoded_exec_args, except for:
                    # - here we don't have the initial "sifnoded"
                    # TODO Do not use .env file here
                    "envFile": "${workspaceFolder}/smart-contracts/.env",
                    "args": [
                        "start",
                        "--log_format", "json",
                        "--log_level", "debug",
                        "--minimum-gas-prices", f"0.5{ROWAN}",
                        "--rpc.laddr", tcp_url,
                        "--home", sifnode_validator0_home
                    ]
                }, {
                    "name": "Pytest",
                    "type": "python",
                    "request": "launch",
                    "stopOnEntry": False,
                    "python": "${command:python.interpreterPath}",
                    "module": "pytest",
                    "args": [
                        "-olog_cli=false",
                        "-olog_level=DEBUG",
                        "-olog_file=/tmp/pytest.out",
                        "-v",
                        "test/integration/src/py/test_eth_transfers.py::test_eth_to_ceth_and_back_to_eth"
                    ],
                    "cwd": "${workspaceRoot}",
                    "env": dot_env,
                    "debugOptions": ["WaitOnAbnormalExit", "WaitOnNormalExit", "RedirectOutput"]
                }
            ]
        }

        # IntelliJ
        def render_intellij_run_xml(name, joined_args, package, filepath, envs):
            def q(s): return s  # TODO Qoute/escape

            # since the contents is being fed from launch_json, we have ${workspaceFolder} in joined_args
            joined_args = re.sub("\\${workspaceFolder}/", "", joined_args)

            return [
                "<component name=\"ProjectRunConfigurationManager\">",
                "  <configuration default=\"false\" name=\"" + q(name) + "\" type=\"GoApplicationRunConfiguration\" factoryName=\"Go Application\">",
                "    <module name=\"sifnode\" />",
                "    <working_directory value=\"$PROJECT_DIR$\" />",
                "    <parameters value=\"" + q(joined_args) + "\" />",
            ] + ((
                ["    <envs>"] +
                ["      <env name=\"" + q(k) + "\" value=\"" + q(v) + "\" />" for k, v in envs.items()] +
                ["    </envs>"]
            ) if envs else []) + [
                "    <kind value=\"PACKAGE\" />",
                "    <package value=\"" + q(package) + "\" />",
                "    <directory value=\"$PROJECT_DIR$\" />",
                "    <filePath value=\"" + q(filepath) + "\" />",
                "    <method v=\"2\" />",
                "  </configuration>",
                "</component>",
            ]

        intellij_sifnoded_configs = []
        intellij_ebrelayer_configs = []
        intellij_witness_configs = []
        for config in launch_json["configurations"]:
            if config["name"].startswith("Debug Relayer"):
                intellij_ebrelayer_configs.append(render_intellij_run_xml(
                    "ebrelayer devenv",
                    " ".join(config["args"]),
                    "github.com/Sifchain/sifnode/cmd/ebrelayer",
                    "$PROJECT_DIR$/cmd/ebrelayer/main.go",
                    {"ETHEREUM_PRIVATE_KEY": dot_env["ETHEREUM_PRIVATE_KEY"]}))
            elif config["name"].startswith("Debug Witness"):
                intellij_witness_configs.append(render_intellij_run_xml(
                    "witness devenv",
                    " ".join(config["args"]),
                    "github.com/Sifchain/sifnode/cmd/ebrelayer",
                    "$PROJECT_DIR$/cmd/ebrelayer/main.go",
                    {"ETHEREUM_PRIVATE_KEY": dot_env["ETHEREUM_PRIVATE_KEY"]}))
            elif config["name"].startswith("Debug Sifnoded"):
                intellij_sifnoded_configs.append(render_intellij_run_xml(
                    "sifnoded devenv",
                    " ".join(config["args"]),
                    "github.com/Sifchain/sifnode/cmd/sifnoded",
                    "$PROJECT_DIR$/cmd/sifnoded/main.go",
                    {}))

        intellij_ebrelayer_config = exactly_one(intellij_ebrelayer_configs)
        intellij_witness_config = intellij_witness_configs[0]
        intellij_sifnoded_config = exactly_one(intellij_sifnoded_configs)

        run_dir = self.project.project_dir(".run")
        self.cmd.mkdir(run_dir)
        vscode_dir = self.project.project_dir(".vscode")
        self.cmd.mkdir(vscode_dir)

        self.cmd.write_text_file(self.project.project_dir("smart-contracts/environment.json"), json.dumps(environment_json, indent=2))
        self.cmd.write_text_file(self.project.project_dir("smart-contracts/env.json"), json.dumps(dot_env, indent=2))
        self.cmd.write_text_file(self.project.project_dir("smart-contracts/.env"), joinlines(format_as_shell_env_vars(dot_env, export=True)))
        self.cmd.write_text_file(os.path.join(vscode_dir, "launch.json"), json.dumps(launch_json, indent=2))
        self.cmd.write_text_file(os.path.join(run_dir, "ebrelayer.run.xml"), joinlines(intellij_ebrelayer_config))
        self.cmd.write_text_file(os.path.join(run_dir, "witness.run.xml"), joinlines(intellij_witness_config))
        self.cmd.write_text_file(os.path.join(run_dir, "sifnoded.run.xml"), joinlines(intellij_sifnoded_config))

        return environment_json, dot_env, launch_json, intellij_ebrelayer_config, intellij_witness_config, intellij_sifnoded_config


class IBCEnvironment(IntegrationTestsEnvironment):
    def __init__(self, cmd):
        super().__init__(cmd)

    def run(self):
        chainnet0 = "sifchain-ibc-0"
        chainnet1 = "sifchain-ibc-1"
        ipaddr0 = "192.168.65.2"
        ipaddr1 = "192.168.65.3"
        subnet = "192.168.65.1/24"
        # Mnemonics can be generated with "sifgen key generate" or "sifnoded keys mnemonic", but that gives us 24 words
        # and here are only 12.
        # A mnemonic contains both public and private key. Public key is the address, there can only be one such entry
        # in the keyring.
        mnemonic = "toddler spike waste purpose neutral beach science dawn joke stock help beyond"

        sifgen = Sifgen(self.cmd)
        # This does not work - "--keyring-backend" is not supported
        x = sifgen.create_standalone(chainnet0, "chain1", mnemonic, ipaddr0, keyring_backend=None)
