# ADR 001: CLP integration with Sifnode

## Changelog

- 2020/11/02 Initial version

## Status

*Proposed*

## Context
This ADR contains a summary of decisions taken during the implementation of the clp module and a summary of the final outcome
### Summary

For the Sifchain MVP, CLP module provides the following functionalities 
- Create New Liquidity Pool
- Add Liquidity to an Existing Liquidity pool 
- Remove Liquidity from an Existing Liquidity pool 
- Swap tokens  
        -Swap an External token for Native or vice versa (single swap)    
        -Swap an External token for another External Token (double swap) 
- Decommission an Existing Liquidity pool 

### Basic Terminology 
-Asset : An asset is most basic unit of a CLP . It Contains symbol to identify a token .
```golang
Symbol: ETH

Symbol: rowan
```
-Pool  : Every Liquidity pool for CLP is created by pairing an External asset with the Native asset .
````golang
ExternalAsset: SourceChain: ETHEREUM
              Symbol: ETH
              Ticker: ceth
ExternalAssetBalance: 1000
NativeAssetBalance: 1000
PoolUnits : 1000
PoolAddress :sif1vdjxzumgtae8wmstpv9skzctpv9skzct72zwra
````
-Liquidity provider : Any user adding liquidity to a pool becomes a liquidity provider for that pool. 
````golang
ExternalAsset: SourceChain: ETHEREUM
               Symbol: ETH
               Ticker: ceth
LiquidityProviderUnits: 1000
liquidityProviderAddress: sif15tyrwghfcjszj7sckxvqh0qpzprup9mhksmuzm 
````
    
