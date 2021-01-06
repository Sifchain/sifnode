import logging

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest, SifchaincliCredentials
from integration_env_credentials import credentials_for_test
from test_utilities import get_required_env_var, get_shell_output


def test_transfer_eth_to_new_sifchain_account_and_back():
    new_account_key = get_shell_output("uuidgen")
    credentials = credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)

    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=get_required_env_var("SMART_CONTRACTS_DIR"),
        ethereum_address=get_required_env_var("ETHEREUM_ADDRESS"),
        amount=90000
    )

    result = burn_lock_functions.transfer_ethereum_to_sifchain(request, credentials)
    logging.info(f"transfer_ethereum_to_sifchain_result_json: {result}")
    return_result = burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    logging.info(f"transfer_sifchain_to_ethereum__result_json: {return_result}")


if __name__ == "__main__":
    test_transfer_eth_to_new_sifchain_account_and_back()