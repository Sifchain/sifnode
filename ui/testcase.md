## Failing test case for frontend which is blocking us being able to peg tokens in the interface

Run all of these in separate terminals

### Terminal 1

Run ethereum

```bash
yarn chain:eth
```

Effectively just runs ganache and does the following:

```bash
yarn ganache-cli -m "$ETHEREUM_ROOT_MNEMONIC" -p 7545 --networkId 5777
```

---

### Terminal 2

Run sifchain

```bash
y chain:sif
```

Effectively does the following:

<details><summary>yarn chain:sif</summary>
<p>

```bash
#!/usr/bin/env bash

. ../credentials.sh


parallelizr() {
  for cmd in "$@"; do {
    echo "Process \"$cmd\" started";
    $cmd & pid=$!
    PID_LIST+=" $pid";
  } done

  trap "kill $PID_LIST" SIGINT

  echo "Parallel processes have started";

  wait $PID_LIST

  echo "All processes have completed";
}

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

sifnoded init test --chain-id=sifchain

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id sifchain
sifnodecli config keyring-backend test

echo "Generating deterministic account - ${SIFUSER1_NAME}"
echo "${SIFUSER1_MNEMONIC}" | sifnodecli keys add ${SIFUSER1_NAME} --recover

echo "Generating deterministic account - ${SIFUSER2_NAME}"
echo "${SIFUSER2_MNEMONIC}" | sifnodecli keys add ${SIFUSER2_NAME} --recover

sifnoded add-genesis-account $(sifnodecli keys show ${SIFUSER1_NAME} -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000stake
sifnoded add-genesis-account $(sifnodecli keys show ${SIFUSER2_NAME} -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000stake

sifnoded gentx --name ${SIFUSER1_NAME} --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis

echo "Starting test chain"

parallelizr "sifnoded start" "sifnodecli rest-server  --unsafe-cors --trace"

```

</p>
</details>

---

### Terminal 3

create liquidity pools and deploy peggy contracts then run ebrelayer

```bash
y chain:migrate && y chain:peggy
```

<details><summary>yarn chain:migrate</summary><p>

```bash
#!/bin/bash

# runs migrate in the smart-contracts folder against our local version of ganache
. ../credentials.sh

cd ../../../smart-contracts

yarn && yarn migrate
```

```bash
# Sets up a couple of liquidity pools to test against
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol catk \
 --nativeAmount 1000000 \
 --externalAmount 1000000 \
 --yes

sifnodecli tx clp create-pool \
 --from akasha \
 --symbol cbtk \
 --nativeAmount 1000000 \
 --externalAmount 1000000 \
 --yes
```

```bash
# deploys a few test erc-20 tokens for us to test swapping with in the interface

cd ui/chains/ethereum
yarn migrate
```

</p>
</details>

<details><summary>yarn chain:peggy</summary>
<p>

```bash
#!/bin/bash

# This script should be run with a CWD that is the local folder
. ../credentials.sh
. ../../../smart-contracts/.env


BRIDGE_TOKEN_ADDRESS="0x82D50AD3C1091866E258Fd0f1a7cC9674609D254"

echo "ETHEREUM_PRIVATE_KEY=$ETHEREUM_PRIVATE_KEY"
echo "BRIDGE_TOKEN_ADDRESS=$BRIDGE_TOKEN_ADDRESS"
echo "SIFUSER1_NAME=$SIFUSER1_NAME"

ETHEREUM_PRIVATE_KEY=$ETHEREUM_PRIVATE_KEY ebrelayer init \
  tcp://localhost:26657 \
  ws://localhost:7545/ \
  "$BRIDGE_TOKEN_ADDRESS" \
  "$SIFUSER1_NAME" \
  "$SIFUSER1_MNEMONIC" \
  --chain-id=sifchain \
  --gas 300000 \
  --gas-adjustment 1.5
```

</p>
</details>

---

### Terminal 4

#### 1) Check ethereum balance

```bash
yarn peggy:getTokenBalance 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
```

returns:

```
Eth balance for 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 is 99.70074700843 Eth (99700747008430000000 Wei)
```

---

#### 2) Check ceth balance

```bash
sifnodecli query account sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
```

Returns

```json
{
  "type": "cosmos-sdk/Account",
  "value": {
    "address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    "coins": [
      {
        "denom": "catk",
        "amount": "999000000"
      },
      {
        "denom": "cbtk",
        "amount": "999000000"
      },
      {
        "denom": "ceth",
        "amount": "1000000000"
      },
      {
        "denom": "rowan",
        "amount": "998000000"
      },
      {
        "denom": "stake",
        "amount": "100000000"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"
    },
    "account_number": "4",
    "sequence": "2"
  }
}
```

Note `ceth` amount is `1000000000`

---

#### 3) Execute the lock

```bash
yarn peggy:lock sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 0x0000000000000000000000000000000000000000 2000000000000000000
```

Try to lock 2.0 eth in bridgebank

```

Expected usage:
truffle exec scripts/sendLockTx.js --network ropsten sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace eth 100

Connecting to contract....
Connected to contract, sending lock...
Sent lock...
{
  to: '0x736966316c376879706d716b327963333334766336766d64777a70357364656679676a32616439337035',
  from: '0x627306090abaB3A6e1400e9345bC60c78a8BEf57',
  symbol: 'ETH',
  token: '0x0000000000000000000000000000000000000000',
  value: 2000000000000000000,
  nonce: 1
}
```

Expected usage message should be expected according to script.

This adds the following to the output of ebrelayer:

```
I[2021-01-03|12:01:19.336]
Chain ID: 5777
Bridge contract address: 0xf204a4Ef082f5c04bB89F7D5E6568B796096735a
Token symbol: ETH
Token contract address: 0x0000000000000000000000000000000000000000
Sender: 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
Recipient: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
Value: 2000000000000000000
Nonce: 1
Claim type: lock
I[2021-01-03|12:01:19.336] Add event into buffer
```

---

#### 4) run advance to advance the ethereum blockchain

```
y advance 100
```

```
Advanced 100 blocks
current block number is 126
{"nBlocks":"100","currentBlockNumber":126}
```

---

#### 5) Query account but `ceth` balance has not updated

After waiting for some time...

```
sifnodecli query account sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
```

```
{
  "type": "cosmos-sdk/Account",
  "value": {
    "address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    "coins": [
      {
        "denom": "catk",
        "amount": "999000000"
      },
      {
        "denom": "cbtk",
        "amount": "999000000"
      },
      {
        "denom": "ceth",
        "amount": "1000000000"
      },
      {
        "denom": "rowan",
        "amount": "998000000"
      },
      {
        "denom": "stake",
        "amount": "100000000"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"
    },
    "account_number": "4",
    "sequence": "2"
  }
}
```
