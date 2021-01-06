import json
import logging
import os
import re
import subprocess
import time

SIF_ETH = "ceth"
ETHEREUM_ETH = "eth"
SIF_ROWAN = "rowan"
ETHEREUM_ROWAN = "erowan"

n_wait_blocks = 50  # number of blocks to wait for the relayer to act


def print_error_message(error_message):
    raise Exception(error_message)


def get_required_env_var(name):
    result = os.environ.get(name)
    if not result:
        print_error_message(f"{name} env var is required")
    return result


bridge_bank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
cd_smart_contracts_dir = f"cd {smart_contracts_dir}; "
moniker = get_required_env_var("MONIKER")
owner_addr = get_required_env_var("OWNER_ADDR")
user1_addr = get_required_env_var("USER1ADDR")
test_integration_dir = get_required_env_var("TEST_INTEGRATION_DIR")
datadir = get_required_env_var("datadir")


def get_shell_output(command_line):
    logging.debug(f"execute shell command: {command_line}")
    sub = subprocess.Popen(command_line, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    subprocess_return = sub.stdout.read().rstrip().decode("utf-8")
    error_return = sub.stderr.read().rstrip().decode("utf-8")
    if error_return and error_return != "incorrect passphrase":
        print_error_message(f"error running command: {command_line}\n{error_return}")
    logging.debug(f"shell command result: {subprocess_return}")
    if error_return:
        logging.debug(f"shell command error: {error_return}")
    return subprocess_return


def get_shell_output_json(command_line):
    output = get_shell_output(command_line)
    if not output:
        print_error_message(f"no result returned from {command_line}")
    result = json.loads(output)
    logoutput = json.dumps({"command": command_line, "result": result})
    logging.debug(f"shell_json: {logoutput}")
    return result


def run_yarn_command(command_line):
    output = get_shell_output(command_line)
    if not output:
        print_error_message(f"no result returned from {command_line}")
    # the actual last line from yarn is Done in XXX, so we want the one before that
    json_line = output.split('\n')[-2]
    return json.loads(json_line)


# converts a key to a sif address.
def get_user_account(user, network_password):
    command_line = "yes " + network_password + " | sifnodecli keys show " + user + " -a"
    return get_shell_output(command_line)


def get_password(network_definition_file_json):
    command_line = f"cat {network_definition_file_json} | jq '.[0].password'"
    password = get_shell_output(command_line)
    print(f"password is {password}")
    return password


# get the balance for user in the denom currency from sifnodecli
def get_sifchain_balance(user, denom, network_password):
    sif_address = get_user_account(user, network_password)
    return get_sifchain_addr_balance(sif_address, denom)


def get_eth_balance(account, symbol, smart_contracts_dir: str):
    command_line = f"yarn --cwd {smart_contracts_dir} peggy:getTokenBalance {account} {symbol}"
    result = get_shell_output(command_line)
    lines = result.split('\n')
    for line in lines:
        balance = re.match("Eth balance for.*\((.*) Wei\).*", line)
        if balance:
            return int(balance.group(1))
    return 0


def get_sifchain_addr_balance(sifaddress, denom):
    command_line = f"sifnodecli q auth account {sifaddress} -o json"
    json_str = get_shell_output_json(command_line)
    coins = json_str["value"]["coins"]
    for coin in coins:
        if coin["denom"] == denom:
            return int(coin["amount"])
    return 0


def get_transaction_result(tx_hash):
    command_line = f"sifnodecli q tx {tx_hash} -o json"
    json_str = get_shell_output_json(command_line)
    print(json_str)


# balance_fn is a lambda that takes no arguments
# and returns a result.  Runs the function up to
# max_attempts times, or until the result is equal to target_balance
def wait_for_balance(balance_fn, target_balance, max_attempts=30, debug_prefix="") -> int:
    attempts = 0
    while True:
        balance = balance_fn()
        if balance == target_balance:
            return int(target_balance)
        else:
            attempts += 1
            if attempts >= max_attempts:
                errmsg = f"{debug_prefix} Failed to get target balance of {target_balance}, balance is {balance}"
                logging.critical(errmsg)
                raise Exception(errmsg)
            else:
                logging.debug(
                    f"waiting for target balance {debug_prefix}: {target_balance}, current balance is {balance}, attempt {attempts}"
                )
                time.sleep(1)


def wait_for_eth_balance(ethereum_address, symbol, target_balance, smart_contracts_dir: str, max_attempts=30):
    wait_for_balance(
        lambda: get_eth_balance(
            ethereum_address,
            symbol,
            smart_contracts_dir=smart_contracts_dir
        ),
        target_balance,
        max_attempts
    )


def wait_for_sifchain_balance(user, denom, network_password, target_balance, max_attempts=30):
    wait_for_balance(lambda: int(get_sifchain_balance(user, denom, network_password)), target_balance, max_attempts)


def wait_for_sifchain_addr_balance(sifchain_address, symbol, target_balance, max_attempts=30, debug_prefix=""):
    if not max_attempts: max_attempts = 30
    return wait_for_balance(lambda: int(get_sifchain_addr_balance(sifchain_address, symbol)), target_balance,
                            max_attempts,
                            debug_prefix)


def sif_tx_send(from_address, to_address, amount, currency, yes_password):
    cmd = f"yes {yes_password} | sifnodecli tx send {from_address} {to_address} {amount}{currency} -y"
    return get_shell_output(cmd)


def burn_peggy_coin(user, eth_user, amount):
    chaindir = get_required_env_var("CHAINDIR")
    command_line = f"""yes {network_password} | sifnodecli tx ethbridge burn {get_user_account(moniker, network_password)} \
    {eth_user} {amount} {SIF_ETH} \
    --ethereum-chain-id=5777 \
    --home {chaindir}/.sifnodecli/ --from={moniker} \
    --yes"""
    return get_shell_output(command_line)


# this does not wait for the transaction to complete
def send_from_sifchain_to_ethereum(
        sifchain_address,
        ethereum_address,
        amount,
        token,
        ethereum_chain_id,
        chain_id,
        keyring_password,
        from_key: str,
        homedir: str,
        keyring_backend: str
):
    yes_entry = f"yes {keyring_password} |" if keyring_password else ""
    keyring_backend_entry = f"--keyring-backend {keyring_backend}" if keyring_backend else ""
    command_line = f"{yes_entry} sifnodecli tx ethbridge burn {sifchain_address} {ethereum_address} {amount} {token} --ethereum-chain-id={ethereum_chain_id} --home {homedir} --from={from_key} {keyring_backend_entry} --chain-id {chain_id} --yes -o json"
    return get_shell_output_json(command_line)


# this does not wait for the transaction to complete
def send_ethereum_currency_to_sifchain_addr(sif_addr, symbol, amount, smart_contracts_dir: str):
    command_line = f"yarn --cwd {smart_contracts_dir} peggy:lock {sif_addr} {symbol} {amount}"
    return run_yarn_command(command_line)


currency_pairs = {
    "eth": "ceth",
    "ceth": "eth",
    "rowan": "erowan",
    "erowan": "rowan"
}


def mirror_of(currency):
    return currency_pairs.get(currency)


def wait_for_sif_account(sif_addr, max_attempts=30):
    def fn():
        try:
            get_sifchain_addr_balance(sif_addr, "eth")
            return True
        except:
            return False

    wait_for_predicate(lambda: fn(), True, max_attempts, f"wait for account {sif_addr}")


def wait_for_predicate(predicate, success_result, max_attempts=30, debug_prefix="") -> int:
    attempts = 0
    while True:
        if predicate():
            return success_result
        else:
            attempts += 1
            if attempts >= max_attempts:
                msg = f"{debug_prefix} wait_for_predicate failed"
                logging.debug(msg)
                raise Exception(msg)
            else:
                logging.debug(
                    f"{debug_prefix} waiting for predicate, attempt {attempts}")
                time.sleep(1)


smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")


def transact_ethereum_currency_to_sifchain_addr(sif_addr, ethereum_symbol, amount):
    sifchain_symbol = mirror_of(ethereum_symbol)
    try:
        starting_balance = get_sifchain_addr_balance(sif_addr, sifchain_symbol)
    except:
        # Sometimes we're creating an account by sending it currency for the
        # first time, so you can't get a balance.
        print("exception is OK, we are creating the account now")
        starting_balance = 0
    print(f"starting balance for {sif_addr} is {starting_balance}")
    send_ethereum_currency_to_sifchain_addr(sif_addr, ethereum_symbol, amount, smart_contracts_dir=smart_contracts_dir)
    advance_n_ethereum_blocks(n_wait_blocks, smart_contracts_dir=smart_contracts_dir)
    wait_for_sif_account(sif_addr)
    wait_for_sifchain_addr_balance(sif_addr, sifchain_symbol, starting_balance + amount, 6,
                                   f"{sif_addr} / {sifchain_symbol} / {amount}")


def ganache_transactions_json():
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


def wait_for_ethereum_block_number(required_block_number: int, smart_contracts_dir: str, max_attempts=30):
    initial_block_number = current_ethereum_block_number(smart_contracts_dir)
    logging.debug(f"wait for ethereum block {required_block_number}, current_block is {initial_block_number}")

    def fn():
        current_block_number = current_ethereum_block_number(smart_contracts_dir)
        logging.debug(
            f"wait for ethereum blocks, current_block_number is {current_block_number}, required is {required_block_number}, initial_block_number is {initial_block_number}")
        return current_block_number >= required_block_number

    result = wait_for_predicate(lambda: fn(), f"wait_for_ethereum_block_number {required_block_number}")
    return result


network_password = get_required_env_var("OWNER_PASSWORD")


def amount_in_wei(amount):
    return amount * 10 ** 18
