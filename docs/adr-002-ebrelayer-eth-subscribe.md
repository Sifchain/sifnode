# ADR 002: Ebrelayer Ethereum Subscribe

## Changelog

- 2020/10/21: Initial version

## Status

*Proposed*

## Context
In this ADR, we discuss the solution for ebrelayer how to subscribe the events from Ethereum and process these events.
### Summary

For ebrelayer, it just needs to subscribe to the events from the BridgeBank smart contract, both LogLock event and LogBurn event, then process both events and send transaction to Sifchain. Basically, there are two problems.
1. The block produced in Ethereum maybe be reorganized, most of system opt to confirm the finalization of block after 6 blocks. So it is better to process the event with some dalay. For the events once we receive from subscription, we store them in the buffer and wait for more blocks produced.
2. How to store the events, there are two options. First one is store them in memory, but events will be lost if ebrelayer restarted. Second solution is store in local db, message queue like Kafka, but will increase the complexity of ebrelayer's deployment.

## Current solution
We start to process the events happened 50 blocks before. 50 blocks can guarantee the finalization of block. Then there is no impact from block reorganization. Repeated events and invalid events are abondoned in events buffer. We choose memory to store events for now, it is simple and easy to implements. The ebrelayer will miss some events if ebrelayer restarted. Considering the whole Sifchain system is a decentralized, we can tolerate some ebrelayer offline and not send transaction to Sifnode.  

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

