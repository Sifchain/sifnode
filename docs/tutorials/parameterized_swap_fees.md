# Paramaterized swap fees

This tutorial demonstrates the behaviour of the parameterized swap fee functionality.

1. Start and run the chain:

```bash
make init
make run
```

2. Create a pool:

```bash
sifnoded tx clp create-pool \
  --from sif \
  --keyring-backend test \
  --symbol ceth \
  --nativeAmount 2000000000000000000 \
  --externalAmount 2000000000000000000 \
  --fees 100000000000000000rowan \
  --broadcast-mode block \
  --chain-id localnet \
  -y
```

3. Confirm pool has been created:

```bash
sifnoded q clp pools --output json | jq
```

returns:

```json
{
  "pools": [
    {
      "external_asset": {
        "symbol": "ceth"
      },
      "native_asset_balance": "2000000000000000000",
      "external_asset_balance": "2000000000000000000",
      "pool_units": "2000000000000000000",
      "swap_price_native": "1.000000000000000000",
      "swap_price_external": "1.000000000000000000",
      "reward_period_native_distributed": "0",
      "external_liabilities": "0",
      "external_custody": "0",
      "native_liabilities": "0",
      "native_custody": "0",
      "health": "0.000000000000000000",
      "interest_rate": "0.000000000000000000",
      "last_height_interest_rate_computed": "0",
      "unsettled_external_liabilities": "0",
      "unsettled_native_liabilities": "0",
      "block_interest_native": "0",
      "block_interest_external": "0"
    }
  ],
  "clp_module_address": "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
  "height": "5",
  "pagination": {
    "next_key": null,
    "total": "0"
  }
}
```

4. Query the current swap fee params:

```bash
sifnoded q clp swap-fee-params --output json | jq
```

```json
{
  "default_swap_fee_rate": "0.003000000000000000",
  "token_params": []
}
```

5. Set new swap fee params

```bash
sifnoded tx clp set-swap-fee-params \
   --from sif \
   --keyring-backend test \
   --chain-id localnet \
   --broadcast-mode block \
   --fees 100000000000000000rowan \
   -y \
   --path <( echo '{
    "default_swap_fee_rate": "0.003",
    "token_params": [{
            "asset": "ceth",
            "swap_fee_rate": "0.004"
        },
        {
            "asset": "rowan",
            "swap_fee_rate": "0.002"
        }
    ]
    }' )
```


6. Check swap fee params have been updated:

```bash
sifnoded q clp swap-fee-params --output json | jq
```

```json
{
  "default_swap_fee_rate": "0.003000000000000000",
  "token_params": [
    {
      "asset": "ceth",
      "min_swap_fee": "0",
      "swap_fee_rate": "0.004000000000000000"
    },
    {
      "asset": "rowan",
      "min_swap_fee": "600000000000",
      "swap_fee_rate": "0.002000000000000000"
    }
  ]
}
```

7. Do a swap:

```bash
sifnoded tx clp swap \
  --from sif \
  --keyring-backend test \
  --sentSymbol ceth \
  --receivedSymbol rowan \
  --sentAmount 200000000000000 \
  --minReceivingAmount 0 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  --broadcast-mode block \
  --output json \
  -y | jq '.logs[0].events[] | select(.type=="swap_successful").attributes[] | select(.key=="swap_amount" or .key=="liquidity_fee")'
```

Returns:

```json
{
  "key": "swap_amount",
  "value": "199180081991801"
}
{
  "key": "liquidity_fee",
  "value": "799920007999"
}

```

We've swapped ceth for rowan, so the ceth swap fee rate, `0.004`, should have been used. So the expected `swap_amount` and `liquidity_fee` are:

```
adjusted_output = x * Y  / ((x + X)(1 + r))
                = 200000000000000 * 2000000000000000000 / ((200000000000000 + 2000000000000000000) * (1 + 0))
                = 199980001999800

liquidity_fee = f * adjusted_output, min_swap_fee
              = 0.004 * 199980001999800
	            = 799920007999

y = adjusted_amount - liquidity_fee
  = 199980001999800 - 799920007999
  = 199180081991801
```

Which match the vales returned by the swap command.

8. Confirm that setting swap fee rate greater than one fails:

```bash
sifnoded tx clp set-swap-fee-params \
   --from sif \
   --keyring-backend test \
   --chain-id localnet \
   --broadcast-mode block \
   --fees 100000000000000000rowan \
   -y \
   --path <( echo '{
    "default_swap_fee_rate": "0.003",
    "token_params": [{
            "asset": "ceth",
            "swap_fee_rate": "1.2"
        },
        {
            "asset": "rowan",
            "swap_fee_rate": "0.002"
        }
    ]
    }' )
```
