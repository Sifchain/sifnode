# How to deploy IbcTokens / BridgeTokens

This script will deploy N new BridgeTokens to an EVM network.

Before executing this script, add the following variables to your .env, changing the values to your actual mainnet Alchemy URL and Private Key:

```
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/XXXXXXXXXXXXXXXXXXXXXXXX
MAINNET_PRIVATE_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
ACTIVE_PRIVATE_KEY=MAINNET_PRIVATE_KEY
```

Then, create or edit the file data/ibc_tokens_to_deploy.json so that it has only the IbcTokens that you want to deploy.  
Example:

```
[
 {
   "name": "Alice Token",
   "symbol": "ALI",
   "decimals": 10,
   "denom": ""
 },
 {
   "name": "Bob Token",
   "symbol": "BOB",
   "decimals": 18,
   "denom": "Bob denom"
 }
]
```

Note that the `denom` field is optional. If you don't have that information, you may leave it as an empty string.

Finally, run the command `yarn deployIbcTokens:run`.

A new file will be created with the results. It's name will be something like data/deployed_ibc_tokens_07_Oct_2021.json, but with today's date.
