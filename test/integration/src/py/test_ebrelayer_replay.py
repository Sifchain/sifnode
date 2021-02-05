import time
import logging
import os
from copy import copy, deepcopy
from functools import lru_cache
from json import JSONDecodeError

import pytest

import burn_lock_functions
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
import test_utilities
from test_utilities import get_required_env_var, get_shell_output, SifchaincliCredentials, amount_in_wei, \
    get_optional_env_var, ganache_owner_account

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")

ethereum_address = get_optional_env_var(
    "ETHEREUM_ADDRESS",
    ganache_owner_account(smart_contracts_dir)
)


def build_request() -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=smart_contracts_dir,
        ethereum_address=ethereum_address,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=get_required_env_var("BRIDGE_BANK_ADDRESS"),
        ethereum_network=(os.environ.get("ETHEREUM_NETWORK") or ""),
        amount=9 * 10 ** 18,
        ceth_amount=2 * (10 ** 16)
    )
    return request, credentials


@pytest.mark.skip(reason="ebrelayer replay not implemented yet")
def test_transfer_eth_to_ceth(integration_dir, ensure_relayer_restart):
    logging.info("stopping ebrelayer")
    test_utilities.get_shell_output("pkill -9 ebrelayer || true")
    request, credentials = build_request()
    logging.info("(no transactions should happen without a relayer)")
    # test_utilities.whitelist_token(request.ethereum_symbol, request.smart_contracts_dir, True)
    logging.info(f"send {request.amount / 10**18} eth ({request.amount} wei) to {request.sifchain_address}")
    test_utilities.send_from_ethereum_to_sifchain(request)
    # test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")

    logging.info("replay blocks using ebrelayer replayEthereum")
    ews=test_utilities.get_required_env_var("ETHEREUM_WEBSOCKET_ADDRESS")
    bra=test_utilities.get_required_env_var("BRIDGE_REGISTRY_ADDRESS")
    mon=test_utilities.get_required_env_var("MONIKER")
    mn=test_utilities.get_required_env_var("MNEMONIC")
    cn=test_utilities.get_required_env_var("CHAINNET")
    # we should get the block numbers programatically, but we know there aren't many yet,
    # so 1 50 is fine
    cmd = f"""ebrelayer replayEthereum tcp://0.0.0.0:26657 {ews} {bra} {mon} '{mn}' 1 50 1 50 --chain-id {cn}"""
    test_utilities.get_shell_output(cmd)
    time.sleep(5)
    ending_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address,request.sifnodecli_node,request.sifchain_symbol)
    assert(ending_balance == request.amount)

    # now do it again
    test_utilities.get_shell_output(cmd)
    time.sleep(5)
    ending_balance2 = test_utilities.get_sifchain_addr_balance(request.sifchain_address,request.sifnodecli_node,request.sifchain_symbol)
    assert(ending_balance2 == request.amount)

    # now start ebrelayer and do another transfer
    test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 15)
    ending_balance3 = test_utilities.get_sifchain_addr_balance(request.sifchain_address,request.sifnodecli_node,request.sifchain_symbol)
    assert(ending_balance3 == request.amount * 2)
