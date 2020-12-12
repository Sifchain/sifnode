import json
import json
import re
import time

from test_utilities import get_shell_output
from test_utilities import print_error_message, get_user_account, get_sifchain_balance, network_password, \
    bridge_bank_address, \
    smart_contracts_dir, wait_for_sifchain_balance, wait_for_balance

# define users
USER = "user1"
ROWAN = "rowan"
PEGGYETH = "ceth"
PEGGYROWAN = "erowan"
ETH = "eth"
ETH_CONTRACT = "0x0000000000000000000000000000000000000000"
SLEEPTIME = 5
AMOUNT = 3 * 10 ** 18
ROWAN_AMOUNT = 5
CLAIMLOCK = "lock"
CLAIMBURN = "burn"

ETH_OPERATOR = "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
ETH_ACCOUNT = "0xf17f52151EbEF6C7334FAD080c5704D77216b732"
ROWAN_CONTRACT = "0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA"

if smart_contracts_dir is None:
    print_error_message("SMART_CONTRACTS_DIR env var is required")

cd_smart_contracts_dir = f"cd {smart_contracts_dir}; "


def get_eth_balance(account, symbol):
    command_line = cd_smart_contracts_dir + "yarn peggy:getTokenBalance {} {}".format(
        account, symbol)
    result = get_shell_output(command_line)
    lines = result.split('\n')
    for line in lines:
        balance = re.match("Eth balance for.*\((.*) Wei\).*", line)
        if balance:
            return int(balance.group(1))
    return 0


def wait_for_eth_balance(account, symbol, target_balance, max_attempts=30):
    wait_for_balance(lambda: get_eth_balance(account, symbol), target_balance, max_attempts)


def get_peggyrwn_balance(account, symbol):
    command_line = cd_smart_contracts_dir + "yarn peggy:getTokenBalance {} {}".format(
        account, symbol)
    result = get_shell_output(command_line)
    lines = result.split('\n')
    for line in lines:
        balance = re.match("Balance of eRWN for.*\((.*) eRWN.*\).*",
                           line)
        if balance:
            return balance.group(1)
    return 0


# Send eth from ETHEREUM_PRIVATE_KEY to BridgeBank, lock the eth on bridgebank, ceth should end up in sifchain_user
def send_eth_lock(sifchain_user, symbol, amount):
    command_line = cd_smart_contracts_dir + "yarn peggy:lock {} {} {}".format(
        get_user_account(sifchain_user, network_password), symbol, amount)
    get_shell_output(command_line)


def burn_peggyrwn(sifchain_user, peggyrwn_contract, amount):
    command_line = cd_smart_contracts_dir + "yarn peggy:burn {} {} {}".format(
        get_user_account(sifchain_user, network_password), peggyrwn_contract, amount)
    get_shell_output(command_line)


def get_operator_account(user):
    command_line = "sifnodecli keys show " + user + " -a --bech val"
    return get_shell_output(command_line)


def get_account_nonce(user):
    command_line = "sifnodecli q auth account " + get_user_account(user, network_password)
    output = get_shell_output(command_line)
    json_str = json.loads(output)
    return json_str["value"]["sequence"]


def burn_peggy_coin(user, eth_user, amount):
    command_line = f"""yes {network_password} | sifnodecli tx ethbridge burn {get_user_account(user, network_password)} \
    {eth_user} {amount} {PEGGYETH} \
    --ethereum-chain-id=3 --from={user} \
    --yes"""
    return get_shell_output(command_line)


def lock_rowan(user, eth_user, amount):
    command_line = f"""yes {network_password} | sifnodecli tx ethbridge lock {get_user_account(user, network_password)} \
        {eth_user} {amount} rwn \
        --ethereum-chain-id=3 --from={user} --yes    
    """
    return get_shell_output(command_line)


def test_case_1():
    print(
        "########## Test Case One Start: lock eth in ethereum then mint ceth in sifchain"
    )
    bridge_bank_balance_before_tx = get_eth_balance(bridge_bank_address, ETH)
    user_balance_before_tx = get_sifchain_balance(USER, PEGGYETH, network_password)

    print(f"send_eth_lock({USER}, {ETH}, {AMOUNT})")
    send_eth_lock(USER, ETH, AMOUNT)

    wait_for_eth_balance(bridge_bank_address, ETH, bridge_bank_balance_before_tx + AMOUNT)
    wait_for_sifchain_balance(USER, PEGGYETH, network_password, user_balance_before_tx + AMOUNT)

    print("########## Test Case One Over ##########")


def test_case_2():
    print(
        "########## Test Case Two Start: burn ceth in sifchain then eth back to ethereum"
    )
    operator_balance_before_tx = int(get_eth_balance(ETH_ACCOUNT, ETH))
    contract_balance_before_tx = int(get_eth_balance(bridge_bank_address, ETH))
    balance_before_tx = int(get_sifchain_balance(USER, PEGGYETH, network_password))
    print("Before lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, ETH, operator_balance_before_tx))
    print("Before lock transaction contract {}'s balance of {} is {}".format(
        bridge_bank_address, ETH, contract_balance_before_tx))
    print("Before burn transaction {}'s balance of {} is {}".format(
        USER, PEGGYETH, balance_before_tx))
    print("Send lock claim to Sifchain...")
    if balance_before_tx < AMOUNT:
        print_error_message("No enough peggyeth to burn")
        return
    burn_peggy_coin(USER, ETH_ACCOUNT, AMOUNT)
    time.sleep(SLEEPTIME)
    operator_balance_after_tx = int(get_eth_balance(ETH_ACCOUNT, ETH))
    contract_balance_after_tx = int(get_eth_balance(bridge_bank_address, ETH))
    balance_after_tx = int(get_sifchain_balance(USER, PEGGYETH, network_password))
    print("After lock transaction {}'s balance of {} is {}".format(
        ETH_ACCOUNT, ETH, operator_balance_after_tx))
    print("After lock transaction contract {}'s balance of {} is {}".format(
        bridge_bank_address, ETH, contract_balance_after_tx))
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
    balance_before_tx = int(get_sifchain_balance(USER, ROWAN, network_password))
    print("Before lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_before_tx))
    if balance_before_tx < ROWAN_AMOUNT:
        print_error_message("No enough rowan for lock")
        return
    print("Send lock transaction to Sifchain...")
    lock_rowan(USER, ETH_ACCOUNT, ROWAN_AMOUNT)
    time.sleep(SLEEPTIME)
    balance_after_tx = int(get_sifchain_balance(USER, ROWAN, network_password))
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
    balance_before_tx = int(get_sifchain_balance(USER, ROWAN, network_password))
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
    balance_after_tx = int(get_sifchain_balance(USER, ROWAN, network_password))
    print("After lock transaction {}'s balance of {} is {}".format(
        USER, ROWAN, balance_after_tx))
    if balance_after_tx != balance_before_tx + ROWAN_AMOUNT:
        print_error_message("balance is wrong after send eth lock claim")
        return
    print("########## Test Case Four Over ##########")


test_case_1()

# TODO enable more test cases

# test_case_2()
# test_case_3()
# test_case_4()
