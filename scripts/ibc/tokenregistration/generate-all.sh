#!/bin/sh

# sh ./deregister-all.sh testnet

. ./envs/$1.sh 

mkdir -p ./$SIFCHAIN_ID
rm -f ./$SIFCHAIN_ID/temp.json
rm -f ./$SIFCHAIN_ID/temp2.json
rm -f ./$SIFCHAIN_ID/tokenregistry.json

sifnoded q tokenregistry add-all ./$SIFCHAIN_ID/registry.json | jq > $SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/cosmos.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/akash.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/sentinel.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/iris.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/persistence.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/crypto-org.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/regen.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/terra.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/osmosis.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/juno.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/ixo.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/emoney.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/likecoin.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/bitsong.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/band.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/tokenregistry.json ./$SIFCHAIN_ID/emoney-eeur.json | jq > $SIFCHAIN_ID/temp.json
rm ./$SIFCHAIN_ID/tokenregistry.json
sifnoded q tokenregistry add ./$SIFCHAIN_ID/temp.json ./$SIFCHAIN_ID/terra-uusd.json | jq > $SIFCHAIN_ID/tokenregistry.json
rm ./$SIFCHAIN_ID/temp.json