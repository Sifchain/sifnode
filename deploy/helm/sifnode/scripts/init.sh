#!/bin/sh
set -eu

: "${SIFNODE_HOME:=/data}"
: "${GENESIS_FILE:=/genesis/genesis.json}"
: "${CHAIN_ID:?CHAIN_ID is required}"
: "${MONIKER:?MONIKER is required}"

# Never reinitialize keys or overwrite genesis on a restarted persistent node.
if [ -f "$SIFNODE_HOME/.initialized" ]; then
    exit 0
fi
if [ -e "$SIFNODE_HOME/data/application.db" ]; then
    echo "Refusing to replace genesis on an existing unmarked database; review/migrate the volume first." >&2
    exit 1
fi
if [ ! -f "$SIFNODE_HOME/config/config.toml" ]; then
    sifnoded init "$MONIKER" --chain-id "$CHAIN_ID" --home "$SIFNODE_HOME"
fi
cp "$GENESIS_FILE" "$SIFNODE_HOME/config/genesis.json"
sifnoded validate-genesis --home "$SIFNODE_HOME"
touch "$SIFNODE_HOME/.initialized"
