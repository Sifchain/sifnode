import base64
import json
import os
import random
import time
import typing
from typing import Iterable, Mapping, Union, List, Callable
import web3
from web3.eth import Contract
from hexbytes import HexBytes
from web3.types import TxReceipt

from siftool import eth, truffle, hardhat, run_env, sifchain, cosmos
from siftool.common import *

# These are utilities to interact with running environment (running agains local ganache-cli/hardhat/sifnoded).
# This is to replace test_utilities.py, conftest.py, burn_lock_functions.py and integration_test_context.py.
# Also to replace smart-contracts/scripts/...


CETH = "ceth"  # Peggy1 only (Peggy2.0 uses denom hash)
ROWAN = "rowan"

sifnode_funds_for_transfer_peggy1 = 10**17  # rowan

log = siftool_logger(__name__)


# This is called from test fixture and will optionally set a snapshot to run the test in.
def get_test_env_ctx(snapshot_name=None):
    assert snapshot_name is None, "Not implemented yet"
    return get_env_ctx()

# This returns an EnvCtx connected to a running environment.
def get_env_ctx(cmd=None, env_file=None, env_vars=None):
    assert cmd is None
    assert env_file is None
    assert env_vars is None
    if on_peggy2_branch:
        ctx = get_env_ctx_peggy2()
    else:
        ctx = get_env_ctx_peggy1()

    # Add any Ethereum private keys to memory
    eth_user_private_keys = ctx.cmd.project.read_peruser_config_file("eth-keys")
    if eth_user_private_keys:
        available_test_accounts = []
        for address, key in [[e["address"], e["key"]] for e in eth_user_private_keys]:
            address, key = eth.validate_address_and_private_key(address, key)
            available_test_accounts.append(address)
            ctx.eth.set_private_key(address, key)
        ctx.available_test_eth_accounts = available_test_accounts

    # Add any Sifchain private keys to test keystore
    sif_user_private_keys = ctx.cmd.project.read_peruser_config_file("sif-keys")
    if sif_user_private_keys:
        available_sif_accounts = ctx.sifnode.keys_list()
        for name, address, mnemonic in [[e["name"], e["address"], e["mnemonic"].split(" ")] for e in sif_user_private_keys]:
            existing_acct = [a for a in available_sif_accounts if a["address"] == address]
            if not existing_acct:
                acct = ctx.sifnode.keys_add(name, mnemonic)
                assert acct["address"] == address, "Invalid address for sif account {}".format(name)
    return ctx

def get_env_ctx_peggy2():
    cmd = run_env.Integrator()
    dot_env_vars = json.loads(cmd.read_text_file(cmd.project.project_dir("smart-contracts/env.json")))
    environment_vars = json.loads(cmd.read_text_file(cmd.project.project_dir("smart-contracts/environment.json")))

    deployed_contract_address_overrides = get_overrides_for_smart_contract_addresses(dot_env_vars)
    tmp = environment_vars["contractResults"]["contractAddresses"]
    deployed_contract_addresses = dict_merge({
        "BridgeBank": tmp["bridgeBank"],
        "CosmosBridge": tmp["cosmosBridge"],
        "BridgeRegistry": tmp["bridgeRegistry"],
        "Rowan": tmp["rowanContract"],
        "Blocklist": tmp["blocklist"],
    }, deployed_contract_address_overrides)
    abi_provider = hardhat.HardhatAbiProvider(cmd, deployed_contract_addresses)

    # TODO We're mixing "OPERATOR" vs. "OWNER"
    # TODO Addressses from dot_env_vars are not in correct EIP55 "checksum" format
    # operator_address = web3.Web3.toChecksumAddress(dot_env_vars["ETH_ACCOUNT_OPERATOR_ADDRESS"])
    # operator_private_key = dot_env_vars["ETH_ACCOUNT_OPERATOR_PRIVATEKEY"][2:]
    owner_address = web3.Web3.toChecksumAddress(dot_env_vars["ETH_ACCOUNT_OWNER_ADDRESS"])
    owner_private_key = dot_env_vars.get("ETH_ACCOUNT_OWNER_PRIVATEKEY")
    if (owner_private_key is not None) and (owner_private_key.startswith("0x")):
        owner_private_key = owner_private_key[2:]  # TODO Remove
    owner_address, owner_private_key = eth.validate_address_and_private_key(owner_address, owner_private_key)

    rowan_source = dot_env_vars["ROWAN_SOURCE"]

    w3_url = eth.web3_host_port_url(dot_env_vars["ETH_HOST"], int(dot_env_vars["ETH_PORT"]))
    w3_conn = eth.web3_connect(w3_url)

    sifnode_url = dot_env_vars["TCP_URL"]
    sifnode_chain_id = "localnet"  # TODO Mandatory, but not present either in environment_vars or dot_env_vars
    assert dot_env_vars["CHAINDIR"] == dot_env_vars["HOME"]
    sifnoded_home = os.path.join(dot_env_vars["CHAINDIR"], ".sifnoded")
    ethereum_network_descriptor = int(dot_env_vars["ETH_CHAIN_ID"])

    eth_node_is_local = True
    generic_erc20_contract = "BridgeToken"
    ceth_symbol = sifchain.sifchain_denom_hash(ethereum_network_descriptor, eth.NULL_ADDRESS)
    assert ceth_symbol == "sifBridge99990x0000000000000000000000000000000000000000"

    ctx_eth = eth.EthereumTxWrapper(w3_conn, eth_node_is_local)
    ctx = EnvCtx(cmd, w3_conn, ctx_eth, abi_provider, owner_address, sifnoded_home, sifnode_url, sifnode_chain_id,
        rowan_source, ceth_symbol, generic_erc20_contract)
    if owner_private_key:
        ctx.eth.set_private_key(owner_address, owner_private_key)

    ctx.eth.fixed_gas_args = {
        # For ganache:
        # 10000000 exceeds default block limit 6721975 ("--gasLimit")
        # 1000000 out of gas
        "gas": 5000000,
        "gasPrice": ctx.eth.w3_conn.eth.gas_price,
    }
    # Hardhat uses base fee of 7 + 1 GWEI
    # assert ctx.eth.fixed_gas_args["gasPrice"] == 1 * eth.GWEI + 7

    # Monkeypatching for peggy2 extras
    # TODO These are set in run_env.py:Peggy2Environment.init_sifchain(), specifically "sifnoded tx ethbridge set-cross-chain-fee"
    # Consider passing them via environment
    ctx.eth.cross_chain_fee_base = 1
    ctx.eth.cross_chain_lock_fee = 1
    ctx.eth.cross_chain_burn_fee = 1
    ctx.eth.ethereum_network_descriptor = ethereum_network_descriptor

    return ctx

