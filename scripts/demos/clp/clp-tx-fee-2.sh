#!/usr/bin/env bash

echo "Creating pools ceth and cdash"
sifnodecli tx clp create-pool --from sif --symbol ceth --nativeAmount 100000000000000000000 --externalAmount 100000000000000000000  --yes --fees 20000000rowan

sleep 5
sifnodecli tx clp create-pool --from sif --symbol cdash --nativeAmount 100000000000000000000 --externalAmount 100000000000000000000  --yes --fees 20000000rowan

echo "swap"
sleep 8
sifnodecli tx clp swap --from akasha --sentSymbol ceth --receivedSymbol cdash --sentAmount 1000000000000000000 --minReceivingAmount 0 --yes --fees 20000000rowan




