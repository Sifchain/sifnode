import concurrent
import logging
import os
from concurrent.futures.thread import ThreadPoolExecutor

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output


def test_transfer_eth_to_ceth_in_parallel():
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = {executor.submit(execute_one_transfer, x) for x in range(0, 3)}
        for f in concurrent.futures.as_completed(futures):
            # As a side effect, this will raise any exception that happened in the future
            logging.info(f"The result is {f.result()}")


def execute_one_transfer(id_number: int):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=get_required_env_var("SMART_CONTRACTS_DIR"),
        ethereum_address=get_required_env_var("ETHEREUM_ADDRESS"),
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=2000,
        ceth_amount=2 * (10 ** 16)
    )
    logging.info(f"execute request #{id_number}")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 10)
    return f"transfered eth to ceth: {request}"
