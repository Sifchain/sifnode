import json
import logging
import os
import re
import shutil
import subprocess
import sys
import time
import urllib.request
import yaml  # pip install pyyaml


log = logging.getLogger(__name__)

def stdout_lines(res):
    return res[1].splitlines()

def exactly_one(items):
    if len(items) == 0:
        raise ValueError("Zero items")
    elif len(items) > 1:
        raise ValueError("Multiple items")
    else:
        return items[0]

def project_dir(*paths):
    return os.path.abspath(os.path.join(os.path.normpath(os.path.join(__file__, *([os.path.pardir]*3))), *paths))

def yaml_load(s):
    return yaml.load(s, Loader=yaml.SafeLoader)

def sif_format_amount(amount, denom):
    return "{}{}".format(amount, denom)

def http_get(url):
    with urllib.request.urlopen(url) as r:
        return r.read()

def popen(args, env=None):
    if env:
        env = dict_merge(os.environ, env)
    return subprocess.Popen(args, env=env)

def dict_merge(*dicts):
    result = {}
    for d in dicts:
        for k, v in d.items():
            result[k] = v
    return result

NULL_ADDRESS = "0x0000000000000000000000000000000000000000"


class Command:
    def execst(self, args, cwd=None, env=None, stdin=None, binary=False, pipe=True, check_exit=True):
        if stdin is not None:
            if type(stdin) == list:
                stdin = "".join([line + "\n" for line in stdin])
        p = subprocess.PIPE if pipe else None
        if env:
            env = dict_merge(os.environ, env)
        popen = subprocess.Popen(args, cwd=cwd, env=env, stdin=subprocess.PIPE, stdout=p, stderr=p, text=not binary)
        stdout_data, stderr_data = popen.communicate(input=stdin)
        if check_exit and (popen.returncode != 0):
            raise Exception("Command '{}' exited with returncode {}: {}".format(" ".join(args), popen.returncode, repr(stderr_data)))
        return popen.returncode, stdout_data, stderr_data

    def rm(self, path):
        os.remove(path)

    def read_text_file(self, path):
        with open(path, "rt") as f:
            return f.read()  # TODO Convert to exec

    def write_text_file(self, path, s):
        with open(path, "wt") as f:
            f.write(s)

    def mkdir(self, path):
        os.makedirs(path, exist_ok=True)

    def rmdir(self, path):
        if os.path.exists(path):
            shutil.rmtree(path)  # TODO Convert to exec

    def copy_file(self, src, dst):
        shutil.copy(src, dst)

    def exists(self, path):
        return os.path.exists(path)

    def get_user_home(self, *paths):
        return os.path.join(os.environ["HOME"], *paths)

    def mktempdir(self):
        return exactly_one(stdout_lines(self.execst(["mktemp", "-d"])))

    def mktempfile(self):
        return exactly_one(stdout_lines(self.execst(["mktemp"])))

    def tar_create(self, path, tar_path, compression=None):
        comp_opts = {"gz": "z"}
        comp = comp_opts[compression] if compression in comp_opts else ""
        # tar on 9p filesystem reports "file shrank by ... bytes" and exits with errorcode 1
        tar_quirks = True
        if tar_quirks:
            tmpdir = self.mktempdir()
            try:
                shutil.copytree(path, tmpdir, dirs_exist_ok=True)
                self.execst(["tar", "cf" + comp, tar_path, "."], cwd=tmpdir)
            finally:
                self.rmdir(tmpdir)
        else:
            self.execst(["tar", "cf" + comp, tar_path, "."], cwd=path)


