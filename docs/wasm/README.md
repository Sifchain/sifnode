# Sifchain - Wasm 

Sifnode imports the `x/wasm` module directly from the official CosmWasm 
[repo](https://github.com/CosmWasm/wasmd/tree/master/x/wasm). This module 
provides a sandboxed environment for the execution of smart-contracts with 
guardrails to prevent rogue programs from undermining the system. CosmWasm ships
with native bindings for standard SDK modules like `bank`, `staking`, or `ibc`, 
that enable developers to extend the capabilities of a blockchain without 
requiring the coordination of a full version upgrade. Along with the `wasm` 
virtual-machine, the module automatically adds a set of CLI sub-commands to 
deploy and interract with smart-contract. Additionnaly, sifnode has implemented
custom bindings to its own `clp` module to enable smart-contracts to call into
the functionnality of liquidity pools (swap, add/remove liquidity, etc). The
purpose of this document is to explain this code so as to enable other 
developers to maintain or extend those bindings. 

To understand how the custom bindings work, it is necessary to have a high-level
understanding of the CosmWasm execution model. CosmWasm has chosen to implement
an [Actor Model](https://docs.cosmwasm.com/docs/1.0/architecture/actor/), 
whereby separate components (contract-contract, contract-module) communicate 
with each other via message passing, instead of directly referencing each other.
This is instrumental in preventing re-entrency bugs that plague Ethereum. When a
contract wants to call into another contract or module, it emits a message that 
gets picked up by the wasm handler, and routed to the appropriate actor. All 
state that is carried over between one call and the next happens in storage and 
not in memory.

Smart-contracts are usually written in Rust and compile down to WASM byte-code.
The `wasm` module, and in particular the handler that dispatches messages from
one component to another (smart-contract or module), is implmented in `go`.
So how are messages exchanged between the WASM environment and the `go` 
environment? The answer is that messages are serialized in JSON, and encoded/
decoded on both sides. The `wasm` handler (in go) decodes incoming messages, and
tries to match them to a set of known types. Incoming messages are then 
forwarded to the appropriate component. 

On the `go` side, the set of known types are: 
```
	Bank         *BankMsg         
	Custom       json.RawMessage  
	Distribution *DistributionMsg 
	Gov          *GovMsg          
	IBC          *IBCMsg          
	Staking      *StakingMsg      
	Stargate     *StargateMsg     
	Wasm         *WasmMsg         
```
It is fairly obvious from the names, which component corresponds to each type of
message. Note the existence of the `Custom` Type; this is where we can extend 
the handler. 

When a raw message is received by the `wasm` module, it is first matched against
a set of `Encoders` that are registered with the module's `Keeper` upon 
initialization. These encoders, try to convert the raw JSON encoded messages 
into a known message type. Note that the handling of a message stops at the 
first encoder that successfully matches the message to a known type. The encoder
then does some processing of the message and forwards it onto the appropriate 
component. The `wasm` module obviously ships with encoders for the known types 
above, but we can add custom encoders to handle our own `clp` related messages.

The `sifnode/wasm` package contains our custom types and encoders. The encoders
are registered with the `wasm.Keeper` in `app.go`. The semantics are slighlty 
different for read-only queries and for messages that are meant to update the 
state, but the idea is roughly the same; we define custom query types and 
register custom handlers with the `wasm.Keeper`.

```
type SifchainMsg struct {
	Swap            *Swap            `json:"swap,omitempty"`
	AddLiquidity    *AddLiquidity    `json:"add_liquidity,omitempty"`
	RemoveLiquidity *RemoveLiquidity `json:"remove_liquidity,omitempty"`
}
```

```
type SifchainQuery struct {
	Pool *PoolQuery `json:"pool,omitempty"`
}
```

When our custom encoders succesfully match an incoming message (query), they 
`encode` into an `sdk.Msg` (`sdk.Query`) and forward it to the same backbone 
router that handles regular Cosmos SDK messages (queries). The SDK router 
matches the incoming SDK message with a registered handler, and executes the 
handler. 

It would be technically possible to reference the `clp.Keeper` directly in our 
custom `wasm` encoders to manipulate the `clp` state directly from there, but
it is preferable to plug back into the main routing system. This enables us to
keep all clp code in one place and to rely on the security model implemented by
the SDK. 

On the Rust side, we have developed a reusable library that mirrors the custom
types mentionned above. It can be imported into a smart-contract to create and
emmit `clp` messages and queries. Please refer to the `clp` tutorial for some
guidance on how to do this.