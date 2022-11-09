# Siftest User Guide

## About

This can be used after running transactions to verify expected vs real changes.

## Installation

```shell
git clone git@github.com:Sifchain/sifnode.git
cd sifnode
make install
```

## Verify Add Liquidity

1. Execute add liquidity transaction
2. Verify by passing in the following arguments
   1. --height [height of transaction]
   2. --from [address of transactor]
   3. --external-asset [external asset of pool]
   4. --nativeAmount [native amount requested to add]
   5. --externalAmount [external amount requested to add]
   6. --node [node to connect to]
```shell
siftest verify add --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --height=43516 --external-asset=ceth --nativeAmount=96176925423929435353999282 --externalAmount=488436982990 --node tcp://localhost:26657
```

Output
```shell
verifying add...

Wallet native balance before 499999807448801156767565321374116
Wallet external balance before 499999999999999999999022148656694

Wallet native balance after 499999711271875632838129967374834 
Wallet external balance after 499999999999999999998533711673704 

Wallet native diff -96176925523929435353999282 (expected: -96176925423929435353999282 unexpected: -100000000000000000)
Wallet external diff -488436982990 (expected: -488436982990 unexpected: 0)

LP units before 192542049745763466715763665 
LP units after 288716849962488234851636499 
LP units diff 96174800216724768135872834 (expected: 96174800216724768135872834)

Pool share before 1.000000000000000000
Pool share after 1.000000000000000000
```

## Verify Remove Liquidity

1. Execute remove liquidity transaction
2. Verify by passing in the following parameters:
   1. --height [height of transaction]
   2. --from [address of transactor]
   3. --external-asset [external asset of pool]
   4. --units [units requested for removal]
   5. --node [node to connect to]
Command
```shell
siftest verify remove --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --units 1000000000000000000 --height=33068 --external-asset=ceth --node tcp://localhost:26657
```

Output
```shell
verifying removal...

Wallet native balance before 499999807448701559072997868307527
Wallet external balance before 499999999999999999999022148651615

Wallet native balance after 499999807448702459100223367291740 
Wallet external balance after 499999999999999999999022148656694 

Wallet native diff 900027225498984213 (expected: 1000032419169645384 unexpected: -100005193670661171)
Wallet external diff 5079 (expected: 5079 unexpected: 0)

LP units before 192542050745763466715763665 
LP units after 192542049745763466715763665 
LP units diff -1000000000000000000 (expected: -1000000000000000000)

Pool share before 1.000000000000000000
Pool share after 1.000000000000000000

```

## Verify Close Position

1. Execute close margin position
2. Verify by passing in the following params:
   1. --height [height of transaction]
   2. --id [mtp id]
   3. --from [owner of mtp]
   4. --node [node to connect to]

Run command using height of close transaction and MTP id.
```shell
siftest verify close --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --height=72990 --id=4 --node tcp://localhost:26657
```
Output
```shell
verifying close...

MTP collateral 500000000000 (stake)
MTP leverage 2.000000000000000000
MTP liability 500000000000
MTP health 1.988017999020000000
MTP interest paid custody 539
MTP interest paid collateral 550
MTP interest unpaid collateral 0

Wallet collateral balance before: 488999999999999500000000000
Wallet custody balance before: 499999211271873832838129967124878

confirmed MTP does not exist at close height 72990


Pool health before 0.999999999999999000
Pool native custody before 996999999460
Pool external custody before 0
Pool native liabilities before 0
Pool external liabilities before 500000000000
Pool native depth (including liabilities) before 499999999999999003000000496
Pool external depth (including liabilities) before 500000000000001000000000000

Pool health after 0.999999999999999000
Pool native custody after 0
Pool external custody after 0
Pool native liabilities after 0
Pool external liabilities after 0

Return amount: 494008999461
Loss: 0

Wallet collateral balance after: 488999999999999994008999412 (diff: 494008999412)
Wallet custody balance after: 499999211271873732838129967124882 (diff: -99999999999999996)

```