# How to: Add IBC ERC20 tokens

## Setup

Modify the .env file to include:

```
MAINNET_URL=<mainnet_url>
MAINNET_PROXY_ADMIN_PRIVATE_KEY=<private_key>
```

Where:

|Item|Description|
|----|-----------|
|`<mainnet_url>`|Replace with the Infura Mainnet URL|
|`<private_key>`|Replace with the ETH Private Key for the BridgeBank operator|


## Execution

    cd smart-contracts
    npm install
    hardhat run scripts/create_ibc_matching_token.ts | grep -v 'No need to generate' > test_data/ibc_token_addresses.jsonl  
    hardhat run scripts/attach_ibc_matching_token.ts < 
