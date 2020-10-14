#!/usr/bin/env bash
echo "Creating a pool"
yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain ROWAN --symbol EOSROWAN --ticker EOS --nativeAmount 1000 --externalAmount 1

echo "Query all pools"
sleep 3
sifnodecli query clp pools

echo "Query specific pool"
sleep 3
sifnodecli query clp pool EOS ROWAN

echo "Query Liquidity Provider / Pool creator is the first lp for the pool"
sleep 3
sifnodecli query clp lp EOS $(sifnodecli keys show user1 -a)


echo "adding more liquidity"
sleep 3
yes Y | sifnodecli tx clp add-liquidity --from user1 --sourceChain ROWAN --symbol EOSROWAN --ticker EOS --nativeAmount 1000 --externalAmount 1


echo "removing Liquidity"
sleep 3
yes Y | sifnodecli tx clp remove-liquidity --from user1 --sourceChain ROWAN --symbol EOSROWAN --ticker EOS --wBasis 1000 --asymmetry 1


#This swap will not work in the future as sent and received asset is the same ,but that validation will be added in a subsequent pr
echo "swap"
sleep 3
yes Y |  sifnodecli tx clp swap --from user1 --sentSourceChain ROWAN --sentSymbol EOSROWAN --sentTicker EOS --receivedSourceChain ROWAN --receivedSymbol EOSROWAN --receivedTicker EOS --sentAmount 1000

echo "decommission pool"
sleep 3
yes Y | sifnodecli tx clp decommission-pool --from user1 --ticker EOS --sourceChain ROWAN