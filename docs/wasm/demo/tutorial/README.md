# How to Write a Sifchain Smart Contract

This tutorial shows how to write a smart contract targeted for Sifnode. This is different to writing a regular smart contract since it includes custom messages for the CLP module of Sifnode. However, the tutorial assumes no prior knowledge of writing CosmWasm smart contracts.

## Prerequisites

1. Install Rust toolchain. Follow instructions of the Rust book https://doc.rust-lang.org/book/ch01-01-installation.html

TL;DR:

```
$ curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh
```

2. Install wasm backend for the rust compiler:

```
rustup target add wasm32-unknown-unknown
```

Check output of `rustup target list --installed` includes `wasm32-unknown-unknown`


## Contract setup

1. Use Cargo (build and package managing tool installed alongside the rust compiler) to create a contract

```
cargo new --lib my-contract
```

2. In the same directory add `lib` and `profile-release` sections to the `Config.toml` file:

```toml
[lib]
crate-type = ["cdylib", "rlib"]

[profile.release]
opt-level = 3
debug = false
rpath = false
lto = true
debug-assertions = false
codegen-units = 1
panic = 'abort'
incremental = false
overflow-checks = true
```

## Contract Instantiation

A contract typically has three entry points, instantiate, execute and query. A contract must be instantiated before execute or query can be called. In this section we'll focus on instantiation.

1. Open `src/lib.rs`. If you find test code a distraction then remove the boilerplate test code. Add code to instantiate the contract:

```rust
use cosmwasm_std::{entry_point, DepsMut, Env, MessageInfo, Response};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
pub struct InstantiateMsg {}

#[entry_point]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, SwapperError> {
    Ok(Response::default())
}
```

In this simple example nothing is done (inside the smart contract) during instantiation.

2. We're missing dependencies, which we'll add to the `Cargo.toml` file:

```toml
[dependencies]
cosmwasm-std = { version = "1.0.0-beta" }
serde = { version = "1.0", default-features = false, features = ["derive"] }
```

3. Compile the contract:

```
cargo build --release --target wasm32-unknown-unknown
```

This should create a `target/wasm32-unknown-unknown/release/my_contract.wasm` file

This is 1.7M! We can make it smaller:

```
RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown
```

This takes us down to 129k.

### Deploy Contract

1. Start and run `sifnoded`. From `sifnode` directory `make init` then `make start`

CosmWasm separates contract deployment and instantiation into two distinct steps. Firstly the contract is stored on the blockchain then an instance is created and the instantiate function is called. This allows multiple instances of the same contract to be instantiated without needing to store multiple copies of the same code.

2. Store the contract on the blockchain:

```
sifnoded tx wasm store ./my_contract.wasm \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet  \
-y
```

3. Instantiate the contract:

```
sifnoded tx wasm instantiate 1 '{}' \
--amount 50000rowan \
--label "swapper" \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y
```

We've not included execute or query entry points, so nothing more can be done with this contract.

## Contract Execution

1. Add code to execute a transaction:

```rust
use schemars::JsonSchema;
use sif_std::{SifchainMsg, SifchainQuery};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecuteMsg {
    Swap { amount: u32 },
}

#[entry_point]
pub fn execute(
    _deps: DepsMut<SifchainQuery>,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<SifchainMsg>, String> {
    match msg {
        ExecuteMsg::Swap { amount } => {
            let swap_msg = SifchainMsg::Swap {
                sent_asset: "rowan".to_string(),
                received_asset: "ceth".to_string(),
                sent_amount: amount.to_string(),
                min_received_amount: "0".to_string(),
            };

            Ok(Response::new()
                .add_attribute("action", "swap")
                .add_message(swap_msg))
        }
    }
}

```

We've added the `ExecuteMsg` enum which defines the different messages that must be handled by the execute function. In this example the ExecuteMsg enum has one variant, you can easily add more variants but you must handle each variant in the match statement in the body of the execute function.

A successful response to an execute call returns a `Response` struct. The `add_message` function allows messages to be added to the response. The messages must be of `CosmosMsg` type which is an enum with the following variants:

```rust
pub enum CosmosMsg<T> {
    Bank(BankMsg),
    Custom(T),
    Staking(StakingMsg),
    Distribution(DistributionMsg),
    Stargate {
        type_url: String,
        value: Binary,
    },
    Ibc(IbcMsg),
    Wasm(WasmMsg),
    Gov(GovMsg),
}
```