## Decisions 
 - **Create new liquidity pool**
    - Creating a pool has a minimum threshold for the amount of liquidity provided. This is a genesis parameter and can be tweaked later.
    - The user who creates a new pool automatically becomes its first liquidity provider.
    - Every pool has been decided to have a different pool address .The pool address the created from the string (External_Asset_Ticker)_(Native_Asset_Ticker) .
    - Pool units are calculated based on the external and native asset amount . the formula used is . The pool creator gets his share of units.
    ````
      r = Native amount added
      a = External amount added
      R = Native Balance (before)
      A = External Balance (before)
      P = existing Pool Units
      slipAdjustment = (1 - ABS((R a - r A)/((r + R) (a + A))))
      units = ((P (a R + A r))/(2 A R))*slipAdjustment
   ````
    ***Consequences***
    - Positive - Every pool has a unique address
    - Negative - We need to provide interpool transfers to facilitate double swaps
    - Neutral  - The pool address generated from the string is not always 20Bytes , it might need to be padded with bytes to confirm to the cosmos standards.
                 The padding is done by a one to one copy and adding extra bytes to the string ,before deriving an address.
 - **Decommission a liquidity pool** 
    - Decommission requires the net balance of the pool to be under the minimum threshold . 
    - If successful a decommission transaction returns balances to its liquidity providers and deletes the liquidity pool.
    - We use the same function as remove liquidity to calculate withdrawal for each liquidity provider.
    
    ***Consequences***
    - Positive - The pool threshold can be managed differently for every pool.
    - Negative - Since we are using floating point calculations ,rounding off might result in some tokens being left in the pool.These need to be managed separately.
    - Neutral  - None
    
 - **Add Liquidity to a pool** 
    - User can add liquidity to the native and external tokens .
    - Liquidity can be added asymmetrically .
    - The same formula is used to calculate pool units ,and the sender is allocated his share.
    
    ***Consequences***
    - Positive - Liquidity can be added asymmetrically.
    - Negative - None
    - Neutral  - None  
 - **Remove liquidity**
    - Remove liquidity consists of a composition of withdraw , and a swap if required
    - Liquidity can be removed in three ways
    
        -Native and external - Withdraw to native and external tokens .   
        -Only Native -  Withdraw to native and external tokens ,and then a swap from external to native.   
        -Only External  - Withdraw to native and external tokens ,and then a swap from native to external.   
   - For asymmetric removal , (option 2 and 3), the user incurs a tradeslip and liquidity fee similar to a swap.
   - The pool is checked for being shallow ( Amount of an asset either native or external dropping to zero ),and the transaction is rejected if that happens.
   - The check is done after withdraw and then again after swap .
   - The range for wBasisPoints was decided to be 0 - 10000
   - The range for Asymmetry was decided to be -10000 to 10000 .
   - To calculate withdrawal amount ,the function converts all values to float .This is done to avoid divide by zero errors.
     The calculation formula used is
     ````
        {
          unitsToClaim = lpUnits / (10000 / wBasisPoints) 
          withdrawExternalAssetAmount = externalAssetBalance / (poolUnits / unitsToClaim)
          withdrawNativeAssetAmount = nativeAssetBalance / (poolUnits / unitsToClaim)
          
          swapAmount = 0
          //if asymmetry is positive we need to swap from native to external
          if asymmetry > 0
            unitsToSwap = (unitsToClaim / (10000 / asymmetry))
            swapAmount = nativeAssetBalance / (poolUnits / unitsToSwap)
        
          //if asymmetry is negative we need to swap from external to native
          if asymmetry < 0
            unitsToSwap = (unitsToClaim / (10000 / asymmetry))
            swapAmount = externalAssetBalance / (poolUnits / unitsToSwap)
        
          //if asymmetry is 0 we don't need to swap
          
          lpUnitsLeft = lpUnits - unitsToClaim
          
          return withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount
        }
     ````  
        ***Consequences***
        - Positive - This gives flexibility for the suers to withdraw
        - Negative - None
        - Neutral  - None
     
 - **Swap**
    
    - The system supports two types of swaps          
        -Swap between external and native tokens - This is a single swap        
        -Swap between external and external tokens - This swap is combination of two single swaps.
        
    - A double swap also includes a transfer between the two pools to maintain pool balances.
    - The formula used for a token swap is as follows .
    ````
   swapOne(keeper, sentAsset, sentAmount, receivedAsset): Swaps amount of sent asset for received asset from specified asset pool.
   
   Uses Thorchain's slip based Continuous Liquidity Pool model to calculate trade slip, liquidity fee, and resulting swap.
   
   {
     if pool for sentAsset exists
       getPool(sentAsset)
     if receivedAsset is native
       Y = pool.nativeAssetBalance
       X = pool.externalAssetBalance
     else
       X = pool.nativeAssetBalance
       Y = pool.externalAssetBalance
     x = sentAmount
     liquidityFee = calcLiquidityFee(X, x, Y)
     tradeSlip = calcTradeSlip(X, x)
     swapResult = calcSwapResult(X, x, Y)
     if swapResult >= Y
       fail, “not enough received asset tokens to swap”
     if sentAsset is native
       pool.nativeAssetBalance = X + x
       pool.externalAssetBalance = Y - swapResult
     else
       pool.externalAssetBalance = X + x
       pool.nativeAssetBalance = Y - swapResult
     setPool()
   
     record liquidityFee for use in liquidity rewards
     return emitAmount
   }
   
   
   calcLiquidityFee(X, x, Y): calculates liquidity fee per Thorchain's CLP model
   
   {
     return (x * x * Y) / ((x + X) * (x + X))
   }
   
   
   calcPriceImpact(X, x): calculates price impact per Thorchain's CLP model
   
   {
     return x / (x + X)
   }
   
   
   calcSwapResult(X, x, Y): calculates final swap received token amount
   
   {
     return (x * X * Y) / ((x + X) * (x + X))
   }
   ````
   ***Consequences***
   - Positive - Users have the flexibility to swap any asset for any different asset 
   - Negative - None
   - Neutral  - Node


  ## References
 https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
 https://docs.thorchain.org/how-it-works/continuous-liquidity-pools
 
