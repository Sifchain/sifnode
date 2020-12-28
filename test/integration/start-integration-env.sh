#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure

BASEDIR=$(pwd)/$(dirname $0)/../..

. ${BASEDIR}/test/integration/shell_utilities.sh

export envexportfile=$BASEDIR/test/integration/vagrantenv.sh
rm -f $envexportfile
echo "export envexportfile=$envexportfile" >> $envexportfile

# datadir contains all the telemetry about the run; docker logs, etc
export datadir=$BASEDIR/test/integration/vagrant/data
echo "export datadir=$datadir" >> $envexportfile

bash $BASEDIR/test/integration/start_watchers.sh

export CONTAINER_NAME="integration_sifnode1_1"
echo "export CONTAINER_NAME=$CONTAINER_NAME" >> $envexportfile

echo "export BASEDIR=$BASEDIR" >> $envexportfile

NETWORKDIR=$BASEDIR/deploy/networks
echo "export NETWORKDIR=$NETWORKDIR" >> $envexportfile

#
# Remove prior generations Config
#
sudo rm -rf $NETWORKDIR && mkdir $NETWORKDIR
rm -rf ${BASEDIR}/smart-contracts/build ${BASEDIR}/smart-contracts/.openzeppelin
make -C ${BASEDIR} install

# ===== Everything from here on down is executed in the $BASEDIR/smart-contracts directory
cd $BASEDIR/smart-contracts

# Startup ganache-cli (https://github.com/trufflesuite/ganache)

cp $BASEDIR/test/integration/.env.ciExample .env

yarn --cwd $BASEDIR/smart-contracts install
export YARN_CACHE_DIR=$(yarn cache dir)
echo "export YARN_CACHE_DIR=$YARN_CACHE_DIR" >> $envexportfile

docker-compose --project-name genesis -f $BASEDIR/test/integration/docker-compose-ganache.yml up -d --force-recreate

# https://www.trufflesuite.com/docs/truffle/overview
# and note that truffle migrate and truffle deploy are the same command
truffle compile
truffle deploy --network develop --reset
# ETHEREUM_CONTRACT_ADDRESS is used for the BridgeRegistry address in many places, so we
# set it and BRIDGE_REGISTRY_ADDRESS to the same value
BRIDGE_REGISTRY_ADDRESS=$(cat $BASEDIR/smart-contracts/build/contracts/BridgeRegistry.json | jq '.networks["5777"].address')
ETHEREUM_CONTRACT_ADDRESS=$BRIDGE_REGISTRY_ADDRESS
if [ -z "$ETHEREUM_CONTRACT_ADDRESS" ]; then
  echo ETHEREUM_CONTRACT_ADDRESS cannot be empty
  exit 1
fi
echo "export ETHEREUM_CONTRACT_ADDRESS=$ETHEREUM_CONTRACT_ADDRESS" >> $envexportfile
echo "# BRIDGE_REGISTRY_ADDRESS and ETHEREUM_CONTRACT_ADDRESS are synonyms">> $envexportfile
echo "export BRIDGE_REGISTRY_ADDRESS=$BRIDGE_REGISTRY_ADDRESS" >> $envexportfile

export BRIDGE_BANK_ADDRESS=$(cat $BASEDIR/smart-contracts/build/contracts/BridgeBank.json | jq '.networks["5777"].address')
if [ -z "BRIDGE_BANK_ADDRESS" ]; then
  echo BRIDGE_BANK_ADDRESS cannot be empty
  exit 1
fi
echo "export BRIDGE_BANK_ADDRESS=$BRIDGE_BANK_ADDRESS" >> $envexportfile

ADD_VALIDATOR_TO_WHITELIST=1 bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile

docker exec ${CONTAINER_NAME} bash -c "cd /smart-contracts && yarn install"

#
# Add keys for a second account to test functions against
#
docker exec ${CONTAINER_NAME} bash -c "/test/integration/add-second-account.sh"

export USER1ADDR=$(cat $NETDEF | yq r - "[1].address")
echo "export USER1ADDR=$USER1ADDR" >> $envexportfile

#
# Run the python tests
#
echo run python tests

docker exec ${CONTAINER_NAME} bash -c ". /test/integration/vagrantenv.sh; cd /sifnode; SMART_CONTRACTS_DIR=/smart-contracts python3 /test/integration/initial_test_balances.py /network-definition.yml"
sleep 15
docker exec ${CONTAINER_NAME} bash -c ". /test/integration/vagrantenv.sh; cd /sifnode; SMART_CONTRACTS_DIR=/smart-contracts python3 /test/integration/peggy-basic-test-docker.py /network-definition.yml"
docker exec ${CONTAINER_NAME} bash -c '. /test/integration/vagrantenv.sh; cd /sifnode; SMART_CONTRACTS_DIR=/smart-contracts python3 /test/integration/peggy-e2e-test.py /network-definition.yml'

# Rebuild sifchain, but this time don't use validators

sudo rm -rf $NETWORKDIR && mkdir $NETWORKDIR
ADD_VALIDATOR_TO_WHITELIST= bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile

docker exec ${CONTAINER_NAME} bash -c ". /test/integration/vagrantenv.sh; cd /sifnode; SMART_CONTRACTS_DIR=/smart-contracts python3 /test/integration/no_whitelisted_validators.py /network-definition.yml"
