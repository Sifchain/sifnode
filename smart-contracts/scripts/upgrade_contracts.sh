#!/bin/sh

set -e

# usage: 
# 
#    scripts/update_contracts.sh $DEPLOYMENT_NAME
#
# 
# must run this from the smart-contracts directory, and must update .env with
# the appropriate values for:

# MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/...
# ROPSTEN_URL=https://eth-ropsten.alchemyapi.io/v2/...
# ROPSTEN_PROXY_ADMIN_PRIVATE_KEY=aaaa...
# DEPLOYMENT_NAME="sifchain-testnet-042-ibc"

deploymentDir=deployments/$1/
rm -f .openzeppelin/
ln -s $deploymentDir .
npx hardhat run scripts/upgrade_contracts.ts
git commit -m "update deployment" $deploymentDir
