#!/usr/bin/env bash

# get the peggy need token, it is a private repo
git clone -b develop https://github.com/Sifchain/peggy

# 
cd peggy
make install
cd testnet-contracts
cp .env.example .env
yarn

yarn develop

# in other console
yarn migrate
yarn peggy:setup

# start the sifchain
./init.sh 

# start the relayer
cd peggy
cp .env.example .env
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 user1 --chain-id=sifchain
