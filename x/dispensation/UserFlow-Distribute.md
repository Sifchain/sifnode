# USER FLOW 
## Context 
- Amara wants to distribute funds to some recipients . She already has a list of recipients formatted as below .
```json
{
  "Output": [
    {
      "address": "sif1acdh3rca2elta9jdg5a6mjsw2cv3map6d8uc0x",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif1g0ecn4l05rdtzd8vcxpnt8283wxrnx4p3g7s3e",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif12xyxcdvxg8xqydu2lejadvmycuryuxxckg84p3",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif1u0yj66x98sshaddfww5dtjx34apjsqvqkzxnjy",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif1egzcve0udyxnakeq9vw9ynzle2qj3awf0zlny2",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif1qx72w5t2g2gv7htmt57kff0j6rrv4vxsmz2g8p",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif1cvp23q8hkx0mqy923s46q5dwv0c7us8c0ntda8",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif104gd36rr8t3mkxtspv2hl4e3w365hkl46m9qj9",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    },
    {
      "address": "sif1ka2euq8p6ymadgz9g9wcc34p84xs4ndp6gkwnr",
      "coins": [
        {
          "denom": "rowan",
          "amount": "10000000000000000000"
        }
      ]
    }
  ]
}
```

## Steps to follow 
### Setup Multi-Sig Key

- Check local wallet to verify keys 
```shell
sifnoded keys list --keyring-backend file
```
Sample output ( address will be different )
```json
[
  {
    "name": "amara",
    "type": "local",
    "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
    "pubkey": "sifpub1addwnpepqt6sfvz3mwetudyaxjn958kztxz9j8rvrlsu55fw6fjkjyac2s9z5sc8npe"
  }
]
```
### Create Dispensation transactions
Amara wants to create an Airdrop which will be executed by Zane. The create transaction would just crete the drops . Zane would need to run the "run" transaction multiple times to distribute the rewards.

Amara can assign Zane `sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5` to be the authorized runner during the create transactions
```shell
sifnoded tx dispensation create Airdrop output.json sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan
```
Sample output
```json
{
  "height": "0",
  "txhash": "A9D019E1080ECD6A012B20B3058534AC6643BD17634F181FBE7F8F5C43B94D8E",
  "raw_log": "[]"
}
```

### Run Dispensation transactions
Zane can now run the dispensation transactions (The configuration is set to distribute 10 rewards every block)
Distribution name is automatically assigned to Height_DistributerAddress `2_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd`
```shell
sifnoded tx dispensation run 2_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd Airdrop --from sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan
```
Sample output
```json
{
  "height": "0",
  "txhash": "A9D019E1080ECD6A012B20B3058534AC6643BD17634F181FBE7F8F5C43B94D8E",
  "raw_log": "[]"
}
```