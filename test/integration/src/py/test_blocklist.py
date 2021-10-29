import logging
import copy
import time

import pytest

import burn_lock_functions
import test_utilities
from pytest_utilities import generate_test_account, get_shell_output, sifchain_cli_credentials_for_test
from test_utilities import EthereumToSifchainTransferRequest, SifchaincliCredentials

import json
import web3
from web3.exceptions import ContractLogicError
import logging
from integration_test_context import main
from integration_test_context import common
from common import *

cmd = main.Integrator()


def web3_connect_ws(host, port):
    return web3.Web3(web3.Web3.WebsocketProvider("ws://{}:{}".format(host, port)))

def get_web3_connection_for_test():
    ethereum_ws_url = test_utilities.get_required_env_var("ETHEREUM_WEBSOCKET_ADDRESS")
    return web3.Web3(web3.WebsocketProvider(ethereum_ws_url))

def get_compiled_sc_ganache(sc_name):
    path = main.project_dir("smart-contracts/build/contracts/{}.json".format(sc_name))
    return json.loads(cmd.read_text_file(path))

def get_sc_abi_ganache(sc_name):
    network_id = 5777
    tmp = get_compiled_sc_ganache(sc_name)
    return tmp["networks"][str(network_id)]["address"], tmp["abi"]

def get_blocklist_sc(w3):
    address, abi = get_sc_abi_ganache("Blocklist")
    result = w3.eth.contract(address=address, abi=abi)
    return result

def get_bridge_bank_sc(w3):
    address, abi = get_sc_abi_ganache("BridgeBank")
    # assert address == test_utilities.get_required_env_var("BRIDGE_BANK_ADDRESS")
    result = w3.eth.contract(address=address, abi=abi)
    return result

def set_blocklist_to(w3, blocklist_sc, addrs):
    addrs = [w3.toChecksumAddress(addr) for addr in addrs]
    current = blocklist_sc.functions.getFullList().call()
    to_add = [addr for addr in addrs if addr not in current]
    to_remove = [addr for addr in current if addr not in addrs]
    txhash1 = blocklist_sc.functions.batchAddToBlocklist(to_add).transact()
    txrcpt1 = w3.eth.wait_for_transaction_receipt(txhash1)
    txhash2 = blocklist_sc.functions.batchRemoveFromBlocklist(to_remove).transact()
    txrcpt2 = w3.eth.wait_for_transaction_receipt(txhash2)
    current = blocklist_sc.functions.getFullList().call()
    assert set(addrs) == set(current)

def create_sifchain_addr():
    mnemonic = random_string(20)
    acct = cmd.sifnoded_keys_add_1(mnemonic)
    return acct["address"]

def send_ether(w3, from_account, to_account, amount):
    logging.info(f"Send {amount} from {from_account} to {to_account}...")
    txhash = w3.eth.send_transaction({
        "from": from_account,
        "to": to_account,
        "value": amount,
        "gas": 30000,
    })
    return w3.eth.wait_for_transaction_receipt(txhash)

def lock_eth(w3, bridge_bank_sc, from_eth_acct, to_sif_acct, amount):
    # Ethereum deposit to
    recipient = to_sif_acct.encode("UTF-8")
    coin_denom = NULL_ADDRESS  # For "eth", otherwise use coin's address
    q = bridge_bank_sc.functions.lock(recipient, coin_denom, amount)
    txhash = q.transact({"from": from_eth_acct, "value": amount})
    txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
    return txrcpt

def create_eth_account(w3, password=""):
    # This creates local account, but does not register it (w3.eth.accounts shows the same number)
    # account = w3.eth.account.create()
    # This creates account in the external node that we're connected to. The node has to support geth extensions.
    # These accounts shouw up in w3.eth.accounts and can be used wih transact().
    # duration must be specified because the method expects 3 parameters.
    account = w3.geth.personal.new_account(password)
    w3.geth.personal.unlock_account(account, password, 0)
    return account

def advance_block_truffle(number):
    args = ["npx", "truffle", "exec", "scripts/advanceBlock.js", str(number)]
    res = cmd.execst(args, cwd=main.project_dir("smart-contracts"))
    return res

def advance_block_w3(w3, number):
    for _ in range(number):
        w3.provider.make_request("evm_mine", [])

def get_sifchain_balance(sif_addr):
    args = ["sifnoded", "query", "bank", "balances", sif_addr, "--limit", str(100000000), "--output", "json"]
    stdout = cmd.execst(args)[1]
    res = json.loads(stdout)["balances"]
    return dict(((x["denom"], int(x["amount"])) for x in res))

def sif_balances_equal(dict1, dict2):
    d2k = set(dict2.keys())
    for k in dict1.keys():
        if (k not in dict2) or (dict1[k] != dict2[k]):
            return False
        d2k.remove(k)
    return len(d2k) == 0

