# Peggy Admin Smart Contract Actions

In order to give sifchain more control of peggy once it gets released into the wild, pause functionality has been implemented on lock, burn, unlock and mintBridgeTokens functions. This way, if a critical bug is discovered in the bridge bank contract, no funds will be able to be withdrawn once the contract is paused.

The function to call to pause contracts is pause(). Only the bridgebank pauser is allowed to call this function. Once the contract is paused, it can be unpaused by calling the unpause() function as the owner. The pause function is not callable while the contract is already paused. The unpause function is not callable if the contract is not paused. 