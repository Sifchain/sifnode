import logging
import time
import web3
import eth_typing
from hexbytes import HexBytes
from web3.types import TxReceipt
from typing import Sequence, Tuple

from siftool.common import *


ETH = 10**18
GWEI = 10**9
NULL_ADDRESS = "0x0000000000000000000000000000000000000000"
MIN_TX_GAS = 21000
Address = eth_typing.AnyAddress
PrivateKey = eth_typing.HexStr

log = siftool_logger(__name__)


def web3_ropsten_alchemy_url(alchemy_id: str) -> str:
    return "wss://eth-ropsten.alchemyapi.io/v2/{}".format(alchemy_id)

def web3_host_port_url(host: str, port: int) -> str:
    return "ws://{}:{}".format(host, port)

def web3_create_account():
    account = web3.Web3().eth.account.create()
    return account.address, account.key.hex()[2:]

def web3_connect(url: str) -> web3.Web3:
    if url.startswith("ws://") or url.startswith("wss://"):
        return web3.Web3(web3.Web3.WebsocketProvider(url, websocket_timeout=90))
    elif url.startswith("http://"):
        return web3.Web3(web3.Web3.HTTPProvider(url))
    else:
        raise Exception("Invalid web3 URL '{}', at the moment only http:// and ws:// are supported.".format(url))

def web3_inject_geth_poa_middleware(w3_conn: web3.Web3):
    # https://web3py.readthedocs.io/en/stable/middleware.html#geth-style-proof-of-authority
    from web3.middleware import geth_poa_middleware
    # inject the poa compatibility middleware to the innermost layer
    w3_conn.middleware_onion.inject(geth_poa_middleware, layer=0)
    # confirm that the connection succeeded
    # log.debug("Injected custom middleware for 'geth --dev' connection: {}".format(w3_conn.clientVersion))

def web3_wait_for_connection_up(w3_conn: web3.Web3, polling_time: int = 1, timeout: int = 90):
    start_time = time.time()
    while True:
        try:
            w3_conn.eth.block_number
            return
        except OSError:
            pass
        now = time.time()
        if now - start_time > timeout:
            raise Exception("Timeout when trying to connect to {}".format(w3_conn.provider.endpoint_uri))
        time.sleep(polling_time)

def validate_address_and_private_key(addr: Optional[Address], private_key: Optional[PrivateKey]
) -> Tuple[Address, Optional[PrivateKey]]:
    a = web3.Web3().eth.account
    addr = web3.Web3.to_checksum_address(addr) if addr else None
    if private_key:
        match_hex = re.match("^(0x)?([0-9a-fA-F]{64})$", private_key)
        private_key = match_hex[2].lower() if match_hex else _mnemonic_to_private_key(private_key)
        account = a.from_key(private_key)
        addr = addr or account.address
        assert addr == account.address, "Address does not correspond to private key"
        assert (not private_key.startswith("0x")) and (private_key == private_key.lower()), "Private key must be in lowercase hex without '0x' prefix"
    else:
        private_key = None
    return addr, private_key

