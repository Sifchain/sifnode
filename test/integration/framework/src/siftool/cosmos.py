import datetime
from typing import Union, Iterable, Mapping, List
from siftool.common import *

akash_binary = "akash"

LegacyBalance = List[List[Union[int, str]]]  # e.g. [[3, "rowan"], [2, "ibc/xxxxx"]]
Balance = Mapping[str, int]
CompatBalance = Union[LegacyBalance, Balance]
Address = str
Bank = Mapping[Address, Balance]
BechAddress = str
KeyName = str  # Name of key in the keyring


def balance_normalize(bal: CompatBalance = None) -> Balance:
    if type(bal) == list:
        bal = dict(((k, v) for v, k in bal))
    elif type(bal) == dict:
        pass
    else:
        assert False, "Balances should be either a dict or a list"
    return {k: v for k, v in bal.items() if v != 0}


def balance_add(*bal: Balance) -> Balance:
    result = {}
    all_denoms = set(flatten([*b.keys()] for b in bal))
    for denom in all_denoms:
        val = sum(b.get(denom, 0) for b in bal)
        if val != 0:
            result[denom] = val
    return result


def balance_mul(bal: Balance, multiplier: Union[int, float]) -> Balance:
    result = {}
    for denom, value in bal.items():
        val = value * multiplier
        if val != 0:
            result[denom] = val
    return result


def balance_neg(bal: Balance) -> Balance:
    return {k: -v for k, v in bal.items()}


def balance_sub(bal1: Balance, *bal2: Balance) -> Balance:
    return balance_add(bal1, *[balance_neg(b) for b in bal2])


def balance_zero(bal: Balance) -> bool:
    return len(bal) == 0


def balance_equal(bal1: Balance, bal2: Balance) -> bool:
    return balance_zero(balance_sub(bal1, bal2))


def balance_format(bal: Balance) -> str:
    return ",".join("{}{}".format(v, k) for k, v in bal.items())


def balance_exceeds(bal: Balance, min_changes: Balance) -> bool:
    have_all = True
    for denom, required_value in min_changes.items():
        actual_value = bal.get(denom, 0)
        if required_value < 0:
            have_all &= actual_value <= required_value
        elif required_value > 0:
            have_all &= actual_value >= required_value
        else:
            assert False
    return have_all


def balance_sum_by_address(*maps_of_balances: Bank) -> Bank:
    totals = {}
    for item in maps_of_balances:
        for address, balance in item.items():
            if address not in totals:
                totals[address] = {}
                totals[address] = balance_add(totals[address], balance)
    return totals


_iso_time_patterm = re.compile("(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.)(\\d+)Z$")

def parse_iso_timestamp(strtime: str):
    m = _iso_time_patterm.match(strtime)
    assert m
    strtime = m[1] + (m[2] + "000")[:3] + "+00:00"
    return datetime.datetime.fromisoformat(strtime)

# <editor-fold>

# This is for Akash, but might be useful for other cosmos-based chains as well. (If not, it should be moved to separate
# class/module.)
# Source: https://sifchain.slack.com/archives/C01T05LPFEG/p1632822677353400?thread_ts=1632735716.332000&cid=C01T05LPFEG

def query_account_balance(cmd, account, node, chain_id):
    # account = "akash19q2swhcxkxlc6va3pz5jz42jfsfv2ly4767kj7"
    # node = "http://147.75.32.35:26657"
    # chain_id = "akash-testnet-6"
    args = [akash_binary, "query", "bank", "balances", account, "--node", node, "--chain-id", chain_id]
    res = yaml_load(stdout(cmd.execst(args)))
    # balances:
    # - amount: "100000000"
    #   denom: uakt
    return res

def transfer(cmd, channel, address, amount, from_addr, chain_id, node, gas_prices, gas, packet_timeout_timestamp):
    # akash tx ibc-transfer transfer transfer channel-66
    # sif19q2swhcxkxlc6va3pz5jz42jfsfv2ly4kuu8y0
    # 100ibc/10CD333A451FAE602172F612E6F0D695476C8A0C4BEC6E0A9F1789A599B9F135
    # --from akash19q2swhcxkxlc6va3pz5jz42jfsfv2ly4767kj7
    # --keyring-backend test
    # --chain-id akash-testnet-6
    # --node http://147.75.32.35:26657
    # -y --gas-prices 2.0uakt --gas 500000 --packet-timeout-timestamp 600000000000
    # channel = "channel-66"
    # address = "sif19q2swhcxkxlc6va3pz5jz42jfsfv2ly4kuu8y0
    # amount = "100ibc/10CD333A451FAE602172F612E6F0D695476C8A0C4BEC6E0A9F1789A599B9F135"
    # from_addr = "akash19q2swhcxkxlc6va3pz5jz42jfsfv2ly4767kj7"
    # chain_id = "akash-testnet-6"
    # node = "http://147.75.32.35:26657"
    # gas_prices = "2.0uakt"
    # gas = "500000"
    # packet_timeout_timestamp = 600000000000
    keyring_backend = "test"
    args = [akash_binary, "tx", "ibc-transfer", "transfer", "transfer", channel,
        address, amount, "--from", from_addr, "--keyring-backend", keyring_backend,
        "--chain-id", chain_id, "--node", node, "-y", "--gas-prices", gas_prices,
        "--gas", gas, "--packet-timeout-timestam[", str(packet_timeout_timestamp)]
    res = cmd.execst(args)
    return res

# </editor-fold>
