#!/usr/bin/env bash

echo "Creating pools ceth and cdash"
sleep 8
yes Y | sifnodecli tx clp create-pool --from user2 --symbol ceth --nativeAmount 200 --externalAmount 200
sleep 8
yes Y | sifnodecli tx clp create-pool --from user2 --symbol cdash --nativeAmount 100 --externalAmount 100

echo "Query all pools"
sleep 8
sifnodecli query clp pools

echo "Query specific pool"
sleep 8
sifnodecli query clp pool ceth

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp ceth $(sifnodecli keys show user2 -a)

echo "adding more liquidity"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from user2 --symbol ceth --nativeAmount 1 --externalAmount 1

echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from user2 --sentSymbol ceth --receivedSymbol cdash --sentAmount 20


echo "removing Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user2 --symbol ceth --wBasis 5001 --asymmetry -1

echo "removing more Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user2 --symbol ceth --wBasis 5001 --asymmetry -1



echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from user2 --symbol ceth

echo "sifnodecli query clp pools -> should list both pools / Decommission can only be done by admin users"