#!/usr/bin/env bash

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

sifnoded init test --chain-id=sifchain

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id sifchain
sifnodecli config keyring-backend test

sifnodecli keys add user1
sifnodecli keys add user2

sifnoded add-genesis-account $(sifnodecli keys show user1 -a) 100000000trwn,100000000stake

sifnoded gentx --name user1 --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis