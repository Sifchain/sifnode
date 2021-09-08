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

def stdout(res):
    return res[1]

def stdout_lines(res):
    return stdout(res).splitlines()

def joinlines(lines):
    return "".join([x + os.linesep for x in lines])

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

# Not used yet
def mkcmd(args, env=None, cwd=None, stdin=None):
    result = {"args": args}
    if env is not None:
        result["env"] = env
    if cwd is not None:
        result["cwd"] = cwd
    if stdin is not None:
        result["stdin"] = stdin
    return result

# stdin will always be redirected to the returned process' stdin.
# If pipe, the stdout and stderr will be redirected and available as stdout and stderr of the returned object.
# If not pipe, the stdout and stderr will not be redirected and will inherit sys.stdout and sys.stderr.
def _popen(args, cwd=None, env=None, binary=False, pipe=True):
    logging.debug(f"execst(): args={repr(args)}, cwd={repr(cwd)}")
    if env:
        env = dict_merge(os.environ, env)
    p = subprocess.PIPE if pipe else None
    return subprocess.Popen(args, cwd=cwd, env=env, stdin=subprocess.PIPE, stdout=p, stderr=p, text=not binary)

def popen(args, cwd=None, env=None):
    #     if env:
    #         env = dict_merge(os.environ, env)
    #     return subprocess.Popen(args, env=env, cwd=cwd)
    return _popen(args, cwd=cwd, env=env, pipe=False)

# popen_obj must be open in binary mode
def log_to_file(popen_obj, path):
    def copy_buf(in_buf, out_buf):
        while True:
            q = in_buf.read(4096)
            if not len(q):
                break
            out_buf.write(q)  # No locking on output file (GIL assumed)

    def run():
        with open(path, "wb") as f:
            t_stdout = Thread(target=copy_buf, args=(popen_obj.stdout, f))
            t_stderr = Thread(target=copy_buf, args=(popen_obj.stderr, f))
            t_stdout.start()
            t_stderr.start()
            t_stdout.join()
            t_stderr.join()

    from threading import Thread

    master_thread = Thread(target=run)
    master_thread.start()
    return master_thread

def dict_merge(*dicts):
    result = {}
    for d in dicts:
        for k, v in d.items():
            result[k] = v
    return result

def format_as_shell_env_vars(env, export=True):
    return ["{}{}=\"{}\"".format("export " if export else "", k, v) for k, v in env.items()]

NULL_ADDRESS = "0x0000000000000000000000000000000000000000"


class Command:
    def execst(self, args, cwd=None, env=None, stdin=None, binary=False, pipe=True, check_exit=True):
        proc = _popen(args, env=env, cwd=cwd, binary=binary, pipe=pipe)
        if stdin is not None:
            if type(stdin) == list:
                stdin = "".join([line + "\n" for line in stdin])
        stdout_data, stderr_data = proc.communicate(input=stdin)
        if check_exit and (proc.returncode != 0):
            raise Exception("Command '{}' exited with returncode {}: {}".format(" ".join(args), proc.returncode, repr(stderr_data)))
        return proc.returncode, stdout_data, stderr_data

    # Default implementation of popen for environemnts to start long-lived processes
    def popen(self, args, cwd=None, env=None):
        return popen(args, cwd=cwd, env=env)

    def rm(self, path):
        if os.path.exists(path):
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

    def __tar_compression_option(self, tarfile):
        filename = os.path.basename(tarfile).lower()
        if filename.endswith(".tar"):
            return ""
        elif filename.endswith(".tar.gz"):
            return "z"
        else:
            raise ValueError(f"Unknown extension for tar file: {tarfile}")

    def tar_create(self, path, tarfile):
        comp = self.__tar_compression_option(tarfile)
        # tar on 9p filesystem reports "file shrank by ... bytes" and exits with errorcode 1
        tar_quirks = True
        if tar_quirks:
            tmpdir = self.mktempdir()
            try:
                shutil.copytree(path, tmpdir, dirs_exist_ok=True)
                self.execst(["tar", "cf" + comp, tarfile, "."], cwd=tmpdir)
            finally:
                self.rmdir(tmpdir)
        else:
            self.execst(["tar", "cf" + comp, tarfile, "."], cwd=path)

    def tar_extract(self, tarfile, path):
        comp = self.__tar_compression_option(tarfile)
        if not self.exists(path):
            self.mkdir(path)
        self.execst(["tar", "xf" + comp, tarfile], cwd=path)


