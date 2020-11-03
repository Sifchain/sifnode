#!/bin/sh
#
# Sifnode entrypoint.
#

#
# Daemon.
#
start_daemon() {
  sifnoded start --rpc.laddr tcp://0.0.0.0:26657 &
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
  expect <<EOD
    spawn ebrelayer init tcp://0.0.0.0:26657 "$ETHEREUM_WEBSOCKET_ADDRESS" \
                                             "$ETHEREUM_CONTRACT_ADDRESS" \
                                             "$MONIKER" \
                                             --chain-id "$CHAINNET" \
                                             --keyring-backend file

    expect "Enter keyring passphrase:"
    send "$PASSWORD\n"
    expect "Enter keyring passphrase:"
    send "$PASSWORD\n"
    expect -timeout -1 eof
EOD
}

#
# Wait for the RPC port to be active.
#
wait_for_rpc() {
  while ! nc -z localhost 26657; do
    sleep 15
  done
}

start_daemon
start_rest_server
start_relayer
