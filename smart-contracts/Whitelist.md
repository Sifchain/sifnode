# Important Notes on the Whitelist

You should only be able to whitelist a token symbol once. There is logic written that will stop you from adding a token smart contract with the same symbol as one that has already been whitelisted. That being said, don't attempt to add a token contract to the whitelist if a token with that same symbol has already been whitelisted.

Do not add a token or smart contract to this contract that is upgradeable or has the ability to change its token symbol name in the future. This could lead to damage on our side. 

# Admin Responsibilities
As the admin, it is your responsibility to make sure that no contract can get added to the whitelist that has the ability to change its symbol. Do not allow proxy ERC20 contracts or any ERC20's that have an admin with the ability to change the token symbol.


# How to use
Here is how to whitelist smart contracts using these scripts locally:
```
truffle exec scripts/sendUpdateWhiteList.js 0x0d8cc4b8d15D4c3eF1d70af0071376fb26B5669b true
```

Insert your own token address that you want to whitelist and use the boolean that false to remove the smart contracts from the whitelist.

Here is how to use this script on a main or testnet:
```
truffle exec scripts/sendUpdateWhiteList.js --network ropsten 0x0d8cc4b8d15D4c3eF1d70af0071376fb26B5669b true
```

# Final warnings
Read the important notes and admin responsibilites.