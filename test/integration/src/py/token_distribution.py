import copy
import logging

import pytest
import time

import burn_lock_functions
import test_utilities
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


@pytest.mark.skipif(not test_utilities.get_optional_env_var("DESTINATION_ACCOUNT", None), reason="run by hand and specify DESTINATION_ACCOUNT")
def test_token_distribution(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        rowan_source_integrationtest_env_credentials: SifchaincliCredentials,
        rowan_source_integrationtest_env_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        source_ethereum_address,
        rowan_source,
        rowan_source_key,
        bridgebank_address,
        bridgetoken_address,
        ethereum_network,
):
    tokens = test_utilities.get_whitelisted_tokens(basic_transfer_request)
    request = basic_transfer_request
    amount_in_tokens = 10000000

    for t in tokens:
        try:
            destination_symbol = "c" + t["symbol"]
            if t["symbol"] == "erowan":
                destination_symbol = "rowan"
            logging.info(f"sending {t}")
            request.amount = int(amount_in_tokens * (10 ** int(t["decimals"])))
            request.ethereum_symbol = t["token"]
            request.sifchain_symbol = destination_symbol
            request.sifchain_address = rowan_source
            request.sifchain_destination_address = test_utilities.get_required_env_var("DESTINATION_ACCOUNT")
            test_utilities.send_from_sifchain_to_sifchain(request, rowan_source_integrationtest_env_credentials)
        except Exception as e:
            logging.error(f"error: {e}")
