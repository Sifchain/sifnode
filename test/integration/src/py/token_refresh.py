import copy
import json
import logging
import os

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


def test_can_mint_token_and_peg_it(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        bridgebank_address,
        solidity_json_path,
        operator_address,
        ethereum_network,
        source_ethereum_address,
):
    logging.info("token_refresh needs to use the operator private key, setting that to ETHEREUM_PRIVATE_KEY")
    os.environ["ETHEREUM_PRIVATE_KEY"] = test_utilities.get_required_env_var("OPERATOR_PRIVATE_KEY")
    sifaddress, credentials = create_new_sifaddr_and_credentials()
    request = copy.deepcopy(basic_transfer_request)
    request.sifchain_address = sifaddress
    request.ethereum_address = source_ethereum_address
    amount_in_tokens = int(test_utilities.get_required_env_var("TOKEN_AMOUNT"))

    tokens = test_utilities.get_whitelisted_tokens(request)
    logging.warning(f"whitelisted tokens: {tokens}")

    for t in tokens:
        destination_symbol = "c" + t["symbol"]
        if t["symbol"] == "erowan":
            destination_symbol = "rowan"
        try:
            logging.info(f"sending {t}")
            request.amount = amount_in_tokens * (10 ** int(t["decimals"]))
            request.ethereum_symbol = t["token"]
            request.sifchain_symbol = destination_symbol
            request.ethereum_address = operator_address
            test_utilities.mint_tokens(request, operator_address)
            test_utilities.send_from_ethereum_to_sifchain(request)
        except Exception as e:
            # try to get as many tokens across the bridge as you can,
            # don't stop if one of them fails
            logging.info(f"failed to mint and send for {t}, error was {e}")
