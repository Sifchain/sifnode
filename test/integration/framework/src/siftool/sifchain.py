import base64
import contextlib
import json
import time
import grpc
import re
import toml
import web3  # TODO Remove dependency
from typing import Mapping, Any, Tuple, AnyStr
from siftool import command, cosmos, eth
from siftool.common import *


log = siftool_logger(__name__)


ROWAN = "rowan"
STAKE = "stake"
ROWAN_DECIMALS = 18
CETH = "ceth"  # Peggy1 only (Peggy2.0 uses denom hash)

# Sifchain public network endpoints
BETANET = {"node": "https://rpc.sifchain.finance", "chain_id": "sifchain-1"}
TESTNET = {"node": "https://rpc-testnet.sifchain.finance", "chain_id": "sifchain-testnet-1"}
DEVNET = {"node": "https://rpc-devnet.sifchain.finance", "chain_id": "sifchain-devnet-1"}

GasFees = Tuple[float, str]  # Special case of cosmos.Balance with only one denom and float amount  # TODO Rename to something more neutral such as Amount

# Format of a single "entry" of output of "sifnoded q tokenregistry generate"; used for "sifnoded tx tokenregistry register"
TokenRegistryParams = JsonDict
RewardsParams = JsonDict
LPPDParams = JsonDict


# Default ports
SIFNODED_DEFAULT_RPC_PORT = 26657  # In config/config.toml; used for --node, can be queried using curl http://...
SIFNODED_DEFAULT_P2P_PORT = 26656  # In config/config.toml
SIFNODED_DEFAULT_API_PORT = 1317  # In config/app.toml, section [api], disabled by default
SIFNODED_DEFAULT_GRPC_PORT = 9090  # In config/app.toml, section [grpc]
SIFNODED_DEFAULT_GRPC_WEB_PORT = 9091  # In config/app.toml, section [grpc-web]


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

def is_cosmos_native_denom(denom: str) -> bool:
    """Returns true if denom is a native cosmos token (Rowan, ibc)
    that was not imported using Peggy"""
    return not str.startswith(denom, "sifBridge")

def ondemand_import_generated_protobuf_sources():
    global cosmos_pb
    global cosmos_pb_grpc
    import cosmos.tx.v1beta1.service_pb2 as cosmos_pb
    import cosmos.tx.v1beta1.service_pb2_grpc as cosmos_pb_grpc

def mnemonic_to_address(cmd: command.Command, mnemonic: Iterable[str]):
    tmpdir = cmd.mktempdir()
    sifnode = Sifnoded(cmd, home=tmpdir)
    try:
       return sifnode.keys_add("tmp", mnemonic)["address"]
    finally:
        cmd.rmdir(tmpdir)

def sifnoded_parse_output_lines(stdout: str) -> Mapping:
    # TODO Some values are like '""'
    pat = re.compile("^(.*?): (.*)$")
    result = {}
    for line in stdout.splitlines():
        m = pat.match(line)
        result[m[1]] = m[2]
    return result

def format_pubkey(pubkey: JsonDict) -> str:
    return "{{\"@type\":\"{}\",\"key\":\"{}\"}}".format(pubkey["@type"], pubkey["key"])

def format_peer_address(node_id: str, hostname: str, p2p_port: int) -> str:
    return "{}@{}:{}".format(node_id, hostname, p2p_port)

def format_node_url(hostname: str, p2p_port: int) -> str:
    return "tcp://{}:{}".format(hostname, p2p_port)

def create_rewards_descriptor(rewards_period_id: str, start_block: int, end_block: int,
    multipliers: Iterable[Tuple[str, int]], allocation: int, reward_period_default_multiplier: float,
    reward_period_distribute: bool, reward_period_mod: int
) -> RewardsParams:
    return {
        "reward_period_id": rewards_period_id,
        "reward_period_start_block": start_block,
        "reward_period_end_block": end_block,
        "reward_period_allocation": str(allocation),
        "reward_period_pool_multipliers": [{
            "pool_multiplier_asset": denom,
            "multiplier": str(multiplier)
        } for denom, multiplier in multipliers],
        "reward_period_default_multiplier": str(reward_period_default_multiplier),
        "reward_period_distribute": reward_period_distribute,
        "reward_period_mod": reward_period_mod
    }

def create_lppd_params(start_block: int, end_block: int, rate: float, mod: int) -> LPPDParams:
    return {
        "distribution_period_block_rate": str(rate),
        "distribution_period_start_block": start_block,
        "distribution_period_end_block": end_block,
        "distribution_period_mod": mod
    }


