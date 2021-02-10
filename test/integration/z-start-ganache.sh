#!/bin/bash
# must run from the root directory of the sifnode tree

set -xv
set -e # exit on any failure
set -o pipefail

BASEDIR=$(pwd)/$(dirname $0)/../..

. ${BASEDIR}/test/integration/shell_utilities.sh

# ===== Everything is executed with a working directory of $BASEDIR/smart-contracts
cd $BASEDIR/smart-contracts

export envexportfile=$BASEDIR/test/integration/vagrantenv.sh
rm -f $envexportfile

set_persistant_env_var ETHEREUM_PRIVATE_KEY c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3 $envexportfile
set_persistant_env_var OWNER 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 $envexportfile
# we may eventually switch things so PAUSER and OWNER aren't the same account, but for now they're the same
set_persistant_env_var PAUSER $OWNER $envexportfile
set_persistant_env_var BASEDIR $(fullpath $BASEDIR) $envexportfile
set_persistant_env_var SIFCHAIN_BIN $BASEDIR/cmd $envexportfile
set_persistant_env_var envexportfile $(fullpath $envexportfile) $envexportfile
set_persistant_env_var TEST_INTEGRATION_DIR ${BASEDIR}/test/integration $envexportfile
set_persistant_env_var TEST_INTEGRATION_PY_DIR ${BASEDIR}/test/integration/src/py $envexportfile
set_persistant_env_var SMART_CONTRACTS_DIR ${BASEDIR}/smart-contracts $envexportfile
set_persistant_env_var datadir ${TEST_INTEGRATION_DIR}/vagrant/data $envexportfile
set_persistant_env_var CONTAINER_NAME integration_sifnode1_1 $envexportfile
set_persistant_env_var NETWORKDIR $BASEDIR/deploy/networks $envexportfile
set_persistant_env_var GANACHE_DB_DIR $(mktemp -d /tmp/ganachedb.XXXX) $envexportfile
set_persistant_env_var ETHEREUM_WEBSOCKET_ADDRESS ws://localhost:7545/ $envexportfile
set_persistant_env_var CHAINNET localnet $envexportfile

mkdir -p $datadir

make -C ${TEST_INTEGRATION_DIR}

cp ${TEST_INTEGRATION_DIR}/.env.ciExample .env

make -C $SMART_CONTRACTS_DIR clean-smartcontracts
yarn --cwd $BASEDIR/smart-contracts install

set -xv

block_delay=$1
if [ ! -z "$block_delay" ]
then
    block_delay="-b $block_delay"
fi

. /Users/junius/github/thor/sifnode/$(dirname $0)/vagrantenv.sh
. $TEST_INTEGRATION_DIR/shell_utilities.sh

set_persistant_env_var GANACHE_LOG $datadir/logs/ganache.$(filenamedate).txt $envexportfile
mkdir -p $(dirname $GANACHE_LOG)

pkill -9 -f ganache-cli || true
while nc -z localhost 7545; do
    sleep 1
done

# ganache really hates running in the background.  Put it in a tmux session to keep all its input code happy.
# If you don't do this, ganache-cli will just exit.
ganache-cli ${block_delay} -h 0.0.0.0 --mnemonic 'candy maple cake sugar pudding cream honey rich smooth crumble sweet treat' --networkId '5777' --port '7545' --db ${GANACHE_DB_DIR}

