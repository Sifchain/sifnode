#!/bin/sh
#
# Sifchain: Mainnet Genesis Entrypoint.
#

#
# Configure the node.
#
setup() {
  sifgen node create "$CHAINNET" "$MONIKER" "$MNEMONIC" --peer-address "$PEER_ADDRESSES" --genesis-url "$GENESIS_URL" --with-cosmovisor
}

#
# Run the node under cosmovisor.
#
run() {
  cosmovisor start --rpc.laddr tcp://0.0.0.0:26657 --minimum-gas-prices "$GAS_PRICE"
}

setup
run
