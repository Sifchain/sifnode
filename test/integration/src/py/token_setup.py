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
    amount_in_tokens = 10 ** 9  # one billion of the token; note that this is not 1/(10 **18) of a token

    existing_whitelist = test_utilities.get_whitelisted_tokens(request)
    logging.info(f"existing whitelist: {existing_whitelist}")
    existing_tokens = set(map(lambda w: "c" + w["symbol"], existing_whitelist))
    logging.info(f"requested tokens: {tokens}")
    for t in tokens["assets"]:
        if t["symbol"] in existing_tokens or t["symbol"] == "rowan":
            logging.info(f"token {t} already whitelisted, skipping")
            continue
        logging.info(f"whitelisting token {t}")
        decimals_int = int(t["decimals"])
        amount_in_fractions_of_a_token = amount_in_tokens * (10 ** decimals_int)
        create_new_currency(
            amount=amount_in_fractions_of_a_token,
            symbol=t["symbol"][1:],
            token_name=t["name"],
            decimals=decimals_int,
            smart_contracts_dir=smart_contracts_dir,
            bridgebank_address=bridgebank_address,
            solidity_json_path=solidity_json_path,
            operator_address=operator_address,
            ethereum_network=ethereum_network
        )
