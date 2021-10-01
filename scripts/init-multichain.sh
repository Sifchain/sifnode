#!/usr/bin/env bash

### chain init script for development purposes only ###
killall sifnoded
rm -rf ~/.sifnode-1
rm -rf ~/.sifnode-2
rm -rf ~/.sifnode-3
make clean install
sifnoded init test --chain-id=localnet-1 -o 

echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test 
echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend=test 

sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test 

sifnoded add-genesis-account $(sifnoded keys show sif -a --keyring-backend=test ) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test 
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend=test ) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test 

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a --keyring-backend=test ) --keyring-backend=test 
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a --keyring-backend=test  ) --keyring-backend=test 
sifnoded set-genesis-whitelister-admin $(sifnoded keys show sif -a --keyring-backend=test ) --keyring-backend=test 
sifnoded set-gen-denom-whitelist scripts/denoms.json 

sifnoded add-genesis-validators $(sifnoded keys show sif -a --bech val --keyring-backend=test ) --keyring-backend=test 

sifnoded gentx sif 1000000000000000000000000stake --keyring-backend=test  --chain-id=localnet-1

echo "Collecting genesis txs..."
sifnoded collect-gentxs 

echo "Validating genesis file..."
sifnoded validate-genesis 

sleep 1
sifnoded start  --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9090 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:27665 >> abci_1.log 2>&1  &