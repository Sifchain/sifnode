#!/usr/bin/env bash

echo "Creating pools ceth and cdash"
sifnoded tx clp create-pool --from sif --symbol ceth --nativeAmount 20000000000000000000 --externalAmount 20000000000000000000  --yes

sleep 5
sifnoded tx clp create-pool --from sif --symbol cdash --nativeAmount 20000000000000000000 --externalAmount 20000000000000000000  --yes


sleep 8
echo "Swap Native for Pegged - Sent rowan Get ceth"
sifnoded tx clp swap --from sif --sentSymbol rowan --receivedSymbol ceth --sentAmount 2000000000000000000 --minReceivingAmount 0 --yes
sleep 8
echo "Swap Pegged for Native - Sent ceth Get rowan"
sifnoded tx clp swap --from sif --sentSymbol ceth --receivedSymbol rowan --sentAmount 2000000000000000000 --minReceivingAmount 0 --yes
sleep 8
echo "Swap Pegged for Pegged - Sent ceth Get cdash"
sifnoded tx clp swap --from sif --sentSymbol ceth --receivedSymbol cdash --sentAmount 2000000000000000000 --minReceivingAmount 0 --yes

sifnoded q clp pools

