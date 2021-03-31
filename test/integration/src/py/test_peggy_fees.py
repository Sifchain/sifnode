import logging

import pytest

import burn_lock_functions
import integration_env_credentials
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from pytest_utilities import generate_test_account
from test_utilities import SifchaincliCredentials


def test_rescue_ceth(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        sifchain_fees_int,
        ethbridge_module_address,
):
    """
    does a lock of rowan (using another test) that should result
    in ceth being sent to a place it can be rescued from
    """
    admin_account = test_utilities.get_required_env_var("SIFCHAIN_ADMIN_ACCOUNT")
    basic_transfer_request.ethereum_address = source_ethereum_address
    admin_user_credentials = SifchaincliCredentials(
        from_key="sifnodeadmin"
    )
    small_amount = 100
    test_account_request, test_account_credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request=rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials=rowan_source_integrationtest_env_credentials,
        target_ceth_balance=test_utilities.burn_gas_cost + small_amount,
        target_rowan_balance=sifchain_fees_int
    )
    test_account_request.amount = small_amount
    burn_lock_functions.transfer_sifchain_to_ethereum(test_account_request, test_account_credentials)
    logging.info(
        f"test account {test_account_request.sifchain_address} should now have no ceth")
    logging.info("ethbridge should have the fee that was paid")
    test_utilities.wait_for_sifchain_addr_balance(
        ethbridge_module_address,
        "ceth",
        test_utilities.burn_gas_cost,
        test_account_request.sifnodecli_node
    )
    logging.info(f"now rescue ceth into {test_account_request.sifchain_address}")
    test_utilities.rescue_ceth(
        receiver_account=test_account_request.sifchain_address,
        admin_account=admin_account,
        amount=test_utilities.burn_gas_cost,
        transfer_request=basic_transfer_request,
        credentials=admin_user_credentials
    )
    test_utilities.wait_for_sifchain_addr_balance(
        test_account_request.sifchain_address,
        "ceth",
        test_utilities.burn_gas_cost,
        test_account_request.sifnodecli_node,
        max_seconds=10,
        debug_prefix="wait for rescue ceth"
    )


def test_ceth_receiver_account(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgetoken_address,
        validator_address,
):
    admin_account = test_utilities.get_required_env_var("SIFCHAIN_ADMIN_ACCOUNT")
    ceth_rescue_account, ceth_rescue_account_credentials = integration_env_credentials.create_new_sifaddr_and_credentials()
    basic_transfer_request.sifchain_address = validator_address
    admin_user_credentials = SifchaincliCredentials(
        from_key="sifnodeadmin"
    )
    test_utilities.update_ceth_receiver_account(
        receiver_account=ceth_rescue_account,
        admin_account=admin_account,
        transfer_request=basic_transfer_request,
        credentials=admin_user_credentials
    )
    test_fee_charged_to_transfer_rowan_to_erowan(
        basic_transfer_request=basic_transfer_request,
        source_ethereum_address=source_ethereum_address,
        rowan_source_integrationtest_env_credentials=rowan_source_integrationtest_env_credentials,
        rowan_source_integrationtest_env_transfer_request=rowan_source_integrationtest_env_transfer_request,
        ethereum_network=ethereum_network,
        smart_contracts_dir=smart_contracts_dir,
        bridgetoken_address=bridgetoken_address,
    )
    received_ceth_charges = test_utilities.get_sifchain_addr_balance(ceth_rescue_account,
                                                                     basic_transfer_request.sifnodecli_node, "ceth")
    assert received_ceth_charges == test_utilities.burn_gas_cost


def test_fee_charged_to_transfer_rowan_to_erowan(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgetoken_address,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    logging.info(f"credentials: {rowan_source_integrationtest_env_credentials}")
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=10 ** 18,
        target_rowan_balance=10 ** 18
    )
    # send some test account ceth back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    logging.info(f"sending rowan to erowan and checking that a ceth fee was charged")
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = 31500

    # get the starting ceth balance, transfer some rowan to erowan, get the ending ceth
    # balance.  The difference is the fee charged and should be equal to request.ceth_amount

    starting_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                     "ceth")
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                   "ceth")
    fee = starting_ceth_balance - ending_ceth_balance
    assert fee == request.ceth_amount


def test_do_not_transfer_if_fee_allowed_is_too_low(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        source_ethereum_address: str,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        ethereum_network,
        smart_contracts_dir,
        bridgetoken_address,
):
    basic_transfer_request.ethereum_address = source_ethereum_address
    target_ceth_balance = 10 ** 18
    target_rowan_balance = 10 ** 18
    request, credentials = generate_test_account(
        basic_transfer_request,
        rowan_source_integrationtest_env_transfer_request,
        rowan_source_integrationtest_env_credentials,
        target_ceth_balance=target_ceth_balance,
        target_rowan_balance=target_rowan_balance
    )
    # send some test account ceth back to a new ethereum address
    request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    request.sifchain_symbol = "rowan"
    request.ethereum_symbol = bridgetoken_address
    request.amount = 31500

    logging.info("try to transfer rowan to erowan with a ceth_amount that's too low")
    with pytest.raises(Exception):
        request.ceth_amount = test_utilities.lock_gas_cost - 1
        burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                   "ceth")
    assert ending_ceth_balance == target_ceth_balance

    logging.info("try with not owning enough ceth to cover the offer")
    with pytest.raises(Exception):
        request.ceth_amount = target_ceth_balance + 1
        burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_ceth_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                   "ceth")
    assert ending_ceth_balance == target_ceth_balance
