#!/usr/bin/env bash

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

sifnoded init test --chain-id=namechain

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id namechain
sifnodecli config keyring-backend test

sifnodecli keys add jack
sifnodecli keys add alice

sifnoded add-genesis-account $(sifnodecli keys show jack -a) 1000nametoken,100000000stake
sifnoded add-genesis-account $(sifnodecli keys show alice -a) 1000nametoken,100000000stake

sifnoded gentx --name jack --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis