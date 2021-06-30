# Update whitelisted tokens on mainnet

If you are trying to whitelist many token addresses at once, you will need to use the bulkSetTokenLockBurnLimit.js script.

Before running the following script go to the data folder and open create a file called `whitelist_<network>_<date>.json`, for example `whitelist_ethereum_feb_21_2021.json`. Change the name of the file to remove the date and insert the current date. Copy the contents from `whitelistUpdate.json` into your newly created file and change the addresses to the address you want to whitelist.

Make sure ETHEREUM_PRIVATE_KEY in your .env file is the private key matching the OPERATOR address and ensure INFURA_PROJECT_ID is set correctly. Get the bridgebank address and set it in the env var when running the script. To bulk update the whitelist and add tokens, use bulkSetTokenLockBurnLimit.js like so:
```
ETHEREUM_PRIVATE_KEY=$eth_key_operator BRIDGEBANK_ADDRESS=$bridge_bank_address truffle exec scripts/bulkSetTokenLockBurnLimit.js --network develop ../data/whitelist_<network>_<date>.json
```

## Note
In the previous version of the smart contracts, there was a concept of a max token lock or burn amount. This functionality has been completely removed from the codebase. There is no longer a max lock or burn amount, there is only a token whitelist. Whitelisted tokens may be locked and burned if they are on the whitelist to perform that action.