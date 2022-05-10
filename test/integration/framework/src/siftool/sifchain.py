import base64
import json
import time
import grpc
import re
import web3
from typing import Mapping, Any, Tuple
from siftool import command, cosmos, eth
from siftool.common import *

def sifchain_denom_hash(network_descriptor: int, token_contract_address: eth.Address) -> str:
    assert on_peggy2_branch
    assert token_contract_address.startswith("0x")
    assert type(network_descriptor) == int
    assert network_descriptor in range(1, 10000)
    denom = f"sifBridge{network_descriptor:04d}{token_contract_address.lower()}"
    return denom

def sifchain_denom_hash_to_token_contract_address(token_hash: str) -> Tuple[int, eth.Address]:
    m = re.match("^sifBridge(\\d{4})0x([0-9a-fA-F]{40})$", token_hash)
    if not m:
        raise Exception("Invalid sifchain denom '{}'".format(token_hash))
    network_descriptor = int(m[1])
    token_address = web3.Web3.toChecksumAddress(m[2])
    return network_descriptor, token_address

# Deprecated
def balance_delta(balances1: cosmos.Balance, balances2: cosmos.Balance) -> cosmos.Balance:
    return cosmos.balance_sub(balances2, balances1)

# Deprecated
def balance_zero(balances: cosmos.Balance) -> bool:
    return cosmos.balance_zero(balances)

def is_cosmos_native_denom(denom: str) -> bool:
    """Returns true if denom is a native cosmos token (Rowan, ibc)
    that was not imported using Peggy"""
    return not str.startswith(denom, "sifBridge")

def import_generated_protobuf_sources():
    import cosmos.tx.v1beta1.service_pb2 as cosmos_pb
    import cosmos.tx.v1beta1.service_pb2_grpc as cosmos_pb_grpc


class Sifnoded:
    def __init__(self, cmd, home: str = None):
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

    def set_genesis_token_registry_admin(self, address):
        self.sifnoded_exec(["set-genesis-token-registry-admin", address], sifnoded_home=self.home)

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
        args = ["tx", "tokenregistry", "set-registry", registry_path, "--gas-prices", sif_format_amount(*gas_prices),
            "--gas-adjustment", str(gas_adjustment), "--from", from_account, "--chain-id", chain_id, "--output", "json",
            "--yes"]
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

    def peggy2_update_consensus_needed(self, admin_account_address, hardhat_chain_id, chain_id):
        consensus_needed = "49"
        args = ["tx", "ethbridge", "update-consensus-needed", admin_account_address, str(hardhat_chain_id),
            consensus_needed, "--from", admin_account_address, "--chain-id", chain_id, "--gas-prices",
            "0.5rowan", "--gas-adjustment", "1.5", "-y"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return res

    def sifnoded_start(self, tcp_url=None, minimum_gas_prices=None, log_format_json=False, log_file=None):
        sifnoded_exec_args = self.build_start_cmd(tcp_url=tcp_url, minimum_gas_prices=minimum_gas_prices,
            log_format_json=log_format_json)
        return self.cmd.spawn_asynchronous_process(sifnoded_exec_args, log_file=log_file)

    def build_start_cmd(self, tcp_url: str = None, minimum_gas_prices=None, log_format_json=False):
        args = [self.binary, "start", "--trace"] + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else []) + \
            (["--rpc.laddr", tcp_url] if tcp_url else []) + \
            (["--log_level", "debug"] if log_format_json else []) + \
            (["--log_format", "json"] if log_format_json else []) + \
            (["--home", self.home] if self.home else [])
        return command.buildcmd(args)

    def sifnoded_exec(self, args, sifnoded_home=None, keyring_backend=None, stdin=None, cwd=None):
        args = [self.binary] + args + \
            (["--home", sifnoded_home] if sifnoded_home else []) + \
            (["--keyring-backend", keyring_backend] if keyring_backend else [])
        res = self.cmd.execst(args, stdin=stdin, cwd=cwd)
        return res

    def _rpc_get(self, host, port, relative_url):
        url = "http://{}:{}/{}".format(host, port, relative_url)
        return json.loads(http_get(url).decode("UTF-8"))

    def get_status(self, host, port):
        return self._rpc_get(host, port, "node_info")

    def wait_for_last_transaction_to_be_mined(self, count=1):
        # TODO return int(self._rpc_get(host, port, abci_info)["response"]["last_block_height"])
        def latest_block_height():
            args = ["status"]  # TODO --node
            return int(json.loads(stderr(self.sifnoded_exec(args)))["SyncInfo"]["latest_block_height"])
        initial_block = latest_block_height()
        while latest_block_height() < initial_block + count:
            time.sleep(1)

    def wait_up(self, host, port):
        while True:
            from urllib.error import URLError
            try:
                return self.get_status(host, port)
            except URLError:
                time.sleep(1)


