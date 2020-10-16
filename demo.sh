#!/usr/bin/env bash
echo "Creating a pool"
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ETH --nativeAmount 1000 --externalAmount 1000
sleep 8
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain TEZOS --symbol XTZ --ticker XTZ --nativeAmount 1000 --externalAmount 1000

echo "Query all pools"
sleep 8
sifnodecli query clp pools

echo "Query specific pool"
sleep 8
sifnodecli query clp pool ETH

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 8
sifnodecli query clp lp ETH $(sifnodecli keys show user1 -a)

echo "adding more liquidity"
sleep 8
yes Y | sifnodecli tx clp add-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ETH --nativeAmount 1000 --externalAmount 1000


echo "removing Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ETH --wBasis 10000 --asymmetry 1
#

#This swap will not work in the future as sent and received asset is the same ,but that validation will be added in a subsequent pr
echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from user1 --sentSourceChain ETHEREUM --sentSymbol ETH --sentTicker ETH --receivedSourceChain TEZOS --receivedSymbol XTZ --receivedTicker XTZ --sentAmount 10

echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from user1 --ticker ETH