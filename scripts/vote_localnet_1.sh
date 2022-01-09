#!/bin/zsh

echo "vote for $1 on localnet-1"
sifnoded tx gov submit-proposal update-client 07-tendermint-0 07-tendermint-2 --from sif --keyring-backend test --home ~/.sifnode-1  --node tcp://127.0.0.1:27665 --title "vote for $1" --description "vote for $1" --chain-id localnet-1  --deposit 100000000stake --broadcast-mode block --yes 
echo "proposal made"
sifnoded tx gov vote $1 yes --chain-id localnet-1 --from sif --keyring-backend test --home ~/.sifnode-1 --node tcp://127.0.0.1:27665 --yes --broadcast-mode block 
echo "sleeping"
sleep 30
sifnoded tx gov vote $1 yes --chain-id localnet-1 --from akasha --keyring-backend test --home ~/.sifnode-1 --node tcp://127.0.0.1:27665 --yes --broadcast-mode block 
echo "done"

