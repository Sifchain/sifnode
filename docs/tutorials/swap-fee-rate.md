# Fixed Rate Fee Model

The swap formula uses the fixed rate fee model to calculate swap fees and consequently
the amount of received token.

## Formulas

Due to ratio shifting the formula depends on whether the swap is from Rowan or to Rowan.

In the following formulas:

```
X - Input balance
Y - Output balance
x - Input amount
y - Output amount
r - Current ratio shifting running rate
f - Swap fee rate
```

### Swapping to Rowan:

```
y = (1 - f) * x * Y / ((x + X)(1 + r))
fee = f * x * Y / ((x + X)(1 + r))
```
### Swapping rom Rowan:

```
y = (1 - f) * (1 + r) * x * Y / (x + X)
fee = f * (1 + r) * x * Y / (x + X)
```

## Examples

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
	"pools": [{
		"external_asset": {
			"symbol": "ceth"
		},
		"native_asset_balance": "2000000000000000000",
		"external_asset_balance": "2000000000000000000",
		"pool_units": "2000000000000000000",
		"swap_price_native": "1.000000000000000000",
		"swap_price_external": "1.000000000000000000",
		"reward_period_native_distributed": "0"
	}],
	"clp_module_address": "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
	"height": "50",
	"pagination": {
		"next_key": null,
		"total": "0"
	}
}
```

4. Query the current swap fee rate:

```bash
sifnoded q clp swap-fee-rate --output json | jq
```

```json
{
	"swap_fee_rate": "0.003000000000000000"
}
```

5. Do a swap:

```
sifnoded tx clp swap \
  --from sif \
  --keyring-backend test \
  --sentSymbol ceth \
  --receivedSymbol rowan \
  --sentAmount 200000000000000 \
  --minReceivingAmount 0 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

This will return a tx hash.

6. Use the tx hash to get the swap amount and liquidity fee:

```bash
TXHASH=1AB7D2B0C4EDC2B18893334E60BFCF3C3F9587314082D314CA641D895F216E62
sifnoded q tx $TXHASH --output json | jq '.logs[0].events[] | select(.type=="swap_successful").attributes[] | select(.key=="swap_amount" or .key=="liquidity_fee")'
```

which returns:

```json
{
  "key": "swap_amount",
  "value": "199380061993800"
}
{
  "key": "liquidity_fee",
  "value": "599940005999"
}
```

The swap amount is as expected:
```
y = (1 - f) * x * Y  / ((x + X)(1 + r))
  = (1 - 0.003) * 200000000000000 * 2000000000000000000 / ((200000000000000 + 2000000000000000000) * (1 + 0))
  = 199380061993800
```

And the swap fee is as expected
```
fee = f * x * Y  / ((x + X)(1 + r))
    = 0.003 * 200000000000000 * 2000000000000000000 / ((200000000000000 + 2000000000000000000) * (1 + 0))
    = 599940005999
```

7. Check the pool balances

```bash
sifnoded q clp pools --output json | jq
```

```json
 {
 	"native_asset_balance": "1999800619938006200",
 	"external_asset_balance": "2000200000000000000"
 }
```

These are as expected:
```
native_asset_balance = init_native - y
                     = 2000000000000000000 - 199380061993800
                     = 1999800619938006200
                     

external_asset_balance = init_external + x
                       = 2000000000000000000 + 200000000000000
                       = 2000200000000000000
```
8. Change the swap fee rate to 0.01

```bash
sifnoded tx clp set-swap-fee-rate \
  --from sif \
  --swapFeeRate 0.01 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

Confirm the new rate:

```bash
sifnoded q clp swap-fee-rate --output json | jq
```

```json
{
	"swap_fee_rate": "0.010000000000000000"
}
```

9. Do another swap, this time the other way around (rowan to ceth):

```
sifnoded tx clp swap \
  --from sif \
  --keyring-backend test \
  --sentSymbol rowan \
  --receivedSymbol ceth \
  --sentAmount 200000000000000 \
  --minReceivingAmount 0 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

10. Repeat steps 6 & 7, confirm the results:

```json
{
  "key": "swap_amount",
  "value": "198019738620019"
}
{
  "key": "liquidity_fee",
  "value": "2000199380000"
}
```

```json
{
  "native_asset_balance": "2000000619938006200",
  "external_asset_balance": "2000001980261379981"
}
```

Which are as expected:
```
y = (1 - f) * (1 + r) * x * Y / (x + X)
  = (1 - 0.01) * (1 + 0) * 200000000000000 * 2000200000000000000 / (200000000000000 + 1999800619938006200)
  = 198019738620019

fee = f * (1 + r) * x * Y / (x + X)
    = 0.01 * (1 + 0) * 200000000000000 * 2000200000000000000 / ((200000000000000 + 1999800619938006200))
    = 2000199380000

native_asset_balance = init_native + x
                     = 1999800619938006200 + 200000000000000
                     = 2000000619938006200

external_asset_balance = init_external - y
                       = 2000200000000000000 - 198019738620019
                       = 2000001980261379981
```
10. Try to change swap fee > 1. This should fail:

```
sifnoded tx clp set-swap-fee-rate \
  --from sif \
  --swapFeeRate 1.00001 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

Which returns:

`Error: swap rate fee must be less than or equal to one`

11. Try to change swap fee < 0. This should fail:

```
sifnoded tx clp set-swap-fee-rate \
  --from sif \
  --swapFeeRate -0.0001 \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

`Error: swap rate fee must be greater than or equal to zero`
