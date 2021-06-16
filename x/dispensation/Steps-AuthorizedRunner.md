#Steps to follow by Authorized Runner

##Requirements
1. A dispensation should already be created : A dispensation should have been created, and the authorized runner address assigned to it.
2. Enough funds to pay for gas fees : The authorized runner would need to pay for the gas fee everytime they run a new transaction.

##Steps
###Run Dispensation
Note The run dispensation only executes 10 transfers at a time. 
If the dispensation contains 3000 records, the authorized runner would need to submit a run transaction `3001/10 = 301 times`

The CLI command 
```shell
sifnodecli tx dispensation run [Distribution_Name] [Distribution Type Airdrop/ValidatorSubsidy/LiquidityMining] --from [Authorised runner neeeds to sign the transaction] --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan --node tcp://rpc.sifchain.finance:80 --chain-id sifchain
```
sample CLI command
```shell
sifnodecli tx dispensation run 2_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd LiquidityMining --from sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan --node tcp://rpc.sifchain.finance:80 --chain-id sifchain
```
sample output
```json
{
  "height": "0",
  "txhash": "A9D019E1080ECD6A012B20B3058534AC6643BD17634F181FBE7F8F5C43B94D8E",
  "raw_log": "[]"
}
```
The Tx hash can then be used to query the blockchain to get the related events
```shell
sifnodecli q tx A9D019E1080ECD6A012B20B3058534AC6643BD17634F181FBE7F8F5C43B94D8E --node tcp://rpc.sifchain.finance:80 --chain-id sifchain
```
The relevant event 
```json
{
          "type": "distribution_record_0",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif12jcegqefrulfcxp565lyjyv4s2ja82ndpjrcse"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "11000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_1",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif15pp0lqwq7squk9rdjejfdcqkf07apmcu2ym2zp"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "8000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_2",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif17pnxcmm2de3j4v3wmzwzrwx0vz5trchc2ysfmt"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "9000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_3",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif18ejp2jzgue7a4jfes3qq3l9n7q6yvztlk85ypu"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "5000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_4",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif199t07akr2cv6rhr8rlw50v4sz9lmkwyl42xlw4"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "12000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_5",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif1dfn94fp0z0dg3cvte4pqgjkj5rucr8x2j5avt6"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "14000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_6",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif1kxgqfh3a5wpzvqlu5tys2prs2wf5xemzzlpyfs"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "7000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_7",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif1pyu5ctnumet3s0hpy3jlqk32pgfsvr46v7v7uj"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "13000rowan"
            }
          ]
        },
        {
          "type": "distribution_record_8",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif1vm3g7fepdrygfwwc5gdfc9mv2zmvlkdga87fhv"
            },
            {
              "key": "type",
              "value": "LiquidityMining"
            },
            {
              "key": "amount",
              "value": "10000rowan"
            }
          ]
        },
        {
          "type": "distribution_run",
          "attributes": [
            {
              "key": "distribution_name",
              "value": "3_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
            },
            {
              "key": "distribution_runner",
              "value": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
            }
          ]
        }
```