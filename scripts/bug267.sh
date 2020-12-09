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

yes Y | sifnodecli tx clp create-pool --from akasha --symbol catk --nativeAmount 1000 --externalAmount 1000
sleep 8
yes Y | sifnodecli tx clp create-pool --from akasha --symbol cbtk --nativeAmount 1000 --externalAmount 1000


echo "Query specific pool"
sleep 8
sifnodecli query clp pool catk

echo "adding new liquidity provider"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from shadowfiend --symbol catk --nativeAmount 10000 --externalAmount 10000

echo "Query 1st Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp catk $(sifnodecli keys show akasha -a)

echo "Query 2nd Liquidity Provider "
sleep 8
sifnodecli query clp lp catk $(sifnodecli keys show shadowfiend -a)


pkill sifnoded