def get_env_ctx_peggy1(cmd=None, env_file=None, env_vars=None):
    cmd = cmd or run_env.Integrator()

    if "ENV_FILE" in os.environ:
        env_file = os.environ["ENV_FILE"]
        env_vars = json.loads(cmd.read_text_file(env_file))
    else:
        env_file = cmd.project.project_dir("test/integration/vagraneenv.json")
        if cmd.exists(env_file):
            env_vars = json.loads(cmd.read_text_file(env_file))
        else:
            # Legacy mode - assume data is in OS environment variables, i.e. running from "start-integration-env.sh"
            # such as CI/CD.
            # For some reason, we get different exceptions and we need to set different parameters.
            # TODO Check: eirher web3.py or ganache-cli might be different
            env_vars = os.environ
            env_vars = cmd.primitive_parse_env_file(cmd.project.project_dir("test/integration/vagrantenv.sh"))
            # is_legacy = True

    collected_private_keys = {}

    deployment_name = env_vars.get("DEPLOYMENT_NAME")

    if "CHAINNET" in env_vars:
        sifnode_chain_id = env_vars["CHAINNET"]
    elif deployment_name:
        sifnode_chain_id = deployment_name
    else:
        sifnode_chain_id = "localnet"

    if "WEB3_PROVIDER" in env_vars:
        w3_url = env_vars["WEB3_PROVIDER"]
    elif "ETHEREUM_WEBSOCKET_ADDRESS" in env_vars:
        # Compatibility with vagrantenv.sh
        w3_url = env_vars["ETHEREUM_WEBSOCKET_ADDRESS"]
    else:
        w3_url = "ws://localhost:7545/"

    if "OWNER" in env_vars:
        # vagrantenv.sh uses OWNER and ETHEREUM_PRIVATE_KEY
        operator_address = env_vars["OWNER"]
        operator_private_key = env_vars.get("ETHEREUM_PRIVATE_KEY")
    else:
        operator_address = env_vars["OPERATOR_ADDRESS"]
        operator_private_key = env_vars.get("OPERATOR_PRIVATE_KEY")
    operator_address, operator_private_key = eth.validate_address_and_private_key(operator_address, operator_private_key)

    # Already added below
    # collected_private_keys[operator_address] = operator_private_key

    if "PAUSER" in env_vars:
        assert env_vars["PAUSER"] == operator_address

    if "ROWAN_SOURCE" in env_vars:
        rowan_source = env_vars["ROWAN_SOURCE"]
    elif "VALIDATOR1_ADDR" in env_vars:
        rowan_source = env_vars["VALIDATOR1_ADDR"]
    else:
        rowan_source = None

    ethereum_network_id = int(env_vars.get("ETHEREUM_NETWORK_ID", 5777))

    generic_erc20_contract_name = "SifchainTestToken"
    if "SMART_CONTRACT_ARTIFACT_DIR" in env_vars:
        artifacts_dir = env_vars["SMART_CONTRACT_ARTIFACT_DIR"]
    elif deployment_name:
        artifacts_dir = cmd.project.project_dir("smart-contracts/deployments/{}/build/contracts".format(deployment_name))
        if deployment_name == "sifchain-1":
            # Special case for Betanet because SifchainTestToken is not deployed there.
            # It's only available on Testnet, Devnet and in local environment.
            # However, BridgeToken will work on Betanet meaning that name(), symbol() and decimals() return meaningful values.
            generic_erc20_contract_name = "BridgeToken"
    else:
        artifacts_dir = cmd.project.project_dir("smart-contracts/build/contracts")

    sifnode_url = env_vars.get("SIFNODE")  # Defaults to "tcp://localhost:26657"
    sifnoded_home = None  # Implies default ~/.sifnoded
    deployed_smart_contract_address_overrides = get_overrides_for_smart_contract_addresses(env_vars)

    w3_conn = eth.web3_connect(w3_url)

    # This variable enables behaviour that is specific to running local Ethereum node (ganache, hardhat):
    # - low-level "advance blocks" command that forces mining of 50 blocks
    # - using fixed gas and gasPrice since we don't care about cost and since ganache doesn't support fee history etc.
    # The following differences might also be considered even though we're not using them yet:
    # - one can use hosted private keys (i.e. using just "transact()" on web3 connection instead of explicit sign_transaction()
    # - additional cleanup after running tests (reclaiming ether from temporary accounts, restoring whitelists/blocklists etc.)
    eth_node_is_local = deployment_name is None

    ctx_eth = eth.EthereumTxWrapper(w3_conn, eth_node_is_local)
    abi_provider = truffle.GanacheAbiProvider(cmd, artifacts_dir, ethereum_network_id, deployed_smart_contract_address_overrides)
    ctx = EnvCtx(cmd, w3_conn, ctx_eth, abi_provider, operator_address, sifnoded_home, sifnode_url, sifnode_chain_id,
        rowan_source, CETH, generic_erc20_contract_name)
    if operator_private_key:
        ctx.eth.set_private_key(operator_address, operator_private_key)

    for addr, private_key in collected_private_keys.items():
        ctx.eth.set_private_key(addr, private_key)

    if eth_node_is_local:
        ctx.eth.fixed_gas_args = {
            # For ganache
            # 10000000 exceeds default block limit 6721975 ("--gasLimit")
            # 1000000 out of gas
            "gas": 5000000,
            "gasPrice": ctx.eth.w3_conn.eth.gas_price,
        }
        assert ctx.eth.fixed_gas_args["gasPrice"] == 20 * eth.GWEI
        # For Ropsten etc. (takes ~30 seconds):
        # web3.gas_strategies.time_based.fast_gas_price_strategy(ctx.eth.w3_conn, {})
    else:
        max_gas = 5000000
        estimator = eth.ExponentiallyWeightedAverageFeeEstimator(w3_conn)
        estimator.coeffs = [
            # Inputs: [1, avg_base_fee, avg_reward, max_priority_fee, gas_price, estimated_gas]
            [max_gas, 0, 0, 0, 0, 0],  # gas returned = max_gas
            [0, 2, 1, 0, 0, 0],  # max_fee_per_gas returned = avg_reward + 2*avg_base_fee
            [0, 0, 1, 0, 0, 0],  # max_priority_fee_per_gas returned = avg_reward
            [0, 0, 0, 0, 1, 0],  # gas_price returned = gas_price
        ]
        ctx.eth.gas_estimate_fn = estimator.estimate_fees

    return ctx