class EthereumTxWrapper:
    """
    This class wraps a Web3 connection in a way that makes calling web3 functions and sending
    transactions simpler and more consistent. It avoids using features of web3 that take
    advantage of implicit accounts and private keys which are not portable between local
    (hardhat) vs. hosted (Alchemy) nodes. The recommended usage pattern is to prefer whatever
    is already in this class over writing it yourself, and for anything else to use w3_conn
    directly.
    """

    def __init__(self, w3_conn, is_local_node):
        self.w3_conn: web3.Web3 = w3_conn
        self.use_eip_1559 = True
        self.private_keys = {}
        self.default_timeout = 600

        # Differences:
        # local node (ganache, hardhat) - use sign_transaction, do not have to bid and specify gas
        # hosted node (Alchemy) - we have to sign transactions ourselves and do the bidding and fee calculation
        self.is_local_node = is_local_node
        self.is_legacy = False
        self.fixed_gas_args = None
        self.gas_estimate_fn = None
        self.used_tx_nonces = {}

        # These are only set in get_env_ctx_peggy2(), otherwise they are undefined.
        # self.cross_chain_fee_base = None
        # self.cross_chain_lock_fee = None
        # self.cross_chain_burn_fee = None
        # self.ethereum_network_descriptor = None

    def _get_private_key(self, addr):
        addr = web3.Web3.to_checksum_address(addr)
        if addr not in self.private_keys:
            raise Exception(f"No private key set for address {addr}")
        return self.private_keys[addr]

    def set_private_key(self, addr: Address, private_key: PrivateKey):
        a = web3.Web3().eth.account
        addr = web3.Web3.to_checksum_address(addr)
        if private_key is None:
            self.private_keys.pop(addr)  # Remove
        else:
            assert re.match("^([0-9a-f]{64})$", private_key)
            assert addr == a.from_key(private_key).address, f"Private key does not correspond to given address {addr}"
            self.private_keys[addr] = private_key
        if self.is_local_node:
            # existing_accounts = self.w3_conn.geth.personal.list_accounts()
            # a = self.w3_conn.eth.account.from_key(private_key)
            # # TODO This does not work, we get
            # # Error: Expected private key to be an Uint8Array with length 32
            # self.w3_conn.geth.personal.import_raw_key(private_key, "")
            pass

    # Obsolete
    def create_new_eth_account(self):
        return web3_create_account()

    # TODO This only works for local nodes (i.e. geth, ganache).
    # It does not work with hosted nodes such as Alchemy, because they don't hold users' private keys.
    def __disabled__create_eth_account_geth_personal(self, password=""):
        # This creates local account, but does not register it (w3.eth.accounts shows the same number)
        # account = w3.eth.account.create()
        # This creates account in the external node that we're connected to. The node has to support geth extensions.
        # These accounts show up in w3.eth.accounts and can be used wih transact().
        # duration must be specified because the method expects 3 parameters.
        account = self.w3_conn.geth.personal.new_account(password)
        self.w3_conn.geth.personal.unlock_account(account, password, 0)
        return account

    def get_eth_balance(self, eth_addr):
        return self.w3_conn.eth.get_balance(eth_addr)

    def _fill_in_gas(self, tx, from_addr):
        if self.fixed_gas_args:
            tx_gas_args = self.fixed_gas_args
            upfront_cost = tx_gas_args["gas"] * tx_gas_args["gasPrice"]
            balance = self.get_eth_balance(from_addr)
            difference = balance - upfront_cost
            if difference < 0:
                log.warning("Logacy transaction will likely fail: upfront_cost={}, balance={}, difference={}, transaction={}"
                .format(upfront_cost, balance, difference, repr(tx_gas_args)))
        else:
            if self.is_local_node:
                # sendTransaction() works with local node (ganache, hardhat) but not with Alchemy. From Alchemy we get an
                # error: Unsupported method: eth_sendTransaction. Alchemy does not hold users' private keys. See available
                # methods at https://docs.alchemy.com/alchemy/documentation/apis
                # There is no private key here, so for this to work "from" has to be one of "known and unlocked" accounts
                # in self.w3_conn.geth.personal.list_accounts().

                # TODO Cannot use eth.send_transaction because geth.personal.import_raw_key() seems not to work.
                #      Fall back to sign_transaction(). We only submit legacy transactions in this case.
                # txhash = self.w3_conn.eth.send_transaction(tx)

                tx_gas_args = {
                    # Transaction must include these fields: {'nonce', 'gas', 'gasPrice'}
                    "gas": 500000,
                    "gasPrice": self.w3_conn.eth.gas_price,
                }
            else:
                if self.use_eip_1559:
                    # Typical Ropsten values:
                    # max_priority_fee: 1.5 GWEI
                    # gas_price: 1.5 GWEI

                    gas, max_fee_per_gas, max_priority_fee_per_gas, gas_price = self.gas_estimate_fn(tx)

                    # For a transaction to be EIC-1559 compliant (type 0x2), remove "gasPrice" and set "maxFeePerGas" and
                    # "MaxPriorityFeePerGas"
                    # See: How to Send Transactions with EIP 1559: https://docs.alchemy.com/alchemy/guides/eip-1559/send-tx-eip-1559
                    # See: A Definitive Guide to Ethereum EIP-1559 Gas Fee Calculations: Base Fee, Priority Fee, Max Fee: https://www.blocknative.com/blog/eip-1559-fees
                    # Empirical:
                    # - gas: mandatory, must be >= 21000
                    # - maxFeePerGas: mandatory, must be >= maxPriorityFeePerGas
                    # - maxPriorityFeePerGas: mandatory
                    tx_gas_args = {
                        "gas": gas,
                        "maxFeePerGas": max_fee_per_gas,
                        "maxPriorityFeePerGas": max_priority_fee_per_gas,
                        "chainId": self.w3_conn.eth.chain_id,
                    }
                else:
                    # TODO This is experimental, do not use it
                    # gas and gasPrice are required
                    tx_gas_args = {
                        "gas": self.w3_conn.eth.estimate_gas(tx),
                        "gasPrice": self.w3_conn.eth.gas_price,
                    }

        return tx_gas_args

    def get_tx_nonce(self, addr):
        # TODO
        # We need to keep a count of nonces if we're not waiting for transaction to complete before we send the next one
        # As a limitation, this has to be shared and synchronized for anybody making transactions in the name of addr.
        if addr in self.used_tx_nonces:
            nonce = self.used_tx_nonces[addr]
        else:
            nonce = self.w3_conn.eth.get_transaction_count(addr)
        self.used_tx_nonces[addr] = nonce + 1
        return nonce

    def _send_raw_transaction(self, smart_contract_call_obj, from_addr, tx_opts=None):
        # This assumes that the one who is sending transactions (eth_addr) is not sending them
        # from anywhere else at the same time (otherwise we might get a duplicate nonce).
        # Any pending transactions with the same nonce would typically result in an error
        # "transaction replacement fee too low".
        # nonce = self.w3_conn.eth.get_transaction_count(from_addr)
        # tx_args = {
        #     # TODO For some reason we don't need to provide gas/gasPrice/maxFeePerGas/maxPriorityFeePerGas when calling
        #     #      smart contract methods. We only have to provide them for sending eth.
        #     "from": eth_addr,
        #     "nonce": nonce,
        # }
        # tx_args = dict_merge(tx_args, self._fill_in_gas(tx_args))
        tx = tx_opts or {}

        if "from" in tx:
            assert tx["from"] == from_addr

        tx = dict_merge(tx, {
            "from":  from_addr,
            "nonce": self.get_tx_nonce(from_addr)
        })

        a, b, c, d = [x in tx for x in ["gas", "gasPrice", "maxFeePerGas", "maxPriorityFeePerGas"]]
        if a and b and (not c) and (not d):
            have_valid_gas_specs = True
        elif (not a) and (not b) and c and d:
            have_valid_gas_specs = True
        elif (not a) and (not b) and (not c) and (not d):
            have_valid_gas_specs = False
        else:
            assert False, "Invalid gas specification in transaction: {}".format(tx)
        if not have_valid_gas_specs:
            tx = dict_merge(tx, self._fill_in_gas(tx, from_addr), override=False)
        else:
            assert False  # TODO At the moment there is no code that uses it so it can be taken out

        if smart_contract_call_obj is not None:
            # With no gas/gasPrice
            tx = smart_contract_call_obj.build_transaction(tx)

        private_key = self._get_private_key(from_addr)
        signed_tx = self.w3_conn.eth.account.sign_transaction(tx, private_key=private_key)
        txhash = self.w3_conn.eth.send_raw_transaction(signed_tx.rawTransaction)
        return txhash

    def wait_for_all_transaction_receipts(self, tx_hashes: Sequence[HexBytes], sleep_time: int = 5,
        timeout: Optional[int] = None
    ) -> Sequence[TxReceipt]:
        result = []
        for txhash in tx_hashes:
            txrcpt = self.wait_for_transaction_receipt(txhash, sleep_time=sleep_time, timeout=timeout)
            result.append(txrcpt)
        return result

    def wait_for_transaction_receipt(self, txhash: HexBytes, sleep_time: int = 5, timeout: Optional[int] = None
    ) -> TxReceipt:
        return self.w3_conn.eth.wait_for_transaction_receipt(txhash, timeout=timeout, poll_latency=sleep_time)

    def transact_sync(self, smart_contract_function, eth_addr, tx_opts=None, timeout=None):
        timeout = timeout if timeout is not None else self.default_timeout
        def wrapped_fn(*args, **kwargs):
            call_obj = smart_contract_function(*args, **kwargs)
            txhash = self._send_raw_transaction(call_obj, eth_addr, tx_opts=tx_opts)
            txrcpt = self.wait_for_transaction_receipt(txhash, timeout=timeout)
            return txrcpt
        return wrapped_fn

    def transact(self, smart_contract_function, eth_addr, tx_opts=None):
        def wrapped_fn(*args, **kwargs):
            call_obj = smart_contract_function(*args, **kwargs)
            txhash = self._send_raw_transaction(call_obj, eth_addr, tx_opts=tx_opts)
            return txhash
        return wrapped_fn

    def send_eth(self, from_addr: Address, to_addr: Address, amount: int) -> TxReceipt:
        log.info(f"Sending {amount} wei from {from_addr} to {to_addr}...")
        tx = {"to": to_addr, "value": amount}
        txhash = self._send_raw_transaction(None, from_addr, tx)
        txrcpt = self.wait_for_transaction_receipt(txhash)
        return txrcpt

    def advance_block_w3(self, number):
        for _ in range(number):
            # See smart-contracts/node_modules/@openzeppelin/test-helpers/src/time.js:advanceBlockTo()
            self.w3_conn.provider.make_request("evm_mine", [])

    def advance_blocks(self, number=50):
        if self.is_local_node:
            previous_block = self.w3_conn.eth.block_number
            self.advance_block_w3(number)
            assert self.w3_conn.eth.block_number - previous_block >= number
        # Otherwise do nothing (e.g. wait for balance change takes longer)

    def is_contract_logic_error(self, exception, expected_message):
        if on_peggy2_branch:
            # Hardhat
            import re
            return isinstance(exception, ValueError) and \
                len(exception.args) == 1 and \
                re.match(expected_message, exception.args[0]["message"])
        if self.is_legacy or True:
            return isinstance(exception, ValueError) and \
                len(exception.args) == 1 and \
                expected_message in exception.args[0]["message"]
        else:
            return \
                isinstance(exception, web3.exceptions.ContractLogicError) and \
                len(exception.args) == 1 and \
                expected_message in exception.args[0]

    def is_contract_logic_error_method_not_found(self, exception, method_name):
        if on_peggy2_branch:
            # Hardhat
            return self.is_contract_logic_error(exception, "Method {} not found".format(method_name))
        else:
            return self.is_contract_logic_error(exception, "not supported")

    def is_contract_logic_error_not_in_minter_role(self, exception):
        if on_peggy2_branch:
            return self.is_contract_logic_error(exception, "^Error: VM Exception while processing transaction: reverted with reason string 'AccessControl: account 0x(.{40}) is missing role 0x(.{64})'$")
        return self.is_contract_logic_error(exception, "MinterRole: caller does not have the Minter role")

    def is_contract_logic_error_amount_exceeds_balance(self, exception):
        if on_peggy2_branch:
            return self.is_contract_logic_error(exception, "^Error: VM Exception while processing transaction: reverted with reason string 'ERC20: transfer amount exceeds balance'$")
        else:
            return self.is_contract_logic_error(exception, "transfer amount exceeds balance")


