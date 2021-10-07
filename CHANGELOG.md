 # v0.9.10
> October 8, 2021
>
## ⭐ Features

- [Sifnode] Support for transferring high precision tokens over IBC via denom conversion.

# v0.9.7
> September 17, 2021
>
## ⭐ Features

- [Sifnode] Enabling of inflationary block rewards.

# v0.9.5
> August 24, 2021
>
## ⭐ Features

- [Sifnode] IBC Support is now available.

# v0.9.1
> August 6, 2021
>
## ⭐ Features

- [Sifnode] Numerous updates to allow for the enabling of IBC in the coming weeks.

# v0.9.0
> July 28, 2021
>
## ⭐ Features

- [Sifnode] Upgraded to version 0.42 of the Cosmos SDK.

# v0.8.8
> July 16, 2021
>
## 🐛 Bug Fixes

- [Sifnode] Updates to the Oracle module to allow for exporting state.

----

# v0.8.7
> June 9, 2021
>
## 🛠 Improvements

- [UI] Introduction of a footer! User's can sign-up for our newsletter, and link to our privacy policy, roadmap, and legal disclaimer.
- [UI] Removal of negative signs in the Pool Stats>Arbitrage as these are confusing and not necessary.
- [UI] Removal of the words Peg and Unpeg. We have updated our language across our entire application to use 'Import' and 'Export' instead. We did this to be more clear with the action that is being done, as well as to prepare for more chains being integrated into Sifchain.
- [UI] Included the ability for a user to see their net gain/loss on their liquidity pool position. This includes earnings from swap fees AND any gains or losses associated with changes in the tokens' prices. This number is represented as USDT.
- [Peggy] Relayer Upgrade - Implementation of Retry Logic.
- [Sifnode] Claims module - The claims module is done and ready! This will allow users to be able to submit a claim for their liquidity mining & validator subsidy rewards.
- [UI] Users can now claim their Liquidity Mining and Validator Subsidy Rewards in the DEX! Feel free to navigate to the 'rewards' tab, see details on your rewards and claim them if desired. We recommend that you keep your liqudity, stake, and rewards untouched until you reach your full maturation date. To read more about this, please reference our documentation [here](https://docs.sifchain.finance/resources/rewards-programs#liquidity-mining-and-validator-subsidy-rewards-on-sifchain)

## 🐛 Bug Fixes

- [Sifnode] Additional updates/fixes to the dispensation module (used for airdrops).
- [UI] When a user had 0 balances, the sorting logic we were using in our token listings was not accurate.

----

# v0.8.4
> May 13, 2021

## 🐛 Bug Fixes

- [Sifnode] Fixes to the dispensation module (used for airdrops).

----

# v0.8.2
> April 29, 2021

## 🛠 Improvements

- [UI] New design elements for the DEX (header/typorgraphy/buttons).

----

# v0.8.1
> April 28, 2021

## ❓ Miscellaneous

- [Peggy] Removed the previously added retry logic, for when infura gives us a not found error.

----

# v0.8.0
> April 27, 2021

## ⭐ Features

- [UI] Prevent pegging when not connected to a supported EVM network.
- [UI] Arbitrage Opportunity - Now color-coded to indicate in which direction the opportunity is present.

## 🛠 Improvements

- [UI] Added logic for when clicking on max button in ROWAN to take into consideration necessary gas fees.
- [UI] In swap confirmation screen, built in cleaner UX logic around the way we display numbers.
- [UI] Included logic in the token list pop-ups when doing a swap or liquidity add for how we sort the displayed tokens. We are also now calling in user's balances directly in this pop-up as well for easy viewing.
- [Peggy] Add retry logic when infura gives us a not found error. Add additional retry logic to try to retrieve the tx if it cannot be found on the first query.
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
