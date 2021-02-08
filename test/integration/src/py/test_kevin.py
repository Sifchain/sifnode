import logging
import os
import time

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account, get_shell_output_json, get_shell_output, detect_errors_in_sifnodecli_output

def get_faucet_balance(sifnodecli_node):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli q faucet balance {node} -o json"
    json_str = get_shell_output_json(command_line)
    print(json_str)

def add_faucet_coins(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"add_faucet_coins {transfer_request}")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx faucet add-coins",
        f"{transfer_request.amount}{transfer_request.sifchain_symbol}",
        f"--from {transfer_request.sifchain_address}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y"
    ])
    result = get_shell_output(cmd)
    detect_errors_in_sifnodecli_output(result)
    return result

def request_faucet_coins(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"request_faucet_coins {transfer_request}")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx faucet request-coins",
        f"{transfer_request.amount}{transfer_request.sifchain_symbol}",
        f"--from {transfer_request.sifchain_address}",
        keyring_backend_entry,
        chain_id_entry,
        node,
        sifchain_fees_entry,
        f"--home {credentials.sifnodecli_homedir} ",
        "-y"
    ])
    result = get_shell_output(cmd)
    detect_errors_in_sifnodecli_output(result)
    return result

def request_faucet_coins1(sifaddress, sifnodecli_node):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli tx faucet request-coins 1rowan --from {sifaddress} {node} -o json"
    json_str = get_shell_output_json(command_line)
    print(json_str)

def test_kevin(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        bridgetoken_address,
        smart_contracts_dir
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=0,
        target_rowan_balance=target_rowan_balance
    )
    time.sleep(10)

    #sifaddress = rowan_source_integrationtest_env_transfer_request.sifchain_address
    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")

    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10 ** 17
    basic_transfer_request.sifchain_symbol = "rowan"
    basic_transfer_request.sifchain_address = sifaddress
    add_faucet_coins(basic_transfer_request, credentials)
    time.sleep(10)
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    request_faucet_coins(basic_transfer_request, credentials)
    time.sleep(10)
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
