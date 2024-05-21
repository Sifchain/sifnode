# Authz, Min Commission, Max Voting Power

The max voting power, min commission restrictions could be bypassed by wrapping the relevant
messages inside an authz message. Updates have been made to prevent this. This tutorial
shows that a delegate message inside an authz message is blocked.

1. Initialize and start the chain:

```shell
make init
make run
```

2. Attempt to delegate to `sif_val` from the `sif` account:

```shell
sifnoded tx staking delegate sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna 100stake \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block -y
```

This fails with the message:

`This validator has a voting power of 100.000000000000000000%. Delegations not allowed to a validator whose post-delegation voting power is more than 6.600000000000000000%. Please delegate to a validator with less bonded tokens: invalid request`

3. Authorize `akasha` to send delegate messages on behalf of `sif`

```shell
sifnoded tx authz grant sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 delegate \
  --allowed-validators sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna \
  --from=sif \
  --keyring-backend=test \
  --chain-id=localnet \
  --broadcast-mode block -y
```

4. Confirm that `akasha` is now authorized:

```shell
sifnoded query authz grants sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 /cosmos.staking.v1beta1.MsgDelegate
```

5. Create a delegate message from `sif` - this won't be signed and won't be sent to the chain:

```shell
sifnoded tx staking delegate sifvaloper1syavy2npfyt9tcncdtsdzf7kny9lh777dzsqna 100stake \
  --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
  --keyring-backend test \
  --chain-id localnet \
  --generate-only > tx.json
```

6. Confirm `sif_val` has 1000000000000000000000000 tokens:

```shell
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="sif_val").tokens'
```

7. Send an authz message from `akasha` with a delegate message from `sif`:

```shell
sifnoded tx authz exec tx.json \
  --from akasha \
  --keyring-backend=test \
  --chain-id=localnet \
  --yes
```

Which fails with the message:

`This validator has a voting power of 100.000000000000000000%. Delegations not allowed to a validator whose post-delegation voting power is more than 6.600000000000000000%. Please delegate to a validator with less bonded tokens: invalid request`

NOTE: previously this would have succeeded.

8. Confirm `sif_val` still has 1000000000000000000000000 tokens:

```shell
sifnoded query staking validators --output=json  | jq '.validators[] | select(.description.moniker=="sif_val").tokens'
```

9. Clenaup - remove tx.json:

```shell
rm tx.json
```
