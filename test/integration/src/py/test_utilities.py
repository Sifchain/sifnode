import json
import logging
import os
import subprocess
from dataclasses import dataclass
from functools import lru_cache

import time

n_wait_blocks = 50  # number of blocks to wait for the relayer to act
burn_gas_cost = 65000000000 * 248692  # see x/ethbridge/types/msgs.go for gas
lock_gas_cost = 65000000000 * 282031
highest_gas_cost = max(burn_gas_cost, lock_gas_cost)

@dataclass
class EthereumToSifchainTransferRequest:
    sifchain_address: str = ""
    sifchain_destination_address: str = ""
    ethereum_address: str = ""
    ethereum_private_key_env_var: str = "not required for localnet"
    sifchain_symbol: str = "ceth"
    ethereum_symbol: str = "eth"
    ethereum_network: str = ""  # mainnet, ropsten, http:// for localnet
    amount: int = 0
    ceth_amount: int = 0
    sifchain_fees: str = ""
    smart_contracts_dir: str = ""
    ethereum_chain_id: str = "5777"
    chain_id: str = "localnet"  # cosmos chain id
    manual_block_advance: bool = True
    n_wait_blocks: int = n_wait_blocks
    bridgebank_address: str = ""
    bridgetoken_address: str = ""
    sifnodecli_node: str = "tcp://localhost:26657"
    solidity_json_path: str = ""
    # set to true if you want to fail if the balance changes before
    # the block waiting period has elapsed.  If you're runing
    # transactions in parallel, and you're doing manual block
    # advances, you can't set this to True.  You also can't set
    # this to true if the block time is really short, since
    # you may get a block advance as soon as you submit the transaction.
    check_wait_blocks: bool = False

    def as_json(self):
        return json.dumps(self.__dict__)

    @staticmethod
    def from_args(args):
        return EthereumToSifchainTransferRequest(
            sifchain_address=args.sifchain_address[0],
            sifchain_destination_address=args.sifchain_destination_address[0],
            ethereum_address=args.ethereum_address[0],
            sifchain_symbol=args.sifchain_symbol[0],
            ethereum_symbol=args.ethereum_symbol[0],
            bridgebank_address=args.bridgebank_address[0],
            amount=int(args.amount[0]),
            smart_contracts_dir=args.smart_contracts_dir[0],
            ethereum_chain_id=args.ethereum_chain_id[0],
            manual_block_advance=args.manual_block_advance,
            n_wait_blocks=int(args.n_wait_blocks[0]),
        )


@dataclass
class SifchaincliCredentials:
    keyring_passphrase: str
    keyring_backend: str
    from_key: str
    sifnodecli_homedir: str

    def printable_entries(self):
        return {**(self.__dict__), "keyring_passphrase": "** hidden **"}

    def as_json(self):
        return json.dumps(self.printable_entries())


@dataclass
class RequestAndCredentials:
    transfer_request: EthereumToSifchainTransferRequest
    credentials: SifchaincliCredentials
    args: object


SIF_ETH = "ceth"
ETHEREUM_ETH = "eth"
SIF_ROWAN = "rowan"
ETHEREUM_ROWAN = "erowan"
NULL_ADDRESS = "0x0000000000000000000000000000000000000000"


def print_error_message(error_message):
    raise Exception(error_message)


def get_required_env_var(name, why: str = "by the system"):
    result = os.environ.get(name)
    if not result:
        print_error_message(f"{name} env var is required {why}")
    return result


def get_optional_env_var(name: str, default_value: str):
    result = os.environ.get(name)
    return result if result else default_value


cmdfile = open("/tmp/testcmds.txt", "w")