def get_overrides_for_smart_contract_addresses(env_vars):
    mappings = {
        "BridgeBank": "BRIDGE_BANK_ADDRESS",
        "BridgeRegistry": "BRIDGE_REGISTRY_ADDRESS",
        "CosmosBridge": "COSMOS_BRIDGE_ADDRESS",  # Peggy2 only?
        "Rowan": "ROWAN_ADDRESS",  # Peggy2 only?
        "BridgeToken": "BRIDGE_TOKEN_ADDRESS",  # Peggy1 only
    }
    return dict(((k, web3.Web3.toChecksumAddress(env_vars[v])) for k, v in mappings.items() if v in env_vars))


def sif_addr_to_evm_arg(sif_address):
    return sif_address.encode("UTF-8")


class EnvCtx:
    def __init__(self, cmd, w3_conn: web3.Web3, ctx_eth, abi_provider, operator, sifnoded_home, sifnode_url, sifnode_chain_id,
        rowan_source, ceth_symbol, generic_erc20_contract
    ):
        self.cmd = cmd
        self.w3_conn = w3_conn
        self.eth: eth.EthereumTxWrapper = ctx_eth
        self.abi_provider: hardhat.HardhatAbiProvider = abi_provider
        self.operator = operator
        self.sifnode = sifchain.Sifnoded(self.cmd, home=sifnoded_home)
        self.sifnode_url = sifnode_url
        self.sifnode_chain_id = sifnode_chain_id
        # Refactoring in progress: moving stuff into separate client that encapsulates things like url, home and chain_id
        self.sifnode_client = sifchain.SifnodeClient(self.cmd, self, node=sifnode_url, home=sifnoded_home, chain_id=sifnode_chain_id, grpc_port=9090)
        self.rowan_source = rowan_source
        self.ceth_symbol = ceth_symbol
        self.generic_erc20_contract = generic_erc20_contract
        self.available_test_eth_accounts = None

    def get_current_block_number(self) -> int:
        return self.eth.w3_conn.eth.block_number

    # TODO Redirect callers and remove
    def advance_blocks(self, number=50):
        return self.eth.advance_block_w3(number)

    def get_blocklist_sc(self):
        abi, _, address = self.abi_provider.get_descriptor("Blocklist")
        result = self.w3_conn.eth.contract(address=address, abi=abi)
        return result

    def get_bridge_bank_sc(self):
        abi, _, address = self.abi_provider.get_descriptor("BridgeBank")
        assert address, "No address for BridgeBank"
        result = self.w3_conn.eth.contract(address=address, abi=abi)
        return result

    def get_cosmos_bridge_sc(self) -> Contract:
        abi, _, address = self.abi_provider.get_descriptor("CosmosBridge")
        assert address, "No address for CosmosBridge"
        result = self.w3_conn.eth.contract(address=address, abi=abi)
        return result

    def get_generic_erc20_sc(self, address):
        abi, _, _ = self.abi_provider.get_descriptor(self.generic_erc20_contract)
        return self.w3_conn.eth.contract(abi=abi, address=address)

    def get_erc20_token_balance(self, token_addr: eth.Address, eth_addr: eth.Address) -> int:
        token_sc = self.get_generic_erc20_sc(token_addr)
        return token_sc.functions.balanceOf(eth_addr).call()

    def send_erc20_tokens(self, token_addr, from_addr, to_addr, amount):
        token_sc = self.get_generic_erc20_sc(token_addr)
        return self.eth.transact_sync(token_sc.functions.transfer, from_addr)(to_addr, amount)

    # Tries to return any ether on the account to operator
    def scavenge_ether(self, account_addr):
        pass  # TODO

    # <editor-fold desc="Refactored">

    def tx_deploy(self, sc_name, deployer, constructor_args):
        abi, bytecode, _ = self.abi_provider.get_descriptor(sc_name)
        token_sc = self.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
        return self.eth.transact(token_sc.constructor, deployer, tx_opts={"from": deployer})(*constructor_args)

    def tx_get_sc_at(self, sc_name, address):
        abi, _, deployed_address = self.abi_provider.get_descriptor(sc_name)
        address = address if address is not None else deployed_address
        return self.w3_conn.eth.contract(abi=abi, address=address)

    def smart_contract_get_past_events(self, sc, event_name, from_block=None, to_block=None):
        from_block = from_block if from_block is not None else 1
        to_block = str(to_block) if to_block is not None else "latest"
        filter = sc.events[event_name].createFilter(fromBlock=from_block, toBlock=to_block)
        try:
            return filter.get_all_entries()
        finally:
            self.w3_conn.eth.uninstall_filter(filter.filter_id)

    def tx_deploy_new_generic_erc20_token(self, deployer_addr: str, name: str, symbol: str, decimals: int, cosmosDenom: str = None) -> Contract:
        # return self.tx_deploy("SifchainTestToken", self.operator, [name, symbol, decimals])
        if on_peggy2_branch:
            # Use BridgeToken
            assert self.generic_erc20_contract == "BridgeToken"
            if cosmosDenom is None:
                cosmosDenom = "erc20denom"  # TODO Dummy variable since we're using BridgeToken instead of SifchainTestToken

            constructor_args = [name, symbol, decimals, cosmosDenom]
        else:
            # Use SifchainTestToken for TestNet and Devnet, and BridgeToken for Betanet
            token_sc_name = self.generic_erc20_contract
            constructor_args = [name, symbol, decimals]
        abi, bytecode, _ = self.abi_provider.get_descriptor(self.generic_erc20_contract)
        token_sc = self.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
        return self.eth.transact(token_sc.constructor, deployer_addr)(*constructor_args)

    def tx_testing_token_mint(self, token_sc, minter_account, amount, minted_tokens_recipient):
        return self.eth.transact(token_sc.functions.mint, minter_account)(minted_tokens_recipient, amount)

    def tx_update_bridge_bank_whitelist(self, token_addr, value=True):
        bridge_bank = self.get_bridge_bank_sc()
        return self.eth.transact(bridge_bank.functions.updateEthWhiteList, self.operator)(token_addr, value)

    def tx_grant_minter_role(self, token_sc: Contract, minter_addr: str):
        self.get_erc20_token_minter_role(token_sc, minter_addr)
        minter_role_hash = token_sc.functions.MINTER_ROLE().call()
        self.eth.transact(token_sc.functions.grantRole, self.operator)(minter_role_hash, minter_addr)
        assert self.get_erc20_token_minter_role(token_sc, minter_addr) is True

    def get_erc20_token_minter_role(self, token_sc: Contract, minter_addr: str) -> bool:
        minter_role_hash = token_sc.functions.MINTER_ROLE().call()
        return token_sc.functions.hasRole(minter_role_hash, minter_addr).call()

    def tx_approve(self, token_sc, from_addr, to_addr, amount):
        return self.eth.transact(token_sc.functions.approve, from_addr)(to_addr, amount)

    def tx_bridge_bank_lock_eth(self, from_eth_acct, to_sif_acct, amount):
        recipient = sif_addr_to_evm_arg(to_sif_acct)
        bridge_bank = self.get_bridge_bank_sc()
        token_addr = eth.NULL_ADDRESS  # For "eth", otherwise use coin's address
        # Mandatory tx_opts: {"from": from_eth_acct, "gas": max_gas_required, "value": amount}
        # If "value" is missing, we get "call to non-contract"
        tx_opts = {"value": amount}
        return self.eth.transact(bridge_bank.functions.lock, from_eth_acct, tx_opts=tx_opts)(recipient, token_addr, amount)

    def tx_bridge_bank_lock_erc20(self, token_addr, from_eth_acct, to_sif_acct, amount):
        recipient = sif_addr_to_evm_arg(to_sif_acct)
        bridge_bank = self.get_bridge_bank_sc()
        # When transfering ERC20, the amount needs to be passed as argument, and the "message.value" should be 0
        tx_opts = {"value": 0}
        return self.eth.transact(bridge_bank.functions.lock, from_eth_acct, tx_opts=tx_opts)(recipient, token_addr, amount)

    def tx_bridge_bank_burn_erc20(self, token_addr: str, from_eth_acct: str, to_sif_acct: str, amount: int) -> HexBytes:
        recipient = sif_addr_to_evm_arg(to_sif_acct)
        bridge_bank = self.get_bridge_bank_sc()
        # When transfering ERC20, the amount needs to be passed as argument, and the "message.value" should be 0
        tx_opts = {"value": 0}
        return self.eth.transact(bridge_bank.functions.burn, from_eth_acct, tx_opts=tx_opts)(recipient, token_addr, amount)

    def tx_bridge_bank_add_existing_bridge_token(self, token_addr: str) -> HexBytes:
        bridge_bank = self.get_bridge_bank_sc()
        tx_opts = {"value": 0}
        return self.eth.transact(bridge_bank.functions.addExistingBridgeToken, self.operator, tx_opts=tx_opts)(token_addr)

    def tx_approve_and_lock(self, token_sc, from_eth_acct, to_sif_acct, amount):
        bridge_bank_sc = self.get_bridge_bank_sc()
        txhash1 = self.tx_approve(token_sc, self.operator, bridge_bank_sc.address, amount)
        txhash2 = self.tx_bridge_bank_lock_erc20(token_sc.address, from_eth_acct, to_sif_acct, amount)
        log.debug("tx_approve_and_lock: {} '{}' ({}) from {} to {}".format(amount, token_sc.functions.name().call(),
            token_sc.functions.symbol().call(), from_eth_acct, to_sif_acct))
        return txhash1, txhash2

    # </editor-fold>

    # Used from test_integration_framework.py, test_eth_transfers.py
    def deploy_new_generic_erc20_token(self, name: str, symbol: str, decimals: int, owner: str = None, mint_amount: int = None, mint_recipient: str = None, cosmosDenom: str = None) -> Contract:
        owner = self.operator if owner is None else owner
        txhash = self.tx_deploy_new_generic_erc20_token(owner, name, symbol, decimals, cosmosDenom)
        txrcpt = self.eth.wait_for_transaction_receipt(txhash)
        token_addr = txrcpt.contractAddress
        token_sc = self.get_generic_erc20_sc(token_addr)
        assert token_sc.functions.name().call() == name
        assert token_sc.functions.symbol().call() == symbol
        assert token_sc.functions.decimals().call() == decimals
        if mint_amount:
            mint_recipient = mint_recipient or owner
            self.mint_generic_erc20_token(token_sc.address, mint_amount, mint_recipient, minter=owner)
        if not on_peggy2_branch:
            self.update_bridge_bank_whitelist(token_sc.address, True)
        return token_sc

    def mint_generic_erc20_token(self, token_addr, amount, recipient, minter=None):
        minter = minter or self.operator
        token_sc = self.get_generic_erc20_sc(token_addr)
        balance_before = self.get_erc20_token_balance(token_addr, recipient)
        total_supply_before = token_sc.functions.totalSupply().call()
        txhash = self.tx_testing_token_mint(token_sc, minter, amount, recipient)
        txrcpt = self.eth.wait_for_transaction_receipt(txhash)
        assert self.get_erc20_token_balance(token_addr, recipient) == balance_before + amount
        assert token_sc.functions.totalSupply().call() == total_supply_before + amount
        return txrcpt

    # Token symbol must be unique on the blocklist
    def update_bridge_bank_whitelist(self, token_addr, value):
        assert not on_peggy2_branch
        # Token needs to be whitelisted, if it is not, then the transaction will be reverted with a message like this:
        # "Only token in whitelist can be transferred to cosmos"
        # Call of updateEthWhiteList will fail if we try to remove an item from whitelist which is not on the whitelist.
        return self.eth.wait_for_transaction_receipt(self.tx_update_bridge_bank_whitelist(token_addr, value))

    # This function walks through all historical events LogWhiteListUpdate of a BridgeBanksmart contract and builds the
    # current whitelist from live on-chain data.
    def get_whitelisted_tokens_from_bridge_bank_past_events(self):
        bridge_bank = self.get_bridge_bank_sc()
        past_events = self.smart_contract_get_past_events(bridge_bank, "LogWhiteListUpdate")
        result = {}
        for e in past_events:
            token_addr = e.args["_token"]
            value = e.args["_value"]
            assert web3.Web3.toChecksumAddress(token_addr) == token_addr
            # Logically the whitelist only consists of entries that have the last value of True.
            # If the data is clean, then for each token_addr we should first see a True event, possibly
            # followed by alternating False and True. The last value is the active one.
            # However, we want to also preserve False values in the dict since this data is used
            # for inflate_tokens where it matters which tokens should be deployed and which not.
            if token_addr in result:
                if result[token_addr] == value:
                    log.warning(f"Redundant event in BridgeBank's past LogWhiteListUpdate: token_addr={token_addr}, value={value}, blockNumber={e.blockNumber}")
            else:
                if not value:
                    log.warning(f"Redundant event in BridgeBank's past LogWhiteListUpdate: token_addr={token_addr}, value={value}, blockNumber={e.blockNumber}")
            result[token_addr] = value
        return result

    def generate_random_erc20_token_data(self):
        id = random_string(6)
        return ERC20TokenData("test-{}".format(id.lower()), "Test Token {}".format(id), random.choice([0, 4, 6, 9, 18]))

    def get_generic_erc20_token_data(self, token_address):
        token_sc = self.get_generic_erc20_sc(token_address)
        return {
            "symbol": token_sc.functions.symbol().call(),
            "name": token_sc.functions.name().call(),
            "decimals": token_sc.functions.decimals().call()
        }

    def approve_erc20_token(self, token_sc, account_owner, amount):
        bridge_bank_sc = self.get_bridge_bank_sc()
        self.eth.transact_sync(token_sc.functions.approve, account_owner)(bridge_bank_sc.address, amount)

    # TODO Used from integration tests and several other places - do not change
    def create_new_currency(self, symbol, name, decimals, amount, minted_tokens_recipient):
        """
        As in smart-contracts/scripts/test/enableNewToken.js:
        1. Deploys a new instance of SifchainTestToken
        2. Calls BridgeBank.updateEthWhiteList with new token's address
        3. Mint amount to mint_recipient_addr
        4. Approve amount to BridgeBank
        """
        assert self.generic_erc20_contract == "SifchainTestToken"  # Preserve compatibiliy with integration tests and inflate_tokens.sh
        token_sc = self.deploy_new_generic_erc20_token(name, symbol, decimals)
        self.update_bridge_bank_whitelist(token_sc.address, True)
        self.eth.transact_sync(token_sc.functions.mint, self.operator)(minted_tokens_recipient, amount)
        self.approve_erc20_token(token_sc, self.operator, amount)
        return token_sc.address

    # TODO Obsolete, use self.bridge_bank_lock_eth()
    def send_eth_from_ethereum_to_sifchain(self, from_eth_addr, to_sif_addr, amount):
        # recipient = to_sif_addr.encode("UTF-8")
        # coin_denom = eth.NULL_ADDRESS  # For "eth", otherwise use coin's address
        #
        # max_gas_required = 200000
        #
        # bridge_bank = self.get_bridge_bank_sc()
        # txhash = bridge_bank.functions.lock(recipient, coin_denom, amount) \
        #     .transact({"from": from_eth_addr, "gas": max_gas_required, "value": amount})
        # txrcpt = self.w3_conn.eth.wait_for_transaction_receipt(txhash)
        # return txrcpt
        assert False  # TODO

    # TODO Obsolete, use self.bridge_bank_lock_eth()
    def send_erc20_from_ethereum_to_sifchain(self, from_eth_addr, dest_sichain_addr, erc20_token_addr, amount):
        # recipient = dest_sichain_addr.encode("UTF-8")
        #
        # max_gas_required = 200000
        #
        # bridge_bank = self.get_bridge_bank_sc()
        # # When transfering ERC20, the amount needs to be passed as argument, and the "message.value" should be 0
        # # TODO Error handling
        # #      "web3.exceptions.ContractLogicError: execution reverted: SafeERC20: low-level call failed" in case that amount is more than what is available / what was "approved" to BridgeBank
        # tx = bridge_bank.functions.lock(recipient, erc20_token_addr, amount).buildTransaction({
        #     "from": self.operator,
        #     "nonce": self.w3_conn.eth.get_transaction_count(self.operator)
        # })
        # txhash = self.eth.transact(bridge_bank.functions.lock, self.operator)(recipient, erc20_token_addr, amount)
        # # .transact({"from": from_eth_addr, "gas": max_gas_required})
        # txrcpt = self.w3_conn.eth.wait_for_transaction_receipt(txhash)
        # return txrcpt
        token_sc = self.get_generic_erc20_sc(erc20_token_addr)
        self.approve_erc20_token(token_sc, from_eth_addr, amount)
        self.bridge_bank_lock_eth(from_eth_addr, dest_sichain_addr, amount)

    def create_sifchain_addr(self, moniker: str = None, fund_amounts: Union[cosmos.Balance, cosmos.LegacyBalance] = None):
        """
        Generates a new sifchain address in test keyring. If moniker is given, uses it, otherwise
        generates a random one 'test-xxx'. If fund_amounts is given, the sifchain funds are transferred
        from rowan_source to the account before returning.
        """
        moniker = moniker or "test-" + random_string(20)
        acct = self.sifnode.keys_add(moniker)
        sif_address = acct["address"]
        if fund_amounts:
            fund_amounts = cosmos.balance_normalize(fund_amounts)  # Convert from old format if neccessary
            rowan_source_balances = self.get_sifchain_balance(self.rowan_source)
            for denom, required_amount in fund_amounts.items():
                available_amount = rowan_source_balances.get(denom, 0)
                assert available_amount >= required_amount, "Rowan source {} would need {}, but only has {}".format(
                    self.rowan_source, sif_format_amount(required_amount, denom), sif_format_amount(available_amount, denom))
            old_balances = self.get_sifchain_balance(sif_address)
            self.send_from_sifchain_to_sifchain(self.rowan_source, sif_address, fund_amounts)
            self.wait_for_sif_balance_change(sif_address, old_balances, min_changes=fund_amounts)
            new_balances = self.get_sifchain_balance(sif_address)
            assert cosmos.balance_zero(cosmos.balance_sub(new_balances, fund_amounts))
        return sif_address

    def send_from_sifchain_to_sifchain(self, from_sif_addr: cosmos.Address, to_sif_addr: cosmos.Address,
        amounts: cosmos.Balance
    ):
        amounts = cosmos.balance_normalize(amounts)
        amounts_string = cosmos.balance_format(amounts)
        args = ["tx", "bank", "send", from_sif_addr, to_sif_addr, amounts_string] + \
            self.sifnode_client._chain_id_args() + \
            self.sifnode_client._node_args() + \
            self.sifnode_client._fees_args() + \
            ["--yes", "--output", "json"]
        res = self.sifnode.sifnoded_exec(args, sifnoded_home=self.sifnode.home, keyring_backend=self.sifnode.keyring_backend)
        retval = json.loads(stdout(res))
        raw_log = retval["raw_log"]
        for bad_thing in ["insufficient funds", "signature verification failed"]:
            if bad_thing in raw_log:
                raise Exception(raw_log)
        return retval

    def get_sifchain_balance(self, sif_addr: cosmos.Address, limit: Optional[int] = 1000000,
        offset: Optional[int] = None, disable_log: bool = False
    ) -> cosmos.Balance:
        args = ["query", "bank", "balances", sif_addr, "--output", "json"] + \
            (["--limit", str(limit)] if limit is not None else []) + \
            (["--offset", str(offset)] if offset is not None else []) + \
            self.sifnode_client._chain_id_args() + \
            self.sifnode_client._node_args()
        res = self.sifnode.sifnoded_exec(args, sifnoded_home=self.sifnode.home, disable_log=disable_log)
        res = json.loads(stdout(res))
        if res["pagination"]["next_key"] is not None:
            raise Exception("More than {} results in balances".format(limit))
        return {denom: amount for denom, amount in ((x["denom"], int(x["amount"])) for x in res["balances"]) if amount != 0}

    # Experimental
    def get_sifchain_balance_large(self, sif_addr: cosmos.Address, height: Optional[int] = None,
        disable_log: bool = False, retries_on_error: int = 3, delay_on_error: int = 3
    ) -> cosmos.Balance:
        all_balances = {}
        desired_page_size = 5000  # The actual limit might be capped to a lower value, in this case we'll get fewer results
        page_key = None
        while True:
            args = ["query", "bank", "balances", sif_addr, "--output", "json"] + \
                (["--height", str(height)] if height is not None else []) + \
                (["--limit", str(desired_page_size)] if desired_page_size is not None else []) + \
                (["--page-key", page_key] if page_key is not None else []) + \
                self.sifnode_client._chain_id_args() + \
                self.sifnode_client._node_args()
            retries_left = retries_on_error
            while True:
                try:
                    res = self.sifnode.sifnoded_exec(args, sifnoded_home=self.sifnode.home, disable_log=disable_log)
                    break
                except Exception as e:
                    retries_left -= 1
                    log.error("Error reading balances, retries left: {}".format(retries_left))
                    if retries_left > 0:
                        time.sleep(delay_on_error)
                    else:
                        raise e
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

    def get_current_block(self):
        return int(self.status()["SyncInfo"]["latest_block_height"])

    def status(self):
        args = ["status"] + self.sifnode_client._node_args()
        res = self.sifnode.sifnoded_exec(args)
        return json.loads(stderr(res))

    # Unless timed out, this function will exit:
    # - if min_changes are given: when changes are greater.
    # - if expected_balance is given: when balances are equal to that.
    # - if neither min_changes nor expected_balance are given: when anything changes.
    # You cannot use min_changes and expected_balance at the same time.
    def wait_for_sif_balance_change(self, sif_addr: cosmos.Address, old_balance: cosmos.Balance,
        min_changes: cosmos.CompatBalance = None, expected_balance: cosmos.CompatBalance = None, polling_time: int = 1,
        timeout: Optional[int] = 90, change_timeout: int = None, disable_log: bool = True
    ) -> cosmos.Balance:
        assert (min_changes is None) or (expected_balance is None), "Cannot use both min_changes and expected_balance"
        log.debug("Waiting for balance to change for account {}...".format(sif_addr))
        min_changes = None if min_changes is None else cosmos.balance_normalize(min_changes)
        expected_balance = None if expected_balance is None else cosmos.balance_normalize(expected_balance)
        start_time = time.time()
        last_change_time = None
        last_changed_balance = None
        while True:
            new_balance = self.get_sifchain_balance(sif_addr, disable_log=disable_log)
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
            if (timeout is not None) and (now - start_time > timeout):
                raise Exception("Timeout waiting for sif balance to change")
            if last_change_time is None:
                last_changed_balance = new_balance
                last_change_time = now
            else:
                delta = cosmos.balance_sub(new_balance, last_changed_balance)
                if not cosmos.balance_zero(delta):
                    last_changed_balance = new_balance
                    last_change_time = now
                    log.debug("New state detected ({} denoms changed)".format(len(delta)))
                if (change_timeout is not None) and (now - last_change_time > change_timeout):
                    raise Exception("Timeout waiting for sif balance to change")
            time.sleep(polling_time)

    def eth_symbol_to_sif_symbol(self, eth_token_symbol):
        # TODO sifchain.use sifchain_denom_hash() if on_peggy2_branch
        # E.g. "usdt" -> "cusdt"
        if eth_token_symbol == "erowan":
            return ROWAN
        else:
            return "c" + eth_token_symbol.lower()

    # from_sif_addr has to be the address which was used at genesis time for "set-genesis-whitelister-admin".
    # You need to have its private key in the test keyring.
    # This is needed when creating pools for the token or when doing IBC transfers.
    def token_registry_register(self, address, symbol, token_name, decimals, from_sif_addr):
        # Check that we have the private key in test keyring. This will throw an exception if we don't.
        self.cmd.sifnoded_keys_show(from_sif_addr)
        sifchain_symbol = self.eth_symbol_to_sif_symbol(symbol)
        upper_symbol = symbol.upper()  # Like "USDT"
        # See scripts/ibc/tokenregistration for more information and examples.
        # JSON file can be generated with "sifnoded q tokenregistry generate"
        token_data = {"entries": [{
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
            "permissions": ["CLP", "IBCEXPORT", "IBCIMPORT"],
            "unit_denom": "",
            "ibc_counterparty_denom": "",
            "ibc_counterparty_chain_id": "",
        }]}
        tmp_registry_json = self.cmd.mktempfile()
        try:
            self.cmd.write_text_file(tmp_registry_json, json.dumps(token_data, indent=4))
            args = ["tx", "tokenregistry", "register", tmp_registry_json] + \
                self.sifnode_client._chain_id_args() + \
                self.sifnode_client._node_args() + \
                self.sifnode_client._fees_args() + [
                "--from", from_sif_addr,
                "--output", "json",
                "--broadcast-mode", "block",  # One of sync|async|block; block will actually get us raw_message
                "--yes"
            ]
            res = self.cmd.sifnoded_exec(args, keyring_backend="test")
            res = json.loads(stdout(res))
            # Example of successful output: {"height":"196804","txhash":"C8252E77BCD441A005666A4F3D76C99BD35F9CB49AA1BE44CBE2FFCC6AD6ADF4","codespace":"","code":0,"data":"0A270A252F7369666E6F64652E746F6B656E72656769737472792E76312E4D73675265676973746572","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"/sifnode.tokenregistry.v1.MsgRegister\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"/sifnode.tokenregistry.v1.MsgRegister"}]}]}],"info":"","gas_wanted":"200000","gas_used":"115149","tx":null,"timestamp":""}
            if res["raw_log"].startswith("signature verification failed"):
                raise Exception(res["raw_log"])
            if res["raw_log"].startswith("failed to execute message"):
                raise Exception(res["raw_log"])
            return res
        finally:
            self.cmd.rm(tmp_registry_json)

    # Deprecated: sifnoded accepts --gas-prices=0.5rowan along with --gas-adjustment=1.5 instead of a fixed fee.
    # Using those parameters is the best way to have the fees set robustly after the .42 upgrade.
    # See https://github.com/Sifchain/sifnode/pull/1802#discussion_r697403408
    # The corresponding denom should be "rowan".
    @property
    def sifchain_fees(self):
        return 200000

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        # If we'running Ropsten and not local hardhat/ganache, try to reclaim any remaining eth.
        if not self.eth.is_local_node:
            # self.scavenge_ether()
            pass

    def wait_for_eth_balance_change(self, eth_addr, old_balance: int, timeout=90, polling_time=1, token_addr=None):
        start_time = time.time()
        while True:
            new_balance = self.get_erc20_token_balance(token_addr, eth_addr) if token_addr \
                else self.eth.get_eth_balance(eth_addr)
            # log.debug("wait_for_eth_balance_change(): {}={}".format(eth_addr, new_balance))
            if new_balance != old_balance:
                return new_balance
            time.sleep(polling_time)
            now = time.time()
            if now - start_time > timeout:
                raise Exception("Timeout waiting for Ethereum balance to change")

    def wait_for_new_bridge_token_created(self, cosmos_denom: str, timeout: int = 90, polling_time: int = 1) -> str:
        start_time = time.time()
        while True:
            cosmos_bridge_sc = self.get_cosmos_bridge_sc()
            events = self.smart_contract_get_past_events(cosmos_bridge_sc, "LogNewBridgeTokenCreated")

            if len(events) > 0:
                for e in events:
                    if e.args["cosmosDenom"] == cosmos_denom:
                        return e.args["bridgeTokenAddress"]

            time.sleep(polling_time)
            now = time.time()
            if now - start_time > timeout:
                raise Exception("Timeout waiting for Ethereum balance to change")

    def create_and_fund_eth_account(self, fund_from=None, fund_amount=None):
        if self.available_test_eth_accounts is not None:
            address = self.available_test_eth_accounts.pop(0)
        else:
            # If None, we're generating non-repeatable accounts.
            address, key = self.eth.create_new_eth_account()
            self.eth.set_private_key(address, key)
            assert self.eth.get_eth_balance(address) == 0
        if fund_amount is not None:
            fund_from = fund_from or self.operator
            funder_balance_before = self.eth.get_eth_balance(fund_from)
            assert funder_balance_before >= fund_amount, "Cannot fund created account with ETH: {} needs {}, but has {}" \
                .format(fund_from, fund_amount, funder_balance_before)
            target_balance_before = self.eth.get_eth_balance(address)
            difference = fund_amount - target_balance_before
            if difference > 0:
                self.eth.send_eth(fund_from, address, difference)
                assert self.eth.get_eth_balance(address) == fund_amount
                assert self.eth.get_eth_balance(fund_from) < funder_balance_before - difference
        return address

    def bridge_bank_lock_eth(self, from_eth_acct, to_sif_acct, amount):
        """ Sends ETH from Ethereum to Sifchain (lock) """
        txhash = self.tx_bridge_bank_lock_eth(from_eth_acct, to_sif_acct, amount)
        return self.eth.wait_for_transaction_receipt(txhash)

    def bridge_bank_lock_erc20(self, token_sc, from_eth_acct, to_sif_acct, amount):
        txhash = self.tx_bridge_bank_lock_erc20(token_sc.address, from_eth_acct, to_sif_acct, amount)
        return self.eth.wait_for_transaction_receipt(txhash)

    def bridge_bank_burn_erc20(self, token_sc: Contract, from_eth_acct: str, to_sif_acct: str, amount: int) -> TxReceipt:
        txhash = self.tx_bridge_bank_burn_erc20(token_sc.address, from_eth_acct, to_sif_acct, amount)
        return self.eth.wait_for_transaction_receipt(txhash)

    def bridge_bank_add_existing_bridge_token(self, token_addr: str):
        txhash = self.tx_bridge_bank_add_existing_bridge_token(token_addr)
        self.eth.wait_for_transaction_receipt(txhash)
        final_value = self.get_cosmos_token_in_white_list(token_addr)
        assert final_value is True

    def get_cosmos_token_in_white_list(self, token_addr: str) -> bool:
        bridge_bank_sc = self.get_bridge_bank_sc()
        return bridge_bank_sc.functions.getCosmosTokenInWhiteList(token_addr).call()

    def get_destination_contract_address(self, cosmos_denom: str) -> Contract:
        cosmos_bridge_sc = self.get_cosmos_bridge_sc()
        return cosmos_bridge_sc.functions.cosmosDenomToDestinationAddress(cosmos_denom).call()

    # TODO At the moment this is only for Ethereum-native assets (ETH and ERC20 tokens) which always use "lock".
    # For Sifchain-native assets (rowan) we need to use "burn".
    # Compare: smart-contracts/scripts/test/{sendLockTx.js OR sendBurnTx.js}
    # sendBurnTx is called when sifchain_symbol == "rowan", sendLockTx otherwise
    def send_from_ethereum_to_sifchain(self, from_eth_acct: str, to_sif_acct: str, amount: int, token_sc: Contract = None, isLock: bool = True) -> TxReceipt:
        if token_sc is None:
            # ETH transfer
            self.bridge_bank_lock_eth(from_eth_acct, to_sif_acct, amount)
        else:
            # ERC20 token transfer
            self.approve_erc20_token(token_sc, from_eth_acct, amount)
            if isLock:
                self.bridge_bank_lock_erc20(token_sc, from_eth_acct, to_sif_acct, amount)
            else:
                self.bridge_bank_burn_erc20(token_sc, from_eth_acct, to_sif_acct, amount)

    # Peggy1-specific
    def set_ofac_blocklist_to(self, addrs):
        blocklist_sc = self.get_blocklist_sc()
        addrs = [web3.Web3.toChecksumAddress(addr) for addr in addrs]
        existing_entries = blocklist_sc.functions.getFullList().call()
        to_add = [addr for addr in addrs if addr not in existing_entries]
        to_remove = [addr for addr in existing_entries if addr not in addrs]
        result = [None, None]
        if to_add:
            result[0] = self.eth.transact_sync(blocklist_sc.functions.batchAddToBlocklist, self.operator)(to_add)
        if to_remove:
            result[1] = self.eth.transact_sync(blocklist_sc.functions.batchRemoveFromBlocklist, self.operator)(to_remove)
        current_entries = blocklist_sc.functions.getFullList().call()
        assert set(addrs) == set(current_entries)
        return result

    def sanity_check(self):
        """ Tries to catch some common configurtion errors. """
        bridge_bank_sc = self.get_bridge_bank_sc()
        if on_peggy2_branch:
            pass
        else:
            assert (self.sifnode_chain_id != "sifchain-testnet-1") or (bridge_bank_sc.address == "0x6CfD69783E3fFb44CBaaFF7F509a4fcF0d8e2835")
            assert (self.sifnode_chain_id != "sifchain-devnet-1") or (bridge_bank_sc.address == "0x96DC6f02C66Bbf2dfbA934b8DafE7B2c08715A73")
            assert (self.sifnode_chain_id != "localnet") or (bridge_bank_sc.address == "0x30753E4A8aad7F8597332E813735Def5dD395028")
        assert bridge_bank_sc.functions.owner().call() == self.operator, \
            "BridgeBank owner is {}, but OPERATOR is {}".format(bridge_bank_sc.functions.owner().call(), self.operator)
        operator_balance = self.eth.get_eth_balance(self.operator) / eth.ETH
        assert operator_balance >= 1, "Insufficient operator balance, should be at least 1 ETH"

        available_accounts = self.sifnode.keys_list()
        rowan_source_account = [x for x in available_accounts if x["address"] == self.rowan_source]
        assert len(rowan_source_account) == 1, "There should be exactly one key in test keystore corresponding to " \
            "ROWAN_SOURCE {}".format(self.rowan_source)
        if len(rowan_source_account) != 1:
            raise Exception
        rowan_source_balance = self.get_sifchain_balance(self.rowan_source).get(ROWAN, 0)
        min_rowan_source_balance = 10 * 10**18
        assert rowan_source_balance > min_rowan_source_balance, "ROWAN_SOURCE should have at least {}rowan balance, " \
            "but has only {}rowan".format(min_rowan_source_balance, rowan_source_balance)


