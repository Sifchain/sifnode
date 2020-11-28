#!/usr/bin/env bash


#original="voting_period\": \"172800000000000"
#new="voting_period\": \"10"
#while read a; do
#    echo "${a//$original/$new}"
#done < $DAEMON_HOME/config/genesis.json > $DAEMON_HOME/config/genesis.json.t
#mv $DAEMON_HOME/config/genesis.json{.t,}

cosmovisor start >> sifnode.log 2>&1  &
sleep 10
yes Y | sifnodecli tx gov submit-proposal software-upgrade testupgrade --from shadowfiend --deposit 100000000stake --upgrade-height 20 --title testupgrade --description testupgrade
sleep 5
yes Y | sifnodecli tx gov vote 1 yes --from shadowfiend --keyring-backend test --chain-id sifchain
clear
sleep 5
sifnodecli query gov proposal 1

