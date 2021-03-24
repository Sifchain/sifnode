import copy
import logging
import time

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
):
    tokens = test_utilities.get_required_env_var("TOKENS", "ceth,rowan").split(",")
    logging.info("create new ethereum and sifchain addresses")
    basic_transfer_request.ethereum_address = source_ethereum_address
    n_transfers = int(test_utilities.get_optional_env_var("NTRANSFERS", 2))
    amount = "{:d}".format(5 * test_utilities.highest_gas_cost)
    new_addresses_and_keys = list(map(lambda x: create_new_sifaddr_and_key(), range(n_transfers)))
    logging.debug(f"new_addresses_and_keys: {new_addresses_and_keys}")
    new_addresses = list(map(lambda a: a[0], new_addresses_and_keys))
    new_eth_addrs = test_utilities.create_ethereum_addresses(
        smart_contracts_dir,
        basic_transfer_request.ethereum_network,
        len(new_addresses)
    )
    logging.debug(f"new eth addrs: {new_eth_addrs}")
    credentials_for_account_with_ceth = SifchaincliCredentials(from_key=rowan_source_key)
    request: EthereumToSifchainTransferRequest = copy.deepcopy(basic_transfer_request)
    request.amount = (len(tokens) + 1) * test_utilities.highest_gas_cost
    request.ethereum_address = source_ethereum_address
    request.sifchain_address = rowan_source
    for a in new_addresses:
        for t in tokens:
            request.sifchain_destination_address = a
            request.sifchain_symbol = t
            logging.info(f"send {t} to {a}, request is {request}")
            test_utilities.send_from_sifchain_to_sifchain(request, credentials_for_account_with_ceth)
            time.sleep(3)

    for a in new_addresses:
        for t in tokens:
            test_utilities.wait_for_sif_account(a, basic_transfer_request.sifnodecli_node, 90)
            test_utilities.wait_for_sifchain_addr_balance(a, t, amount, basic_transfer_request.sifnodecli_node, 180)

    text_file = open("pfile.cmds", "w")
    simple_credentials = SifchaincliCredentials(
        keyring_passphrase=None,
        keyring_backend="test",
        from_key=None,
        sifnodecli_homedir=None
    )
    logging.info(f"all accounts are on sifchain and have the correct balance")
    for sifaddr, ethaddr in zip(new_addresses_and_keys, new_eth_addrs):
        r = copy.deepcopy(basic_transfer_request)
        r.sifchain_address = sifaddr[0]
        r.ethereum_address = ethaddr["address"]
        r.amount = 100
        simple_credentials.from_key = sifaddr[1]
        c = test_utilities.send_from_sifchain_to_ethereum_cmd(r, simple_credentials)
        text_file.write(f"{c}\n")
    text_file.close()
    test_utilities.get_shell_output("cat pfile.cmds | parallel --trim lr -v {}")
    test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks, smart_contracts_dir)
    # test_utilities.get_shell_output("bash -x pfile.cmds")
    for sifaddr, ethaddr in zip(new_addresses_and_keys, new_eth_addrs):
        r = copy.deepcopy(basic_transfer_request)
        r.ethereum_address = ethaddr["address"]
        r.amount = 100
        test_utilities.wait_for_eth_balance(r, 100, 6000 * (n_transfers + 1))