def deploy_bridge_token_for_testing(w3, token_symbol, owner_address, mint_amount):
    # Get BridgeToken contract; on peggy1 branch it is already deployed by truffle migrate.
    sc_json = get_compiled_sc_ganache("BridgeToken")
    abi = sc_json["abi"]
    bytecode = sc_json["bytecode"]
    bridge_token = w3.eth.contract(abi=abi, bytecode=bytecode)
    txhash = bridge_token.constructor(token_symbol).transact()
    txrcpt = w3.eth.get_transaction_receipt(txhash)
    address = txrcpt.contractAddress

    bridge_token = w3.eth.contract(address=address, abi=abi)
    txhash = bridge_token.functions.mint(owner_address, mint_amount).transact()
    txrcpt = w3.eth.get_transaction_receipt(txhash)

    assert bridge_token.functions.balanceOf(owner_address).call() == mint_amount
    assert bridge_token.functions.totalSupply().call() == mint_amount
    assert bridge_token.functions.symbol().call() == token_symbol
    assert bridge_token.address == address

    return bridge_token

def wait_for_sif_balance_change(sif_addr, old_balances, polling_time=1, timeout=90):
    advance_block_truffle(50)
    start_time = time.time()
    result = None
    while result is None:
        new_balances = get_sifchain_balance(sif_addr)
        if not sif_balances_equal(old_balances, new_balances):
            return new_balances
        time.sleep(polling_time)
        now = time.time()
        if now - start_time > timeout:
            raise Exception("Timeout waiting for sif balance to change")

def test_blocklist_eth(basic_transfer_request: EthereumToSifchainTransferRequest, source_ethereum_address: str):
    w3 = get_web3_connection_for_test()
    return _test_blocklist_eth(w3, source_ethereum_address)

def _test_blocklist_eth(w3, source_ethereum_address):
    default_account = w3.eth.accounts[0]  # Should be deployer
    assert default_account == test_utilities.get_required_env_var("OWNER"), "OWNER account is not the same as default"
    assert default_account.lower() == source_ethereum_address.lower(), "source_ethereum_address account is not the same as default"
    w3.eth.defaultAccount = default_account

    all_accounts = []
    for i in range(10):
        account = create_eth_account(w3)
        all_accounts.append(account)

    blocked_accounts = [x for x in all_accounts[:3]]
    nonblocked_accounts = [x for x in all_accounts if x not in blocked_accounts]

    amount_to_fund = ETHER
    amount_to_lock = ETHER // 100
    assert w3.eth.get_balance(default_account) > len(all_accounts) * amount_to_fund, \
        f"Source account {default_account} has insufficient ether balance"

    # Transfer 1 eth to every account
    for acct in all_accounts:
        start_balance = w3.eth.get_balance(acct)
        send_ether(w3, default_account, acct, amount_to_fund)
        assert w3.eth.get_balance(acct) == start_balance + amount_to_fund

    blocklist_sc = get_blocklist_sc(w3)

    set_blocklist_to(w3, blocklist_sc, [])
    currently_blocked = blocklist_sc.functions.getFullList().call()
    assert len(currently_blocked) == 0

    set_blocklist_to(w3, blocklist_sc, blocked_accounts)
    currently_blocked = blocklist_sc.functions.getFullList().call()
    assert len(currently_blocked) == len(blocked_accounts)
    assert set(currently_blocked) == set(blocked_accounts)

    sif_acct1 = create_sifchain_addr()
    sif_symbol = "ceth"

    bridge_bank_sc = get_bridge_bank_sc(w3)

    for acct in all_accounts:
        if acct in blocked_accounts:
            try:
                lock_eth(w3, bridge_bank_sc, acct, sif_acct1, amount_to_lock)
                assert False
            except ContractLogicError as e:
                # Valid negative test outcome: transaction rejected with the string "Address is blocklisted"
                assert "Address is blocklisted" in e.args[0]
        else:
            # Valid positive test outcome: event emitted, optionally: funds appear on sifchain
            balances_before = get_sifchain_balance(sif_acct1)
            txrcpt3 = lock_eth(w3, bridge_bank_sc, acct, sif_acct1, amount_to_lock)
            balances_after = wait_for_sif_balance_change(sif_acct1, balances_before)
            assert balances_after.get(sif_symbol, 0) == balances_before.get(sif_symbol, 0) + amount_to_lock

def test_blocklist_erc20(basic_transfer_request: EthereumToSifchainTransferRequest, source_ethereum_address: str):
    web3_host = "127.0.0.1"
    web3_port = 7545
    # w3 = web3_connect_ws(web3_host, web3_port)

    w3 = get_web3_connection_for_test()

    return _test_blocklist_erc20(w3, basic_transfer_request, source_ethereum_address)

