# README



## Changelog

- 2020/10/22 Initial version

### Summary

For the Sifchain MVP ,CLP module provides the following functionalities 
- Create New Liquidity Pool
- Add Liquidity to an Existing Liquidity pool 
- Remove Liquidity from an Existing Liquidity pool 
- Swap tokens  
        -Swap an External token for Native or vice versa (single swap)    
        -Swap an External token for another External Token (double swap) 
- Decommission an Existing Liquidity pool 

### Data structures 
-Asset : An asset is most basic unit of a CLP . It Contains source chain, symbol and ticker to identify a token .
```golang
SourceChain: ETHEREUM
Symbol: ETH
Ticker: ceth

SourceChain: SIFCHAIN
Symbol: RWN
Ticker: rwn
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
    
## Trasactions supported
 - **Create new liquidity pool**
    - Creating a pool has a minimum threshold for the amount of liquidity provided. This is a genensis parameter and can be tweaked later.
    - The user who creates a new pool automatically becomes its first liquidity provider.
 - **Decommission a liquidity pool** 
    - Decommission requires the net balance of the pool to be under the minimum threshold . 
    - If successful a decommission transaction returns balances to its liquidity providers and deletes the liquidity pool. 
 - **Add Liquidity to a pool** 
    - User can add liquidity to the native and external tokens 
 - **Remove liquidity**
    - Remove liquidity consists of a composition of withdraw , and a swap if required
    - Liquidity can be removed in three ways
    
        -Native and external - Withdraw to native and external tokens .   
        -Only Native -  Withdraw to native and external tokens ,and then a swap from external to native.   
        -Only External  - Withdraw to native and external tokens ,and then a swap from native to external.   
   - For asymmetric removal , (option 2 and 3), the user incurs a tradeslip and liquidity fee similar to a swap.
 - **Swap**
    
    - The system supports two types of swaps          
        -Swap between external and native tokens - This is a single swap        
        -Swap between external and external tokens - This swap is combination of two single swaps.
        
    - A double swap also includes a transfer between the two pools to maintain pool balances.