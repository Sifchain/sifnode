import json
import logging
import os
import shutil
import subprocess
import yaml  # pip install pyyaml
import urllib.request
import time


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
    return subprocess.Popen(args, env=env)

NULL_ADDRESS = "0x0000000000000000000000000000000000000000"


class Command:
    def execst(self, args, cwd=None, env=None, stdin=None, binary=False):
        if stdin is not None:
            if type(stdin) == list:
                stdin = "".join([line + "\n" for line in stdin])
        popen = subprocess.Popen(args, cwd=cwd, env=env, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=not binary)
        stdout_data, stderr_data = popen.communicate(input=stdin)
        if popen.returncode != 0:
            raise Exception("Command '{}' exited with returncode {}: {}".format(" ".join(args), popen.returncode, repr(stderr_data)))
        return popen.returncode, stdout_data, stderr_data

    def read_text_file(self, path):
        with open(path, "rt") as f:
            return f.read()  # TODO Convert to exec

    def rmdir(self, path):
        if os.path.exists(path):
            shutil.rmtree(path)  # TODO Convert to exec

    def copy_file(self, src, dst):
        shutil.copy(src, dst)

    def get_user_home(self, *paths):
        return os.path.join(os.environ["HOME"], *paths)

    def tar_create(self, path, tar_path, compression=None):
        comp_opts = {"gz": "z"}
        comp = comp_opts[compression] if compression in comp_opts else ""
        args = ["tar", "cf" + comp, tar_path, "."]
        self.execst(args, cwd=path)


class Ganache(Command):
    def start_ganache_cli(self, mnemonic=None, db=None, port=None, host=None, network_id=None, gas_price=None, gas_limit=None):
        args = ["ganache-cli"] + \
            (["--mnemonic", " ".join(mnemonic)] if mnemonic else []) + \
            (["--db", db] if db else []) + \
            (["--port", str(port)] if port is not None else []) + \
            (["--host", host] if host else []) + \
            (["--networkId", str(network_id)] if network_id is not None else []) + \
            (["--gasPrice", str(gas_price)] if gas_price is not None else []) + \
            (["--gasLimit", str(gas_limit)] if gas_limit is not None else [])
        return popen(args)


class Sifnoded(Command):
    def sifnoded_init(self, moniker, chain_id):
        args = ["sifnoded", "init", moniker, "--chain-id={}".format(chain_id)]
        self.execst(args)

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

    def sifnoded_add_genesis_account(self, address, tokens):
        tokens_str = ",".join([sif_format_amount(amount, denom) for amount, denom in tokens])
        args = ["sifnoded", "add-genesis-account", address, tokens_str]
        self.execst(args)

    def sifnoded_add_genesis_validators(self, address):
        args = ["sifnoded", "add-genesis-validators", address]
        res = self.execst(args)
        return res

    def sifnoded_tx_clp_create_pool(self, chain_id, keyring_backend, from_name, symbol, fees, native_amount, external_amount):
        args = ["sifnoded", "tx", "clp", "create-pool", "--chain-id={}".format(chain_id),
            "--keyring-backend={}".format(keyring_backend), "--from", from_name, "--symbol", symbol, "--fees",
            sif_format_amount(*fees), "--nativeAmount", str(native_amount), "--externalAmount", str(external_amount),
            "--yes"]
        self.execst(args)

    def sifnoded_launch(self, minimum_gas_prices=None):
        args = ["sifnoded", "start"] + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else [])
        return popen(args)

    def sifnoded_get_status(self, host, port):
        url = "http://{}:{}/node_info".format(host, port)
        return json.loads(http_get(url).decode("UTF-8"))


class Integrator(Ganache, Sifnoded, Command):
    def ebrelayer_init(self, ethereum_private_key, tendermind_node, web3_provider, bridge_registry_contract_address,
        validator_moniker, validator_mnemonic, chain_id, gas, gas_prices):
        env = {"ETHEREUM_PRIVATE_KEY": ethereum_private_key}
        args = ["ebrelayer", "init", tendermind_node, web3_provider, bridge_registry_contract_address,
            validator_moniker, validator_mnemonic, "--chain-id={}".format(chain_id), "--gas", str(gas), "--gas-prices",
            sif_format_amount(*gas_prices)]
        return popen(args, env=env)

    def sif_wait_up(self, host, port):
        while True:
            from urllib.error import URLError
            try:
                return self.sifnoded_get_status(host, port)
            except URLError:
                time.sleep(1)

    def yarn(self, args, cwd=None, env=None):
        return self.execst(["yarn"] + args, cwd=cwd, env=env)


