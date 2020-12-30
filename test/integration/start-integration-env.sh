#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure
set -o pipefail

BASEDIR=$(pwd)/$(dirname $0)/../..

. ${BASEDIR}/test/integration/shell_utilities.sh

export envexportfile=$BASEDIR/test/integration/vagrantenv.sh
rm -f $envexportfile

set_persistant_env_var ETHEREUM_PRIVATE_KEY c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3 $envexportfile
set_persistant_env_var BASEDIR $(fullpath $BASEDIR) $envexportfile
set_persistant_env_var SIFCHAIN_BIN $BASEDIR/cmd $envexportfile
set_persistant_env_var envexportfile $(fullpath $envexportfile) $envexportfile
set_persistant_env_var TEST_INTEGRATION_DIR ${BASEDIR}/test/integration $envexportfile
set_persistant_env_var SMART_CONTRACTS_DIR ${BASEDIR}/smart-contracts $envexportfile
set_persistant_env_var datadir ${TEST_INTEGRATION_DIR}/vagrant/data $envexportfile
set_persistant_env_var CONTAINER_NAME integration_sifnode1_1 $envexportfile
set_persistant_env_var NETWORKDIR $BASEDIR/deploy/networks $envexportfile
set_persistant_env_var GANACHE_DB_DIR $(mktemp -d --tmpdir ganachedb.XXXX) $envexportfile
set_persistant_env_var ETHEREUM_WEBSOCKET_ADDRESS ws://localhost:7545/ $envexportfile
set_persistant_env_var CHAINNET localnet $envexportfile

mkdir -p $datadir

#rm -rf ${BASEDIR}/smart-contracts/build ${BASEDIR}/smart-contracts/.openzeppelin
make -C ${BASEDIR} install

# ===== Everything from here on down is executed in the $BASEDIR/smart-contracts directory
cd $BASEDIR/smart-contracts

# Startup ganache-cli (https://github.com/trufflesuite/ganache)

cp ${TEST_INTEGRATION_DIR}/.env.ciExample .env

yarn --cwd $BASEDIR/smart-contracts install
#set_persistant_env_var YARN_CACHE_DIR $(yarn cache dir) $envexportfile

# Uses GANACHE_DB_DIR for the --db argument to the chain
docker-compose --project-name genesis -f ${TEST_INTEGRATION_DIR}/docker-compose-ganache.yml up -d --force-recreate

# https://www.trufflesuite.com/docs/truffle/overview
# and note that truffle migrate and truffle deploy are the same command
truffle deploy --network develop --reset

# ETHEREUM_CONTRACT_ADDRESS is used for the BridgeRegistry address in many places, so we
# set it and BRIDGE_REGISTRY_ADDRESS to the same value
echo "# BRIDGE_REGISTRY_ADDRESS and ETHEREUM_CONTRACT_ADDRESS are synonyms">> $envexportfile
set_persistant_env_var BRIDGE_REGISTRY_ADDRESS $(cat $BASEDIR/smart-contracts/build/contracts/BridgeRegistry.json | jq -r '.networks["5777"].address') $envexportfile required
set_persistant_env_var ETHEREUM_CONTRACT_ADDRESS $BRIDGE_REGISTRY_ADDRESS $envexportfile required

set_persistant_env_var BRIDGE_BANK_ADDRESS $(cat $BASEDIR/smart-contracts/build/contracts/BridgeBank.json | jq -r '.networks["5777"].address') $envexportfile required

ADD_VALIDATOR_TO_WHITELIST=true bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile

#
# Add keys for a second account to test functions against

yes $OWNER_PASSWORD | sifnodecli keys add user1 --home $CHAINDIR/.sifnodecli || true
yes $OWNER_PASSWORD | sifnodecli keys show user1 --home $CHAINDIR/.sifnodecli >> $NETDEF || true
cat $NETDEF | to_json > $NETDEF_JSON
set_persistant_env_var USER1ADDR $(cat $NETDEF_JSON | jq -r ".[1].address") $envexportfile

# TODO the python tests use a lot of calls that require the owner password,
# and they look it up every time.  Need to change that so they just have the
# password passed in, but for now modify the user's ~/.sifnodecli directly

rm -rf ~/.sifnodecli
ln -s $CHAINDIR/.sifnodecli ~

logecho finished $0