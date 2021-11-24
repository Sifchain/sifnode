# EVM events

## LogLock
- `_from`: Ethereum address that initiated the lock 
- `_to`: the sifchain address that the imported assets should be credited to (UTF-8 encoded string)
- `_token`: the token's contract address or the null address for EVM-native currency
- `_value`: the quantity of asset being transferred (a uint256 representing the smallest unit of the base value)
- `_nonce`: the current transaction sequence number which is indexed as a topic (_nonce) (this value increments automatically for each `lock`)
- `_decimals`: the decimals of the asset which defaults to 18 if not found
- `_symbol`: the symbol of the asset which defaults to empty string if not found
- `_name`: the name of the asset which defaults to empty string if not found (_name)
- `_networkDescriptor`: the network descriptor for the chain this asset is on


# Cosmos events


## NewEthBridgeClaim
- ...
- ...

## EventTypeBurn
- ...
- ...