class Sifnoded:
    def __init__(self, cmd, /, home: Optional[str] = None, node: Optional[str] = None, chain_id: Optional[str] = None):
        self.cmd = cmd
        self.binary = "sifnoded"
        self.home = home
        self.node = node
        self.chain_id = chain_id
        self.keyring_backend = "test"

        # Firing transactions with "sifnoded tx bank send" in rapid succession does not work. This is currently a
        # known limitation of Cosmos SDK, see https://github.com/cosmos/cosmos-sdk/issues/4186
        # Instead, we take advantage of batching multiple denoms to single account with single send command (amounts
        # separated by by comma: "sifnoded tx bank send ... 100denoma,100denomb,100denomc") and wait for destination
        # account to show changes for all denoms after each send. But also batches don't work reliably if they are too
        # big, so we limit the maximum batch size here.
        self.max_send_batch_size = 5

        self.broadcast_mode = None
        # self.sifnoded_burn_gas_cost = 16 * 10**10 * 393000  # see x/ethbridge/types/msgs.go for gas
        # self.sifnoded_lock_gas_cost = 16 * 10**10 * 393000
        self.get_balance_default_retries = 0

        # Defaults
        self.wait_for_balance_change_default_timeout = 90
        self.wait_for_balance_change_default_change_timeout = None
        self.wait_for_balance_change_default_polling_time = 2

    # Returns what looks like genesis file data
    def init(self, moniker):
        args = [self.binary, "init", moniker] + self._home_args() + self._chain_id_args()
        res = self.cmd.execst(args)
        return json.loads(stderr(res))

    def keys_list(self):
        args = ["keys", "list", "--output", "json"] + self._home_args() + self._keyring_backend_args()
        res = self.sifnoded_exec(args)
        return json.loads(stdout(res))

    def keys_show(self, name, bech=None):
        args = ["keys", "show", name] + \
            (["--bech", bech] if bech else []) + \
            self._home_args() + \
            self._keyring_backend_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def get_val_address(self, moniker) -> cosmos.BechAddress:
        args = ["keys", "show", "-a", "--bech", "val", moniker] + self._home_args() + self._keyring_backend_args()
        res = self.sifnoded_exec(args)
        expected = exactly_one(stdout_lines(res))
        result = exactly_one(self.keys_show(moniker, bech="val"))["address"]
        assert result == expected
        return result

    def keys_add(self, moniker: Optional[str] = None, mnemonic: Optional[Iterable[str]] = None) -> JsonDict:
        moniker = self.__fill_in_moniker(moniker)
        if mnemonic is None:
            args = ["keys", "add", moniker] + self._home_args() + self._keyring_backend_args()
            res = self.sifnoded_exec(args, stdin=["y"])
            _unused_mnemonic = stderr(res).splitlines()[-1].split(" ")
        else:
            args = ["keys", "add", moniker, "--recover"] + self._home_args() + self._keyring_backend_args()
            res = self.sifnoded_exec(args, stdin=[" ".join(mnemonic)])
        account = exactly_one(yaml_load(stdout(res)))
        return account

    def generate_mnemonic(self) -> List[str]:
        args = ["keys", "mnemonic"] + self._home_args() + self._keyring_backend_args()
        res = self.sifnoded_exec(args)
        return exactly_one(stderr(res).splitlines()).split(" ")

    def create_addr(self, moniker: Optional[str] = None, mnemonic: Optional[Iterable[str]] = None) -> cosmos.Address:
        return self.keys_add(moniker=moniker, mnemonic=mnemonic)["address"]

    def keys_add_multisig(self, moniker: Optional[str], signers: Iterable[cosmos.KeyName], multisig_threshold: int):
        moniker = self.__fill_in_moniker(moniker)
        args = ["keys", "add", moniker, "--multisig", ",".join(signers), "--multisig-threshold",
            str(multisig_threshold)] + self._home_args() + self._keyring_backend_args()
        res = self.sifnoded_exec(args)
        account = exactly_one(yaml_load(stdout(res)))
        return account

    def __fill_in_moniker(self, moniker):
        return moniker if moniker else "temp-{}".format(random_string(10))

    def keys_delete(self, name: str):
        self.cmd.execst(["sifnoded", "keys", "delete", name] + self._home_args() + self._keyring_backend_args(),
            stdin=["y"], check_exit=False)

    def query_account(self, addr: cosmos.Address) -> JsonDict:
        args = ["query", "auth", "account", addr, "--output", "json"] + self._node_args() + self._chain_id_args()
        res = self.sifnoded_exec(args)
        return json.loads(stdout(res))

    def get_acct_seq(self, addr: cosmos.Address) -> Tuple[int, int]:
        account = self.query_account(addr)
        account_number = account["account_number"]
        account_sequence = int(account["sequence"])
        return account_number, account_sequence

    def add_genesis_account(self, sifnodeadmin_addr: cosmos.Address, tokens: cosmos.Balance):
        tokens_str = cosmos.balance_format(tokens)
        self.sifnoded_exec(["add-genesis-account", sifnodeadmin_addr, tokens_str] + self._home_args() + self._keyring_backend_args())

    # TODO Obsolete
    def add_genesis_account_directly_to_existing_genesis_json(self,
        extra_balances: Mapping[cosmos.Address, cosmos.Balance]
    ):
        genesis = self.load_genesis_json()
        self.add_accounts_to_existing_genesis(genesis, extra_balances)
        self.save_genesis_json(genesis)

    def add_accounts_to_existing_genesis(self, genesis: JsonDict, extra_balances: Mapping[cosmos.Address, cosmos.Balance]):
        bank = genesis["app_state"]["bank"]
        # genesis.json uses a bit different structure for balances so we need to convert to and from our balances.
        # Whatever is in extra_balances will be added to the existing amounts.
        # We must also update supply which must be the sum of all balances. We assume that it initially already is.
        # Cosmos SDK wants coins to be sorted or it will panic during chain initialization.
        balances = {b["address"]: {c["denom"]: int(c["amount"]) for c in b["coins"]} for b in bank["balances"]}
        supply = {b["denom"]: int(b["amount"]) for b in bank["supply"]}
        accounts = genesis["app_state"]["auth"]["accounts"]
        for addr, bal in extra_balances.items():
            b = cosmos.balance_add(balances.get(addr, {}), bal)
            balances[addr] = b
            supply = cosmos.balance_add(supply, bal)
        accounts.extend([{
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": addr,
          "pub_key": None,
          "account_number": "0",
          "sequence": "0"
        } for addr in set(balances).difference(set(x["address"] for x in accounts))])
        bank["balances"] = [{"address": a, "coins": [{"denom": d, "amount": str(c[d])} for d in sorted(c)]} for a, c in balances.items()]
        bank["supply"] = [{"denom": d, "amount": str(supply[d])} for d in sorted(supply)]

    def load_genesis_json(self) -> JsonDict:
        genesis_json_path = os.path.join(self.get_effective_home(), "config", "genesis.json")
        return json.loads(self.cmd.read_text_file(genesis_json_path))

    def save_genesis_json(self, genesis: JsonDict):
        genesis_json_path = os.path.join(self.get_effective_home(), "config", "genesis.json")
        self.cmd.write_text_file(genesis_json_path, json.dumps(genesis))

    def load_app_toml(self) -> JsonDict:
        app_toml_path = os.path.join(self.get_effective_home(), "config", "app.toml")
        with open(app_toml_path, "r") as app_toml_file:
            return toml.load(app_toml_file)

    def save_app_toml(self, data: JsonDict):
        app_toml_path = os.path.join(self.get_effective_home(), "config", "app.toml")
        with open(app_toml_path, "w") as app_toml_file:
            app_toml_file.write(toml.dumps(data))

    def load_config_toml(self) -> JsonDict:
        config_toml_path = os.path.join(self.get_effective_home(), "config", "config.toml")
        with open(config_toml_path, "r") as config_toml_file:
            return toml.load(config_toml_file)

    def save_config_toml(self, data: JsonDict):
        config_toml_path = os.path.join(self.get_effective_home(), "config", "config.toml")
        with open(config_toml_path, "w") as config_toml_file:
            config_toml_file.write(toml.dumps(data))

    def enable_rpc_port(self):
        app_toml = self.load_app_toml()
        app_toml["api"]["enable"] = True
        app_toml["api"]["address"] = format_node_url(ANY_ADDR, SIFNODED_DEFAULT_API_PORT)
        self.save_app_toml(app_toml)

    def get_effective_home(self) -> str:
        return self.home if self.home is not None else self.cmd.get_user_home(".sifnoded")

    def add_genesis_clp_admin(self, address: cosmos.Address):
        args = ["add-genesis-clp-admin", address] + self._home_args() + self._keyring_backend_args()
        self.sifnoded_exec(args)

    # Modifies genesis.json and adds the address to .oracle.address_whitelist array.
    def add_genesis_validators(self, address: cosmos.BechAddress):
        args = ["add-genesis-validators", address] + self._home_args() + self._keyring_backend_args()
        res = self.sifnoded_exec(args)
        return res

    # At the moment only on future/peggy2 branch, called from PeggyEnvironment
    def add_genesis_validators_peggy(self, evm_network_descriptor: int, valoper: cosmos.BechAddress, validator_power: int):
        assert on_peggy2_branch
        args = ["add-genesis-validators", str(evm_network_descriptor), valoper, str(validator_power)] + \
            self._home_args()
        self.sifnoded_exec(args)

    def set_genesis_oracle_admin(self, address):
        self.sifnoded_exec(["set-genesis-oracle-admin", address] + self._home_args() + self._keyring_backend_args())

    def set_genesis_token_registry_admin(self, address):
        self.sifnoded_exec(["set-genesis-token-registry-admin", address] + self._home_args())

    def set_genesis_whitelister_admin(self, address):
        self.sifnoded_exec(["set-genesis-whitelister-admin", address] + self._home_args() + self._keyring_backend_args())

    def set_gen_denom_whitelist(self, denom_whitelist_file):
        self.sifnoded_exec(["set-gen-denom-whitelist", denom_whitelist_file] + self._home_args())

    def tendermint_show_node_id(self) -> str:
        args = ["tendermint", "show-node-id"] + self._home_args()
        res = self.sifnoded_exec(args)
        return exactly_one(stdout(res).splitlines())

    def tendermint_show_validator(self):
        args = ["tendermint", "show-validator"] + self._home_args()
        res = self.sifnoded_exec(args)
        return json.loads(stdout(res))

    # self.node ("--node") should point to existing validator (i.e. node 0) which must be up.
    # The balance of from_acct (from node 0's perspective) must be greater than the staking amount.
    # amount must be a single denom, and must denominated as per config/app_state.toml::staking.params.bond_denom
    # pubkey must be from "tendermint show validator", NOT from "keys add"
    def staking_create_validator(self, amount: cosmos.Balance, pubkey: JsonDict, moniker: str, commission_rate: float,
        commission_max_rate: float, commission_max_change_rate: float, min_self_delegation: int,
        from_acct: cosmos.Address, broadcast_mode: Optional[str] = None
    ):
        assert len(amount) == 1  # Maybe not? We haven't seen staking with more than one denom yet though.
        assert cosmos.balance_exceeds(self.get_balance(from_acct), amount)
        assert pubkey["@type"] == "/cosmos.crypto.ed25519.PubKey"
        args = ["tx", "staking", "create-validator", "--amount", cosmos.balance_format(amount), "--pubkey",
            format_pubkey(pubkey), "--moniker", moniker, "--commission-rate", str(commission_rate),
            "--commission-max-rate", str(commission_max_rate), "--commission-max-change-rate",
            str(commission_max_change_rate), "--min-self-delegation", str(min_self_delegation), "--from", from_acct] + \
            self._home_args() + self._chain_id_args() + self._node_args() + self._keyring_backend_args() + \
            self._fees_args() + self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def query_staking_validators(self) -> JsonObj:
        args = ["query", "staking", "validators", "--output", "json"] + self._home_args() + self._node_args()
        res = self.sifnoded_exec(args)
        res = json.loads(stdout(res))
        # TODO Implement paging (here we don't have "block" so we can't use self._paged_read())
        #      But probably have to go the same way as self.get_balance()
        assert res["pagination"]["next_key"] is None, "TODO: Implement pagination"
        return res["validators"]

    # See scripts/ibc/tokenregistration for more information and examples.
    # JSON file can be generated with "sifnoded q tokenregistry generate"
    def create_tokenregistry_entry(self, symbol: str, sifchain_symbol: str, decimals: int,
        permissions: Iterable[str] = None
    ) -> TokenRegistryParams:
        permissions = permissions if permissions is not None else ["CLP", "IBCEXPORT", "IBCIMPORT"]
        upper_symbol = symbol.upper()  # Like "USDT"
        return {
            "decimals": str(decimals),
            "denom": sifchain_symbol,
            "base_denom": sifchain_symbol,
            "path": "",
            "ibc_channel_id": "",
            "ibc_counterparty_channel_id": "",
            "display_name": upper_symbol,
            "display_symbol": "",
            "network": "",
            "address": "",
            "external_symbol": upper_symbol,
            "transfer_limit": "",
            "permissions": list(permissions),
            "unit_denom": "",
            "ibc_counterparty_denom": "",
            "ibc_counterparty_chain_id": "",
        }

    # from_sif_addr has to be the address which was used at genesis time for "set-genesis-whitelister-admin".
    # You need to have its private key in the test keyring.
    # This is needed when creating pools for the token or when doing IBC transfers.
    # If you are calling this for several tokens, you need to call it synchronously
    # (i.e. wait_for_current_block_to_be_mined(), or broadcast_mode="block"). Otherwise this will silently fail.
    # This is used in test_many_pools_and_liquidity_providers.py
    def token_registry_register(self, entry: TokenRegistryParams, from_sif_addr: cosmos.Address, broadcast_mode=None
    ) -> JsonDict:
        # Check that we have the private key in test keyring. This will throw an exception if we don't.
        assert self.keys_show(from_sif_addr)
        # This command requires a single TokenRegistryEntry, even though the JSON file has "entries" as a list.
        # If not: "Error: exactly one token entry must be specified in input file"
        token_data = {"entries": [entry]}
        with self._with_temp_json_file(token_data) as tmp_registry_json:
            args = ["tx", "tokenregistry", "register", tmp_registry_json, "--from", from_sif_addr, "--output",
                "json"] + self._home_args() + self._keyring_backend_args() + self._chain_id_args() + \
                self._node_args() + self._fees_args() + self._broadcast_mode_args(broadcast_mode=broadcast_mode) + \
                self._yes_args()
            res = self.sifnoded_exec(args)
            res = json.loads(stdout(res))
            # Example of successful output: {"height":"196804","txhash":"C8252E77BCD441A005666A4F3D76C99BD35F9CB49AA1BE44CBE2FFCC6AD6ADF4","codespace":"","code":0,"data":"0A270A252F7369666E6F64652E746F6B656E72656769737472792E76312E4D73675265676973746572","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/sifnode.tokenregistry.v1.MsgRegister\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"/sifnode.tokenregistry.v1.MsgRegister"}]}]}],"info":"","gas_wanted":"200000","gas_used":"115149","tx":null,"timestamp":""}
            if res["raw_log"].startswith("signature verification failed"):
                raise Exception(res["raw_log"])
            if res["raw_log"].startswith("failed to execute message"):
                raise Exception(res["raw_log"])
            return res

    def query_tokenregistry_entries(self):
        args = ["query", "tokenregistry", "entries"] + self._node_args() + self._chain_id_args()
        res = self.sifnoded_exec(args)
        return json.loads(stdout(res))["entries"]

    # Creates file config/gentx/gentx-*.json
    def gentx(self, name: str, stake: cosmos.Balance):
        # TODO Make chain_id an attribute
        args = ["gentx", name, cosmos.balance_format(stake)] + self._home_args() + self._keyring_backend_args() + \
            self._chain_id_args()
        res = self.sifnoded_exec(args)
        return exactly_one(stderr(res).splitlines())

    # Modifies genesis.json and adds .genutil.gen_txs (presumably from config/gentx/gentx-*.json)
    def collect_gentx(self) -> JsonDict:
        args = ["collect-gentxs"] + self._home_args()  # Must not use --keyring-backend
        res = self.sifnoded_exec(args)
        return json.loads(stderr(res))

    def validate_genesis(self):
        args = ["validate-genesis"] + self._home_args()  # Must not use --keyring-backend
        res = self.sifnoded_exec(args)
        res = exactly_one(stdout(res).splitlines())
        assert res.endswith(" is a valid genesis file")

    # At the moment only on future/peggy2 branch, called from PeggyEnvironment
    # This was split from init_common
    def peggy2_add_account(self, name: str, tokens: cosmos.Balance, is_admin: bool = False):
        # TODO Peggy2 devenv feed "yes\nyes" into standard input, we only have "y\n"
        account = self.keys_add(name)
        account_address = account["address"]

        self.add_genesis_account(account_address, tokens)
        if is_admin:
            self.set_genesis_oracle_admin(account_address)
            self.set_genesis_whitelister_admin(account_address)
        return account_address

    def peggy2_add_relayer_witness_account(self, name: str, tokens: cosmos.Balance, evm_network_descriptor: int,
        validator_power: int, denom_whitelist_file: str
    ):
        account_address = self.peggy2_add_account(name, tokens)  # Note: is_admin=False
        # Whitelist relayer/witness account
        valoper = self.get_val_address(name)
        self.set_gen_denom_whitelist(denom_whitelist_file)
        self.add_genesis_validators_peggy(evm_network_descriptor, valoper, validator_power)
        return account_address

    def peggy2_token_registry_set_registry(self, registry_path: str, gas_prices: GasFees, gas_adjustment: float,
        from_account: cosmos.Address, chain_id: str
    ) -> List[Mapping[str, Any]]:
        args = ["tx", "tokenregistry", "set-registry", registry_path, "--gas-prices", sif_format_amount(*gas_prices),
            "--gas-adjustment", str(gas_adjustment), "--from", from_account, "--chain-id", chain_id, "--output", "json"] + \
            self._home_args() + self._keyring_backend_args() + self._yes_args()
        res = self.sifnoded_exec(args)
        return [json.loads(x) for x in stdout(res).splitlines()]

    def peggy2_set_cross_chain_fee(self, admin_account_address, network_id, ethereum_cross_chain_fee_token,
        cross_chain_fee_base, cross_chain_lock_fee, cross_chain_burn_fee, admin_account_name, gas_prices,
        gas_adjustment
    ):
        args = ["tx", "ethbridge", "set-cross-chain-fee", str(network_id),
            ethereum_cross_chain_fee_token, str(cross_chain_fee_base), str(cross_chain_lock_fee),
            str(cross_chain_burn_fee), "--from", admin_account_name, "--gas-prices", sif_format_amount(*gas_prices),
            "--gas-adjustment", str(gas_adjustment)] + self._home_args() + self._keyring_backend_args() + \
            self._chain_id_args() + self._yes_args()
        return self.sifnoded_exec(args)

    def peggy2_update_consensus_needed(self, admin_account_address, hardhat_chain_id, consensus_needed):
        args = ["tx", "ethbridge", "update-consensus-needed", str(hardhat_chain_id),
            str(consensus_needed), "--from", admin_account_address] + self._home_args() + \
               self._keyring_backend_args() + self._gas_prices_args() + self._chain_id_args() + self._yes_args()
        return self.sifnoded_exec(args)

    # TODO Rename tcp_url to rpc_laddr + remove dependency on self.node
    def sifnoded_start(self, tcp_url: Optional[str] = None, minimum_gas_prices: Optional[GasFees] = None,
        log_format_json: bool = False, log_file: Optional[IO] = None, log_level: Optional[str] = None,
        trace: bool = False, p2p_laddr: Optional[str] = None, grpc_address: Optional[str] = None,
        grpc_web_address: Optional[str] = None, address: Optional[str] = None
    ):
        sifnoded_exec_args = self.build_start_cmd(tcp_url=tcp_url, p2p_laddr=p2p_laddr, grpc_address=grpc_address,
            grpc_web_address=grpc_web_address, address=address, minimum_gas_prices=minimum_gas_prices,
            log_format_json=log_format_json, log_level=log_level, trace=trace)
        return self.cmd.spawn_asynchronous_process(sifnoded_exec_args, log_file=log_file)

    # TODO Rename tcp_url to rpc_laddr + remove dependency on self.node
    def build_start_cmd(self, tcp_url: Optional[str] = None, p2p_laddr: Optional[str] = None,
        grpc_address: Optional[str] = None, grpc_web_address: Optional[str] = None, address: Optional[str] = None,
        minimum_gas_prices: Optional[GasFees] = None, log_format_json: bool = False, log_level: Optional[str] = None,
        trace: bool = False
    ):
        args = [self.binary, "start"] + \
            (["--trace"] if trace else []) + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else []) + \
            (["--rpc.laddr", tcp_url] if tcp_url else []) + \
            (["--p2p.laddr", p2p_laddr] if p2p_laddr else []) + \
            (["--grpc.address", grpc_address] if grpc_address else []) + \
            (["--grpc-web.address", grpc_web_address] if grpc_web_address else []) + \
            (["--address", address] if address else []) + \
            (["--log_level", log_level] if log_level else []) + \
            (["--log_format", "json"] if log_format_json else []) + \
            self._home_args()
        return command.buildcmd(args)

    def send(self, from_sif_addr: cosmos.Address, to_sif_addr: cosmos.Address, amounts: cosmos.Balance) -> Mapping:
        amounts = cosmos.balance_normalize(amounts)
        assert len(amounts) > 0
        # TODO Implement batching (factor it out of inflate_token and put here)
        if self.max_send_batch_size > 0:
            assert len(amounts) <= self.max_send_batch_size, \
                "Currently only up to {} balances can be send at the same time reliably.".format(self.max_send_batch_size)
        amounts_string = cosmos.balance_format(amounts)
        args = ["tx", "bank", "send", from_sif_addr, to_sif_addr, amounts_string, "--output", "json"] + \
            self._home_args() + self._keyring_backend_args() + self._chain_id_args() + self._node_args() + \
            self._fees_args() + self._broadcast_mode_args() + self._yes_args()
        res = self.sifnoded_exec(args)
        retval = json.loads(stdout(res))
        raw_log = retval["raw_log"]
        for bad_thing in ["insufficient funds", "signature verification failed"]:
            if bad_thing in raw_log:
                raise Exception(raw_log)
        return retval

    def get_balance(self, sif_addr: cosmos.Address, height: Optional[int] = None,
        disable_log: bool = False, retries_on_error: Optional[int] = None, delay_on_error: int = 3
    ) -> cosmos.Balance:
        retries_on_error = retries_on_error if retries_on_error is not None else self.get_balance_default_retries
        all_balances = {}
        # The actual limit might be capped to a lower value (100), in this case everything will still work but we'll get
        # fewer results
        desired_page_size = 5000
        page_key = None
        while True:
            args = ["query", "bank", "balances", sif_addr, "--output", "json"] + \
                (["--height", str(height)] if height is not None else []) + \
                (["--limit", str(desired_page_size)] if desired_page_size is not None else []) + \
                (["--page-key", page_key] if page_key is not None else []) + \
                self._home_args() + \
                self._chain_id_args() + \
                self._node_args()
            retries_left = retries_on_error
            while True:
                try:
                    res = self.sifnoded_exec(args, disable_log=disable_log)
                    break
                except Exception as e:
                    if retries_left == 0:
                        raise e
                    else:
                        retries_left -= 1
                        log.error("Error reading balances, retries left: {}".format(retries_left))
                        time.sleep(delay_on_error)
            res = json.loads(stdout(res))
            balances = res["balances"]
            next_key = res["pagination"]["next_key"]
            if next_key is not None:
                if height is None:
                    # There are more results than fit on a page. To ensure we get all balances as a consistent
                    # snapshot, retry with "--height" fised to the current block. This wastes one request.
                    # We could optimize this by starting with explicit "--height" in the first place, but the current
                    # assumption is that most of results will fit on one page and that this will be faster without
                    # "--height".
                    height = self.get_current_block()
                    log.debug("Large balance result, switching to paged mode using height of {}.".format(height))
                    continue
                page_key = base64.b64decode(next_key).decode("UTF-8")
            for bal in balances:
                denom, amount = bal["denom"], int(bal["amount"])
                assert denom not in all_balances
                all_balances[denom] = amount
            log.debug("Read {} balances, all={}, first='{}', next_key={}".format(len(balances), len(all_balances),
                balances[0]["denom"] if len(balances) > 0 else None, next_key))
            if next_key is None:
                break
        return all_balances

    # Unless timed out, this function will exit:
    # - if min_changes are given: when changes are greater.
    # - if expected_balance is given: when balances are equal to that.
    # - if neither min_changes nor expected_balance are given: when anything changes.
    # You cannot use min_changes and expected_balance at the same time.
    def wait_for_balance_change(self, sif_addr: cosmos.Address, old_balance: cosmos.Balance,
        min_changes: cosmos.CompatBalance = None, expected_balance: cosmos.CompatBalance = None,
        polling_time: Optional[int] = None, timeout: Optional[int] = None, change_timeout: Optional[int] = None,
        disable_log: bool = True
    ) -> cosmos.Balance:
        polling_time = polling_time if polling_time is not None else self.wait_for_balance_change_default_polling_time
        timeout = timeout if timeout is not None else self.wait_for_balance_change_default_timeout
        change_timeout = change_timeout if change_timeout is not None else self.wait_for_balance_change_default_change_timeout
        assert (min_changes is None) or (expected_balance is None), "Cannot use both min_changes and expected_balance"
        log.debug("Waiting for balance to change for account {}...".format(sif_addr))
        min_changes = None if min_changes is None else cosmos.balance_normalize(min_changes)
        expected_balance = None if expected_balance is None else cosmos.balance_normalize(expected_balance)
        start_time = time.time()
        last_change_time = None
        last_changed_balance = None
        while True:
            new_balance = self.get_balance(sif_addr, disable_log=disable_log)
            delta = cosmos.balance_sub(new_balance, old_balance)
            if expected_balance is not None:
                should_return = cosmos.balance_equal(expected_balance, new_balance)
            elif min_changes is not None:
                should_return = cosmos.balance_exceeds(delta, min_changes)
            else:
                should_return = not cosmos.balance_zero(delta)
            if should_return:
                return new_balance
            now = time.time()
            if (timeout is not None) and (timeout > 0) and (now - start_time > timeout):
                raise Exception("Timeout waiting for sif balance to change ({}s)".format(timeout))
            if last_change_time is None:
                last_changed_balance = new_balance
                last_change_time = now
            else:
                delta = cosmos.balance_sub(new_balance, last_changed_balance)
                if not cosmos.balance_zero(delta):
                    last_changed_balance = new_balance
                    last_change_time = now
                    log.debug("New state detected ({} denoms changed)".format(len(delta)))
                if (change_timeout is not None) and (change_timeout > 0) and (now - last_change_time > change_timeout):
                    raise Exception("Timeout waiting for sif balance to change ({}s)".format(change_timeout))
            time.sleep(polling_time)

    # TODO Refactor - consolidate with test_inflate_tokens.py
    def send_batch(self, from_addr: cosmos.Address, to_addr: cosmos.Address, amounts: cosmos.Balance):
        to_go = list(amounts.keys())
        while to_go:
            cnt = len(to_go)
            if (self.max_send_batch_size > 0) and (cnt > self.max_send_batch_size):
                cnt = self.max_send_batch_size
            batch_denoms = to_go[:cnt]
            to_go = to_go[cnt:]
            batch_balance = {denom: amounts.get(denom, 0) for denom in batch_denoms}
            balance_before = self.get_balance(to_addr)
            self.send(from_addr, to_addr, batch_balance)
            expected_balance = cosmos.balance_add(balance_before, batch_balance)
            self.wait_for_balance_change(to_addr, balance_before, expected_balance=expected_balance)

    def get_current_block(self):
        return int(self.status()["SyncInfo"]["latest_block_height"])

    # TODO Deduplicate
    # TODO The /node_info URL does not exist any more, use /status?
    def get_status(self, host, port):
        return self._rpc_get(host, port, "node_info")

    # TODO Deduplicate
    def status(self):
        args = ["status"] + self._node_args()
        res = self.sifnoded_exec(args)
        return json.loads(stderr(res))

    def query_block(self, block: Optional[int] = None) -> JsonDict:
        args = ["query", "block"] + \
            ([str(block)] if block is not None else []) + self._node_args()
        res = self.sifnoded_exec(args)
        return json.loads(stdout(res))

    def query_pools(self, height: Optional[int] = None) -> List[JsonDict]:
        return self._paged_read(["query", "clp", "pools"], "pools", height)

    def query_clp_liquidity_providers(self, denom: str, height: Optional[int] = None) -> List[JsonDict]:
        # Note: this paged result is slightly different than `query bank balances`. Here we always get "height"
        return self._paged_read(["query", "clp", "lplist", denom], "liquidity_providers", height)

    def _paged_read(self, base_args: List[str], result_key: str, height: Optional[int]) -> List[JsonDict]:
        # Note: this paged result is slightly different than `query bank balances`. Here we always get "height"
        all_results = []
        page_key = None
        while True:
            args = base_args + \
                (["--height", str(height)] if height is not None else []) + \
                (["--page-key", page_key] if page_key is not None else []) + \
               self._chain_id_args() + self._node_args()
            res = self.sifnoded_exec(args)
            res = yaml_load(stdout(res))
            all_results.extend(res[result_key])
            next_key = res["pagination"]["next_key"]
            if next_key is None:
                break
            page_key = base64.b64decode(next_key)
            if height is None:
                height = int(res["height"])
        return all_results

    def tx_clp_create_pool(self, from_addr: cosmos.Address, symbol: str, native_amount: int, external_amount: int
    ) -> JsonDict:
        # For more examples see ticket #2470, e.g.
        args = ["tx", "clp", "create-pool", "--from", from_addr, "--symbol", symbol, "--nativeAmount",
            str(native_amount), "--externalAmount", str(external_amount)] + self._chain_id_args() + \
            self._node_args() + self._fees_args() + self._home_args() + self._keyring_backend_args() + \
            self._broadcast_mode_args() + self._yes_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def tx_clp_add_liquidity(self, from_addr: cosmos.Address, symbol: str, native_amount: int, external_amount: int, /,
        account_seq: Optional[Tuple[int, int]] = None
    ) -> JsonDict:
        args = ["tx", "clp", "add-liquidity", "--from", from_addr, "--symbol", symbol, "--nativeAmount",
            str(native_amount), "--externalAmount", str(external_amount)] + self._home_args() + \
            self._keyring_backend_args() + self._chain_id_args() + self._node_args() + self._fees_args() + \
            self._account_number_and_sequence_args(account_seq) + self._broadcast_mode_args() + \
            self._yes_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def tx_clp_unbond_liquidity(self, from_addr: cosmos.Address, symbol: str, units: int) -> JsonDict:
        args = ["tx", "clp", "unbond-liquidity", "--from", from_addr, "--symbol", symbol, "--units", str(units)] + \
            self._home_args() + self._keyring_backend_args() + self._chain_id_args() + self._node_args() + \
            self._fees_args() + self._broadcast_mode_args() + self._yes_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    # TODO
    #
    # CANCEL-UNBOND-LIQUIDITY
    # ===========================
    # sifnoded tx clp cancel-unbond\
    #   --from $SIF_ACT \
    #   --keyring-backend test \
    #   --symbol cusdc \
    #   --units 77971121144444445014456057 \
    #   --fees 100000000000000000rowan \
    #   --chain-id $SIFNODE_CHAIN_ID \
    #   --broadcast-mode block \
    #   -y

    # REMOVE-LIQUIDITY
    # ===========================
    # sifnoded tx clp remove-liquidity-units \
    #   --from $SIF_ACT \
    #   --keyring-backend test \
    #   --symbol cusdc \
    #   --withdrawUnits 77971121144444445014456057 \
    #   --fees 100000000000000000rowan \
    #   --chain-id $SIFNODE_CHAIN_ID \
    #   --broadcast-mode block \
    #   -y

    # SWAP-IN
    # ============================
    # sifnoded tx clp swap \
    #   --from $SIF_ACT \
    #   --keyring-backend test \
    #   --sentSymbol cusdc \
    #   --receivedSymbol rowan \
    #   --sentAmount 184515000000 \
    #   --minReceivingAmount 0 \
    #   --fees 100000000000000000rowan \
    #   --chain-id $SIFNODE_CHAIN_ID \
    #   --broadcast-mode block \
    #   -y

    # SWAP-OUT
    # ============================
    # sifnoded tx clp swap \
    #   --from $SIF_ACT \
    #   --keyring-backend test \
    #   --sentSymbol rowan \
    #   --receivedSymbol cusdc \
    #   --sentAmount 77971121000000000000000000 \
    #   --minReceivingAmount 0 \
    #   --fees 100000000000000000rowan \
    #   --chain-id $SIFNODE_CHAIN_ID \
    #   --broadcast-mode block \
    #   -y

    def tx_clp_cancel_unbound(self):
        assert False, "Not implemented yet"  # TODO

    def tx_clp_remove_liquidity_units(self):
        assert False, "Not implemented yet"  # TODO

    def tx_clp_swap(self):
        assert False, "Not implemented yet"  # TODO

    def clp_reward_period(self, from_addr: cosmos.Address, rewards_params: RewardsParams):
        with self._with_temp_json_file([rewards_params]) as tmp_rewards_json:
            args = ["tx", "clp", "reward-period", "--from", from_addr, "--path", tmp_rewards_json] + \
                self._home_args() + self._keyring_backend_args() + self._node_args() + self._chain_id_args() + \
                self._fees_args() + self._broadcast_mode_args() + self._yes_args()
            res = self.sifnoded_exec(args)
            res = sifnoded_parse_output_lines(stdout(res))
            return res

    def query_reward_params(self):
        args = ["query", "reward", "params"] + self._node_args()
        res = self.sifnoded_exec(args)
        return res

    def clp_set_lppd_params(self, from_addr: cosmos.Address, lppd_params: LPPDParams):
        with self._with_temp_json_file([lppd_params]) as tmp_distribution_json:
            args = ["tx", "clp", "set-lppd-params", "--path", tmp_distribution_json, "--from", from_addr] + \
                self._home_args() + self._keyring_backend_args() + self._node_args() + self._chain_id_args() + \
                self._fees_args() + self._broadcast_mode_args() + self._yes_args()
            res = self.sifnoded_exec(args)
            res = sifnoded_parse_output_lines(stdout(res))
            return res

    def sign_transaction(self, tx: JsonDict, from_sif_addr: cosmos.Address, sequence: int = None,
        account_number: int = None
    ) -> Mapping:
        assert (sequence is not None) == (account_number is not None)  # We need either both or none
        with self._with_temp_json_file(tx) as tmp_tx_file:
            args = ["tx", "sign", tmp_tx_file, "--from", from_sif_addr] + \
                (["--sequence", str(sequence), "--offline", "--account-number", str(account_number)] if sequence else []) + \
                self._home_args() + self._keyring_backend_args() + self._chain_id_args() + self._node_args()
            res = self.sifnoded_exec(args)
            signed_tx = json.loads(stderr(res))
            return signed_tx

    def encode_transaction(self, tx: JsonObj) -> bytes:
        with self._with_temp_json_file(tx) as tmp_file:
            res = self.sifnoded_exec(["tx", "encode", tmp_file])
            encoded_tx = base64.b64decode(stdout(res))
            return encoded_tx

    @contextlib.contextmanager
    def _with_temp_json_file(self, json_obj: JsonObj) -> str:
        with self.cmd.with_temp_file() as tmpfile:
            self.cmd.write_text_file(tmpfile, json.dumps(json_obj))
            yield tmpfile

    def sifnoded_exec(self, args: List[str], stdin: Union[str, bytes, Sequence[str], None] = None,
        cwd: Optional[str] = None, disable_log: bool = False
    ) -> command.ExecResult:
        args = [self.binary] + args
        res = self.cmd.execst(args, stdin=stdin, cwd=cwd, disable_log=disable_log)
        return res

    # Block has to be mined, does not work for block 0
    def get_block_results(self, height: Optional[int] = None):
        path = "block_results{}".format("?height={}".format(height) if height is not None else "")
        host, port = self._get_host_and_port()
        return self._rpc_get(host, port, path)["result"]

    def _get_host_and_port(self) -> Tuple[str, int]:
        # TODO HACK
        # TODO Refactor ports
        # TODO Better store self.host and self.port and make self.node a calculated property
        if self.node is None:
            return LOCALHOST, SIFNODED_DEFAULT_RPC_PORT
        else:
            m = re.compile("^tcp://(.+):(.+)$").match(self.node)
            assert m, "Not implemented"
            host, port = m[1], int(m[2])
            if host == ANY_ADDR:
                host = LOCALHOST
            return host, port

    def _rpc_get(self, host, port, relative_url):
        url = "http://{}:{}/{}".format(host, port, relative_url)
        http_result_payload = http_get(url)
        log.debug("Result for {}: {} bytes".format(url, len(http_result_payload)))
        return json.loads(http_result_payload.decode("UTF-8"))

    def wait_for_last_transaction_to_be_mined(self, count: int = 1, disable_log: bool = True, timeout: int = 90):
        log.debug("Waiting for last sifnode transaction to be mined...")
        start_time = time.time()
        initial_block = self.get_current_block()
        while self.get_current_block() < initial_block + count:
            time.sleep(1)
            if time.time() - start_time > timeout:
                raise Exception("Timeout expired while waiting for last sifnode transaction to be mined")

    # TODO Refactor wait_up() / _wait_up()
    def wait_up(self, host, port):
        from urllib.error import URLError
        while True:
            try:
                return self.get_status(host, port)
            except URLError:
                time.sleep(1)

    # TODO Refactor wait_up() / _wait_up()
    def _wait_up(self, timeout: int=30):
        host, port = self._get_host_and_port()
        from urllib.error import URLError
        start_time = time.time()
        while True:
            try:
                response = self._rpc_get(host, port, "status")
                result = response["result"]
                return result
            except URLError:
                if time.time() - start_time > timeout:
                    raise Exception("Timeout waiting for sifnoded to come up")
                time.sleep(1)

    def _home_args(self) -> Optional[List[str]]:
        return ["--home", self.home] if self.home else []

    def _keyring_backend_args(self) -> Optional[List[str]]:
        return ["--keyring-backend", self.keyring_backend] if self.keyring_backend else []

    def _gas_prices_args(self) -> List[str]:
        return ["--gas-prices", "0.5rowan", "--gas-adjustment", "1.5"]

    def _fees_args(self) -> List[str]:
        sifnode_tx_fees = [10 ** 17, "rowan"]
        return [
            # Deprecated: sifnoded accepts --gas-prices=0.5rowan along with --gas-adjustment=1.5 instead of a fixed fee.
            # "--gas-prices", "0.5rowan", "--gas-adjustment", "1.5",
            "--fees", sif_format_amount(*sifnode_tx_fees)]

    def _chain_id_args(self) -> List[str]:
        assert self.chain_id
        return ["--chain-id", self.chain_id]

    def _node_args(self) -> Optional[List[str]]:
        return ["--node", self.node] if self.node else []

    def _account_number_and_sequence_args(self, account_seq: Optional[Tuple[int, int]] = None) -> Optional[List[str]]:
            return ["--account-number", str(account_seq[0]), "--sequence", str(account_seq[1])] if account_seq is not None else []

    # One of sync|async|block; block will actually get us raw_message
    def _broadcast_mode_args(self, broadcast_mode: Optional[str] = None) -> Optional[List[str]]:
        broadcast_mode = broadcast_mode if broadcast_mode is not None else self.broadcast_mode
        return ["--broadcast-mode", broadcast_mode] if broadcast_mode is not None else []

    def _yes_args(self) -> List[str]:
        return ["--yes"]


