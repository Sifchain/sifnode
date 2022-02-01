# Sifchain - Wasm (CW1 Subkeys example)

[CW1](https://docs.cosmwasm.com/cw-plus/0.9.0/cw1/intro/) is a specification for
proxy contracts that can hold assets or perform actions on behalf of other 
accounts. It is a fundamental building block for things like "1 of N" multisig
or conditional-approval where sub-accounts have a right to spend a limited 
amount from the contract's account. 

This subkey example implements a system where an admin creates and initializes a
contract with an amount of native token (rowan). The admin can specifiy a list 
of accounts with an "allowance" to spend from the contract account. When an 
authorized account sends rowan through this contract, their own rowan balance is
not affected, instead the funds are taken from the contract's account and some
internal accounting is performed to decrease the calling account's remaining 
allowance.

## Setup

1. Initialize the local chain: `make init`

2. Start the chain: `make run`

## Store and Initialize

Store the contract:

```
sifnoded tx wasm store docs/tutorials/sc/cw1_subkeys.wasm \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet  \
-y
```

Initialize the contract from Sif's account, with 50000rowan. The rowan will be
transferred into the smart-contract account, and keys with an allowance will be
able to spend rowan from that account.

```
sifnoded tx wasm instantiate 1 '{"admins":["sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"],"mutable":false}' \
--amount 5000000rowan \
--label "CW1 Subkey" \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y
```

Check contract balance:

```
sifnoded q bank balances sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
```

Output:

```
balances:
- amount: "5000000"
  denom: rowan
pagination:
  next_key: null
  total: "0"
```

## Query

Check Akasha's allowance:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6  \
 '{"allowance":{"spender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}}' \
 --chain-id localnet
```

Check all allowances:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
 '{"all_allowances":{}}' \
 --chain-id localnet
```

Check which keys are admins:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
 '{"admin_list":{}}' \
 --chain-id localnet
```

## Running Commands

As the admin key (Sif), set an allowance for a key (Akasha)
As the key with an allowance, send tokens from key (Akasha) to key (C)
See tokens arrive at key (C)
See allowance decrease for key (Akasha)

### 1) Add allowance for Akasha

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
  '{"increase_allowance":{"spender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5","amount":{"denom":"rowan","amount":"2000000"}}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y 
```

Check Akasha's allowance again:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6  \
 '{"allowance":{"spender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}}' \
 --chain-id localnet
```

Output:

```
data:
  balance:
  - amount: "2000000"
    denom: rowan
  expires:
    never: {}
```

### 2) Send tokens from Akasha to Sif

First check Sif's balance before the operation:

```
sifnoded q bank balances $(sifnoded keys show -a sif --keyring-backend test) | grep -B1 rowan
```

Output:

```
- amount: "499999999999999999950000"
  denom: rowan
```

Send the tokens via the cw1-subkey proxy from Akasha to Sif:

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
  '{"execute":{"msgs":[{"bank":{"send":{"to_address":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd","amount":[{"denom":"rowan","amount":"999"}]}}}]}}' \
  --from akasha \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```

### 3) See balance increase for Sif

```
sifnoded q bank balances $(sifnoded keys show -a sif --keyring-backend test) | grep -B1 rowan
```

Output:

```
- amount: "499999999999999999950999"
  denom: rowan

```

### 4) See allowance decrease for Akasha

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6  \
 '{"allowance":{"spender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}}' \
 --chain-id localnet
```

Output:

```
data:
  balance:
  - amount: "1999001"
    denom: rowan
  expires:
    never: {}
```

### 5) See Rowan balance decrease for smart-contract


```
sifnoded q bank balances sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
```

Output:

```
balances:
- amount: "49500"
  denom: rowan
pagination:
  next_key: null
  total: "0"
```
