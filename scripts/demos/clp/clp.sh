#!/usr/bin/env bash

echo "Creating pools ceth and cdash"
sifnoded tx clp create-pool --from sif --symbol ceth --nativeAmount 1000000000000000000 --externalAmount 1000000000000000000  --yes

sleep 5
sifnoded tx clp create-pool --from sif --symbol cdash --nativeAmount 1000000000000000000 --externalAmount 1000000000000000000  --yes

echo "Query all pools"
sleep 8
sifnoded query clp pools

echo "Query specific pool"
sleep 8
sifnoded query clp pool ceth

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnoded query clp lp ceth $(sifnoded keys show sif -a)

echo "adding more liquidity"
sleep 8
sifnoded tx clp add-liquidity --from sif --symbol ceth --nativeAmount 1 --externalAmount 1 --yes

echo "swap"
sleep 8
sifnoded tx clp swap --from sif --sentSymbol ceth --receivedSymbol cdash --sentAmount 1000000000000000000 --minReceivingAmount 0 --yes

echo "removing Liquidity"
sleep 8
sifnoded tx clp remove-liquidity --from sif --symbol ceth --wBasis 5001 --asymmetry -1 --yes

echo "removing more Liquidity"
sleep 8
sifnoded tx clp remove-liquidity --from sif --symbol ceth --wBasis 5001 --asymmetry -1 --yes

echo "decommission pool"
sleep 8
sifnoded tx clp decommission-pool --from sif --symbol ceth --yes

echo "sifnoded query clp pools "
echo "Should list only 1 pool , user has been added as admin"



