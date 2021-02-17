# Peggy Admin Smart Contract Threshold

In CosmosBridge contract, the validators can send newProphecyClaim transaction to the contract. And each validator has its voting power. Threshold is important argument, represents how many validators need send newProphecyClaim transactions before prophecy claim be finalized.

Sometimes, we need update the threshold. But the operation may lead to side effect, like someone might not receive a payout. So we should puase the contract first, then update the threshold.
