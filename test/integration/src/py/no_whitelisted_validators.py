from test_utilities import owner_addr, amount_in_wei, transact_ethereum_currency_to_sifchain_addr, print_error_message


def no_whitelisted_validators():
    print(f"adding eth to {owner_addr}, expect to fail")
    try:
        transact_ethereum_currency_to_sifchain_addr(owner_addr, "eth", amount_in_wei(10))
        print(f"should have failed to send eth to {owner_addr} - transfer should have failed")
        transferred = True
    except:
        # this is the good path; we don't want the transfer to happen
        transferred = False
        print(f"correctly failed to change balance of {owner_addr}")
    if transferred:
        print_error_message("balance should not have tranferred, we should have no validators")


no_whitelisted_validators()