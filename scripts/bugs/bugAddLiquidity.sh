#!/usr/bin/env bash

#Sample script to test the issue mentioned in
#https://app.asana.com/0/1199697235740010/1199903639901927/f

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli
rm -rf sifnode.log
rm -rf testlog.log

cd "$(dirname "$0")"

./init.sh
sleep 5
sifnoded start >> sifnode.log 2>&1  &
sleep 5

yes Y | sifnodecli tx clp create-pool --from sif --symbol cacoin --nativeAmount 144219482657918838950052 --externalAmount 95306982476314799709783920 --fees 1300000rowan
sleep 5

sifnodecli q clp pools

echo "adding new liquidity provider"
sleep 5
yes Y | sifnodecli tx clp add-liquidity --from akasha --symbol cacoin --nativeAmount 1000000000000000000000 --externalAmount 999990445060966787235986082005 --fees 1300000rowan

