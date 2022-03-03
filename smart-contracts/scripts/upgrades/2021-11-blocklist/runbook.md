# BridgeBank Upgrade: Blocklist

If you are trying to upgrade the BridgeBank contract so that it knows how to use OFAC's Blocklist, you will need to use the upgrade script that was specifically created for that.

1. Before running the script, edit your .env file adding the following variables:

```
MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/your-alchemy-id
MAINNET_PRIVATE_KEYS=XXXXXXXX,YYYYYYYY,ZZZZZZZZ
ACTIVE_PRIVATE_KEY=MAINNET_PRIVATE_KEYS
```

Where:

```
XXXXXXXX is the PROXY ADMIN private key
YYYYYYYY is BridgeBank's and CosmosBridge's OPERATOR private key
ZZZZZZZZ is the PAUSER's private key
```

They should be separated by a comma, and they have to be in that order (admin first, operator second, pauser third).  
Please also make sure you changed `your-alchemy-id` for your actual Alchemy id in `MAINNET_URL`.

To run the script, use the following command:

```
npx hardhat run scripts/upgrades/2021-11-blocklist/run.js
```

#

## Devnotes

You may want to run the script in test mode before executing it in production. To do that, simply run the command

```
USE_FORKING=1 npx hardhat run scripts/upgrades/2021-11-blocklist/run.js
```

It will fork the mainnet and impersonate the account of the proxy admin to do the upgrade, than it will impersonate the BridgeBank's pauser and operator to register the blocklist in it. You will be able to see all the steps as if you were executing it on the mainnet. If there are no errors, it means it's safe to execute the script in production.

Also, you might want to run the proper test that was created to make sure everything will work after the upgrade. To do that, use

```
USE_FORKING=1 npx hardhat run scripts/upgrades/2021-11-blocklist/test.js
```

It will do everything, impersonating the correct accounts, and then it will simulate prophecies and check whether they are fulfilled or blocked as expected.

#

## Next Steps

After executing the upgrade, we should sync the blocklist. To do that, please consult Blocklist_Sync.md.
