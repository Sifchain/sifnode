# USER FLOW 
## Context 
- Two users amara and zane want to do an airdrop to distribute rowan to some recipients. The already have a list of recipients formatted as below .
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
 - Amara and Zane also have their mnemonic keys.

## Steps to follow 
### Setup Multi-Sig Key
- They first need add their mnemonic keys to a local wallet .They will generate a multi-sig key using their individual keys. This multi-sig key generation is a one time process only and needs to be done with consent from both amara and zane. For this generation to work both their keys need to be present in the same local wallet.
```shell
#Amara adds her key .
sifnodecli keys add amara -i
##will be prompted to create a bip39 passphrase and a keyring passphrase
# Note the keyring passphrase will be used by zane as well .

# Zane adds his key to the wallet
sifnodecli keys add zane -i
```
- Check local wallet to verify keys 
```shell
sifnodecli keys list --keyring-backend file
```
Sample output ( address will be different )
```json
[
  {
    "name": "amara",
    "type": "local",
    "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
    "pubkey": "sifpub1addwnpepqt6sfvz3mwetudyaxjn958kztxz9j8rvrlsu55fw6fjkjyac2s9z5sc8npe"
  },
  {
    "name": "zane",
    "type": "local",
    "address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    "pubkey": "sifpub1addwnpepqdycrc8usnjh0yk7cd532ushualgsderdqj8jr9m2rzy8stqrlpj5vymlww"
  }
]
```
- Amara and Zane then create a multi-sig key 
```shell
# multi-sig-threshold refers to the min signatures required for this multi-sig key to work
# In this case 2 means that both amara and zane will have to sign for this to work 
sifnodecli keys add mkey --multisig amara,zane --multisig-threshold 2  
```
- After this step ideally both amara and zane can take a copy of the `mkey` which they can use to create transactions ,and send it over to the other person to sign ( More details on this later). Amara and Zane can then leave ,they would be able to sign and broadcast tx from their own locations.
### Create offline tx
- Zane and amara needs to decide how much of the total funding they would provide individually . Based on the output list above we need a total of 10rowan . They decide to fund 5rowan each.Based on this they will create an input file which looks like this.
```json
{
 "Input": [
  {
   "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
   "coins": [
    {
     "denom": "rowan",
     "amount": "5000000000000000000"
    }
   ]
  },
  {
   "address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
   "coins": [
    {
     "denom": "rowan",
     "amount": "5000000000000000000"
    }
   ]
  }
 ]
}
```
- Zane decides to create the offline tx . He has the mkey in his local wallet. He runs the following command.Note the tx creator will have to pay for the gas fee
```shell
sifnodecli tx dispensation create mkey airdrop-1 input.json output.json --gas 200064128 --gas-prices 1.0rowan --generate-only --from zane --keyring-backend file --node tcp://rpc-mainnet.sifchain.finance:80 --chain-id sifchain-mainnet >> offlinetx.json
```

### Zane Signs the transaction 
```shell
sifnodecli tx sign --multisig sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --from zane offlinetx.json --keyring-backend file --node tcp://rpc-mainnet.sifchain.finance:80 --chain-id sifchain-mainnet >> sig-zane.json
```
Sample sig-zane.json
```json
{
  "pub_key": {
    "type": "tendermint/PubKeySecp256k1",
    "value": "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq"
  },
  "signature": "3ZvEoBPCwFqroslYxFXOY+rMcbpnVpFidwVPNNsCEx4GsdZ3RCgsaSXnrbbVN3BHG94q6goOpPIkk7EPtWmK4g=="
}

```

### Amara Signs the transaction
- Zane then sends the two json files offlinetx.json and sig-zane.json to amara
- Amara creates her signature 
```shell
sifnodecli tx sign --multisig sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --from amara offlinetx.json --keyring-backend file --node tcp://rpc-mainnet.sifchain.finance:80 --chain-id sifchain-mainnet >> sig-amara.json
```
Sample sig-amara.json
```json
{
  "pub_key": {
    "type": "tendermint/PubKeySecp256k1",
    "value": "AvUEsFHbsr40nTSmWh7CWYRZHGwf4cpRLtJlaRO4VAoq"
  },
  "signature": "2l3nNYfAXMPJedsowdXUS6ssLN57vk73Q4ahrA8saLI51Ieo7xHQbkLCMyLeuWejkK3EYWzOeaYu1qNhOrNWFw=="
}
```

### Create multi-sig 
- Either amara or zane can do this step . They just need the two signature files and multi-sig key
```shell
sifnodecli tx multisign offlinetx.json mkey sig1.json sig2.json --keyring-backend file --node tcp://rpc-mainnet.sifchain.finance:80 --chain-id sifchain-mainnet >> signedtx.json
```

### Broadcast transaction to network
```shell
sifnodecli tx broadcast signedtx.json --node tcp://rpc-mainnet.sifchain.finance:80 --chain-id sifchain-mainnet 
```
Sample output
```json
{
  "height": "0",
  "txhash": "A9D019E1080ECD6A012B20B3058534AC6643BD17634F181FBE7F8F5C43B94D8E",
  "raw_log": "[]"
}

```