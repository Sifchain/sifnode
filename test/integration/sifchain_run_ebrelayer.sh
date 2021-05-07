#!/bin/bash

#
# Runs ebrelayer.  Normally, this is run by sifchain_start_ebrelayer.sh;
# that file sets up the logs and runs sifchain_run_ebrelayer in the background.
# Normally, you don't run this script directly.
#

set -e

. $TEST_INTEGRATION_DIR/vagrantenv.sh
. ${TEST_INTEGRATION_DIR}/shell_utilities.sh


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
USER1ADDR=nothing python3 $TEST_INTEGRATION_PY_DIR/wait_for_sif_account.py $NETDEF_JSON $VALIDATOR1_ADDR
sleep 10

echo ETHEREUM_WEBSOCKET_ADDRESS $ETHEREUM_WEBSOCKET_ADDRESS
echo BRIDGE_REGISTRY_ADDRESS $BRIDGE_REGISTRY_ADDRESS
echo MONIKER $MONIKER
echo MNEMONIC $MNEMONIC

if [ -z "${EBDEBUG}" ]; then
  runner=ebrelayer
else
  cd $BASEDIR/cmd/ebrelayer
  runner="dlv exec $GOBIN/ebrelayer -- "
fi

TCP_URL=tcp://0.0.0.0:26657

sifnoded keys delete $MONIKER --keyring-backend test || true
echo $MNEMONIC | sifnoded keys add $MONIKER --keyring-backend test --recover

ETHEREUM_PRIVATE_KEY=$EBRELAYER_ETHEREUM_PRIVATE_KEY $runner init $TCP_URL "$ETHEREUM_WEBSOCKET_ADDRESS" \
  "$BRIDGE_REGISTRY_ADDRESS" \
  "$MONIKER" \
  "$MNEMONIC" \
  --chain-id $CHAINNET \
  --rpc-url $TCP_URL \
  --node $TCP_URL \
  --keyring-backend test \
  --from $MONIKER \
  # --home $CHAINDIR/.sifnoded \
