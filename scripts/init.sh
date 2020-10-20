#!/usr/bin/env bash

killall sifnoded sifnodecli

rm $(which sifnoded) 2> /dev/null || echo sifnoded not install yet ...
rm $(which sifnodecli) 2> /dev/null || echo sifnodecli not install yet ...

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

make install

sifnoded init test --chain-id=sifchain

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id sifchain
sifnodecli config keyring-backend test

sifnodecli keys add user1
sifnodecli keys add user2

sifnoded add-genesis-account $(sifnodecli keys show user1 -a) 1000rwn,100000000stake
sifnoded add-genesis-account $(sifnodecli keys show user2 -a) 1000rwn,100000000stake

sifnoded gentx --name user1 --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis
