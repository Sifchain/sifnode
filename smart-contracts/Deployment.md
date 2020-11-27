# How to: Deploy Peggy

1. You will first have to set up the accounts that you want to use as validators on the bridge. Each of these accounts should run a relayer that listens to sifchain for incoming transactions that need to go to ethereum.

2. Currently there are 2 user roles in peggy, more are planned to come. There is an operator role, this person can add and remove validators from the ethereum smart contract so that they can no longer unlock or mint assets on the ethereum peggy smart contracts. When you set up your env file, remove the mnemonic and all local variables. Your readme should now look like this.

```
# ------------
#    General
# ------------

# This number is how much total voting power is needed before a prophecy is completed on the ethereum side of peggy and a user gets their funds released or minted to them.
# Ideally the number shouldn't be over 100, but it can be as long as the sum of INITIAL_VALIDATOR_POWERS is greater than or equal to it
CONSENSUS_THRESHOLD = 100

# This is the address of the operator. This will essentially be the admin and owner of the peggy smart contracts and will perform administrative function for the system.
ETHEREUM_PRIVATE_KEY=["Your raw hex private key without 0x prepended"]

# Replace example INFURA_PROJECT_ID with your Infura project's ID
INFURA_PROJECT_ID='replace with your infura project id'

# Replace example OPERATOR with the address that corresponds to your ETHEREUM_PRIVATE_KEY
OPERATOR = "0xb52c8fb9416611806027a2675c4B748726Ce85Fa"

# Replace example INITIAL_VALIDATOR_ADDRESSES with the desired validator addresses.
# These are the addresses of the users who will validate prophecies and run relayers
INITIAL_VALIDATOR_ADDRESSES = "0x515d8ab15EB94d64b6E2a2878520651BA50d8F7f,0x2fA4F2EB8104af7Dd9A8a4BCa573b6757877F4f8,0x6119c0D7c840038F61E7167b674212A1df5c73E8,0x7B8f616ecf0cE23E0d8564E90c5038a0D8862e58"
# You must have the private keys for these addresses or know that the relayers own the private keys so that they can sign transactions that will be sent to the smart contract

# Replace example INITIAL_VALIDATOR_POWERS with the desired validator powers
INITIAL_VALIDATOR_POWERS = 25,25,25,25

# The array of INITIAL_VALIDATOR_POWERS and INITIAL_VALIDATOR_ADDRESSES should both be the same length otherwise the smart contract deployment will fail.
# The arrays need to correspond at the index you want a user to have a certain power. e.g. 
# INITIAL_VALIDATOR_ADDRESSES = "0x515d8ab15EB94d64b6E2a2878520651BA50d8F7f,0x2fA4F2EB8104af7Dd9A8a4BCa573b6757877F4f8"
# INITIAL_VALIDATOR_POWERS = 50,30
# This configuration would mean 0x515d8ab15EB94d64b6E2a2878520651BA50d8F7f would have 50 power and 0x2fA4F2EB8104af7Dd9A8a4BCa573b6757877F4f8 would have 30 power because the first address is 0x5 and the first power is 50, the same is true for the second address and power.

```

Copy paste the above setup into a .env file inside the testnet-contracts folder and replace with your infura project id, validator powers, consensus threshold, private key for the operator and addresses.

3. Make sure that the gas price in the truffle-config file for the network you are trying to deploy to is priced at the current market price or higher so that your contracts actually get deployed.

4. Now it's time to actually deploy, first run an npm install.

Then copy the file .env.example to the .env file.

Now go to eth gas station or etherscan and find the current gas price on mainnet.
https://ethgasstation.info/
https://etherscan.io/gastracker

Go to your .env file and assign the current gas price in wei to the variable MAINNET_GAS_PRICE

Set the ETHEREUM_PRIVATE_KEY to your private key you want to be the operator for the smart contract.

Ensure that OPERATOR is the address that corresponds to the ETHEREUM_PRIVATE_KEY.

Ensure that OWNER is the address that will be the admin for the bridge bank and can wire eRowan in.

Ensure that INITIAL_VALIDATOR_ADDRESSES and INITIAL_VALIDATOR_POWERS are set correctly.

then run the following command:
```
truffle deploy --network mainnet
```
You can replace mainnet with ropsten or local, whichever network you would like to deploy to.

5. After you have deployed the contracts to your network of choice, you will need to run this command from the smart-contracts folder:
```
DIRECTORY_NAME="your_deployment_name_here" node scripts/saveContracts.js
```
Save the deployment folder and the subdirectory you just made and all of its files into git so that other users can interact with the smart contracts as well and know where the addresses are located. The output files will be located in the deployment directory inside the subdirectory of the name you passed. The fallback directory within the deployment directory is the default directory which is where things automatically get saved if you run this script without a directory name.

6. Grab the eRowan token address on whatever network you are on. Then run the setup_eRowan.js file to properly hook eRowan into the contracts. Make sure that the EROWAN_ADDRESS variable in the .env file is set to the eRowan token address. Make sure that the OWNER address is set properly in the env file so that you have an owner for the bridgebank contract that can use the admin api. Then run the command from the testnet-contracts folder:
```
truffle exec scripts/setup_eRowan.js --network mainnet
```

You can replace the network with whatever network you would like, ropsten, develop or mainnet.

When running this command, make sure that the private key in the .env file corresponds to the operator of the bridgebank contract as well as the admin and deployer of the eRowan token contract.

7. Make sure to record the smart contract addresses and share them with the team so that you can interact with them in the future. Then git add the deployment folder, commit it and push it to your own branch, then merge it with develop so that we have the records documented of the smart contract address in the version control history.
