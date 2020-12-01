import subprocess
import json
import time
import re

# define users
VALIDATOR = "user1"
USER = "user2"
ROWAN = "rwn"
PEGGYETH = "ceth"
PEGGYROWAN = "erwn"
ETH = "eth"
ETH_CONTRACT = "0x0000000000000000000000000000000000000000"
SLEEPTIME = 5
AMOUNT = 10**18
ROWAN_AMOUNT = 100
CLAIMLOCK = "lock"
CLAIMBURN = "burn"

ETH_OPERATOR = "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
ETH_ACCOUNT = "0xf17f52151EbEF6C7334FAD080c5704D77216b732"
BRIDGE_CONTRACT = "0x75c35C980C0d37ef46DF04d31A140b65503c0eEd"
ROWAN_CONTRACT = "0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA"

GOTO_TESTNET_FOLDER = "cd ../smart-contracts/;\n"


def print_error_message(error_message):
    print("#################################")
    print("!!!!Error: ", error_message)
    print("#################################")


def get_shell_output(command_line):
    sub = subprocess.Popen(command_line, shell=True, stdout=subprocess.PIPE)
    subprocess_return = sub.stdout.read()
    return subprocess_return.rstrip()


def get_eth_balance(account, symbol):
    command_line = GOTO_TESTNET_FOLDER + "yarn peggy:getTokenBalance {} {}".format(
        account, symbol)
    result = get_shell_output(command_line).decode("utf-8")
    lines = result.split('\n')
    for line in lines:
        balance = re.match("Eth balance for.*\((.*) Wei\).*", line)
        if balance:
            return balance.group(1)
    return 0


def get_peggyrwn_balance(account, symbol):
    command_line = GOTO_TESTNET_FOLDER + "yarn peggy:getTokenBalance {} {}".format(
        account, symbol)
    result = get_shell_output(command_line).decode("utf-8")
    lines = result.split('\n')
    for line in lines:
        balance = re.match("Balance of eRWN for.*\((.*) eRWN.*\).*",
                           line)
        if balance:
            return balance.group(1)
    return 0


def send_eth_lock(sifchain_user, symbol, amount):
    command_line = GOTO_TESTNET_FOLDER + "yarn peggy:lock {} {} {}".format(
        get_user_account(sifchain_user), symbol, amount)
    result = get_shell_output(command_line).decode("utf-8")


def burn_peggyrwn(sifchain_user, peggyrwn_contract, amount):
    command_line = GOTO_TESTNET_FOLDER + "yarn peggy:burn {} {} {}".format(
        get_user_account(sifchain_user), peggyrwn_contract, amount)
    get_shell_output(command_line)


def get_user_account(user):
    command_line = "sifnodecli keys show " + user + " -a"
    return get_shell_output(command_line).decode("utf-8")


def get_operator_account(user):
    command_line = "sifnodecli keys show " + user + " -a --bech val"
    return get_shell_output(command_line).decode("utf-8")


def get_account_nonce(user):
    command_line = "sifnodecli q auth account " + get_user_account(user)
    output = get_shell_output(command_line).decode("utf-8")
    json_str = json.loads(output)
    return json_str["value"]["sequence"]


def get_balance(user, denom):
    command_line = "sifnodecli q auth account " + get_user_account(user)
    output = get_shell_output(command_line).decode("utf-8")
    json_str = json.loads(output)
    coins = json_str["value"]["coins"]
    for coin in coins:
        if coin["denom"] == denom:
            return coin["amount"]
    return 0


def burn_peggy_coin(user, eth_user, amount):
    command_line = """sifnodecli tx ethbridge burn {} \
    {} {} {} \
    --ethereum-chain-id=3 --from={} \
    --yes""".format(get_user_account(user), eth_user, amount, PEGGYETH, user)
    return get_shell_output(command_line)


def lock_rowan(user, eth_user, amount):
    command_line = """sifnodecli tx ethbridge lock {} \
        {} {} rwn \
        --ethereum-chain-id=3 --from={} --yes    
    """.format(get_user_account(user), eth_user, amount, user)
    return get_shell_output(command_line)


