#!/usr/bin/env bash

pkill sifnoded
sleep 5
sifnoded export --height -1 > exported_state.json
sleep 1
sifnoded migrate v0.38 exported_state.json --chain-id new-chain > new-genesis.json  2>&1
sleep 1
sifnoded unsafe-reset-all
sleep 1
cp new-genesis.json ~/.sifnoded/config/genesis.json
sleep 2
sifnoded start