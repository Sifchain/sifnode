# Dispensation Multiple Types
## Concepts

- We need the ability to add and remove dispensation types.The present design uses an enum for the types .

- The plan is to replace this enum with a register , which can be controlled by an admin account or through a governance proposal .


## Changes

- All are checks for claimtype are stateless, we do validations by matching the string to the enum value, this would need to be changed to statefull changes , therefore can only be processed in checkTx and deliverTx , and not ValidateBasic.This will change the User Experience a bit becuase the user might not get error responses immediatly , would need to be checked though `raw_log` or `events`    .

> Example : This function called from `dispensation/client/tx.go`
```go=
func GetDistributionTypeFromShortString(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "Airdrop":
		return DistributionType_DISTRIBUTION_TYPE_AIRDROP, true
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}
```

> Will change and be called from `dispensation/keeper/msg_server.go`

```go=
(func (keeper) GetClaimByType (ctx sdk.Context,claimType string) ClaimType {}).IsClaimable
```



- Changes in validation logic in msg_server for
```go=
func (srv msgServer) CreateUserClaim(ctx context.Context,
	claim *types.MsgCreateUserClaim) (*types.MsgCreateClaimResponse, error)
    
func (srv msgServer) CreateDistribution(ctx context.Context,
	msg *types.MsgCreateDistribution) (*types.MsgCreateDistributionResponse, error)
    
```

- Changes in
```go=
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64, distributionName string, authorizedRunner string, distributionType types.DistributionType) (*types.DistributionRecords, error) {
if record.DoesTypeSupportClaim() {}
}
```

## State

```go=
type ClaimType struct {
    Type : string ,
    IsClaimable : bool ,
    IsActive : bool ,
}
```

## State Transitions

```go=
func (keeper) SetNewClaimType(ctx sdk.Context,c Claimtype){} {
  // Can only be called by admin
}
```

## Helpers

```go=
func (keeper) IsClaimTypeActive (ctx sdk.Context,c Claimtype) bool {}
func (keeper) IsClaimTypeClaimabale (ctx sdk.Context,c Claimtype) bool {}
```

## Queries

```go=
func (keeper) GetAllClaimTypes(ctx sdk.Context)[]ClaimType {}
func (keeper) GetClaimByType(ctx sdk.Context,claimtype string)ClaimType {}

```


## Migrations

- Claimtypes for Existing Distribution records would need to be modified, to use `ClaimType` struct instead of the enum.
- Set Claimtypes in InitGenesis

## Misc Questions

- We use the same type for Claims and Distributions . It was called distributionType earlier . I am using ClaimType in this doc , but the name is debatable . I would prefer to use a more generic name if possible , and any suggestions for be great .

- Should we test and move Distribute Records back to block ender / beginner . It would save a lot of the logistics we are dealing with , using the whole run distribution . The mechanism is used in other places in the module , and if tested properly should be doable .
