#!/bin/bash
BASEDIR=$(pwd)
set -e
echo "apologies for this sudo, it is to delete non-persisent cryptographic keys that usually has enhanced permissions"
sudo rm -rf ./build/networks
make clean install
cd $BASEDIR/smart-contracts && yarn install
docker-compose -f $BASEDIR/build/genesis/docker-compose-ganache.yml up -d
cp .env.example .env ; truffle deploy  ;
echo 'ETHEREUM_CONTRACT_ADDRESS='$(cat build/contracts/BridgeRegistry.json | grep '"address": "0x' | awk -F ": " '/``"address": "0x`/ -F ":" {print $2}' | sed 's/.$//') 
ETHEREUM_CONTRACT_ADDRESS=$(cat build/contracts/BridgeRegistry.json | grep '"address": "0x' | awk -F ": " '/``"address": "0x`/ -F ":" {print $2}' | sed 's/.$//' | sed 's/"//g' )
cd $BASEDIR/build/
rake genesis:network:scaffold['localnet']
rake genesis:network:boot["localnet,${ETHEREUM_CONTRACT_ADDRESS},ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f 0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1 c88b703fb08cbea894b6aeff5a544fb92e78a18e19814cd85da83b71f772aa6c 388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418,ws://192.168.2.6:7545/"]
echo "======================"
echo 'if force killed remember to stop the services, remove non-running containers, network and untagged images'
echo "======================"
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker network rm genesis_sifchain
# Image built is untagged at 3.21 GB, this removes them to prevent devouring ones disk space
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")


