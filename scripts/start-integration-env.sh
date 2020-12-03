#!/bin/bash
BASEDIR=$(pwd)
set -e
echo "apologies for this sudo, it is to delete non-persisent cryptographic keys that usually has enhanced permissions"
rm -rf ./deploy/networks
make clean install
cd $BASEDIR/smart-contracts && yarn install
docker-compose -f $BASEDIR/deploy/genesis/docker-compose-ganache.yml up -d
cp .env.example .env ; truffle deploy  ;
echo 'ETHEREUM_CONTRACT_ADDRESS='$(cat build/contracts/BridgeRegistry.json | grep '"address": "0x' | awk -F ": " '/``"address": "0x`/ -F ":" {print $2}' | sed 's/.$//') 
ETHEREUM_CONTRACT_ADDRESS=$(cat build/contracts/BridgeRegistry.json | grep '"address": "0x' | awk -F ": " '/``"address": "0x`/ -F ":" {print $2}' | sed 's/.$//' | sed 's/"//g' )
cd $BASEDIR/deploy/
rake genesis:network:scaffold['localnet']
rake genesis:network:boot["localnet,${ETHEREUM_CONTRACT_ADDRESS},c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3 ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f aca63a7d8138e36b68e08f4eabeea9bc051d9cb11924393b67ecd61a0292f689 765ed5b36cf22755b16240c5552050e31db0e043ec6c1398ab7349108e32d807,ws://192.168.2.6:7545/"]
echo "======================"
echo 'if force killed remember to stop the services, remove non-running containers, network and untagged images'
echo "======================"
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker network rm genesis_sifchain
# Image built is untagged at 3.21 GB, this removes them to prevent devouring ones disk space
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")


