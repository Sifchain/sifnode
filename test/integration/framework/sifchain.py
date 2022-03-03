import hashlib
import json
import time
from command import buildcmd
from common import *


def sifchain_denom_hash(network_descriptor, token_contract_address):
    assert on_peggy2_branch
    assert token_contract_address.startswith("0x")
    s = str(network_descriptor) + token_contract_address.lower()
    return "sif" + hashlib.sha256(s.encode("UTF-8")).digest().hex()


class Sifnoded:
    def __init__(self, cmd, home=None):
        self.cmd = cmd
        self.binary = "sifnoded"
        self.home = home
        self.keyring_backend = "test"
        # self.sifnoded_burn_gas_cost = 16 * 10**10 * 393000  # see x/ethbridge/types/msgs.go for gas
        # self.sifnoded_lock_gas_cost = 16 * 10**10 * 393000

    def init(self, moniker, chain_id):
        args = [self.binary, "init", moniker, "--chain-id", chain_id]
        res = self.cmd.execst(args)
        return json.loads(res[2])  # output is on stderr

    def keys_list(self):
        args = ["keys", "list", "--output", "json"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return json.loads(stdout(res))

    def keys_show(self, name, bech=None):
        args = ["keys", "show", name] + \
            (["--bech", bech] if bech else [])
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return yaml_load(stdout(res))

    def get_val_address(self, moniker):
        res = self.sifnoded_exec(["keys", "show", "-a", "--bech", "val", moniker], keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        expected = exactly_one(stdout_lines(res))
        result = exactly_one(self.keys_show(moniker, bech="val"))["address"]
        assert result == expected
        return result

    # How "sifnoded keys add <name> --keyring-backend test" works:
    # If name does not exist yet, it creates it and returns a yaml
    # If name alredy exists, prompts for overwrite (y/n) on standard input, generates new address/pubkey/mnemonic
    # Directory used is xxx/keyring-test if "--home xxx" is specified, otherwise $HOME/.sifnoded/keyring-test

    def keys_add(self, moniker, mnemonic):
        stdin = [" ".join(mnemonic)]
        res = self.sifnoded_exec(["keys", "add", moniker, "--recover"], keyring_backend=self.keyring_backend,
            sifnoded_home=self.home, stdin=stdin)
        account = exactly_one(yaml_load(stdout(res)))
        return account

    # Creates a new key in the keyring and returns its address ("sif1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx").
    # Since this is a test keyring, we don't need to save the generated private key.
    # If we wanted to recreate it, we can capture the mnemonic from the message that is printed to stderr.
    def keys_add_1(self, moniker):
        res = self.sifnoded_exec(["keys", "add", moniker], keyring_backend=self.keyring_backend, sifnoded_home=self.home, stdin=["y"])
        account = exactly_one(yaml_load(stdout(res)))
        unused_mnemonic = stderr(res).splitlines()[-1].split(" ")
        return account

    def keys_delete(self, name):
        self.cmd.execst(["sifnoded", "keys", "delete", name, "--keyring-backend", self.keyring_backend], stdin=["y"], check_exit=False)

    def add_genesis_account(self, sifnodeadmin_addr, tokens):
        tokens_str = ",".join([sif_format_amount(amount, denom) for amount, denom in tokens])
        self.sifnoded_exec(["add-genesis-account", sifnodeadmin_addr, tokens_str], sifnoded_home=self.home)

    def add_genesis_validators(self, address):
        args = ["sifnoded", "add-genesis-validators", address]
        res = self.cmd.execst(args)
        return res

    # At the moment only on future/peggy2 branch, called from PeggyEnvironment
    def add_genesis_validators_peggy(self, evm_network_descriptor, valoper, validator_power):
        self.sifnoded_exec(["add-genesis-validators", str(evm_network_descriptor), valoper, str(validator_power)],
            sifnoded_home=self.home)

    def set_genesis_oracle_admin(self, address):
        self.sifnoded_exec(["set-genesis-oracle-admin", address], sifnoded_home=self.home)

    def set_genesis_whitelister_admin(self, address):
        self.sifnoded_exec(["set-genesis-whitelister-admin", address], sifnoded_home=self.home)

    def set_gen_denom_whitelist(self, denom_whitelist_file):
        self.sifnoded_exec(["set-gen-denom-whitelist", denom_whitelist_file], sifnoded_home=self.home)

    # At the moment only on future/peggy2 branch, called from PeggyEnvironment
    # This was split from init_common
    def peggy2_add_account(self, name, tokens, is_admin=False):
        # TODO Peggy2 devenv feed "yes\nyes" into standard input, we only have "y\n"
        account = self.keys_add_1(name)
        account_address = account["address"]

        self.add_genesis_account(account_address, tokens)
        if is_admin:
            self.set_genesis_oracle_admin(account_address)
        self.set_genesis_whitelister_admin(account_address)
        return account_address

    def peggy2_add_relayer_witness_account(self, name, tokens, evm_network_descriptor, validator_power, denom_whitelist_file):
        account_address = self.peggy2_add_account(name, tokens)  # Note: is_admin=False
        # Whitelist relayer/witness account
        valoper = self.get_val_address(name)
        self.set_gen_denom_whitelist(denom_whitelist_file)
        self.add_genesis_validators_peggy(evm_network_descriptor, valoper, validator_power)
        return account_address

    def tx_clp_create_pool(self, chain_id, from_name, symbol, fees, native_amount, external_amount):
        args = ["tx", "clp", "create-pool", "--chain-id", chain_id, "--from", from_name, "--symbol", symbol,
            "--fees", sif_format_amount(*fees), "--nativeAmount", str(native_amount), "--externalAmount",
            str(external_amount), "--yes"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend)  # TODO home?
        return yaml_load(stdout(res))

    def peggy2_token_registry_register_all(self, registry_path, gas_prices, gas_adjustment, from_account,
        chain_id
    ):
        args = ["tx", "tokenregistry", "register-all", registry_path, "--gas-prices", sif_format_amount(*gas_prices),
            "--gas-adjustment", str(gas_adjustment), "--from", from_account, "--chain-id", chain_id, "--yes"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return [json.loads(x) for x in stdout(res).splitlines()]

    def peggy2_set_cross_chain_fee(self, admin_account_address, network_id, ethereum_cross_chain_fee_token,
        cross_chain_fee_base, cross_chain_lock_fee, cross_chain_burn_fee, admin_account_name, chain_id, gas_prices,
        gas_adjustment
    ):
        # Checked OK
        args = ["tx", "ethbridge", "set-cross-chain-fee", admin_account_address, str(network_id),
            ethereum_cross_chain_fee_token, str(cross_chain_fee_base), str(cross_chain_lock_fee),
            str(cross_chain_burn_fee), "--from", admin_account_name, "--chain-id", chain_id, "--gas-prices",
            sif_format_amount(*gas_prices), "--gas-adjustment", str(gas_adjustment), "-y"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return res

    def sifnoded_start(self, tcp_url=None, minimum_gas_prices=None, log_format_json=False, log_file=None):
        sifnoded_exec_args = self.build_start_cmd(tcp_url=tcp_url, minimum_gas_prices=minimum_gas_prices,
            log_format_json=log_format_json)
        return self.cmd.spawn_asynchronous_process(sifnoded_exec_args, log_file=log_file)

    def build_start_cmd(self, tcp_url=None, minimum_gas_prices=None, log_format_json=False):
        args = [self.binary, "start"] + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else []) + \
            (["--rpc.laddr", tcp_url] if tcp_url else []) + \
            (["--log_level", "debug"] if log_format_json else []) + \
            (["--log_format", "json"] if log_format_json else []) + \
            (["--home", self.home] if self.home else [])
        return buildcmd(args)

    def sifnoded_exec(self, args, sifnoded_home=None, keyring_backend=None, stdin=None, cwd=None):
        args = [self.binary] + args + \
            (["--home", sifnoded_home] if sifnoded_home else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend else [])
        res = self.cmd.execst(args, stdin=stdin, cwd=cwd)
        return res

    def get_status(self, host, port):
        url = "http://{}:{}/node_info".format(host, port)
        return json.loads(http_get(url).decode("UTF-8"))

    def wait_up(self, host, port):
        while True:
            from urllib.error import URLError
            try:
                return self.get_status(host, port)
            except URLError:
                time.sleep(1)


class Sifgen:
    def __init__(self, cmd):
        self.cmd = cmd
        self.binary = "sifgen"

    # Reference: docker/localnet/sifnode/root/scripts/sifnode.sh (branch future/peggy2):
    # sifgen node create "$CHAINNET" "$MONIKER" "$MNEMONIC" --bind-ip-address "$BIND_IP_ADDRESS" --standalone --keyring-backend test
    def create_standalone(self, chainnet, moniker, mnemonic, bind_ip_address, keyring_backend=None):
        args = ["node", "create", chainnet, moniker, mnemonic, bind_ip_address]
        return self.sifgen_exec(args, keyring_backend=keyring_backend)

    def sifgen_exec(self, args, keyring_backend=None, cwd=None, env=None):
        args = [self.binary] + args + \
            (["--keyring-backend", keyring_backend] if keyring_backend else [])
        return self.cmd.execst(args, cwd=cwd, env=env)


class Ebrelayer:
    def __init__(self, cmd):
        self.cmd = cmd
        self.binary = "ebrelayer"

    def peggy2_build_ebrelayer_cmd(self, init_what, network_descriptor, tendermint_node, web3_provider,
        bridge_registry_contract_address, validator_mnemonic, chain_id, node=None, keyring_backend=None,
        sign_with=None, symbol_translator_file=None, relayerdb_path=None, log_format=None, extra_args=None,
        ethereum_private_key=None, ethereum_address=None, home=None, cwd=None
    ):
        env = _env_for_ethereum_address_and_key(ethereum_address, ethereum_private_key)
        args = [
            self.binary,
            init_what,
            "--network-descriptor", str(network_descriptor),  # Network descriptor for the chain (31337)
            "--tendermint-node", tendermint_node,  # URL to tendermint node
            "--web3-provider", web3_provider,  # Ethereum web3 service address (ws://localhost:8545/)
            "--bridge-registry-contract-address", bridge_registry_contract_address,
            "--validator-mnemonic", validator_mnemonic,
            "--chain-id", chain_id  # chain ID of tendermint node (localnet)
        ] + \
            (extra_args if extra_args else []) + \
            (["--node", node] if node else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend else []) + \
            (["--from", sign_with] if sign_with else []) + \
            (["--relayerdb-path", relayerdb_path] if relayerdb_path else []) + \
            (["--home", home] if home else []) + \
            (["--symbol-translator-file", symbol_translator_file] if symbol_translator_file else []) + \
            (["--log_format", log_format] if log_format else [])
        return buildcmd(args, env=env, cwd=cwd)

    # Legacy stuff - pre-peggy2
    # Called from IntegrationContext
    def init(self, tendermind_node, web3_provider, bridge_registry_contract_address, validator_moniker,
        validator_mnemonic, chain_id, ethereum_private_key=None, ethereum_address=None, gas=None, gas_prices=None,
        node=None, keyring_backend=None, sign_with=None, symbol_translator_file=None, relayerdb_path=None,
        trace=True, cwd=None, log_file=None
    ):
        env = _env_for_ethereum_address_and_key(ethereum_address, ethereum_private_key)
        args = [self.binary, "init", tendermind_node, web3_provider, bridge_registry_contract_address,
            validator_moniker, " ".join(validator_mnemonic), "--chain-id={}".format(chain_id)] + \
            (["--gas", str(gas)] if gas is not None else []) + \
            (["--gas-prices", sif_format_amount(*gas_prices)] if gas_prices is not None else []) + \
            (["--node", node] if node is not None else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend is not None else []) + \
            (["--from", sign_with] if sign_with is not None else []) + \
            (["--symbol-translator-file", symbol_translator_file] if symbol_translator_file else []) + \
            (["--relayerdb-path", relayerdb_path] if relayerdb_path else []) + \
            (["--trace"] if trace else [])
        return self.cmd.popen(args, env=env, cwd=cwd, log_file=log_file)


# This is probably useful for any program that uses web3 library in the same way
# ETHEREUM_ADDRESS has to start with "0x" and ETHEREUM_PRIVATE_KEY has to be without "0x".
def _env_for_ethereum_address_and_key(ethereum_address, ethereum_private_key):
    env = {}
    if ethereum_private_key:
        assert not ethereum_private_key.startswith("0x")
        env["ETHEREUM_PRIVATE_KEY"] = ethereum_private_key
    if ethereum_address:
        assert ethereum_address.startswith("0x")
        env["ETHEREUM_ADDRESS"] = ethereum_address
    return env or None  # Avoid passing empty environment
