# DISPENSATION MODULE (Claims)

## Overview
- The module allows a user to create a Claim .
- A claim is record , stating that the user want to collect their rewards at the end of the week.
- A claim is deleted once the user is paid out.
- The `create claim` is the on-chain record for the off-chain components to calculate ValidatorSubsidy and LiquidityMining rewards

## General use case 
This is the general use case , the actual implementation might vary a bit
Any day of the week
 - create claims through this api (on chain)
Off-Chain components
 - Create claim snapshots every 200 hours .  
  `@Niko To add the steps here`
On friday
 - Create a dispensation of type Liquidity Mining or Validator Subsidy . 
 - Run the created dispensation using the authorize runner


## Technicals
### Data structures
```go
type UserClaim struct {
	UserAddress   sdk.AccAddress   `json:"user_address"`
	UserClaimType DistributionType `json:"user_claim_type"`
	UserClaimTime time.Time        `json:"user_claim_time"`
}
```
### High Level Flow
- User creates claims of type LiquidityMining or ValidatorSubsidy.
- The claims are deleted when the user receives the funds . 
- A user can only have 1 claim of a particular type ,at any point in time .

### User Flow
- Create Claim
```shell
sifnodecli tx dispensation claim ValidatorSubsidy --from akasha --keyring-backend test --yes
```

### Events Emitted
- The following event is emitted on every claim
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

### Queries
- Query to get claims by type
```shell
sifnodecli q dispensation claims-by-type ValidatorSubsidy
```
Response 
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