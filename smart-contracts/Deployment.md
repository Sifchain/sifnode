# How to: Deploy Peggy

## Setup
Before we start this guide, you must have the following dependencies installed on your system:
- node v16.6.1
- yarn 1.22.10
- latest version of truffle

cd into the smart-contracts directory and run ```yarn install```. Once those dependencies have been installed you can move onto the next step.

## Deployments

1. You will first have to set up the accounts that you want to use as validators on the bridge. Each of these accounts should run a relayer that listens to sifchain for incoming transactions that need to go to ethereum.

2. Currently there are 2 user roles in peggy, more are planned to come. There is an operator role, this person can add and remove validators from the ethereum smart contract so that they can no longer unlock or mint assets on the ethereum peggy smart contracts. When you set up your env file, remove the mnemonic and all local variables. Your readme should now look like this.

```

# ------------
#    General
# ------------
# This number is how much total voting power is needed before a prophecy is completed on the ethereum side of peggy and a user gets their funds released or minted to them.

CONSENSUS_THRESHOLD=75
# This is the address of the owner of the upgradeable proxy
ETHEREUM_PRIVATE_KEY="c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3"

# ------------
#   Network
# ------------
# Owner of the bridgebank, this must be set when we are deploying to tesnet or mainnet
OWNER='0x627306090abaB3A6e1400e9345bC60c78a8BEf57'
# Address of the pauser. This is the user that can pause locking, burning, minting and unlocking from the bridgebank. This does not need to be a multisig address that owns this as it would be cumbersome to have to rely on a multisig for this.
PAUSER='0x627306090abaB3A6e1400e9345bC60c78a8BEf57'
# Replace example INFURA_PROJECT_ID with your Infura project's ID
INFURA_PROJECT_ID="JFSH7439sjsdtqTM23Dz"
# Replace example OPERATOR with the intended address
# The operator is allowed to tell the cosmos bridge contract where the bridgebank contract is.
# Additionally, the operator is allowed to update the relayer addresses and their associated powers. This means they can add or delete relayers from the whitelist in cosmos bridge.
# In the bridgebank, the operator is allowed to update the whitelist and limits of each token. This means that this operator should not be a multisig, at least for the bridgebank as this would be incredible cumbersome to have to sign a tx multiple times before a limit or whitelist action goes through.
OPERATOR='0x627306090abaB3A6e1400e9345bC60c78a8BEf57'
# Replace example INITIAL_VALIDATOR_ADDRESSES with the desired validator addresses
# You must have the private keys for these addresses or know that the relayers own the private keys so that they can sign transactions that will be sent to the smart contract
INITIAL_VALIDATOR_ADDRESSES = "0x515d8ab15EB94d64b6E2a2878520651BA50d8F7f,0x2fA4F2EB8104af7Dd9A8a4BCa573b6757877F4f8,0x6119c0D7c840038F61E7167b674212A1df5c73E8,0x7B8f616ecf0cE23E0d8564E90c5038a0D8862e58"
# Replace example INITIAL_VALIDATOR_POWERS with the desired validator powers
INITIAL_VALIDATOR_POWERS = 25,25,25,25
# On the mainnet, set the price of gas based on the current gas prices
# This is not needed locally
MAINNET_GAS_PRICE=10000000
# If you are deploying a new instance of peggy that needs to use eRowan,
# set this variable to the address of the eRowan smart contract
EROWAN_ADDRESS='0x0d8cc4b8d15D4c3eF1d70af0071376fb26B5669b'
```

Copy paste the above setup into a .env file inside the smart-contracts folder and replace with your infura project id, validator powers, consensus threshold, private key for the user that will be deploying, and operator, owner and pauser addresses. The private key you use will be the owner of the proxy admin, so be sure to hold onto that private key as it is the only one authorized to do a transaction to upgrade the peggy bridge smart contracts.

3. Make sure that the gas price in the truffle-config file for the network you are trying to deploy to is priced at the current market price or higher so that your contracts actually get deployed.

4. Now it's time to actually deploy, first run a yarn install.

Now go to eth gas station or etherscan and find the current gas price on mainnet.
https://ethgasstation.info/
https://etherscan.io/gastracker

Go to your .env file and assign the current gas price in wei to the variable MAINNET_GAS_PRICE

Set the ETHEREUM_PRIVATE_KEY to your private key you want to be the upgradeable proxy contract owner for the smart contract.

Ensure that OWNER is the address that will be the user who deployed the erowan bridgetoken contract so that they can wire eRowan into the bridgebank easily.

Ensure that INITIAL_VALIDATOR_ADDRESSES and INITIAL_VALIDATOR_POWERS are set correctly. For mainnet, this should be 4 validator addresses and each validator power should be 25 with the consensus threshold being 75 so that it only takes 3/4 validators to sign off on a transaction for funds to be released.

then run the following commands:
```
truffle migrate --network mainnet -f 1 --to 1
truffle migrate --network mainnet -f 2 --to 2
truffle migrate --network mainnet -f 3 --to 3
truffle migrate --network mainnet -f 4 --to 4
```
You can replace mainnet with ropsten or local, whichever network you would like to deploy to. If you are deploying to local or testnet, you can instead run:
```
truffle migrate --network <ropsten or develop>
```

4.5 Now you will need to manually set the bridgebank address on the cosmos bridge by calling setBridgeBank as the operator to get the smart contract fully wired up to talk with each other.
Do this by running:
```
BRIDGEBANK_ADDRESS='insert bridgebank address' COSMOS_BRIDGE_ADDRESS='insert cosmosbridge address' truffle exec scripts/setBridgeBank.js --network mainnet
```

5. After you have deployed the contracts to your network of choice and run the setBridgeBank script, you will need to run this command from the smart-contracts folder:
```
DIRECTORY_NAME="your_deployment_name_here" node scripts/saveContracts.js
```
Save the deployment folder and the subdirectory you just made and all of its files into git so that other users can interact with the smart contracts as well and know where the addresses are located. The output files will be located in the deployment directory inside the subdirectory of the name you passed. The fallback directory within the deployment directory is the default directory which is where things automatically get saved if you run this script without a directory name.

6. Grab the eRowan token address on whatever network you are on. Then run the setup_eRowan.js file to properly hook eRowan into the contracts. Make sure that the EROWAN_ADDRESS variable in the .env file is set to the eRowan token address. Make sure that the OWNER address is set properly in the env file so that you have an owner for the bridgebank contract that can use the admin api. Then run the command from the smart-contracts folder:
```
truffle exec scripts/setup_eRowan.js --network mainnet
```

You can replace the network with whatever network you would like, ropsten, develop or mainnet.

When running this command, make sure that the private key in the .env file corresponds to the operator of the bridgebank contract as well as the admin and deployer of the eRowan token contract.

7. Make sure to record the smart contract addresses and share them with the team so that you can interact with them in the future. Then git add the deployment folder, commit it and push it to your own branch, then merge it with develop so that we have the records documented of the smart contract address in the version control history. If you are deploying to mainnet, make sure to add the build folder to git and commit it to develop.
