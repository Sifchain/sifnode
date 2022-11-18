This tutorial demonstrates the ability to add asymmetrically to a pool.
It also shows how adding asymmetrically to a
pool then removing liquidity is equivalent to performing a swap, that is the liquidity
provider does not achieve a cheap swap by adding then removing from the pool.

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
  "height": "7",
  "pagination": {
    "next_key": null,
    "total": "0"
  }
}
```

4. Query akasha balances:

```bash
sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test)
```

ceth: 500000000000000000000000
rowan: 500000000000000000000000

5. Query pool share estimation:

```
sifnoded query clp estimate-pool-share --externalAmount=0 --nativeAmount=1000000000000000000 --symbol ceth --output json
```

```json
{
	"percentage": "0.183227703974514619",
	"native_asset_amount": "549683111923543857",
	"external_asset_amount": "366455407949029238",
	"swap_info": {
		"status": "SELL_NATIVE",
		"fee": "1102674246586848",
		"fee_rate": "0.003000000000000000",
		"amount": "450316888076456143",
		"result": "366455407949029237"
	}
}
```

6. Add liquidity asymmetrically from akasha account to the ceth pool

```bash
sifnoded tx clp add-liquidity \
  --from akasha \
  --keyring-backend test \
  --symbol ceth \
  --nativeAmount 1000000000000000000 \
  --externalAmount 0 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

7. Query akasha balances:

```
sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test)
```

ceth: 500000000000000000000000
rowan: 499998900000000000000000


8. Query ceth lps:

```bash
sifnoded q clp lplist ceth
```

9. Remove the liquidity added by akasha in the previous step

```bash
sifnoded tx clp remove-liquidity \
  --from akasha \
  --keyring-backend test \
  --symbol ceth \
  --wBasis 10000 \
  --asymmetry 0 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

10. Query akasha balances:

```bash
sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test)
```

ceth: 500000366455407949029238
rowan: 499999349683111923543856

# Confirm share estimation calculation

We'll confirm the pool share estimation in step five was correct. The estimated pool share was given as:

```json
{
	"percentage": "0.183227703974514619",
	"native_asset_amount": "549683111923543857",
	"external_asset_amount": "366455407949029238"
}
```

Before removing liquidity (but after adding liquidity) akasha's balances were:

ceth: 500000000000000000000000
rowan: 499998900000000000000000

After removing liquidity akasha's balances were:

ceth: 500000366455407949029238
rowan: 499999349683111923543856

The gain in ceth after removing liquidity equals, 500000366455407949029238 - 500000000000000000000000 = 366455407949029238. This equals the estimated
external asset amount.

The gain in rowan after removing liquidity equals, 499999349683111923543856 - 499998900000000000000000 = 449683111923543856. 100000000000000000rowan was spent on tx fees, so if we add this to the difference, 449683111923543856 + 100000000000000000 = 549683111923543856. This equals the estimated
external asset amount.

# Confirm correct pool accounting

akasha started with 500000000000000000000000rowan and now has 499999349683111923543856rowan. So akasha has
500000000000000000000000rowan - 499999349683111923543856rowan = 650316888076456144rowan less rowan.
200000000000000000rowan was spent on tx fees. So 650316888076456144rowan - 200000000000000000rowan = 450316888076456144rowan
was given to the pool by akasha. In return akasha has gained 500000366455407949029238 - 500000000000000000000000 = 366455407949029238ceth
from the pool.

11. Check akash's gains/losses are reflected in the pool balances

```
sifnoded q clp pool ceth
```

external_asset_balance: "1633544592050970762"
native_asset_balance: "2450316888076456144"

We can confirm that what akasha has lost the pool has gained and vice versa

native_asset_balance = original_balance + amount_added_by_akasha
                     = 2000000000000000000rowan + 450316888076456144rowan
                     = 2450316888076456144rowan

Which equals the queried native asset pool balance.

external_asset_balance = original_balance - amount_gained_by_akasha
                       = 2000000000000000000ceth - 366455407949029238ceth
                       = 1633544592050970762ceth

Which equals the queried external asset pool balance.

# Equivalence with swapping

Has akasha had a "cheap" swap? How much would akasha have if instead of adding and removing from the pool
they had simply swapped 450316888076456144rowan for ceth?

12. Reset the chain

```bash
make init
make run
```

13. Recreate the ceth pool

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

14. Swap 450316888076456144rowan for ceth from akasha:

```bash
sifnoded tx clp swap \
  --from akasha \
  --keyring-backend test \
  --sentSymbol rowan \
  --receivedSymbol ceth \
  --sentAmount 450316888076456144 \
  --minReceivingAmount 0 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  -y
```

15. Query akasha balances:

```bash
sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test)
```

ceth: 500000366455407949029237
rowan: 499999449683111923543856

akasha has swapped 450316888076456144rowan for 366455407949029237ceth.
By adding then removing from the pool, akasha gained 366455407949029238ceth and provided 450316888076456144rowan to the pool.
So both actions are almost identical, except akasha gains 1ceth more by adding then removing from the pool rather than swapping.
This is a rounding error. Which means, as expected, adding asymmetrically then removing
liquidity is equivalent to swapping.


