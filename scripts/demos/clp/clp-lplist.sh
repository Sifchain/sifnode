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
yes Y | sifnodecli tx clp add-liquidity --from sif --symbol catk --nativeAmount 5000000000000000000000 --externalAmount 5000000000000000000
sleep 8
echo "Getting from CLI"
sifnodecli query clp lplist catk

echo "Getting from REST"
curl --request GET -sL \
     --url 'http://localhost:1317/clp/getLpList?symbol=catk'\
