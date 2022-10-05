# Minimum swap fees

This is a proposal to set a minimum swap fee when selling Rowan.

## Current swap behaviour

Due to ratio shifting the current behaviour depends on whether the user is buying or selling Rowan.

In the following formulas:

```
X - input depth (balance + liabilities)
Y - output depth (balance + liabilities)
x - input amount
y - output amount
r - current ratio shifting running rate
f - swap fee rate, this must satisfy 0 =< f =< 1
```

### Swapping to Rowan:

```
y = (1 - f) * x * Y / ((x + X) * (1 + r))
fee = f * x * Y / ((x + X) * (1 + r))
```

Equivalently this can be written as:

```
raw_XYK_output = x * Y / ((x + X) * (1 + r))

(1) fee = f * raw_XYK_output
y = raw_XYK_output - fee
```

### Swapping from Rowan:

```
y = (1 - f) * (1 + r) * x * Y / (x + X)
fee = f * (1 + r) * x * Y / (x + X)
```

## Proposed Change

Apply a minimum fee when swapping to Rowan.

The fee calculation in equation (1) becomes:

```
fee = max(f * raw_XYK_output, min_fee)
```

Where `min_fee` is a minimum fee parameter which is set via an admin key. See CLI
section for more details.

## Events

There are no new events or updates to existing events.

## CLI

The CLI option for querying the swap fee rate (`sifnoded q clp swap-fee-rate`) and setting the swap fee
rate (`sifnoded tx clp set-swap-fee-rate`), must be renamed to `sifnoded q clp swap-params`
and `sifnoded tx clp set-swap-params` and updated to include the min fee.

### Setting

The CLI should validate that the min fee is a valid cosmos Uint256.

```bash
sifnoded tx clp set-swap-params \
  --from sif \
  --swapFeeRate 0.01 \
  --minFee 100 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

### Querying

```bash
sifnoded q clp swap-params --output json
```

```json
{
    "swap_fee_rate": "0.010000000000000000",
    "min_fee": "100"
}
```
