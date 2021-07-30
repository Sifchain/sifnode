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
    return res[0].splitlines()

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
    # From ui/chains/credentials.sh
    SHADOWFIEND_NAME = "shadowfiend"
    SHADOWFIEND_MNEMONIC = ["race", "draft", "rival", "universe", "maid", "cheese", "steel", "logic", "crowd", "fork",
        "comic", "easy", "truth", "drift", "tomorrow", "eye", "buddy", "head", "time", "cash", "swing", "swift",
        "midnight", "borrow"]
    AKASHA_NAME = "akasha"
    AKASHA_MNEMONIC = ["hand", "inmate", "canvas", "head", "lunar", "naive", "increase", "recycle", "dog", "ecology",
        "inhale", "december", "wide", "bubble", "hockey", "dice", "worth", "gravity", "ketchup", "feed", "balance",
        "parent", "secret", "orchard"]
    JUNIPER_NAME = "juniper"
    JUNIPER_MNEMONIC = ["clump", "genre", "baby", "drum", "canvas", "uncover", "firm", "liberty", "verb", "moment",
        "access", "draft", "erupt", "fog", "alter", "gadget", "elder", "elephant", "divide", "biology", "choice",
        "sentence", "oppose", "avoid"]
    ETHEREUM_ROOT_MNEMONIC = ["candy", "maple", "cake", "sugar", "pudding", "cream", "honey", "rich", "smooth",
        "crumble", "sweet", "treat"]
    CHAIN_ID = "sifchain-local"

    def __init__(self, cmd):
        self.cmd = cmd

    def eth_start(self):
        return

    def run(self):
        # self.cmd.yarn([], cwd=project_dir())  # TODO Where?
        # ganache_cli_proc = self.eth_start()
        # print(repr(ganache_cli_proc))
        # self.ui_sif_launch()
        # sifnoded_proc = self.ui_sif_start()
        # print(repr(sifnoded_proc))
        # self.cmd.sif_wait_up("localhost", 1317)
        self.stack_save_snapshot()

    def stack_save_snapshot(self):
        # ui-stack.yml
        # cd .; go get -v -t -d ./...
        # cd ui; yarn install --frozen-lockfile --silent
        # Compile smart contracts:
        # cd ui; yarn build

        keyring_backend = "test"

        # yarn stack --save-snapshot -> ui/scripts/stack.sh -> ui/scripts/stack-save-snapshot.sh
        # rm ui/node_modules/.migrate-complete

        # yarn stack --save-snapshot -> ui/scripts/stack.sh -> ui/scripts/stack-save-snapshot.sh => ui/scripts/stack-launch.sh
        # ui/scripts/stack-launch.sh -> ui/scripts/_sif-build.sh -> ui/chains/sif/build.sh
        # killall sifnoded
        # rm $(which sifnoded)
        self.cmd.rmdir(self.cmd.get_user_home(".sifnoded"))
        self.cmd(["make", "install"], project_dir())

        # ui/scripts/stack-launch.sh -> ui/scripts/_eth.sh -> ui/chains/etc/launch.sh
        self.cmd.rmdir(self.cmd.get_user_home(".ganachedb"))
        self.cmd.yarn([], cwd=project_dir("ui/chains/eth"))  # Installs ui/chains/eth/node_modules
        # Note that this runs ganache-cli from $PATH whereas scripts start it with yarn in ui/chains/eth
        ganache_proc = self.cmd.start_ganache_cli(mnemonic=UIPlaybook.ETHEREUM_ROOT_MNEMONIC,
            db=self.cmd.get_user_home(".ganachedb"), port=7545, network_id=5777, gas_price=20000000000,
            gas_limit=6721975, host="0.0.0.0")

        # ui/scripts/stack-launch.sh -> ui/scripts/_sif.sh -> ui/chains/sif/launch.sh
        self.cmd.sifnoded_init("test", UIPlaybook.CHAIN_ID)
        self.cmd.copy_file(project_dir("ui", "chains", "sif", "app.toml"), self.cmd.get_user_home(".sifnoded", "config", "app.toml"))
        log.info(f"Generating deterministic account - {UIPlaybook.SHADOWFIEND_NAME}...")
        shadowfiend_account = self.cmd.sifnoded_generate_deterministic_account(UIPlaybook.SHADOWFIEND_NAME, UIPlaybook.SHADOWFIEND_MNEMONIC)
        log.info(f"Generating deterministic account - {UIPlaybook.AKASHA_NAME}...")
        akasha_account = self.cmd.sifnoded_generate_deterministic_account(UIPlaybook.AKASHA_NAME, UIPlaybook.AKASHA_MNEMONIC)
        log.info(f"Generating deterministic account - {UIPlaybook.JUNIPER_NAME}...")
        juniper_account = self.cmd.sifnoded_generate_deterministic_account(UIPlaybook.JUNIPER_NAME, UIPlaybook.JUNIPER_MNEMONIC)
        # shadowfiend_address = self.sifnoded_keys_show(SHADOWFIEND_NAME)[0]["address"]
        # akasha_address = self.sifnoded_keys_show(AKASHA_NAME)[0]["address"]
        # juniper_address = self.sifnoded_keys_show(JUNIPER_NAME)[0]["address"]
        shadowfiend_address = shadowfiend_account["address"]
        akasha_address = akasha_account["address"]
        juniper_address = juniper_account["address"]

        tokens_shadowfiend = [[10**29, "rowan"], [10**29, "catk"], [10**29, "cbtk"], [10**29, "ceth"], [10**29, "cusdc"], [10**29, "clink"], [10**26, "stake"]]
        tokens_akasha = [[10**29, "rowan"], [10**29, "catk"], [10**29, "cbtk"], [10**29, "ceth"], [10**29, "cusdc"], [10**29, "clink"], [10**26, "stake"]]
        tokens_juniper = [[10**22, "rowan"], [10**22, "cusdc"], [10**20, "clink"], [10**20, "ceth"]]
        self.cmd.sifnoded_add_genesis_account(shadowfiend_address, tokens_shadowfiend)
        self.cmd.sifnoded_add_genesis_account(akasha_address, tokens_akasha)
        self.cmd.sifnoded_add_genesis_account(juniper_address, tokens_juniper)

        shadowfiend_address_bech_val = self.cmd.sifnoded_keys_show(UIPlaybook.SHADOWFIEND_NAME, bech="val")[0]["address"]
        self.cmd.sifnoded_add_genesis_validators(shadowfiend_address_bech_val)

        amount = sif_format_amount(10**24, "stake")
        self.cmd.execst(["sifnoded", "gentx", UIPlaybook.SHADOWFIEND_NAME, amount, "--chain-id={}".format(UIPlaybook.CHAIN_ID), "--keyring-backend={}".format(keyring_backend)])

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
        self.cmd.sifnoded_tx_clp_create_pool(UIPlaybook.CHAIN_ID, keyring_backend, "akasha", "catk", [10**5, "rowan"], 10**25, 10**25)
        time.sleep(5)
        log.info("Creating liquidity pool from cbtk:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(UIPlaybook.CHAIN_ID, keyring_backend, "akasha", "cbtk", [10**5, "rowan"], 10**25, 10**25)
        # should now be able to swap from catk:cbtk
        time.sleep(5)
        log.info("Creating liquidity pool from ceth:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(UIPlaybook.CHAIN_ID, keyring_backend, "akasha", "ceth", [10**5, "rowan"], 10**25, 83*10**20)
        # should now be able to swap from x:ceth
        time.sleep(5)
        log.info("Creating liquidity pool from cusdc:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(UIPlaybook.CHAIN_ID, keyring_backend, "akasha", "cusdc", [10**5, "rowan"], 10**25, 10**25)
        time.sleep(5)
        log.info("Creating liquidity pool from clink:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(UIPlaybook.CHAIN_ID, keyring_backend, "akasha", "clink", [10**5, "rowan"], 10**25, 588235*10**18)
        time.sleep(5)
        log.info("Creating liquidity pool from ctest:rowan...")
        self.cmd.sifnoded_tx_clp_create_pool(UIPlaybook.CHAIN_ID, keyring_backend, "akasha", "ctest", [10**5, "rowan"], 10**25, 10**13)

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
            bridge_registry_address, UIPlaybook.SHADOWFIEND_NAME, UIPlaybook.SHADOWFIEND_MNEMONIC,
            "--chain-id={}".format(UIPlaybook.CHAIN_ID), 5*10**12, [0.5, "rowan"])

        # At this point we have 3 running processes - ganache_proc, sifnoded_proc and ebrelayer_proc
        # await sif-node-up and migrate-complete

        time.sleep(30)
        # ui/scripts/_snapshot.sh

        # ui/scripts/stack-pause.sh
        # killall sifnoded sifnoded ebrelayer ganache-cli
        sifnoded_proc.kill()
        ebrelayer_proc.kill()
        ganache_proc.kill()
        time.sleep(10)

        # ui/chains/peggy/snapshot.sh:
        # mkdir -p ui/chains/snapshots
        # mkdir -p ui/chains/peggy/relayerdb
        # cd ui/chains/peggy/relayerdb && tar -zcvf ui/chains/snapshots/peggy.tar.gz
        self.cmd.tar_create(project_dir("ui/chains/peggy/relayerdb"), project_dir("ui/chains/snapshots/peggy.tar.gz"), compression="gz")
        # mkdir -p smart-contracts/build
        # cd smart-contracts/build && tar -zcvf ui/chains/snapshots/peggy_build.tar.gz
        self.cmd.tar_create(project_dir("smart-contracts/build"), project_dir("ui/chains/snapshots/peggy_build.tar.gz"), compression="gz")

        # ui/chains/sif/snapshot.sh:
        # mkdir -p ui/chains/snapshots
        # cd ~/.sifnoded && tar -zcvf ui/chains/snapshots/sif.tar.gz
        self.cmd.tar_create(self.cmd.get_user_home(".sifnoded"), project_dir("ui/chains/snapshots/sif.tar.gz"), compression="gz")

        # ui/chains/etc/snapshot.sh:
        # cd ~/.ganachedb && tar -zcvf ui/chains/snapshots/eth.tar.gz
        self.cmd.tar_create(self.cmd.get_user_home(".ganachedb"), project_dir("ui/chains/snapshots/eth.tar.gz"), compression="gz")


def main():
    cmd = Integrator()
    ui_playbook = UIPlaybook(cmd)
    ui_playbook.run()

if __name__ == "__main__":
    main()