There are handlers in the node which will handle all the variants, the exception is the custom variant, which has code in the wasmd module to handle the unwrapping but requires additional handlers for the content of the custom message.

The type T can be any structure (provided it implements CustomMsg) however there need to be handlers in the node for the custom types. The sif_std crate defines a SifchainMsg enum which represent Sifchain specific transactions which have corresponding handlers in the golang code:

```rust
pub enum SifchainMsg {
    Swap {
        sent_asset: String,
        received_asset: String,
        sent_amount: String,
        min_received_amount: String,
    },
    AddLiquidity {
        external_asset: String,
        native_asset_amount: String,
        external_asset_amount: String,
    },
    RemoveLiquidity {
        external_asset: String,
        w_basis_points: String,
        asymmetry: String,
    },
}
```
We could have explicitly wrapped the SifchainMsg in the call to add_message:

```rust
add_message(CosmosMsg::Custom(SifchainMsg::Swap{sent_asset: "rowan", received_asset: "ceth", sent_amount: "12", min_received_amount: "3"}))
```

However since the sif_std crate also implements `impl From<SifchainMsg> for CosmosMsg<SifchainMsg>` which defines how a `CosmosMsg` is derived from a `SifchainMsg`, the call to the add_message function can be reduced to:

```rust
add_message(SifchainMsg::Swap{sent_asset: "rowan", received_asset: "ceth", sent_amount: "12", min_received_amount: "3"})
```

2. We're missing dependencies, which we'll add to the `Cargo.toml` file:

```toml
[dependencies]
schemars = "0.8"
sif-std = { path = "../../sif_std/" }
```

3. Compile and deploy this contract as above (re-initialize the chain), create a rowan/ceth liquidity pool, then execute the contract:

```
sifnoded tx clp create-pool \
--symbol ceth \
--nativeAmount 2000000000000000000 \
--externalAmount 2000000000000000000 \
--from sif --keyring-backend test \
--fees 100000000000000000rowan \
--chain-id localnet \
--broadcast-mode block \
-y
```

Execute a swap via the smart contract:

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
  '{"swap":{"amount": 200}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```

## Querying a Contract

1. Add code to perform a query:

```rust
use cosmwasm_std::{to_binary, Deps, QueryResponse, StdResult};
use sif_std::PoolResponse;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    Pool { external_asset: String },
}

#[entry_point]
pub fn query(deps: Deps<SifchainQuery>, _env: Env, msg: QueryMsg) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::Pool { external_asset } => {
            let req = SifchainQuery::Pool { external_asset }.into();
            to_binary(&deps.querier.query::<PoolResponse>(&req)?)
        }
    }
}

```

Similarly to the ExecuteMsg of the previous section, the QueryMsg defines the different messages that must be handled by the query function. As before you can add additional variants provided they're handled in the match statement in the body of the query function.

deps.querier.query() is what performs the query. The type passed to the querier is a QueryRequest enum, which has several variants:

```rust
pub enum QueryRequest<C> {
    Bank(BankQuery),
    Custom(C),
    Staking(StakingQuery),
    Stargate {
        path: String,
        data: Binary,
    },
    Ibc(IbcQuery),
    Wasm(WasmQuery),
}
```

There are handlers in the node which will handle all the variants, the exception is the custom variant, which has code in the wasmd module to handle the unwrapping but requires additional handlers for the content of the custom message.

The type C can be any structure (provided it implements CustomQuery) however there need to be handlers in the node for the custom types. The sif_std crate defines a SifchainQuery enum which represent Sifchain specific queries which have corresponding handlers in the golang code:

We could have explicitly wrapped SifchainQuery in the call to query:

```rust
deps.querier.query::<PoolResponse>(&QueryRequest::Custom(req))
```

However the cosmwasm_std crate implements `impl<C: CustomQuery> From<C> for QueryRequest<C>` which defines how a QueryRequest can be derived for any type, allowing us to write the call as

```rust
deps.querier.query::<PoolResponse>(&req)
```

The other thing to note is that the sif_std crate defines PoolResponse which is the structure that this query is expected to return. This corresponds to the corresponding structure in the golang code.

2. Compile and deploy this contract as above (re-initialize the chain), create a liquidity pool, then run a query via the smart contract:

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
  '{"pool":{"external_asset": "ceth"}}' \
  --chain-id localnet
```
