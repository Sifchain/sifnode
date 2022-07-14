#!/usr/bin/env bash

if ! type "rly" > /dev/null; then
  # install foobar here
  echo "You need the go-ibc relayer to run this script. You can find it here https://github.com/strangelove-ventures/relayer"
  exit 0
fi

### chain init script for development purposes only ###
killall sifnoded
killall hermes
rm -rf ~/.sifnode-1
rm -rf ~/.sifnode-2
rm -rf ~/.sifnode-3
make clean install
sifnoded init test --chain-id=localnet-1 -o --home ~/.sifnode-1

echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test --home ~/.sifnode-1

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend=test --home ~/.sifnode-1


sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test --home ~/.sifnode-1

sifnoded add-genesis-account $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-1) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home ~/.sifnode-1
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-1) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home ~/.sifnode-1

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-1) --keyring-backend=test --home ~/.sifnode-1
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-1 ) --keyring-backend=test --home ~/.sifnode-1
sifnoded set-genesis-whitelister-admin $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-1) --keyring-backend=test --home ~/.sifnode-1
sifnoded set-gen-denom-whitelist scripts/denoms.json --home ~/.sifnode-1

sifnoded add-genesis-validators $(sifnoded keys show sif -a --bech val --keyring-backend=test --home ~/.sifnode-1) --keyring-backend=test --home ~/.sifnode-1

sifnoded gentx sif 1000000000000000000000000stake --keyring-backend=test --home ~/.sifnode-1 --chain-id=localnet-1

echo "Collecting genesis txs..."
sifnoded collect-gentxs --home ~/.sifnode-1

echo "Validating genesis file..."
sifnoded validate-genesis --home ~/.sifnode-1



sifnoded init test --chain-id=localnet-2 -o --home ~/.sifnode-2


echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test --home ~/.sifnode-2

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend=test --home ~/.sifnode-2


sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test --home ~/.sifnode-2

sifnoded add-genesis-account $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-2 ) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home ~/.sifnode-2
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-2) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home ~/.sifnode-2

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-2 ) --keyring-backend=test --home ~/.sifnode-2
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-2) --keyring-backend=test --home ~/.sifnode-2
sifnoded set-genesis-whitelister-admin $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-2) --keyring-backend=test --home ~/.sifnode-2
sifnoded set-gen-denom-whitelist scripts/denoms.json --home ~/.sifnode-2
sifnoded add-genesis-validators $(sifnoded keys show sif -a --bech val --keyring-backend=test --home ~/.sifnode-2 ) --keyring-backend=test --home ~/.sifnode-2

sifnoded gentx sif 1000000000000000000000000stake --chain-id=localnet --keyring-backend=test --home ~/.sifnode-2 --chain-id=localnet-2

echo "Collecting genesis txs..."
sifnoded collect-gentxs --home ~/.sifnode-2

echo "Validating genesis file..."
sifnoded validate-genesis --home ~/.sifnode-2



sifnoded init test --chain-id=localnet-3 -o --home ~/.sifnode-3


echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test --home ~/.sifnode-3

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend=test --home ~/.sifnode-3


sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test --home ~/.sifnode-3

sifnoded add-genesis-account $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-3 ) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home ~/.sifnode-3
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-3) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test --home ~/.sifnode-3

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-3 ) --keyring-backend=test --home ~/.sifnode-3
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-3) --keyring-backend=test --home ~/.sifnode-3
sifnoded set-genesis-whitelister-admin $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-3) --keyring-backend=test --home ~/.sifnode-3
sifnoded set-gen-denom-whitelist scripts/denoms.json --home ~/.sifnode-3
sifnoded add-genesis-validators $(sifnoded keys show sif -a --bech val --keyring-backend=test --home ~/.sifnode-3 ) --keyring-backend=test --home ~/.sifnode-3

sifnoded gentx sif 1000000000000000000000000stake --chain-id=localnet-3 --keyring-backend=test --home ~/.sifnode-3 --chain-id=localnet-3

echo "Collecting genesis txs..."
sifnoded collect-gentxs --home ~/.sifnode-3

echo "Validating genesis file..."
sifnoded validate-genesis --home ~/.sifnode-3

rm -rf abci_*.log
rm -rf hermes.log
rm -rf ~/.hermes

echo "Starting sifnoded's"

sleep 1
sifnoded start --home ~/.sifnode-1 --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9090 --grpc-web.address 0.0.0.0:9093 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:27665 >> abci_1.log 2>&1  &
sleep 1
sifnoded start --home ~/.sifnode-2 --p2p.laddr 0.0.0.0:27656  --grpc.address 0.0.0.0:9091 --grpc-web.address 0.0.0.0:9094 --address tcp://0.0.0.0:27660 --rpc.laddr tcp://127.0.0.1:27666 >> abci_2.log 2>&1  &
sleep 1
sifnoded start --home ~/.sifnode-3 --p2p.laddr 0.0.0.0:27657  --grpc.address 0.0.0.0:9092 --grpc-web.address 0.0.0.0:9095 --address tcp://0.0.0.0:27661 --rpc.laddr tcp://127.0.0.1:27667 >> abci_3.log 2>&1 &
sleep 5

echo "Setting go-relayer"
rly config init
echo "Relay config init ~/.relayer"

echo "Loading config from -> config-ibc-go-relayer/localnet/chains"
echo "into -> ~/.relayer/config/config.yaml"
rly config add-chains scripts/config-ibc-go-relayer/localnet/chains
echo "Checking config loaded..."
cat ~/.relayer/config/config.yaml

echo "Restoring keys akasha into localnet relayer internal account (that pays fees for relayer registration) for each chain."
rly keys restore localnet-1 akasha "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
rly keys restore localnet-2 akasha "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
rly keys restore localnet-3 akasha "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
echo "Mnemonic for all tree chains: ->"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

echo "Testing query localnet-1"
rly q balance localnet-1
echo "Testing query localnet-2"
rly q balance localnet-2
echo "Testing query localnet-3"
rly q balance localnet-3

rly paths generate localnet-1 localnet-2 transfer1 --port=transfer
rly paths generate localnet-2 localnet-3 transfer2 --port=transfer
rly paths generate localnet-3 localnet-1 transfer3 --port=transfer


