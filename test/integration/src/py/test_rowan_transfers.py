import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
import test_utilities
from test_utilities import get_required_env_var, SifchaincliCredentials, get_optional_env_var, ganache_owner_account

bridgetoken_address = get_required_env_var("BRIDGE_TOKEN_ADDRESS")
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")
smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
owner_password = get_required_env_var("OWNER_PASSWORD", "because we need to get rowan from the owner")
ethereum_address = get_optional_env_var(
    "ETHEREUM_ADDRESS",
    ganache_owner_account(smart_contracts_dir)
)


# this transfers rowan that's already in the owner account back to the ethereum side,
# so we don't need to mint new rowan
def test_transfer_rowan_to_erowan_and_back():
    # we need to use the credentials that were created for the owner to get rowan
    credentials = SifchaincliCredentials(
        keyring_passphrase=owner_password,
        keyring_backend="file",
        from_key=get_required_env_var("MONIKER"),
        sifnodecli_homedir=f"""{get_required_env_var("CHAINDIR")}/.sifnodecli"""
    )
    request = EthereumToSifchainTransferRequest(
        ethereum_symbol="eth",
        sifchain_symbol="ceth",
        sifchain_address=get_required_env_var("OWNER_ADDR"),
        smart_contracts_dir=smart_contracts_dir,
        ethereum_address=ethereum_address,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=bridgebank_address,
        bridgetoken_address=bridgetoken_address,
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=10 ** 17,
        ceth_amount=2 * 10 ** 16
    )
    logging.info(f"get initial ceth to cover fees: {request}")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 30)

    request.ethereum_symbol = bridgetoken_address
    request.sifchain_symbol = "rowan"
    request.amount = 12000
    request.ethereum_address = test_utilities.ganache_second_account(smart_contracts_dir)

    starting_balance = burn_lock_functions.get_eth_balance(request)

    logging.info(f"transfer rowan to erowan: {request}")
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)

    ending_balance = burn_lock_functions.get_eth_balance(request)
    assert(ending_balance == starting_balance + request.amount)

    test_utilities.whitelist_token(bridgetoken_address, smart_contracts_dir)
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 20)


@pytest.mark.skip(reason="not implemented")
def test_transfer_erowan_to_another_sifchain_address():
    assert False
