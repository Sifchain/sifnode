#!/usr/bin/env bash
echo "Creating a pool"
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain ETH --symbol cLINK --ticker cLINK --nativeAmount 1000 --externalAmount 1

echo "Query all pools"
sleep 8
sifnodecli query clp pools

echo "Query specific pool"
sleep 3
sifnodecli query clp pool cLINK ETH

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 3
sifnodecli query clp lp cLINK $(sifnodecli keys show user1 -a)

echo "Adding more liquidity"
sleep 3
yes Y | sifnodecli tx clp add-liquidity --from user1 --sourceChain ETH --symbol cLINK --ticker cLINK --nativeAmount 1000 --externalAmount 1


echo "Removing Liquidity"
sleep 3
yes Y | sifnodecli tx clp remove-liquidity --from user1 --sourceChain ETH --symbol cLINK --ticker cLINK --wBasis 1000 --asymmetry 1


#This swap will not work in the future as sent and received asset is the same ,but that validation will be added in a subsequent pr
echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from user1 --sentSourceChain ETH --sentSymbol cLINK --sentTicker cLINK --receivedSourceChain ETH --receivedSymbol cLINK --receivedTicker cLINK --sentAmount 1000

echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from user1 --ticker cLINK --sourceChain ETH
