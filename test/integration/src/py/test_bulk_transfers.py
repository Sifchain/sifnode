import copy
import json
import logging

import pytest

import burn_lock_functions
import test_utilities
from integration_env_credentials import sifchain_cli_credentials_for_test
from pytest_utilities import generate_minimal_test_account
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


@pytest.mark.skip(reason="run manually")
def test_bulk_transfers(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        source_ethereum_address,
        bridgebank_address,
        bridgetoken_address,
        ethereum_network,
):
    n_transfers = int(test_utilities.get_optional_env_var("NTRANSFERS", 2))
    ganache_delay = test_utilities.get_optional_env_var("GANACHE_DELAY", 1)
    # test_utilities.get_shell_output(f"{integration_dir}/ganache_start.sh {ganache_delay}")
    amount = "{:d}".format(5 * test_utilities.highest_gas_cost)
    new_addresses_and_keys = list(map(lambda x: create_new_sifaddr_and_key(), range(n_transfers)))
    logging.info(f"aandk: {new_addresses_and_keys}")
    new_addresses = list(map(lambda a: a[0], new_addresses_and_keys))
    logging.debug(f"new_addresses: {new_addresses}")
    new_eth_addrs = test_utilities.create_ethereum_addresses(smart_contracts_dir, basic_transfer_request.ethereum_network, len(new_addresses))
    logging.info(f"new eth addrs: {new_eth_addrs}")
    request: EthereumToSifchainTransferRequest = copy.deepcopy(basic_transfer_request)
    requests = list(map(lambda addr: {
        "amount": amount,
        "symbol": test_utilities.NULL_ADDRESS,
        "sifchain_address": addr
    }, new_addresses))
    json_requests = json.dumps(requests)
    test_utilities.run_yarn_command(
        " ".join([
            f"yarn --cwd {smart_contracts_dir}",
            "integrationtest:sendBulkLockTx",
            f"--amount {amount}",
            f"--symbol eth",
            f"--json_path {request.solidity_json_path}",
            f"--sifchain_address {new_addresses[0]}",
            f"--transactions \'{json_requests}\'",
            f"--ethereum_address {source_ethereum_address}",
            f"--bridgebank_address {bridgebank_address}",
            f"--ethereum_network {ethereum_network}",
        ])
    )
    requests = list(map(lambda addr: {
        "amount": amount,
        "symbol": bridgetoken_address,
        "sifchain_address": addr
    }, new_addresses))
    json_requests = json.dumps(requests)
    yarn_result = test_utilities.run_yarn_command(
        " ".join([
            f"yarn --cwd {smart_contracts_dir}",
            "integrationtest:sendBulkLockTx",
            f"--amount {amount}",
            "--lock_or_burn burn",
            f"--symbol {bridgetoken_address}",
            f"--json_path {request.solidity_json_path}",
            f"--sifchain_address {new_addresses[0]}",
            f"--transactions \'{json_requests}\'",
            f"--ethereum_address {source_ethereum_address}",
            f"--bridgebank_address {bridgebank_address}",
            f"--ethereum_network {ethereum_network}",
        ])
    )
    logging.info(f"bulk result: {yarn_result}")
    manual_advance = False
    if manual_advance:
        test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks, smart_contracts_dir)
    test_utilities.wait_for_ethereum_block_number(yarn_result["blockNumber"] + test_utilities.n_wait_blocks, basic_transfer_request);
    for a in new_addresses:
        test_utilities.wait_for_sif_account(a, basic_transfer_request.sifnodecli_node, 90)
        test_utilities.wait_for_sifchain_addr_balance(a, "ceth", amount, basic_transfer_request.sifnodecli_node, 180)
        test_utilities.wait_for_sifchain_addr_balance(a, "rowan", amount, basic_transfer_request.sifnodecli_node, 180)
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
    # test_utilities.get_shell_output("cat pfile.cmds | parallel --trim lr -v {}")
    test_utilities.get_shell_output("bash -x pfile.cmds")
    for sifaddr, ethaddr in zip(new_addresses_and_keys, new_eth_addrs):
        r = copy.deepcopy(basic_transfer_request)
        r.ethereum_address = ethaddr["address"]
        r.amount = 100
        test_utilities.wait_for_eth_balance(r, 100, 300)