def get_shell_output(command_line):
    cmdfile.write(command_line)
    if "sifnodecli" in command_line and not "q auth account" in command_line:
        time.sleep(2)
    logging.debug(f"execute shell command:\n{command_line}")
    sub = subprocess.run(command_line, shell=True, capture_output=True)
    stdout_string = sub.stdout.decode("utf-8").rstrip()
    stderr_string = sub.stderr.decode("utf-8").rstrip()
    logging.debug(f"execute shell command stdout:\n{stdout_string}\n")
    if stderr_string:
        logging.debug(f"\nexecute shell command stderr:\n{stderr_string}")
    if sub.returncode != 0:
        raise Exception(
            f"error running command: {sub.returncode} for command\n{command_line}\nstdout:\n{stdout_string}\nstderr:\n{stderr_string}")
    return stdout_string


def get_shell_output_json(command_line):
    output = get_shell_output(command_line)
    if not output:
        print_error_message(f"no result returned from {command_line}")
    try:
        result = json.loads(output)
        return result
    except:
        logging.critical(f"failed to decode json.  cmd is: {command_line}, output is: {output}")
        raise


def run_yarn_command(command_line):
    output = get_shell_output(command_line)
    if not output:
        print_error_message(f"no result returned from {command_line}")
    # If you don't use silent mode, the last line from yarn is Done in XXX,
    # so we want the one before that
    lines = output.split('\n')
    json_line = lines[-2] if lines[-1].startswith("Done in") else lines[-1]
    try:
        return json.loads(json_line)
    except Exception as e:
        raise Exception(f"json error from command:\n{command_line}\noutput:\n{lines}\noriginal exception: {e}")


# converts a key to a sif address.
def get_user_account(user, network_password):
    command_line = "yes " + network_password + " | sifnodecli keys show " + user + " -a"
    return get_shell_output(command_line)


def get_password(network_definition_file_json):
    command_line = f"cat {network_definition_file_json} | jq '.[0].password'"
    password = get_shell_output(command_line)
    print(f"password is {password}")
    return password


def get_eth_balance(transfer_request: EthereumToSifchainTransferRequest):
    network_element = f"--ethereum_network {transfer_request.ethereum_network} " if transfer_request.ethereum_network else ""
    symbol_element = f"--symbol {transfer_request.ethereum_symbol} " if transfer_request.ethereum_symbol else ""
    private_element = f"--ethereum_private_key_env_var \"{transfer_request.ethereum_private_key_env_var}\"" if transfer_request.ethereum_private_key_env_var else ""
    command_line = " ".join(
        [f"yarn -s --cwd {transfer_request.smart_contracts_dir}",
         f"integrationtest:getTokenBalance",
         f"--ethereum_address {transfer_request.ethereum_address}",
         f"--json_path {transfer_request.solidity_json_path}",
         private_element,
         symbol_element,
         network_element]
    )
    result = run_yarn_command(command_line)
    return int(result["balanceWei"])


def get_sifchain_addr_balance(sifaddress, sifnodecli_node, denom):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli q auth account {node} {sifaddress} -o json"
    json_str = get_shell_output_json(command_line)
    coins = json_str["value"]["coins"]
    for coin in coins:
        if coin["denom"] == denom:
            return int(coin["amount"])
    return 0


def wait_for_success(success_fn, max_seconds=80, debug_prefix=""):
    done_at_time = time.time() + max_seconds
    while True:
        try:
            return success_fn()
        except Exception as e:
            if time.time() >= done_at_time:
                errmsg = f"{debug_prefix} Failed to wait for success, waited for {max_seconds} seconds"
                logging.critical(errmsg)
                raise Exception(errmsg)
            else:
                logging.debug(f"waiting for success...")
                time.sleep(1)


def wait_for_successful_command(command_line, max_seconds=80):
    return wait_for_success(
        lambda: get_shell_output_json(command_line),
        max_seconds
    )


def get_transaction_result(tx_hash, sifnodecli_node, chain_id):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    chain_id_entry = f"--chain-id {chain_id}" if chain_id else ""
    command_line = f"sifnodecli q tx {node} {tx_hash} {chain_id_entry} -o json"
    json_str = wait_for_successful_command(command_line)
    return json_str


