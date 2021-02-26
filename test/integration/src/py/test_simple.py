import copy

import test_utilities
from pytest_utilities import generate_test_account
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def all_balances(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        sifaddr,
        ethaddr,
        tokens
):
    request = copy.deepcopy(basic_transfer_request)
    request.ethereum_address = ethaddr
    for t in tokens:
        request.ethereum_symbol = t
        test_utilities.get_eth_balance(request)
    test_utilities.get_sifchain_addr_balance(sifaddr, request.sifnodecli_node, "ceth")


def test_create_account(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
):
    basic_transfer_request.ethereum_symbol = "eth"
    basic_transfer_request.ethereum_address = source_ethereum_address
    amount = 5 * test_utilities.highest_gas_cost
    basic_transfer_request.amount = 5 * test_utilities.highest_gas_cost
    basic_transfer_request.ethereum_address = source_ethereum_address
    generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=amount,
        target_rowan_balance=amount
    )
