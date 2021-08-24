# How to: Upgrade existing Peggy

## Setup

Modify the .env file to include:

MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/...
ROPSTEN_URL=https://eth-ropsten.alchemyapi.io/v2/...
ROPSTEN_PROXY_ADMIN_PRIVATE_KEY=aaaa...
MAINNET_PROXY_ADMIN_PRIVATE_KEY=aaaa...
DEPLOYMENT_NAME="sifchain-testnet-042-ibc"

## Execution

cd to the smart-contracts directory

Run:

    bash scripts/upgrade_contracts.sh sifchain-testnet-042-ibc

Replacing sifchain-testnet-042-ibc with whatever deployment you're upgrading.
