# Update whitelisted tokens on mainnet

If you are trying to whitelist many token addresses at once, you will need to use the `yarn whitelist:run` command.

1) Before running the script, go to the data folder and create a file called address_list_source.json, or edit it so
that it has only the addresses that you want to whitelist.  
The contents of the file should have a list of addresses, like so:  
```
[
 "0x217ddead61a42369a266f1fb754eb5d3ebadc88a",
 "0x9e32b13ce7f2e80a01932b42553652e053d6ed8e"
]
```

2) Now, edit you .env file adding the following variables:
```
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/ZGe5q0xD06oCAHiwf6ZAexnzGKSPrS5P
MAINNET_PRIVATE_KEY_OPERATOR=e67825808c9642d98d16b5794da4582432cb159610ff3934e8a0bac074e725f2
ACTIVE_PRIVATE_KEY=MAINNET_PRIVATE_KEY_OPERATOR
BRIDGEBANK_ADDRESS=0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8
DEPLOYMENT_NAME=sifchain-1
ADDRESS_LIST_SOURCE="data/address_list_source.json"
```
_Please note that the values of MAINNET_URL and MAINNET_PRIVATE_KEY_OPERATOR have been redacted and won't work on the mainnet. You should change them to your actual mainnet Alchemy URL and the BridgeBank OPERATOR's private key. You may use the other values exactly as they are above._

Important:
- Make sure MAINNET_PRIVATE_KEY_OPERATOR in your .env file is the private key matching the OPERATOR address.

- Ensure MAINNET_URL is set correctly.  

- Make sure ACTIVE_PRIVATE_KEY equals MAINNET_PRIVATE_KEY_OPERATOR (exactly as in the above example).  

- Ensure the BridgeBank address is set correctly.

To bulk update the whitelist and add tokens, use `yarn whitelist:run` like so:

```
yarn whitelist:run
```

## More details and next steps
The command above will run two scripts sequentially.  
The first script (fetchTokenDetails.js) will fetch metadata from each token address in the initial address list. It will try to fetch name, symbol, decimals and imageUrl from each token.  
- Symbols that contain spaces or special characters will be rejected and that token will NOT be added to the whitelist.
- A new file will be created in the data folder. Its name will be something like "whitelist_mainnet_update_14_sep_2021.json", but with today's date.  
- If the script fails to fetch imageUrl for a token, it will set imageUrl to `null` in that token data. You may edit it manually later.

The second script (bulk_set_whitelist.ts) will communicate with the BridgeBank, adding all tokens to the whitelist.  

After both scripts are completed, you will se the message "~~~ DONE ~~~" in your terminal/console. Then, you should verify whether all tokens have been successfully added to the whitelist. All you need to do is read the logs that the scripts have generated, directly in your terminal/console.  

The last step is to create a new UI PR with the newly added tokens. Verify the generated file in the data folder and make sure all tokens have an imageUrl property with an URL assigned to it. If any token doesn't, you'll need to manually find out that token's icon URL and add it there.
Finally, copy all tokens from that file and add them to this file:  
`https://github.com/Sifchain/sifchain-ui/blob/develop/ui/core/src/config/networks/ethereum/assets.ethereum.mainnet.json`  
(it's in a different repo). Open a PR there and you're done.

## Testing with a mainnet fork
If you want to test the whitelisting flow, add this variable to your .env:
```
USE_FORKING=1
```

And run the following command:
```
yarn whitelist:test
```

## Note

In the previous version of the smart contracts, there was a concept of a max token lock or burn amount. This functionality has been completely removed from the codebase. There is no longer a max lock or burn amount, there is only a token whitelist. Whitelisted tokens may be locked and burned if they are on the whitelist to perform that action.