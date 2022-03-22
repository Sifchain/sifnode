# denom mapping from peggy1.0 to peggy2.0
We have different structure for denom in peggy1.0 and peggy2.0. To upgrade the sifnoded, we need migrate the denom. At first, we need to get a complete map to list what's the denom in peggy1.0 and its counterpart in peggy2.0

## mapping algorithm
case 1: for rowan token, it is sifnode native toke, not changed
case 2: ibc token, it is not changed
case 3: etheruem imported token, the denom is the 'c' + contract's symbol in peggy1.0. its denom will be "sifBridge{network_descriptor:04d}{token_contract_address.lower()}"

## how to get the mapping step by step

### get the all denom in production environment
In sifnoded, we have a sub-command to get all denom from tokenregistry x module

```
sifnoded query tokenregistry entries
```

### get the all smart contract addresses from Ethereum network
In develop branch, we have a script to get all whitelisted token. We can run following command to get all contracts and their denom
```
yarn integrationtest:whitelistedTokens --json_path /Users/junius/github/sifnode/smart-contracts/deployments/sifchain-1  --ethereum_network mainnet --bridgebank_address 0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8 --network mainnet
```
The result file is put the same folder as denom_contracts.json

### run denom_mapping.py to get the denom mapping
```
python3 denom_mapping.py
```
the result files is denom_mapping_peggy1_to_peggy2.json and denom_mapping_peggy2_to_peggy1.json.
denom_mapping_peggy1_to_peggy2.json: key is denom in peggy1.0, value is denom in peggy2.0.
denom_mapping_peggy2_to_peggy1.json: it is reverse to first json file.

The script also prints out all denom not mapped, the main reason is the contract not found in the Ethereum.
