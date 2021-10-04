import json
import time
from command import Command
from common import *


class Sifnoded(Command):
    def __init__(self):
        self.binary = "sifnoded"
        # home = None
        # keyring_backend = None

    def sifnoded_init(self, moniker, chain_id):
        args = [self.binary, "init", moniker, "--chain-id", chain_id]
        res = self.execst(args)
        return json.loads(res[2])  # output is on stderr

    def sifnoded_keys_show(self, name, bech=None, keyring_backend=None, home=None):
        keyring_backend = keyring_backend or "test"
        args = ["keys", "show", name] + \
            (["--bech", bech] if bech else [])
        res = self.sifnoded_exec(args, keyring_backend=keyring_backend, sifnoded_home=home)
        return yaml_load(stdout(res))

    def sifnoded_get_val_address(self, moniker):
        expected = exactly_one(stdout_lines(self.sifnoded_exec(["keys", "show", "-a", "--bech", "val", moniker], keyring_backend="test")))
        result = exactly_one(self.sifnoded_keys_show(moniker, bech="val", keyring_backend="test"))["address"]
        assert result == expected
        return result

    def sifnoded_keys_add(self, moniker, mnemonic):
        stdin = [" ".join(mnemonic)]
        res = self.sifnoded_exec(["keys", "add", moniker, "--recover"], keyring_backend="test", stdin=stdin)
        return yaml_load(stdout(res))

    # How "sifnoded keys add <name> --keyring-backend test" works:
    # If name does not exist yet, it creates it and returns a yaml
    # If name alredy exists, prompts for overwrite (y/n) on standard input, generates new address/pubkey/mnemonic
    # Directory used is xxx/keyring-test if "--home xxx" is specified, otherwise $HOME/.sifnoded/keyring-test

    def sifnoded_keys_add_1(self, moniker):
        res = self.sifnoded_exec(["keys", "add", moniker], keyring_backend="test", stdin=["y"])
        return exactly_one(yaml_load(stdout(res)))

    # From peggy
    # @TODO Passing mnemonic to stdin is useless, only "y/n" makes sense, probably could use sifnoded_keys_add_1
    # See smart-contracts/src/devenv/sifnoded.ts:addValidatorKeysToTestKeyring
    def sifnoded_keys_add_2(self, moniker, mnemonic):
        stdin = [" ".join(mnemonic)]
        res = self.sifnoded_exec(["keys", "add", moniker], keyring_backend="test", stdin=stdin)
        result = exactly_one(yaml_load(stdout(res)))
        # {"name": "<moniker>", "type": "local", "address": "sif1...", "pubkey": "sifpub1...", "mnemonic": "", "threshold": 0, "pubkeys": []}
        return result

    def sifnoded_keys_delete(self, name):
        self.execst(["sifnoded", "keys", "delete", name, "--keyring-backend", "test"], stdin=["y"], check_exit=False)

    def sifnoded_add_genesis_account(self, sifnodeadmin_addr, tokens, sifnoded_home=None):
        tokens_str = ",".join([sif_format_amount(amount, denom) for amount, denom in tokens])
        self.sifnoded_exec(["add-genesis-account", sifnodeadmin_addr, tokens_str], sifnoded_home=sifnoded_home)

    def sifnoded_add_genesis_validators(self, address):
        args = ["sifnoded", "add-genesis-validators", address]
        res = self.execst(args)
        return res

    # At the moment only on future/peggy2 branch, called from PeggyEnvironment
    def sifnoded_add_genesis_validators_peggy(self, unknown_parameter_1, valoper, unknown_parameter_2, sifnoded_home):
        self.sifnoded_exec(["add-genesis-validators", str(unknown_parameter_1), valoper, str(unknown_parameter_2)],
            sifnoded_home=sifnoded_home)

    def sifnoded_tx_clp_create_pool(self, chain_id, keyring_backend, from_name, symbol, fees, native_amount, external_amount):
        args = [self.binary, "tx", "clp", "create-pool", "--chain-id={}".format(chain_id),
            "--keyring-backend={}".format(keyring_backend), "--from", from_name, "--symbol", symbol, "--fees",
            sif_format_amount(*fees), "--nativeAmount", str(native_amount), "--externalAmount", str(external_amount),
            "--yes"]
        res = self.execst(args)
        return yaml_load(stdout(res))

    def sifnoded_start(self, tcp_url=None, minimum_gas_prices=None, sifnoded_home=None, log_file=None,
        log_format_json=False
    ):
        args = [self.binary, "start"] + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else []) + \
            (["--rpc.laddr", tcp_url] if tcp_url else []) + \
            (["--log_format", "json"] if log_format_json else []) + \
            (["--home", sifnoded_home] if sifnoded_home else [])
        return self.popen(args, log_file=log_file)

    def sifnoded_exec(self, args, sifnoded_home=None, keyring_backend=None, stdin=None, cwd=None):
        args = [self.binary] + args + \
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

    def peggy2_init_relayer(
            self,
            network_descriptor,
            tendermint_node,
            web3_provider,
            bridge_registry_contract_address,
            validator_moniker,
            validator_mnemonic,
            chain_id,
            symbol_translator_file,
            ethereum_address,
            ethereum_private_key,
            keyring_backend=None,
            log_file=None,
            cwd=None,
    ):
        # Usage:
        #   ebrelayer init-relayer [networkDescriptor] [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMnemonic] [flags]
        #
        # Examples:
        # ebrelayer init-relayer 1 tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 mnemonic --chain-id=peggy
        return self.__peggy2_init_common("init-relayer", network_descriptor, tendermint_node, web3_provider,
            bridge_registry_contract_address, validator_mnemonic, chain_id=chain_id, node=tendermint_node,
            sign_with=validator_moniker, symbol_translator_file=symbol_translator_file,
            ethereum_address=ethereum_address, ethereum_private_key=ethereum_private_key,
            keyring_backend=keyring_backend, cwd=cwd, log_file=log_file)

    def peggy2_init_witness(
            self,
            network_descriptor,
            tendermint_node,
            web3_provider,
            bridge_registry_contract_address,
            validator_moniker,
            validator_mnemonic,
            chain_id,
            symbol_translator_file,
            ethereum_address,
            ethereum_private_key,
            relayerdb_path=None,
            keyring_backend=None,
            log_file=None,
            cwd=None,
    ):
        # Usage:
        #   ebrelayer init-witness [networkDescriptor] [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMnemonic] [flags]
        #
        # Examples:
        # ebrelayer init-witness 1 tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 mnemonic --chain-id=peggy
        extra_args = [] + \
            (["--relayerdb-path", relayerdb_path] if relayerdb_path else [])
        return self.__peggy2_init_common("init-witness", network_descriptor, tendermint_node, web3_provider,
            bridge_registry_contract_address, validator_mnemonic, chain_id=chain_id, node=tendermint_node,
            sign_with=validator_moniker, symbol_translator_file=symbol_translator_file,
            ethereum_address=ethereum_address, ethereum_private_key=ethereum_private_key,
            keyring_backend=keyring_backend, cwd=cwd, log_file=log_file,
            log_format="json", extra_args=extra_args)

    def __peggy2_init_common(self, init_what, network_descriptor, tendermint_node, web3_provider,
        bridge_registry_contract_address, validator_mnemonic, chain_id, node=None, keyring_backend=None,
        sign_with=None, symbol_translator_file=None, log_format=None, extra_args=None, ethereum_private_key=None,
        ethereum_address=None, cwd=None, log_file=None
    ):
        env = {}
        if ethereum_private_key:
            assert not ethereum_private_key.startswith("0x")
            env["ETHEREUM_PRIVATE_KEY"] = ethereum_private_key
        if ethereum_address:
            assert ethereum_address.startswith("0x")
            env["ETHEREUM_ADDRESS"] = ethereum_address
        env = env or None  # Avoid passing empty environment
        args = [
            self.binary,
            init_what,
            network_descriptor,
            tendermint_node,
            web3_provider,
            bridge_registry_contract_address,
            " ".join(validator_mnemonic),
            "--chain-id", chain_id
        ] + \
            (extra_args if extra_args else []) + \
            (["--node", node] if node else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend else []) + \
            (["--from", sign_with] if sign_with else []) + \
            (["--symbol-translator-file", symbol_translator_file] if symbol_translator_file else []) + \
            (["--log_format", log_format] if log_format else [])
        return self.cmd.popen(args, env=env, cwd=cwd, log_file=log_file)

    # Legacy stuff - pre-peggy2
    # Called from IntegrationContext
    def init(self, tendermind_node, web3_provider, bridge_registry_contract_address, validator_moniker,
        validator_mnemonic, chain_id, ethereum_private_key=None, ethereum_address=None, gas=None, gas_prices=None,
        node=None, keyring_backend=None, sign_with=None, symbol_translator_file=None, relayerdb_path=None,
        cwd=None, log_file=None
    ):
        env = {}
        if ethereum_private_key:
            assert not ethereum_private_key.startswith("0x")
            env["ETHEREUM_PRIVATE_KEY"] = ethereum_private_key
        if ethereum_address:
            assert ethereum_address.startswith("0x")
            env["ETHEREUM_ADDRESS"] = ethereum_address
        env = env or None  # Avoid passing empty environment
        args = [self.binary, "init", tendermind_node, web3_provider, bridge_registry_contract_address,
            validator_moniker, " ".join(validator_mnemonic), "--chain-id={}".format(chain_id)] + \
            (["--gas", str(gas)] if gas is not None else []) + \
            (["--gas-prices", sif_format_amount(*gas_prices)] if gas_prices is not None else []) + \
            (["--node", node] if node is not None else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend is not None else []) + \
            (["--from", sign_with] if sign_with is not None else []) + \
            (["--symbol-translator-file", symbol_translator_file] if symbol_translator_file else []) + \
            (["--relayerdb-path", relayerdb_path] if relayerdb_path else [])
        return self.cmd.popen(args, env=env, cwd=cwd, log_file=log_file)
