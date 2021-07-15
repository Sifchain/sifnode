#!/usr/bin/env bash

seq=77
seq2=21
tx_in_block=2000
for (( i=2; i <= $tx_in_block; ++i ))
do
sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27665 --chain-id=localnet-1 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-1 --yes --broadcast-mode async --sequence $seq --account-number 1 --offline
seq=$((seq+1))
sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-2 --yes --broadcast-mode async --sequence $seq2 --account-number 1 --offline
seq2=$((seq2+1))
#sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27665 --chain-id=localnet-1 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-1 --yes --packet-timeout-timestamp 0

done