#!/usr/bin/env bash

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli
rm -rf sifnode.log
rm -rf testlog.log

cd "$(dirname "$0")"

./init.sh
sleep 8
sifnoded start >> sifnode.log 2>&1  &
sleep 8

echo "Request coins from faucet"
yes Y | sifnodecli tx faucet request-coins 10000rowan --from sif
sleep 5

echo "Add Coins to faucet"
yes Y | sifnodecli tx faucet add-coins 100rowan --from sif
sleep 5
sifnodecli query faucet balance
