# ADR 002: Ebrelayer Ethereum Subscribe

## Changelog

- 2020/10/21: Initial version

## Status

*Proposed*

## Context
In this ADR, we discuss the solution for ebrelayer how to subscribe the events from Ethereum and process these events.
### Summary

For ebrelayer, it just needs to subscribe to the events from the BridgeBank smart contract, both LogLock event and LogBurn event, then process both events and send transaction to Sifchain. Basically, there are two problems.
1. The block produced in Ethereum maybe be reorganized, most of system opt to confirm the finalization of block after 6 blocks.
2. How to store the events. The way that this is accomplished is anytime we see a new block, we look back 50 blocks, look at all of the events in that block, then package them up as prophecy claims and send them to sifchain.

## Current solution
We start to process the events happened 50 blocks before. 50 blocks can guarantee the finalization of block. Then there is no impact from block reorganization. Whenever a new block comes in, we look back 50 blocks ago, and relay events from that block to sifchain.

### Pros and Cons

Pros: The solution is easy to implement.

Cons: The events lost if ebrelayer restart. We will store the events in persistent storage like local database or message queue system. It depends on the requirement of product in the future.

## Consequences
We will see obvious of delay for Sifchain get the message of cross-chain asset transfer. The end to end tests are impacted, need extra operations and transactions to verify the cases both transfering eth/erc20 asset to Sifchain and burn pegged Cosmos asset back to Sifchain.

### Positive

- We can quickly deliver our MVP

### Negative

- Nothing major

### Neutral

- Nothing major

