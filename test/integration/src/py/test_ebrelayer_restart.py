import copy
import json
import logging
import time

import pytest

import test_utilities
from pytest_utilities import create_new_sifaddr
from pytest_utilities import generate_test_account, generate_minimal_test_account
from test_utilities import EthereumToSifchainTransferRequest


def test_ebrelayer_restart(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        integration_dir,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    request, credentials = generate_minimal_test_account(
        base_transfer_request=basic_transfer_request,
        target_ceth_balance=10 ** 15
    )
    balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, "ceth")
    logging.info("restart ebrelayer normally, leaving the last block db in place")
    test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")
    test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks * 2, request.smart_contracts_dir)
    time.sleep(5)
    assert balance == test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                               "ceth")


@pytest.mark.usefixtures("ensure_relayer_restart")
def test_ethereum_transactions_with_offline_relayer(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        source_ethereum_address,
        bridgebank_address,
        integration_dir,
):
    logging.debug("need one transaction to make sure ebrelayer writes out relaydb")
    basic_transfer_request.ethereum_address = source_ethereum_address
    generate_minimal_test_account(
        base_transfer_request=basic_transfer_request,
        target_ceth_balance=100
    )

    logging.info("shut down ebrelayer")
    time.sleep(10)
    test_utilities.get_shell_output(f"pkill -9 ebrelayer || true")

    logging.info("prepare transactions to be sent while ebrelayer is offline")
    amount = 9000
    new_addresses = list(map(lambda x: create_new_sifaddr(), range(3)))
    logging.debug(f"new_addresses: {new_addresses}")
    request: EthereumToSifchainTransferRequest = copy.deepcopy(basic_transfer_request)
    requests = list(map(lambda addr: {
        "amount": amount,
        "symbol": test_utilities.NULL_ADDRESS,
        "sifchain_address": addr
    }, new_addresses))
    json_requests = json.dumps(requests)

    logging.info("send ethereum transactions while ebrelayer is offline")
    yarn_result = test_utilities.run_yarn_command(
        " ".join([
            f"yarn --cwd {smart_contracts_dir}",
            "integrationtest:sendBulkLockTx",
            f"--amount {amount}",
            f"--symbol eth",
            f"--json_path {request.solidity_json_path}",
            f"--sifchain_address {new_addresses[0]}",
            f"--transactions \'{json_requests}\'",
            f"--ethereum_address {source_ethereum_address}",
            f"--bridgebank_address {bridgebank_address}"
        ])
    )
    logging.debug(f"bulk result: {yarn_result}")
    logging.info("restart ebrelayer with outstanding locks on the ethereum side")
    test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")
    time.sleep(5)
    for _ in new_addresses:
        # ebrelayer only reads blocks if there are new blocks generated
        test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks, request.smart_contracts_dir)
    for a in new_addresses:
        test_utilities.wait_for_sif_account(a, basic_transfer_request.sifnodecli_node, 90)
        test_utilities.wait_for_sifchain_addr_balance(a, "ceth", amount, basic_transfer_request.sifnodecli_node, 90)


@pytest.mark.usefixtures("ensure_relayer_restart")
def test_sifchain_transactions_with_offline_relayer(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_credentials: test_utilities.SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source,
        smart_contracts_dir,
        source_ethereum_address,
        integration_dir,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=10 ** 19,
        target_rowan_balance=10 ** 19,
    )
    logging.info("shut down ebrelayer")
    time.sleep(10)
    test_utilities.get_shell_output(f"pkill -9 ebrelayer || true")

    logging.info("prepare transactions to be sent while ebrelayer is offline")
    amount = 9000

    new_eth_addrs = test_utilities.create_ethereum_addresses(
        smart_contracts_dir,
        basic_transfer_request.ethereum_network,
        2
    )

    request.amount = amount
    request.sifchain_symbol = "ceth"
    request.ethereum_symbol = "eth"

    logging.info("send transactions while ebrelayer is offline")

    for a in new_eth_addrs:
        request.ethereum_address = a["address"]
        sifchain_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                    "ceth")
        logging.info(f"sifchain balance is {sifchain_balance}, request is {request}")
        test_utilities.send_from_sifchain_to_ethereum(
            transfer_request=request,
            credentials=credentials
        )
        time.sleep(5)

    logging.info("restart ebrelayer")
    test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")
    time.sleep(15)
    test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks * 2, request.smart_contracts_dir)
    for a in new_eth_addrs:
        request.ethereum_address = a["address"]
        test_utilities.wait_for_eth_balance(request, amount, 600)
