from siftool.common import *

akash_binary = "akash"

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
