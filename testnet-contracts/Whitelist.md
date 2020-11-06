# Important Notes on the Whitelist

You should only be able to whitelist a token symbol once. There is logic written that will stop you from adding a token smart contract with the same symbol as one that has already been whitelisted. That being said, don't attempt to add a token contract to the whitelist if a token with that same symbol has already been whitelisted.

Do not add a token or smart contract to this contract that is upgradeable or has the ability to change its token symbol name in the future. This could lead to damage on our side. 

# Admin Responsibilities
As the admin, it is your responsibility to make sure that no contract can get added to the whitelist that has the ability to change its symbol.

# Final warnings
Read the important notes and admin responsibilites.