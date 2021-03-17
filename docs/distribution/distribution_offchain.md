# Distribution Module
## Requirements
- Create mechanism which allows us to distribute different types of rewards ( liquidity mining rewards being one of them ) .
- Option to pause distributions.


## Concepts
- The distribution calculation is off-chain . The json input is parsed and stored to use as a basis for distribution.
- Every distribution will start with a distibutionList. The distribution list contains all parameters to facilitate a distribution .
- A list of distibutionLists is stored in the keeper.
- At the next blockbeginner we check if we need to distribute rewards in this block or not . If we do we distribute the rewards and update the map.


## Data-Structures
```go
type DistributionList struct {
	IsActive bool
	Identifier string // or id 
	ListType listType
	FundingAddress sdk.Address  // Assuming that the funding address would activate a distribution list
	Receivers receiverType
	TotalRewards sdk.Uint
	ReceivingList map([]sdk.Address)sdk.Coins  // I have see some issues with map and amino before .Not 100 % if this is the best idea . Will need to look into for cosmos handles the deserialization
	DistributionFunction distributionFunction
	DistributionFrequency int64 // in num of blocks
	DistributionTokens []string // Type of token to be distributed
	// IntermediaryAddress sdk.Address or string 
}
```

We can add a field here to store an intermediary address , to which all funds are transferred at the block ender .
I chose not to , because the funding address for all use cases would be controlled by us  .This simplifies the logic .
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

- Iterate over active distribution lists

- Iterate over receiverTypes , and call getAddressList() on each 

- Append all receiving address to an address list .

- Use getRewards() function of Rewards interface to get total rewards for the current block

- Call distributionFunction on complete address list . Pass the address list , and any other parameters this function might need.

- Store the returned map in ReceivingList , append values if addresses are present it the list . Create new entries if they are not. 

- When creating the receiving list use TokenName to create a Coin from the generated amount

## BlockBeginnerLogic

- Iterate over active distribution lists

- Use DistributionFrequency parameter and block height to create to check if we need to distribute in the present block.

- If true from above iterate over the ReceivingList and distribute coins.

- Use Funding address to as from address for distribution.

- When distributed deduct amount from total rewards

## How the developed module would behave 
- We would be defining types by Implementing the various interfaces .
- We create distribution lists , by using the types defined 
```go
func GetLiquidityMiningDistibutionList(/* Cli input for values*/) types.DistributionList {
        return DistributionList {
        IsActive false
        Identifier "CLI INPUT"
        ListType LiquidityMining
        FundingAddress "CLI INPUT"
        Receivers []ReceiverType{Validators{}, LiquidityProviders{}}
        TotalRewards "CLI INPUT"
        ReceivingList make(map([]sdk.Address)sdk.Coins)
        DistributionFunction DistributeLiquidityMiningRewards
        DistributionFrequency 1 
        DistributionTokens []string{"rowan"} 
        }
      }
```
- The creation of the list would take some inputs from the user, and some default values, I feel we can take only Identifier ,Funding address and Total funds as input and use a switch case for creating the lists. We would have a maximum of 3 different types of lists which can be defined easily. 
- Sample cli command to create the above list 
```shell
sifnodecli tx distribution create --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd LiquidityMining SifLiquidityRewards  21000000000000000000000
```
- The command would run through a switch case like 
```go
switch ListType {
 case LiquidityMining : GetLiquidityMiningDistibutionList
 _
}
```
- Save the list with " active = true " . Key will be the identifier .
### Additional functionality

- Deactivate distribution - funding address can send a transaction to pause distribution
```shell
sifnodecli tx distribution deactivate --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd SifLiquidityRewards
```

- Activate distribution - funding address can send a transaction to pause distribution
```shell
sifnodecli tx distribution activate --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd SifLiquidityRewards
```

- (Optional )While distributing in the BlockBeginner we can add an entry to the keeper for rewards earned by an address .This can be saved as
  
  Key : Address_ListType_Height
  
  Value : Reward_Earned

   This can be used to provide various queries to users such as 
- Query rewards earned at a particular height / for listtype /combination of both etc
```shell
sifnodecli q distribution rewardsforheight sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd 1
sifnodecli q distribution rewardsforlistname sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd SifLiquidityRewards
```

  
### Points to note
- This is a very high level document for the overall logic. 
- The module would need to handle state export etc , for upgrades to happen , which would be part of a subsequent document.

## Sample Code 
We can use functions instead of interfaces , Choosing to use interfaces because it would allow us to differentiate types easily .
Also the structs implementing the interface can hold parameters required fo the underlying functions .

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
