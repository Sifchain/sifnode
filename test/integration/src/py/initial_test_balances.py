import time

from test_utilities import amount_in_wei, \
    send_from_sifchain_to_sifchain, transact_ethereum_to_sifchain, get_required_env_var


def setup_currencies():
    owner_addr = get_required_env_var("OWNER_ADDR")
    user1_addr = get_required_env_var("USER1ADDR")
    print(f"adding eth to {owner_addr}")
    transact_ethereum_to_sifchain(owner_addr, "eth", amount_in_wei(10))
    print(f"adding eth to {user1_addr}")
    transact_ethereum_to_sifchain(user1_addr, "eth", amount_in_wei(13))
    time.sleep(15)
    send_from_sifchain_to_sifchain(owner_addr, user1_addr, amount_in_wei(23), "rowan", get_required_env_var("OWNER_PASSWORD"))


setup_currencies()