class Ganache(Command):
    def start_ganache_cli(self, mnemonic=None, db=None, port=None, host=None, network_id=None, gas_price=None,
        gas_limit=None, block_time=None, account_keys_path=None, popen_args=None):
        args = ["ganache-cli"] + \
            (["--mnemonic", " ".join(mnemonic)] if mnemonic else []) + \
            (["--db", db] if db else []) + \
            (["--port", str(port)] if port is not None else []) + \
            (["--host", host] if host else []) + \
            (["--networkId", str(network_id)] if network_id is not None else []) + \
            (["--gasPrice", str(gas_price)] if gas_price is not None else []) + \
            (["--gasLimit", str(gas_limit)] if gas_limit is not None else []) + \
            (["--blockTime", str(block_time)] if block_time is not None else []) + \
            (["--account_keys_path", account_keys_path] if account_keys_path is not None else [])
        return popen(args, **(popen_args if popen_args is not None else dict()))


class Sifnoded(Command):
    def sifnoded_init(self, moniker, chain_id):
        args = ["sifnoded", "init", moniker, "--chain-id={}".format(chain_id)]
        res = self.execst(args)
        return json.loads(res[2])  # output is on stderr

    def sifnoded_generate_deterministic_account(self, name, mnemonic):
        args = ["sifnoded", "keys", "add", name, "--keyring-backend={}".format("test"), "--recover"]
        stdin = [" ".join(mnemonic)]
        res = self.execst(args, stdin=stdin)
        return yaml_load(res[1])[0]

    def sifnoded_keys_show(self, name, bech=None, keyring_backend=None):
        keyring_backend = "test"
        args = ["sifnoded", "keys", "show", name] + \
               (["--bech", bech] if bech else []) + \
               (["--keyring-backend={}".format(keyring_backend)] if keyring_backend else [])
        res = self.execst(args)
        return yaml_load(res[1])

    def sifnoded_keys_add(self, args, stdin=None):
        return yaml_load(self.execst(["sifnoded", "keys", "add"] + args, stdin=stdin)[1])

    def sifnoded_add_genesis_account(self, address, tokens):
        tokens_str = ",".join([sif_format_amount(amount, denom) for amount, denom in tokens])
        args = ["sifnoded", "add-genesis-account", address, tokens_str]
        self.execst(args, pipe=False)

    def sifnoded_add_genesis_validators(self, address):
        args = ["sifnoded", "add-genesis-validators", address]
        res = self.execst(args)
        return res

    def sifnoded_tx_clp_create_pool(self, chain_id, keyring_backend, from_name, symbol, fees, native_amount, external_amount):
        args = ["sifnoded", "tx", "clp", "create-pool", "--chain-id={}".format(chain_id),
            "--keyring-backend={}".format(keyring_backend), "--from", from_name, "--symbol", symbol, "--fees",
            sif_format_amount(*fees), "--nativeAmount", str(native_amount), "--externalAmount", str(external_amount),
            "--yes"]
        res = self.execst(args)
        return yaml_load(res[1])

    def sifnoded_launch(self, minimum_gas_prices=None):
        args = ["sifnoded", "start"] + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else [])
        return popen(args)

    def sifnoded_get_status(self, host, port):
        url = "http://{}:{}/node_info".format(host, port)
        return json.loads(http_get(url).decode("UTF-8"))

    def tcp_probe_connect(self, host, port):
        res = self.execst(["nc", "-z", host, str(port)], check_exit=False)
        return res[0] == 0

    def wait_for_file(self, path):
        while not self.exists(path):
            time.sleep(1)

