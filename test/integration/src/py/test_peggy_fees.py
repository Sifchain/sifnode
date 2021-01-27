# to test against ropsten, define:
# ETHEREUM_ADDRESS
# ETHEREUM_PRIVATE_KEY
# ETHEREUM_NETWORK = ropsten
#
import logging
import os

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from test_utilities import get_required_env_var, SifchaincliCredentials


def build_request():
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
        amount=2 * (10 ** 17),
        ceth_amount=2 * 10 ** 16 - 37
    )
    return request, credentials


def test_charge_a_fee_on_a_sifchain_burn():
    logging.info(f"get initial ceth to cover fees")
    (request, credentials) = build_request()
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 10)

    request.ethereum_symbol = "erowan"
    request.sifchain_symbol = "rowan"
    request.amount = 12000
    logging.info(f"transfer rowan to erowan: {request}")
    ceth_balance = burn_lock_functions.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                 "ceth")
    logging.info(f"initial ceth balance is {ceth_balance}")
    burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    new_ceth_balance = burn_lock_functions.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                                     "ceth")
    ceth_fee = ceth_balance - new_ceth_balance
    logging.info(f"ceth fee is {ceth_fee}")
    logging.info(f"final ceth balance change due to fees should be <= the request ceth_amount {request.ceth_amount}")
    assert (ceth_fee <= request.ceth_amount)


def test_do_not_transfer_if_fee_allowed_is_too_low():
    logging.info(f"get initial ceth to cover fees")
    (request, credentials) = build_request()
    request.ethereum_symbol = "erowan"
    request.sifchain_symbol = "rowan"
    request.amount = 12000
    request.ceth_amount = 65000000000 * 248692 - 1  # from x/ethbridge/types/msgs.go

    starting_rowan_balance = burn_lock_functions.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, request.sifchain_symbol)
    with pytest.raises(Exception):
        burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)
    ending_rowan_balance = burn_lock_functions.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node, request.sifchain_symbol)
    logging.info(f"starting_rowan_balance {starting_rowan_balance} and ending_rowan_balance {ending_rowan_balance} should be equal")
    assert(starting_rowan_balance == ending_rowan_balance)