# balance_fn is a lambda that takes no arguments
# and returns a result.  Runs the function until
# max_seconds have passed, or until the result is equal to target_balance
def wait_for_balance(balance_fn, target_balance, max_seconds=80, debug_prefix="") -> int:
    done_at_time = time.time() + max_seconds
    while True:
        balance = balance_fn()
        if balance == target_balance:
            return int(target_balance)
        else:
            if time.time() >= done_at_time:
                difference = target_balance - balance
                errmsg = f"{debug_prefix} Failed to get target balance of {target_balance}, balance is {balance}, difference is {difference} ({float(difference) / 10 ** 18}), waited for {max_seconds} seconds"
                logging.critical(errmsg)
                raise Exception(errmsg)
            else:
                difference = target_balance - balance
                logging.debug(
                    f"waiting for target balance {debug_prefix}: {target_balance}, current balance is {balance}, difference is {difference} ({difference / 10 ** 18})"
                )
                time.sleep(1)


def wait_for_eth_balance(transfer_request: EthereumToSifchainTransferRequest, target_balance, max_seconds=80):
    wait_for_balance(
        lambda: get_eth_balance(transfer_request),
        int(target_balance),
        max_seconds
    )


def normalize_symbol(symbol: str):
    return symbol.lower()


def wait_for_sifchain_addr_balance(
        sifchain_address,
        symbol,
        target_balance,
        sifchaincli_node,
        max_seconds=90,
        debug_prefix=""
):
    normalized_symbol = normalize_symbol(symbol)
    if not max_seconds:
        max_seconds = 90
    logging.debug(f"wait_for_sifchain_addr_balance for node {sifchaincli_node}, {normalized_symbol}, {target_balance}")
    return wait_for_balance(
        lambda: int(get_sifchain_addr_balance(sifchain_address, sifchaincli_node, normalized_symbol)),
        int(target_balance),
        max_seconds,
        debug_prefix
    )


def detect_errors_in_sifnodecli_output(result):
    result_lines = result.split("\n")
    for line in result_lines:
        line: str
        if "rawlog:" in line and line != "rawlog: '[]'":
            raise Exception(f"should not have error in output: {result}")


def send_from_sifchain_to_sifchain_cmd(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"send_from_sifchain_to_sifchain {transfer_request}")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    home_entry = f"--home {credentials.sifnodecli_homedir}" if credentials.sifnodecli_homedir else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx send",
        transfer_request.sifchain_address,
        transfer_request.sifchain_destination_address,
        keyring_backend_entry,
        chain_id_entry,
        node,
        f"{transfer_request.amount}{transfer_request.sifchain_symbol}",
        sifchain_fees_entry,
        home_entry,
        "-y"
    ])
    return cmd


