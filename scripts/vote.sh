#!/bin/zsh

echo "vote for $1"
sifnoded tx gov submit-proposal update-client 07-tendermint-1 07-tendermint-2 --from sif --keyring-backend test --home ~/.sifnode-2  --node tcp://127.0.0.1:27666 --title "vote for $1" --description "vote for $1" --chain-id localnet-2  --deposit 100000000stake --broadcast-mode block --yes --timeout-height 2000000
echo "proposal made"
sifnoded tx gov vote $1 yes --chain-id localnet-2 --from sif --keyring-backend test --home ~/.sifnode-2 --node tcp://127.0.0.1:27666 --yes --broadcast-mode block 
echo "sleeping"
sleep 30
sifnoded tx gov vote $1 yes --chain-id localnet-2 --from akasha --keyring-backend test --home ~/.sifnode-2 --node tcp://127.0.0.1:27666 --yes --broadcast-mode block 
echo "done"
