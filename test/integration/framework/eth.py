import logging


ETH = 10**18

def send_ether(w3, from_account, to_account, amount):
    logging.info(f"Send {amount} from {from_account} to {to_account}...")
    txhash = w3.eth.send_transaction({
        "from": from_account,
        "to": to_account,
        "value": amount,
        "gas": 30000,
    })
    return w3.eth.wait_for_transaction_receipt(txhash)

def get_eth_balance(w3, addr):
    return w3.eth.get_balance(addr)

def get_erc20_token_balance(w3, token_sc, addr):
    return token_sc.functions.balanceOf(addr).call()