def _test_blocklist_erc20(w3, basic_transfer_request, source_ethereum_address):
    # For ERC20 tokens, we need to create a new instance of Blocklist smart contract, deploy it and whitelist it with
    # BridgeBank. In peggy1, the token matching in BridgeBank is done by symbol, so we need to give our token a unique
    # symbol such as TEST or MOCK + random suffix + call updateEthWtiteList() + mint() + approve().
    # See smart-contracts/test/test_bridgeBank.js:131-160 for example.

    default_account = w3.eth.accounts[0]  # Should be deployer
    w3.eth.defaultAccount = default_account

    # Must not exist on EVM yet
    eth_token_symbol = "MOCK" + random_string(6)
    sif_token_symbol = "c" + eth_token_symbol.lower()

    bridge_token = deploy_bridge_token_for_testing(w3, eth_token_symbol, default_account, 1000)
    bridge_bank = get_bridge_bank_sc(w3)

    symbol = bridge_token.functions.symbol().call()

    acct1 = create_eth_account(w3)
    acct2 = create_eth_account(w3)

    amount_to_fund = 2 * ETHER
    send_ether(w3, default_account, acct1, amount_to_fund)
    send_ether(w3, default_account, acct2, amount_to_fund)

    # b0 = [bridge_token.functions.balanceOf(x).call() for x in [bridge_token.address, default_account, acct1, acct2]]
    #
    # txhash = bridge_token.functions.transfer(acct1, 20).transact()
    # txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
    #
    # b1 = [bridge_token.functions.balanceOf(x).call() for x in [bridge_token.address, default_account, acct1, acct2]]
    #
    # txhash = bridge_token.functions.transfer(acct2, 15).transact()
    # txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
    #
    # b2 = [bridge_token.functions.balanceOf(x).call() for x in [bridge_token.address, default_account, acct1, acct2]]
    #
    # try:
    #     txhash = bridge_token.functions.transfer(acct2, 11).transact({"from": acct1, "gas": 30000})
    #     txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
    #     assert False, "Should fail as only 10 tokens are available"
    # except:
    #     pass
    #
    # txhash = bridge_token.functions.transfer(acct2, 10).transact({"from": acct1, "gas": 50000})
    # txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
    #
    # b3 = [bridge_token.functions.balanceOf(x).call() for x in [bridge_token.address, default_account, acct1, acct2]]

    txhash = bridge_token.functions.transfer(acct1, 10).transact()
    txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
    txhash = bridge_token.functions.transfer(acct2, 10).transact()
    txrcpt = w3.eth.wait_for_transaction_receipt(txhash)

    # Set allowance for BridgeBank to send 10 tokens on behalf of acct1
    txhash = bridge_token.functions.approve(bridge_bank.address, 10).transact({"from": acct1})
    txrcpt = w3.eth.wait_for_transaction_receipt(txhash)

    blocklist_sc = get_blocklist_sc(w3)
    set_blocklist_to(w3, blocklist_sc, [acct2])

    to_sif_acct = create_sifchain_addr()

    coin_denom = bridge_token.address

    # At this point the token needs to be whitelisted, if not:
    # "revert Only token in whitelist can be transferred to cosmos"
    # TODO First we try to remove it from whitelist and it needs to fail.

    txhash = bridge_bank.functions.updateEthWhiteList(coin_denom, True).transact()
    txrcpt = w3.eth.wait_for_transaction_receipt(txhash)

    def bridgebank_lock_erc20(from_addr, to_sif_addr, token_contract_addr, amount):
        entries = filter.get_new_entries()
        assert len(entries) == 0

        nonce = w3.eth.get_transaction_count(from_addr)

        # When transfering ERC20, the amount needs to be passed as argument, and the "message.value" should be 0
        # nonce seems to be not neccessary, but it is in sendLockTx.js
        recipient = to_sif_addr.encode("UTF-8")
        txhash = bridge_bank.functions.lock(recipient, token_contract_addr, amount) \
            .transact({"from": from_addr, "gas": 200000, "nonce": nonce}) # + {"nonce": nonce}
        txrcpt = w3.eth.wait_for_transaction_receipt(txhash)

        entries = filter.get_new_entries()
        assert len(entries) == 1
        assert entries[0].event == "LogLock"
        assert entries[0].transactionHash == txhash
        assert entries[0].address == bridge_bank.address
        assert entries[0].args["_from"] == from_addr
        assert entries[0].args["_to"] == recipient
        assert entries[0].args["_value"] == amount

    filter = bridge_bank.events.LogLock.createFilter(fromBlock="latest")
    try:
        amount_to_send = 1

        # Should fail because of blocklist
        try:
            bridgebank_lock_erc20(acct2, to_sif_acct, bridge_token.address, amount_to_send)
            assert False
        except ValueError as e:
            assert "Address is blocklisted" in e.args[0]["message"]

        # Should succeed
        balances_before = get_sifchain_balance(to_sif_acct)
        bridgebank_lock_erc20(acct1, to_sif_acct, bridge_token.address, amount_to_send)
        balances_after = wait_for_sif_balance_change(to_sif_acct, balances_before)
        assert balances_after[sif_token_symbol] == amount_to_send
    finally:
        w3.eth.uninstall_filter(filter.filter_id)