class UIPlaybook:
    def __init__(self, cmd):
        self.cmd = cmd
        self.chain_id = "sifchain-local"
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
        self.cmd(["make", "install"], project_dir())

        # ui/scripts/stack-launch.sh -> ui/scripts/_eth.sh -> ui/chains/etc/launch.sh
        self.cmd.rmdir(self.ganache_db_path)
        self.cmd.yarn([], cwd=project_dir("ui/chains/eth"))  # Installs ui/chains/eth/node_modules
        # Note that this runs ganache-cli from $PATH whereas scripts start it with yarn in ui/chains/eth
        ganache_proc = self.cmd.start_ganache_cli(mnemonic=self.ethereum_root_mnemonic, db=self.ganache_db_path,
            port=7545, network_id=5777, gas_price=20000000000, gas_limit=6721975, host="0.0.0.0")

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
        smart_contracts_dir = project_dir("smart-contracts")
        self.cmd.copy_file(os.path.join(smart_contracts_dir, ".env.ui.example"), os.path.join(smart_contracts_dir, ".env"))
        self.cmd.yarn([], cwd=smart_contracts_dir)
        self.cmd.yarn(["migrate"], cwd=smart_contracts_dir)

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
        def get_smart_contract_address(path):
            return json.loads(self.cmd.read_text_file(project_dir(path)))["networks"]["5777"]["address"]

        atk_address, btk_address, usdc_address, link_address = [
            get_smart_contract_address(project_dir(f"ui/chains/eth/build/contracts/{x}.json"))
            for x in ["AliceToken", "BobToken", "UsdCoin", "LinkCoin"]
        ]
        bridge_token_address = get_smart_contract_address(project_dir("smart-contracts/build/contracts/BridgeToken.json"))
        bridge_registry_address = get_smart_contract_address(project_dir("smart-contracts/build/contracts/BridgeRegistry.json"))

        # From smart-contracts/.env.ui.example
        smart_contracts_env_ui_example_vars = {
            "ETHEREUM_PRIVATE_KEY": "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3",
            "INFURA_PROJECT_ID": "JFSH7439sjsdtqTM23Dz",
            "LOCAL_PROVIDER": "http://localhost:7545",
        }

        def set_token_lock_burn_limit(update_address, amount):
            env = smart_contracts_env_ui_example_vars.copy()
            env["UPDATE_ADDRESS"] = update_address
            # Needs: ETHEREUM_PRIVATE_KEY, INFURA_PROJECT_ID, LOCAL_PROVIDER, UPDATE_ADDRESS
            # Hint: call web3 directly, avoid npx + truffle + script
            # Maybe: self.cmd.yarn(["integrationtest:setTokenLockBurnLimit", str(amount)])
            self.cmd.execst(["npx", "truffle", "exec", "scripts/setTokenLockBurnLimit.js", str(amount)], env=env, cwd=smart_contracts_dir)

        set_token_lock_burn_limit(NULL_ADDRESS, 31*10**18)
        set_token_lock_burn_limit(bridge_token_address, 10**25)
        set_token_lock_burn_limit(atk_address, 10**25)
        set_token_lock_burn_limit(btk_address, 10**25)
        set_token_lock_burn_limit(usdc_address, 10**25)
        set_token_lock_burn_limit(link_address, 10**25)
        # signal migrate-complete

        # Whitelist test tokens
        for addr in [atk_address, btk_address, usdc_address, link_address]:
            self.cmd.yarn(["peggy:whitelist", addr, "true"], cwd=smart_contracts_dir)

        # ui/scripts/stack-launch.sh -> ui/scripts/_peggy.sh -> ui/chains/peggy/launch.sh
        # rm -rf ui/chains/peggy/relayerdb
        ethereum_private_key = smart_contracts_env_ui_example_vars["ETHEREUM_PRIVATE_KEY"]
        ebrelayer_proc = self.cmd.ebrelayer_init(ethereum_private_key, "tcp://localhost:26657", "ws://localhost:7545/",
            bridge_registry_address, self.shadowfiend_name, self.shadowfiend_mnemonic, f"--chain-id={self.chain_id}",
            5*10**12, [0.5, "rowan"])

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
        self.cmd.mkdir(snapshots_dir)
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
        commit = exactly_one(stdout_lines(self.cmd.execst(["git", "rev-parse", "HEAD"], cwd=project_dir())))
        branch = exactly_one(stdout_lines(self.cmd.execst(["git", "rev-parse", "--abbrev-ref", "HEAD"], cwd=project_dir())))

        image_root = "ghcr.io/sifchain/sifnode/ui-stack"
        image_name = "{}:{}".format(image_root, commit)
        stable_tag = "{}:{}".format(image_root, branch.replace("/", "__"))

        running_in_ci = bool(os.environ.get("CI"))

        if running_in_ci:
            # # reverse grep for go.mod because on CI this can be altered by installing go dependencies
            # if [[ -z "$CI" && ! -z "$(git status --porcelain --untracked-files=no)" ]]; then
            #   echo "Git workspace must be clean to save git commit hash"
            #   exit 1
            # fi
            pass

        log.info("Github Registry Login found.")
        log.info("Building new container...")
        log.info(f"New image name: {image_name}")

        self.cmd.execst(["docker", "build", "-f", project_dir("ui/scripts/stack.Dockerfile"), "-t", image_name, "."],
            cwd=project_dir(), env={"DOCKER_BUILDKIT", "1"})

        if running_in_ci:
            log.info(f"Tagging image as {stable_tag}...")
            self.cmd.execst(["docker", "tag", image_name, stable_tag])
            self.cmd.execst(["docker", "push", image_name])
            self.cmd.execst(["docker", "push", stable_tag])


def main():
    cmd = Integrator()
    ui_playbook = UIPlaybook(cmd)
    ui_playbook.stack_save_snapshot()
    ui_playbook.stack_push()

if __name__ == "__main__":
    main()