class Ganache(Command):
    def start_ganache_cli(self, mnemonic=None, db=None, port=None, host=None, network_id=None, gas_price=None,
        gas_limit=None, default_balance_ether=None, block_time=None, account_keys_path=None, popen_args=None):
        args = ["ganache-cli"] + \
            (["--mnemonic", " ".join(mnemonic)] if mnemonic else []) + \
            (["--db", db] if db else []) + \
            (["--port", str(port)] if port is not None else []) + \
            (["--host", host] if host else []) + \
            (["--networkId", str(network_id)] if network_id is not None else []) + \
            (["--gasPrice", str(gas_price)] if gas_price is not None else []) + \
            (["--gasLimit", str(gas_limit)] if gas_limit is not None else []) + \
            (["--defaultBalanceEther", str(default_balance_ether)] if default_balance_ether is not None else []) + \
            (["--blockTime", str(block_time)] if block_time is not None else []) + \
            (["--account_keys_path", account_keys_path] if account_keys_path is not None else [])
        return self.popen(args, **(popen_args if popen_args is not None else dict()))


class Sifnoded(Command):
    def sifnoded_init(self, moniker, chain_id):
        args = ["sifnoded", "init", moniker, "--chain-id={}".format(chain_id)]
        res = self.execst(args)
        return json.loads(res[2])  # output is on stderr

    def sifnoded_generate_deterministic_account(self, name, mnemonic):
        args = ["sifnoded", "keys", "add", name, "--keyring-backend={}".format("test"), "--recover"]
        stdin = [" ".join(mnemonic)]
        res = self.execst(args, stdin=stdin)
        return yaml_load(stdout(res))[0]

    def sifnoded_keys_show(self, name, bech=None, keyring_backend=None, home=None):
        keyring_backend = keyring_backend or "test"
        args = ["keys", "show", name] + \
            (["--bech", bech] if bech else [])
        res = self.sifnoded_exec(args, keyring_backend=keyring_backend, sifnoded_home=home)
        return yaml_load(stdout(res))

    def sifnoded_get_val_address(self, name):
        expected = exactly_one(stdout_lines(self.sifnoded_exec(["keys", "show", "-a", "--bech", "val", name], keyring_backend="test")))
        result = exactly_one(self.sifnoded_keys_show(name, bech="val", keyring_backend="test"))["address"]
        assert result == expected
        return result

    def sifnoded_keys_add(self, moniker, mnemonic):
        args = ["sifnoded", "keys", "add", moniker, "--keyring-backend", "test", "--recover"]
        stdin = [" ".join(mnemonic)]
        return yaml_load(stdout(self.execst(args, stdin=stdin)))

    # How "sifnoded keys add <name> --keyring-backend test" works:
    # If name does not exist yet, it creates it and returns a yaml
    # If name alredy exists, prompts for overwrite (y/n) on standard input, generates new address/pubkey/mnemonic
    # Directory used is xxx/keyring-test if "--home xxx" is specified, otherwise $HOME/.sifnoded/keyring-test

    def sifnoded_keys_add_1(self, name):
        res = self.sifnoded_exec(["keys", "add", name], keyring_backend="test", stdin=["y"])
        return exactly_one(yaml_load(stdout(res)))

    # From peggy
    # @TODO Passing mnemonic to stdin is useless, only "y/n" makes sense, probably could use sifnoded_keys_add_1
    # See smart-contracts/src/devenv/sifnoded.ts:addValidatorKeysToTestKeyring
    def sifnoded_keys_add_2(self, name, mnemonic):
        stdin = [" ".join(mnemonic)]
        res = self.sifnoded_exec(["keys", "add", name], keyring_backend="test", stdin=stdin)
        result = exactly_one(yaml_load(stdout(res)))
        # {"name": "<moniker>", "type": "local", "address": "sif1...", "pubkey": "sifpub1...", "mnemonic": "", "threshold": 0, "pubkeys": []}
        return result

    def sifnoded_keys_delete(self, name):
        self.execst(["sifnoded", "keys", "delete", name, "--keyring-backend", "test"], stdin=["y"], check_exit=False)

    def sifnoded_add_genesis_account(self, address, tokens, sifnoded_home=None):
        tokens_str = ",".join([sif_format_amount(amount, denom) for amount, denom in tokens])
        self.sifnoded_exec(["add-genesis-account", address, tokens_str], sifnoded_home=sifnoded_home)

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
        return yaml_load(stdout(res))

    def sifnoded_start(self, tcp_url=None, minimum_gas_prices=None, sifnoded_home=None):
        args = ["sifnoded", "start"] + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else []) + \
            (["--rpc.laddr", tcp_url] if tcp_url else []) + \
            (["--home", sifnoded_home] if sifnoded_home else [])
        return self.popen(args)

    def sifnoded_exec(self, args, sifnoded_home=None, keyring_backend=None, stdin=None, cwd=None):
        args = ["sifnoded"] + args + \
            (["--home", sifnoded_home] if sifnoded_home else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend else [])
        res = self.execst(args, stdin=stdin, cwd=cwd)
        return res

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

    def ebrelayer_init(self, tendermind_node, web3_provider, bridge_registry_contract_address,
        validator_moniker, validator_mnemonic, chain_id, ethereum_private_key=None, ethereum_address=None, gas=None,
        gas_prices=None, node=None, keyring_backend=None, sign_with=None, symbol_translator_file=None,
        relayerdb_path=None, cwd=None
    ):
        env = {}
        if ethereum_private_key:
            assert ethereum_private_key.startswith("0x")
            env["ETHEREUM_PRIVATE_KEY"] = ethereum_private_key[2:]
        if ethereum_address:
            assert ethereum_address.startswith("0x")
            env["ETHEREUM_ADDRESS"] = ethereum_address[2:]
        args = ["ebrelayer", "init", tendermind_node, web3_provider, bridge_registry_contract_address,
            validator_moniker, " ".join(validator_mnemonic), "--chain-id={}".format(chain_id)] + \
            (["--gas", str(gas)] if gas is not None else []) + \
            (["--gas-prices", sif_format_amount(*gas_prices)] if gas_prices is not None else []) + \
            (["--node", node] if node is not None else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend is not None else []) + \
            (["--from", sign_with] if sign_with is not None else []) + \
            (["--symbol-translator-file", symbol_translator_file] if symbol_translator_file else []) + \
            (["--relayerdb-path", relayerdb_path] if relayerdb_path else [])
        return self.popen(args, env=env, cwd=cwd)

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

    def install_smart_contracts_dependencies(self):
        self.execst(["make", "clean-smartcontracts"], cwd=self.smart_contracts_dir)  # = rm -rf build .openzeppelin
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
            env["INITIAL_VALIDATOR_POWERS"] = ",".join([str(x) for x in initial_validator_powers])
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
            cwd=self.smart_contracts_dir, pipe=False)  # TODO Why pipe=False?

    def deploy_smart_contracts_for_ui_stack(self):
        self.copy_file(os.path.join(self.smart_contracts_dir, ".env.ui.example"), os.path.join(self.smart_contracts_dir, ".env"))
        # TODO Might not be neccessary
        self.yarn([], cwd=self.smart_contracts_dir)
        self.yarn(["migrate"], cwd=self.smart_contracts_dir)

    # truffle
    def get_smart_contract_address(self, compiled_json_path, network_id):
        return json.loads(self.read_text_file(compiled_json_path))["networks"][str(network_id)]["address"]

    # truffle
    def get_bridge_smart_contract_addresses(self, network_id):
        return [self.get_smart_contract_address(os.path.join(
            self.smart_contracts_dir, f"build/contracts/{x}.json"), network_id)
            for x in ["BridgeToken", "BridgeRegistry", "BridgeBank"]]

    def npx(self, args, env=None, cwd=None):
        return self.execst(["npx"] + args, env=env, cwd=cwd)

    def truffle_exec(self, script_name, *script_args, env=None):
        self._check_env_vs_file(env, os.path.join(self.smart_contracts_dir, ".env"))
        script_path = os.path.join(self.smart_contracts_dir, f"scripts/{script_name}.js")
        # Hint: call web3 directly, avoid npx + truffle + script
        # Maybe: self.cmd.yarn(["integrationtest:setTokenLockBurnLimit", str(amount)])
        self.npx(["truffle", "exec", script_path] + list(script_args), env=env, cwd=self.smart_contracts_dir)

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

        adminuser_addr = self.sifnoded_keys_add_1("sifnodeadmin")["address"]
        self.sifchain_init_common(adminuser_addr, denom_whitelist_file, sifnoded_home)
        return adminuser_addr

    # @parameter validator_moniker - from network config
    # @parameter validator_mnemonic - from network config
    def sifchain_init_peggy(self, validator_moniker, validator_mnemonic, sifnoded_home, denom_whitelist_file):
        # Add validator key to test keyring
        _tmp0 = self.sifnoded_keys_add_2(validator_moniker, validator_mnemonic)
        valoper = self.sifnoded_get_val_address(validator_moniker)

        # (0, '', '2021/09/07 05:55:33 AddGenesisValidatorCmd, adding addr: sifvaloper1f5vj6j2mnkaw0yec3ut9at4rkl2u23k2fxtrsv to whitelist: []\n')
        self.sifnoded_exec(["add-genesis-validators", valoper], sifnoded_home=sifnoded_home)

        # Get whitelisted validator
        # TODO Value is not being used
        _whitelisted_validator = self.sifnoded_get_val_address(validator_moniker)
        assert valoper == _whitelisted_validator

        adminuser_addr = self.sifnoded_keys_add_1("sifnodeadmin")["address"]
        self.sifchain_init_common(adminuser_addr, denom_whitelist_file, sifnoded_home)
        return adminuser_addr

    def sifchain_init_common(self, adminuser_addr, denom_whitelist_file, sifnoded_home):
        tokens = [[10**20, "rowan"]]
        # Original from peggy:
        # self.cmd.execst(["sifnoded", "add-genesis-account", sifnoded_admin_address, "100000000000000000000rowan", "--home", sifnoded_home])
        self.sifnoded_add_genesis_account(adminuser_addr, tokens, sifnoded_home=sifnoded_home)
        self.sifnoded_exec(["set-genesis-oracle-admin", adminuser_addr], sifnoded_home=sifnoded_home)
        self.sifnoded_exec(["set-genesis-whitelister-admin", adminuser_addr], sifnoded_home=sifnoded_home)
        self.sifnoded_exec(["set-gen-denom-whitelist", denom_whitelist_file], sifnoded_home=sifnoded_home)

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
        sifnoded_proc = self.cmd.sifnoded_start(minimum_gas_prices=[0.5, "rowan"])  # TODO sifnoded_home=???

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
            self.cmd.yarn(["peggy:whiteList", addr, "true"], cwd=self.cmd.smart_contracts_dir)

        # ui/scripts/stack-launch.sh -> ui/scripts/_peggy.sh -> ui/chains/peggy/launch.sh
        # rm -rf ui/chains/peggy/relayerdb
        # ebrelayer is in $GOBIN, gets installed by "make install"
        ethereum_private_key = smart_contracts_env_ui_example_vars["ETHEREUM_PRIVATE_KEY"]
        ebrelayer_proc = self.cmd.ebrelayer_init("tcp://localhost:26657", "ws://localhost:7545/",
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

    def make_go_binaries(self):
        # make go binaries (TODO Makefile needs to be trimmed down, especially "find")
        self.cmd.execst(["make"], cwd=self.test_integration_dir, env={"BASEDIR": project_dir()})

    def run(self):
        self.cmd.mkdir(self.data_dir)

        self.make_go_binaries()

        self.cmd.install_smart_contracts_dependencies()

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
            ganache_proc = self.cmd.start_ganache_cli(block_time=block_time, host="0.0.0.0",
                mnemonic=self.ganache_mnemonic, network_id=self.network_id, port=7545, db=ganache_db_path,
                account_keys_path=account_keys_path)

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

        validator_moniker = netdef["moniker"]
        validator1_address = netdef["address"]
        validator1_password = netdef["password"]
        validator_mnemonic = netdef["mnemonic"].split(" ")
        chaindir = os.path.join(networks_dir, f"validators/{self.chainnet}/{validator_moniker}")
        sifnoded_home = os.path.join(chaindir, ".sifnoded")
        denom_whitelist_file = os.path.join(self.test_integration_dir, "whitelisted-denoms.json")
        # SIFNODED_LOG=$datadir/logs/sifnoded.log

        adminuser_addr = self.cmd.sifchain_init_integration(validator_moniker, validator_mnemonic, sifnoded_home, denom_whitelist_file, validator1_password)

        # Start sifnoded
        sifnoded_proc = self.cmd.sifnoded_start(tcp_url=self.tcp_url, minimum_gas_prices=[0.5, "rowan"], sifnoded_home=sifnoded_home)

        # TODO: should we wait for sifnoded to come up before continuing? If so, how do we do it?

        # TODO Process exits immediately with returncode 1
        # TODO Why does it not stop start-integration-env.sh?
        # rest_server_proc = self.cmd.popen(["sifnoded", "rest-server", "--laddr", "tcp://0.0.0.0:1317"])  # TODO cwd

        # test/integration/sifchain_start_ebrelayer.sh -> test/integration/sifchain_run_ebrelayer.sh
        # This script is also called from tests

        relayer_db_path = os.path.join(self.test_integration_dir, "sifchainrelayerdb")
        ebrelayer_proc = self.run_ebrelayer(netdef_json, validator1_address, validator_moniker, validator_mnemonic,
            ebrelayer_ethereum_private_key, bridge_registry_sc_addr, relayer_db_path)

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
            "SMART_CONTRACTS_DIR": project_dir("smart-contracts"),
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
            "MONIKER": validator_moniker,
            "VALIDATOR1_PASSWORD": validator1_password,
            "VALIDATOR1_ADDR": validator1_address,
            "MNEMONIC": " ".join(validator_mnemonic),
            "CHAINDIR": os.path.join(networks_dir, "validators", self.chainnet, validator_moniker),
            "SIFCHAIN_ADMIN_ACCOUNT": adminuser_addr,  # Needed by test_peggy_fees.py (via conftest.py)
            "EBRELAYER_DB": relayer_db_path,  # Created by sifchain_run_ebrelayer.sh, does not appear to be used anywhere at the moment
        }
        self.write_vagrantenv_sh()

        return ganache_proc, sifnoded_proc, ebrelayer_proc

    def write_vagrantenv_sh(self):
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
        # ETHEREUM_WEBSOCKET_ADDRESS (required), value=ws://localhost:7545/
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
        env = dict_merge(self.state_vars, {
            # For running test/integration/execute_integration_tests_against_*.sh
            "TEST_INTEGRATION_DIR": project_dir("test/integration"),
            "TEST_INTEGRATION_PY_DIR": project_dir("test/integration/src/py"),
            "SMART_CONTRACTS_DIR": self.cmd.smart_contracts_dir,
            "datadir": self.data_dir,  # Needed by test_rollback_chain.py that calls ganache_start.sh
            "GANACHE_KEYS_JSON": os.path.join(self.data_dir, "ganachekeys.json"),  # Needed by test_rollback_chain.py that calls ganache_start.sh
            "ETHEREUM_WEBSOCKET_ADDRESS": self.ethereum_websocket_address,   # Needed by test_ebrelayer_replay.py (and possibly others)
            "CHAINNET": self.chainnet,   # Needed by test_ebrelayer_replay.py (and possibly others)
        })
        vagrantenv_path = project_dir("test/integration/vagrantenv.sh")
        self.cmd.write_text_file(vagrantenv_path, joinlines(format_as_shell_env_vars(env)))
        self.cmd.write_text_file(project_dir("test/integration/vagrantenv.json"), json.dumps(env))

    def wait_for_sif_account(self, netdef_json, validator1_address):
        return self.cmd.execst(["python3", os.path.join(self.test_integration_dir, "src/py/wait_for_sif_account.py"),
            netdef_json, validator1_address], env={"USER1ADDR": "nothing"})

    def remove_and_add_sifnoded_keys(self, validator_moniker, validator_mnemonic):
        # Error: The specified item could not be found in the keyring
        # This is not neccessary during start-integration-env.sh (as the key does not exist yet), but is neccessary
        # during tests that restart ebrelayer
        # res = self.cmd.execst(["sifnoded", "keys", "delete", moniker, "--keyring-backend", "test"], stdin=["y"])
        self.cmd.sifnoded_keys_delete(validator_moniker)
        self.cmd.sifnoded_keys_add(validator_moniker, validator_mnemonic)

    def process_netdef(self, network_definition_file):
        # networks_dir = deploy/networks
        # File deploy/networks/network-definition.yml is created by "rake genesis:network:scaffold", specifically by
        # "sifgen network create"
        # We read it and convert to test/integration/vagrant/data/netdef.json
        netdef = exactly_one(yaml_load(self.cmd.read_text_file(network_definition_file)))
        netdef_json = os.path.join(self.data_dir, "netdef.json")
        self.cmd.write_text_file(netdef_json, json.dumps(netdef))
        return netdef, netdef_json

    def run_ebrelayer(self, netdef_json, validator1_address, validator_moniker, validator_mnemonic,
        ebrelayer_ethereum_private_key, bridge_registry_sc_addr, relayer_db_path):
        while not self.cmd.tcp_probe_connect("localhost", 26657):
            time.sleep(1)
        self.wait_for_sif_account(netdef_json, validator1_address)
        time.sleep(10)
        self.remove_and_add_sifnoded_keys(validator_moniker, validator_mnemonic)  # Creates ~/.sifnoded/keyring-tests/xxxx.address
        ebrelayer_proc = self.cmd.ebrelayer_init(self.tcp_url, self.ethereum_websocket_address, bridge_registry_sc_addr,
            validator_moniker, validator_mnemonic, self.chainnet, ethereum_private_key=ebrelayer_ethereum_private_key,
            node=self.tcp_url, keyring_backend="test", sign_with=validator_moniker,
            symbol_translator_file=os.path.join(self.test_integration_dir, "config/symbol_translator.json"),
            relayerdb_path=relayer_db_path, cwd=self.test_integration_dir)
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
        self.write_vagrantenv_sh()
        self.cmd.mkdir(self.data_dir)

        return self.restart_processes()

    def restart_processes(self):
        block_time = None
        ganache_db_path = self.state_vars["GANACHE_DB_DIR"]
        account_keys_path = os.path.join(self.data_dir, "ganachekeys.json")  # TODO this is in test/integration/vagrant/data, which is supposed to be cleared

        ganache_proc = self.cmd.start_ganache_cli(block_time=block_time, host="0.0.0.0", mnemonic=self.ganache_mnemonic,
            network_id=self.network_id, port=7545, db=ganache_db_path, account_keys_path=account_keys_path)

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
        self.smart_contracts_dir = project_dir("smart-contracts")
        gobin_dir = os.environ["GOBIN"]
        self.sifgen_cmd = os.path.join(gobin_dir, "sifgen")
        self.sifnoded_cmd = os.path.join(gobin_dir, "sifnoded")
        self.ebrelater_cmd = os.path.join(gobin_dir, "ebrelayer")

    # TODO Merge with super.make_go_binaries
    # Main Makefile requires GOBIN to be set to an absolute path. Compiled executables ebrelayer, sifgen and
    # sifnoded will be written there. The directory will be created if it doesn't exist yet.
    #
    # c.f. IntegrationEnvironment:
    # cd test/integration; BASEDIR=... make
    # (checks all *.go files and, runs make in $BASEDIR, touches sifnoded, removes ~/.sifnoded/localnet
    def _make_go_binaries(self):
        # Original: cd smart-contracts; make -C .. install
        res = self.cmd.execst(["make", "install"], cwd=project_dir())
        print(repr(res))

    def start_hardhat(self, hostname, port):
        # TODO We need to manaege smart-contracts/hardhat.config.ts + it also reads smart-contracts/.env via dotenv
        # TODO Handle failures, e.g. if the process is already running we get exit value 1 and
        # "Error: listen EADDRINUSE: address already in use 127.0.0.1:8545"
        proc = self.cmd,popen([os.path.join("node_modules", ".bin", "hardhat"), "node", "--hostname", hostname, "--port",
            str(port)], cwd=self.smart_contracts_dir)
        return proc

    def signer_array_to_ethereum_accounts(self, accounts, n_validators):
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

    def default_hardhat_accounts(self):
        # Hardhat doesn't provide a way to get the private keys of its default accounts, so just hardcode them for now.
        # TODO hardhat prints 20 accounts upon startup
        # Keep synced to smart-contracts/src/devenv/hardhatNode.ts:defaultHardhatAccounts
        # Format: [address, private_key]
        return [[
            "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
            "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
        ], [
            "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
            "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
        ], [
            "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc",
            "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
        ], [
            "0x90f79bf6eb2c4f870365e785982e1f101e93b906",
            "0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
        ], [
            "0x15d34aaf54267db7d7c367839aaf71a00a2c6a65",
            "0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
        ], [
            "0x9965507d1a55bcc2695c58ba16fb37d819b0a4dc",
            "0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
        ], [
            "0x976ea74026e726554db657fa54763abd0c3a0aa9",
            "0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
        ], [
            "0x14dc79964da2c08b23698b3d3cc7ca32193d9955",
            "0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
        ], [
            "0x23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f",
            "0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
        ], [
            "0xa0ee7a142d267c1f36714e4a8f75612f20a79720",
            "0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
        ]]

    def deploy_smart_contracts_hardhat(self):
        res = self.cmd.execst(["npx", "hardhat", "run", "scripts/deploy_contracts.ts", "--network", "localhost"],
            cwd=project_dir("smart-contracts"))
        # Skip first line "No need to generate any newer types"
        m = json.loads(stdout(res).splitlines()[1])
        return m["bridgeBank"], m["bridgeRegistry"], m["rowanContract"]

    def run_ebrelayer_peggy(self, tcp_url, websocket_address, bridge_registry_sc_addr, validator_moniker,
        validator_mnemonic, chain_id, symbol_translator_file, relayerdb_path, ethereum_address, ethereum_private_key
    ):
        return self.cmd.ebrelayer_init(tcp_url, websocket_address, bridge_registry_sc_addr, validator_moniker,
            validator_mnemonic, chain_id, ethereum_private_key=ethereum_private_key, ethereum_address=ethereum_address,
            node=tcp_url, keyring_backend="test", sign_with=validator_moniker,
            symbol_translator_file=symbol_translator_file, relayerdb_path=relayerdb_path)

    # Override
    def run(self):
        # self._make_go_binaries()

        self.cmd.rmdir(self.cmd.get_user_home(".sifnoded"))  # Purge test keyring backend

        hardhat_hostname = "localhost"
        hardhat_port = 8545
        hardhat_proc = self.start_hardhat(hardhat_hostname, hardhat_port)

        hardhat_validator_count = 1
        hardhat_network_id = 1
        hardhat_chain_id = 1
        hardhat_accounts = self.signer_array_to_ethereum_accounts(self.default_hardhat_accounts(), hardhat_validator_count)

        bridgebank_sc_addr, bridge_registry_sc_addr, rowan_sc_addr = self.deploy_smart_contracts_hardhat()

        chain_id = "localnet"
        sifnoded_network_dir = "/tmp/sifnodedNetwork"
        self.cmd.rmdir(sifnoded_network_dir)
        self.cmd.mkdir(sifnoded_network_dir)
        network_config_file = "/tmp/sifnodedConfig.yml"
        validator_count = 1
        seed_ip_address = "10.10.1.1"
        self.cmd.sifgen_create_network(chain_id, validator_count, sifnoded_network_dir, network_config_file, seed_ip_address)

        netdev_yml = exactly_one(yaml_load(self.cmd.read_text_file(network_config_file)))
        validator_moniker = netdev_yml["moniker"]
        validator_mnemonic = netdev_yml["mnemonic"].split(" ")

        chain_dir = os.path.join(sifnoded_network_dir, "validators", chain_id, validator_moniker)
        sifnoded_home = os.path.join(chain_dir, ".sifnoded")
        denom_whitelist_file = project_dir("test", "integration", "whitelisted-denoms.json")

        self.cmd.sifchain_init_peggy(validator_moniker, validator_mnemonic, sifnoded_home, denom_whitelist_file)

        tcp_url = "tcp://0.0.0.0:26657"
        sifnoded_proc = self.cmd.sifnoded_start(minimum_gas_prices=[0.5, "rowan"], tcp_url=tcp_url, sifnoded_home=sifnoded_home)

        # TODO Wait for account (somewhere)

        relayerdb_path = self.cmd.mktempdir()
        ethereum_address, ethereum_private_key = hardhat_accounts["validators"][0]
        ebrelayer_proc = self.run_ebrelayer_peggy(
            tcp_url,
            f"ws://{hardhat_hostname}:{hardhat_port}/",
            bridge_registry_sc_addr,
            validator_moniker,
            validator_mnemonic,
            chain_id,
            os.path.join(self.test_integration_dir, "config/symbol_translator.json"),
            relayerdb_path,
            ethereum_address,
            ethereum_private_key,
        )

        return hardhat_proc, sifnoded_proc, ebrelayer_proc

