# CLP Withdrawal Scenario

The purpose of this demo is to go through various attempts at withdrawing 
liquidity from a liquidity pool. We will see what happens when we attempt to 
withdraw liquidity before it is fully unlocked. We will see how to send an 
unocking request, and what happens if we try to withdraw before the unlocking
period expires. We will also see how to update the reward parameters.

## Setup

First, intialize a local node from the `sifnode` root directory:

1. Initialize the local chain: `make init`
2. start the chain: `make run`
3. Change the working directory `cd test/scenarios/clp`

## Create rowan/ceth liquidity pool

1. Create pool: `make create-pool`
2. Check pool: `make show-pool`

Notice how the native balance increases during the course of this test. This is the rewards program at work.

## Show CLP parameters

1. Show params: `make show-params`

```
params:
  default_multiplier: "0.000000000000000000"
  liquidity_removal_cancel_period: "518400"
  liquidity_removal_lock_period: "120960"
  reward_periods:
  - allocation: "10000000000000000000000"
    end_block: "120960"
    id: RP_1
    multipliers:
    - asset: ceth
      multiplier: "1.500000000000000000"
    start_block: "1"
```

There is a non-zero multiplier associated with the `ceth` pool as part of a 
reward program between blocks 1 and 120960

## Locked Liquidity

1. Remove Liquidity (first attempt)

Try to remove half of our liquidity from the `ceth` pool: 
`make remove-liquidity`

This fails because we need to unlock/unbond liquidity before it can be removed

```
failed to execute message; message index: 0: user does not have enough balance
  of the required coin
```

2. Check our status as a liquidity provider: `make show-lp`

```
./scripts/show_lp.sh
external_asset_balance: "2000000000000000000"
height: "718"
liquidity_provider:
  asset:
    symbol: ceth
  liquidity_provider_address: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
  liquidity_provider_units: "2000000000000000000"
  unlocks: []
native_asset_balance: "50693783068783068219"
```

note that we have 2*10^18 units in the pool (100%)

2. Unlock Liquidity

Let's try to unlock half of our liquidity units: `make unlock-liquidity`

3. Check Registered Unlocks: `make show-lp`

```
external_asset_balance: "2000000000000000000"
height: "801"
liquidity_provider:
  asset:
    symbol: ceth
  liquidity_provider_address: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
  liquidity_provider_units: "2000000000000000000"
  unlocks:
  - request_height: "793"
    units: "1000000000000000000"
native_asset_balance: "57555555555555554912"
```

This time we can see that we have an unlock request for 1*10^18 units

4. Remove Liquidity (second attempt): `make remove-liquidity`

```
failed to execute message; message index: 0: user does not have enough balance
  of the required coin
```

It still doesn't work because the locking period has not passed.

As we saw with `make show-params`, the locking period is 120960 blocks, so we 
have to wait 120960 blocks (which roughly corresponds to 7 days on BetaNet), 
before being able to withdraw our unlocked liquidity.

## Change Locking Period

Let us try to reduce the locking period so that we don't have to wait until 
block 121753. We do this by submitting a special transaction from the admin 
account:

1. Submit parameter change transaction: `make change-params`

2. Check params again: `make show-params`

You will notice that `liquidity_removal_cancel_period: "720"` and `liquidity_removal_lock_period: "10"`

## Withdraw Liquidity

Now we only have to wait 10 blocks before being able to remove liquidity.

1. Try removing liquidity again: `make remove-liquidity`

2. Show status of liquidity provider: `make show-lp`

```
external_asset_balance: "1000000000000000000"
height: "92"
liquidity_provider:
  asset:
    symbol: ceth
  liquidity_provider_address: sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
  liquidity_provider_units: "1000000000000000000"
  unlocks:
  - request_height: "64"
    units: "0"
native_asset_balance: "3273478835978835953"
```

Note that our liquidity units have dropped from 2*10^18 to 1*10^18 and the
unlock request is fully consumed.
