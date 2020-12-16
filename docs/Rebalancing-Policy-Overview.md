# **Sifchain Rebalancing Policy: Overview**

## Changelog



## Context
This document is intended to provide high level context on Sifchain's `Dynamic Rebalancing Policy` for its economic system.

## Economic System
Sifchain's existing economic system is comprised of two subsystems: the `Validator Subsystem` where validators and delegators participate in staking in order to earn revenue and the `Liquidity Pool Subsystem` where liquidity providers lock assets to earn revenue. 

Revenue in the `Validator Subsystem` is comprised of inflationary block rewards and revenue in the `Liquidity Pool Subsystem` is comprised of swap fees. In order for the overall system to remain healthy the revenues earned through these two subsystems must remain balanced and equivalent. If one subsystem is more profitable than the other users will be incentivized to participate in that subsystem, upsetting the balance between the two and compromising the security of the network.

## Rebalancing Policy
In order to maintain subsystem balance we introduce dynamic controls on the two revenue streams: inflationary block rewards and swap fees. By throttling these streams dynamically we can correct any deviations from revenue balance. The system we will employ to facilitate this dynamic control is the `Rebalancing Policy`.

**Rebalancing Policy Flow:** Every block the following flow will be executed.

- The `Rebalancing Policy` takes as inputs the total supply of Rowan and the supplies of Rowan currently observed in each subsystem.

- Each subsystem's supply is then compared to the total supply to obtain a ratio. 

- Each ratio is then compared with predetermined target ratios to obtain an error value. These target ratios can be set through an external governance process. 

- From the error values we then calculate each subsystem's control parameter. The control parameter is essentially the throttle described above. 

- The subsystems then retrieve their control parameters and incorporate them into their revenue calculations (inflationary block rewards and swap fees). The control parameters are in the range of 0 to 1 so the result of incorporating them will either be a decrease to the unaltered revenue or will maintain the unaltered revenue. The control parameters will never boost revenues beyond their unaltered state.

**Changes Needed to Implement This Policy:**

In practice the `Rebalancing Policy` will live in the new `Rebalancer Module`. This module will be responsible for executing the flow described above. Each of our subsystems will also need to be altered to bring in and utilize the new control parameters. This means changes to Cosmos's mint module, where inflationary block rewards are calculated, as well as changes to the CLP module, where swap fees are calculated.

This system is also extensible, meaning that further subsystems can be incorporated as desired.

## Links

Blockscience's math specifications and implementation guides can be found here: https://hackmd.io/@shrutiappiah/By450P3ID