import json
import re
import sys
import time
from truffle import Ganache
from command import Command
from hardhat import Hardhat
from sifchain import Sifnoded, Ebrelayer
from project import Project, killall, force_kill_processes
from common import *


class Integrator(Ganache, Sifnoded, Command):
    def __init__(self):
        super().__init__()  # TODO Which super is this? All of them?
        self.project = Project(self, project_dir())

    def sif_wait_up(self, host, port):
        while True:
            from urllib.error import URLError
            try:
                return self.sifnoded_get_status(host, port)
            except URLError:
                time.sleep(1)

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
            self.execst(["bash", "-c", "set -o posix; IFS=' '; set > {}; source {}; set > {}".format(tmp1, path, tmp2)])
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

    # @TODO Merge
    def sifchain_init_integration(self, validator_moniker, validator_mnemonic, sifnoded_home, denom_whitelist_file, validator1_password):
        # now we have to add the validator key to the test keyring so the tests can send rowan from validator1
        self.sifnoded_keys_add(validator_moniker, validator_mnemonic)
        valoper = self.sifnoded_keys_show(validator_moniker, bech="val", keyring_backend="test", home=sifnoded_home)[0]["address"]
        assert valoper == self.sifnoded_get_val_address(validator_moniker)  # This does not use "home"; if it the assertion holds it could be grouped with sifchain_init_peggy

        self.execst(["sifnoded", "add-genesis-validators", valoper, "--home", sifnoded_home])

        try:
            # Probable bug in test/integration/sifchain_start_daemon.sh:
            # whitelisted_validator=$(yes $VALIDATOR1_PASSWORD | sifnoded keys show --keyring-backend file -a \
            #     --bech val $MONIKER --home $CHAINDIR/.sifnoded)
            # TODO We probably don't need validator1_passsword
            # TODO This could then be merged with "sifnoded_keys_show"
            whitelisted_validator = exactly_one(stdout_lines(self.execst(["sifnoded", "keys", "show",
                "--keyring-backend", "file", "-a", "--bech", "val", validator_moniker, "--home", sifnoded_home],
                stdin=[validator1_password])))
            assert False
            log.info(f"Whitelisted validator: {whitelisted_validator}")
            self.cmd.execst(["sifnoded", "add-genesis-validators", whitelisted_validator, "--home", sifnoded_home])
        except:
            log.error("Failed to get whitelisted validator (probable bug)", exc_info=True)
            assert True

        adminuser_addr = self.sifchain_init_common(denom_whitelist_file, sifnoded_home)
        return adminuser_addr

    # @parameter validator_moniker - from network config
    # @parameter validator_mnemonic - from network config
    def sifchain_init_peggy(self, validator_moniker, validator_mnemonic, sifnoded_home, denom_whitelist_file):
        # Add validator key to test keyring
        _tmp0 = self.sifnoded_keys_add_2(validator_moniker, validator_mnemonic)
        valoper = self.sifnoded_get_val_address(validator_moniker)

        # (0, '', '2021/09/07 05:55:33 AddGenesisValidatorCmd, adding addr: sifvaloper1f5vj6j2mnkaw0yec3ut9at4rkl2u23k2fxtrsv to whitelist: []\n')
        unknown_parameter_1 = 1  # Likely "network_descriptor"
        unknown_parameter_2 = 100  # Likely "power"
        self.sifnoded_add_genesis_validators_peggy(unknown_parameter_1, valoper, unknown_parameter_2, sifnoded_home)

        # Get whitelisted validator
        # TODO Value is not being used
        _whitelisted_validator = self.sifnoded_get_val_address(validator_moniker)
        assert valoper == _whitelisted_validator

        adminuser_addr = self.sifchain_init_common(denom_whitelist_file, sifnoded_home)
        return adminuser_addr

    # Shared between IntegrationEnvironment and PeggyEnvironment
    def sifchain_init_common(self, denom_whitelist_file, sifnoded_home):
        sifnodeadmin_addr = self.sifnoded_keys_add_1("sifnodeadmin")["address"]
        tokens = [[10**20, "rowan"]]
        # Original from peggy:
        # self.cmd.execst(["sifnoded", "add-genesis-account", sifnoded_admin_address, "100000000000000000000rowan", "--home", sifnoded_home])
        self.sifnoded_add_genesis_account(sifnodeadmin_addr, tokens, sifnoded_home=sifnoded_home)
        self.sifnoded_exec(["set-genesis-oracle-admin", sifnodeadmin_addr], sifnoded_home=sifnoded_home)
        self.sifnoded_exec(["set-genesis-whitelister-admin", sifnodeadmin_addr], sifnoded_home=sifnoded_home)
        self.sifnoded_exec(["set-gen-denom-whitelist", denom_whitelist_file], sifnoded_home=sifnoded_home)
        return sifnodeadmin_addr

    def sifgen_create_network(self, chain_id, validator_count, networks_dir, network_definition_file, seed_ip_address, mint_amount=None):
        # Old call (no longer works either):
        # sifgen network create localnet 1 /mnt/shared/work/projects/sif/sifnode/local-tmp/my/deploy/rake/../networks \
        #     192.168.1.2 /mnt/shared/work/projects/sif/sifnode/local-tmp/my/deploy/rake/../networks/network-definition.yml \
        #     --keyring-backend file
        # self.cmd.execst(["sifgen", "network", "create", "localnet", str(validator_count), networks_dir, seed_ip_address,
        #     os.path.join(networks_dir, "network-definition.yml"), "--keyring-backend", "file"])
        # TODO Most likely, this should be "--keyring-backend file"
        args = ["sifgen", "network", "create", chain_id, str(validator_count), networks_dir, seed_ip_address,
            network_definition_file, "--keyring-backend", "test"] + \
            (["--mint-amount", ",".join([sif_format_amount(*x) for x in mint_amount])] if mint_amount else [])
        self.execst(args)


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
            port=7545, network_id=self.network_id, gas_price=20000000000, gas_limit=6721975, host="0.0.0.0")

        # ui/scripts/stack-launch.sh -> ui/scripts/_sif.sh -> ui/chains/sif/launch.sh
        self.cmd.sifnoded_init("test", self.chain_id)
        self.cmd.copy_file(project_dir("ui/chains/sif/app.toml"), os.path.join(self.sifnoded_path, "config/app.toml"))
        log.info(f"Generating deterministic account - {self.shadowfiend_name}...")
        shadowfiend_account = self.cmd.sifnoded_keys_add(self.shadowfiend_name, self.shadowfiend_mnemonic)
        log.info(f"Generating deterministic account - {self.akasha_name}...")
        akasha_account = self.cmd.sifnoded_keys_add(self.akasha_name, self.akasha_mnemonic)
        log.info(f"Generating deterministic account - {self.juniper_name}...")
        juniper_account = self.cmd.sifnoded_keys_add(self.juniper_name, self.juniper_mnemonic)
        shadowfiend_address = shadowfiend_account["address"]
        akasha_address = akasha_account["address"]
        juniper_address = juniper_account["address"]
        assert shadowfiend_address == self.cmd.sifnoded_keys_show(self.shadowfiend_name)[0]["address"]
        assert akasha_address == self.cmd.sifnoded_keys_show(self.akasha_name)[0]["address"]
        assert juniper_address == self.cmd.sifnoded_keys_show(self.juniper_name)[0]["address"]

        tokens_shadowfiend = [[10**29, "rowan"], [10**29, "catk"], [10**29, "cbtk"], [10**29, "ceth"], [10**29, "cusdc"], [10**29, "clink"], [10**26, "stake"]]
        tokens_akasha = [[10**29, "rowan"], [10**29, "catk"], [10**29, "cbtk"], [10**29, "ceth"], [10**29, "cusdc"], [10**29, "clink"], [10**26, "stake"]]
        tokens_juniper = [[10**22, "rowan"], [10**22, "cusdc"], [10**20, "clink"], [10**20, "ceth"]]
        self.cmd.sifnoded_add_genesis_account(shadowfiend_address, tokens_shadowfiend)
        self.cmd.sifnoded_add_genesis_account(akasha_address, tokens_akasha)
        self.cmd.sifnoded_add_genesis_account(juniper_address, tokens_juniper)

        shadowfiend_address_bech_val = self.cmd.sifnoded_keys_show(self.shadowfiend_name, bech="val")[0]["address"]
        self.cmd.sifnoded_add_genesis_validators(shadowfiend_address_bech_val)

        amount = sif_format_amount(10**24, "stake")
        self.cmd.execst(["sifnoded", "gentx", self.shadowfiend_name, amount, f"--chain-id={self.chain_id}",
            f"--keyring-backend={self.keyring_backend}"])

        log.info("Collecting genesis txs...")
        self.cmd.execst(["sifnoded", "collect-gentxs"])
        log.info("Validating genesis file...")
        self.cmd.execst(["sifnoded", "validate-genesis"])

        log.info("Starting test chain...")
        sifnoded_proc = self.cmd.sifnoded_start(minimum_gas_prices=[0.5, "rowan"])  # TODO sifnoded_home=???

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
        self.cmd.sifnoded_tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "catk", [10**5, "rowan"], 10**25, 10**25)
        time.sleep(5)
        log.info("Creating liquidity pool from cbtk:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "cbtk", [10**5, "rowan"], 10**25, 10**25)
        # should now be able to swap from catk:cbtk
        time.sleep(5)
        log.info("Creating liquidity pool from ceth:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "ceth", [10**5, "rowan"], 10**25, 83*10**20)
        # should now be able to swap from x:ceth
        time.sleep(5)
        log.info("Creating liquidity pool from cusdc:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "cusdc", [10**5, "rowan"], 10**25, 10**25)
        time.sleep(5)
        log.info("Creating liquidity pool from clink:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "clink", [10**5, "rowan"], 10**25, 588235*10**18)
        time.sleep(5)
        log.info("Creating liquidity pool from ctest:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(self.chain_id, self.keyring_backend, "akasha", "ctest", [10**5, "rowan"], 10**25, 10**13)

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
            [NULL_ADDRESS, 31 * 10 ** 18],
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
            ethereum_private_key=ethereum_private_key, gas=5*10**12, gas_prices=[0.5, "rowan"])

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
        self.using_ganache_gui = False
        self.peruser_storage_dir = self.cmd.get_user_home(".sifnode-integration")
        self.state_vars = {}
        self.test_integration_dir = project_dir("test/integration")
        self.data_dir = project_dir("test/integration/vagrant/data")
        self.chainnet = "localnet"
        self.tcp_url = "tcp://0.0.0.0:26657"
        self.ethereum_websocket_address = "ws://localhost:7545/"
        self.ganache_mnemonic = ["candy", "maple", "cake", "sugar", "pudding", "cream", "honey", "rich", "smooth",
                "crumble", "sweet", "treat"]

    def prepare(self):
        self.project.make_go_binaries()
        self.project.install_smart_contracts_dependencies()

    def run(self):
        self.cmd.mkdir(self.data_dir)

        self.prepare()

        log_dir = "/tmp/sifnode"
        self.cmd.mkdir(log_dir)
        ganache_log_file = open(os.path.join(log_dir, "ganache.log"), "w")  # TODO close
        sifnoded_log_file = open(os.path.join(log_dir, "sifnoded.log"), "w")  # TODO close
        ebrelayer_log_file = open(os.path.join(log_dir, "ebrelayer.log"), "w")  # TODO close

        if self.using_ganache_gui:
            ebrelayer_ethereum_addr = "0x8e2bE12daDbCcbf7c98DBb59f98f22DFF0eF3F2c"
            ebrelayer_ethereum_private_key = "2eaddbc0bca859ff5b09c5a48a2feaeaf464f7cbf8ddbfa4a32a625a8322fe99"
            ganache_db_path = None
            ganache_proc = None
        else:
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
            ganache_proc = Ganache.start_ganache_cli(self.cmd, block_time=block_time, host="0.0.0.0",
                mnemonic=self.ganache_mnemonic, network_id=self.network_id, port=7545, db=ganache_db_path,
                account_keys_path=account_keys_path, log_file=ganache_log_file)

            self.cmd.wait_for_file(account_keys_path)  # Created by ganache-cli
            time.sleep(2)

            ganache_keys = json.loads(self.cmd.read_text_file(account_keys_path))
            ebrelayer_ethereum_addr = list(ganache_keys["private_keys"].keys())[9]
            ebrelayer_ethereum_private_key = ganache_keys["private_keys"][ebrelayer_ethereum_addr]
            # TODO Check for possible non-determinism of dict().keys() ordering (c.f. test/integration/vagrantenv.sh)
            # TODO ebrelayer_ethereum_private_key is NOT the same as in test/integration/.env.ciExample
            assert ebrelayer_ethereum_addr == "0x5aeda56215b167893e80b4fe645ba6d5bab767de"
            assert ebrelayer_ethereum_private_key == "8d5366123cb560bb606379f90a0bfd4769eecc0557f1b362dcae9012b548b1e5"

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
        #     [NULL_ADDRESS, 31*10**18],
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
        mint_amount = [[999999 * 10**21, "rowan"], [137 * 10**16, "ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE"]]

        self.cmd.sifgen_create_network(chain_id, validator_count, networks_dir, network_definition_file, seed_ip_address, mint_amount=mint_amount)

        netdef, netdef_json = self.process_netdef(network_definition_file)

        validator1_moniker = netdef["moniker"]
        validator1_address = netdef["address"]
        validator1_password = netdef["password"]
        validator1_mnemonic = netdef["mnemonic"].split(" ")
        chaindir = os.path.join(networks_dir, f"validators/{self.chainnet}/{validator1_moniker}")
        sifnoded_home = os.path.join(chaindir, ".sifnoded")
        denom_whitelist_file = os.path.join(self.test_integration_dir, "whitelisted-denoms.json")
        # SIFNODED_LOG=$datadir/logs/sifnoded.log

        adminuser_addr = self.cmd.sifchain_init_integration(validator1_moniker, validator1_mnemonic, sifnoded_home, denom_whitelist_file, validator1_password)

        # Start sifnoded
        sifnoded_proc = self.cmd.sifnoded_start(tcp_url=self.tcp_url, minimum_gas_prices=[0.5, "rowan"],
            sifnoded_home=sifnoded_home, log_file=sifnoded_log_file)

        # TODO: should we wait for sifnoded to come up before continuing? If so, how do we do it?

        # TODO Process exits immediately with returncode 1
        # TODO Why does it not stop start-integration-env.sh?
        # rest_server_proc = self.cmd.popen(["sifnoded", "rest-server", "--laddr", "tcp://0.0.0.0:1317"])  # TODO cwd

        # test/integration/sifchain_start_ebrelayer.sh -> test/integration/sifchain_run_ebrelayer.sh
        # This script is also called from tests

        relayer_db_path = os.path.join(self.test_integration_dir, "sifchainrelayerdb")
        ebrelayer_proc = self.run_ebrelayer(netdef_json, validator1_address, validator1_moniker, validator1_mnemonic,
            ebrelayer_ethereum_private_key, bridge_registry_sc_addr, relayer_db_path, log_file=ebrelayer_log_file)

        vagrantenv_path = project_dir("test/integration/vagrantenv.sh")
        self.state_vars = {
            "ETHEREUM_PRIVATE_KEY": self.ethereum_private_key,
            "OWNER": self.owner,
            "PAUSER": self.pauser,
            "BASEDIR": project_dir(),
            # export SIFCHAIN_BIN="/home/jurez/work/projects/sif/sifnode/local/cmd"
            "envexportfile": vagrantenv_path,
            # export TEST_INTEGRATION_DIR="/home/jurez/work/projects/sif/sifnode/local/test/integration"
            # export TEST_INTEGRATION_PY_DIR="/home/jurez/work/projects/sif/sifnode/local/test/integration/src/py"
            "SMART_CONTRACTS_DIR": self.project.smart_contracts_dir,
            # export datadir="/home/jurez/work/projects/sif/sifnode/local/test/integration/vagrant/data"
            # export CONTAINER_NAME="integration_sifnode1_1"
            "NETWORKDIR": networks_dir,
            # export ETHEREUM_WEBSOCKET_ADDRESS="ws://localhost:7545/"
            # export CHAINNET="localnet"
            "GANACHE_DB_DIR": ganache_db_path,
            # export GANACHE_KEYS_JSON="/home/jurez/work/projects/sif/sifnode/local/test/integration/vagrant/data/ganachekeys.json"
            "EBRELAYER_ETHEREUM_ADDR": ebrelayer_ethereum_addr,
            "EBRELAYER_ETHEREUM_PRIVATE_KEY": ebrelayer_ethereum_private_key,  # Needed by sifchain_run_ebrelayer.sh
            # # BRIDGE_REGISTRY_ADDRESS and ETHEREUM_CONTRACT_ADDRESS are synonyms
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

        return ganache_proc, sifnoded_proc, ebrelayer_proc

    def wait_for_sif_account(self, netdef_json, validator1_address):
        return self.cmd.execst(["python3", os.path.join(self.test_integration_dir, "src/py/wait_for_sif_account.py"),
            netdef_json, validator1_address], env={"USER1ADDR": "nothing"})

    def remove_and_add_sifnoded_keys(self, moniker, mnemonic):
        # Error: The specified item could not be found in the keyring
        # This is not neccessary during start-integration-env.sh (as the key does not exist yet), but is neccessary
        # during tests that restart ebrelayer
        # res = self.cmd.execst(["sifnoded", "keys", "delete", moniker, "--keyring-backend", "test"], stdin=["y"])
        self.cmd.sifnoded_keys_delete(moniker)
        self.cmd.sifnoded_keys_add(moniker, mnemonic)

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
        self.wait_for_sif_account(netdef_json, validator1_address)
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

        ganache_proc = Ganache.start_ganache_cli(self.cmd, block_time=block_time, host="0.0.0.0",
            mnemonic=self.ganache_mnemonic, network_id=self.network_id, port=7545, db=ganache_db_path,
            account_keys_path=account_keys_path)  # TODO log_file

        self.cmd.wait_for_file(account_keys_path)  # Created by ganache-cli
        time.sleep(2)

        validator_moniker = self.state_vars["MONIKER"]
        networks_dir = project_dir("deploy/networks")
        chaindir = os.path.join(networks_dir, f"validators/{self.chainnet}/{validator_moniker}")
        sifnoded_home = os.path.join(chaindir, ".sifnoded")
        sifnoded_proc = self.cmd.sifnoded_start(tcp_url=self.tcp_url, minimum_gas_prices=[0.5, "rowan"], sifnoded_home=sifnoded_home)

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


class PeggyEnvironment(IntegrationTestsEnvironment):
    def __init__(self, cmd):
        super().__init__(cmd)
        self.project = cmd.project
        # gobin_dir = os.environ["GOBIN"]
        self.hardhat = Hardhat(cmd)

    def signer_array_to_ethereum_accounts(self, accounts, n_validators):
        assert len(accounts) >= n_validators + 3
        operator, owner, pauser, *rest = accounts
        validators, available = rest[:n_validators], rest[n_validators:]
        return {
            "proxy_admin": operator,
            "operator": operator,
            "owner": owner,
            "pauser": pauser,
            "validators": validators,
            "available": available,
        }

    def run_ebrelayer_peggy(self, tcp_url, websocket_address, bridge_registry_sc_addr, validator_moniker,
        validator_mnemonic, chain_id, symbol_translator_file, relayerdb_path, ethereum_address, ethereum_private_key,
        log_file=None
    ):
        return Ebrelayer(self.cmd).init(tcp_url, websocket_address, bridge_registry_sc_addr, validator_moniker,
            validator_mnemonic, chain_id, ethereum_private_key=ethereum_private_key, ethereum_address=ethereum_address,
            node=tcp_url, keyring_backend="test", sign_with=validator_moniker,
            symbol_translator_file=symbol_translator_file, relayerdb_path=relayerdb_path, log_file=log_file)

    # Override
    def run(self):
        # self.project._make_go_binaries()

        log_dir = "/tmp/sifnode"
        self.cmd.mkdir(log_dir)
        hardhat_log_file = open(os.path.join(log_dir, "ganache.log"), "w")  # TODO close + use a different name
        sifnoded_log_file = open(os.path.join(log_dir, "sifnoded.log"), "w")  # TODO close
        ebrelayer_log_file = open(os.path.join(log_dir, "evmrelayer.log"), "w")  # TODO close
        witness_log_file = open(os.path.join(log_dir, "witness.log"), "w")  # TODO close; will be empty on non-peggy2 branch

        self.cmd.rmdir(self.cmd.get_user_home(".sifnoded"))  # Purge test keyring backend

        hardhat_hostname = "localhost"
        hardhat_port = 8545
        hardhat_proc = self.hardhat.start(hardhat_hostname, hardhat_port, log_file=hardhat_log_file)

        hardhat_validator_count = 1
        hardhat_network_id = 1  # Not used in smart-contracts/src/devenv/hardhatNode.ts
        hardhat_chain_id = 1  # Not used in smart-contracts/src/devenv/hardhatNode.ts
        # This value is actually returned from HardhatNodeRunner. It comes from smart-contracts/hardhat.config.ts.
        # In Typescript, its value is obtained by 'require("hardhat").hre.network.config.chainId'.
        # See https://hardhat.org/advanced/hardhat-runtime-environment.html
        hardhat_chain_id = 31337  # From smart-contracts/hardhat.config.ts, a dynamically set in Typescript by 'requre("hardhat")'
        hardhat_accounts = self.signer_array_to_ethereum_accounts(Hardhat.default_accounts(), hardhat_validator_count)

        self.hardhat.compile_smart_contracts()
        peggy_sc_addrs = self.hardhat.deploy_smart_contracts()

        self.write_compatibility_json_file_with_smart_contract_addresses({
            "BridgeRegistry": peggy_sc_addrs.bridge_registry,
            "BridgeBank": peggy_sc_addrs.bridge_bank,
            # TODO There is no BridgeToken smart contract on Peggy2.0 branch, but there are "cosmos_bridge" and "rowan"
        })

        chain_id = "localnet"
        sifnoded_network_dir = "/tmp/sifnodedNetwork"
        self.cmd.rmdir(sifnoded_network_dir)
        self.cmd.mkdir(sifnoded_network_dir)
        network_config_file = "/tmp/sifnodedConfig.yml"
        validator_count = 1
        seed_ip_address = "10.10.1.1"
        self.cmd.sifgen_create_network(chain_id, validator_count, sifnoded_network_dir, network_config_file, seed_ip_address)
        netdef_yml = yaml_load(self.cmd.read_text_file(network_config_file))

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
        assert len(netdef_yml) == validator_count
        netdef = exactly_one(netdef_yml)
        validator_moniker = netdef["moniker"]
        validator_mnemonic = netdef["mnemonic"].split(" ")
        # Not used
        # validator_password = netdef["password"]

        chain_dir = os.path.join(sifnoded_network_dir, "validators", chain_id, validator_moniker)
        sifnoded_home = os.path.join(chain_dir, ".sifnoded")
        denom_whitelist_file = project_dir("test", "integration", "whitelisted-denoms.json")

        self.cmd.sifchain_init_peggy(validator_moniker, validator_mnemonic, sifnoded_home, denom_whitelist_file)

        tendermint_port = 26657
        tcp_url = "tcp://{}:{}".format("0.0.0.0", tendermint_port)
        sifnoded_proc = self.cmd.sifnoded_start(minimum_gas_prices=[0.5, "rowan"], tcp_url=tcp_url,
            sifnoded_home=sifnoded_home, log_format_json=True, log_file=sifnoded_log_file)

        def _wait_for_sif_validator_up():
            # TODO Deduplicate: this is also in run_ebrelayer()
            # netdef_json is path to file containing json_dump(netdef)
            # while not self.tcp_probe_connect("localhost", tendermint_port):
            #     time.sleep(1)
            # self.wait_for_sif_account(netdef_json, validator1_address)
            pass
        _wait_for_sif_validator_up()

        relayerdb_path = self.cmd.mktempdir()
        web3_provider = "ws://{}:{}/".format(hardhat_hostname, str(hardhat_port))
        ethereum_address, ethereum_private_key = hardhat_accounts["validators"][0]
        symbol_translator_file = os.path.join(self.test_integration_dir, "config", "symbol_translator.json")
        ebrelayer = Ebrelayer(self.cmd)

        # validator_moniker and validator_mnemonic is the validator on sifchain side
        # ethereum_address and ethereum_private_key is the validator on ethereum side

        if on_peggy2_branch:
            hardcoded_network_descriptor = "1"  # TODO Don't hardcode me
            ebrelayer_proc = ebrelayer.peggy2_init_relayer(
                hardcoded_network_descriptor,
                tcp_url,
                web3_provider,
                peggy_sc_addrs.bridge_registry,
                validator_moniker,
                validator_mnemonic,
                chain_id,
                symbol_translator_file,
                ethereum_address,
                ethereum_private_key,
                keyring_backend="test",
                log_file=ebrelayer_log_file,
                cwd=None,
            )
            log.debug("Started ebrelayer: pid={}".format(ebrelayer_proc.pid))
            witness_proc = ebrelayer.peggy2_init_witness(
                hardcoded_network_descriptor,
                tcp_url,
                web3_provider,
                peggy_sc_addrs.bridge_registry,
                validator_moniker,
                validator_mnemonic,
                chain_id,
                symbol_translator_file,
                ethereum_address,
                ethereum_private_key,
                relayerdb_path=relayerdb_path,
                keyring_backend="test",
                log_file=witness_log_file,
                cwd=None,
            )
            log.debug("Started witness: pid={}".format(witness_proc.pid))

            env_vars = {
                # Computed
                "BASEDIR": self.project.base_dir,
                "GOBIN": self.project.go_bin_dir,
                "CHAINDIR": chain_dir,

                # Ethereum
                "ETHEREUM_ADDRESS": hardhat_accounts["available"][0][0],
                "ETHEREUM_PRIVATE_KEY": "0x" + hardhat_accounts["available"][0][1],
                "ROWAN_SOURCE": "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
                "ETH_ACCOUNT_OPERATOR_ADDRESS": hardhat_accounts["operator"][0],
                "ETH_ACCOUNT_OPERATOR_PRIVATEKEY": "0x" + hardhat_accounts["operator"][1],
                "ETH_ACCOUNT_OWNER_ADDRESS": hardhat_accounts["owner"][0],
                "ETH_ACCOUNT_OWNER_PRIVATEKEY": "0x" + hardhat_accounts["owner"][1],
                "ETH_ACCOUNT_PAUSER_ADDRESS": hardhat_accounts["pauser"][0],
                "ETH_ACCOUNT_PAUSER_PRIVATEKEY": "0x" + hardhat_accounts["pauser"][1],
                "ETH_ACCOUNT_PROXYADMIN_ADDRESS": hardhat_accounts["proxy_admin"][0],
                "ETH_ACCOUNT_PROXYADMIN_PRIVATEKEY": "0x" + hardhat_accounts["proxy_admin"][1],
                "ETH_ACCOUNT_VALIDATOR_ADDRESS": "0x90f79bf6eb2c4f870365e785982e1f101e93b906",
                "ETH_ACCOUNT_VALIDATOR_PRIVATEKEY": "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
                "ETH_CHAIN_ID": hardhat_chain_id,
                "ETH_HOST": hardhat_hostname,
                "ETH_PORT": hardhat_port,

                # Smart contracts
                "BRIDGE_BANK_ADDRESS": peggy_sc_addrs.bridge_bank,
                "BRIDGE_REGISTERY_ADDRESS": peggy_sc_addrs.bridge_registry,
                "COSMOS_BRIDGE_ADDRESS": peggy_sc_addrs.cosmos_bridge,
                "ROWANTOKEN_ADDRESS": peggy_sc_addrs.rowan,
                "BRIDGE_TOKEN_ADDRESS": peggy_sc_addrs.rowan,

                # Sifnode
                "TCP_URL": tcp_url,
                "VALIDATOR_ADDRESS": netdef["validator_address"],
                "VALIDATOR_CONSENSUS_ADDRESS": netdef["validator_consensus_address"],
                "VALIDATOR_MENOMONIC": " ".join(validator_mnemonic),  # == netdef["validator_mnemonic"]
                "VALIDATOR_MONIKER": validator_moniker,  # == netdef["validator_moniker"]
                "VALIDATOR_PASSWORD": netdef["password"],
                "VALIDATOR_PUB_KEY": netdef["pub_key"],
            }

            log.debug("env_vars: " + repr(env_vars))
            dotenv_file = os.path.join(self.project.smart_contracts_dir, ".env")
            env_json_file = os.path.join(self.project.smart_contracts_dir, "env.json")
            environment_json_file = os.path.join(self.project.smart_contracts_dir, "environment.json")
            self.cmd.write_text_file(dotenv_file, joinlines([
                "export {}=\"{}\"".format(k, v) for k, v in env_vars.items()]))
            self.cmd.write_text_file(env_json_file, json.dumps(env_vars, indent=4))
            self.cmd.write_text_file(environment_json_file, json.dumps(env_vars, indent=4))
        else:
            ebrelayer_proc = ebrelayer.init(
                tcp_url,
                web3_provider,
                peggy_sc_addrs.bridge_registry,
                validator_moniker,
                validator_mnemonic,
                chain_id,
                ethereum_private_key=ethereum_private_key,
                ethereum_address=ethereum_address,
                node=tcp_url,
                keyring_backend="test",
                sign_with=validator_moniker,
                symbol_translator_file=symbol_translator_file,
                relayerdb_path=relayerdb_path,
                log_file=ebrelayer_log_file,
            )
            witness_proc = None

            state_vars = {
                "BASEDIR": self.project.base_dir,
                "TEST_INTEGRATION_DIR": project_dir("test/integration"),
                "ETHEREUM_WEBSOCKET_ADDRESS": web3_provider,
                "TEST_INTEGRATION_PY_DIR": project_dir("test/integration/src/py"),
            }
            vagrantenv_path = project_dir("test/integration/vagrantenv.sh")
            self.cmd.write_text_file(vagrantenv_path, joinlines(format_as_shell_env_vars(state_vars)))
            self.cmd.write_text_file(project_dir("test/integration/vagrantenv.json"), json.dumps(state_vars, indent=4))

        return hardhat_proc, sifnoded_proc, ebrelayer_proc, witness_proc

    # Write compatibility JSON files with smart contract addresses so that test_utilities.py:contract_artifact() keeps
    # working. We're not using truffle, so we need to create files with the same names and structure as it's used for
    # integration tests: .["networks"]["5777"]["address"]
    def write_compatibility_json_file_with_smart_contract_addresses(self, smart_contract_addresses):
        integration_tests_expected_network_id = 5777
        d = project_dir("smart-contracts", "build", "contracts")
        self.cmd.mkdir(d)
        for sc_name, sc_addr in smart_contract_addresses.items():
            self.cmd.write_text_file(os.path.join(d, f"{sc_name}.json"), json.dumps({
                "networks": {str(integration_tests_expected_network_id): {"address": sc_addr}}}, indent=4))


def main(argv):
    # tmux usage:
    # tmux new-session -d -s env1
    # tmux main-pane-height -t env1 10
    # tmux split-window -h -t env1
    # tmux split-window -h -t env1
    # tmux select-layout -t env1 even-vertical
    # OR: tmux select-layout main-horizontal
    logging.basicConfig(stream=sys.stdout, level=logging.DEBUG, format="%(message)s")
    what = argv[0] if argv else None
    cmd = Integrator()
    project = cmd.project
    if what == "project-init":
        project.init()
    elif what == "project-clean":
        project.cleanup_and_reset_state()
    elif what == "project-fullclean":
        project.fullclean()
    elif what == "run-ui-env":
        e = UIStackEnvironment(cmd)
        e.stack_save_snapshot()
        e.stack_push()
    elif what == "run-integration-env":
        env = IntegrationTestsEnvironment(cmd)
        project.cleanup_and_reset_state()
        # deploy/networks already included in run()
        processes = env.run()
        input("Press ENTER to exit...")
        killall(processes)
        # TODO Cleanup:
        # - rm -rf test/integration/sifnoderelayerdb
        # - rm -rf networks/validators/localnet/$moniker/.sifnoded
        # - If you ran the execute_integration_test_*.sh you need to kill ganache-cli for proper cleanup
        #   as it might have been killed and started outside of our control
    elif what == "create_snapshot":
        snapshot_name = argv[1]
        project.cleanup_and_reset_state()
        env = IntegrationTestsEnvironment(cmd)
        processes = env.run()
        # Give processes some time to settle, for example relayerdb must init and create its "relayerdb"
        time.sleep(45)
        killall(processes)
        # processes1 = e.restart_processes()
        env.create_snapshot(snapshot_name)
    elif what == "restore_snapshot":
        snapshot_name = argv[1]
        env = IntegrationTestsEnvironment(cmd)
        env.restore_snapshot(snapshot_name)
        processes = env.restart_processes()
        input("Press ENTER to exit...")
        killall(processes)
    elif what == "run-peggy-env":
        # Equivalent to future/devenv - hardhat, sifnoded, ebrelayer
        # I.e. cd smart-contracts; GOBIN=/home/anderson/go/bin npx hardhat run scripts/devenv.ts
        env = PeggyEnvironment(cmd)
        processes = env.run()
        input("Press ENTER to exit...")
        killall(processes)
    elif what == "run-integration-tests":
        # TODO After switching the branch,: cd smart-contracts; rm -rf node_modules; + cmd.install_smart_contract_dependencies() (yarn clean + yarn install)
        scripts = [
            "execute_integration_tests_against_test_chain_peg.sh",
            "execute_integration_tests_against_test_chain_clp.sh",
            "execute_integration_tests_against_any_chain.sh",
            "execute_integration_tests_with_snapshots.sh",
        ]
        for script in scripts:
            force_kill_processes(cmd)
            e = IntegrationTestsEnvironment(cmd)
            processes = e.run()
            cmd.execst(script, cwd=project_dir("test", "integration"))
            killall(processes)
            force_kill_processes(cmd)  # Some processes are restarted during integration tests so we don't own them
        log.info("Everything OK")
    elif what == "test-logging":
        ls_cmd = mkcmd(["ls", "-al", "."], cwd="/tmp")
        res = stdout_lines(cmd.execst(**ls_cmd))
        print(ls_cmd)
    else:
        raise Exception("Missing/unknown command")


if __name__ == "__main__":
    main(sys.argv[1:])
