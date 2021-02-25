# scripts description
The document describes the scripts in smart contract. 

At first, you need set up the execution environment variable by `cp .env.example .env`, replace the ETHEREUM_PRIVATE_KEY and INFURA_PROJECT_ID. Then you can call `yarn develop` to start ganache-cli to start local Ethereum node.

For each script, you can call `yarn XXX`, the XXX list can be found in the package.json scripts. For the prefix with integrationtest, just used in the integration test. You can send transaction to different Ethereum network via parameter `--network ropsten`, default is local network. All scripts for integration test not explained in the document, they are similar to scripts listed below, just applied to the integration environment, for some special account, or load test and so on.

## usage for each script

### advance 
Call scripts/advanceBlock.js to generate some blocks in local test, default is to advance 5 blocks. It is very useful since for Ethereum we can confirm the transaction after some blocks. For ebrelayer, just process the events 50 blocks ago.

Usage:

yarn advance --network ropsten [blocks]

### peggy:address
Call scripts/getBridgeRegistryAddress.js to get the address of bridgeRegistry contract. Other two important contracts BridgeBank and CosmosBridge's addresses are stored in bridgeRegistry.

Usage:

yarn peggy:address --network ropsten

### peggy:validators
Call scripts/getValidators.js to get all validators address and its power. The script try to check if accounts derived from MNEMONIC is active validator and its power. 

Usage:

yarn peggy:validators --network ropsten

### peggy:hasLocked 
Call scripts/hasLockedTokens.js to get the ERC20 contract address according to the locked token's symbol. 

Usage:

yarn peggy:hasLocked --network ropsten [token_symbol]

### peggy:getTx 
scripts/getTxReceipt.js to get the receipt of a transaction.

Usage:

yarn peggy:getTx [tx_hash]

### peggy:setup
Call scripts/setOracleAndBridgeBank.js to set the oracle and bridge bank smart contracts' addresses in bridgeRegistry. You must deploy thr oracle and bridge bank contracts before setup.

Usage:

BRIDGEBANK_ADDRESS=[address] COSMOS_BRIDGE_ADDRESS=[address] yarn peggy:setup

### peggy:lock
scripts/sendLockTx.js to lock some tokens to smart contract.

Usage:

BRIDGEBANK_ADDRESS=[address] yarn peggy:lock --network ropsten [cosmos_address] [ERC20_address] [amount]

### peggy:whiteList
Call scripts/sendUpdateWhiteList.js to update whitelist. Set the ERC20 contract as true or false in the whitelist. Only enabled ERC20 contract can be used for lock/burn.

Usage:

yarn peggy:whiteList --network ropsten [ERC20-address] [true/false]

### peggy:burn
Call scripts/sendBurnTx.js to burn some pegged token from Cosmos, then token will go to account of cosmos-address in the transaction.

Usage:

yarn peggy:burn --network ropsten [cosmos-address] [ERC20-address] [amount]

### peggy:check
scripts/sendCheckProphecy.js to check prophecy power and its threshold according to prophecy ID.

Usage:

yarn peggy:check --network ropsten [prophecy-ID]

### peggy:getTokenBalance
Call scripts/getTokenBalance.js to get the balance of an account on an ERC20 smart contract. If ERC20-address not provided, will query the balance of eth.

Usage:

yarn peggy:getTokenBalance --network ropsten [account] [ERC20-address]

### token:address
Call scripts/getTokenContractAddress.js to get the BridgeToken contract address. BridgeToken is an ERC20 contract, it is used for testing.

Usage:

yarn token:address --network ropsten

### token:mint
Call scripts/mintTestTokens.js to mint 1^20 token for deployed BridgeToken contract. BridgeToken is an ERC20 contract, it is used for testing.

Usage:

yarn token:mint --network ropsten

### token:approve
Call scripts/sendApproveTx.js to approve some amount of token to BridgeBank contract via send transaction to an ERC20 smart contract, then BridgeBank can call transferFrom method to lock tokens. If the ERC20 address not provide, will approve eth.

Usage:

yarn token:approve --network ropsten [amount] [ERC20-address]

