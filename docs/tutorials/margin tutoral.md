# Sifchain - Margin Trading Basics Tutorial

#### demo video

TBD

#### Previous tutorial

* Clp Basics: https://github.com/Sifchain/sifnode/blob/develop/docs/tutorials/clp%20tutorial.md

#### Dependencies:

    0. `git clone git@github.com:Sifchain/sifnode.git`

#### What is Margin Trading or trading with margin or leverage

Margin trading refers to the use of borrowed funds from continuous liquidity pools (CLP) providers 
to trade a financial asset. The margin trader can bet against an asset to go up (long) or down (short) and 
relies upon the collateral amount to form the loan.
At the same time providing a continous funding interest rate to the CLP providers. 

#### Setup

1. Initialize the local chain run; `./scripts/init.sh`

2. Start the chain; `./scripts/run.sh`

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

4. Check your seed account balance/s;
   `sifnoded q bank balances $(sifnoded keys show sif -a --keyring-backend=test)`
   `sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test)`

#### Create and query pools

note: 
* the minimum threshold for native amount is 10^18 rowan.
* the minimum transaction fee for these operations is 10^17 rowan.

1. Create the first pool for ceth; 
`sifnoded tx clp create-pool --from sif --keyring-backend test --symbol ceth --nativeAmount 2000000000000000000 --externalAmount 2000000000000000000 --fees 100000000000000000rowan --chain-id localnet -y`

2. Create another pool for cdash with a different account; 
`sifnoded tx clp create-pool --from akasha --keyring-backend test --symbol cdash --nativeAmount 3000000000000000000 --externalAmount 3000000000000000000 --fees 100000000000000000rowan --chain-id localnet -y`

3. Query all clp pools; `sifnoded q clp pools`

#### Enable margin on pools

The set of pools that have margin enabled is managed through governance.

The param change proposal takes the format:

```
{
"title": "Margin Pools Param Change",
"description": "Update enabled margin pools",
"changes": [
{
"subspace": "margin",
"key": "Pools",
"value": ["ceth","cusdt"]
}
],
"deposit": "10000000stake"
}
```

To submit a param change proposal:

` sifnoded tx gov submit-proposal param-change proposal.json --from sif --keyring-backend test --chain-id localnet`

To vote on proposal

`sifnoded tx gov vote 1 yes --from sif --chain-id localnet --keyring-backend test`

#### Create and query margin trading positions (MTP)

1. Create margin long position against ceth;
`sifnoded tx margin open-long --from sif --keyring-backend test --borrow_asset ceth --collateral_asset rowan --collateral_amount 1000 --chain-id localnet`

Result:
```
code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: 372EDDE367EE22B3E0D2034F6429BDAB082D756D5223848F2F3A722ADE808615
```

2. Add up to an existing margin position by creating a second margin position for the same asset ceth;
`sifnoded tx margin open-long --from sif --keyring-backend test --borrow_asset ceth --collateral_asset rowan --collateral_amount 500 --chain-id localnet`

Result:
```
code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: 372EDDE367EE22B3E0D2034F6429BDAB082D756D5223848F2F3A722ADE808615
```

3. Query all the existing margin positions (same asset ceth);
`sifnoded q margin positions-for-address $(sifnoded keys show sif -a --keyring-backend=test)`

Result:
```json
{
  "mtps": [
    {
      "id": "xxxx",
      "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
      "collateral_asset": "rowan",
      "collateral_amount": "1000",
      "liabilities_p": "1000",
      "liabilities_i": "0",
      "custody_asset": "ceth",
      "custody_amount": "4000",
      "leverage": "1",
      "mtp_health": "0.100000000000000000"
    },
    {
      "id": "yyyy",
      "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
      "collateral_asset": "rowan",
      "collateral_amount": "500",
      "liabilities_p": "500",
      "liabilities_i": "0",
      "custody_asset": "ceth",
      "custody_amount": "2000",
      "leverage": "1",
      "mtp_health": "0.100000000000000000"
    }
  ]
}
```

#### Reduce size and close existing margin positions

1. Reduce the size of an existing margin position for ceth by closing one of the existing MTPs;
`sifnoded tx margin close-long --from sif --keyring-backend test --id yyyy --borrow_asset ceth --collateral_asset rowan --chain-id localnet`

Result:
```
code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: 372EDDE367EE22B3E0D2034F6429BDAB082D756D5223848F2F3A722ADE808615
```

2. Query remaining margin positions for ceth;
`sifnoded q margin positions-for-address $(sifnoded keys show sif -a --keyring-backend=test)`

Result:
```json
{
  "mtps": [
    {
      "id": "xxxx",
      "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
      "collateral_asset": "rowan",
      "collateral_amount": "1000",
      "liabilities_p": "1000",
      "liabilities_i": "0",
      "custody_asset": "ceth",
      "custody_amount": "4000",
      "leverage": "1",
      "mtp_health": "0.100000000000000000"
    }
  ]
}
```

3. Close the margin long position entirely for ceth;
`sifnoded tx margin close-long --from sif --keyring-backend test --id xxxx --borrow_asset ceth --collateral_asset rowan --chain-id localnet`

Result:
```
code: 0
codespace: ""
data: ""
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: '[]'
timestamp: ""
tx: null
txhash: 372EDDE367EE22B3E0D2034F6429BDAB082D756D5223848F2F3A722ADE808615
```

4. Query existing margin positions (none);
`sifnoded q margin positions-for-address $(sifnoded keys show sif -a --keyring-backend=test)`

Result:
```json
{
  "mtps": []
}
```
