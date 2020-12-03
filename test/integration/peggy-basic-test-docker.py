import subprocess
import json
import time
import sys

# define users
USER = "user1"
ROWAN = "rowan"
PEGGYETH = "ceth"
PEGGYROWAN = "erwn"
ETH = "eth"
SLEEPTIME = 5
AMOUNT = 100000
CLAIMLOCK = "lock"
CLAIMBURN = "burn"
ETHEREUM_SENDER_ADDRESS='0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9'
ETHEREUM_NULL_ADDRESS='0x0000000000000000000000000000000000000000'
ETHEREUM_CHAIN_ID='5777'
def print_error_message(error_message):
    print("#################################")
    print("!!!!Error: ", error_message)
    print("#################################")
    sys.exit(error_message)

def get_shell_output(command_line):
    sub = subprocess.Popen(command_line, shell=True, stdout=subprocess.PIPE)
    subprocess_return = sub.stdout.read()
    return subprocess_return.rstrip()

def get_password():
    command_line = "yq r network-definition.yml \"(*==$MONIKER).password\""
    output = get_shell_output(command_line).decode("utf-8")
    return f"{output}"

def get_moniker():
    command_line = "echo $MONIKER"
    return get_shell_output(command_line).decode("utf-8")

def get_ethereum_contract_address():
    command_line = "echo $ETHEREUM_CONTRACT_ADDRESS"
    return get_shell_output(command_line).decode("utf-8")

VALIDATOR = get_moniker()
ETHEREUM_CONTRACT_ADDRESS = get_ethereum_contract_address()

def get_user_account(user):
    password = get_password()
    command_line = "yes " + password + " | sifnodecli keys show " + user + " -a"
    return get_shell_output(command_line).decode("utf-8")


def get_operator_account(user):
    password = get_password()
    command_line = "yes " + password + " | sifnodecli keys show " + user + " -a --bech val"
    return get_shell_output(command_line).decode("utf-8")


def get_account_nonce(user):
    command_line = "sifnodecli q auth account " + get_user_account(user) + ' -o json'
    output = get_shell_output(command_line).decode("utf-8")
    json_str = json.loads(output)
    return json_str["value"]["sequence"]


def get_balance(user, denom):
    command_line = "sifnodecli q auth account " + get_user_account(user) + ' -o json'
    output = get_shell_output(command_line).decode("utf-8")
    json_str = json.loads(output)
    coins = json_str["value"]["coins"]
    for coin in coins:
        if coin["denom"] == denom:
            return coin["amount"]
    return 0


def create_claim(user, validator, amount, denom, claim_type):
    print(amount)
    print('----- params')
    password = get_password()
    print(password)
    print(validator)
    print(get_account_nonce(validator))
    print(get_user_account(user))
    print(get_operator_account(validator))
    print(get_ethereum_contract_address())
    print('----- params')
    print(get_password())
    command_line = f""" yes {get_password()} | sifnodecli tx ethbridge create-claim \
            {ETHEREUM_CONTRACT_ADDRESS} {get_account_nonce(validator)} {denom} \
            {ETHEREUM_SENDER_ADDRESS} {get_user_account(user)} {get_operator_account(validator)} \
            {amount} {claim_type} --token-contract-address={ETHEREUM_NULL_ADDRESS} \
            --ethereum-chain-id={ETHEREUM_CHAIN_ID} --from={validator} --yes"""
    print(command_line)
    #command_line = """yes  | sifnodecli tx ethbridge create-claim $ETHEREUM_CONTRACT_ADDRESS {} {} \
    #0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 {} {} {} {} \
    #--token-contract-address=0x0000000000000000000000000000000000000000 \
    #--ethereum-chain-id=3 --from={} \
    #--yes""".format(get_account_nonce(validator), denom, get_user_account(user),
     #               get_operator_account(validator), amount, claim_type, validator)
    # print(command_line)
    return get_shell_output(command_line)

def burn_peggy_coin(user, validator, amount):
    command_line = """yes {} | sifnodecli tx ethbridge burn {} \
    0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 {} {} \
    --ethereum-chain-id=5777 --from={} \
    --yes""".format(get_password(), get_user_account(user),
                    amount, PEGGYETH, user)
    # print(command_line)
    return get_shell_output(command_line)

def lock_rowan(user, amount):
    print('lock')
    command_line = """yes {} |sifnodecli tx ethbridge lock {} \
            0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 {} rowan \
            --ethereum-chain-id=5777 --from={} --yes    
    """.format(get_password(), get_user_account(user), amount, user)
    # print(command_line)
    return get_shell_output(command_line)

def test_case_1():
    print(
        "########## Test Case One Start: lock eth in ethereum then mint ceth in sifchain"
    )
    balance_before_tx = int(get_balance(USER, PEGGYETH))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_before_tx))
    print("Send lock claim to Sifchain...")

    print(create_claim(USER, VALIDATOR, AMOUNT, ETH, CLAIMLOCK))
    time.sleep(SLEEPTIME)
    balance_after_tx = int(get_balance(USER, PEGGYETH))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_after_tx))
    if balance_after_tx != balance_before_tx + AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    print("########## Test Case One Over ##########")

def test_case_2():
    print(
        "########## Test Case Two Start: burn ceth in sifchain then eth back to ethereum"
    )
    balance_before_tx = int(get_balance(USER, PEGGYETH))
    print('before_tx', balance_before_tx)
    print("Before burn transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_before_tx))
    if balance_before_tx < AMOUNT:
        print_error_message("No enough ceth to burn")
        return
    print("Send burn claim to Sifchain...")
    burn_peggy_coin(USER, VALIDATOR, AMOUNT)
    time.sleep(SLEEPTIME)
    balance_after_tx = int(get_balance(USER, PEGGYETH))
    print("After burn transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_after_tx))
    if balance_after_tx != balance_before_tx - AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    print("########## Test Case Two Over ##########")

def test_case_3():
    print(
        "########## Test Case Three Start: lock rowan in sifchain transfer to ethereum"
    )
    balance_before_tx = int(get_balance(USER, ROWAN))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    if balance_before_tx < AMOUNT:
        print_error_message("No enough rowan to lock")
    print("Send lock claim to Sifchain...")
    lock_rowan(USER, AMOUNT)
    time.sleep(SLEEPTIME)
    balance_after_tx = int(get_balance(USER, ROWAN))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_after_tx))
    if balance_after_tx != balance_before_tx - AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    print("########## Test Case Three Over ##########")

def test_case_4():
    print(
        "########## Test Case Four Start: burn erwn in ethereum then transfer rwn back to sifchain"
    )
    balance_before_tx = int(get_balance(USER, ROWAN))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    print("Send burn claim to Sifchain...")
    create_claim(USER, VALIDATOR, AMOUNT, ROWAN, CLAIMBURN)
    time.sleep(SLEEPTIME)
    balance_after_tx = int(get_balance(USER, ROWAN))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_after_tx))
    if balance_after_tx != balance_before_tx + AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
    print("########## Test Case Four Over ##########")

test_case_1()
test_case_2()
test_case_3()
test_case_4()
