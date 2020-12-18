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
persistantLog = open("/tmp/testrun.sh", "a")
n_wait_blocks = 50  # number of blocks to wait for the relayer to act


def print_error_message(error_message):
    print("#################################")
    print("!!!!Error: ", error_message)
    print("#################################")
    traceback.print_stack()
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


def get_password(network_definition_file):
    if not os.environ.get("MONIKER"):
        print_error_message("MONIKER environment var not set")
    f = get_shell_output(f"cat {network_definition_file}")
    command_line = f"cat {network_definition_file} | yq r - \"(*==$MONIKER).password\""
    output = get_shell_output(command_line)
    return output


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


# balance_fn is a lambda that takes no arguments
# and returns a result.  Runs the function up to
# max_attempts times, or until the result is equal to target_balance
def wait_for_balance(balance_fn, target_balance, max_attempts=30):
    attempts = 0
    while True:
        balance = balance_fn()
        if balance == target_balance:
            return target_balance
        else:
            attempts += 1
            if attempts >= max_attempts:
                print_error_message(f"Failed to get target balance of {target_balance}, balance is {balance}")
            else:
                time.sleep(1)


def wait_for_sifchain_balance(user, denom, network_password, target_balance, max_attempts=30):
    wait_for_balance(lambda: int(get_sifchain_balance(user, denom, network_password)), target_balance, max_attempts)


def wait_for_sifchain_addr_balance(sif_addr, denom, target_balance, max_attempts=30):
    wait_for_balance(lambda: int(get_sifchain_addr_balance(sif_addr, denom)), target_balance, max_attempts)


def burn_peggy_coin(user, eth_user, amount):
    command_line = f"""yes {network_password} | sifnodecli tx ethbridge burn {get_user_account(moniker, network_password)} \
    {eth_user} {amount} {SIF_ETH} \
    --ethereum-chain-id=5777 \
    --home deploy/networks/validators/localnet/{moniker}/.sifnodecli/ --from={moniker} \
    --yes"""
    return get_shell_output(command_line)


def advance_n_ethereum_blocks(n=50):
    return run_yarn_command(f"{cd_smart_contracts_dir} yarn advance {n}")


def amount_in_wei(amount):
    return amount * 10 ** 18


network_definition_file = sys.argv[1]
if not network_definition_file:
    print_error_message("missing network_definition_file")

network_password = get_password(network_definition_file)
if not network_password:
    print_error_message(f"unable to read network password from {network_definition_file}")
