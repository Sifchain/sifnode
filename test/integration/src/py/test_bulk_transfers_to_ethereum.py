import copy
import logging

import pytest

import burn_lock_functions
import test_utilities
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials


def create_new_sifaddr():
    new_account_key = test_utilities.get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    return new_addr["address"]


def create_new_sifaddr_and_key():
    new_account_key = test_utilities.get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    return new_addr["address"], new_addr["name"]


@pytest.mark.skipif(
    not test_utilities.get_optional_env_var("NTRANSFERS", None),
    reason="run by hand and specify NTRANSFERS"
)
def test_bulk_transfers_from_sifchain(
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
        sifchain_fees_int,
):
    test_transfer_amount = 100  # just a tiny number of wei to move to confirm things are working
    tokens = test_utilities.get_required_env_var("TOKENS", "ceth,rowan").split(",")
    logging.info(f"tokens to be transferred are: {tokens}")
    logging.info("create new ethereum and sifchain addresses")
    basic_transfer_request.ethereum_address = source_ethereum_address
    n_transfers = int(test_utilities.get_optional_env_var("NTRANSFERS", 2))
    n_transactions = n_transfers * len(tokens)
    new_addresses_and_keys = list(map(lambda x: create_new_sifaddr_and_key(), range(n_transactions)))
    logging.debug(f"new_addresses_and_keys: {new_addresses_and_keys}")
    credentials_for_account_with_ceth = SifchaincliCredentials(from_key=rowan_source_key)
    request: EthereumToSifchainTransferRequest = copy.deepcopy(basic_transfer_request)
    ceth_amount = n_transactions * (test_utilities.highest_gas_cost + 100)
    request.amount = ceth_amount
    request.ethereum_address = source_ethereum_address
    request.sifchain_address = rowan_source
    addresses_to_populate = copy.deepcopy(new_addresses_and_keys)
    test_transfers = []
    for a in range(n_transfers):
        for t in tokens:
            request.sifchain_destination_address, from_key = addresses_to_populate.pop()

            # send ceth to pay for the burn
            request.amount = test_utilities.burn_gas_cost
            request.sifchain_symbol = "ceth"
            burn_lock_functions.transfer_sifchain_to_sifchain(request, credentials_for_account_with_ceth)

            # send rowan to pay the fee
            request.amount = sifchain_fees_int
            request.sifchain_symbol = "rowan"
            burn_lock_functions.transfer_sifchain_to_sifchain(request, credentials_for_account_with_ceth)

            # send the token itself
            request.amount = test_transfer_amount
            request.sifchain_symbol = t
            burn_lock_functions.transfer_sifchain_to_sifchain(request, credentials_for_account_with_ceth)
            transfer = (request.sifchain_destination_address, from_key, request.sifchain_symbol, request.amount)

            test_utilities.get_sifchain_addr_balance(request.sifchain_destination_address, request.sifnoded_node, t)

            test_transfers.append(transfer)

    logging.debug(f"test_transfers is {test_transfers}")

    text_file = open("pfile.cmds", "w")
    simple_credentials = SifchaincliCredentials(
        keyring_passphrase=None,
        keyring_backend="test",
        from_key=None,
        sifnoded_homedir=None
    )

    logging.info(f"all accounts are on sifchain and have the correct balance")

    new_eth_addrs = test_utilities.create_ethereum_addresses(
        smart_contracts_dir,
        basic_transfer_request.ethereum_network,
        n_transactions
    )
    logging.debug(f"new eth addrs: {new_eth_addrs}")

    ethereum_transfers = []
    for sifaddr, from_key, sifsymbol, amount in test_transfers:
        destination_ethereum_address_element = new_eth_addrs.pop()
        r = copy.deepcopy(basic_transfer_request)
        r.sifchain_symbol = sifsymbol
        r.sifchain_address = sifaddr
        r.ethereum_address = destination_ethereum_address_element["address"]
        r.amount = amount
        simple_credentials.from_key = from_key
        c = test_utilities.send_from_sifchain_to_ethereum_cmd(r, simple_credentials)
        ethereum_symbol = test_utilities.sifchain_symbol_to_ethereum_symbol(sifsymbol)
        transfer = (r.ethereum_address, ethereum_symbol, amount)
        ethereum_transfers.append(transfer)
        text_file.write(f"{c}\n")
    text_file.close()
    test_utilities.get_shell_output("cat pfile.cmds | parallel --trim lr -v {}")
    whitelist = test_utilities.get_whitelisted_tokens(basic_transfer_request)
    test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks, smart_contracts_dir)
    for ethereum_address, ethereum_symbol, amount in ethereum_transfers:
        r = copy.deepcopy(basic_transfer_request)
        r.ethereum_address = ethereum_address
        r.ethereum_symbol = test_utilities.get_token_ethereum_address(
            ethereum_symbol,
            whitelist
        )
        r.amount = amount
        test_utilities.wait_for_eth_balance(
            transfer_request=r,
            target_balance=amount,
            max_seconds=60 * 60 * 10
        )
