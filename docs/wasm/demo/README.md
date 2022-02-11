# Demo

In this demo we deploy the `swapper` smart-contract which calls the `clp` module
to do a swap. First we manually create a rowan/ceth pool using the `sifnoded` 
cli commands. Then we deploy and instantiate the smart-contract, giving it 50000
rowan. Finally, we call the smart-contract with a specified `amount` to trigger 
a swap from rowan to ceth. 

Note that it is the contract's bank balances that are updated, not the balances 
of the signer of the transaction.

## Setup

First, initialize a local node from the `sifnode` root directory:

1. Initialize the local chain: `make init`

2. Start the chain: `make run`

The rest of the commands are to be executed from the same directory as this
`README` file.

## Create rowan/ceth liquidity pool

1. Create pool: `make create-pool`

2. Check pool: `make show-pool`

## Store and Initialize

1. Store and initialize the `swapper` contract: `make deploy-contract`

2. Check contract balance: `make show-contract-balance`

## Swap rowan for ceth

1. Swap: `make swap`

2. Check balances again: `make show-contract-balance`

3. Check pool: `make show-pool`

### Add liquidity to ceth pool

1. Add liquidity: `make add-liquidity`

2. Check balances again: `make show-contract-balance`

3. Check pool: `make show-pool`

### Remove liquidity from ceth pool

1. Remove liquidity: `make remove-liquidity`

2. Check balances again: `make show-contract-balance`

3. Check pool: `make show-pool`