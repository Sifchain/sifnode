#!/usr/bin/env bash
rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli
rm -rf sifnode.log
rm -rf testlog.log


sifnoded init test --chain-id=sifchain

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id sifchain
sifnodecli config keyring-backend test

echo "Generating deterministic account - shadowfiend"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnodecli keys add shadowfiend --recover

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnodecli keys add akasha --recover

sifnoded add-genesis-account $(sifnodecli keys show shadowfiend -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,10000000000stake,1000000000cdash
sifnoded add-genesis-account $(sifnodecli keys show akasha -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000000stake,1000000000cdash

sifnoded add-genesis-clp-admin $(sifnodecli keys show shadowfiend -a)
sifnoded add-genesis-clp-admin $(sifnodecli keys show akasha -a)


sifnoded gentx --name shadowfiend --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis



#contents="$(jq '.gov.voting_params.voting_period = 10' $DAEMON_HOME/config/genesis.json)" && \
#echo "${contents}" > $DAEMON_HOME/config/genesis.json