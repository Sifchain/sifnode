#!/usr/bin/env bash

yes Y | sifnodecli tx clp create-pool --from user1 --sourceChain ETHEREUM --symbol ETH --ticker ceth --nativeAmount 200 --externalAmount 200
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


echo "Query all asset for the liquidity provider "
sleep 8
sifnodecli query clp assets $(sifnodecli keys show user1 -a)

