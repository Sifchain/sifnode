# Sifchain - Wasm 

This folder contains code used to demonstrate how to bind wasm contracts to 
custom SDK modules. 

## Reflect

The `reflect` package contains the `reflect.wasm` smart-contract as well as the
Go bindings that enable other sifnode code to exchange messages with this 
contract.

The `reflect` contract essentially forwards incoming Cosmos messages and queries
to the SDK's underlying routing mechanism.

In the demo below, we send a wrapped `swap` message to the smart-contract which 
relays it onto the `clp` module via the SDK's message passing system.

Most of the custom `go` code relates to encoding and decoding messages.

We define a `ReflectCustomMsg` (which corresponds to the output of the 
`reflect` contract) and add an encoder to the wasm keeper that enables the 
wasm keeper to recognize these messages, unpack their contents and convert them 
into SDK messages so that they can be picked up by other modules. A similar 
pattern is implemented for queries.

## Demo



### Setup

First, initialize a local node from the `sifnode` root directory:

1. Initialize the local chain: `make init`

2. Start the chain: `make run`

The rest of the commands are to be executed from the same directory as this
`README` file.

### Create rowan/ceth liquidity pool

1. Create pool: `make create-pool`

2. Check pool: `make show-pool`

### Store and Initialize

1. Store and initialize `reflect` contract: `make deploy-contract`

2. Check contract balance: `make show-contract-balance`

### Swap rowan for ceth

1. Swap: `make swap`

2. Check balances again: `make show-contract-balance`

3. Check pool: `make show-pool`