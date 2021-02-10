import logging
import os
import time

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account, \
    get_shell_output_json, get_shell_output, detect_errors_in_sifnodecli_output


def get_faucet_balance(sifnodecli_node):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli q faucet balance {node} -o json"
    json_str = get_shell_output_json(command_line)
    print(json_str)


def get_pools(sifnodecli_node):
    node = f"--node {sifnodecli_node}" if sifnodecli_node else ""
    command_line = f"sifnodecli q clp pools {node} -o json"
    json_str = get_shell_output_json(command_line)
    print(json_str)


# sifnodecli tx clp create-pool --from sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --symbol ceth --nativeAmount 100000 --externalAmount 1000000
def create_pool(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"create_pool")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp create-pool",
        f"--from {transfer_request.sifchain_address}",
        f"--symbol {transfer_request.sifchain_symbol}",
        f"--nativeAmount {transfer_request.amount}",
        f"--externalAmount {transfer_request.amount}",
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

# sifnodecli tx clp remove-liquidity --from sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5 --symbol cacoin --wBasis 5001 --asymmetry -1 --yes --node tcp://54.218.170.168:26657 --trust-node --chain-id sandpit --gas-prices "0.5rowan"
def remove_pool_liquidity(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"remove_pool_liquidity")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp remove-liquidity",
        f"--from {transfer_request.sifchain_address}",
        f"--symbol {transfer_request.sifchain_symbol}",
        f"--wBasis 5001",
        f"--asymmetry -1",
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

# sifnodecli tx clp add-liquidity --from sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5 --symbol cacoin --nativeAmount 10000000000000000000000 --externalAmount 1000000000000000000000000 --node tcp://54.218.170.168:26657 --trust-node --chain-id sandpit --gas-prices "0.5rowan"
def add_pool_liquidity(
        transfer_request: EthereumToSifchainTransferRequest,
        credentials: SifchaincliCredentials
):
    logging.debug(f"add_pool_liquidity")
    yes_entry = f"yes {credentials.keyring_passphrase} | " if credentials.keyring_passphrase else ""
    keyring_backend_entry = f"--keyring-backend {credentials.keyring_backend}" if credentials.keyring_backend else ""
    chain_id_entry = f"--chain-id {transfer_request.chain_id}" if transfer_request.chain_id else ""
    node = f"--node {transfer_request.sifnodecli_node}" if transfer_request.sifnodecli_node else ""
    sifchain_fees_entry = f"--fees {transfer_request.sifchain_fees}" if transfer_request.sifchain_fees else ""
    cmd = " ".join([
        yes_entry,
        "sifnodecli tx clp add-liquidity",
        f"--from {transfer_request.sifchain_address}",
        f"--symbol {transfer_request.sifchain_symbol}",
        f"--nativeAmount {transfer_request.amount}",
        f"--externalAmount {transfer_request.amount}",
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


@pytest.mark.skip(reason="only test removal")
def test_add_faucet_coins(
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

    # sifaddress = rowan_source_integrationtest_env_transfer_request.sifchain_address
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


@pytest.mark.skip(reason="not now")
def test_request_faucet_coins(
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

    # sifaddress = rowan_source_integrationtest_env_transfer_request.sifchain_address
    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")

    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10 ** 17
    basic_transfer_request.sifchain_symbol = "rowan"
    basic_transfer_request.sifchain_address = sifaddress
    request_faucet_coins(basic_transfer_request, credentials)
    time.sleep(10)
    get_faucet_balance(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")


def test_pools(
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
        target_ceth_balance=10 ** 18,
        target_rowan_balance=target_rowan_balance
    )
    time.sleep(10)

    # sifaddress = rowan_source_integrationtest_env_transfer_request.sifchain_address
    sifaddress = request.sifchain_address
    from_key = credentials.from_key
    logging.info("get balances just to have those commands in the history")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")

    #get_pools(basic_transfer_request.sifnodecli_node)
    basic_transfer_request.amount = 10 ** 17
    basic_transfer_request.sifchain_symbol = "ceth"
    basic_transfer_request.sifchain_address = sifaddress

    create_pool(basic_transfer_request, credentials)
    time.sleep(10)
    get_pools(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")

    add_pool_liquidity(basic_transfer_request, credentials)
    time.sleep(10)
    get_pools(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")

    remove_pool_liquidity(basic_transfer_request, credentials)
    time.sleep(10)
    get_pools(basic_transfer_request.sifnodecli_node)
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
    balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")

    # This has problems, querying tx never returns.
    if True:
        basic_transfer_request.amount = 10 ** 19
        add_pool_liquidity(basic_transfer_request, credentials)
        time.sleep(10)
        get_pools(basic_transfer_request.sifnodecli_node)
        balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "rowan")
        balance = test_utilities.get_sifchain_addr_balance(sifaddress, basic_transfer_request.sifnodecli_node, "ceth")
