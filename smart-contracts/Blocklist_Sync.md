# Synchronize the Blocklist

If you are trying to synchronize our EVM blocklist with OFAC's blocklist, you will need to use the `yarn blocklist:run` command.

1. Before running the script, edit your .env file adding the following variables:

```
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/ZGe5q0xD06oCAHiwf6ZAexnzGKSPrS5P
MAINNET_PRIVATE_KEY_BLOCKLIST=e67825808c9642d98d16b5794da4582432cb159610ff3934e8a0bac074e725f2
ACTIVE_PRIVATE_KEY=MAINNET_PRIVATE_KEY_BLOCKLIST
```

_Please note that the values of MAINNET_URL and MAINNET_PRIVATE_KEY_BLOCKLIST have been redacted and won't work on the mainnet. You should change them to your actual mainnet Alchemy URL and the Blocklist admin's private key._

Important:

- Make sure MAINNET_PRIVATE_KEY_BLOCKLIST in your .env file is the private key matching the Blocklist admin's address.

- Ensure MAINNET_URL is set correctly.

To synchronize the blocklist, use the following command:

```
yarn blocklist:run
```

## Devnotes

If you just want to test the sync, all you have to do is run the command

```
yarn blocklist:test
```

It will fork the mainnet and simulate the sync.
