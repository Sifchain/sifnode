# Min commission

The min commission-rate feature is intended to prevent validators charging a commission below
a hard-coded threshold, `minCommission` (currently 5%).
There are three elements to enforcing the min commission:

1. Ensuring that validators cannot be created with a `commission-rate` below 5%. This is achieved by blocking
`MsgCreateValidator` messages with a `commission-rate` < `minCommission`.
2. Ensuring that validators cannot edit their configuration to set the `commission-rate` to less than 5%. This is achieved by blocking
`MsgEditValidator` messages with a `commission-rate` < `minCommission`.
3. Ensuring that all current validators have their `commission-rate` set to at least 5%. This is
achieved with an upgrade handler on the release which introduces the min commission feature.

NOTE: There is no blocking of `MsgDelegate` and `MsgBeginRedelegate` messages to prevent attempts to
delegate/redelegate to validator's with a `commission-rate` below `minCommision` since such validators
cannot exist due to the enforcement steps outlined above.

This tutorial demonstrates these three components working on a localnet.

## Prep

1. Initialize the chain then start a node

```
make init
make run
```

## Blocking `MsgCreateValidator` messages

This section demonstrates the blocking behaviour on `MsgCreateValidator` messages.

### Success

If the `commission-rate` is set above 5% creating a validator succeeds.

1. Query the list of validators and confirm there is currently one validator `sif_val`

```
sifnoded query staking validators
```

2. Create a `akasha_val` validator:

```
sifnoded tx staking create-validator \
  --amount=92000000000000000000000stake \
  --pubkey='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"+uo5x4+nFiCBt2MuhVwT5XeMfj6ttkjY/JC6WyHb+rE="}' \
  --moniker="akasha_val" \
  --chain-id=localnet \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.1" \
  --min-self-delegation="1000000" \
  --from=akasha \
  --keyring-backend=test \
  --broadcast-mode block \
  -y
```

NOTE: The pubkey used in this example was found by running the `sifnoded tendermint show-validator` command
on a different localnet instance.

NOTE: Each validator must have a unique public key.

2. Query the list of validators and confirm that `akasha_val` has been added and there are now two validators:

```
sifnoded query staking validators
```

### Failure

Attempting to create a validator with a `commission-rate` below 5% fails (even if the max rate > 5%).

1. Attempt to create an a validator but with a 3% commission rate (this would have succeeded if the commission-rate had been set > 5%):

```
sifnoded tx staking create-validator \
  --amount=92000000000000000000000stake \
  --pubkey='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"/7LUsFhIdP0jj36wToOwY3zWC75YXxVd1vxp7YAc1Gs="}' \
  --moniker="alice_fail_val" \
  --chain-id=localnet \
  --commission-rate="0.03" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --from=alice \
  --keyring-backend=test \
  --broadcast-mode block \
  -y
```

Which fails with the message `validator commission 0.030000000000000000 cannot be lower than minimum of 0.050000000000000000: invalid request`

## Blocking `MsgEditValidator` messages

This section demonstrates the blocking behaviour on `MsgEditValidator` messages.

### Success

If editing the `commission-rate` to a value above 5% (by less than commission-max-change-rate) the edit succeeds.

1. The commission rate can only be updated once within 24hrs, so wait 24 hrs then edit the `akasha_val` validator to set the commission-rate to 7%. NOTE: there doesn't seem to be a way to shorten the wait here, the 24hr check is hardcoded into the sdk, see https://github.com/cosmos/cosmos-sdk/blob/3f8596c1955e40ef30e4abcd06f2237d132db3a9/x/staking/types/commission.go#L85:

```
sifnoded tx staking edit-validator \
  --from=akasha \
  --commission-rate="0.07" \
  --chain-id=localnet \
  --keyring-backend=test \
  --broadcast-mode block \
  -y
```

2. Query the validators and observe the `akasha_val` validator's `commission-rate` is now 7%

```
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="akasha_val").commission.commission_rates.rate'
```

### Failure

1. Attempt to set the commission to 3%:

```
sifnoded tx staking edit-validator \
  --from=akasha \
  --commission-rate="0.03" \
  --chain-id=localnet \
  --keyring-backend=test \
  --broadcast-mode block \
  -y
```

Which fails with the message `validator commission 0.030000000000000000 cannot be lower than minimum of 0.050000000000000000: invalid request`


## Min commission upgrade handler

Demonstrating the upgrade handler, will require running the previous release then upgrading to the
release which introduces the min-commission feature. It also requires that there's a validator
pre the upgrade with a `commission-rate` below `minCommission`.

1. Checkout the previous release

```
git checkout v0.14.0
```

2. Update the init script to set validator `commission-rate` to 3% and `commission-max` to 4%

