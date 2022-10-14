# Minimum swap fees

This is a proposal to set a minimum swap fee.

## Current swap behaviour

Due to ratio shifting the current swap behaviour depends on whether the user is buying or selling Rowan.

In the following formulas:

```
X - input depth (balance + liabilities)
Y - output depth (balance + liabilities)
x - input amount
y - output amount
r - current ratio shifting running rate
f - swap fee rate, this must satisfy 0 =< f =< 1
min_fee - minimum fee for a swap
```

### Swapping to Rowan:

```
y = (1 - f) * x * Y / ((x + X) * (1 + r))
fee = f * x * Y / ((x + X) * (1 + r))
```

Equivalently this can be written as:

```
raw_XYK_output = x * Y / (x + X)
adjusted_output = raw_XYK_output / (1 + r)

(1) fee = f * adjusted_output
y = adjusted_output - fee
```

### Swapping from Rowan:

```
y = (1 - f) * (1 + r) * x * Y / (x + X)
fee = f * (1 + r) * x * Y / (x + X)
```

Similar to the case of swapping to rowan, this can be written as:

```
raw_XYK_output = x * Y / (x + X)
adjusted_output = raw_XYK_output * (1 + r)

(2) fee = f * adjusted_output
y = adjusted_output - fee
```

## Proposed Change

Apply a minimum fee when swapping.

The fee calculation in equation (1) and (2) becomes:

```
fee = min(max(f * adjusted_output, min_fee), adjusted_output)
```

The min function is required to ensure that the fee is not greater than the adjusted output.

### Parameters

The `min_fee` must be in the same denomination as the buy token, in which case a single value for `min_fee` (greater than zero)
cannot be defined across all tokens. Consequently a default `min_fee` of zero must be applied and an admin account
must be able to specify override `min_fee` values for specific tokens. The admin account must
not be able to change the default `min_fee` value since zero is the only reasonable value.

The same `min_fee` must be used when buying a token
regardless of the sell token e.g. it will not be possible to set one `min_fee` for swapping atom to Rowan
and another `min_fee` for swapping usdc to Rowan.

Unlike the `min_fee` it is possible to define a swap fee rate, `f`, which can be applied on all swaps.
The admin account must be able to specify a default swap fee rate and specify override values for specific tokens.
As with the `min_fee` the same swap fee rate must be used regardless of the sell token.

See the CLI section for commands for setting and querying the swap fee params.

## Events

There are no new events or updates to existing events.

## CLI

The CLI option for querying the swap fee rate (`sifnoded q clp swap-fee-rate`) and setting the swap fee
rate (`sifnoded tx clp set-swap-fee-rate`), must be renamed to `sifnoded q clp swap-fee-params`
and `sifnoded tx clp set-swap-fee-params` and updated to include the min fee.

### Setting

The CLI should validate that the min fees are valid cosmos Uint256.

```bash
sifnoded tx clp set-swap-fee-params \
  --from sif \
  --path ./swap-fee-params.json \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

```json
{
	"default_swap_fee_rate": "0.003",
	"token_params": [{
			"asset": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			"swap_fee_rate": "0.002",
			"min_swap_fee": "12"
		},
		{
			"asset": "cusdc",
			"swap_fee_rate": "0.002",
			"min_swap_fee": "800"
		},
		{
			"asset": "rowan",
			"swap_fee_rate": "0.001",
			"min_swap_fee": "12"
		}
	]
}
```

### Querying

```bash
sifnoded q clp swap-fee-params --output json
```

```json
{
	"default_swap_fee_rate": "0.003",
	"token_params": [{
			"asset": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			"swap_fee_rate": "0.002",
			"min_swap_fee": "12"
		},
		{
			"asset": "cusdc",
			"swap_fee_rate": "0.002",
			"min_swap_fee": "800"
		},
		{
			"asset": "rowan",
			"swap_fee_rate": "0.001",
			"min_swap_fee": "12"
		}
	]
}
```