class Integrator(Ganache, Sifnoded, Command):
    def __init__(self):
        self.smart_contracts_dir = project_dir("smart-contracts")

    def ebrelayer_init(self, ethereum_private_key, tendermind_node, web3_provider, bridge_registry_contract_address,
        validator_moniker, validator_mnemonic, chain_id, gas=None, gas_prices=None, node=None, keyring_backend=None,
        sign_with=None):
        env = {"ETHEREUM_PRIVATE_KEY": ethereum_private_key}
        args = ["ebrelayer", "init", tendermind_node, web3_provider, bridge_registry_contract_address,
            validator_moniker, " ".join(validator_mnemonic), "--chain-id={}".format(chain_id)] + \
            (["--gas", str(gas)] if gas is not None else []) + \
            (["--gas-prices", sif_format_amount(*gas_prices)] if gas_prices is not None else []) + \
            (["--node", node] if node is not None else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend is not None else []) + \
            (["--from", sign_with] if sign_with is not None else [])
        return popen(args, env=env)

    def sif_wait_up(self, host, port):
        while True:
            from urllib.error import URLError
            try:
                return self.sifnoded_get_status(host, port)
            except URLError:
                time.sleep(1)

    def yarn(self, args, cwd=None, env=None):
        return self.execst(["yarn"] + args, cwd=cwd, env=env, pipe=False)

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

    def build_smart_contracts_for_integration_tests(self):
        self.execst(["make", "clean-smartcontracts"], cwd=self.smart_contracts_dir)
        self.yarn(["install"], cwd=self.smart_contracts_dir)

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
            env["INITIAL_VALIDATOR_POWERS"] = ",".join(initial_validator_powers)
        if pauser is not None:
            env["PAUSER"] = pauser
        if mainnet_gas_price is not None:
            env["MAINNET_GAS_PRICE"] = mainnet_gas_price

        env_path = os.path.join(self.smart_contracts_dir, ".env")
        if env_file is not None:
            self.copy_file(env_file, env_path)

        self._check_env_vs_file(env, env_path)

        # TODO ui scripts use just "yarn; yarn migrate" alias "npx truffle migrate --reset",
        self.execst(["npx", "truffle", "deploy", "--network", network_name, "--reset"], env=env,
            cwd=self.smart_contracts_dir, pipe=False)

    def deploy_smart_contracts_for_ui_stack(self):
        self.copy_file(os.path.join(self.smart_contracts_dir, ".env.ui.example"), os.path.join(self.smart_contracts_dir, ".env"))
        # TODO Might not be neccessary
        self.yarn([], cwd=self.smart_contracts_dir)
        self.yarn(["migrate"], cwd=self.smart_contracts_dir)

    def get_smart_contract_address(self, compiled_json_path, network_id):
        return json.loads(self.read_text_file(compiled_json_path))["networks"][str(network_id)]["address"]

    def get_bridge_smart_contract_addresses(self, network_id):
        return [self.get_smart_contract_address(os.path.join(
            self.smart_contracts_dir, f"build/contracts/{x}.json"), network_id)
            for x in ["BridgeToken", "BridgeRegistry", "BridgeBank"]]

    def truffle_exec(self, script_name, *script_args, env=None):
        self._check_env_vs_file(env, os.path.join(self.smart_contracts_dir, ".env"))
        script_path = os.path.join(self.smart_contracts_dir, f"scripts/{script_name}.js")
        # Hint: call web3 directly, avoid npx + truffle + script
        # Maybe: self.cmd.yarn(["integrationtest:setTokenLockBurnLimit", str(amount)])
        self.execst(["npx", "truffle", "exec", script_path] + list(script_args), env=env, cwd=self.smart_contracts_dir)

    def set_token_lock_burn_limit(self, update_address, amount, ethereum_private_key, infura_project_id, local_provider):
        env = {
            "ETHEREUM_PRIVATE_KEY": ethereum_private_key,
            "UPDATE_ADDRESS": update_address,
            "INFURA_PROJECT_ID": infura_project_id,
            "LOCAL_PROVIDER": local_provider,
        }
        # Needs: ETHEREUM_PRIVATE_KEY, INFURA_PROJECT_ID, LOCAL_PROVIDER, UPDATE_ADDRESS
        self.truffle_exec("setTokenLockBurnLimit", str(amount), env=env)


