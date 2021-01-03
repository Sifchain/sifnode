## Run chains

Run all of these in separate terminals

Terminal 1

Run ethereum

```bash
y chain:eth
```

Terminal 2

Run sifchain

```bash
y chain:eth
```

Terminal 3

create liquidity pools and deploy peggy contracts then run ebrelayer

```bash
y chain:migrate && y chain:peggy
```

## Test case

Check ethereum balance

```bash
yarn peggy:getTokenBalance 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
```

returns:

```
Eth balance for 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 is 99.70074700843 Eth (99700747008430000000 Wei)
```

---

Check ceth balance

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

run advance to advance the ethereum blockchain

```
y advance 100
```

```
Advanced 100 blocks
current block number is 126
{"nBlocks":"100","currentBlockNumber":126}
```

---

Query account but `ceth` balance has not updated

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
