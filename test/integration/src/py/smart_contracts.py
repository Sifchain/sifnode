import json
import web3
import logging
from integration_test_context import make_py_module as main


cmd = main.Integrator()


def web3_connect_ws(host, port):
    return web3.Web3(web3.Web3.WebsocketProvider("ws://{}:{}".format(host, port)))

def get_sc_abi_ganache(sc_name):
    network_id = 5777
    path = main.project_dir("smart-contracts/build/contracts/{}.json".format(sc_name))
    tmp = json.loads(cmd.read_text_file(path))
    return tmp["networks"][str(network_id)]["address"], tmp["abi"]

def get_blocklist_sc(w3):
    address, abi = get_sc_abi_ganache("Blocklist")
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

def random_string(length):
    import string, random
    chars = string.ascii_letters + string.digits
    return "".join([chars[random.randrange(len(chars))] for _ in range(length)])

def create_sifchain_addr():
    mnemonic = random_string(20)
    acct = cmd.sifnoded_keys_add_1(mnemonic)
    return acct["address"]

def test():
    w3 = web3_connect_ws("127.0.0.1", 7545)
    default_account = w3.eth.accounts[0] # Should be deployer
    w3.eth.defaultAccount = default_account

    accounts = [w3.eth.account.create() for _ in range(10)]

    all_accounts = [x.address for x in accounts]
    blocked_accounts = [x for x in all_accounts[:3]]

    balance = w3.eth.get_balance(default_account)
    gwei = 10**18
    amount_to_send = 1 * gwei
    assert balance > len(all_accounts) * amount_to_send, f"Source account {default_account} has insufficient ether balance"

    # Transfer 1 eth to every account
    for acct in all_accounts:
        logging.info(f"Send {amount_to_send} from {default_account} to {acct}...")
        txhash = w3.eth.send_transaction({
            "from": default_account,
            "to": acct,
            "value": amount_to_send,
            "gas": 30000,
        })
        txrcpt = w3.eth.wait_for_transaction_receipt(txhash)
        assert w3.eth.get_balance(acct) == amount_to_send

    blocklist_sc = get_blocklist_sc(w3)

    set_blocklist_to(w3, blocklist_sc, [])
    currently_blocked = blocklist_sc.functions.getFullList().call()
    assert len(currently_blocked) == 0

    set_blocklist_to(w3, blocklist_sc, blocked_accounts)
    currently_blocked = blocklist_sc.functions.getFullList().call()
    assert len(currently_blocked) == len(currently_blocked)

    sif_acct1 = create_sifchain_addr()



    print(repr(w3))


if __name__ == "__main__":
    test()