class UIStackPlaybook:
    def __init__(self, cmd):
        self.cmd = cmd
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
        self.akasha_name = "akasha"
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
        self.cmd.execst(["make", "install"], cwd=project_dir(), pipe=False)

        # ui/scripts/stack-launch.sh -> ui/scripts/_eth.sh -> ui/chains/etc/launch.sh
        self.cmd.rmdir(self.ganache_db_path)
        self.cmd.yarn([], cwd=project_dir("ui/chains/eth"))  # Installs ui/chains/eth/node_modules
        # Note that this runs ganache-cli from $PATH whereas scripts start it with yarn in ui/chains/eth
        ganache_proc = self.cmd.start_ganache_cli(mnemonic=self.ethereum_root_mnemonic, db=self.ganache_db_path,
            port=7545, network_id=self.network_id, gas_price=20000000000, gas_limit=6721975, host="0.0.0.0")

        # ui/scripts/stack-launch.sh -> ui/scripts/_sif.sh -> ui/chains/sif/launch.sh
        self.cmd.sifnoded_init("test", self.chain_id)
        self.cmd.copy_file(project_dir("ui/chains/sif/app.toml"), os.path.join(self.sifnoded_path, "config/app.toml"))
        log.info(f"Generating deterministic account - {self.shadowfiend_name}...")
        shadowfiend_account = self.cmd.sifnoded_generate_deterministic_account(self.shadowfiend_name, self.shadowfiend_mnemonic)
        log.info(f"Generating deterministic account - {self.akasha_name}...")
        akasha_account = self.cmd.sifnoded_generate_deterministic_account(self.akasha_name, self.akasha_mnemonic)
        log.info(f"Generating deterministic account - {self.juniper_name}...")
        juniper_account = self.cmd.sifnoded_generate_deterministic_account(self.juniper_name, self.juniper_mnemonic)
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
        sifnoded_proc = self.cmd.sifnoded_launch(minimum_gas_prices=[0.5, "rowan"])

        # sifnoded must be up before continuing
        self.cmd.sif_wait_up("localhost", 1317)

        # ui/scripts/_migrate.sh -> ui/chains/peggy/migrate.sh
        self.cmd.deploy_smart_contracts_for_ui_stack()

        # ui/scripts/_migrate.sh -> ui/chains/eth/migrate.sh
        # send through atk and btk tokens to eth chain
        self.cmd.yarn(["migrate"], cwd=project_dir("ui/chains/eth"))

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
            self.cmd.yarn(["peggy:whiteList", addr, "true"], cwd=self.cmd.smart_contracts_dir)

        # ui/scripts/stack-launch.sh -> ui/scripts/_peggy.sh -> ui/chains/peggy/launch.sh
        # rm -rf ui/chains/peggy/relayerdb
        # ebrelayer is in $GOBIN, gets installed by "make install"
        ethereum_private_key = smart_contracts_env_ui_example_vars["ETHEREUM_PRIVATE_KEY"]
        ebrelayer_proc = self.cmd.ebrelayer_init(ethereum_private_key, "tcp://localhost:26657", "ws://localhost:7545/",
            bridge_registry_address, self.shadowfiend_name, self.shadowfiend_mnemonic, self.chain_id, gas=5*10**12,
            gas_prices=[0.5, "rowan"])

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
        self.cmd.tar_create(project_dir("ui/chains/peggy/relayerdb"), os.path.join(snapshots_dir, "peggy.tar.gz"), compression="gz")
        # mkdir -p smart-contracts/build
        self.cmd.tar_create(project_dir("smart-contracts/build"), os.path.join(snapshots_dir, "peggy_build.tar.gz"), compression="gz")

        # ui/chains/sif/snapshot.sh:
        self.cmd.tar_create(self.sifnoded_path, os.path.join(snapshots_dir, "sif.tar.gz"), compression="gz")

        # ui/chains/etc/snapshot.sh:
        self.cmd.tar_create(self.ganache_db_path, os.path.join(snapshots_dir, "eth.tar.gz"), compression="gz")

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


