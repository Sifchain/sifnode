import copy
import logging

import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_shell_output, SifchaincliCredentials


def generate_minimal_test_account(
        base_transfer_request: EthereumToSifchainTransferRequest,
        target_ceth_balance: int = 10 ** 18,
        timeout=burn_lock_functions.default_timeout_for_ganache
) -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
    """Creates a test account with ceth and rowan.  The address for the new account is in request.sifchain_address"""
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    new_sifaddr = new_addr["address"]
    credentials.from_key = new_addr["name"]

    request: EthereumToSifchainTransferRequest = copy.deepcopy(base_transfer_request)
    request.sifchain_address = new_sifaddr
    request.amount = target_ceth_balance
    request.sifchain_symbol = "ceth"
    request.ethereum_symbol = "eth"
    logging.debug(f"transfer {target_ceth_balance} eth to {new_sifaddr} from {base_transfer_request.ethereum_address}")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, timeout)
    logging.info(
        f"created sifchain addr {new_sifaddr} with {test_utilities.display_currency_value(target_ceth_balance)} ceth")
    return request, credentials


def generate_sifchain_account(
        base_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        target_ceth_balance: int = 10 ** 18,
        target_rowan_balance: int = 10 ** 18
):
    assert False


def generate_test_account(
        base_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        target_ceth_balance: int = 10 ** 18,
        target_rowan_balance: int = 10 ** 18
) -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
    """Creates a test account with ceth and rowan"""
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    new_sifaddr = new_addr["address"]
    credentials.from_key = new_addr["name"]

    if target_rowan_balance > 0:
        rowan_request: EthereumToSifchainTransferRequest = copy.deepcopy(
            rowan_source_integrationtest_env_transfer_request)
        rowan_request.sifchain_destination_address = new_sifaddr
        rowan_request.amount = target_rowan_balance
        logging.debug(f"transfer {target_rowan_balance} rowan to {new_sifaddr} from {rowan_request.sifchain_address}")
        test_utilities.send_from_sifchain_to_sifchain(rowan_request, rowan_source_integrationtest_env_credentials)

    request: EthereumToSifchainTransferRequest = copy.deepcopy(base_transfer_request)
    request.sifchain_address = new_sifaddr
    request.amount = target_ceth_balance
    request.sifchain_symbol = "ceth"
    request.ethereum_symbol = "eth"
    if target_ceth_balance > 0:
        logging.debug(f"transfer {target_ceth_balance} eth to {new_sifaddr} from {base_transfer_request.ethereum_address}")
        burn_lock_functions.transfer_ethereum_to_sifchain(request)

    logging.info(
        f"created sifchain addr {new_sifaddr} "
        f"with {test_utilities.display_currency_value(target_ceth_balance)} ceth "
        f"and {test_utilities.display_currency_value(target_rowan_balance)} rowan"
    )

    return request, credentials


def create_new_sifaddr() -> str:
    new_account_key = test_utilities.get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    return new_addr["address"]