def run_with_logging_until_keypress(procs):
    o = []
    for i, p in enumerate(procs):
        path = f"/tmp/proc-{i}.log"
        pass
    procs_and_files = [[p, f"/tmp/proc-{i}.log"] for i, p in enumerate(procs)]

    pass

def cleanup_and_reset_state():
    # git checkout 4cb7322b6b282babd93a0d0aedda837c9134e84e deploy
    # pkill node; pkill ebrelayer; pkill sifnoded; rm -rvf $HOME/.sifnoded; rm -rvf ./vagrant/data; mkdir vagrant/data
    cmd = Command()
    cmd.execst(["pkill", "node"], check_exit=False)
    cmd.execst(["pkill", "ebrelayer"], check_exit=False)
    cmd.execst(["pkill", "sifnoded"], check_exit=False)

    # rm -rvf /tmp/tmp.xxxx (ganache DB, unique for every run)
    cmd.rmdir(project_dir("test/integration/sifchainrelayerdb"))  # TODO move to /tmp
    cmd.rmdir(project_dir("smart-contracts/build"))
    cmd.rmdir(project_dir("test/integration/vagrant/data"))
    cmd.rmdir(cmd.get_user_home(".sifnoded"))  # Probably needed for "--keyring-backend test"

    # Additional cleanup (not neccessary to make it work)
    # cmd.rm(project_dir("smart-contracts/combined.log"))
    # cmd.rmdir(project_dir("test/integration/.pytest_cache"))
    # cmd,rm(project_dir("smart-contracts/.env"))
    # cmd.rmdir(project_dir("deploy/networks"))
    # cmd.rmdir(project_dir("smart-contracts/.openzeppelin"))

    # rmdir ~/.cache/yarn
    time.sleep(3)

