#!/bin/bash

#
# Sifnode entrypoint.
#

fromBlock=$1
shift
toBlock=$1
shift

set -xe

. $TEST_INTEGRATION_DIR/vagrantenv.sh

echo ETHEREUM_WEBSOCKET_ADDRESS $ETHEREUM_WEBSOCKET_ADDRESS
echo MONIKER $MONIKER
echo MNEMONIC $MNEMONIC

if [ -z "${EBDEBUG}" ]; then
  runner=ebrelayer
else
  cd $BASEDIR/cmd/ebrelayer
  runner="dlv exec $GOBIN/ebrelayer -- "
fi

tendermintNode=tcp://0.0.0.0:26657
web3Provider="$ETHEREUM_WEBSOCKET_ADDRESS"

$runner replayEthereum $tendermintNode $web3Provider $BRIDGE_REGISTRY_ADDRESS $MONIKER "$MNEMONIC" $fromBlock $toBlock $fromBlock $toBlock
#  ebrelayer replayEthereum [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker] [validatorMnemonic] [fromBlock] [toBlock] [flags]
