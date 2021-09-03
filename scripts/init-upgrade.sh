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


sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2
sifnoded add-genesis-account $(sifnoded keys show sif -a) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink
sifnoded add-genesis-account $(sifnoded keys show akasha -a) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink
sifnoded add-genesis-account $(sifnoded keys show mkey -a) 500000000000000000000000rowan

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a)
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a)

sifnoded  add-genesis-validators $(sifnoded keys show sif -a --bech val)

sifnoded gentx --name sif --amount 1000000000000000000000000stake --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis


mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/release-20210414000000/bin

cp $GOPATH/bin/old/sifnoded $DAEMON_HOME/cosmovisor/genesis/bin
cp $GOPATH/bin/sifnoded $DAEMON_HOME/cosmovisor/upgrades/release-20210414000000/bin/

#contents="$(jq '.gov.voting_params.voting_period = 10' $DAEMON_HOME/config/genesis.json)" && \
#echo "${contents}" > $DAEMON_HOME/config/genesis.json
