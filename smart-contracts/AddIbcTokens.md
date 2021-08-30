# How to: Add IBC ERC20 tokens

## Setup

Modify the .env file to include:

```
MAINNET_URL=<mainnet_url>
MAINNET_PROXY_ADMIN_PRIVATE_KEY=<private_key>
DEPLOYMENT_NAME=<deployment_name>
```

Where:

|Item|Description|
|----|-----------|
|`<mainnet_url>`|Replace with the Infura Mainnet URL|
|`<private_key>`|Replace with the ETH Private Key for the BridgeBank operator|
|`<deployment_name>`|Replace with the deployment name like sifchain |
|`<token_file>`|Replace with the path to the token file containing json|

## Execution

    cd smart-contracts
    npm install
    hardhat run scripts/create_ibc_matching_token.ts --network mainnet | grep -v 'No need to generate' > test_data/ibc_token_addresses.jsonl  
    hardhat run scripts/attach_ibc_matching_token.ts --network mainnet < test_data/ibc_token_addresses.jsonl 
 
