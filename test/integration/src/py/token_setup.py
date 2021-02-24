import copy
import json
import logging

import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import create_new_sifaddr_and_credentials
from test_utilities import create_new_currency


# This file is for setting up an installation with a set of currencies given
# in a json file like ui/core/src/assets.sifchain.mainnet.json.
#
# Run example:
#
#  TOKENS_FILE=.../sifnode/ui/core/src/assets.sifchain.mainnet.json python3 -m pytest src/py/token_setup.py

def build_request_for_new_sifchain_address(basic_transfer_request, source_ethereum_address, new_currency, amount):
    sifaddress, credentials = create_new_sifaddr_and_credentials()
    request = copy.deepcopy(basic_transfer_request)
    request.ethereum_symbol = new_currency["newtoken_address"]
    request.ethereum_address = source_ethereum_address
    request.sifchain_symbol = "c" + new_currency["newtoken_symbol"]
    request.sifchain_address = sifaddress
    request.amount = amount
    return request, credentials


def test_can_create_a_new_token_and_peg_it(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        operator_address,
        ethereum_network,
        source_ethereum_address,
):
    json_filename = test_utilities.get_required_env_var("TOKENS_FILE")

    with open(json_filename, mode="r") as json_file:
        contents = json_file.read()
        tokens = json.loads(contents)

    sifaddress, credentials = create_new_sifaddr_and_credentials()
    request = copy.deepcopy(basic_transfer_request)
    request.sifchain_address = sifaddress
    request.ethereum_address = source_ethereum_address
    amount = 10 ** 9 * 10 ** 18  # one billion of the token
    request.amount = amount

    for t in tokens["assets"]:
        if t["symbol"] == "rowan" or t["symbol"] == 'ceth':
            continue
        logging.info(f"token is: {t}")
        try:
            new_currency = create_new_currency(
                amount=amount,
                symbol=t["symbol"][1:],
                token_name=t["name"],
                decimals=int(t["decimals"]),
                smart_contracts_dir=smart_contracts_dir,
                bridgebank_address=bridgebank_address,
                solidity_json_path=solidity_json_path,
                operator_address=operator_address,
                ethereum_network=ethereum_network
            )
            request.ethereum_symbol = new_currency["newtoken_address"]
            request.sifchain_symbol = "c" + new_currency["newtoken_symbol"]
            burn_lock_functions.transfer_ethereum_to_sifchain(request)
        except Exception as e:
            # it might already exist, so do nothing
            logging.info(f"failed to create token {t}, error was {e}")
