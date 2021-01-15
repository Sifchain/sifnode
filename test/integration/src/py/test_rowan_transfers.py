import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from test_utilities import get_required_env_var, SifchaincliCredentials


# to test against ropsten, define:
# ETHEREUM_ADDRESS
# ETHEREUM_PRIVATE_KEY
# ETHEREUM_NETWORK = ropsten
#


def test_transfer_rowan_to_erowan():
    credentials = SifchaincliCredentials(
        keyring_passphrase=get_required_env_var("OWNER_PASSWORD"),
        keyring_backend="file",
        from_key=get_required_env_var("MONIKER"),
        sifnodecli_homedir=f"""{get_required_env_var("CHAINDIR")}/.sifnodecli"""
    )
    request = EthereumToSifchainTransferRequest(
        ethereum_symbol="eth",
        sifchain_symbol="ceth",
        sifchain_address=get_required_env_var("OWNER_ADDR"),
        smart_contracts_dir=get_required_env_var("SMART_CONTRACTS_DIR"),
        ethereum_address=get_required_env_var("ETHEREUM_ADDRESS"),
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        bridgetoken_address=get_required_env_var("BRIDGE_TOKEN_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=10 ** 17,
        ceth_amount=2 * 10 ** 16
    )
    logging.info(f"get initial ceth to cover fees: {request}")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 10)

    request.ethereum_symbol = "erowan"
    request.sifchain_symbol = "rowan"
    request.amount = 12000
    logging.info(f"transfer rowan to erowan: {request}")
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)


@pytest.mark.skip(reason="not implemented")
def test_transfer_erowan_to_rowan():
    assert False


@pytest.mark.skip(reason="not implemented")
def test_transfer_erowan_to_another_sifchain_address():
    assert False
