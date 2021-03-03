import copy
import json
import logging

import burn_lock_functions
import test_utilities
from integration_env_credentials import sifchain_cli_credentials_for_test
from pytest_utilities import generate_minimal_test_account
from test_utilities import EthereumToSifchainTransferRequest


def create_new_sifaddr():
    new_account_key = test_utilities.get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    return new_addr["address"]


def test_bulk_transfers(
        basic_transfer_request: EthereumToSifchainTransferRequest,
        smart_contracts_dir,
        source_ethereum_address,
        bridgebank_address,
        integration_dir,
        ganache_timed_blocks,
):
    n_transfers = int(test_utilities.get_optional_env_var("NTRANSFERS", 1))
    ganache_delay = test_utilities.get_optional_env_var("GANACHE_DELAY", 1)
    # test_utilities.get_shell_output(f"{integration_dir}/ganache_start.sh {ganache_delay}")
    amount = 9000
    new_addresses = list(map(lambda x: create_new_sifaddr(), range(n_transfers)))
    logging.debug(f"new_addresses: {new_addresses}")
    request: EthereumToSifchainTransferRequest = copy.deepcopy(basic_transfer_request)
    requests = list(map(lambda addr: {
        "amount": amount,
        "symbol": test_utilities.NULL_ADDRESS,
        "sifchain_address": addr
    }, new_addresses))
    json_requests = json.dumps(requests)
    yarn_result = test_utilities.run_yarn_command(
        " ".join([
            f"yarn --cwd {smart_contracts_dir}",
            "integrationtest:sendBulkLockTx",
            f"--amount {amount}",
            f"--symbol eth",
            f"--json_path {request.solidity_json_path}",
            f"--sifchain_address {new_addresses[0]}",
            f"--transactions \'{json_requests}\'",
            f"--ethereum_address {source_ethereum_address}",
            f"--bridgebank_address {bridgebank_address}"
        ])
    )
    logging.info(f"bulk result: {yarn_result}")
    test_utilities.wait_for_ethereum_block_number(yarn_result["blockNumber"] + test_utilities.n_wait_blocks, basic_transfer_request);
    for a in new_addresses:
        test_utilities.wait_for_sif_account(a, basic_transfer_request.sifnodecli_node, 90)
        test_utilities.wait_for_sifchain_addr_balance(a, "ceth", amount, basic_transfer_request.sifnodecli_node, 90)
