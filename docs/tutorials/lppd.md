1. Initialize and start the chain. From the repo root run the following:

```
make init
make run
```

2. Create pool - add one million dollars of usdt and price rowan at 10c:

```
sifnoded tx clp create-pool \
  --from sif \
  --keyring-backend test \
  --symbol cusdt \
  --nativeAmount 10000000000000000000000000 \
  --externalAmount 1000000000000 \
  --fees 100000000000000000rowan \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```

3. Query the lppd parameters:

```
sifnoded q clp lppd-params
```

4. Set a new lppd policy:

```
sifnoded tx clp set-lppd-params \
   --from sif \
   --keyring-backend test \
   --chain-id localnet \
   --broadcast-mode block \
   -y \
   --path <( echo '[
    {
        "distribution_period_block_rate": "0.01",
        "distribution_period_start_block": 1,
        "distribution_period_mod": 1,
        "distribution_period_end_block": 433000
    }
]' )
```

5. Query block results and observe the lppd success event:

```
curl http://localhost:26657/block_results
```