NOTE: The commission-rates of the validator could be permanently set to 3% and 4% in the `scripts/init.sh` script. This
however causes the chain to fail to start on newer versions of the code which have the new min-commission feature.

```
sed -i 's/sifnoded gentx.*/sifnoded gentx sif 1000000000000000000000000stake --chain-id=localnet --keyring-backend=test --commission-max-rate=0.04 --commission-rate=0.03/g' scripts/init.sh
```

3. Initialize the chain

```
make init
```

4. Decrease the governance voting period time before first start:

```
echo "$(jq '.app_state.gov.voting_params.voting_period = "60s"' $HOME/.sifnoded/config/genesis.json)" > $HOME/.sifnoded/config/genesis.json
```

5. Start the chain

```
make run
```

6. Query the commission and observe that the rate is 3% and the max is 4%

```
sifnoded query staking validators --output=json | jq .validators[0].commission.commission_rates
```

7. Raise an upgrade proposal:

```
sifnoded tx gov submit-proposal software-upgrade 0.15.0 \
  --from sif \
  --deposit 10000000000000000000stake \
  --upgrade-height 30 \
  --upgrade-info '{"binaries":{"linux/amd64":"url_with_checksum"}}' \
  --title test_release \
  --description "Test Release" \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  --fees 100000000000000000rowan \
  -y
```

8. Vote on proposal:

```
sifnoded tx gov vote 1 yes --from sif --chain-id localnet --keyring-backend test -y --broadcast-mode block
```

The node will have a consensus failure when it reaches the `upgrade-height` set in the upgrade proposal.
Hit `ctrl-c` to kill the stuck node.

9. Checkout the release with the min-commission upgrade handler:

```
git checkout . # to drop the changes made to the init script
git checkout v0.15.0
```

10. Run the release:

```
make run
```

11. Query the commission and observe that the rate is 5% and the max is 5%:

```
sifnoded query staking validators --output=json | jq .validators[0].commission.commission_rates
```

Repeating the above steps but with `commission-rate` initiated to 10% and `commission-max` initiated to 20% will show no
change after the upgrade handler has run. 

# Max voting power

The max voting power feature is intended to prevent delegations or redelegations which would result in a validator having
more than a hard coded threshold voting power (currently 6.6%). This is done by blocking `MsgDelegate` and
`MsgBeginRedelegate` messages which would give the targeted validator more than 6.6% of the voting power. The projected
voting power is defined as the amount of projected token delegated to the validator divided by the projected total amount of **delegated** token.

The SDK and mintscan calculate the voting power as zero for validators outside the validator set and for validators inside the validator set the voting power is the amount of token delegated to the validator divided by
the total amount of **bonded** token (bonded tokens are the subset of delegated tokens which are delegated to validators
in the validator set (i.e. those validators which vote on blocks)) - see https://github.com/cosmos/cosmos-sdk/blob/d0043914ba7c37c3a0d7039d2c2a2aca6b38a93b/x/staking/types/validator.go#L350-L356 and https://www.mintscan.io/sifchain/validators - the cumulative share of the validators in
the validator set (115 validators) add up to 100%, so the tokens delegated to validators outside the validator set (some do exist) are not
included in the calculation.

The reason for choosing to use delegated tokens rather than bonded tokens is that it significantly simplifies the calculation.
Calculating the projected number of bonded tokens means calculating which validators will be in the validator set after the
delegation/redelegation. This would require replicating (and testing) the logic inside the staking module for determining changes to the validator set.

Most tokens are delegated to validators which are bonded, so the real world difference between using bonded vs delegated
tokens is negligible. Given this and the complexity of calculating the projected amount of bonded tokens, the delegated amount of tokens
is used to calculate the projected voting power.

This tutorial demonstrates the max voting power restriction in action.

## Delegate

1. Initialize and start the chain

```
make init
make run
```

2. Confirm that there's one validator, `sif_val`, with 1000000000000000000000000 tokens:

```
sifnoded query staking validators
```

3. Create a second validator with 62000000000000000000000 tokens, this will give it ~5.838% of the voting power:

```
sifnoded tx staking create-validator \
  --amount=62000000000000000000000stake \
  --pubkey='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"+uo5x4+nFiCBt2MuhVwT5XeMfj6ttkjY/JC6WyHb+rE="}' \
  --moniker="akasha_val" \
  --chain-id=localnet \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.1" \
  --min-self-delegation="1000000" \
  --from=akasha \
  --keyring-backend=test \
  --broadcast-mode block \
  -y
```

4. There should now be two validators:

```
sifnoded query staking validators
```

### Failure - Current (and projected) voting power too big

