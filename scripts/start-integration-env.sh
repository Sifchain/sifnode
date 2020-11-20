#!/bin/bash
BASEDIR=$(pwd)
NETWORKDIR=$BASEDIR/deploy/networks
CONTAINER_NAME="genesis_sifnode1_1"
echo "apologies for this sudo, it is to delete non-persisent cryptographic keys that usually has enhanced permissions"
sudo rm -rf $NETWORKDIR
make clean install
cd $BASEDIR/smart-contracts && yarn install
docker-compose -f $BASEDIR/deploy/genesis/docker-compose-ganache.yml up -d
cp .env.example .env ; truffle deploy  ;
ETHEREUM_CONTRACT_ADDRESS=$(cat $BASEDIR/smart-contracts/build/contracts/BridgeRegistry.json | jq '.networks["5777"].address')
echo "bridgeregistry address"
echo $ETHEREUM_CONTRACT_ADDRESS
echo "====================="
rake genesis:network:scaffold['localnet']
rake genesis:network:boot["localnet,${ETHEREUM_CONTRACT_ADDRESS},ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f,ws://192.168.2.6:7545/"]
#rake genesis:network:boot["localnet,$ETHEREUM_CONTRACT_ADDRESS,ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f 0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1 c88b703fb08cbea894b6aeff5a544fb92e78a18e19814cd85da83b71f772aa6c 388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418,ws://192.168.2.6:7545/"]


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

docker exec -it ${CONTAINER_NAME} bash -c "/test-scripts/add-second-account.sh"

SUBSCRIBED=$(docker logs ${CONTAINER_NAME} | grep "Subscribed")
while [ ! -n "$SUBSCRIBED" ];
do
  sleep 2
  echo "waiting for websocket subscription"
  SUBSCRIBED=$(docker logs ${CONTAINER_NAME} | grep "Subscribed")
done

cd $BASEDIR/smart-contracts
yarn peggy:lock ${ADDR} 0x0000000000000000000000000000000000000000 1000000000000000000

USER1ADDR=$(yq r $NETDEF "[1].address")
echo $USER1ADDR
sleep 5
yarn peggy:lock ${USER1ADDR} 0x0000000000000000000000000000000000000000 1000000000000000000
sleep 5

docker exec -it ${CONTAINER_NAME} bash -c "/test-scripts/add-rowan-to-second-account.sh"
sleep 5
docker exec -it ${CONTAINER_NAME} bash -c "python3 /test-scripts/peggy-basic-test-docker.py"
docker logs -f ${CONTAINER_NAME}
echo "======================"
echo 'if force killed remember to stop the services, remove non-running containers, network and untagged images'
echo "======================"
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker network rm genesis_sifchain
# Image built is untagged at 3.21 GB, this removes them to prevent devouring ones disk space
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")
