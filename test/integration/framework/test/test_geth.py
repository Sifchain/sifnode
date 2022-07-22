import web3
from siftool import test_utils, eth
import siftool.main as mod_main
import siftool.geth as mod_geth
from siftool.eth import ETH



def test1():
    w3_url = "ws://localhost:8545/"
    w3_conn = eth.web3_connect(w3_url)
    current_block_number = w3_conn.eth.block_number
    print(w3_conn.eth.current_block)


def geth_proof_of_concept():
    cmd = mod_main.Integrator()
    geth = mod_geth.Geth(cmd)
    geth_datadir = cmd.mktempdir()
    geth_ipcpath = cmd.mktempfile()
    geth_ipcpath = "/tmp/geth.ipc"
    f = geth.attach_eval_fn(geth_ipcpath)
    coinbase_addr = f.coinbase_addr
    acct1_password = "password"
    acct2_password = "password"
    acct1_addr = f.create_account(acct1_password)
    acct2_addr = f.create_account(acct2_password)
    f.send(coinbase_addr, acct1_addr, 1000 * ETH)
    acct1_balance = f.get_balance(acct1_addr)
    assert acct1_balance == 1000 * ETH
    f.unlock_account(acct1_addr, acct1_password)
    f.send(acct1_addr, acct2_addr, 1 * ETH)
    acct1_balance = f.get_balance(acct1_addr)
    acct2_balance = f.get_balance(acct2_addr)
    max_gas_used = 3*10**12
    assert (acct1_balance - (1000 - 1) * ETH < 0) and (acct1_balance - (1000 - 1) * ETH > -max_gas_used)
    assert acct2_balance == 1 * ETH
