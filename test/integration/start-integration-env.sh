#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure

# add 18 zeros to a number
to_wei () { echo "${1}000000000000000000" ; }

BASEDIR=$(pwd)/$(dirname $0)/../..
NETWORKDIR=$BASEDIR/deploy/networks
CONTAINER_NAME="integration_sifnode1_1"

envexportfile=$BASEDIR/test/integration/vagrantenv.sh
rm -f $envexportfile
echo "export BASEDIR=$BASEDIR" >> $envexportfile

#
# Remove prior generations Config
#
echo "apologies for this sudo, it is to delete non-persisent cryptographic keys that usually has enhanced permissions"
sudo rm -rf $NETWORKDIR && mkdir $NETWORKDIR
rm -rf ${BASEDIR}/smart-contracts/build ${BASEDIR}/smart-contracts/.openzeppelin
make -C ${BASEDIR} install

# ===== Everything from here on down is executed in the $BASEDIR/smart-contracts directory
cd $BASEDIR/smart-contracts

# Startup ganache-cli (https://github.com/trufflesuite/ganache)

yarn --cwd $BASEDIR/smart-contracts install

docker-compose --project-name genesis -f $BASEDIR/test/integration/docker-compose-ganache.yml up -d --force-recreate

# deploy peggy smart contracts
if [ ! -f .env ]; then
  # if you haven't created a .env file, use .env.example
  cp .env.example .env
fi

# https://www.trufflesuite.com/docs/truffle/overview
# and note that truffle migrate and truffle deploy are the same command
truffle compile
truffle deploy --network develop --reset
ETHEREUM_CONTRACT_ADDRESS=$(cat $BASEDIR/smart-contracts/build/contracts/BridgeRegistry.json | jq '.networks["5777"].address')
if [ -z "$ETHEREUM_CONTRACT_ADDRESS" ]; then
  echo ETHEREUM_CONTRACT_ADDRESS cannot be empty
  exit 1
fi
echo "export ETHEREUM_CONTRACT_ADDRESS=$ETHEREUM_CONTRACT_ADDRESS" >> $envexportfile

export BRIDGE_BANK_ADDRESS=$(cat $BASEDIR/smart-contracts/build/contracts/BridgeBank.json | jq '.networks["5777"].address')
if [ -z "BRIDGE_BANK_ADDRESS" ]; then
  echo BRIDGE_BANK_ADDRESS cannot be empty
  exit 1
fi
echo "export BRIDGE_BANK_ADDRESS=$BRIDGE_BANK_ADDRESS" >> $envexportfile

#
# scaffold and boot the dockerized localnet
#
BASEDIR=${BASEDIR} rake genesis:network:scaffold['localnet']
# see deploy/rake/genesis.rake for the description of the args to genesis:network:boot
# :chainnet, :eth_bridge_registry_address, :eth_keys, :eth_websocket
BASEDIR=${BASEDIR} rake genesis:network:boot["localnet,${ETHEREUM_CONTRACT_ADDRESS},c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3,ws://192.168.2.6:7545/"]

# those rake commands generate yaml that provides useful usernames and passwords
# wait for it to appear

NETDEF=$NETWORKDIR/network-definition.yml
while [ ! -f $NETWORKDIR/network-definition.yml ]
do
  sleep 2
done

PASSWORD=$(cat $NETDEF | yq r - ".password")
ADDR=$(cat $NETDEF | yq r - ".address")
echo $PASSWORD
echo $ADDR

#
# Add keys for a second account to test functions against
#
docker exec ${CONTAINER_NAME} bash -c "/test/integration/add-second-account.sh"

#
# Wait for the Websocket subscriptions to be initialized (like 10 seconds)
#
docker logs -f ${CONTAINER_NAME} | grep -m 1 "Subscribed"

#
# Transfer Eth into Ceth in our validator account
#
yarn --cwd $BASEDIR/smart-contracts peggy:lock ${ADDR} 0x0000000000000000000000000000000000000000 $(to_wei 10)

# balance:
#
# Transfer Eth into Ceth on our User account 
# This also makes the account visible to sifnodecli q auth account <addr>

export USER1ADDR=$(cat $NETDEF | yq r - "[1].address")
echo "export USER1ADDR=$USER1ADDR" >> $envexportfile

sleep 5
yarn --cwd $BASEDIR/smart-contracts peggy:lock ${USER1ADDR} 0x0000000000000000000000000000000000000000 $(to_wei 10)
sleep 5

#
# Transfer Rowan from validator account to user account
#
docker exec ${CONTAINER_NAME} bash -c "/test/integration/add-rowan-to-account.sh $(to_wei 23) user1"
sleep 5

# We need to forward the port used by ganache, since adding new network didn't allow
# using the cli
docker exec ${CONTAINER_NAME} bash -c "bash /test/integration/start-ganache-port-forwarding.sh"
docker exec ${CONTAINER_NAME} bash -c "cd /smart-contracts && yarn install"

#
# Run the python tests
#
echo run python tests
echo ADDR $ADDR
echo USER1ADDR $USER1ADDR
docker exec ${CONTAINER_NAME} bash -c ". /test/integration/vagrantenv.sh; SMART_CONTRACTS_DIR=/smart-contracts python3 /test/integration/peggy-basic-test-docker.py /network-definition.yml"
docker exec ${CONTAINER_NAME} bash -c '. /test/integration/vagrantenv.sh; SMART_CONTRACTS_DIR=/smart-contracts python3 /test/integration/peggy-e2e-test.py /network-definition.yml'

# killing script will not end network use stop-integration-env.sh for that
# and note that we just allow the github actions environment to be cleaned
# up by their scripts
