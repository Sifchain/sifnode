# Distribution Module
## Data-Structures
```go
type DistributionList struct {
	Identifier string // or id 
	FundingAddress sdk.Address
	Receivers receiverType 
	ReceivingList map([]sdk.Address)sdk.Amount 
	RewardsAllocated rewards
	DistributionFunction distributionFunction
	DistributionFrequency int64 // in num of blocks
	TokenName string // Type of token to be distributed
}
```
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
	getAmounts([]sdk.address) map[sdk.address]amounts
}

```
```go
type DistributeLiquidityMiningRewards struct {}
func (DistributeLiquidityMiningRewards)getAmounts([]sdk.address) map[sdk.address]amounts {
	//Iterate over all addresses .
	// Logically allocate tokens .
	// Create output map
}
```

###rewards Interface

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


Keeper stores a list of distributionList 

##Keeper Functions
```go
func SetDistributionList() {}
func GetDistributionList() {}
func IterateAllList() {}
```


##BlockEnderLogic
-Iterate over all distribution lists

-Iterate over receiverTypes , and call getAddressList() on each

-Append all receiving address to an address list .

-Use getRewards() function of Rewards interface to get total rewards for the current block

-Call distributionFunction on complete address list .

-Store the returned map in ReceivingList , append values if addresses are present it the list . Create new entries if they are not. 
##BlockBeginnerLogic

-Iterate over distribution lists

-Use DistributionFrequency parameter and block height to create to check if we need to distribute in the present block.

-If true from above iterate over the ReceivingList and distribute tokens . Use Token name for type of token to distribute.

-Use Funding address to as from address for distribution.

##Sample Code 

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