# Refactoring in progress
class SifnodeClient:
    def __init__(self, cmd, node=None, home=None, chain_id=None, grpc_port=None):
        self.cmd = cmd
        self.binary = "sifnoded"
        self.node = node
        self.home = home
        self.chain_id = chain_id
        self.grpc_port = grpc_port

    def query_account(self, sif_addr):
        result = json.loads(stdout(self.sifnoded_exec(["query", "account", sif_addr, "--output", "json"])))
        return result

    def send_from_sifchain_to_ethereum(self, from_sif_addr: cosmos.Address, to_eth_addr: str, amount: int, denom: str,
        generate_only: bool = False
    ) -> Mapping:
        """ Sends ETH from Sifchain to Ethereum (burn) """
        assert on_peggy2_branch, "Only for Peggy2.0"
        assert self.ctx.eth
        eth = self.ctx.eth

        direction = "lock" if is_cosmos_native_denom(denom) else "burn"
        cross_chain_ceth_fee = eth.cross_chain_fee_base * eth.cross_chain_burn_fee  # TODO
        args = ["tx", "ethbridge", direction, from_sif_addr, to_eth_addr, str(amount), denom, str(cross_chain_ceth_fee),
                "--network-descriptor", str(eth.ethereum_network_descriptor),  # Mandatory
                "--from", from_sif_addr,  # Mandatory, either name from keyring or address
                "--output", "json",
                "-y"
            ] + \
            (["--generate-only"] if generate_only else []) + \
            self._gas_prices_args() + \
            self._home_args() + \
            self._chain_id_and_node_args() + \
            (self._keyring_backend_args() if not generate_only else [])
        res = self.sifnoded_exec(args)
        result = json.loads(stdout(res))
        if not generate_only:
            assert "failed to execute message" not in result["raw_log"]
        return result

    def send_from_sifchain_to_ethereum_grpc(self, from_sif_addr: cosmos.Address, to_eth_addr: str, amount: int,
        denom: str
    ):
        tx = self.send_from_sifchain_to_ethereum(from_sif_addr, to_eth_addr, amount, denom, generate_only=True)
        signed_tx = self.sign_transaction(tx, from_sif_addr)
        encoded_tx = self.encode_transaction(signed_tx)
        result = self.broadcast_tx(encoded_tx)
        return result

    def sign_transaction(self, tx: Mapping, from_sif_addr: cosmos.Address, sequence: int = None,
        account_number: int = None
    ) -> Mapping:
        tmp_tx_file = self.cmd.mktempfile()
        assert (sequence is not None) == (account_number is not None)  # We need either both or none
        try:
            self.cmd.write_text_file(tmp_tx_file, json.dumps(tx))
            args = ["tx", "sign", tmp_tx_file, "--from", from_sif_addr] + \
                (["--sequence", str(sequence), "--offline", "--account-number", str(account_number)] if sequence else []) + \
                self._home_args() + \
                self._chain_id_and_node_args() + \
                self._keyring_backend_args()
            res = self.sifnoded_exec(args)
            signed_tx = json.loads(stderr(res))
            return signed_tx
        finally:
            self.cmd.rm(tmp_tx_file)

    def encode_transaction(self, tx: Mapping[str, Any]) -> bytes:
        tmp_file = self.cmd.mktempfile()
        try:
            self.cmd.write_text_file(tmp_file, json.dumps(tx))
            res = self.sifnoded_exec(["tx", "encode", tmp_file])
            encoded_tx = base64.b64decode(stdout(res))
            return encoded_tx
        finally:
            self.cmd.rm(tmp_file)

    def open_grpc_channel(self) -> grpc.Channel:
        # See https://docs.cosmos.network/v0.44/core/proto-docs.html
        # See https://docs.cosmos.network/v0.44/core/grpc_rest.html
        # See https://app.swaggerhub.com/apis/Ivan-Verchenko/sifnode-swagger-api/1.1.1
        # See https://raw.githubusercontent.com/Sifchain/sifchain-ui/develop/ui/core/swagger.yaml
        return grpc.insecure_channel("127.0.0.1:9090")

    def broadcast_tx(self, encoded_tx: bytes):
        import_generated_protobuf_sources()
        broadcast_mode = cosmos_pb.BROADCAST_MODE_ASYNC
        with self.open_grpc_channel() as channel:
            tx_stub = cosmos_pb_grpc.ServiceStub(channel)
            req = cosmos_pb.BroadcastTxRequest(tx_bytes=encoded_tx, mode=broadcast_mode)
            resp = tx_stub.BroadcastTx(req)
            return resp

    def _gas_prices_args(self):
        return ["--gas-prices", "0.5rowan", "--gas-adjustment", "1.5"]

    def _chain_id_and_node_args(self):
        return \
           (["--node", self.node] if self.node else []) + \
           (["--chain-id", self.chain_id] if self.chain_id else [])

    def _keyring_backend_args(self):
        keyring_backend = self.ctx.sifnode.keyring_backend
        return ["--keyring-backend", keyring_backend] if keyring_backend else []

    def _home_args(self):
        return ["--home", self.home] if self.home else []

    def sifnoded_exec(self, *args, **kwargs):
        return self.ctx.sifnode.sifnoded_exec(*args, **kwargs)


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
        keyring_dir=None, sign_with=None, symbol_translator_file=None, log_format=None, extra_args=None,
        ethereum_private_key=None, ethereum_address=None, home=None, cwd=None
    ):
        env = _env_for_ethereum_address_and_key(ethereum_address, ethereum_private_key)
        args = [
            self.binary,
            init_what,
            "--network-descriptor", str(network_descriptor),  # Network descriptor for the chain (9999)
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
            (["--home", home] if home else []) + \
            (["--keyring-dir", keyring_dir] if keyring_dir else []) + \
            (["--symbol-translator-file", symbol_translator_file] if symbol_translator_file else []) + \
            (["--log_format", log_format] if log_format else [])
        return command.buildcmd(args, env=env, cwd=cwd)

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