class IntegrationTestsPlaybook:
    def __init__(self, cmd):
        self.cmd = cmd
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

    def run(self):
        test_integration_dir = project_dir("test/integration")
        data_dir = project_dir("test/integration/vagrant/data")

        chainnet = "localnet"
        tcp_url = "tcp://0.0.0.0:26657"
        ethereum_websocket_address = "ws://localhost:7545/"

        # make go binaries (a lot of nonsense!)
        self.cmd.execst(["make"], cwd=test_integration_dir, env={"BASEDIR": project_dir()})

        self.cmd.build_smart_contracts_for_integration_tests()

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
        validator_mnemonic = ["candy", "maple", "cake", "sugar", "pudding", "cream", "honey", "rich", "smooth", "crumble", "sweet", "treat"]
        block_time = None  # TODO
        account_keys_path = os.path.join(data_dir, "ganachekeys.json")
        ganache_db_path = self.cmd.mktempdir()
        ganache_proc = self.cmd.start_ganache_cli(block_time=block_time, host="0.0.0.0", mnemonic=validator_mnemonic,
            network_id=self.network_id, port=7545, db=ganache_db_path, account_keys_path=account_keys_path)

        self.cmd.wait_for_file(account_keys_path)
        time.sleep(2)

        ganache_keys = json.loads(self.cmd.read_text_file(account_keys_path))
        ebrelayer_ethereum_addr = list(ganache_keys["private_keys"].keys())[9]
        ebrelayer_ethereum_private_key = ganache_keys["private_keys"][ebrelayer_ethereum_addr]
        # TODO Check for possible non-determinism of dict().keys() ordering (c.f. test/integration/vagrantenv.sh)
        # TODO ebrelayer_ethereum_private_key is NOT the same as in test/integration/.env.ciExample
        assert ebrelayer_ethereum_addr == "0x5aeda56215b167893e80b4fe645ba6d5bab767de"
        assert ebrelayer_ethereum_private_key == "8d5366123cb560bb606379f90a0bfd4769eecc0557f1b362dcae9012b548b1e5"

        env_file = project_dir("test/integration/.env.ciExample")
        self.cmd.deploy_smart_contracts_for_integration_tests(self.network_name, owner=self.owner, pauser=self.pauser,
            initial_validator_addresses=[ebrelayer_ethereum_addr], env_file=env_file)

        bridge_token_sc_addr, bridge_registry_sc_addr, bridge_bank_sc_addr = \
            self.cmd.get_bridge_smart_contract_addresses(self.network_id)

        # TODO This should be last (after return from setup_sifchain.sh)
        burn_limits = [
            [NULL_ADDRESS, 31*10**18],
            [bridge_token_sc_addr, 10**25],
        ]
        env_file_vars = self.cmd.primitive_parse_env_file(env_file)
        for address, amount in burn_limits:
            self.cmd.set_token_lock_burn_limit(
                address,
                amount,
                env_file_vars["ETHEREUM_PRIVATE_KEY"],  # != ebrelayer_ethereum_private_key
                env_file_vars["INFURA_PROJECT_ID"],
                env_file_vars["LOCAL_PROVIDER"]
            )

        # test/integration/setup_sifchain.sh:
        networks_dir = project_dir("deploy/networks")
        self.cmd.rmdir(networks_dir)  # networks_dir has many directories without write permission, so change those before deleting it
        self.cmd.mkdir(networks_dir)
        self.cmd.execst(["rake", f"genesis:network:scaffold[{chainnet}]"], env={"BASEDIR": project_dir()}, pipe=False)
        netdef = exactly_one(yaml_load(self.cmd.read_text_file(project_dir(networks_dir, "network-definition.yml"))))
        netdef_json = os.path.join(data_dir, "netdef.json")
        self.cmd.write_text_file(netdef_json, json.dumps(netdef))

        validator_moniker = netdef["moniker"]
        validator1_address = netdef["address"]
        validator1_password = netdef["password"]
        validator_mnemonic = netdef["mnemonic"].split(" ")
        chaindir = os.path.join(networks_dir, f"validators/{chainnet}/{validator_moniker}")
        # SIFNODED_LOG=$datadir/logs/sifnoded.log

        # test/integration/sifchain_start_daemon.sh:
        sifchaind_home = os.path.join(chaindir, ".sifnoded")
        whitelisted_validator = exactly_one(stdout_lines(self.cmd.execst(["sifnoded", "keys", "show",
            "--keyring-backend", "file", "-a", "--bech", "val", validator_moniker, "--home", sifchaind_home],
            stdin=[validator1_password])))
        log.info(f"Whitelisted validator: {whitelisted_validator}")
        self.cmd.execst(["sifnoded", "add-genesis-validators", whitelisted_validator, "--home", sifchaind_home])
        adminuser_addr = json.loads(self.cmd.execst(["sifnoded", "keys", "add", "sifnodeadmin", "--keyring-backend",
            "test", "--output", "json"], stdin=["y"])[1])["address"]
        self.cmd.execst(["sifnoded", "add-genesis-account", adminuser_addr, sif_format_amount(10**20, "rowan"),
            "--home", sifchaind_home], pipe=False)
        self.cmd.execst(["sifnoded", "set-genesis-oracle-admin", adminuser_addr, "--home", sifchaind_home], pipe=False)
        sifnoded_proc = popen(["sifnoded", "start", "--minimum-gas-prices", sif_format_amount(0.5, "rowan"),
            "--rpc.laddr", tcp_url, "--home", sifchaind_home])

        rest_server_proc = popen(["sifnoded", "rest-server", "--laddr", "tcp://0.0.0.0:1317"])  # TODO cwd

        # test/integration/sifchain_start_ebrelayer.sh -> test/integration/sifchain_run_ebrelayer.sh
        while not self.cmd.tcp_probe_connect("localhost", 26657):
            time.sleep(1)
        self.cmd.execst(["python3", os.path.join(test_integration_dir, "src/py/wait_for_sif_account.py"),
            netdef_json, validator1_address], env={"USER1ADDR": "nothing"})
        time.sleep(10)
        # Error: The specified item could not be found in the keyring
        # res = self.cmd.execst(["sifnoded", "keys", "delete", moniker, "--keyring-backend", "test"], stdin=["y"])
        self.cmd.sifnoded_keys_add([validator_moniker, "--keyring-backend", "test", "--recover"],
            stdin=[" ".join(validator_mnemonic)])

        ebrelayer_proc = self.cmd.ebrelayer_init(ebrelayer_ethereum_private_key, tcp_url, ethereum_websocket_address,
            bridge_registry_sc_addr, validator_moniker, validator_mnemonic, chainnet, node=tcp_url,
            keyring_backend="test", sign_with=validator_moniker)

        return ganache_proc, sifnoded_proc, ebrelayer_proc



def main():
    logging.basicConfig(stream=sys.stdout, level=logging.DEBUG, format="%(message)s")
    cmd = Integrator()
    # Reset state:
    # pkill node; pkill ebrelayer; pkill sifnoded; rm -rvf $HOME/.sifnoded; rm -rvf ./vagrant/data; mkdir vagrant/data
    cmd.execst(["pkill", "node"], check_exit=False)
    cmd.execst(["pkill", "ebrelayer"], check_exit=False)
    cmd.execst(["pkill", "sifnoded"], check_exit=False)
    cmd.rmdir(cmd.get_user_home(".sifnoded"))
    cmd.rmdir(project_dir("test/integration/vagrant/data"))
    cmd.mkdir(project_dir("test/integration/vagrant/data"))
    time.sleep(3)
    ui_playbook = UIStackPlaybook(cmd)
    # ui_playbook.stack_save_snapshot()
    # ui_playbook.stack_push()
    it_playbook = IntegrationTestsPlaybook(cmd)
    it_playbook.run()


if __name__ == "__main__":
    main()
