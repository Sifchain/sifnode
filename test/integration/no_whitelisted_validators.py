from test_utilities import owner_addr, amount_in_wei, transact_ethereum_currency_to_sifchain_addr, print_error_message


def no_whitelisted_validators():
    print(f"adding eth to {owner_addr}, expect to fail")
    try:
        transact_ethereum_currency_to_sifchain_addr(owner_addr, "eth", amount_in_wei(10))
        print("did not fail")
        transferred = True
    except:
        # this is the good path; we don't want the transfer to happen
        transferred = False
        print("did fail")
    if transferred:
        print_error_message("balance should not have tranferred, we should have no validators")


no_whitelisted_validators()