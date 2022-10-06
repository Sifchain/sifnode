## Devnotes

The whitelist command (`yarn whitelist:run`) will run two scripts sequentially.

The first script (fetchTokenDetails.js) will fetch metadata from each token address in the initial address list. It will try to fetch name, symbol, decimals and imageUrl from each token.

- Symbols that contain spaces or special characters will be rejected and that token will NOT be added to the whitelist.
- A new file will be created in the data folder. Its name will be something like "whitelist_mainnet_update_14_sep_2021.json", but with today's date.
- If the script fails to fetch imageUrl for a token, it will set imageUrl to `null` in that token data. You may edit it manually later.

The second script (bulk_set_whitelist.ts) will communicate with the BridgeBank, adding all tokens to the whitelist.

After both scripts are completed, you will se the message "~~~ DONE ~~~" in your terminal/console. Then, you should verify whether all tokens have been successfully added to the whitelist. All you need to do is read the logs that the scripts have generated, directly in your terminal/console.

## Testing with a mainnet fork

If you want to test the whitelisting flow, add this variable to your .env:

```
USE_FORKING=1
```

And run the following command:

```
yarn whitelist:test
```

This is also useful if you just want to generate the json file without actually updating the whitelist in production.

## Note

In the previous version of the smart contracts, there was a concept of a max token lock or burn amount. This functionality has been completely removed from the codebase. There is no longer a max lock or burn amount, there is only a token whitelist. Whitelisted tokens may be locked and burned if they are on the whitelist to perform that action.
