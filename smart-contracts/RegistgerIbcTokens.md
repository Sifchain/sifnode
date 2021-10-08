# How to register IbcTokens / BridgeTokens

After having deployed IbcTokens, we have to register them in BridgeBank.  
If you have been asked to run this command, your contact will have given you
a filename that holds all the information needed by the system.

Add it to your .env like so:

```
REGISTER_TOKENS_SOURCE_FILENAME=./data/deployed_ibc_tokens_08_Oct_2021.json
```

_Please note that you should change the filename accordingly._

You will also need a mainnet Private Key and a mainnet URL. Please add them to tour .env like so:

```
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/XXXXXXXXXXXXXXXXXXXXXXXX
MAINNET_PRIVATE_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

Now, run the following command:

```
yarn registerIbcTokens:run
```

The final results will be added to your source file (in this example, "./data/deployed_ibc_tokens_08_Oct_2021.json").

This file should be pushed to git. If you cannot do that, please send it to the Peggy team.
