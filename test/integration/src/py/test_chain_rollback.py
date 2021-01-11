import time

from test_utilities import amount_in_wei, test_integration_dir, wait_for_sifchain_addr_balance, \
    transact_ethereum_currency_to_sifchain_addr, print_error_message, get_required_env_var
from test_utilities import get_shell_output, SIF_ETH, ETHEREUM_ETH, get_sifchain_addr_balance, \
    advance_n_ethereum_blocks, n_wait_blocks, \
    user1_addr, send_ethereum_currency_to_sifchain_addr


def test_chain_rollback():
    print("########## test_chain_rollback")

    amount = amount_in_wei(1)
    snapshot = get_shell_output(f"{test_integration_dir}/snapshot_ganache_chain.sh")
    user_balance_before_tx = get_sifchain_addr_balance(user1_addr, "", SIF_ETH)
    print(f"user_balance_before_tx {user_balance_before_tx}")
    smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
    send_ethereum_currency_to_sifchain_addr(user1_addr, ETHEREUM_ETH, amount, smart_contracts_dir)

    advance_n_ethereum_blocks(n_wait_blocks / 2, smart_contracts_dir)

    # the transaction should not have happened on the sifchain side.
    # roll back ganache to the snapshot and try transfer #2, and only
    # transfer #2 should succeed.

    get_shell_output(f"{test_integration_dir}/apply_ganache_snapshot.sh {snapshot} 2>&1")
    print("snapshot applied")

    advance_n_ethereum_blocks(n_wait_blocks * 2, smart_contracts_dir)

    if get_sifchain_addr_balance(user1_addr, "", SIF_ETH) != user_balance_before_tx:
        print_error_message("balance should be the same after applying snapshot and rolling forward n_wait_blocks * 2")

    new_amount = amount + 1000

    print(f"transact_ethereum_currency_to_sifchain_addr {user1_addr} {new_amount}")
    transact_ethereum_currency_to_sifchain_addr(user1_addr, ETHEREUM_ETH, new_amount)

    # TODO we need to wait for ebrelayer directly
    time.sleep(10)

    balance_after_sleep = get_sifchain_addr_balance(user1_addr, "", SIF_ETH)
    print(f"get_sifchain_addr_balance after sleep is {balance_after_sleep} for {user1_addr}")

    expected_balance = user_balance_before_tx + new_amount
    wait_for_sifchain_addr_balance(user1_addr, SIF_ETH, expected_balance)

    print(f"test_chain_rollback complete, balance is correct at {expected_balance}")


test_chain_rollback()
