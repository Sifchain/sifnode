#!/usr/bin/env bash


yes Y | sifnodecli tx clp create-pool --from akasha --sourceChain ETH --symbol ETH --ticker catk --nativeAmount 200 --externalAmount 200
sleep 8
yes Y | sifnodecli tx clp create-pool --from akasha --sourceChain ETH --symbol ETH --ticker cbtk --nativeAmount 200 --externalAmount 200


echo "Query specific pool"
sleep 8
sifnodecli query clp pool catk

echo "Query 1st Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp catk $(sifnodecli keys show akasha -a)

echo "adding more liquidity"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from shadowfiend --sourceChain ETH --symbol ETH --ticker catk --nativeAmount 100 --externalAmount 100

echo "Query 2nd Liquidity Provider "
sleep 8
sifnodecli query clp lp catk $(sifnodecli keys show shadowfiend -a)

