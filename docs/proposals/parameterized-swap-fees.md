# Parameterized swap fee rates

This is a proposal to remove the minimum swap fee and introduce parameterized swap fee rates.

## Proposed behaviour

The swap formulas should be reverted to how they were before the introduction of min fees.

When swapping to rowan:

```
raw_XYK_output = x * Y / (x + X)
adjusted_output = raw_XYK_output / (1 + r)

(1) fee = f * adjusted_output
y = adjusted_output - fee
```

Swapping from rowan:

```
raw_XYK_output = x * Y / (x + X)
adjusted_output = raw_XYK_output * (1 + r)

(2) fee = f * adjusted_output
y = adjusted_output - fee
```

Where:

```
X - input depth (balance + liabilities)
Y - output depth (balance + liabilities)
x - input amount
y - output amount
r - current ratio shifting running rate
f - swap fee rate, this must satisfy 0 =< f =< 1
```

The admin account must be able to specify a default swap fee rate and specify override values for specific tokens. See the CLI section for commands for setting and querying the swap fee params.

The swap fee rate of the **sell** token must be used in the swap calculations.

When swapping between two non native tokens, TKN1:TKN2, the system performs two swaps, TKN1:rowan followed by rowan:TKN2. In this case the swap fee rate of TKN1 must be used for both swaps.

The swaps occuring during an open or close of a margin position use the default swap fee rate.

## Events

There are no new events or updates to existing events.

## CLI

### Setting

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
  "token_params": [
    {
      "asset": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
      "swap_fee_rate": "12"
    },
    {
      "asset": "cusdc",
      "swap_fee_rate": "800"
    },
    {
      "asset": "rowan",
      "swap_fee_rate": "12"
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
  "default_swap_fee_rate": "0.003000000000000000",
  "token_params": [
    {
      "asset": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
      "swap_fee_rate": "12"
    },
    {
      "asset": "cusdc",
      "swap_fee_rate": "800"
    },
    {
      "asset": "rowan",
      "swap_fee_rate": "12"
    }
  ]
}
```

## References

Product spec https://hackmd.io/MhRTYAsfR2qtP83jvmDdmQ
