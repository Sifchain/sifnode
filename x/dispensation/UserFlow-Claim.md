# Claims Flow

### Users will create claims throughout the week . 
The wallet can use 
```/dispensation/createClaim```  
This expects the following input 
```go
package main

type DistributionType int64
const Airdrop DistributionType = 1
const LiquidityMining DistributionType = 2
const ValidatorSubsidy DistributionType = 3

type CreateClaimReq struct {
	BaseReq   rest.BaseReq     `json:"base_req"`
	Signer    string           `json:"signer"`
	ClaimType DistributionType `json:"claim_type"`   
}
```

### On friday we get a list of all the claims (After Cut-off time) for the week. Any claims submitted after cutoff would be processed next week.
This query through the cli would look like
```shell
sifnodecli q dispensation claims-by-type ValidatorSubsidy --chain-id sifchain --node tcp://rpc.sifchain.finance:80
```
Which returns 
```json
{
  "claims": [
    {
      "user_address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
      "user_claim_type": "3",
      "user_claim_time": "2021-05-02T02:43:10.593125Z",
    },
    {
      "user_address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
      "user_claim_type": "3",
      "user_claim_time": "2021-05-02T02:43:10.593125Z",
    }
  ],
  "height": "7"
}

```
The relevant event would be in the same block as the one which has the dispensation/createClaim request
```json
 {"type": "userClaim_new",
            "attributes": [
              {
                "key": "userClaim_creator",
                "value": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
              },
              {
                "key": "userClaim_type",
                "value": "ValidatorSubsidy"
              },
              {
                "key": "userClaim_creationTime",
                "value": "2021-05-02T02:43:10.593125Z"
              }
            ]
}

```

### This list obtained above is run through the parsing API ,to get output list
```json
{
 "Output": [
  {
   "address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
   "coins": [
    {
     "denom": "rowan",
     "amount": "10000000000000000000"
    }
   ]
  },
  {
   "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
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

### This file is then used to create and run a dispensation
Create
```shell
sifnodecli tx dispensation create ValidatorSubsidy output.json sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan
```
Run
```shell
sifnodecli tx dispensation run 2_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd ValidatorSubsidy --from sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan
```
