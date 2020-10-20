#!/usr/bin/env bash
echo "Creating a pool"
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain ETHEREUM --symbol ETH --ticker eth --nativeAmount 99 --externalAmount 99
sleep 8
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain DASH --symbol DASH --ticker dash --nativeAmount 100 --externalAmount 100

echo "Query all pools"
sleep 8
sifnodecli query clp pools

echo "Query specific pool"
sleep 8
sifnodecli query clp pool eth

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp eth $(sifnodecli keys show user1 -a)

echo "adding more liquidity"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker eth --nativeAmount 1 --externalAmount 1

echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from user1 --sentSourceChain ETHEREUM --sentSymbol ETH --sentTicker eth --receivedSourceChain DASH --receivedSymbol DASH --receivedTicker dash --sentAmount 20


echo "removing Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker eth --wBasis 5001 --asymmetry -1


echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from user1 --ticker eth