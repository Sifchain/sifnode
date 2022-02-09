# Sifchain - Wasm 

This folder contains code that demonstrates how to bind wasm contracts to 
custom SDK modules. 

## Reflect

The `reflect` package contains the `reflect.wasm` smart-contract as well as the
Go bindings that enable other sifnode code to exchange messages with this 
contract.

The `reflect` contract essentially forwards incoming messages to the `wasm`
module's handler on the `go` side. 

Internally, the `wasm` module's handler matches each incoming messages to one of 
the following types: 

    - bank
    - staking
    - distribution
    - stargate
    - IBC 
    - gov
    - wasm
    - custom

and forwards the message to the appropriate module. In the case of a `custom`
message, it tries to match the message against the registered custom decoders,
and  **this is where we plug in our custom logic to process custom messages**.

When we create the `wasm` module's keeper, we pass it our custom decoder as an
option. 

```go
	wasmOpts = append(wasmOpts,
		wasmkeeper.WithGasRegister(NewJunoWasmGasRegister()),
		// the reflect options are added for testing only
		wasmkeeper.WithMessageEncoders(reflect.ReflectEncoders(codec)),
		wasmkeeper.WithQueryPlugins(reflect.ReflectPlugins()),
	)
```

Our custom decoder decodes an incoming CustomMsg (from json format) to our 
`ReflectCustomMsg` type. This type contains some fields that we can use to 
create a `clp` message. The `clp` message is returned by our decoder and further
 relayed by the `wasm` module's handler.

At the moment, we have hardcoded the `clp` message that gets created when a
`ReflectCustomMsg` is received, and we don't use the fields in the 
`ReflectCustomMsg`. Hence, the next step is to encode meaningful information in
`ReflectCustomMsg` and use it to populate a `clp` message.

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