class ExponentiallyWeightedAverageFeeEstimator:
    def __init__(self, w3_conn, n=10, e=0.8, k=None, percentile=60):
        self.w3_conn = w3_conn
        self.n = n
        self.e = e
        self.k = k
        self.percentile = percentile
        self.coeffs = [
            # Inputs: [1, avg_base_fee, avg_reward, max_priority_fee, gas_price, estimated_gas]
            [0, 0, 0, 0, 0, 2],  # gas returned = 2*estimated_gas
            [0, 2, 1, 0, 0, 0],  # max_fee_per_gas returned = avg_reward + 2*avg_base_fee
            [0, 0, 1, 0, 0, 0],  # max_priority_fee_per_gas returned = avg_reward
            [0, 0, 0, 0, 1, 0],  # gas_price returned = gas_price
        ]
        self.cached_data = None
        self.cached_data_timestamp = None
        self.cached_data_block_number = None
        self.cached_data_max_age = 15  # seconds

    def exp_weighted_avg(self, data):
        # cnt: number of samples, cnt >= 1
        # e: weight of last sample, 0 < e < 1 since we want to have lower weights for older blocks
        #     e = 0: first is 1 and rest are 0 (0**0 == 1)
        #     e = 1: equal weights
        # k: exponent factor, k >= 0
        #     k = 0: all weights are 1
        #     k = 1: next is previous times e
        #     k = 1/(cnt -1): first is 1 and last is e
        #     k = infinity: first is 1 and others are 0 (in this case, better set n = 1 and k = 0)
        cnt = len(data)
        k = self.k if self.k is not None else 1/(cnt - 1)
        weights = [pow(self.e, i*k) for i in range(cnt)]
        return sum([data[i] * weights[i] for i in range(cnt)]) / sum(weights)

    def _refresh_cache_if_necessary(self):
        now = time.time()
        if (self.cached_data_timestamp is None) or (now - self.cached_data_timestamp > self.cached_data_max_age):
            current_block_number = self.w3_conn.eth.block_number
            if (self.cached_data_block_number is None) or (current_block_number > self.cached_data_block_number):
                # Not all web3 implementations support eth.fee_history and eth.max_priority_fee.
                # We deterministically avoid calling those if their values are not used (i.e corresponding
                # multipliers are all 0) so that we can support different scenarios just by means of using
                # zero/nonzero multipliers. (Note: anything less than approx. 1e-324 is considered as zero)
                # Something is disabled if all the coeffs in the columns from which it is calculated are zero.
                disable_eth_fee_history, disable_eth_max_priority_fee, disable_eth_gas_price = \
                    [all([all([ci[i] == 0 for i in d]) for ci in self.coeffs]) for d in [[1, 2], [3], [4]]]

                disable_eth_fee_history = disable_eth_max_priority_fee = disable_eth_gas_price = False

                if disable_eth_fee_history:
                    avg_base_fee = 0
                    avg_reward = 0
                else:
                    fee_history = self.w3_conn.eth.fee_history(self.n - 1, "latest", [self.percentile])
                    avg_base_fee = self.exp_weighted_avg(fee_history.baseFeePerGas)
                    # TODO fee_history.reward can contain zeros. Why? Is it when empty blocks are mined? Investigate.
                    # For us, zeros will make averages wrong. So this is still indeterminate since
                    # we have no guarantee that we won't receive only zeros. Perhaps in this case we need to look at
                    # more blocks.
                    reward_history_without_zeros = [x[0] for x in fee_history.reward if x[0] > 0]
                    if len(reward_history_without_zeros) == 0:
                        log.warning("fee_history.reward contains only zeros")
                        avg_reward = 0
                    else:
                        avg_reward = self.exp_weighted_avg(reward_history_without_zeros)

                max_priority_fee = 0 if disable_eth_max_priority_fee else self.w3_conn.eth.max_priority_fee
                gas_price = 0 if disable_eth_gas_price else self.w3_conn.eth.gas_price
                self.cached_data = [avg_base_fee, avg_reward, max_priority_fee, gas_price]
                self.cached_data_timestamp = now
                self.cached_data_block_number = current_block_number
        return self.cached_data

    def estimate_fees(self, tx):
        avg_base_fee, avg_reward, max_priority_fee, gas_price = self._refresh_cache_if_necessary()
        estimated_gas = self.w3_conn.eth.estimate_gas(tx)

        vals = [1, avg_base_fee, avg_reward, max_priority_fee, gas_price, estimated_gas]

        gas, max_fee_per_gas, max_priority_fee_per_gas, gas_price = \
            [round(sum([v * coeffs[i] for i, v in enumerate(vals)])) for coeffs in self.coeffs]

        return gas, max_fee_per_gas, max_priority_fee_per_gas, gas_price

    @staticmethod
    def estimate_gas_price():
        return 0


__web3_enabled_unaudited_hdwallet_features = False

# https://stackoverflow.com/questions/68050645/how-to-create-a-web3py-account-using-mnemonic-phrase
def _mnemonic_to_private_key(mnemonic, derivation_path="m/44'/60'/0'/0/0"):
    a = web3.Web3().eth.account
    global __web3_enabled_unaudited_hdwallet_features
    if not __web3_enabled_unaudited_hdwallet_features:
        a.enable_unaudited_hdwallet_features()
        __web3_enabled_unaudited_hdwallet_features = True
    return a.from_mnemonic(mnemonic, account_path=derivation_path).privateKey.hex()[2:]
