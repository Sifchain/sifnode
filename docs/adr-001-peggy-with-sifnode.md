# ADR 001: Peggy integration with Sifnode

## Changelog

- 2020/10/21: Initial version

## Status

*Proposed*

## Context
In this ADR, the two possible solutions for Sifchain are listed, and the pros and cons are compared with each other. It explain the reason behind the current implementation and what's the ideal architecture would be like.
### Summary

In Sifchain, there are two major programmes, Peggy and Sifnode. Peggy works as the bridge to other blockchain system like Ethereum, Bitcoin, EOS, Polkadot and so on, transfer assets to Cosmos also for the opposite direction. Sifnode is based on Cosmos SDK, its major functionality is to provide the liquidity and token swap. By combine two componets again, it is possible to provide the liquility with transfered ETH, BTC, EOS and so on. Swap is the same, then Sifchain user can swap the RWN (Sifchain native token) with transfered ETH, BTC, EOS and so on.

From architecture point of view, there are two solutions for Sifchain.
1. Both Peggy and Sifnode have their own ledger, they communicate and transfer value via IBC. It is an ideal solution considering the flexibility and scaling out. But the IBC, at the time write the ADR, is not mature enough for development. 

2. Peggy and Sifnode co-exist in the same ledger, they share the accout and balance. The solution couple the Peggy Sifnode, but it is easier to implement for now. If consider the IBC's availability, it is maybe the only solution to deliver.

### Pros and Cons

1. seperate chain solution

Pros: Peggy and Sifnode can develop and extend seperately, totally decoupled. Peggy will connect more blockchain system like EOS, BTS, Polkadot, ETC and so on, even the other Cosmos based chain instead of Sifchain, with enable IBC. Peggy focus on recording the cross chain assets transfer. Design the its own incentive algorithm, consensus strategy and native token.

Cons: Peggy and Sifnode need IBC support, which not used in production environment yet. For customer, they need two transactions to provide liquidity. At first, transfer asset to Peggy. Then transfer asset from Peggy to Sifnode via IBC.

2. shared ledger solution

Pros: It is much easier to deploy and maintenence since all operations like cross-transfer, add liquidity and swap are processed by single chain. No dependency on the service of IBC.

Cons: For the long term, the system is hard to scale out. For example, all node validators for sifnode must deploy the Ethereum node and pay for high gas fee.
It increase the cose of validators.
## Decision
We choose the second solution to implement now. The major reason is the IBC still in development, not mature for production environment usage for the time we write the ARD. We will keep our eyes on the maturity of IBC, give our judgement when it is could be trid and even be switched.

In our development, We need decouple the cross-chain functions (as peggy) and liquidilty/swap at the module level. It will avoid too much efforts to split them if IBC is availale.

## Consequences

### Positive

- We can quickly deliver our MVP

### Negative

- Not yet

### Neutral

- Not yet

## References

