import logging
import os
import time

import burn_lock_functions
import test_utilities
from burn_lock_functions import EthereumToSifchainTransferRequest
from integration_env_credentials import sifchain_cli_credentials_for_test
from test_utilities import get_required_env_var, get_shell_output, SifchaincliCredentials


def build_request(
        smart_contracts_dir,
        ethereum_address,
        solidity_json_path,
) -> (EthereumToSifchainTransferRequest, SifchaincliCredentials):
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
        solidity_json_path=solidity_json_path
    )
    return request, credentials


def test_transfer_eth_to_ceth_using_replay_blocks(
        integration_dir,
        smart_contracts_dir,
        solidity_json_path,
        source_ethereum_address,
        validator_address,
        ensure_relayer_restart
):
    starting_block = test_utilities.current_ethereum_block_number(smart_contracts_dir)
    logging.info("stopping ebrelayer")
    test_utilities.get_shell_output("pkill -9 ebrelayer || true")
    request, credentials = build_request(smart_contracts_dir, source_ethereum_address, solidity_json_path)
    logging.info("(no transactions should happen without a relayer)")
    # test_utilities.whitelist_token(request.ethereum_symbol, request.smart_contracts_dir, True)
    logging.info(f"send {request.amount / 10 ** 18} eth ({request.amount} wei) to {request.sifchain_address}")
    test_utilities.send_from_ethereum_to_sifchain(request)
    # test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")

    logging.info("replay blocks using ebrelayer replayEthereum")
    ews = test_utilities.get_required_env_var("ETHEREUM_WEBSOCKET_ADDRESS")
    bra = test_utilities.get_required_env_var("BRIDGE_REGISTRY_ADDRESS")
    mon = test_utilities.get_required_env_var("MONIKER")
    mn = test_utilities.get_required_env_var("MNEMONIC")
    cn = test_utilities.get_required_env_var("CHAINNET")
    ending_block = test_utilities.current_ethereum_block_number(smart_contracts_dir) + 1
    cmd = f"""yes | ebrelayer replayEthereum tcp://0.0.0.0:26657 {ews} {bra} {mon} '{mn}' {starting_block} {ending_block} 1 2 --chain-id {cn} --gas 5000000000000 \
 --keyring-backend test --node tcp://0.0.0.0:26657 --from {mon}"""
    test_utilities.get_shell_output(cmd)
    time.sleep(15)
    logging.info(f"check the ending balance of {request.sifchain_address} after replaying blocks")
    ending_balance = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                              request.sifchain_symbol)
    assert (ending_balance == request.amount)

    # now do it again
    test_utilities.get_shell_output(cmd)
    time.sleep(5)
    ending_balance2 = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                               request.sifchain_symbol)
    assert (ending_balance2 == request.amount)

    # now start ebrelayer and do another transfer
    test_utilities.advance_n_ethereum_blocks(test_utilities.n_wait_blocks + 1, smart_contracts_dir)
    test_utilities.get_shell_output(f"{integration_dir}/sifchain_start_ebrelayer.sh")
    burn_lock_functions.transfer_ethereum_to_sifchain(request, 15)
    ending_balance3 = test_utilities.get_sifchain_addr_balance(request.sifchain_address, request.sifnodecli_node,
                                                               request.sifchain_symbol)
    assert (ending_balance3 == request.amount * 2)
