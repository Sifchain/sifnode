#!/bin/sh
#
# Sifnode entrypoint.
#

#
# Daemon.
#
start_daemon() {
  sifnoded start --rpc.laddr tcp://0.0.0.0:26657
}

#
# Rest server.
#
start_rest_server() {
  sifnodecli rest-server --laddr tcp://0.0.0.0:1317 &
}

#
# Start relayer.
#
start_relayer() {
  wait_for_rpc
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
                                             --gas-adjustment 1.5
}

#
# Wait for the RPC port to be active.
#
wait_for_rpc() {
  while ! nc -z localhost 26657; do
    sleep 15
  done
}

# Only start the relayer if enabled.
if [ "$RELAYER_ENABLED" = "true" ]
then
  start_relayer &
fi

start_rest_server
start_daemon
