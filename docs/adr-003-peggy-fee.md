# ADR 003: Peggy fee

## Changelog

- 2021/01/01: Initial version

## Status

*Proposed*

## Context
In this ADR, we discuss the solution for sifnode how to charge the cross chain transaction fee in Ethereum network.
### Summary

Thre are two cases of the cross chain assets transfer from Sifchain to Ethereum, lock Sifchain native token and burn peggy Ethereum token. The ebrelayers need send transaction of Prophecy claim to Ethereum after receive such events. Then there are some cost of transaction fee in Ethereum. The total cost depends on the number of ebrelayers and transaction fee for each Prophecy claim call. There are lots of different solutions to cover the cost and charge the Sifnode account.

## Current solution

In the first implementation, the Sifnode account need pay some Ceth (pegged Eth in Sifnode) for cross chain assets transfer from Sifchain to Ethereum. And we assume there are three permissioned ebrelayer nodes. For this feature, we extend the both lock and burn messages with two more arguments, first one is the Ceth paid for the transaction, second argument is the message type which could be submit, revert and retrun. 
1. Submit is the same transaction as before, Sifnode account send transaction for lock or burn.
2. Revert is a new transaction, the ebrelayer will send transaction to Sifnode if paid Ceth is no enough.
3. Return is other new transaction type, for ebrelayer to gas fee left after send Prophecy claim to Ethereum   

### Pros and Cons

Pros: The solution is easy to implement. 

Cons:
1. The sifnode account must hold some pegged Eth, then can send lock and burn transaction.
2. Current implementation assume there are three ebrelayer nodes, so we can easily to judge if the paid Ceth enough. But in real network, the number of ebrelayer nodes are more than three and variable. It will be difficult to estimate the total cost.
3. The solution depends on the permissioned ebrelayer nodes. We don't have any mechanism to avoid malicious ebrelayer yet.

## Consequences
After the implementation deployed, the interaction between Sifnode and ebrelayer will be more complex. For wallet, explorer and Peggy UI, need track the balance change from extended transactions.

### Positive

- We can quickly deliver our MVP

### Negative

- Nothing major

### Neutral

- Nothing major

