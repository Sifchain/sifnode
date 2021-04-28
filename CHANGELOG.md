
----

# v0.8.0
> April 27, 2021

## ⭐ Features

- [UI] Prevent pegging when not connected to a supported EVM network.
- [UI] Arbitrage Opportunity - Now color-coded to indicate in which direction the opportunity is present.

## 🛠 Improvements

- [Sifnode] Ability to propose and vote on new jailing parameters.
- [ChainOps] MongoDb resource limit increases.
- [Blockexplorer] Improved design elements and cUSDT decimal formatting.
- [UI] Incorporated logic in token listing when swapping and adding liquidity. These will now be sorted by balance (high to low) and then alphebetically. Also added the ability to see your balance in these listings as well. - TBD
- [UI] Added logic for when clicking on max button in ROWAN to take into consideration necessary gas fees. - TBD
- [UI] In swap confirmation screen, built in cleaner UX logic around the way we display numbers. 

## 🐛 Bug Fixes

- [UI] State loader for dispensation.
- [UI] Remove the "select all" functionality when clicking in a field.
- [UI] Token with zero balances would sometimes disappear from the swap options.
- [Peggy] Add retry logic when infura gives us a not found error. Add additional retry logic to try to retrieve the tx if it cannot be found on the first query.

## ❓ Miscellaneous

- [UI] Integration of the Playwright test framework.
- [UI] Amount API for appropriate decimal placement across all token types.
- [Peggy] Ability to export the Ethbridge keeper data (required for when migrating to Cosmos 0.42.x).
- [ChainOps] Automated pipeline deployment.
