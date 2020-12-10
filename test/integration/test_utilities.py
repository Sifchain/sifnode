import json
import os
import subprocess
import sys
import time
import traceback

persistantLog = open("/tmp/testrun.sh", "a")


def print_error_message(error_message):
    print("#################################")
    print("!!!!Error: ", error_message)
    print("#################################")
    traceback.print_stack()
    sys.exit(error_message)


bridge_bank_address = os.environ.get("BRIDGE_BANK_ADDRESS")
if not bridge_bank_address:
    print_error_message("BRIDGE_BANK_ADDRESS env var is required")

smart_contracts_dir = os.environ.get("SMART_CONTRACTS_DIR")
if not bridge_bank_address:
    print_error_message("SMART_CONTRACTS_DIR env var is required")

BASEDIR = sys.argv[0]


def test_log_line(s):
    test_log(s + "\n")


def test_log(s):
    persistantLog.write(s)


def get_shell_output(command_line):
    # we append all shell commands and output to /tmp/testrun.sh
    test_log_line("\n==========\n")
    test_log_line(command_line)
    sub = subprocess.Popen(command_line, shell=True, stdout=subprocess.PIPE)
    subprocess_return = sub.stdout.read().rstrip().decode("utf-8")
    test_log_line(f"\n  returns:\n{subprocess_return}\n\n")
    return subprocess_return


def get_shell_output_json(command_line):
    output = get_shell_output(command_line)
    if not output:
        print_error_message(f"no result returned from {command_line}")
    return json.loads(output)


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
    command_line = "sifnodecli q auth account " + get_user_account(user, network_password) + ' -o json'
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
            print(f"waiting for target balance t: {target_balance} b:{balance}")
            attempts += 1
            if attempts >= max_attempts:
                print_error_message(f"Failed to get target balance of {target_balance}, balance is {balance}")
            else:
                time.sleep(1)


def wait_for_sifchain_balance(user, denom, network_password, target_balance, max_attempts=30):
    wait_for_balance(lambda: int(get_sifchain_balance(user, denom, network_password)), target_balance, max_attempts)


def amount_in_wei(amount):
    return amount * 10 ** 18


network_definition_file = sys.argv[1]
if not network_definition_file:
    print_error_message("missing network_definition_file argument")
network_password = get_password(network_definition_file)
if not network_password:
    print_error_message(f"unable to read network password from {network_definition_file}")
