import base64
import json
import time
import grpc
import re
import web3
from typing import Mapping, Any, Tuple, AnyStr
from siftool import command, cosmos, eth
from siftool.common import *


log = siftool_logger(__name__)


ROWAN = "rowan"

# Sifchain public network endpoints
BETANET = {"node": "https://rpc.sifchain.finance", "chain_id": "sifchain-1"}
TESTNET = {"node": "https://rpc-testnet.sifchain.finance", "chain_id": "sifchain-testnet-1"}
DEVNET = {"node": "https://rpc-devnet.sifchain.finance", "chain_id": "sifchain-devnet-1"}

GasFees = Tuple[float, str]  # Special case of cosmos.Balance with only one denom and float amount

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


class Sifnoded:
    def __init__(self, cmd, home: Optional[str] = None):
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

    def keys_add(self, moniker: str, mnemonic: Optional[Iterable[str]] = None) -> Mapping[str, Any]:
        if mnemonic is None:
            res = self.sifnoded_exec(["keys", "add", moniker], keyring_backend=self.keyring_backend,
                 sifnoded_home=self.home, stdin=["y"])
            _unused_mnemonic = stderr(res).splitlines()[-1].split(" ")
        else:
            res = self.sifnoded_exec(["keys", "add", moniker, "--recover"], keyring_backend=self.keyring_backend,
                sifnoded_home=self.home, stdin=[" ".join(mnemonic)])
        account = exactly_one(yaml_load(stdout(res)))
        return account

    def keys_delete(self, name: str):
        self.cmd.execst(["sifnoded", "keys", "delete", name, "--keyring-backend", self.keyring_backend], stdin=["y"], check_exit=False)

    def add_genesis_account(self, sifnodeadmin_addr: cosmos.Address, tokens: cosmos.Balance):
        tokens_str = cosmos.balance_format(tokens)
        self.sifnoded_exec(["add-genesis-account", sifnodeadmin_addr, tokens_str], sifnoded_home=self.home)

    def add_genesis_validators(self, address: cosmos.Address):
        args = ["sifnoded", "add-genesis-validators", address]
        res = self.cmd.execst(args)
        return res

    # At the moment only on future/peggy2 branch, called from PeggyEnvironment
    def add_genesis_validators_peggy(self, evm_network_descriptor: int, valoper: str, validator_power: int):
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
            "--gas-adjustment", str(gas_adjustment), "--from", from_account, "--chain-id", chain_id, "--output", "json",
            "--yes"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return [json.loads(x) for x in stdout(res).splitlines()]

    def peggy2_set_cross_chain_fee(self, admin_account_address, network_id, ethereum_cross_chain_fee_token,
        cross_chain_fee_base, cross_chain_lock_fee, cross_chain_burn_fee, admin_account_name, chain_id, gas_prices,
        gas_adjustment
    ):
        args = ["tx", "ethbridge", "set-cross-chain-fee", str(network_id),
            ethereum_cross_chain_fee_token, str(cross_chain_fee_base), str(cross_chain_lock_fee),
            str(cross_chain_burn_fee), "--from", admin_account_name, "--chain-id", chain_id, "--gas-prices",
            sif_format_amount(*gas_prices), "--gas-adjustment", str(gas_adjustment), "-y"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return res

    def peggy2_update_consensus_needed(self, admin_account_address, hardhat_chain_id, chain_id, consensus_needed):
        args = ["tx", "ethbridge", "update-consensus-needed", str(hardhat_chain_id),
            str(consensus_needed), "--from", admin_account_address, "--chain-id", chain_id, "--gas-prices",
            "0.5rowan", "--gas-adjustment", "1.5", "-y"]
        res = self.sifnoded_exec(args, keyring_backend=self.keyring_backend, sifnoded_home=self.home)
        return res

    def sifnoded_start(self, tcp_url=None, minimum_gas_prices: Optional[GasFees] = None,
        log_format_json: bool = False, log_file: Optional[IO] = None
    ):
        sifnoded_exec_args = self.build_start_cmd(tcp_url=tcp_url, minimum_gas_prices=minimum_gas_prices,
            log_format_json=log_format_json)
        return self.cmd.spawn_asynchronous_process(sifnoded_exec_args, log_file=log_file)

    def build_start_cmd(self, tcp_url: Optional[str] = None, minimum_gas_prices: Optional[GasFees] = None,
        log_format_json: bool = False, trace: bool = True
    ):
        args = [self.binary, "start"] + \
            (["--trace"] if trace else []) + \
            (["--minimum-gas-prices", sif_format_amount(*minimum_gas_prices)] if minimum_gas_prices is not None else []) + \
            (["--rpc.laddr", tcp_url] if tcp_url else []) + \
            (["--log_level", "debug"] if log_format_json else []) + \
            (["--log_format", "json"] if log_format_json else []) + \
            (["--home", self.home] if self.home else [])
        return command.buildcmd(args)

    def send(self, from_sif_addr: cosmos.Address, to_sif_addr: cosmos.Address, amounts: cosmos.Balance,
        account_seq: Optional[Tuple[int, int]] = None, broadcast_mode: Optional[str] = None
    ) -> JsonDict:
        amounts = cosmos.balance_normalize(amounts)
        assert len(amounts) > 0
        # TODO Implement batching (factor it out of inflate_token and put here)
        if self.max_send_batch_size > 0:
            assert len(amounts) <= self.max_send_batch_size, \
                "Currently only up to {} balances can be send at the same time reliably.".format(self.max_send_batch_size)
        amounts_string = cosmos.balance_format(amounts)
        args = ["tx", "bank", "send", from_sif_addr, to_sif_addr, amounts_string, "--output", "json"] + \
            self._home_args() + self._keyring_backend_args() + self._chain_id_args() + self._node_args() + \
            self._fees_args() + self._account_number_and_sequence_args(account_seq) + \
            self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        retval = json.loads(stdout(res))
        # raw_log = retval["raw_log"]
        # for bad_thing in ["insufficient funds", "signature verification failed"]:
        #     if bad_thing in raw_log:
        #         raise Exception(raw_log)
        check_raw_log(retval)
        return retval

    def send_and_check(self, from_addr: cosmos.Address, to_addr: cosmos.Address, amounts: cosmos.Balance
    ) -> cosmos.Balance:
        from_balance_before = self.get_balance(from_addr)
        assert cosmos.balance_exceeds(from_balance_before, amounts), \
            "Source account has insufficient balance (excluding transaction fee)"
        to_balance_before = self.get_balance(to_addr)
        expected_balance = cosmos.balance_add(to_balance_before, amounts)
        self.send(from_addr, to_addr, amounts)
        self.wait_for_balance_change(to_addr, to_balance_before, expected_balance=expected_balance)
        from_balance_after = self.get_balance(from_addr)
        to_balance_after = self.get_balance(to_addr)
        expected_tx_fee = {ROWAN: sif_tx_fee_in_rowan}
        assert cosmos.balance_equal(cosmos.balance_sub(from_balance_before, from_balance_after, expected_tx_fee), amounts)
        assert cosmos.balance_equal(to_balance_after, expected_balance)
        return to_balance_after

    def get_balance(self, sif_addr: cosmos.Address, height: Optional[int] = None,
        disable_log: bool = False, retries_on_error: Optional[int] = None, delay_on_error: int = 3
    ) -> cosmos.Balance:
        base_args = ["query", "bank", "balances", sif_addr]
        tmp_result = self._paged_read(base_args, "balances", height=height, limit=5000, disable_log=disable_log,
            retries_on_error=retries_on_error, delay_on_error=delay_on_error)
        return {b["denom"]: int(b["amount"]) for b in tmp_result}

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
                raise SifnodedException("Timeout waiting for sif account {} balance to change ({}s)".format(sif_addr, timeout))
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
                    raise SifnodedException("Timeout waiting for sif account {} balance to change ({}s)".format(sif_addr, change_timeout))
            time.sleep(polling_time)

    # TODO Refactor - consolidate with test_inflate_tokens.py
    def send_batch(self, from_addr: cosmos.Address, to_addr: cosmos.Address, amounts: cosmos.Balance):
        to_go = list(amounts.keys())
        account_number, account_sequence = self.get_acct_seq(from_addr)
        balance_before = self.get_balance(to_addr)
        while to_go:
            cnt = len(to_go)
            if (self.max_send_batch_size > 0) and (cnt > self.max_send_batch_size):
                cnt = self.max_send_batch_size
            batch_denoms = to_go[:cnt]
            to_go = to_go[cnt:]
            batch_balance = {denom: amounts.get(denom, 0) for denom in batch_denoms}
            self.send(from_addr, to_addr, batch_balance, account_seq=(account_number, account_sequence))
            account_sequence += 1
        expected_balance = cosmos.balance_add(balance_before, amounts)
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
        return self._paged_read(["query", "clp", "pools"], "pools", height=height)

    def query_pools_sorted(self, height: Optional[int] = None) -> Mapping[str, JsonDict]:
        pools = self.query_pools(height=height)
        result = {p["external_asset"]["symbol"]: p for p in pools}
        assert len(result) == len(pools)
        return result

    def query_clp_liquidity_providers(self, denom: str, height: Optional[int] = None) -> List[JsonDict]:
        # Note: this paged result is slightly different than `query bank balances`. Here we always get "height"
        return self._paged_read(["query", "clp", "lplist", denom], "liquidity_providers", height=height)

    def _paged_read(self, base_args: List[str], result_key: str, height: Optional[int] = None,
        limit: Optional[int] = None, disable_log: bool = False, retries_on_error: Optional[int] = None,
        delay_on_error: int = 3
    ) -> List[JsonObj]:
        retries_on_error = retries_on_error if retries_on_error is not None else self.get_balance_default_retries
        all_results = []
        page_key = None
        while True:
            args = base_args + ["--output", "json"] + \
                (["--height", str(height)] if height is not None else []) + \
                (["--limit", str(limit)] if limit is not None else []) + \
                (["--page-key", page_key] if page_key is not None else []) + \
                self._home_args() + self._chain_id_args() + self._node_args()
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
                        log.error("Error reading query result for '{}' ({}), retries left: {}".format(repr(base_args), repr(e), retries_left))
                        time.sleep(delay_on_error)
            res = json.loads(stdout(res))
            next_key = res["pagination"]["next_key"]
            if next_key is not None:
                if height is None:
                    # There are more results than fit on a page. To ensure we get all balances as a consistent
                    # snapshot, retry with "--height" fixed to the current block.
                    if "height" in res:
                        # In some cases such as "query clp pools" the response already contains a "height" and we can
                        # use it without incurring a separate request.
                        height = int(res["height"])
                        log.debug("Large query result, continuing in paged mode using height of {}.".format(height))
                    else:
                        # In some cases there is no "height" in the response and we must restart the query which wastes
                        # one request. We could optimize this by starting with explicit "--height" in the first place,
                        # but the assumption is that most of results will fit on one page and that this will be faster
                        # without "--height").
                        height = self.get_current_block()
                        log.debug("Large query result, restarting in paged mode using height of {}.".format(height))
                        continue
                page_key = _b64_decode(next_key)
            chunk = res[result_key]
            all_results.extend(chunk)
            log.debug("Read {} items, all={}, next_key={}".format(len(chunk), len(all_results), next_key))
            if next_key is None:
                break
        return all_results

    def tx_clp_create_pool(self, from_addr: cosmos.Address, symbol: str, native_amount: int, external_amount: int,
        account_seq: Optional[Tuple[int, int]] = None, broadcast_mode: Optional[str] = None,
    ) -> JsonDict:
        # For more examples see ticket #2470, e.g.
        args = ["tx", "clp", "create-pool", "--from", from_addr, "--symbol", symbol, "--nativeAmount",
            str(native_amount), "--externalAmount", str(external_amount)] + self._chain_id_args() + \
            self._node_args() + self._high_gas_prices_args() + self._home_args() + self._keyring_backend_args() + \
            self._account_number_and_sequence_args(account_seq) + \
            self._broadcast_mode_args(broadcast_mode=broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    # Items: (denom, native_amount, external_amount)
    # clp_admin has to have enough balances for the (native? / external?) amounts
    def create_liquidity_pools_batch(self, clp_admin: cosmos.Address, entries: Iterable[Tuple[str, int, int]]):
        account_number, account_sequence = self.get_acct_seq(clp_admin)
        for denom, native_amount, external_amount in entries:
            res = self.tx_clp_create_pool(clp_admin, denom, native_amount, external_amount,
                account_seq=(account_number, account_sequence))
            check_raw_log(res)
            account_sequence += 1
        self.wait_for_last_transaction_to_be_mined()
        assert set(p["external_asset"]["symbol"] for p in self.query_pools()) == set(e[0] for e in entries), \
            "Failed to create one or more liquidity pools"

    def tx_clp_add_liquidity(self, from_addr: cosmos.Address, symbol: str, native_amount: int, external_amount: int, /,
        account_seq: Optional[Tuple[int, int]] = None, broadcast_mode: Optional[str] = None
    ) -> JsonDict:
        args = ["tx", "clp", "add-liquidity", "--from", from_addr, "--symbol", symbol, "--nativeAmount",
            str(native_amount), "--externalAmount", str(external_amount)] + self._home_args() + \
            self._keyring_backend_args() + self._chain_id_args() + self._node_args() + self._fees_args() + \
            self._account_number_and_sequence_args(account_seq) + \
            self._broadcast_mode_args(broadcast_mode=broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    # asymmetry: -10000 = 100% of native asset, 0 = 50% of native asset and 50% of external asset, 10000 = 100% of external asset
    # w_basis 0 = 0%, 10000 = 100%, Remove 0-100% of liquidity symmetrically for both assets of the pair
    # See https://github.com/Sifchain/sifnode/blob/master/docs/tutorials/clp%20tutorial.md
    def tx_clp_remove_liquidity(self, from_addr: cosmos.Address, w_basis: int, asymmetry: int) -> JsonDict:
        assert (w_basis >= 0) and (w_basis <= 10000)
        assert (asymmetry >= -10000) and (asymmetry <= 10000)
        args = ["tx", "clp", "remove-liquidity", "--from", from_addr, "--wBasis", int(w_basis), "--asymmetry",
            str(asymmetry)] + self._node_args() + self._chain_id_args() + self._keyring_backend_args() + \
            self._fees_args() + self._yes_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(res)
        check_raw_log(res)
        return res

    def tx_clp_unbond_liquidity(self, from_addr: cosmos.Address, symbol: str, units: int) -> JsonDict:
        args = ["tx", "clp", "unbond-liquidity", "--from", from_addr, "--symbol", symbol, "--units", str(units)] + \
            self._home_args() + self._keyring_backend_args() + self._chain_id_args() + self._node_args() + \
            self._fees_args() + self._broadcast_mode_args() + self._yes_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(stdout(res))
        check_raw_log(res)
        return res

    def tx_clp_swap(self, from_addr: cosmos.Address, sent_symbol: str, sent_amount: int, received_symbol: str,
        min_receiving_amount: int, broadcast_mode: Optional[str] = None
    ) -> JsonDict:
        args = ["tx", "clp", "swap", "--from", from_addr, "--sentSymbol", sent_symbol, "--sentAmount", str(sent_amount),
            "--receivedSymbol", received_symbol, "--minReceivingAmount", str(min_receiving_amount)] + \
            self._node_args() + self._chain_id_args() + self._home_args() + self._keyring_backend_args() + \
            self._fees_args() + self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(stdout(res))
        check_raw_log(res)
        return res

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

    def tx_margin_update_pools(self, from_addr: cosmos.Address, open_pools: Iterable[str],
        closed_pools: Iterable[str], broadcast_mode: Optional[str] = None
    ) -> JsonDict:
        with self.cmd.with_temp_file() as open_pools_file, self.cmd.with_temp_file() as closed_pools_file:
            self.cmd.write_text_file(open_pools_file, json.dumps(list(open_pools)))
            self.cmd.write_text_file(closed_pools_file, json.dumps(list(closed_pools)))
            args = ["tx", "margin", "update-pools", open_pools_file, "--closed-pools", closed_pools_file,
                "--from", from_addr] + self._home_args() + self._keyring_backend_args() + self._node_args() + \
                self._chain_id_args() + self._fees_args() + self._broadcast_mode_args(broadcast_mode) + self._yes_args()
            res = self.sifnoded_exec(args)
            res = yaml_load(stdout(res))
            check_raw_log(res)
            return res

    def query_margin_params(self, height: Optional[int] = None) -> JsonDict:
        args = ["query", "margin", "params"] + \
            (["--height", str(height)] if height is not None else []) + \
            self._node_args() + self._chain_id_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(stdout(res))
        return res

    def tx_margin_whitelist(self, from_addr: cosmos.Address, address: cosmos.Address,
        broadcast_mode: Optional[str] = None
    ) -> JsonDict:
        args = ["tx", "margin", "whitelist", address, "--from", from_addr] + self._home_args() + \
            self._keyring_backend_args() + self._node_args() + self._chain_id_args() + self._fees_args() + \
            self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(stdout(res))
        return res

    def tx_margin_open(self, from_addr: cosmos.Address, borrow_asset: str, collateral_asset: str, collateral_amount: int,
        leverage: int, position: str, broadcast_mode: Optional[str] = None
    ) -> JsonDict:
        args = ["tx", "margin", "open", "--borrow_asset", borrow_asset, "--collateral_asset", collateral_asset,
           "--collateral_amount", str(collateral_amount), "--leverage", str(leverage), "--position", position, \
            "--from", from_addr] + self._home_args() + self._keyring_backend_args() + self._node_args() + \
            self._chain_id_args() + self._fees_args() + self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(stdout(res))
        check_raw_log(res)
        return res

    def tx_margin_close(self, from_addr: cosmos.Address, id: int, broadcast_mode: Optional[str] = None) -> JsonDict:
        args = ["tx", "margin", "close", "--id", str(id), "--from", from_addr] + self._home_args() + \
            self._keyring_backend_args() + self._node_args() + self._chain_id_args() + self._fees_args() + \
            self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = self.sifnoded_exec(args)
        res = yaml_load(stdout(res))
        check_raw_log(res)
        return res

    def margin_open_simple(self, from_addr: cosmos.Address, borrow_asset: str, collateral_asset: str, collateral_amount: int,
        leverage: int, position: str
    ) -> JsonDict:
        res = self.tx_margin_open(from_addr, borrow_asset, collateral_asset, collateral_amount, leverage, position,
            broadcast_mode="block")
        mtp_open_event = exactly_one([x for x in res["events"] if x["type"] == "margin/mtp_open"])["attributes"]
        result = {_b64_decode(x["key"]): _b64_decode(x["value"]) for x in mtp_open_event}
        return result

    def query_margin_positions_for_address(self, address: cosmos.Address, height: Optional[int] = None) -> JsonDict:
        args = ["query", "margin", "positions-for-address", address] + self._node_args() + self._chain_id_args()
        res = self._paged_read(args, "mtps", height=height)
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

    def version(self) -> str:
        return exactly_one(stderr(self.sifnoded_exec(["version"])).splitlines())

    def gov_submit_software_upgrade(self, version: str, from_acct: cosmos.Address, deposit: cosmos.Balance,
        upgrade_height: int, upgrade_info: str, title: str, description: str, broadcast_mode: Optional[str] = None
    ):
        args = ["tx", "gov", "submit-proposal", "software-upgrade", version, "--from", from_acct, "--deposit",
            cosmos.balance_format(deposit), "--upgrade-height", str(upgrade_height), "--upgrade-info", upgrade_info,
            "--title", title, "--description", description] + self._home_args() +  self._keyring_backend_args() + \
            self._node_args() + self._chain_id_args() + self._fees_args() + \
            self._broadcast_mode_args(broadcast_mode=broadcast_mode) + self._yes_args()
        res = yaml_load(stdout(self.sifnoded_exec(args)))
        return res

    def query_gov_proposals(self) -> JsonObj:
        args = ["query", "gov", "proposals"] + self._node_args() + self._chain_id_args()
        # Check if there are no active proposals, in this case we don't get an empty result but an error
        res = self.sifnoded_exec(args, check_exit=False)
        if res[0] == 1:
            error_lines = stderr(res).splitlines()
            if len(error_lines) > 0:
                if error_lines[0] == "Error: no proposals found":
                    return []
        elif res[0] == 0:
            assert len(stderr(res)) == 0
            # TODO Pagination without initial block
            res = yaml_load(stdout(self.sifnoded_exec(args)))
            assert res["pagination"]["next_key"] is None
            return res["proposals"]

    def gov_vote(self, proposal_id: int, vote: bool, from_acct: cosmos.Address, broadcast_mode: Optional[str] = None):
        args = ["tx", "gov", "vote", str(proposal_id), "yes" if vote else "no", "--from", from_acct] + \
            self._home_args() + self._keyring_backend_args() + self._node_args() + self._chain_id_args() + \
            self._fees_args() + self._broadcast_mode_args(broadcast_mode) + self._yes_args()
        res = yaml_load(stdout(self.sifnoded_exec(args)))
        return res

    @contextlib.contextmanager
    def _with_temp_json_file(self, json_obj: JsonObj) -> str:
        with self.cmd.with_temp_file() as tmpfile:
            self.cmd.write_text_file(tmpfile, json.dumps(json_obj))
            yield tmpfile

    def sifnoded_exec(self, args: List[str], sifnoded_home: Optional[str] = None,
                      keyring_backend: Optional[str] = None, stdin: Union[str, bytes, Sequence[str], None] = None,
                      cwd: Optional[str] = None, disable_log: bool = False
                      ) -> command.ExecResult:
        args = [self.binary] + args + \
               (["--home", sifnoded_home] if sifnoded_home else []) + \
               (["--keyring-backend", keyring_backend] if keyring_backend else [])
        res = self.cmd.execst(args, stdin=stdin, cwd=cwd, disable_log=disable_log)
        return res

    def _rpc_get(self, host, port, relative_url):
        url = "http://{}:{}/{}".format(host, port, relative_url)
        return json.loads(http_get(url).decode("UTF-8"))

    def get_status(self, host, port):
        return self._rpc_get(host, port, "node_info")

    def wait_for_last_transaction_to_be_mined(self, count: int = 1, disable_log: bool = True, timeout: int = 90):
        # TODO return int(self._rpc_get(host, port, abci_info)["response"]["last_block_height"])
        def latest_block_height():
            args = ["status"]  # TODO --node
            return int(json.loads(stderr(self.sifnoded_exec(args, disable_log=disable_log)))["SyncInfo"]["latest_block_height"])
        log.debug("Waiting for last sifnode transaction to be mined...")
        start_time = time.time()
        initial_block = latest_block_height()
        while latest_block_height() < initial_block + count:
            time.sleep(1)
            if time.time() - start_time > timeout:
                raise Exception("Timeout expired while waiting for last sifnode transaction to be mined")

    def wait_up(self, host, port):
        while True:
            from urllib.error import URLError
            try:
                return self.get_status(host, port)
            except URLError:
                time.sleep(1)


# Refactoring in progress - this class is supposed to become the successor of class Sifnode.
# It wraps node, home, chain_id, fees and keyring backend
class SifnodeClient:
    def __init__(self, cmd: command.Command, ctx, node: Optional[str] = None, home:
        Optional[str] = None, chain_id: Optional[str] = None, grpc_port: Optional[int] = None
    ):
        self.cmd = cmd
        self.ctx = ctx  # TODO Remove (currently needed for cross-chain fees for Peggy2)
        self.binary = "sifnoded"
        self.node = node
        self.home = home
        self.chain_id = chain_id
        self.grpc_port = grpc_port

    def query_account(self, sif_addr):
        result = json.loads(stdout(self.sifnoded_exec(["query", "account", sif_addr, "--output", "json"])))
        return result

    def query_tx(self, tx_hash: str) -> Optional[str]:
        time.sleep(6)
        args = [
            "q", "tx", tx_hash,
        ] + self._home_args() + \
            self._chain_id_args()

        try:
            tx = self.sifnoded_exec(args)
        except Exception as e:
            not_found = re.findall(".*(.*tx \(" + tx_hash + "\) not found).*", str(e))

            if len(not_found) > 0:
                return None
            else:
                raise e

        return tx

    def query_tx_exists(self, tx_hash: str) -> bool:
        return self.query_tx(tx_hash) is not None

    def generate_sign_prophecy_tx(self, from_sif_addr: cosmos.Address, prophecy_id: str, to_eth_addr: str, signature: str) -> Mapping:
        assert on_peggy2_branch, "Only for Peggy2.0"
        assert self.ctx.eth
        eth = self.ctx.eth

        args = [
            "tx", "ethbridge", "sign",
            str(eth.ethereum_network_descriptor),
            prophecy_id, to_eth_addr, signature,
            "--from", from_sif_addr,
            "--output", "json", "-y", "--generate-only"
        ] + self._gas_prices_args() + \
            self._home_args() + \
            self._chain_id_args() + \
            self._node_args()

        res = self.sifnoded_exec(args)
        result = json.loads(stdout(res))
        return result

    def send_sign_prophecy_with_wrong_signature_grpc(self, from_sif_addr: cosmos.Address, from_val_addr: cosmos.Address,
        wrong_from_sif_addr: cosmos.Address, prophecy_id: str, to_eth_addr: str, signature_for_sign_prophecy: str) -> bool:

        tx = self.generate_sign_prophecy_tx(from_sif_addr, to_eth_addr, prophecy_id, signature_for_sign_prophecy)

        # need update cosmos sender according to prophecy message
        tx['body']['messages'][0]['cosmos_sender'] = from_val_addr

        signed_tx = self.sign_transaction(tx, from_sif_addr)
        # replace the cosmos sender to simulate wrong signature for cosmos tx
        signed_tx['body']['messages'][0]['cosmos_sender'] = wrong_from_sif_addr
        encoded_tx = self.encode_transaction(signed_tx)
        result = str(self.broadcast_tx(encoded_tx))

        # tx_response {txhash: "C4CDD532E73F3D12335FA30C306452808BEFAC4A4226803585E738EC24D77320"}
        find_txhash = re.findall(".*\"(.*)\"", result)
        if len(find_txhash) > 0:
            # get the tx hash from result
            tx_hash = find_txhash[0]

            result = self.query_tx_exists(tx_hash)
            return result
        else:
            return False

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
                "-y"
            ] + \
            (["--generate-only"] if generate_only else []) + \
            self._gas_prices_args() + \
            self._home_args() + \
            self._chain_id_args() + \
            self._node_args() + \
            (self._keyring_backend_args() if not generate_only else [])
        res = self.sifnoded_exec(args)
        result = json.loads(stdout(res))
        if not generate_only:
            assert "failed to execute message" not in result["raw_log"]
        return result

    def tx_clp_create_pool(self, from_addr: cosmos.Address, symbol: str, native_amount: int, external_amount: int
    ) -> Mapping[str, Any]:
        # For more examples see ticket #2470, e.g.
        # sifnoded tx clp create-pool \
        #   --from $SIF_ACT \
        #   --keyring-backend test \
        #   --symbol ceth \
        #   --nativeAmount 49352380611368792060339203 \
        #   --externalAmount 1576369012576526264262 \
        #   --fees 100000000000000000rowan \
        #   --node ${SIFNODE_NODE} \
        #   --chain-id $SIFNODE_CHAIN_ID \
        #   --broadcast-mode block \
        #   -y
        args = ["tx", "clp", "create-pool", "--from", from_addr, "--symbol", symbol, "--nativeAmount",
            str(native_amount), "--externalAmount", str(external_amount), "--yes"] + \
            self._chain_id_args() + \
            self._node_args() + \
            self._fees_args() + \
            self._home_args() + \
            self._keyring_backend_args() + \
            self._broadcast_mode_args("block")
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def tx_clp_add_liquidity(self, from_addr: cosmos.Address, symbol: str, native_amount: int, external_amount: int
    ) -> Mapping[str, Any]:
        args = ["tx", "clp", "add-liquidity", "--from", from_addr, "--symbol", symbol, "--nativeAmount",
            str(native_amount), "--externalAmount", str(external_amount) + "--yes"] + \
            self._chain_id_args() + \
            self._node_args() + \
            self._fees_args() + \
            self._home_args() + \
            self._keyring_backend_args() + \
            self._broadcast_mode_args("block")
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def tx_clp_unbond_liquidity(self, from_addr: cosmos.Address, symbol: str, units: int):
        args = ["tx", "clp", "unbond-liquidity", "--from", from_addr, "--symbol", symbol, "--units", str(units) +
            "--yes"] + \
            self._chain_id_args() + \
            self._node_args() + \
            self._fees_args() + \
            self._home_args() + \
            self._keyring_backend_args() + \
            self._broadcast_mode_args("block")
        res = self.sifnoded_exec(args)
        return yaml_load(stdout(res))

    def tx_clp_cancel_unbound(self):
        pass  # TODO

    def tx_clp_remove_liquidity_units(self):
        pass  # TODO

    def tx_clp_swap(self):
        pass  # TODO

    def query_pools(self):
        page_key = None
        height = None
        all_pools = []
        while True:
            args = ["query", "clp", "pools"] + \
                (["--height", str(height)] if height is not None else []) + \
                (["--page-key", page_key] if page_key is not None else []) + \
                self._chain_id_args() + \
                self._node_args()
            res = self.sifnoded_exec(args)
            page_of_results = yaml_load(stdout(res))
            clp_module_address = page_of_results["clp_module_address"]  # TODO What is this? Should we return it in results?
            pools = page_of_results["pools"]
            pagination = page_of_results["pagination"]
            next_key = pagination["next_key"]
            all_pools.extend(pools)
            if next_key is None:
                break
            page_key = next_key
        return all_pools

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
                self._chain_id_args() + \
                self._node_args() + \
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
        ondemand_import_generated_protobuf_sources()
        broadcast_mode = cosmos_pb.BROADCAST_MODE_ASYNC
        with self.open_grpc_channel() as channel:
            tx_stub = cosmos_pb_grpc.ServiceStub(channel)
            req = cosmos_pb.BroadcastTxRequest(tx_bytes=encoded_tx, mode=broadcast_mode)
            resp = tx_stub.BroadcastTx(req)
            return resp

    def _gas_prices_args(self):
        return ["--gas-prices", "0.5rowan", "--gas-adjustment", "1.5"]

    def _fees_args(self):
        sifnode_tx_fees = [10 ** 17, "rowan"]
        return [
            # Deprecated: sifnoded accepts --gas-prices=0.5rowan along with --gas-adjustment=1.5 instead of a fixed fee.
            # "--gas-prices", "0.5rowan", "--gas-adjustment", "1.5",
            "--fees", sif_format_amount(*sifnode_tx_fees)]

    def _chain_id_args(self):
        return ["--chain-id", self.chain_id] if self.chain_id else []

    def _node_args(self):
        return ["--node", self.node] if self.node else []

    def _keyring_backend_args(self):
        keyring_backend = self.ctx.sifnode.keyring_backend
        return ["--keyring-backend", keyring_backend] if keyring_backend else []

    def _broadcast_mode_args(self, broadcast_mode):
        return ["--broadcast-mode", broadcast_mode]

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

    def peggy2_build_ebrelayer_cmd(self, init_what: str, network_descriptor: int, tendermint_node: str,
        web3_provider: str, bridge_registry_contract_address: eth.Address, validator_moniker: str, chain_id: str,
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
            "--validator-moniker", validator_moniker,
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

def _b64_decode(s: str):
    return base64.b64decode(s).decode("UTF-8")
