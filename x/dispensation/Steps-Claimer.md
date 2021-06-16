#Steps to follow by User trying to claim rewards

##Requirements
- The user should have rewards to claim.
- A user is allowed to claim rewards of type LiquidityMining or ValidatorSubsidy
- A user can only claim once per distribution cycle , for a particular type. This means a user can have a maximum of two claims (one for ValidatorSubsidy and one for LiquidityMining at any given time)

##Steps
###Create claim
Claims can be created using the CLI or REST interface 
The usual flow for creating claims would be through the wallet ,which uses the following rest API
####REST
```/dispensation/createClaim```  
The API expects the following input
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


####CLI
```shell
sifnodecli tx dispensation claim LiquidityMining --from sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan --node tcp://rpc.sifchain.finance:80 --chain-id sifchain
```

The relevant event 
```json
 {"type": "userClaim_new",
            "attributes": [
              {
                "key": "userClaim_creator",
                "value": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
              },
              {
                "key": "userClaim_type",
                "value": "LiquidityMining"
              },
              {
                "key": "userClaim_creationTime",
                "value": "2021-05-02T02:43:10.593125Z"
              }
            ]
}

```