# How to: Add IBC ERC20 tokens

## Setup

For a mainnet deployment, modify the .env file to include:

```
MAINNET_URL=<mainnet_url>
MAINNET_PRIVATE_KEY=<private_key>
DEPLOYMENT_NAME=<deployment_name>
TOKEN_FILE=data/ibc_mainnet_tokens.json
TOKEN_ADDRESS_FILE=data/ibc_token_addresses.jsonl
```

Where:

| Item                   | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| `<mainnet_url>`        | Replace with the Infura Mainnet URL                          |
| `<private_key>`        | Replace with the ETH Private Key for the BridgeBank operator |
| `<deployment_name>`    | Replace with the deployment name like sifchain               |
| `<token_file>`         | File with information on new tokens to be deployed           |
| `<token_address_file>` | File where new tokens that are created will be written to    |

# Overview

There are two distinct steps to this script.

One is creating the new bridge tokens. 
Two is attaching the bridge tokens to the bridgebank by calling addExistingBridgeToken on the bridgebank.

The script to attach bridge tokens will be run by a user with priviledged access to the bridgebank with the owner role.

## Mainnet Token Deployment

This does not require any special permissions.  Tokens will be handed off to another user
before they're attached to the BridgeBank.

    cd smart-contracts
    npm install
    npx hardhat run scripts/create_ibc_matching_token.ts --network mainnet | grep -v 'No need to generate' | tee data/ibc_token_addresses.jsonl

## Steps to Attach Tokens to BridgeBank

This requires the private key of the BridgeBank owner.

    cd smart-contracts
    npm install
    npx hardhat run scripts/attach_ibc_matching_token.ts --network mainnet < data/ibc_token_addresses.jsonl

## Testing with forked mainnnet

Since you're running two scripts, you'll need a hardhat node running (otherwise the first script will run, execute transactions, then throw them away)

Start a hardhat node in a shell:

    npx hardhat node --verbose

Then run the two scripts in a different shell:

    npx hardhat run scripts/create_ibc_matching_token.ts --network localhost | grep -v 'No need to generate' | tee test_data/ibc_token_addresses.jsonl
    npx hardhat run scripts/attach_ibc_matching_token.ts --network localhost < test_data/ibc_token_addresses.jsonl

or for ropsten

    npx hardhat run scripts/create_ibc_matching_token.ts --network ropsten | grep -v 'No need to generate' | tee test_data/ibc_token_addresses.jsonl
    npx hardhat run scripts/attach_ibc_matching_token.ts --network ropsten < test_data/ibc_token_addresses.jsonl

# Updating dev and test deployments

## Update symbol_translator.json

Modify https://github.com/Sifchain/chainOps/blob/main/.github/workflows/variableMapping/ebrelayer.yaml
with the new symbol_translation.json entries.  It's a yaml file, so you need to escape json - you could
use something like ```jq -aRs``` to do the quoting.

To push that data out and restart the relayers, use https://github.com/Sifchain/chainOps/actions/workflows/peggy-ebrelayer-deployment.yml
