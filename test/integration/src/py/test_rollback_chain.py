import logging
import os
import time

import burn_lock_functions
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_shell_output, get_sifchain_addr_balance, \
    advance_n_ethereum_blocks, n_wait_blocks, \
    send_from_ethereum_to_sifchain
from test_utilities import wait_for_sifchain_addr_balance, \
    get_required_env_var, \
    EthereumToSifchainTransferRequest

test_integration_dir = get_required_env_var("TEST_INTEGRATION_DIR")


def test_rollback_chain(source_ethereum_address, solidity_json_path):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_account = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_account["name"]

    # Any amount will work
    amount = 11000

    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_account["address"],
        smart_contracts_dir=get_required_env_var("SMART_CONTRACTS_DIR"),
        ethereum_address=source_ethereum_address,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=amount,
        solidity_json_path=solidity_json_path
    )

    logging.info(f"create account with a balance of {request.amount}")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 50)

    new_addr = new_account["address"]

    snapshot = get_shell_output(f"{test_integration_dir}/snapshot_ganache_chain.sh")
    logging.info(f"created new account, took ganache snapshot {snapshot}")
    initial_user_balance = get_sifchain_addr_balance(new_addr, "", request.sifchain_symbol)
    logging.info(f"initial_user_balance {initial_user_balance}")

    transfer_1 = send_from_ethereum_to_sifchain(transfer_request=request)
    logging.info(f"transfer started but it will never complete (by design)")

    logging.info("advance less than wait blocks")
    advance_n_ethereum_blocks(n_wait_blocks / 2, request.smart_contracts_dir)

    # the transaction should not have happened on the sifchain side yet
    # since we haven't waited for the right number of blocks.
    # roll back ganache to the snapshot and try another transfer that
    # should succeed.

    logging.info(f"apply snapshot {snapshot} - this eliminates transfer_1 (block {transfer_1})")
    get_shell_output(f"{test_integration_dir}/apply_ganache_snapshot.sh {snapshot} 2>&1")

    logging.info("advance past block wait")
    advance_n_ethereum_blocks(n_wait_blocks * 2, request.smart_contracts_dir)
    time.sleep(5)

    second_user_balance = get_sifchain_addr_balance(new_addr, "", request.sifchain_symbol)
    if second_user_balance == initial_user_balance:
        logging.info(f"got expected outcome of no balance change @ {initial_user_balance}")
    else:
        raise Exception(
            f"balance should be the same after applying snapshot and rolling forward n_wait_blocks * 2.  initial_user_balance: {initial_user_balance} second_user_balance: {second_user_balance}"
        )

    request.amount = 10000

    logging.info(f"sending more eth: {request.amount} to {new_addr}")
    burn_lock_functions.transfer_ethereum_to_sifchain(request)

    # We want to know that ebrelayer will never do a second transaction.
    # We can't know that, so just delay a reasonable amount of time.
    logging.info("delay to give ebrelayer time to make a mistake")
    time.sleep(10)

    balance_after_sleep = get_sifchain_addr_balance(new_addr, "", request.sifchain_symbol)
    logging.info(f"get_sifchain_addr_balance after sleep is {balance_after_sleep} for {new_addr}")

    expected_balance = initial_user_balance + request.amount
    logging.info(f"look for a balance of {expected_balance}")
    wait_for_sifchain_addr_balance(new_addr, request.sifchain_symbol, expected_balance, "")
