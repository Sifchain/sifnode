#!/usr/bin/env bash

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

sifnoded init test --chain-id=sifnode

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id namechain
sifnodecli config keyring-backend test

sifnodecli keys add shadowfiend
sifnodecli keys add akasha

sifnoded add-genesis-account $(sifnodecli keys show shadowfiend -a) 1000nametoken,100000000stake
sifnoded add-genesis-account $(sifnodecli keys show akasha -a) 1000nametoken,100000000stake

sifnoded gentx --name shadowfiend --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis