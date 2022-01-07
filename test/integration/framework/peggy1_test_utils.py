import json
import os
import time

import main
import eth
from common import *


# These are utilities to interact with running Peggy1 environment (running agains local ganache-cli/hardhat/sifnoded).
# This is to replace test_utilities.py, conftest.py, burn_lock_functions.py and integration_test_context.py.
# Also to replace smart-contracts/scripts/...


CETH = "ceth"
ROWAN = "rowan"

sifnode_funds_for_transfer_peggy1 = 10**17  # rowan

def get_peggy1_env_ctx_test(cmd=None, env_file=None, env_vars=None):
    return get_env_ctx(cmd=cmd, env_file=env_file, env_vars=env_vars)

def get_env_ctx(cmd=None, env_file=None, env_vars=None):
    cmd = cmd or main.Integrator()

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
        operator_private_key = env_vars["ETHEREUM_PRIVATE_KEY"]
    else:
        operator_address = env_vars["OPERATOR_ADDRESS"]
        operator_private_key = env_vars["OPERATOR_PRIVATE_KEY"]

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

    if "SMART_CONTRACT_ARTIFACT_DIR" in env_vars:
        artifacts_dir = env_vars["SMART_CONTRACT_ARTIFACT_DIR"]
    elif deployment_name:
        artifacts_dir = cmd.project.project_dir("smart-contracts/deployments/{}/build".format(deployment_name))
    else:
        artifacts_dir = cmd.project.project_dir("smart-contracts/build")

    sifnode_url = env_vars.get("SIFNODE")  # Defaults to "tcp://localhost:26657"

    w3_conn = eth.web3_connect(w3_url, websocket_timeout=30)

    # This variable enables behaviour that is specific to running local Ethereum node (ganache, hardhat):
    # - low-level "advance blocks" command that forces mining of 50 blocks
    # - using fixed gas and gasPrice since we don't care about cost and since ganache doesn't support fee history etc.
    # The following differences might also be considered even though we're not using them yet:
    # - one can use hosted private keys (i.e. using just "transact()" on web3 connection instead of explicit sign_transaction()
    # - additional cleanup after running tests (reclaiming ether from temporary accounts, restoring whitelists/blocklists etc.)
    eth_node_is_local = deployment_name is None

    ctx = get_ctx(w3_conn, cmd, artifacts_dir, ethereum_network_id, operator_address, sifnode_url, sifnode_chain_id,
        rowan_source, operator_private_key, eth_node_is_local)

    for addr, private_key in collected_private_keys.items():
        ctx.w3_tx.set_private_key(addr, private_key)

    if eth_node_is_local:
        ctx.w3_tx.fixed_gas_args = {
            # For ganache
            # 10000000 exceeds default block limit 6721975 ("--gasLimit")
            # 1000000 out of gas
            "gas": 5000000,
            "gasPrice": ctx.w3_tx.w3_conn.eth.gas_price,
        }
        assert ctx.w3_tx.fixed_gas_args["gasPrice"] == 20 * eth.GWEI
        # For Ropsten etc. (takes ~30 seconds):
        # web3.gas_strategies.time_based.fast_gas_price_strategy(ctx.w3_tx.w3_conn, {})
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
        ctx.w3_tx.gas_estimate_fn = estimator.estimate_fees

    user_private_keys = ctx.cmd.project.read_peruser_config_file("user_private_keys")
    if user_private_keys:
        available_test_accounts = []
        for address, key in [[entry["address"], entry["key"]] for entry in user_private_keys]:
            available_test_accounts.append(address)
            ctx.w3_tx.set_private_key(address, key)
        ctx.available_test_eth_accounts = available_test_accounts

    return ctx

