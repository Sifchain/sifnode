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

sifnoded add-genesis-account $(sifnodecli keys show sif -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,10000000000stake,1000000000cdash
sifnoded add-genesis-account $(sifnodecli keys show akasha -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000000stake,1000000000cdash

sifnoded add-faucet 10000000000000000rowan

sifnoded add-genesis-clp-admin $(sifnodecli keys show sif -a)
sifnoded add-genesis-clp-admin $(sifnodecli keys show akasha -a)

sifnoded  add-genesis-validators $(sifnodecli keys show sif -a --bech val)

sifnoded gentx --name sif --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis



#mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin


#cp $GOPATH/src/old/sifnoded $DAEMON_HOME/cosmovisor/genesis/bin
#cp $GOPATH/src/old/sifnodecli $GOPATH/bin/

#mkdir -p $DAEMON_HOME/cosmovisor/upgrades/testupgrade/bin
#cp $GOPATH/src/new/sifnoded $DAEMON_HOME/cosmovisor/upgrades/testupgrade/bin/


#contents="$(jq '.gov.voting_params.voting_period = 10' $DAEMON_HOME/config/genesis.json)" && \
#echo "${contents}" > $DAEMON_HOME/config/genesis.json