1. Try to delegate 100 tokens to `sif_val`. This would give `sif_val` a projected voting power
of (1000000000000000000000000 + 100) / (1000000000000000000000000 + 62000000000000000000000 + 100) = 0.94161. NOTE:
the exact voting power here might vary as `sif_val` earns rewards and `akasha_val` gets slashed:

```
sifnoded tx staking delegate sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna 100stake \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block -y
```

Which fails with the message:

`This validator has a voting power of 94.161958568738229800%. Delegations not allowed to a validator whose post-delegation voting power is more than 6.600000000000000000%. Please delegate to a validator with less bonded tokens: invalid request`

### Failure - Projected voting power too big

1. Try to delegate 100000000000000000000000 tokens to `akasha_val`. This would give `akasha_val` a projected voting power
of (92000000000000000000000 + 100000000000000000000000) / (92000000000000000000000 + 1000000000000000000000000 + 100000000000000000000000) = 0.1394. NOTE:
the exact voting power here might vary as `sif_val` earns rewards and `akasha_val` gets slashed::

```
sifnoded tx staking delegate sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy 100000000000000000000000stake \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block -y
```

Which fails with the message:

`This validator has a voting power of 13.941480206540447500%. Delegations not allowed to a validator whose post-delegation voting power is more than 6.600000000000000000%. Please delegate to a validator with less bonded tokens: invalid request`

### Success

1. Confirm that `akasha_val` has 62000000000000000000000 tokens. NOTE: this may
vary if `akasha_val` has already been slashed:


```
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="akasha_val").tokens'
```

2. Delegate to `akasha_val`

```
sifnoded tx staking delegate sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy 100stake \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block -y
```

3. Confirm that `akasha_val` now has 62000000000000000000100 tokens. NOTE: this may
vary if `akasha_val` has already been slashed:


```
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="akasha_val").tokens'
```

### Redelegate

1. Initialize and start the chain

```
make init
make run
```

2. Confirm that there's one validator, `sif_val`, with 1000000000000000000000000 tokens:

```
sifnoded query staking validators
```

3. Create a second validator with 62000000000000000000000 tokens, this will give it ~5.838% of the voting power:

```
sifnoded tx staking create-validator \
  --amount=62000000000000000000000stake \
  --pubkey='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"+uo5x4+nFiCBt2MuhVwT5XeMfj6ttkjY/JC6WyHb+rE="}' \
  --moniker="akasha_val" \
  --chain-id=localnet \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.1" \
  --min-self-delegation="1000000" \
  --from=akasha \
  --keyring-backend=test \
  --broadcast-mode block \
  -y
```

4. There should now be two validators:

```
sifnoded query staking validators
```

### Failure - Current (and projected) voting power too big

1. Try to redelegate 100 tokens from `akasha_val` to `to sif_val`. This would give `sif_val` a projected voting power of 
(1000000000000000000000000 + 100) / (1000000000000000000000000 + 62000000000000000000000) = 0.9416. NOTE:
the exact voting power here might vary as `sif_val` earns rewards and `akasha_val` gets slashed:

```
sifnoded tx staking redelegate sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna 100stake \
  --from akasha \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block -y
```

Which fails with the message:

`This validator has a voting power of 94.161958568738229800%. Delegations not allowed to a validator whose post-delegation voting power is more than 6.600000000000000000%. Please redelegate to a validator with less bonded tokens: invalid request`

### Failure - Projected voting power too big

1. Try to redelegate 100000000000000000000000 tokens to `akasha_val` from `sif_val`. This would give `akasha_val` a projected voting power
of (92000000000000000000000 + 100000000000000000000000) / (92000000000000000000000 + 1000000000000000000000000) = 0.1525. NOTE:
the exact voting power here might vary as `sif_val` earns rewards and `akasha_val` gets slashed:

```
sifnoded tx staking redelegate sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy 100000000000000000000000stake \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block -y
```

Which fails with the message:

`This validator has a voting power of 15.254237288135593200%. Delegations not allowed to a validator whose post-delegation voting power is more than 6.600000000000000000%. Please redelegate to a validator with less bonded tokens: invalid request`

### Success

1. Confirm that `akasha_val` has 62000000000000000000000 tokens. NOTE: this may
vary if `akasha_val` has already been slashed:

```
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="akasha_val").tokens'
```

2. Rdelegate to `akasha_val` from `sif_val`

```
sifnoded tx staking redelegate sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna sifvaloper1l7hypmqk2yc334vc6vmdwzp5sdefygj250dmpy 100stake \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --gas auto \
  --broadcast-mode block -y
```

3. Confirm that `akasha_val` now has 62000000000000000000100 tokens (previously they had 62000000000000000000000 tokens). NOTE: this may
vary if `akasha_val` has already been slashed:

```
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="akasha_val").tokens'
```
