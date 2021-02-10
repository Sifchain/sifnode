import logging
from copy import deepcopy

import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output, SifchaincliCredentials, get_optional_env_var

smart_contracts_dir = get_required_env_var("SMART_CONTRACTS_DIR")
ethereum_network = get_optional_env_var("ETHEREUM_NETWORK", "")
amount = int(get_optional_env_var("AMOUNT", "20000"))
bridgebank_address = get_required_env_var("BRIDGE_BANK_ADDRESS")


def build_request(sifnodecli_node, source_ethereum_address, chain_id) -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    ceth_fee = 2 * (10 ** 16)
    request = EthereumToSifchainTransferRequest(
        sifchain_address=new_addr["address"],
        smart_contracts_dir=smart_contracts_dir,
        ethereum_address=source_ethereum_address,
        ethereum_private_key_env_var="ETHEREUM_PRIVATE_KEY",
        bridgebank_address=bridgebank_address,
        ethereum_network=ethereum_network,
        amount=amount + ceth_fee,
        ceth_amount=ceth_fee,
        sifnodecli_node=sifnodecli_node,
        manual_block_advance=False,
        chain_id=chain_id,
        sifchain_fees="100000rowan",
    )
    return request, credentials


def test_transfer_eth_to_ceth_and_back(ropsten_wait_time, rowan_source, sifnodecli_node, source_ethereum_address, chain_id):
    eth_transfer, credentials = build_request(sifnodecli_node, source_ethereum_address, chain_id)
    # first we have to give the sif address some rowan so it can pay sifchain gas
    rowan_transfer = deepcopy(eth_transfer)
    logging.info(f"rowan_source is {rowan_source}")
    rowan_transfer.sifchain_address = rowan_source
    rowan_transfer.sifchain_destination_address = eth_transfer.sifchain_address
    rowan_transfer.amount = 2000000
    rowan_transfer.sifchain_symbol = "rowan"
    test_utilities.send_from_sifchain_to_sifchain(rowan_transfer, credentials)
    burn_lock_functions.transfer_ethereum_to_sifchain(eth_transfer, ropsten_wait_time)
    logging.info(f"send ceth back to {eth_transfer.ethereum_address}")
    return_request = deepcopy(eth_transfer)
    # don't transfer ceth => eth to the BridgeBank address since BridgeBank is responsible for paying gas.
    # That means you can't just see if the exact transfer went through.
    return_request.ethereum_address, _ = test_utilities.create_ethereum_address(
        smart_contracts_dir, ethereum_network
    )
    return_request.amount = amount
    burn_lock_functions.transfer_sifchain_to_ethereum(return_request, credentials, ropsten_wait_time)
