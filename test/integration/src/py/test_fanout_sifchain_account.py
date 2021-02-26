import copy
import logging

import pytest

import burn_lock_functions
from integration_env_credentials import create_new_sifaddr_and_credentials
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def test_fanout_sifchain_account(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        n_sifchain_accounts,
        rowan_source,
        rowan_source_key,
        sifnodecli_homedir,
        ceth_amount,
        rowan_amount,
):
    logging.info("create new account")
    credentials = SifchaincliCredentials(
        keyring_backend="test",
        keyring_passphrase=None,
        from_key=None,
        sifnodecli_homedir=None
    )
    new_sifaddr, _ = create_new_sifaddr_and_credentials()

    logging.info("transfer ceth to destination")
    ceth_request = copy.deepcopy(basic_transfer_request)
    ceth_request.sifchain_address = rowan_source
    ceth_request.sifchain_destination_address = new_sifaddr
    ceth_request.amount = ceth_amount
    burn_lock_functions.transfer_sifchain_to_sifchain(ceth_request, credentials)

    logging.info("transfer rowan to destination")
    rowan_request: EthereumToSifchainTransferRequest = copy.deepcopy(
        rowan_source_integrationtest_env_transfer_request
    )
    rowan_request.sifchain_destination_address = new_sifaddr
    rowan_request.amount = rowan_amount
    logging.debug(f"transfer {rowan_amount} to {new_sifaddr} from {rowan_request.sifchain_address}")
    burn_lock_functions.transfer_sifchain_to_sifchain(rowan_request, credentials)


def test_first_fanout_account(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        rowan_amount,
        ceth_amount
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.check_wait_blocks = True
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_rowan_balance=rowan_amount,
        target_ceth_balance=ceth_amount
    )
    test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, "eth")


def test_transfer_eth_to_ceth_over_limit(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
):
    basic_transfer_request.ethereum_symbol = "eth"
    basic_transfer_request.ethereum_address = source_ethereum_address
    basic_transfer_request.amount = 35 * 10 ** 18
    with pytest.raises(Exception):
        basic_transfer_request.ethereum_address = source_ethereum_address
        generate_test_account(
            basic_transfer_request,
            rowan_source_integrationtest_env_transfer_request,
            rowan_source_integrationtest_env_credentials,
            target_ceth_balance=50 * 10 ** 19,
        )
