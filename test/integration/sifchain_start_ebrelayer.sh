#!/bin/bash

#
# Sifnode entrypoint.
#

set -x

. $TEST_INTEGRATION_DIR/vagrantenv.sh

#
# Wait for the RPC port to be active.
#
wait_for_rpc() {
  while ! nc -z localhost 26657; do
    sleep 1
  done
}

wait_for_rpc

echo TEST_INTEGRATION_DIR is $TEST_INTEGRATION_DIR
USER1ADDR=nothing python3 $TEST_INTEGRATION_DIR/wait_for_sif_account.py $NETDEF_JSON $OWNER_ADDR

echo ETHEREUM_WEBSOCKET_ADDRESS $ETHEREUM_WEBSOCKET_ADDRESS
echo ETHEREUM_CONTRACT_ADDRESS $ETHEREUM_CONTRACT_ADDRESS
echo MONIKER $MONIKER
echo MNEMONIC $MNEMONIC

ebrelayer init tcp://0.0.0.0:26657 "$ETHEREUM_WEBSOCKET_ADDRESS" \
                                           "$ETHEREUM_CONTRACT_ADDRESS" \
                                           "$MONIKER" \
                                           "$MNEMONIC" \
                                           --chain-id "$CHAINNET" \
                                           --gas 300000 \
                                           --gas-adjustment 1.5 \
                                            --home $CHAINDIR/.sifnodecli
