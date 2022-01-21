# Sifchain - WASM Tutorial

In this tutorial we will deploy and interact with a wasm smart-contract on a
local sifnoded instance.

Sifchain implements the same wasm functionality as Juno, and this tutorial is 
basically taken from [their docs](https://docs.junonetwork.io/smart-contracts-and-junod-development/tutorial-erc-20).

We will skip the part about writing and compiling the smart-contract. Instead we
will deploy the precompiled file kept in `sc/cw_erc20.wasm` which, as the name 
implies, implements erc20.

## Setup

1. Initialize the local chain: `make init`

2. Start the chain: `make run`

3. Check to see you have two local accounts/keys setup; `sifnoded keys list --keyring-backend=test`

```
- name: akasha
  type: local
  address: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"}'
  mnemonic: ""
- name: mkey
  type: multi
  address: sif1kkdqp4dtqmc7wh59vchqr0zdzk8w2ydukjugkz
  pubkey: '{"@type":"/cosmos.crypto.multisig.LegacyAminoPubKey","threshold":2,"public_keys":[{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AvUEsFHbsr40nTSmWh7CWYRZHGwf4cpRLtJlaRO4VAoq"},{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"}]}'
  mnemonic: ""
- name: sif
  type: local
  address: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AvUEsFHbsr40nTSmWh7CWYRZHGwf4cpRLtJlaRO4VAoq"}'
  mnemonic: ""
```

## Store/Upload the contract to the chain

```
sifnoded tx wasm store docs/tutorials/sc/cw_erc20.wasm \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet  \
-y
```

Output:

```
code: 0
codespace: ""
data: 0A240A1E2F636F736D7761736D2E7761736D2E76312E4D736753746F7265436F646512020801
gas_used: "16622540"
gas_wanted: "1000000000000000000"
height: "35"
info: ""
logs:
- events:
  - attributes:
    - key: action
      value: /cosmwasm.wasm.v1.MsgStoreCode
    - key: module
      value: wasm
    - key: sender
      value: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
    type: message
  - attributes:
    - key: code_id
      value: "1"
    type: store_code
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"/cosmwasm.wasm.v1.MsgStoreCode"},{"key":"module","value":"wasm"},{"key":"sender","value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"}]},{"type":"store_code","attributes":[{"key":"code_id","value":"1"}]}]}]'
timestamp: ""
tx: null
txhash: 24C7E1C514B6280B41C33690208EF873091471FFA3239A87CFD0EAC7C1A3C49D
```

## Initialise

Here we create a Coin called `Poodle Coin` with symbol `POOD` and initial amount
of 1234568000.

We create it from Sif's account (an account that we create with `make init`), so 
the initial balance should go to that account.

```
sifnoded tx wasm instantiate 1 '{"name":"Poodle Coin","symbol":"POOD","decimals":6,"initial_balances":[{"address":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd","amount":"12345678000"}]}' \
--amount 50000rowan \
--label "Poodlecoin erc20" \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y
```

Output;

```
code: 0
codespace: ""
data: 0A6C0A282F636F736D7761736D2E7761736D2E76312E4D7367496E7374616E7469617465436F6E747261637412400A3E7369663134686A32746176713866706573647778786375343472747933686839307668756A7276636D73746C347A723374786D6676773973363263767536
gas_used: "151279"
gas_wanted: "1000000000000000000"
height: "158"
info: ""
logs:
- events:
  - attributes:
    - key: receiver
      value: sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
    - key: amount
      value: 50000rowan
    type: coin_received
  - attributes:
    - key: spender
      value: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
    - key: amount
      value: 50000rowan
    type: coin_spent
  - attributes:
    - key: _contract_address
      value: sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
    - key: code_id
      value: "1"
    type: instantiate
  - attributes:
    - key: action
      value: /cosmwasm.wasm.v1.MsgInstantiateContract
    - key: module
      value: wasm
    - key: sender
      value: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
    type: message
  - attributes:
    - key: recipient
      value: sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
    - key: sender
      value: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
    - key: amount
      value: 50000rowan
    type: transfer
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"coin_received","attributes":[{"key":"receiver","value":"sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6"},{"key":"amount","value":"50000rowan"}]},{"type":"coin_spent","attributes":[{"key":"spender","value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"},{"key":"amount","value":"50000rowan"}]},{"type":"instantiate","attributes":[{"key":"_contract_address","value":"sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6"},{"key":"code_id","value":"1"}]},{"type":"message","attributes":[{"key":"action","value":"/cosmwasm.wasm.v1.MsgInstantiateContract"},{"key":"module","value":"wasm"},{"key":"sender","value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6"},{"key":"sender","value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"},{"key":"amount","value":"50000rowan"}]}]}]'
timestamp: ""
tx: null
txhash: 694D4DBDBE0433806F85AF7EBFF0977A0EB028BEE55F144AE3C411D8F035CBEE
```

We see that the contract address is `sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6`

## Query

### Get contract info:

```
sifnoded query wasm contract sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
```

Output:

```
address: sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
contract_info:
  admin: ""
  code_id: "1"
  created: null
  creator: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
  extension: null
  ibc_port_id: ""
  label: Poodlecoin erc20
```

### Check Sif's balance

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 '{"balance":{"address":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"}}'
```

Output:

```
data:
  balance: "12345678000"
```

## Run Commands

Send 200 POOD from Sif to Akasha:

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 '{"transfer":{"amount":"200","owner":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd","recipient":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}}' --from sif --keyring-backend test --chain-id localnet --broadcast-mode block -y  
```

Output:

```
code: 0
codespace: ""
data: 0A260A242F636F736D7761736D2E7761736D2E76312E4D736745786563757465436F6E7472616374
gas_used: "113938"
gas_wanted: "200000"
height: "355"
info: ""
logs:
- events:
  - attributes:
    - key: _contract_address
      value: sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
    type: execute
  - attributes:
    - key: action
      value: /cosmwasm.wasm.v1.MsgExecuteContract
    - key: module
      value: wasm
    - key: sender
      value: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
    type: message
  - attributes:
    - key: _contract_address
      value: sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
    - key: action
      value: transfer
    - key: sender
      value: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
    - key: recipient
      value: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
    type: wasm
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"execute","attributes":[{"key":"_contract_address","value":"sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6"}]},{"type":"message","attributes":[{"key":"action","value":"/cosmwasm.wasm.v1.MsgExecuteContract"},{"key":"module","value":"wasm"},{"key":"sender","value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"}]},{"type":"wasm","attributes":[{"key":"_contract_address","value":"sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6"},{"key":"action","value":"transfer"},{"key":"sender","value":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"},{"key":"recipient","value":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}]}]}]'
timestamp: ""
tx: null
txhash: 378ABB5D7867EE39CB460FECB24C759ACA063BD9680D50B494A1A0813F86CE64
```

Check Sif's balances again:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 '{"balance":{"address":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"}}'
```

Output:

```
data:
  balance: "12345677800"
```

And Akasha's balance:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 '{"balance":{"address":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"}}'
```

Output:

```
data:
  balance: "200"
```
