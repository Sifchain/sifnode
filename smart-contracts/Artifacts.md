# This is a document to describe how we save our smart contract deployments to a central location in git as a team.

## Problem

Truffle has a directory called build that gets generated whenever contracts are compiled. This directory also keeps track of the networks and the addresses contracts are deployed to on various networks. Unfortunately, truffle only allows you to save 1 smart contract address per network. This presents a problem as we still need to track the addresses of the same contract that has deployed to a single testnet multiple times for different smart contract instances.

## Solution

In order to remedy this problem, we have created a script called ```saveContracts.js``` located inside of the scripts folder. Before you do a deployment to the network of your choosing, delete the build folder, then run your migrations. After the contracts are successfully deployed, then from the smart-contracts directory, run the saveContracts.js script. Be sure to pass a name as an environment variable when you run the script of the new folder where these artififacts and addresses will be deployed. 

## Example:
```
DIRECTORY_NAME="your_deployment_name_here" node scripts/saveContracts.js
```

After you have run the script and specified a new directory name, commit that new folder in the deployments folder into git and include it in your Pull Request so that the entire team can understand where the smart contract addresses are.

### A Word on our Smart Contracts
We are using an upgradeable smart contract pattern with a delegate proxy pattern. This means that we are essentially deploying two contracts for every single contract in our repository. One contract is the storage and proxy contract, the other contract is the logic contract where the proxy contract delegate calls to in order to have its storage changed.

## Upgradeability
Because we are using the openzeppelin framework, after you have deployed to your network of choice, the .openzeppelin folder will have all the data you need. When you run the above script, saveContracts.js, the .openzepplin folder will be copied into the sub directory of deployments. This way, if you ever want to upgrade those contract instances, you will have the addresses to do so.