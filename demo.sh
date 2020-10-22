#!/usr/bin/env bash
echo "Creating a pool"
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ceth --nativeAmount 99 --externalAmount 99
sleep 8
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain DASH --symbol DASH --ticker cdash --nativeAmount 100 --externalAmount 100

echo "Query all pools"
sleep 8
sifnodecli query clp pools

echo "Query specific pool"
sleep 8
sifnodecli query clp pool ceth

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp ceth $(sifnodecli keys show user1 -a)

echo "adding more liquidity"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ceth --nativeAmount 1 --externalAmount 1

echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from user1 --sentSourceChain ETHEREUM --sentSymbol ETH --sentTicker ceth --receivedSourceChain DASH --receivedSymbol DASH --receivedTicker cdash --sentAmount 20


echo "removing Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ceth --wBasis 5001 --asymmetry -1

echo "removing more Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ceth --wBasis 5001 --asymmetry -1



echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from user1 --ticker ceth