def send_from_sifchain_to_sifchain(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    cmd = send_from_sifchain_to_sifchain_cmd(transfer_request, credentials)
    result = get_shell_output(cmd)
    detect_errors_in_sifnodecli_output(result)
    return result


def send_from_sifchain_to_ethereum_cmd(transfer_request: EthereumToSifchainTransferRequest,
                                   credentials: SifchaincliCredentials):
    assert transfer_request.amount > 0
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    direction = "lock" if transfer_request.sifchain_symbol == "rowan" else "burn"
    home_entry = f"--home {credentials.sifnodecli_homedir}" if credentials.sifnodecli_homedir else ""
    from_entry = f"--from {credentials.from_key} " if credentials.from_key else ""
    command_line = f"{yes_entry} " \
                   f"sifnodecli tx ethbridge {direction} {node} " \
                   f"{transfer_request.sifchain_address} " \
                   f"{transfer_request.ethereum_address} " \
                   f"{int(transfer_request.amount):0} " \
                   f"{transfer_request.sifchain_symbol} " \
                   f"{transfer_request.ceth_amount} " \
                   f"{keyring_backend_entry} " \
                   f"{sifchain_fees_entry} " \
                   f"--ethereum-chain-id={transfer_request.ethereum_chain_id} " \
                   f"--chain-id={transfer_request.chain_id} " \
                   f"{home_entry} " \
                   f"{from_entry} " \
                   f"--yes "
    return command_line


def send_from_sifchain_to_ethereum(transfer_request: EthereumToSifchainTransferRequest,
                                   credentials: SifchaincliCredentials):
    command_line = send_from_sifchain_to_ethereum_cmd(transfer_request, credentials)
    result = get_shell_output(command_line)
    detect_errors_in_sifnodecli_output(result)
    return result


# this does not wait for the transaction to complete
def send_from_ethereum_to_sifchain(transfer_request: EthereumToSifchainTransferRequest) -> int:
    direction = "sendBurnTx" if transfer_request.sifchain_symbol == "rowan" else "sendLockTx"
    command_line = f"yarn -s --cwd {transfer_request.smart_contracts_dir} integrationtest:{direction} " \
                   f"--sifchain_address {transfer_request.sifchain_address} " \
                   f"--symbol {transfer_request.ethereum_symbol} " \
                   f"--amount {int(transfer_request.amount):0} " \
                   f"--bridgebank_address {transfer_request.bridgebank_address} " \
                   f"--ethereum_address {transfer_request.ethereum_address} " \
                   f"--ethereum_private_key_env_var \"{transfer_request.ethereum_private_key_env_var}\" " \
                   f"--json_path {transfer_request.solidity_json_path} " \
                   f"--gas estimate "
    command_line += f"--ethereum_network {transfer_request.ethereum_network} " if transfer_request.ethereum_network else ""
    transaction_result = run_yarn_command(command_line)
    if "burn" in transaction_result:
        result = transaction_result["burn"]["receipt"]["blockNumber"]
    else:
        result = transaction_result["receipt"]["blockNumber"]
    return result


def lock_rowan(user, amount):
    command_line = """yes {} |sifnodecli tx ethbridge lock {} \
            0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 {} rowan \
            --ethereum-chain-id=5777 --from={} --yes -o json
    """.format(network_password, get_user_account(user, network_password), amount, user)
    return get_shell_output(command_line)


currency_pairs = {
    "eth": "ceth",
    "ceth": "eth",
    "rowan": "erowan",
    "erowan": "rowan"
}


def mirror_of(currency):
    return currency_pairs.get(currency)


def wait_for_sif_account(sif_addr, sifchaincli_node, max_seconds=90):
    def fn():
        try:
            get_sifchain_addr_balance(sif_addr, sifchaincli_node, "eth")
            return True
        except:
            return False

    wait_for_predicate(lambda: fn(), True, max_seconds, f"wait for account {sif_addr}")


def wait_for_predicate(predicate, success_result, max_seconds=90, debug_prefix="") -> int:
    done_at_time = time.time() + max_seconds
    while True:
        if predicate():
            return success_result
        else:
            t = time.time()
            logging.debug(f"wait_for_predicate: wait for {done_at_time - t} more seconds")
            if t >= done_at_time:
                msg = f"{debug_prefix} wait_for_predicate failed"
                logging.debug(msg)
                raise Exception(msg)
            else:
                time.sleep(5)


def ganache_transactions_json():
    smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
    cd_smart_contracts_dir = f"cd {smart_contracts_dir}; "
    transaction_cmd = f"{cd_smart_contracts_dir} npx truffle exec scripts/getIntegrationTestTransactions.js  |" \
                      f" grep 'result:' | sed -e 's/.*result: //'"
    return get_shell_output(transaction_cmd)


def write_ganache_transactions_to_file(filename):
    json = ganache_transactions_json()
    with open(filename, "w") as text_file:
        print(json, file=text_file)


def advance_n_ethereum_blocks(n: int, smart_contracts_dir: str):
    return run_yarn_command(f"yarn --cwd {smart_contracts_dir} advance {int(n)}")


def current_ethereum_block_number(smart_contracts_dir: str):
    return advance_n_ethereum_blocks(0, smart_contracts_dir)["currentBlockNumber"]


def wait_for_ethereum_block_number(block_number: int, transfer_request: EthereumToSifchainTransferRequest):
    network_element = f"--ethereum_network {transfer_request.ethereum_network} " if transfer_request.ethereum_network else ""
    command_line = f"yarn --cwd {transfer_request.smart_contracts_dir} " \
                   f"integrationtest:waitForBlock " \
                   f"{network_element} " \
                   f"--block_number {block_number} "
    get_shell_output(command_line)


def amount_in_wei(amount):
    return amount * 10 ** 18


@lru_cache(maxsize=1)
def ganache_accounts(smart_contracts_dir: str):
    accounts = run_yarn_command(
        f"yarn -s --cwd {smart_contracts_dir} "
        f"integrationtest:ganacheAccounts"
    )
    return accounts


def ganache_owner_account(smart_contracts_dir: str):
    return ganache_accounts(smart_contracts_dir)["accounts"][0]


def ganache_second_account(smart_contracts_dir: str):
    """
    Returns the second ganache account.

    Useful for doing transfers so you can transfer to an
    ethereum address that doesn't have anything to do with
    paying gas fees.
    """
    return ganache_accounts(smart_contracts_dir)["accounts"][1]


def whitelist_token(token: str, smart_contracts_dir: str, setting: bool = True):
    setting = "true" if setting else "false"
    return get_shell_output(f"yarn --cwd {smart_contracts_dir} peggy:whiteList {token} {setting}")


def approve_token_amount(token_request: EthereumToSifchainTransferRequest):
    cmd = f"yarn --cwd {token_request.smart_contracts_dir} " \
          f"integrationtest:approve " \
          f"--amount {token_request.amount} " \
          f"--ethereum_address {token_request.ethereum_address} " \
          f"--spender_address {token_request.bridgebank_address} " \
          f"--symbol {token_request.ethereum_symbol} "
    return run_yarn_command(cmd)


def set_lock_burn_limit(smart_contracts_dir: str, token: str, amount: int):
    s = NULL_ADDRESS if token == "eth" else token
    cmd = f"UPDATE_ADDRESS={s} " \
          f"yarn --cwd {smart_contracts_dir} " \
          f"integrationtest:setTokenLockBurnLimit {amount}"
    return get_shell_output(cmd)


def create_ethereum_address(smart_contracts_dir: str, ethereum_network: str) -> (str, str):
    cmd = f"yarn -s --cwd {smart_contracts_dir} " \
          "integrationtest:createEthereumAddress " \
          f"--ethereum_network {ethereum_network} "
    result = run_yarn_command(cmd)
    return result["address"], result["privateKey"]


def create_ethereum_addresses(smart_contracts_dir: str, ethereum_network: str, count: int = 1):
    count_element = f"--count {count}" if count > 1 else ""
    cmd = f"yarn -s --cwd {smart_contracts_dir} " \
          "integrationtest:createEthereumAddress " \
          f"{count_element} " \
          f"--ethereum_network {ethereum_network} "
    return run_yarn_command(cmd)


def display_currency_value(x: int) -> str:
    """if x is 19 + 18 zeros, return (19000000000000000000 | 19)"""
    return f"({x} | {x / 10 ** 18})"


def create_new_currency(amount, symbol, smart_contracts_dir, bridgebank_address, solidity_json_path, operator_address = "", ethereum_network: str = ""):
    """returns {'destination': '0x627306090abaB3A6e1400e9345bC60c78a8BEf57', 'amount': '9000000000000000000', 'newtoken_address': '0x74e3FC764c2474f25369B9d021b7F92e8441A2Dc', 'newtoken_symbol': 'a3c626b'}"""
    if not operator_address:
        operator_address = ganache_owner_account(smart_contracts_dir)
    network_element = f"--ethereum_network {ethereum_network} " if ethereum_network else ""
    return run_yarn_command(
        f"yarn --cwd {smart_contracts_dir} "
        f"integrationtest:enableNewToken "
        f"--bridgebank_address {bridgebank_address} "
        f"--json_path {solidity_json_path} "
        f"--symbol {symbol} "
        f"--amount {amount} "
        f"--limit_amount {amount} "
        f"--operator_address {operator_address} "
        f"--ethereum_private_key_env_var ETHEREUM_PRIVATE_KEY "
        f"{network_element} "
    )
