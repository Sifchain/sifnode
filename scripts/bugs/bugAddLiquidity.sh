#!/usr/bin/env bash

#Sample script to test the issue mentioned in
#https://app.asana.com/0/1199697235740010/1199903639901927/f

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli
rm -rf sifnode.log
rm -rf testlog.log

cd "$(dirname "$0")"

./init.sh
sleep 8
sifnoded start >> sifnode.log 2>&1  &
sleep 8

yes Y | sifnodecli tx clp create-pool --from sif --symbol cacoin --nativeAmount 1000 --externalAmount 1000 --fees 13000rowan
sleep 8


echo "adding new liquidity provider"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from akasha --symbol cacoin --nativeAmount 1000000000000000000000 --externalAmount 999990445060966787230000000000 --fees 13000rowan

