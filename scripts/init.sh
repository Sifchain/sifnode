#!/usr/bin/env bash
rm -rf ~/.sifnoded
#rm -rf ~/.sifnodecli
#rm -rf sifnode.log
#rm -rf testlog.log


sifnoded init test --chain-id=sifchain

#sifnoded config output json
#sifnoded config indent true
#sifnoded config trust-node true
#sifnoded config chain-id sifchain
#sifnoded config keyring-backend test

echo "Generating deterministic account - shadowfiend"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add shadowfiend --recover --keyring-backend test

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend test

sifnoded add-genesis-account $(sifnoded keys show shadowfiend -a --keyring-backend test) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,10000000000stake,1000000000cdash
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend test) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000000stake,1000000000cdash
#sifnoded add-genesis-account sif17s95c5jpc6x2l3edwh4dm8yhac68yru7a7kr3x 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000000stake,1000000000cdash

#sifnoded add-faucet 10000000000000000rowan
#
#sifnoded add-genesis-clp-admin $(sifnodecli keys show shadowfiend -a)
#sifnoded add-genesis-clp-admin $(sifnodecli keys show akasha -a)


sifnoded gentx shadowfiend 1000000rowan --chain-id my-test-chain --keyring-backend test
echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis



#contents="$(jq '.gov.voting_params.voting_period = 10' $DAEMON_HOME/config/genesis.json)" && \
#echo "${contents}" > $DAEMON_HOME/config/genesis.json