def killall(processes):
    # TODO Order - ebrelayer, sifnoded, ganache
    for p in processes:
        if p is not None:
            p.kill()
            p.wait()

def main(argv):
    logging.basicConfig(stream=sys.stdout, level=logging.DEBUG, format="%(message)s")
    what = argv[0] if argv else None
    cmd = Integrator()
    if what == "run-ui-env":
        e = UIStackEnvironment(cmd)
        e.stack_save_snapshot()
        e.stack_push()
    elif what == "run-integration-env":
        e = IntegrationTestsEnvironment(cmd)
        processes = e.run()
        input("Press ENTER to exit...")
        killall(processes)
    elif what == "create_snapshot":
        snapshot_name = argv[1]
        cleanup_and_reset_state()
        e = IntegrationTestsEnvironment(cmd)
        processes = e.run()
        # Give processes some time to settle, for example relayerdb must init and create its "relayerdb"
        time.sleep(45)
        killall(processes)
        # processes1 = e.restart_processes()
        e.create_snapshot(snapshot_name)
    elif what == "restore_snapshot":
        snapshot_name = argv[1]
        e = IntegrationTestsEnvironment(cmd)
        processes = e.restore_snapshot(snapshot_name)
        input("Press ENTER to exit...")
        killall(processes)
    elif what == "run-peggy-env":
        # hardhat, sifnoded, ebrelayer
        e = PeggyEnvironment(cmd)
        processes = e.run()
        input("Press ENTER to exit...")
        killall(processes)
    elif what == "fullclean":
        cmd.execst(["chmod", "-R", "+w", cmd.get_user_home("go")])
        cmd.rmdir(cmd.get_user_home("go"))
        cmd.mkdir(cmd.get_user_home("go"))
        cmd.rmdir(cmd.get_user_home(".npm"))
        cmd.rmdir(cmd.get_user_home(".npm-global"))
        cmd.mkdir(cmd.get_user_home(".npm-global"))
        cmd.rmdir(cmd.get_user_home(".cache/yarn"))
        cmd.rmdir(cmd.get_user_home(".sifnoded"))
        cmd.rmdir(cmd.get_user_home(".sifnode-integration"))
        cmd.rmdir(project_dir("smart-contracts/node_modules"))
        cmd.execst(["npm", "install", "-g", "ganache-cli", "dotenv", "yarn"], cwd=project_dir("smart-contracts"))
        cmd.install_smart_contracts_dependencies()
    elif what == "test-logging":
        ls_cmd = mkcmd(["ls", "-al", "."], cwd="/tmp")
        res = stdout_lines(cmd.execst(**ls_cmd))
        print(ls_cmd)

        # log_to_file(proc, "/tmp/test.txt")
    else:
        raise Exception("Missing/unknown command")


if __name__ == "__main__":
    main(sys.argv[1:])
