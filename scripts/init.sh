#!/usr/bin/env bash

### chain init script for development purposes only ###

make clean install
sifnoded init test --chain-id=localnet

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id localnet
sifnodecli config keyring-backend test

echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnodecli keys add sif --recover

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnodecli keys add akasha --recover

sifnoded add-genesis-account $(sifnodecli keys show sif -a) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink
sifnoded add-genesis-account $(sifnodecli keys show akasha -a) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink

sifnoded add-faucet 100000000000000000000000rowan

sifnoded add-genesis-clp-admin $(sifnodecli keys show sif -a)
sifnoded add-genesis-clp-admin $(sifnodecli keys show akasha -a)

sifnoded  add-genesis-validators $(sifnodecli keys show sif -a --bech val)

sifnoded gentx --name sif --amount 1000000000000000000000000stake --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis

#rake "genesis:sifnode:scaffold[localnet, localnode-2, 'clock flock bunker female manage special upset wear modify wire garage ship tide lock muffin tennis vicious high truck web slab medal topic borrow', '', 5D8A947E7F23D08550E0823AF07E701ECD3DE5C4@127.0.0.1:26656, http://localhost:26657/genesis]"
#rake "genesis:sifnode:scaffold[localnet, localnode-2, 'clock flock bunker female manage special upset wear modify wire garage ship tide lock muffin tennis vicious high truck web slab medal topic borrow', '', ff0dd55dffa0e67fe21e2c85c80b0c2894bf2586@52.89.19.109:26656, http://localhost:26657/genesis]"
