#!/bin/sh

set -e

# usage:
#
#    scripts/update_contracts.sh $DEPLOYMENT_NAME $network
#
# For example,
#
#    scripts/update_contracts.sh sifchain-testnet-042-ibc ropsten
#
# must run this from the smart-contracts directory, and must update .env with
# the appropriate values for:

# MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/...
# ROPSTEN_URL=https://eth-ropsten.alchemyapi.io/v2/...
# ROPSTEN_PROXY_ADMIN_PRIVATE_KEY=aaaa...
# DEPLOYMENT_NAME="sifchain-testnet-042-ibc"

deploymentDir=deployments/$1/
rm -f .openzeppelin/
ln -s $deploymentDir/.openzeppelin .openzeppelin
npx hardhat run scripts/upgrade_contracts.ts --network $2
git commit -m "update deployment" $deploymentDir
