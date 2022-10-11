# Siftest User Guide

## Verify Add Liquidity

```shell
siftest verify add --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --height=43516 --external-asset=ceth --nativeAmount=96176925423929435353999282 --externalAmount=488436982990
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

## Verify Close Position
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