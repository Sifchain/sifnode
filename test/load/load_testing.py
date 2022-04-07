import random
from typing import Sequence, Any


# Fees for sifchain -> sifchain transactions, paid by the sender.
sif_tx_fee_in_rowan = 1 * 10**17

# Fees for "ethbridge burn" transactions. Determined experimentally
sif_tx_burn_fee_in_rowan = 100000
sif_tx_burn_fee_in_ceth = 1

# There seems to be a minimum amount of rowan that a sif account needs to own in order for the bridge to do an
# "ethbridge burn". This amount does not seem to be actually used. For example, if you fund the account just with
# sif_tx_burn_fee_in_rowan, We observed that if you try to fund sif accounts with just the exact amount of rowan
# needed to pay fees (sif_tx_burn_fee_in_rowan * number_of_transactions), the bridge would stop forwarding after
# approx. 200 transactions, and you would see in sifnoded logs this message:
# {"level":"debug","module":"mempool","err":null,"peerID":"","res":{"check_tx":{"code":5,"data":null,"log":"0rowan is smaller than 500000000000000000rowan: insufficient funds: insufficient funds","info":"","gas_wanted":"1000000000000000000","gas_used":"19773","events":[],"codespace":"sdk"}},"tx":"...","time":"2022-03-26T10:09:26+01:00","message":"rejected bad transaction"}
sif_tx_burn_fee_buffer_in_rowan = 5 * sif_tx_fee_in_rowan

rowan = "rowan"


def choose_from(distr: Sequence[Any], rnd: random.Random = None) -> int:
    r = (rnd.randrange(sum(distr))) if rnd else 0
    s = 0
    for i, di in enumerate(distr):
        s += di
        if r < s:
            distr[i] -= 1
            return i
    assert False
