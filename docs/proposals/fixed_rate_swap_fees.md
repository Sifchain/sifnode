# Fixed Rate Swap Fees

Sifchain is switching from a slip-based fee model to a fixed-rate fee model. This requires
the swap formula and the liquidity fee formula within sifnode to be updated.

## Fixed-rate formuals

Due to ratio shifting, the updated swap and liquidity fee formulas depend on whether the swap is
from Rowan or to Rowan.

In the following formulas:

```
X - Input balance
Y - Output balance
x - Input amount
y - Output amount
r - Current ratio shifting running rate
f - Swap fee rate. This must satisfy `0 =< f =< 1`
```

### Swapping to Rowan:

```
y = (1 - f) * x * Y / ((x + X)(1 + r))
fee = f * x * Y / ((x + X)(1 + r))
```

### Swapping from Rowan:

```
y = (1 - f) * (1 + r) * x * Y / (x + X)
fee = f * (1 + r) * x * Y / (x + X)
```

## Changing the swap rate fee

The swap fee rate, `f`, in the above formulas must be updatable with a regular Cosmos transaction
however the transaction must be signed by the PMTP/rewards admin key. The swap rate fee must
satisfy `0 =< f =< 1` otherwise the transaction is rejected.

## Events

There are no new events or updates to existing events.

## CLI

CLI options for setting and querying the swap fee rate must be implemented.

### Setting

The CLI should validate that the value of the swap rate fee satisfies `0 =< f =< 1`

```bash
sifnoded tx clp set-swap-fee-rate \
  --from sif \
  --swapFeeRate 0.01 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

### Querying

```bash
sifnoded q clp swap-fee-rate --output json
```

```json
{
	"swap_fee_rate": "0.010000000000000000"
}
```
## References

Background on the use of fixed rate fee swap formula for asymmetric adds https://hackmd.io/NjvaZY1qQiS17s_uEgZmTw?both