import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account

# bridgetoken_address = get_required_env_var("BRIDGE_TOKEN_ADDRESS")
# bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
# smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
# owner_password = get_required_env_var("OWNER_PASSWORD", "because we need to get rowan from the owner")
# ethereum_address = get_optional_env_var(
#     "ETHEREUM_ADDRESS",
#     ganache_owner_account(smart_contracts_dir)
# )


def test_rowan_to_erowan(
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

    logging.info(f"send erowan to ethereum from test account")
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = int(target_rowan_balance / 2)
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
