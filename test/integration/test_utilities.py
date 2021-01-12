import json
import os
import subprocess
import sys
import time
import traceback

SIF_ETH = "ceth"
ETHEREUM_ETH = "eth"
SIF_ROWAN = "rowan"
ETHEREUM_ROWAN = "erowan"

verbose = False
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

persistantLog = open(f"{datadir}/python_commands.txt", "a")


def test_log_line(s):
    test_log(s + "\n")


def test_log(s):
    persistantLog.write(s)
    if verbose:
        print(s)


def get_shell_output(command_line):
    # we append all shell commands and output to /tmp/testrun.sh
    test_log_line("\n==========\n")
    test_log_line(command_line)
    sub = subprocess.Popen(command_line, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    subprocess_return = sub.stdout.read().rstrip().decode("utf-8")
    error_return = sub.stderr.read().rstrip().decode("utf-8")
    if error_return and error_return != "incorrect passphrase":
        print_error_message(f"error running command: {command_line}\n{error_return}")
    test_log_line(f"\n  returns:\n{subprocess_return}")
    if error_return:
        test_log_line(f"\n\nerrors:\n\n{error_return}")
    return subprocess_return


def get_shell_output_json(command_line):
    output = get_shell_output(command_line)
    if not output:
        print_error_message(f"no result returned from {command_line}")
    return json.loads(output)


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
def wait_for_balance(balance_fn, target_balance, max_attempts=30, debug_prefix=""):
    attempts = 0
    while True:
        balance = balance_fn()
        if balance == target_balance:
            return target_balance
        else:
            attempts += 1
            if attempts >= max_attempts:
                print_error_message(
                    f"{debug_prefix} Failed to get target balance of {target_balance}, balance is {balance}")
            else:
                if verbose:
                    print(
                        f"waiting for target balance {debug_prefix}: {target_balance}, current balance is {balance}, attempt {attempts}")
                time.sleep(1)


def wait_for_sifchain_balance(user, denom, network_password, target_balance, max_attempts=30):
    wait_for_balance(lambda: int(get_sifchain_balance(user, denom, network_password)), target_balance, max_attempts)


def wait_for_sifchain_addr_balance(sif_addr, denom, target_balance, max_attempts=30, debug_prefix=""):
    if not max_attempts: max_attempts = 30
    wait_for_balance(lambda: int(get_sifchain_addr_balance(sif_addr, denom)), target_balance, max_attempts,
                     debug_prefix)


def sif_tx_send(from_address, to_address, amount, currency, network_password):
    cmd = f"yes {network_password} | sifnodecli tx send {from_address} {to_address} {amount}{currency} -y"
    return get_shell_output(cmd)


def burn_peggy_coin(user, eth_user, amount):
    command_line = f"""yes {network_password} | sifnodecli tx ethbridge burn {get_user_account(moniker, network_password)} \
    {eth_user} {amount} {SIF_ETH} \
    --ethereum-chain-id=5777 \
    --home deploy/networks/validators/localnet/{moniker}/.sifnodecli/ --from={moniker} \
    --yes"""
    return get_shell_output(command_line)


# Send eth from ETHEREUM_PRIVATE_KEY to BridgeBank, lock the eth on bridgebank, ceth should end up in sifchain_user
def send_eth_lock(sifchain_user, symbol, amount):
    return send_ethereum_currency_to_sifchain_addr(get_user_account(sifchain_user, network_password), symbol, amount)


# this does not wait for the transaction to complete
def send_ethereum_currency_to_sifchain_addr(sif_addr, symbol, amount):
    command_line = f"{cd_smart_contracts_dir} yarn peggy:lock {sif_addr} {symbol} {amount}"
    return get_shell_output(command_line)


currency_pairs = {
    "eth": "ceth",
    "ceth": "eth",
    "rowan": "erowan",
    "erowan": "rowan"
}


def mirror_of(currency):
    return currency_pairs.get(currency)


def wait_for_sif_account(sif_addr, max_attempts=30):
    command = f"sifnodecli q account {sif_addr}"
    attempts = 0
    while True:
        try:
            result = get_shell_output(command)
            print(f"account {sif_addr} is now created")
            return result
        except:
            attempts += 1
            if attempts > max_attempts:
                raise Exception(f"too many attempts to get sif account {sif_addr}")
            time.sleep(1)


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
    send_ethereum_currency_to_sifchain_addr(sif_addr, ethereum_symbol, amount)
    advance_n_ethereum_blocks(n_wait_blocks)
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


def advance_n_ethereum_blocks(n=50):
    return run_yarn_command(f"{cd_smart_contracts_dir} yarn advance {n}")


def amount_in_wei(amount):
    return amount * 10 ** 18


network_definition_file_json = sys.argv[1]
if not network_definition_file_json:
    print_error_message("missing network_definition_file")

network_password = get_password(network_definition_file_json)
if not network_password:
    print_error_message(f"unable to read network password from {network_definition_file_json}")
