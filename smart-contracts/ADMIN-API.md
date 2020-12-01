# This document describes how to hook eRowan into peggy as a Cosmos Native Asset

## Please follow these instructions to the T. If you do not, peggy may not behave in the way you expect it to with eRowan

# Most important thing for deployment to test or mainnet

Make sure that the OWNER address is set properly in the .env file so that you have an owner for the bridgebank contract that can use the admin api. Owner and operator addresses are different roles and have different capabilities.

1. You will need to have created eRowan as an ERC20 token on the mainnet. The token's symbol should be "eRowan"

This symbol is hard coded in the eRowan token repo so as long as you use that implementation, you should be good.

2. You will need to set the BridgeBank contract as an admin and minter role so that it can mint new tokens when assets are locked on the sifchain side.

Follow step 6 in Deployment.md to do this.

3. Call the function addExistingBridgeToken on BridgeBank and pass the address of eRowan as the first parameter.

Follow step 6 in Deployment.md to do this.

4. Whenever you make a new prophecy claim for rowan, you will need to pass the symbol eRowan as a parameter, otherwise it will mess up and create a new token when it should not.