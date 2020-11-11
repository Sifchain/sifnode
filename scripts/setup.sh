!/bin/bash
BASEDIR=$(pwd)
make clean install
rm -rf build/networks/
set -e
cd $BASEDIR/smart-contracts && yarn install
yarn develop > chain.log &
truffle deploy  ;
echo 'ETHEREUM_CONTRACT_ADDRESS='$(cat build/contracts/BridgeRegistry.json | grep '"address": "0x' | awk -F ": " '/``"address": "0x`/ -F ":" {print $2}' | sed 's/.$//') >> $BASEDIR/.env.cicd
ETHEREUM_CONTRACT_ADDRESS=$(cat build/contracts/BridgeRegistry.json | grep '"address": "0x' | awk -F ": " '/``"address": "0x`/ -F ":" {print $2}' | sed 's/.$//' | sed 's/"//g' )
cd ../build/
rake genesis:network:scaffold['localnet']
rake genesis:network:boot["localnet,${ETHEREUM_CONTRACT_ADDRESS},0xae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f 0x0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1 0xc88b703fb08cbea894b6aeff5a544fb92e78a18e19814cd85da83b71f772aa6c 0x388c684f0ba1ef5017716adb5d21a053ea8e90277d0868337519f97bede61418,wss://127.0.0.1:7545"]
