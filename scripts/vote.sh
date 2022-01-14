#!/bin/zsh

proposal_count=1
sifnoded q gov proposals --node tcp://127.0.0.1:27665 >/dev/null 2>&1
if [ $? -eq 0 ]
then
  proposal_count=2
fi

echo "vote for $proposal_count on localnet-2"
terrad tx gov submit-proposal update-client 07-tendermint-0 07-tendermint-2 x --from terra --keyring-backend test --home ~/.terranode-1  --node tcp://127.0.0.1:27665 --title "vote for $proposal_count" --description "vote for $proposal_count from localnet-2" --chain-id localnet-1 --deposit 100000000uluna --broadcast-mode block --yes 
echo "proposal made"
terrad tx gov vote $proposal_count yes --chain-id localnet-1 --from terra --keyring-backend test --home ~/.terranode-1 --node tcp://127.0.0.1:27665 --yes --broadcast-mode block 
echo "sleeping"
sleep 30
terrad tx gov vote $proposal_count yes --chain-id localnet-1 --from akasha --keyring-backend test --home ~/.terranode-1 --node tcp://127.0.0.1:27665 --yes --broadcast-mode block 
echo "done"

echo "vote for $proposal_count on localnet-3"
sifnoded tx gov submit-proposal update-client 07-tendermint-0 07-tendermint-2 --from sif --keyring-backend test --home ~/.sifnode-3  --node tcp://127.0.0.1:27667 --title "vote for $proposal_count" --description "vote for $proposal_count" --chain-id localnet-3 --deposit 100000000stake --broadcast-mode block --yes 
echo "proposal made"
sifnoded tx gov vote $proposal_count yes --chain-id localnet-3 --from sif --keyring-backend test --home ~/.sifnode-3 --node tcp://127.0.0.1:27667 --yes --broadcast-mode block 
echo "sleeping"
sleep 30
sifnoded tx gov vote $proposal_count yes --chain-id localnet-3 --from akasha --keyring-backend test --home ~/.sifnode-3 --node tcp://127.0.0.1:27667 --yes --broadcast-mode block 
echo "done"
