import time

from test_utilities import network_password, owner_addr, amount_in_wei, user1_addr, \
    sif_tx_send, transact_ethereum_currency_to_sifchain_addr


def setup_currencies():
    print(f"adding eth to {owner_addr}")
    transact_ethereum_currency_to_sifchain_addr(owner_addr, "eth", amount_in_wei(10))
    print(f"adding eth to {user1_addr}")
    transact_ethereum_currency_to_sifchain_addr(user1_addr, "eth", amount_in_wei(13))
    time.sleep(15)
    sif_tx_send(owner_addr, user1_addr, amount_in_wei(23), "rowan", network_password)


setup_currencies()