#!/bin/bash

#
# Runs ebrelayer.  Normally, this is run by sifchain_start_ebrelayer.sh;
# that file sets up the logs and runs sifchain_run_ebrelayer in the background.
# Normally, you don't run this script directly.
#

set -e

. $TEST_INTEGRATION_DIR/vagrantenv.sh

#
# Wait for the RPC port to be active.
#
wait_for_rpc() {
  while ! nc -z localhost 26657; do
    sleep 1
  done
}

set -x

wait_for_rpc

echo TEST_INTEGRATION_DIR is $TEST_INTEGRATION_DIR
USER1ADDR=nothing python3 $TEST_INTEGRATION_PY_DIR/wait_for_sif_account.py $NETDEF_JSON $OWNER_ADDR

echo ETHEREUM_WEBSOCKET_ADDRESS $ETHEREUM_WEBSOCKET_ADDRESS
echo ETHEREUM_CONTRACT_ADDRESS $ETHEREUM_CONTRACT_ADDRESS
echo MONIKER $MONIKER
echo MNEMONIC $MNEMONIC

if [ -z "${EBDEBUG}" ]; then
  runner=ebrelayer
else
  cd $BASEDIR/cmd/ebrelayer
  runner="dlv exec $GOBIN/ebrelayer -- "
fi

$runner init tcp://0.0.0.0:26657 "$ETHEREUM_WEBSOCKET_ADDRESS" \
  "$ETHEREUM_CONTRACT_ADDRESS" \
  "$MONIKER" \
  "$MNEMONIC" \
  --chain-id "$CHAINNET" \
  --home $CHAINDIR/.sifnodecli \
  --gas 5000000000000 \
  --gas-prices 0.5rowan