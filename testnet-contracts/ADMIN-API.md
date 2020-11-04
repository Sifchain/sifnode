# This document describes how to hook eRowan into peggy as a Cosmos Native Asset

## Please follow these instructions to the T. If you do not, peggy may not behave in the way you expect it to with eRowan

1. You will need to have created eRowan as an ERC20 token on the mainnet. The token's symbol should be "PEGGYeRowan"

2. You will need to set the BridgeBank contract as an admin and minter role so that it can mint new tokens when assets are locked on the sifchain side.

3. Call the function addExistingBridgeToken on BridgeBank and pass the address of eRowan as the first parameter.

4. Whenever you make a new prophecy claim for rowan, you will need to pass the symbol eRowan as a parameter, otherwise it will mess create a new token when it should not.