# Refactoring in progress - this class is supposed to become the successor of class Sifnode.
# It wraps node, home, chain_id, fees and keyring backend
# TODO Remove 'ctx' (currently needed for cross-chain fees for Peggy2)
class SifnodeClient:
    def __init__(self, ctx, sifnode: Sifnoded, node: Optional[str] = None, chain_id: Optional[str] = None,
        grpc_port: Optional[int] = None
    ):
        self.sifnode = sifnode
        self.ctx = ctx  # TODO Remove (currently needed for cross-chain fees for Peggy2)
        self.grpc_port = grpc_port

    def query_account(self, sif_addr):
        result = json.loads(stdout(self.sifnode.sifnoded_exec(["query", "account", sif_addr, "--output", "json"])))
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
        args = ["tx", "ethbridge", direction, to_eth_addr, str(amount), denom, str(cross_chain_ceth_fee),
                "--network-descriptor", str(eth.ethereum_network_descriptor),  # Mandatory
                "--from", from_sif_addr,  # Mandatory, either name from keyring or address
                "--output", "json",
            ] + \
            (["--generate-only"] if generate_only else []) + \
            self.sifnode._gas_prices_args() + \
            self.sifnode._home_args() + \
            (self.sifnode._keyring_backend_args() if not generate_only else []) + \
            self.sifnode._chain_id_args() + \
            self.sifnode._node_args() + \
            self.sifnode._yes_args()
        res = self.sifnode.sifnoded_exec(args)
        result = json.loads(stdout(res))
        if not generate_only:
            assert "failed to execute message" not in result["raw_log"]
        return result

    def send_from_sifchain_to_ethereum_grpc(self, from_sif_addr: cosmos.Address, to_eth_addr: str, amount: int,
        denom: str
    ):
        tx = self.send_from_sifchain_to_ethereum(from_sif_addr, to_eth_addr, amount, denom, generate_only=True)
        signed_tx = self.sifnode.sign_transaction(tx, from_sif_addr)
        encoded_tx = self.sifnode.encode_transaction(signed_tx)
        result = self.broadcast_tx(encoded_tx)
        return result

    def open_grpc_channel(self) -> grpc.Channel:
        # See https://docs.cosmos.network/v0.44/core/proto-docs.html
        # See https://docs.cosmos.network/v0.44/core/grpc_rest.html
        # See https://app.swaggerhub.com/apis/Ivan-Verchenko/sifnode-swagger-api/1.1.1
        # See https://raw.githubusercontent.com/Sifchain/sifchain-ui/develop/ui/core/swagger.yaml
        return grpc.insecure_channel("{}:{}".format(LOCALHOST,
            self.grpc_port if self.grpc_port is not None else SIFNODED_DEFAULT_GRPC_PORT))

    def broadcast_tx(self, encoded_tx: bytes):
        ondemand_import_generated_protobuf_sources()
        broadcast_mode = cosmos_pb.BROADCAST_MODE_ASYNC
        with self.open_grpc_channel() as channel:
            tx_stub = cosmos_pb_grpc.ServiceStub(channel)
            req = cosmos_pb.BroadcastTxRequest(tx_bytes=encoded_tx, mode=broadcast_mode)
            resp = tx_stub.BroadcastTx(req)
            return resp


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

    def peggy2_build_ebrelayer_cmd(self, init_what: str, network_descriptor: int, tendermint_node: str,
        web3_provider: str, bridge_registry_contract_address: eth.Address, validator_mnemonic: str, chain_id: str,
        node: Optional[str] = None, keyring_backend: Optional[str] = None, keyring_dir: Optional[str] = None,
        sign_with: Optional[str] = None, symbol_translator_file: Optional[str] = None, log_format: Optional[str] = None,
        max_fee_per_gas: Optional[int] = None, max_priority_fee_per_gas: Optional[str] = None,
        extra_args: Mapping[str, Any] = None, ethereum_private_key: Optional[eth.PrivateKey] = None,
        ethereum_address: Optional[eth.Address] = None, home: Optional[str] = None, log_level: Optional[str] = None,
        cwd: Optional[str] = None
    ) -> command.ExecArgs:
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
            (["--log_format", log_format] if log_format else []) + \
            (["--log_level", log_level] if log_level else []) + \
            (["--maxFeePerGasFlag", str(max_fee_per_gas)] if max_fee_per_gas is not None else []) + \
            (["--maxPriorityFeePerGasFlag", str(max_priority_fee_per_gas)] if max_priority_fee_per_gas is not None else [])

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


class RateLimiter:
    def __init__(self, sifnoded, max_tpb):
        self.sifnoded = sifnoded
        self.max_tpb = max_tpb
        self.counter = 0

    def limit(self):
        if self.max_tpb == 0:
            pass
        self.counter += 1
        if self.counter == self.max_tpb:
            self.sifnoded.wait_for_last_transaction_to_be_mined()
            self.counter = 0


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
