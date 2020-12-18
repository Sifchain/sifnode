import time
import json

from test_utilities import print_error_message, get_user_account, get_sifchain_balance, get_shell_output, get_shell_output_json, \
    network_password, amount_in_wei, test_log_line, wait_for_sifchain_balance, get_transaction_result

# define users
USER = "user1"
ROWAN = "rowan"
PEGGYETH = "ceth"
PEGGYROWAN = "erwn"
ETH = "eth"
SLEEPTIME = 5
CLAIMLOCK = "lock"
CLAIMBURN = "burn"
ETHEREUM_SENDER_ADDRESS='0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9'
ETHEREUM_NULL_ADDRESS='0x0000000000000000000000000000000000000000'
ETHEREUM_CHAIN_ID='5777'


def get_moniker():
    command_line = "echo $MONIKER"
    return get_shell_output(command_line)


def get_ethereum_contract_address():
    command_line = "echo $ETHEREUM_CONTRACT_ADDRESS"
    return get_shell_output(command_line)


VALIDATOR = get_moniker()
ETHEREUM_CONTRACT_ADDRESS = get_ethereum_contract_address()


def get_operator_account(user):
    password = network_password
    command_line = "yes " + password + " | sifnodecli keys show " + user + " -a --bech val"
    return get_shell_output(command_line)


def get_account_nonce(user):
    command_line = "sifnodecli q auth account " + get_user_account(user, network_password) + ' -o json'
    return get_shell_output_json(command_line)["value"]["sequence"]


# sifnodecli tx ethbridge create-claim
# claim_type is lock | burn
def create_claim(user, validator, amount, denom, claim_type):
    print(amount)
    print('----- params')
    password = network_password
    print(password)
    print(validator)
    print(get_account_nonce(validator))
    print(get_user_account(user, network_password))
    print(get_operator_account(validator))
    print(get_ethereum_contract_address())
    print('----- params')
    print(network_password)
    command_line = f""" yes {network_password} | sifnodecli tx ethbridge create-claim \
            {ETHEREUM_CONTRACT_ADDRESS} {get_account_nonce(validator)} {denom} \
            {ETHEREUM_SENDER_ADDRESS} {get_user_account(user, network_password)} {get_operator_account(validator)} \
            {amount} {claim_type} --token-contract-address={ETHEREUM_NULL_ADDRESS} \
            --ethereum-chain-id={ETHEREUM_CHAIN_ID} --from={validator} --yes -o json"""
    print(command_line)
    return get_shell_output(command_line)


def burn_peggy_coin(user, validator, amount):
    command_line = """yes {} | sifnodecli tx ethbridge burn {} \
    0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 {} {} \
    --ethereum-chain-id=5777 --from={} \
    --yes -o json""".format(network_password, get_user_account(user, network_password),
                    amount, PEGGYETH, user)
    return get_shell_output(command_line)


def lock_rowan(user, amount):
    print('lock')
    command_line = """yes {} |sifnodecli tx ethbridge lock {} \
            0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 {} rowan \
            --ethereum-chain-id=5777 --from={} --yes -o json
    """.format(network_password, get_user_account(user, network_password), amount, user)
    return get_shell_output(command_line)


def test_case_1():
    test_log_line("########## Test Case One Start: lock eth in ethereum then mint ceth in sifchain\n")
    print(
        f"########## Test Case One Start: lock eth in ethereum then mint ceth in sifchain {network_password}"
    )
    balance_before_tx = int(get_sifchain_balance(USER, PEGGYETH, network_password))
    print(f"Before lock transaction {USER}'s balance of {PEGGYETH} is {balance_before_tx}")

    print("Send lock claim to Sifchain...")
    amount = amount_in_wei(5)
    tx_result = create_claim(USER, VALIDATOR, amount, ETH, CLAIMLOCK)
    tx_hash = json.loads(tx_result)["txhash"]
    time.sleep(SLEEPTIME)
    get_transaction_result(tx_hash)

    wait_for_sifchain_balance(USER, PEGGYETH, network_password, balance_before_tx + amount, 30)

    print("########## Test Case One Over ##########")
    test_log_line("########## Test Case One Over ##########\n")


def test_case_2():
    print(
        "########## Test Case Two Start: burn ceth in sifchain"
    )
    balance_before_tx = int(get_sifchain_balance(USER, PEGGYETH, network_password))
    print('before_tx', balance_before_tx)
    print("Before burn transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_before_tx))
    amount = amount_in_wei(1)
    if balance_before_tx < amount:
        print_error_message("No enough ceth to burn")
        return
    print("Send burn claim to Sifchain...")
    burn_peggy_coin(USER, VALIDATOR, amount)

    wait_for_sifchain_balance(USER, PEGGYETH, network_password, balance_before_tx - amount)

    print("########## Test Case Two Over ##########")


def test_case_3():
    print(
        "########## Test Case Three Start: lock rowan in sifchain transfer to ethereum"
    )
    balance_before_tx = int(get_sifchain_balance(USER, ROWAN, network_password))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    amount = 12
    if balance_before_tx < amount:
        print_error_message("No enough rowan to lock")
    print("Send lock claim to Sifchain...")
    lock_rowan(USER, amount)
    wait_for_sifchain_balance(USER, ROWAN, network_password, balance_before_tx - amount)
    print("########## Test Case Three Over ##########")


def test_case_4():
    print(
        "########## Test Case Four Start: burn erwn in ethereum then transfer rwn back to sifchain"
    )
    balance_before_tx = int(get_sifchain_balance(USER, ROWAN, network_password))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    print("Send burn claim to Sifchain...")
    amount = 13
    create_claim(USER, VALIDATOR, amount, ROWAN, CLAIMBURN)
    wait_for_sifchain_balance(USER, ROWAN, network_password, balance_before_tx + amount)
    print("########## Test Case Four Over ##########")


test_case_1()
test_case_2()
test_case_3()
test_case_4()