def get_ctx(w3_conn, cmd, artifacts_dir, ethereum_network_id, operator_address, sifnode_url, sifnode_chain_id,
    rowan_source, operator_private_key, eth_node_is_local
):
    w3_tx = eth.EthereumTxWrapper(w3_conn, eth_node_is_local)
    abi_provider = GanacheAbiProvider(cmd, artifacts_dir, ethereum_network_id)
    ctx = Peggy1EnvCtx(cmd, w3_conn, w3_tx, abi_provider, operator_address, sifnode_url, sifnode_chain_id, rowan_source)
    ctx.w3_tx.set_private_key(operator_address, operator_private_key)
    return ctx

def sif_addr_to_evm_arg(sif_address):
    return sif_address.encode("UTF-8")


class GanacheAbiProvider:
    def __init__(self, cmd, artifacts_dir, ethereum_network_id):
        self.cmd = cmd
        self.artifacts_dir = artifacts_dir
        self.ethereum_default_network_id = ethereum_network_id

    def get_compiled_sc(self, sc_name):
        path = self.cmd.project.project_dir(self.artifacts_dir, "contracts/{}.json".format(sc_name))
        return json.loads(self.cmd.read_text_file(path))

    def get_sc_abi(self, sc_name):
        tmp = self.get_compiled_sc(sc_name)
        return tmp["networks"][str(self.ethereum_default_network_id)]["address"], tmp["abi"]

    def get_descriptor(self, sc_name):
        tmp = self.get_compiled_sc(sc_name)
        abi = tmp["abi"]
        bytecode = tmp["bytecode"]
        deployed_address = None
        if ("networks" in tmp) and (self.ethereum_default_network_id is not None):
            if self.ethereum_default_network_id in tmp["networks"]:
                deployed_address = tmp["networks"][str(self.ethereum_default_network_id)]["address"]
        return abi, bytecode, deployed_address


