#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure
set -o pipefail

BASEDIR=$(pwd)/$(dirname $0)/../..

. ${BASEDIR}/test/integration/shell_utilities.sh

# ===== Everything is executed with a working directory of $BASEDIR/smart-contracts
cd $BASEDIR/smart-contracts

export envexportfile=$BASEDIR/test/integration/vagrantenv.sh
rm -f $envexportfile

set_persistant_env_var ETHEREUM_PRIVATE_KEY c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3 $envexportfile
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

# Startup ganache-cli (https://github.com/trufflesuite/ganache)
#   Uses GANACHE_DB_DIR for the --db argument to the chain
docker-compose --project-name genesis -f ${TEST_INTEGRATION_DIR}/docker-compose-ganache.yml up -d --force-recreate

# https://www.trufflesuite.com/docs/truffle/overview
# and note that truffle migrate and truffle deploy are the same command
npx truffle deploy --network develop --reset

# ETHEREUM_CONTRACT_ADDRESS is used for the BridgeRegistry address in many places, so we
# set it and BRIDGE_REGISTRY_ADDRESS to the same value
echo "# BRIDGE_REGISTRY_ADDRESS and ETHEREUM_CONTRACT_ADDRESS are synonyms">> $envexportfile
set_persistant_env_var BRIDGE_REGISTRY_ADDRESS $(cat $BASEDIR/smart-contracts/build/contracts/BridgeRegistry.json | jq -r '.networks["5777"].address') $envexportfile required
set_persistant_env_var BRIDGE_TOKEN_ADDRESS $(cat $BASEDIR/smart-contracts/build/contracts/BridgeToken.json | jq -r '.networks["5777"].address') $envexportfile required
set_persistant_env_var ETHEREUM_CONTRACT_ADDRESS $BRIDGE_REGISTRY_ADDRESS $envexportfile required

set_persistant_env_var BRIDGE_BANK_ADDRESS $(cat $BASEDIR/smart-contracts/build/contracts/BridgeBank.json | jq -r '.networks["5777"].address') $envexportfile required

ADD_VALIDATOR_TO_WHITELIST=true bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile

UPDATE_ADDRESS=0x0000000000000000000000000000000000000000 npx truffle exec scripts/setTokenLockBurnLimit.js 31000000000000000000
UPDATE_ADDRESS=$BRIDGE_TOKEN_ADDRESS npx truffle exec scripts/setTokenLockBurnLimit.js 10000000000000000000000000

logecho finished $0
