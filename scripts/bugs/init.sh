#!/usr/bin/env bash

### chain init script for development purposes only ###

make clean install
sifnoded init test --chain-id=localnet

sifnoded config output json
sifnoded config indent true
sifnoded config trust-node true
sifnoded config chain-id localnet
sifnoded config keyring-backend test

echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover

sifnoded add-genesis-account $(sifnoded keys show sif -a) 16205782692902021002506278400rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink,899999867990000000000000000000cacoin
sifnoded add-genesis-account $(sifnoded keys show akasha -a) 5000000000000003407464rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink,8999998679900000000000000000000cacoin

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a)
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a)

sifnoded  add-genesis-validators $(sifnoded keys show sif -a --bech val)

sifnoded gentx --name sif --amount 1000000000000000000000000stake --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis
