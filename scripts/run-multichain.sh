#!/usr/bin/env bash

killall sifnoded 

#sifnodecli rest-server &
#sifnoded start

echo "starting sifnode servers"
sleep 1
sifnoded start --home ~/.sifnode-1 --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9090 --grpc-web.address 0.0.0.0:9093 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:27665 >> abci_1.log 2>&1  &
sleep 1
sifnoded start --home ~/.sifnode-2 --p2p.laddr 0.0.0.0:27656  --grpc.address 0.0.0.0:9091 --grpc-web.address 0.0.0.0:9094 --address tcp://0.0.0.0:27660 --rpc.laddr tcp://127.0.0.1:27666 >> abci_2.log 2>&1  &
sleep 1
sifnoded start --home ~/.sifnode-3 --p2p.laddr 0.0.0.0:27657  --grpc.address 0.0.0.0:9092 --grpc-web.address 0.0.0.0:9095 --address tcp://0.0.0.0:27661 --rpc.laddr tcp://127.0.0.1:27667 >> abci_3.log 2>&1  &
sleep 1

echo "starting IBC connections"
#rm -rf ~/.ibc-setup/last-queried-heights.json

#Reset connections
#rm -rf ~/.ibc-12/last-queried-heights.json
#rm -rf ~/.ibc-23/last-queried-heights.json
#rm -rf ~/.ibc-31/last-queried-heights.json
#rm -rf ~/.ibc-12/app.yaml
#rm -rf ~/.ibc-23/app.yaml
#rm -rf ~/.ibc-31/app.yaml
printf "src: localnet-3\ndest: localnet-1\n" > ~/.ibc-31/app.yaml
printf "src: localnet-1\ndest: localnet-2\n" > ~/.ibc-12/app.yaml
printf "src: localnet-2\ndest: localnet-3\n" > ~/.ibc-23/app.yaml

sleep 1
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | ibc-relayer start -i -v --poll 10 --home ~/.ibc-12 >> ibc_12.log &
sleep 1
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | ibc-relayer start -i -v --poll 10 --home ~/.ibc-23 >> ibc_23.log &
sleep 1
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | ibc-relayer start -i -v --poll 10 --home ~/.ibc-31 >> ibc_31.log &

#Created channel:
#  localnet-1: transfer/channel-0 (connection-0)
#  localnet-2: transfer/channel-0 (connection-0)

#sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
echo "Checking sifnode-1 balances"
sifnoded q bank balances sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --node tcp://127.0.0.1:27665
echo "Checking sifnode-2 balances"
sifnoded q bank balances sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --node tcp://127.0.0.1:27666
echo "Checking sifnode-3 balances"
sifnoded q bank balances sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --node tcp://127.0.0.1:27667

#
#sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27665 --chain-id=localnet-1 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-1
#sifnoded tx ibc-transfer transfer transfer channel-2 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 50ibc/E0263CEED41F926DCE9A805F0358074873E478B515A94DF202E6B69E29DA6178 --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-2
#sifnoded tx ibc-transfer transfer transfer channel-0 sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd 50ibc/4C2B3D3B398FC7B8FFA3A96314006FF0B38E3BFC4CE90D8EE46E9EB6768A482D --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=sif --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-2
#sifnoded tx ibc-transfer transfer transfer channel-1 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 50ibc/5C3977A32007D22B1845B57076D0E27C3159C3067B11B9CEF6FA551D71DAEDD6 --node tcp://127.0.0.1:27667 --chain-id=localnet-3 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-3