class ERC20TokenData:
    def __init__(self, symbol, name, decimals):
        self.symbol: string = symbol
        self.name: string = name
        self.decimals: int = decimals


def recover_eth_from_test_accounts():
    ctx = get_test_env_ctx()
    w = eth.ExponentiallyWeightedAverageFeeEstimator()

    gas_price = 20 * eth.GWEI
    tx_cost = eth.MIN_TX_GAS * gas_price
    total_recovered = 0
    for addr in ctx.available_test_eth_accounts:
        balance = ctx.eth.get_eth_balance(addr)
        to_recover = balance - tx_cost
        if to_recover > 0:
            log.info("Account {}: balance={}, to_recover={}".format(addr, balance//eth.GWEI, to_recover//eth.GWEI))
            ctx.eth.send_eth(addr, ctx.operator, to_recover)
            total_recovered += to_recover
    log.info("Total recovered: {} ETH".format(total_recovered/eth.ETH))


def sifnoded_parse_output_lines(stdout):
    pat = re.compile("^(.*?): (.*)$")
    result = {}
    for line in stdout.splitlines():
        m = pat.match(line)
        result[m[1]] = m[2]
    return result

# Generalized version of "grep -B _ -A _". Can be used as iterator on long streams without loading everything to memory.
def generalized_grep(items: Iterable, match_fn: Callable, before: int = 0, after: int = 0):
    it = iter(items)
    buf = []
    matched = False
    while True:
        try:
            item = next(it)
        except StopIteration:
            break
        if len(buf) > before + 1:
            buf.pop(0)
        buf.append(item)
        if match_fn(item):
            yield from buf
            matched = True
            break
    if matched:
        for _ in range(after):
            try:
                item = next(it)
            except StopIteration:
                break
            yield item

def pytest_ctx_fixture(request):
    # To pass the "snapshot_name" as a parameter with value "foo" from test, annotate the test function like this:
    # @pytest.mark.snapshot_name("foo")
    snapshot_name = request.node.get_closest_marker("snapshot_name")
    if snapshot_name is not None:
        snapshot_name = snapshot_name.args[0]
        logging.debug("Context setup: snapshot_name={}".format(repr(snapshot_name)))
    with get_test_env_ctx() as ctx:
        yield ctx
        logging.debug("Test context cleanup")

def pytest_test_wrapper_fixture():
    disable_noisy_loggers()
