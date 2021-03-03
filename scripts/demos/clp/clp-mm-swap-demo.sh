#!/usr/bin/env bash

echo "Creating pools ceth and cdash"
sifnodecli tx clp create-pool --from sif --symbol ceth --nativeAmount 20000000000000000000 --externalAmount 20000000000000000000  --yes

sleep 5
sifnodecli tx clp create-pool --from sif --symbol cdash --nativeAmount 20000000000000000000 --externalAmount 20000000000000000000  --yes


sleep 8
echo "Swap Native for Pegged - Sent rowan Get ceth"
sifnodecli tx clp swap --from sif --sentSymbol rowan --receivedSymbol ceth --sentAmount 2000000000000000000 --minReceivingAmount 0 --yes
sleep 8
echo "Swap Pegged for Native - Sent ceth Get rowan"
sifnodecli tx clp swap --from sif --sentSymbol ceth --receivedSymbol rowan --sentAmount 2000000000000000000 --minReceivingAmount 0 --yes
sleep 8
echo "Swap Pegged for Pegged - Sent ceth Get cdash"
sifnodecli tx clp swap --from sif --sentSymbol ceth --receivedSymbol cdash --sentAmount 2000000000000000000 --minReceivingAmount 0 --yes

sifnodecli q clp pools

