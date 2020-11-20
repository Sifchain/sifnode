#!/bin/bash

BASEDIR=$(pwd)
NETWORKDIR=$BASEDIR/deploy/networks
CONTAINER_NAME="genesis_sifnode1_1"

#
# Remove prior generations Config
#
echo "apologies for this sudo, it is to delete non-persisent cryptographic keys that usually has enhanced permissions"
sudo rm -rf $NETWORKDIR
make clean install

#
# Startup ganache-cli and deploy peggy smart contracts
#
cd $BASEDIR/smart-contracts && yarn install
docker-compose -f $BASEDIR/deploy/genesis/docker-compose-ganache.yml up -d
cp .env.example .env ; truffle deploy  ;
ETHEREUM_CONTRACT_ADDRESS=$(cat $BASEDIR/smart-contracts/build/contracts/BridgeRegistry.json | jq '.networks["5777"].address')

#
# scaffold and boot the dockerized localnet
#
rake genesis:network:scaffold['localnet']
rake genesis:network:boot["localnet,${ETHEREUM_CONTRACT_ADDRESS},ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f,ws://192.168.2.6:7545/"]

# prior command generations yaml that provides useful usernames and passwords
# wait for it to be complete
#
NETDEF=$NETWORKDIR/network-definition.yml
while [ ! -f "${NETDEF}" ]
do
  sleep 2
  echo "waiting for network-definition on deployment"
  NETDEF=$NETWORKDIR/network-definition.yml
done
PASSWORD=$(yq r $NETDEF ".password")
ADDR=$(yq r $NETDEF ".address")
echo $PASSWORD
echo $ADDR

#
# Add keys for a second account to test functions against
#
docker exec -it ${CONTAINER_NAME} bash -c "/test-scripts/add-second-account.sh"

#
# Wait for the Websocket subscriptions to be initialized (like 10 seconds)
#
SUBSCRIBED=$(docker logs ${CONTAINER_NAME} | grep "Subscribed")
while [ ! -n "$SUBSCRIBED" ];
do
  sleep 2
  echo "waiting for websocket subscription"
  SUBSCRIBED=$(docker logs ${CONTAINER_NAME} | grep "Subscribed")
done

#
# Transfer Eth into Ceth in our validator account
#
cd $BASEDIR/smart-contracts
yarn peggy:lock ${ADDR} 0x0000000000000000000000000000000000000000 1000000000000000000

#
# Transfer Eth into Ceth on our User account 
# This also makes the account visible to sifnodecli q auth account <addr>

USER1ADDR=$(yq r $NETDEF "[1].address")
echo $USER1ADDR
sleep 5
yarn peggy:lock ${USER1ADDR} 0x0000000000000000000000000000000000000000 1000000000000000000
sleep 5

#
# Transfer Rowan from validator account to user account
#
docker exec -it ${CONTAINER_NAME} bash -c "/test-scripts/add-rowan-to-second-account.sh"
sleep 5

#
# Run the python tests
#
docker exec -it ${CONTAINER_NAME} bash -c "python3 /test-scripts/peggy-basic-test-docker.py"

#
# tail logs
# 
docker logs -f ${CONTAINER_NAME}

#
# killing script will not end network use stop-integration-env.sh for that
#
