#!/usr/bin/env bash


echo "Minting peggyeth ( minted from Peggy) using ethbridge"
## Case 1
## 1. send tx to cosmos after get the lock event in ethereum
sifnodecli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 3 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
$(sifnodecli keys show user2 -a) $(sifnodecli keys show user1 -a --bech val) 10000 lock \
--token-contract-address=0x0000000000000000000000000000000000000000 --ethereum-chain-id=3 --from=user1 --yes

# 2. query the tx
#sifnodecli q tx

# 3. check user2 account balance
sifnodecli q auth account $(sifnodecli keys show user2 -a)

# 4. query the prophecy
sifnodecli query ethbridge prophecy 0x30753E4A8aad7F8597332E813735Def5dD395028 3 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 --ethereum-chain-id=3 --token-contract-address=0x0000000000000000000000000000000000000000

## Case 2
## 1. burn peggyetch for user2
sifnodecli tx ethbridge burn $(sifnodecli keys show user2 -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
10 ceth --ethereum-chain-id=3 --from=user2 --yes

## 2. query the tx
#sifnodecli q tx

## 3. check user2 account balance
sifnodecli q auth account $(sifnodecli keys show user2 -a)

## Case 3
## 1. lock user2 rwn in sifchain
sifnodecli tx ethbridge lock $(sifnodecli keys show user2 -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
10 rwn  --ethereum-chain-id=3 --from=user2 --yes

## 2. query the tx
#sifnodecli q tx

## 3. check user2 account balance
sifnodecli q auth account $(sifnodecli keys show user2 -a)

## Case 4
## 1. send tx to cosmos after erwn burn in ethereum
sifnodecli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 1 rwn 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
$(sifnodecli keys show user2 -a) $(sifnodecli keys show user1 -a --bech val) \
10 burn --ethereum-chain-id=3 --token-contract-address=0x345cA3e014Aaf5dcA488057592ee47305D9B3e10 --from=user1 --yes

## 2. query the tx
#sifnodecli q tx

## 3. check user2 account balance
sifnodecli q auth account $(sifnodecli keys show user2 -a)


echo "Creating pools for peggyeth ( minted from Peggy) and cdash"
sleep 8
yes Y | sifnodecli tx clp create-pool --from user2 --sourceChain ETHEREUM --symbol ETH --ticker ceth --nativeAmount 200 --externalAmount 200
sleep 8
yes Y | sifnodecli tx clp create-pool --from user2 --sourceChain DASH --symbol DASH --ticker cdash --nativeAmount 100 --externalAmount 100

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
yes Y | sifnodecli tx clp add-liquidity --from user2 --sourceChain ETHEREUM --symbol ETH --ticker ceth --nativeAmount 1 --externalAmount 1

echo "swap"
sleep 8
yes Y |  sifnodecli tx clp swap --from user2 --sentSourceChain ETHEREUM --sentSymbol ETH --sentTicker ceth --receivedSourceChain DASH --receivedSymbol DASH --receivedTicker cdash --sentAmount 20


echo "removing Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user2 --sourceChain ETHEREUM --symbol ETH --ticker ceth --wBasis 5001 --asymmetry -1

echo "removing more Liquidity"
sleep 8
yes Y | sifnodecli tx clp remove-liquidity --from user2 --sourceChain ETHEREUM --symbol ETH --ticker ceth --wBasis 5001 --asymmetry -1



echo "decommission pool"
sleep 8
yes Y | sifnodecli tx clp decommission-pool --from user2 --ticker ceth

echo "sifnodecli query clp pools -> should list only 1 pool cdash"