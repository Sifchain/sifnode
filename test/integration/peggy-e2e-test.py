import json
import re
import time
import os

from test_utilities import get_shell_output, SIF_ETH, burn_peggy_coin, ETHEREUM_ETH, owner_addr, moniker, \
    get_sifchain_addr_balance, wait_for_sifchain_addr_balance, advance_n_ethereum_blocks, n_wait_blocks, \
    cd_smart_contracts_dir, send_eth_lock
from test_utilities import print_error_message, get_user_account, get_sifchain_balance, network_password, \
    bridge_bank_address, \
    smart_contracts_dir, wait_for_sifchain_balance, wait_for_balance

# define users
USER = "user1"
ETH_CONTRACT = "0x0000000000000000000000000000000000000000"
SLEEPTIME = 5
AMOUNT = 3 * 10 ** 18
ROWAN_AMOUNT = 5
CLAIMLOCK = "lock"
CLAIMBURN = "burn"

operatorAddress = "0xf17f52151EbEF6C7334FAD080c5704D77216b732"

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
    bridge_bank_balance_before_tx = get_eth_balance(bridge_bank_address, ETHEREUM_ETH)
    user_balance_before_tx = get_sifchain_balance(USER, SIF_ETH, network_password)

    print(f"send_eth_lock({USER}, {ETHEREUM_ETH}, {AMOUNT})")
    send_eth_lock(USER, ETHEREUM_ETH, AMOUNT)
    advance_n_ethereum_blocks(n_wait_blocks)

    wait_for_eth_balance(bridge_bank_address, ETHEREUM_ETH, bridge_bank_balance_before_tx + AMOUNT)
    wait_for_sifchain_balance(USER, SIF_ETH, network_password, user_balance_before_tx + AMOUNT)

    print("########## Test Case One Over ##########")


def test_case_2():
    print(
        "########## Test Case Two Start: ceth => eth"
    )

    # send owner ceth to operator eth
    amount = 1 * 10 ** 18

    operator_balance_before_tx = get_eth_balance(operatorAddress, ETHEREUM_ETH)
    owner_sifchain_balance_before_tx = get_sifchain_addr_balance(owner_addr, SIF_ETH)
    print(f"starting user_eth_balance_before_tx {operator_balance_before_tx}, owner_sifchain_balance_before_tx {owner_sifchain_balance_before_tx}, amount {amount}")
    burn_peggy_coin(owner_addr, operatorAddress, amount)

    wait_for_sifchain_addr_balance(owner_addr, SIF_ETH, owner_sifchain_balance_before_tx - amount)
    wait_for_eth_balance(operatorAddress, ETHEREUM_ETH, operator_balance_before_tx + amount)
    print("########## Test Case Two Over ##########")


def test_balance_does_not_change_without_manual_block_advance():
    print("########## test_balance_does_not_change_without_manual_block_advance")

    user_balance_before_tx = get_sifchain_balance(USER, SIF_ETH, network_password)
    send_eth_lock(USER, ETHEREUM_ETH, AMOUNT)

    advance_n_ethereum_blocks(n_wait_blocks / 2)

    # what we really want is to know that ebrelayer has done nothing,
    # but it's not clear how to get that, so we just wait a bit
    time.sleep(6)

    user_balance_before_required_wait = get_sifchain_balance(USER, SIF_ETH, network_password)

    print(f"Starting balance {user_balance_before_tx}, current balance {user_balance_before_required_wait} should be equal")

    if user_balance_before_required_wait != user_balance_before_tx:
        print_error_message(f"balance should not have changed yet.  Starting balance {user_balance_before_tx}, current balance {user_balance_before_required_wait}")

    advance_n_ethereum_blocks(n_wait_blocks)

    wait_for_sifchain_balance(USER, SIF_ETH, network_password, user_balance_before_tx + AMOUNT)
    print(f"final balance is {get_sifchain_balance(USER, SIF_ETH, network_password)}")


test_case_1()
test_case_2()
test_balance_does_not_change_without_manual_block_advance()
