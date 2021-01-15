import os

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var
from test_utilities import get_shell_output


# to test against ropsten, define:
# ETHEREUM_ADDRESS
# ETHEREUM_PRIVATE_KEY
# ETHEREUM_NETWORK = ropsten
#


def generate_test_account(starting_ceth_balance: int = 10 ** 19):
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
        amount=starting_ceth_balance,
        ceth_amount=2 * (10 ** 16)
    )
    burn_lock_functions.transfer_ethereum_to_sifchain(request)
    return {
        "user": new_account_key,
        "request": request,
        "credentials": credentials
    }
