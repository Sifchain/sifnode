# DISPENSATION MODULE (Claims)

## Overview
- The module allows a user to create a Claim .
- A claim is record , stating that the user want to collect their rewards at the end of the week.
- A claim is deleted once the user is paid out.
- The claim_create is an auxiliary functionality for the offline calculation api created to calculate user rewards.

## General use case 
This is the general use case , the actual implementation might vary a bit
Any day of the week
 - create claims through this api (on chain)
On friday
- Run a query to get all claims / Keep reading the events emitted from create_claim
- This is list is an input for a function (This function is off-chain and not part of this module)  which iterates over the list and creates a json which the dispensation module can use.
- Run a dispensation of type Liquidity Mining or Validator Subsidy . 
- After running the above the transfers happen over the next few blocks (10 per block) . An external function (Not part of this module).Reads these events and resets user multipliers accordingly.


## Technicals
### Data structures
```go
type UserClaim struct {
	UserAddress   sdk.AccAddress   `json:"user_address"`
	UserClaimType DistributionType `json:"user_claim_type"`
	UserClaimTime time.Time        `json:"user_claim_time"`
	Locked        bool             `json:"locked"`
}
```
### High Level Flow
- User creates claims of type LiquidityMining or ValidatorSubsidy.
- The claims get locked when the associated distribution is created for the claim  
- The claims are deleted when the user receives the funds . 
- A user can only have 1 claim of a particular type ,at any point in time .

### User Flow
- Create Claim
```shell
sifnoded tx dispensation claim ValidatorSubsidy --from akasha --keyring-backend test --yes
```

### Events Emitted
- The following event is emitted on every claim
```json
 {"type": "claim_created",
            "attributes": [
              {
                "key": "claim_creator",
                "value": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5"
              },
              {
                "key": "claim_type",
                "value": "ValidatorSubsidy"
              },
              {
                "key": "userClaim_creationTime",
                "value": "2021-05-02T02:43:10.593125Z"
              }
            ]
}
```

### Queries
- Query to get claims by type
```shell
sifnoded q dispensation claims-by-type ValidatorSubsidy
```
Response 
```json
{
  "claims": [
    {
      "user_address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
      "user_claim_type": "3",
      "user_claim_time": "2021-05-02T02:43:10.593125Z",
      "locked": false
    },
    {
      "user_address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
      "user_claim_type": "3",
      "user_claim_time": "2021-05-02T02:43:10.593125Z",
      "locked": false
    }
  ],
  "height": "7"
}
```