class Peggy1EnvCtx:
    def __init__(self, cmd, w3_conn, w3_tx, abi_provider, operator, sifnode_url, sifnode_chain_id, rowan_source):
        self.cmd = cmd
        self.w3_conn = w3_conn
        self.w3_tx = w3_tx
        self.abi_provider = abi_provider
        self.operator = operator
        self.sifnode_url = sifnode_url
        self.sifnode_chain_id = sifnode_chain_id
        self.rowan_source = rowan_source
        self.generic_erc20_contract = "SifchainTestToken"  # TODO Cleanup + consolidate
        self.available_test_eth_accounts = None

    def advance_block_w3(self, number):
        for _ in range(number):
            self.w3_conn.provider.make_request("evm_mine", [])

    def advance_block_truffle(self, number):
        args = ["npx", "truffle", "exec", "scripts/advanceBlock.js", str(number)]
        self.cmd.execst(args, cwd=main.project_dir("smart-contracts"))

    def advance_block(self, number):
        self.advance_block_truffle(number)

    def advance_blocks(self):
        if self.w3_tx.is_local_node:
            self.advance_block(50)
        # Otherwise just wait

    def get_blocklist_sc(self):
        address, abi = self.abi_provider.get_sc_abi("Blocklist")
        result = self.w3_conn.eth.contract(address=address, abi=abi)
        return result

    def get_bridge_bank_sc(self):
        address, abi = self.abi_provider.get_sc_abi("BridgeBank")
        # assert address == test_utilities.get_required_env_var("BRIDGE_BANK_ADDRESS")
        result = self.w3_conn.eth.contract(address=address, abi=abi)
        return result

    def get_bridge_token_sc(self, address=None):
        _address, abi = self.abi_provider.get_sc_abi("BridgeToken")
        return self.w3_conn.eth.contract(address=address, abi=abi)

    def get_generic_erc20_sc(self, address):
        return self.get_bridge_token_sc(address=address)

    def get_erc20_token_balance(self, token_addr, eth_addr):
        token_sc = self.get_generic_erc20_sc(token_addr)
        return token_sc.functions.balanceOf(eth_addr).call()

    def send_erc20_tokens(self, token_addr, from_addr, to_addr, amount):
        token_sc = self.get_generic_erc20_sc(token_addr)
        return self.w3_tx.transact_sync(token_sc.functions.transfer, from_addr)(to_addr, amount)

    # Tries to return any ether on the account to operator
    def scavenge_ether(self, account_addr):
        pass  # TODO

    # <editor-fold desc="Refactored">

    def tx_deploy(self, sc_name, deployer, constructor_args):
        abi, bytecode, _ = self.abi_provider.get_descriptor(sc_name)
        sc_json = self.abi_provider.get_compiled_sc(sc_name)
        token_sc = self.w3_conn.eth.contract(abi=sc_json["abi"], bytecode=sc_json["bytecode"])
        return self.w3_tx.transact(token_sc.constructor, deployer, tx_opts={"from": deployer})(*constructor_args)

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

    def tx_deploy_new_token_for_testing(self, name, symbol, decimals):
        return self.tx_deploy("SifchainTestToken", self.operator, [name, symbol, decimals])

    def tx_get_testing_token_at(self, address):
        return self.tx_get_sc_at("SifchainTestToken", address)

    def tx_testing_token_mint(self, token_sc, minter_account, amount, minted_tokens_recipient):
        return self.w3_tx.transact(token_sc.functions.mint, minter_account)(minted_tokens_recipient, amount)

    def tx_update_bridge_bank_whitelist(self, token_addr, value=True):
        bridge_bank = self.get_bridge_bank_sc()
        return self.w3_tx.transact(bridge_bank.functions.updateEthWhiteList, self.operator)(token_addr, value)

    def tx_approve(self, token_sc, from_addr, to_addr, amount):
        return self.w3_tx.transact(token_sc.functions.approve, from_addr)(to_addr, amount)

    def tx_bridge_bank_lock_eth(self, from_eth_acct, to_sif_acct, amount):
        recipient = sif_addr_to_evm_arg(to_sif_acct)
        bridge_bank = self.get_bridge_bank_sc()
        token_addr = eth.NULL_ADDRESS  # For "eth", otherwise use coin's address
        # Mandatory tx_opts: {"from": from_eth_acct, "gas": max_gas_required, "value": amount}
        # If "value" is missing, we get "call to non-contract"
        tx_opts = {"value": amount}
        return self.w3_tx.transact(bridge_bank.functions.lock, from_eth_acct, tx_opts=tx_opts)(recipient, token_addr, amount)

    def tx_bridge_bank_lock_erc20(self, token_addr, from_eth_acct, to_sif_acct, amount):
        recipient = sif_addr_to_evm_arg(to_sif_acct)
        bridge_bank = self.get_bridge_bank_sc()
        # When transfering ERC20, the amount needs to be passed as argument, and the "message.value" should be 0
        tx_opts = {"value": 0}
        return self.w3_tx.transact(bridge_bank.functions.lock, from_eth_acct, tx_opts=tx_opts)(recipient, token_addr, amount)

    def tx_approve_and_lock(self, token_sc, from_eth_acct, to_sif_acct, amount):
        bridge_bank_sc = self.get_bridge_bank_sc()
        txhash1 = self.tx_approve(token_sc, self.operator, bridge_bank_sc.address, amount)
        txhash2 = self.tx_bridge_bank_lock_erc20(token_sc.address, from_eth_acct, to_sif_acct, amount)
        return txhash1, txhash2

    # </editor-fold>

    # TODO Merge with test_ofac_blocklist.py and move into standalone module
    def deploy_new_token(self, erc20_contract_name, deployer_addr, constructor_args):
        new_token_sc_json = self.abi_provider.get_compiled_sc(erc20_contract_name)
        new_token_abi = new_token_sc_json["abi"]
        new_token_bytecode = new_token_sc_json["bytecode"]
        new_token_sc = self.w3_conn.eth.contract(abi=new_token_abi, bytecode=new_token_bytecode)
        txrcpt = self.w3_tx.transact_sync(new_token_sc.constructor, deployer_addr)(*constructor_args)
        new_token_sc = self.w3_conn.eth.contract(abi=new_token_abi, address=txrcpt.contractAddress)
        return new_token_sc

    # Called ONLY from test_ofac_blocklist.py
    # TODO Prefer SifchainTestToken / create_new_currency
    def deploy_bridge_token_for_testing(self, token_symbol, mint_amount):
        # Get BridgeToken contract; on peggy1 branch it is already deployed by truffle migrate.
        sc_json = self.abi_provider.get_compiled_sc("BridgeToken")
        abi = sc_json["abi"]
        bytecode = sc_json["bytecode"]
        bridge_token = self.w3_tx.w3_conn.eth.contract(abi=abi, bytecode=bytecode)
        txrcpt = self.w3_tx.transact_sync(bridge_token.constructor, self.operator)(token_symbol)
        address = txrcpt.contractAddress

        bridge_token = self.w3_tx.w3_conn.eth.contract(address=address, abi=abi)
        self.w3_tx.transact_sync(bridge_token.functions.mint, self.operator)(self.operator, mint_amount)

        # assert bridge_token.functions.balanceOf(owner_address).call() == mint_amount
        assert self.get_erc20_token_balance(address, self.operator) == mint_amount
        assert bridge_token.functions.totalSupply().call() == mint_amount
        assert bridge_token.functions.symbol().call() == token_symbol
        assert bridge_token.address == address

        return bridge_token

    def deploy_new_generic_erc20_token(self, name, symbol, decimals, owner=None):
        owner = self.operator if owner is None else owner
        token_sc = self.deploy_new_token(self.generic_erc20_contract, owner, [name, symbol, decimals])
        # TODO We might want to do self.update_bridge_bank_whitelist() here too
        #      In that case, update test_integration_framework (to assert that it's whitelisted)
        assert token_sc.functions.name().call() == name
        assert token_sc.functions.symbol().call() == symbol
        assert token_sc.functions.decimals().call() == decimals
        return token_sc

    def update_bridge_bank_whitelist(self, token_addr, value):
        return self.w3_tx.wait_for_transaction_receipt(self.tx_update_bridge_bank_whitelist(token_addr, value))

    def get_whitelisted_tokens_from_bridge_bank_past_events(self):
        bridge_bank = self.get_bridge_bank_sc()
        past_events = self.smart_contract_get_past_events(bridge_bank, "LogWhiteListUpdate")
        result = {}
        for e in past_events:
            token_addr = e.args["_token"]
            value = e.args["_value"]
            assert self.w3_tx.w3_conn.toChecksumAddress(token_addr) == token_addr
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
        return ERC20TokenData("test-{}".format(id.lower()), "Test Token {}".format(id), 18)

    def get_generic_erc20_token_data(self, token_address):
        token_sc = self.get_generic_erc20_sc(token_address)
        return {
            "symbol": token_sc.functions.symbol().call(),
            "name": token_sc.functions.name().call(),
            "decimals": token_sc.functions.decimals().call()
        }

    def approve_erc20_token(self, token_sc, account_owner, amount):
        bridge_bank_sc = self.get_bridge_bank_sc()
        self.w3_tx.transact_sync(token_sc.functions.approve, account_owner)(bridge_bank_sc.address, amount)

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
        self.w3_tx.transact_sync(token_sc.functions.mint, self.operator)(minted_tokens_recipient, amount)
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
        # txhash = self.w3_tx.transact(bridge_bank.functions.lock, self.operator)(recipient, erc20_token_addr, amount)
        # # .transact({"from": from_eth_addr, "gas": max_gas_required})
        # txrcpt = self.w3_conn.eth.wait_for_transaction_receipt(txhash)
        # return txrcpt
        token_sc = self.get_generic_erc20_sc(erc20_token_addr)
        self.approve_erc20_token(token_sc, from_eth_addr, amount)
        self.bridge_bank_lock_eth(from_eth_addr, dest_sichain_addr, amount)

    def create_sifchain_addr(self, moniker=None, fund_amounts=None):
        """
        Generates a new sifchain address in test keyring. If moniker is given, uses it, otherwise
        generates a random one 'test-xxx'. If fund_amounts is given, the sifchain funds are transferred
        from rowan_source to the account before returning.
        """
        moniker = moniker or "test-" + random_string(20)
        acct = self.cmd.sifnoded_keys_add_1(moniker)
        sif_address = acct["address"]
        if fund_amounts:
            old_balances = self.get_sifchain_balance(sif_address)
            self.send_from_sifchain_to_sifchain(self.rowan_source, sif_address, fund_amounts)
            self.wait_for_sif_balance_change(sif_address, old_balances, min_changes=fund_amounts)
        return sif_address

    # smart-contracts/scripts/test/{sendLockTx.js OR sendBurnTx.js}
    # sendBurnTx is called when sifchain_symbol == "rowan", sendLockTx otherwise
    def send_from_ethereum_to_sifchain(self):
        assert False,"Not implemented yet"  # TODO

    def send_from_sifchain_to_sifchain(self, from_sif_addr, to_sif_addr, amounts):
        amounts_string = ",".join([sif_format_amount(*a) for a in amounts])
        args = ["tx", "bank", "send", from_sif_addr, to_sif_addr, amounts_string] + \
            self._sifnoded_chain_id_and_node_arg() + \
            self._sifnoded_fees_arg() + \
            ["--yes", "--output", "json"]
        res = self.cmd.sifnoded_exec(args, keyring_backend="test")
        retval = json.loads(stdout(res))
        raw_log = retval["raw_log"]
        if "insufficient funds" in raw_log:
            raise Exception(raw_log)
        return retval

    # TODO
    # def generate_test_account(self, target_ceth_balance=10**18, target_rowan_balance=10**18):
    #     sifchain_addr = self.create_sifchain_addr()
    #     self.send_eth_from_ethereum_to_sifchain(self.operator, sifchain_addr, target_ceth_balance)
    #     self.send_from_sifchain_to_sifchain(self.rowan_source, sifchain_addr, target_rowan_balance)
    #     return sifchain_addr

    def get_sifchain_balance(self, sif_addr):
        args = ["query", "bank", "balances", sif_addr, "--limit", str(100000000), "--output", "json"] + \
            self._sifnoded_chain_id_and_node_arg()
        res = self.cmd.sifnoded_exec(args)
        res = json.loads(stdout(res))["balances"]
        return dict(((x["denom"], int(x["amount"])) for x in res))

    def sif_balances_equal(self, dict1, dict2):
        d2k = set(dict2.keys())
        for k in dict1.keys():
            if (k not in dict2) or (dict1[k] != dict2[k]):
                return False
            d2k.remove(k)
        return len(d2k) == 0

    def wait_for_sif_balance_change(self, sif_addr, old_balances, min_changes=None, polling_time=1, timeout=90):
        start_time = time.time()
        result = None
        while result is None:
            new_balances = self.get_sifchain_balance(sif_addr)
            if min_changes is not None:
                have_all = True
                for amount, denom in min_changes:
                    change = new_balances.get(denom, 0) - old_balances.get(denom, 0)
                    have_all = have_all and change >= amount
                if have_all:
                    return new_balances
            else:
                if not self.sif_balances_equal(old_balances, new_balances):
                    return new_balances
            time.sleep(polling_time)
            now = time.time()
            if now - start_time > timeout:
                raise Exception("Timeout waiting for sif balance to change")

    def eth_symbol_to_sif_symbol(self, eth_token_symbol):
        # E.g. "usdt" -> "cusdt"
        if eth_token_symbol == "erowan":
            return ROWAN
        else:
            return "c" + eth_token_symbol.lower()

    # from_sif_addr has to be the address which was used at genesis time for "set-genesis-whitelister-admin".
    # You need to have its private key in the test keyring.
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
                self._sifnoded_chain_id_and_node_arg() + \
                self._sifnoded_fees_arg() + [
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

    def _sifnoded_chain_id_and_node_arg(self):
        return [] + \
            (["--node", self.sifnode_url] if self.sifnode_url else []) + \
            (["--chain-id", self.sifnode_chain_id] if self.sifnode_chain_id else [])

    def _sifnoded_fees_arg(self):
        sifnode_tx_fees = [10**17, "rowan"]
        return [
            # Deprecated: sifnoded accepts --gas-prices=0.5rowan along with --gas-adjustment=1.5 instead of a fixed fee.
            # "--gas-prices", "0.5rowan", "--gas-adjustment", "1.5",
            "--fees", sif_format_amount(*sifnode_tx_fees)]

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        # If we'running Ropsten and not local hardhat/ganache, try to reclaim any remaining eth.
        if not self.w3_tx.is_local_node:
            # self.scavenge_ether()
            pass

    def pop_test_eth_account(self):
        if self.available_test_eth_accounts is not None:
            address = self.available_test_eth_accounts.pop(0)
        else:
            # If None, we're generating non-repeatable accounts.
            address, key = self.w3_tx.create_new_eth_account()
            self.w3_tx.set_private_key(address, key)
        return address

    def create_and_fund_eth_account(self, amount_to_fund=None):
        address = self.pop_test_eth_account()
        if amount_to_fund is not None:
            balance_before = self.w3_tx.get_eth_balance(address)
            difference = amount_to_fund - balance_before
            if difference > 0:
                self.w3_tx.send_eth(self.operator, address, difference)
                assert self.w3_tx.get_eth_balance(address) == amount_to_fund
        return address

    def bridge_bank_lock_eth(self, from_eth_acct, to_sif_acct, amount):
        txhash = self.tx_bridge_bank_lock_eth(from_eth_acct, to_sif_acct, amount)
        return self.w3_tx.wait_for_transaction_receipt(txhash)

    def bridge_bank_lock_erc20(self, token_addr, from_eth_acct, to_sif_acct, amount):
        txhash = self.tx_bridge_bank_lock_erc20(token_addr, from_eth_acct, to_sif_acct, amount)
        return self.w3_tx.wait_for_transaction_receipt(txhash)

    # Peggy1-specific
    def set_ofac_blocklist_to(self, addrs):
        blocklist_sc = self.get_blocklist_sc()
        addrs = [self.w3_tx.w3_conn.toChecksumAddress(addr) for addr in addrs]
        existing_entries = blocklist_sc.functions.getFullList().call()
        to_add = [addr for addr in addrs if addr not in existing_entries]
        to_remove = [addr for addr in existing_entries if addr not in addrs]
        result = [None, None]
        if to_add:
            result[0] = self.w3_tx.transact_sync(blocklist_sc.functions.batchAddToBlocklist, self.operator)(to_add)
        if to_remove:
            result[1] = self.w3_tx.transact_sync(blocklist_sc.functions.batchRemoveFromBlocklist, self.operator)(to_remove)
        current_entries = blocklist_sc.functions.getFullList().call()
        assert set(addrs) == set(current_entries)
        return result


class ERC20TokenData:
    def __init__(self, symbol, name, decimals):
        self.symbol = symbol
        self.name = name
        self.decimals = decimals


def recover_eth_from_test_accounts():
    ctx = get_peggy1_env_ctx_test()
    w = eth.ExponentiallyWeightedAverageFeeEstimator()

    gas_price = 20 * eth.GWEI
    tx_cost = eth.MIN_TX_GAS * gas_price
    total_recovered = 0
    for addr in ctx.available_test_eth_accounts:
        balance = ctx.w3_tx.get_eth_balance(addr)
        to_recover = balance - tx_cost
        if to_recover > 0:
            log.info("Account {}: balance={}, to_recover={}".format(addr, balance//eth.GWEI, to_recover//eth.GWEI))
            ctx.w3_tx.send_eth(addr, ctx.operator, to_recover)
            total_recovered += to_recover
    log.info("Total recovered: {} ETH".format(total_recovered/eth.ETH))
