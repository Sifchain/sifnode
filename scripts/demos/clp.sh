#!/usr/bin/env bash

echo "Creating pools ceth and cdash"
yes Y | sifnodecli tx clp create-pool --from shadowfiend --symbol ceth --nativeAmount 200 --externalAmount 200
sleep 5
yes Y | sifnodecli tx clp create-pool --from shadowfiend --symbol cdash --nativeAmount 100 --externalAmount 100

echo "Query all pools"
sleep 8
sifnodecli query clp pools

echo "Query specific pool"
sleep 8
sifnodecli query clp pool ceth

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp ceth $(sifnodecli keys show shadowfiend -a)

echo "adding more liquidity"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from shadowfiend --symbol ceth --nativeAmount 1 --externalAmount 1

echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from shadowfiend --sentSymbol ceth --receivedSymbol cdash --sentAmount 20


echo "removing Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from shadowfiend --symbol ceth --wBasis 5001 --asymmetry -1

echo "removing more Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from shadowfiend --symbol ceth --wBasis 5001 --asymmetry -1



echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from shadowfiend --symbol ceth

echo "sifnodecli query clp pools "
echo "Should list only 1 pool , user has been added as admin"