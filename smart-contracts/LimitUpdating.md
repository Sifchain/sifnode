# Update whitelisted tokens on mainnet

If you are trying to whitelist many token addresses at once, you will need to use the `yarn whitelist:run` script.

1) Before runnig the script, go to the data folder and create a file called address_list_source.json, or edit it so
that it has only the addresses you want to whitelist. The contents of the file should have a list of addresses, like so:
[
 "0x217ddead61a42369a266f1fb754eb5d3ebadc88a",
 "0x9e32b13ce7f2e80a01932b42553652e053d6ed8e"
]

2) Now, edit you .env file adding the following variables:
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/ZGe5q0xD06oCAHiwf6ZAexnzGKSPrS5P
MAINNET_PRIVATE_KEY_OPERATOR=e67825808c9642d98d16b5794da4582432cb159610ff3934e8a0bac074e725f2
ACTIVE_PRIVATE_KEY=MAINNET_PRIVATE_KEY_OPERATOR
BRIDGEBANK_ADDRESS=0x6CfD69783E3fFb44CBaaFF7F509a4fcF0d8e2835
DEPLOYMENT_NAME=sifchain-1
ADDRESS_LIST_SOURCE="data/address_list_source.json"

Expected usage: first, set the .env variable ADDRESS_LIST_SOURCE with the path for
 * the file that contains a list of addresses in the following format:
 * 
 * 
 * 
 * EXAMPLE (.env):
 * ADDRESS_LIST_SOURCE="data/testAddressList.json"

Before running the following script go to the data folder and create a file called `whitelist_<network>_<date>.json`, for example `whitelist_ethereum_feb_21_2021.json`. Change the name of the file to remove the date and insert the current date. Copy the contents from `whitelist_mainnet_update_postibc.json` into your newly created file and change the addresses to the addresses you want to whitelist.  
Example:

```
{
  "array": [
    {
      "address": "0xa47c8bf37f92abed4a126bda807a7b7498661acd"
    },
    {
      "address": "0x853d955acef822db058eb8505911ed77f175b99e"
    }
  ]
}
```

Your .env file should have the following variables.  
All values here are mocked examples that won't work on the mainnet.  
Please change them to match your needs.

```
BRIDGEBANK_ADDRESS=0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8
WHITELIST_DATA=data/whitelist_mainnet_update_14_sep_2021.json
DEPLOYMENT_NAME=sifchain-1
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/ZGe5q0xD06oCAHiwf6ZAexnzGKSPrS5P
MAINNET_PRIVATE_KEY_OPERATOR=c8750aa1c067bbde78beb793e8fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
ACTIVE_PRIVATE_KEY=MAINNET_PRIVATE_KEY_OPERATOR
```

Make sure MAINNET_PRIVATE_KEY_OPERATOR in your .env file is the private key matching the OPERATOR address and ensure MAINNET_URL is set correctly.  
Make sure ACTIVE_PRIVATE_KEY equals MAINNET_PRIVATE_KEY_OPERATOR in your .env file (exactly as in the above example).  
Get the bridgebank address and set it in the env var when running the script.  
To bulk update the whitelist and add tokens, use bulk_set_whitelist.ts like so:

```
npx hardhat run scripts/bulk_set_whitelist.ts --network mainnet
```

## Note

In the previous version of the smart contracts, there was a concept of a max token lock or burn amount. This functionality has been completely removed from the codebase. There is no longer a max lock or burn amount, there is only a token whitelist. Whitelisted tokens may be locked and burned if they are on the whitelist to perform that action.
