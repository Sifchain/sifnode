# Upgrade: Peggy 2.0

If you are trying to upgrade the system from Peggy 1.0 to Peggy 2.0, you will need to use the upgrade script that was specifically created for that.

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
npx hardhat run scripts/upgrades/peggy2-2022-02/run.js
```

The script will:

1. Pause the system
2. Upgrade BridgeBank
3. Upgrade CosmosBridge
4. Check if all storage slots are safe
5. Unpause the system

#

## Devnotes

You may want to run the script in test mode before executing it in production. To do that, simply run the command

```
USE_FORKING=1 npx hardhat run scripts/upgrades/peggy2-2022-02/run.js
```

It will fork the mainnet and impersonate the relevant accounts to do the upgrade. You will be able to see all the steps as if you were executing it on the mainnet. If there are no errors, it means it's safe to execute the script in production.