def test_case_1():
    print(
        "########## Test Case One Start: lock eth in ethereum then mint ceth in sifchain"
    )
    operator_balance_before_tx = int(get_eth_balance(ETH_OPERATOR, ETH))
    contract_balance_before_tx = int(get_eth_balance(BRIDGE_CONTRACT, ETH))
    balance_before_tx = int(get_balance(USER, PEGGYETH))
    print("Before lock transaction {}'s balance of {} is {}".format(
        ETH_OPERATOR, ETH, operator_balance_before_tx))
    print("Before lock transaction contract {}'s balance of {} is {}".format(
        BRIDGE_CONTRACT, ETH, contract_balance_before_tx))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_before_tx))
    print("Send lock claim to Sifchain...")
    if operator_balance_before_tx < AMOUNT:
        print_error_message("No enough ETH for the account to lock")
    send_eth_lock(USER, ETH, AMOUNT)
    time.sleep(SLEEPTIME)

    operator_balance_after_tx = int(get_eth_balance(ETH_OPERATOR, ETH))
    contract_balance_after_tx = int(get_eth_balance(BRIDGE_CONTRACT, ETH))
    balance_after_tx = int(get_balance(USER, PEGGYETH))
    print("After lock transaction {}'s balance of {} is {}".format(
        ETH_OPERATOR, ETH, operator_balance_after_tx))
    print("After lock transaction contract {}'s balance of {} is {}".format(
        BRIDGE_CONTRACT, ETH, contract_balance_after_tx))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_after_tx))
    if balance_after_tx != balance_before_tx + AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    if contract_balance_after_tx != contract_balance_before_tx + AMOUNT:
        print_error_message("bridge contract balance is wrong after send eth lock claim")
    print("########## Test Case One Over ##########")


def test_case_2():
    print(
        "########## Test Case Two Start: burn ceth in sifchain then eth back to ethereum"
    )
    operator_balance_before_tx = int(get_eth_balance(ETH_ACCOUNT, ETH))
    contract_balance_before_tx = int(get_eth_balance(BRIDGE_CONTRACT, ETH))
    balance_before_tx = int(get_balance(USER, PEGGYETH))
    print("Before lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, ETH, operator_balance_before_tx))
    print("Before lock transaction contract {}'s balance of {} is {}".format(
        BRIDGE_CONTRACT, ETH, contract_balance_before_tx))
    print("Before burn transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_before_tx))
    print("Send lock claim to Sifchain...")
    if balance_before_tx < AMOUNT:
        print_error_message("No enough peggyeth to burn")
        return
    burn_peggy_coin(USER, ETH_ACCOUNT, AMOUNT)
    time.sleep(SLEEPTIME)
    operator_balance_after_tx = int(get_eth_balance(ETH_ACCOUNT, ETH))
    contract_balance_after_tx = int(get_eth_balance(BRIDGE_CONTRACT, ETH))
    balance_after_tx = int(get_balance(USER, PEGGYETH))
    print("After lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, ETH, operator_balance_after_tx))
    print("After lock transaction contract {}'s balance of {} is {}".format(
        BRIDGE_CONTRACT, ETH, contract_balance_after_tx))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_after_tx))
    if balance_after_tx != balance_before_tx - AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    if contract_balance_before_tx != contract_balance_after_tx + AMOUNT:
        print_error_message("bridge contract's balance is wrong after send eth lock claim")
    print("########## Test Case Two Over ##########")


def test_case_3():
    print(
        "########## Test Case Three Start: lock rowan in sifchain transfer to ethereum"
    )
    operator_balance_before_tx = int(
        get_peggyrwn_balance(ETH_ACCOUNT, ROWAN_CONTRACT))
    print("Before lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, ROWAN_CONTRACT, operator_balance_before_tx))
    balance_before_tx = int(get_balance(USER, ROWAN))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    if balance_before_tx < ROWAN_AMOUNT:
        print_error_message("No enough rowan for lock")
        return
    print("Send lock transaction to Sifchain...")
    lock_rowan(USER, ETH_ACCOUNT, ROWAN_AMOUNT)
    time.sleep(SLEEPTIME)
    balance_after_tx = int(get_balance(USER, ROWAN))
    operator_balance_after_tx = int(
        get_peggyrwn_balance(ETH_ACCOUNT, ROWAN_CONTRACT))
    print("After lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, ETH, operator_balance_after_tx))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_after_tx))
    if balance_after_tx != balance_before_tx - ROWAN_AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    print("########## Test Case Three Over ##########")


def test_case_4():
    print(
        "########## Test Case Four Start: burn erwn in ethereum then transfer rwn back to sifchain"
    )
    operator_balance_before_tx = int(
        get_peggyrwn_balance(ETH_ACCOUNT, ROWAN_CONTRACT))
    print("Before lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, PEGGYROWAN, operator_balance_before_tx))
    balance_before_tx = int(get_balance(USER, ROWAN))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    if operator_balance_before_tx < ROWAN_AMOUNT:
        print_error_message("No enough peggyrowan to burn")
        return
    print("Send burn transaction to Sifchain...")
    burn_peggyrwn(USER, ROWAN_CONTRACT, ROWAN_AMOUNT)
    time.sleep(SLEEPTIME)
    operator_balance_after_tx = int(
        get_peggyrwn_balance(ETH_ACCOUNT, ROWAN_CONTRACT))
    print("After lock transaction operator {}'s balance of {} is {}".format(
        ETH_ACCOUNT, PEGGYROWAN, operator_balance_after_tx))
    balance_after_tx = int(get_balance(USER, ROWAN))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_after_tx))
    if balance_after_tx != balance_before_tx + ROWAN_AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
        return
    print("########## Test Case Four Over ##########")


test_case_1()
test_case_2()
test_case_3()
test_case_4()
