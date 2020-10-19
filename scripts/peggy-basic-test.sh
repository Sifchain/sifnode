#!/usr/bin/env bash

## Case 1
## 1. send tx to cosmos after get the lock event in ethereum
sifnodecli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 3 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
$(sifnodecli keys show user2 -a) $(sifnodecli keys show user1 -a --bech val) 5 lock \
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
1 peggyeth --ethereum-chain-id=3 --from=user2 --yes

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
## 1. send tx to cosmos after peggyrwn burn in ethereum
sifnodecli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 1 rwn 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
$(sifnodecli keys show user2 -a) $(sifnodecli keys show user1 -a --bech val) \
1 burn --ethereum-chain-id=3 --token-contract-address=0x345cA3e014Aaf5dcA488057592ee47305D9B3e10 --from=user1 --yes

## 2. query the tx
#sifnodecli q tx

## 3. check user2 account balance
sifnodecli q auth account $(sifnodecli keys show user2 -a)