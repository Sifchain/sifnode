#!/bin/sh
#
# Sifnode entrypoint.
#

set -x

ADD_VALIDATOR_TO_WHITELIST=$1
shift

NETDEF=/network-definition.yml
PASSWORD=$(cat $NETDEF | yq r - ".password")

if [ -z "${ADD_VALIDATOR_TO_WHITELIST}" ]
then
  # no whitelist validator requested; mostly useful for testing validator whitelisting
  echo $0: no whitelisted validators
else
  whitelisted_validator=$(yes $PASSWORD | sifnodecli keys show -a --bech val $MONIKER)
  echo $0: whitelisted validator $whitelisted_validator
  sifnoded add-genesis-validators $whitelisted_validator
fi

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
