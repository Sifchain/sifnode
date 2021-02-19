# Distribution Module

## Concepts
- Every distribution will start with a distibutionList. The distribution list contains all parameters to facilitate a distribution .
- A list of distibutionLists is stored in the keeper.
- At every blockender we iterate over these lists and calculate the rewards for each receiver for this block .
- At the next blockbeginner we check if we need to distribute rewards in this block or not . If we do we distribute the rewards and update the map.
- If we don't the rewards keep adding up to .

## Data-Structures
```go
type DistributionList struct {
	Identifier string // or id 
	FundingAddress sdk.Address
	Receivers receiverType
	TotalRewards sdk.Uint
	ReceivingList map([]sdk.Address)sdk.Coins
	DistributionFunction distributionFunction
	DistributionFrequency int64 // in num of blocks
	DistributionTokens []string // Type of token to be distributed
}
```

We can add a field here to store and intermediary address , to which all funds are transferred at the block ender 
I chose not to , because the funding address for all use cases would be controlled by us.This simplifies the logic a lot .
We can add that feature in the future
## Interfaces

### receiverType Interface 
```go
type receiverType interface {
	getAddressList() []sdk.address
	setAddressList([]sdk.address)  
}
```
```go
type Validators struct {
	// implements receiverType
}
func (Validators)getAddressList() []sdk.address {
	// Iterate over validator set and return list of validators
}
type LiquidityProviders struct {
	// implements receiverType
}
func (LiquidityProviders)getAddressList() []sdk.address {
// Iterate over liquidityProviders and return list of addresses
}


```

### distributionFunction Interface
```go
type distributionFunction interface {
	getAmounts(sdk.Uint,[]string,[]sdk.address,{}interface) map[sdk.address]amounts
}

```
```go
type DistributeLiquidityMiningRewards struct {}
func (DistributeLiquidityMiningRewards)getAmounts(
	    totalRewards sdk.Uint,
	    distributionTokens []string, 
	    addrlist []sdk.address,
	    i {}interface) map[sdk.address]Coins {
	// parse interface into struct
	// use Total rewards to check if required amount is available.
	// Iterate over all addresses .
	// Logically allocate tokens  
	// Create output map
}
```

### rewards Interface

```go
type Rewards interface {
	getRewards() sdk.Amount
}
```
```go
type CollectLiquidityMiningRewards struct {}
func (CollectLiquidityMiningRewards)getRewards([]sdk.address) map[sdk.address]amounts {
	// Logic to get total rewards for this block
}
```


Note the structs can be used to store data required by the underlying functions.

## Keeper Functions

```go
func SetDistributionList() {}
func GetDistributionList() {}
func IterateAllLists() {}
```


## BlockEnderLogic

- Iterate over all distribution lists

- Iterate over receiverTypes , and call getAddressList() on each 

- Append all receiving address to an address list .

- Use getRewards() function of Rewards interface to get total rewards for the current block

- Call distributionFunction on complete address list . Pass the address list , and any other parameters this function might need.

- Store the returned map in ReceivingList , append values if addresses are present it the list . Create new entries if they are not. 

- When creating the receiving list use TokenName to create a Coin from the generated amount

## BlockBeginnerLogic

- Iterate over distribution lists

- Use DistributionFrequency parameter and block height to create to check if we need to distribute in the present block.

- If true from above iterate over the ReceivingList and distribute coins.

- Use Funding address to as from address for distribution.

- When distributed deduct amount from total rewards

### Points to note
- This is a very high level document for the overall logic. 
- The module would need to handle state export etc , for upgrades to happen , which would be part of a subsequent document.

## Sample Code 
We can use functions instead of interfaces , Choosing to use interfaces becuase it would allow us to differentiate types easily 

```go
package main

import "fmt"

type DistributionList struct {
	Receivers []ReceiverType
}

func main() {

	// Create new Distribution list
	d1 := DistributionList{
		// Add receivers
		Receivers: []ReceiverType{Validators{}, LiquidityProviders{}},
	}
	d2 := DistributionList{
		// Add receivers
		Receivers: []ReceiverType{Validators{}},
	}

	dList := []DistributionList{d1,d2}
	
	// BlockEnder Logic
	for _,dl := range dList {
		for _, r := range dl.Receivers {
			r.GetAddressList()
		}
	}
}

type ReceiverType interface {
	GetAddressList()
}
// Create a receiver type
type Validators struct {
}
func (Validators) GetAddressList() {
	fmt.Println("Validator address list returned here")
}

// Create another receiver type
type LiquidityProviders struct {
}
func (LiquidityProviders) GetAddressList() {
	fmt.Println("LiquidityProvider address list returned here")
}

```
