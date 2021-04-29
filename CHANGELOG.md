
----

# v0.8.1
> April 27, 2021

## ⭐ Features

- [UI] Prevent pegging when not connected to a supported EVM network.
- [UI] Arbitrage Opportunity - Now color-coded to indicate in which direction the opportunity is present.

## 🛠 Improvements

- [UI] Added logic for when clicking on max button in ROWAN to take into consideration necessary gas fees.
- [UI] In swap confirmation screen, built in cleaner UX logic around the way we display numbers. 
- [UI] Updated design elements (typography and general look & feel) across the DEX.
- [Sifnode] Ability to propose and vote on new jailing parameters.
- [ChainOps] MongoDb resource limit increases.

## 🐛 Bug Fixes

- [UI] State loader for dispensation.
- [UI] Remove the "select all" functionality when clicking in a field.
- [UI] Token with zero balances would sometimes disappear from the swap options.

## ❓ Miscellaneous

- [UI] Integration of the Playwright test framework.
- [UI] Amount API for appropriate decimal placement across all token types.
- [Peggy] Ability to export the Ethbridge keeper data (required for when migrating to Cosmos 0.42.x).
- [ChainOps] Automated pipeline deployment.
