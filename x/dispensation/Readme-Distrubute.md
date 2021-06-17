# DISPENSATION MODULE (Distribution)

## Overview
- The module allows a user to create a Distribution which can be of Type [Airdrop/LiquidityMining/ValidatorSubsidy]. 
- It accepts a output list and an autorized runner .
- The records are created in the same block.
- The distribution process can be manually executed by the authorized runner.



## Technicals 
### Data structures
 - The base-level data structure is 
```go
package records
type DistributionType int64
const DistributionTypeUnknown DistributionType = 0
const Airdrop DistributionType = 1
const LiquidityMining DistributionType = 2
const ValidatorSubsidy DistributionType = 3

type Distribution struct {
	DistributionType DistributionType `json:"distribution_type"`
	DistributionName string           `json:"distribution_name"`
	Runner           sdk.AccAddress   `json:"runner"`
}
```
This is stored in the keeper with the key DistributionName_DistributionType for historical records. 
The Distribution name is BlockHeight_Distributor ,Therefore, the combination BlockHeight_Distributor_DistributionType needs to be unique

- Distribution records are created for processing individual transfers to recipients

```go
package records

type DistributionStatus int64

const Pending DistributionStatus = 1
const Completed DistributionStatus = 2



type DistributionRecord struct {
	DistributionStatus          DistributionStatus `json:"distribution_status"`
	DistributionName            string             `json:"distribution_name"`
	DistributionType            DistributionType   `json:"distribution_type"`
	RecipientAddress            sdk.AccAddress     `json:"recipient_address"`
	Coins                       sdk.Coins          `json:"coins"`
	DistributionStartHeight     int64              `json:"distribution_start_height"`
	DistributionCompletedHeight int64              `json:"distribution_completed_height"`
	AuthorizedRunner            sdk.AccAddress     `json:"authorized_runner"`
}
```
This record is also stored in the keeper for historical records .

### High Level Flow
- The `create tx` sends required funds (sum of outputs ) from the distributor address to a module account.
- The program iterates over the output addresses and creates individual records for them in the keeper .
- The `run tx`  iterates over these records and completes 10 records per block .
- Complete refers to sending the specified amount from the  module account to the recipient.
- In case of type LiquidityMining or ValidatorSubsidy the program also deletes the associated claim.

### Events Emitted
- A `create_tx` emits a distribution_started event .
```json
{
        "type": "distribution_started",
        "attributes": [
        {
        "key": "module_account",
        "value": "sif1zvwfuvy3nh949rn68haw78rg8jxjevgm2c820c"
        },
        {
        "key": "distribution_name",
        "value": "1158855_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
        },
        {
        "key": "distribution_type",
        "value": "LiquidityMining"
        }
]
}
```
- The `run tx` emits a list of records rewarded in that block in its events , and a distribution_run event .
```shell
 {
          "type": "distribution_record_0",
          "attributes": [
            {
              "key": "recipient_address",
              "value": "sif12jcegqefrulfcxp565lyjyv4s2ja82ndpjrcse"
            },
            {
              "key": "type",
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "Airdrop"
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
              "value": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
            }
          ]
        }
```

### Queries supported
```shell
#Query all distributions
sifnodecli q dispensation distributions-all
#Query all distribution records by distribution name 
sifnodecli q dispensation records-by-name-all ar1
#Query pending distribution records by distribution name 
sifnodecli q dispensation records-by-name-pending ar1
#Query completed distribution records by distribution name
sifnodecli q dispensation records-by-name-completed ar1
#Query distribution records by address
sifnodecli q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00
```
