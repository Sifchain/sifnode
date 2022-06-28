1. Initialize the chain

```
make init
```

2. Decrease the governance voting period time before first start;


```bash
echo "$(jq '.app_state.gov.voting_params.voting_period = "60s"' $HOME/.sifnoded/config/genesis.json)" > $HOME/.sifnoded/config/genesis.json
```

3. Start the chain:

```
make run
```

4. List upgrade proposals:

```
sifnoded q gov proposals --chain-id localnet
```

5. Raise an upgrade proposal:


```bash
sifnoded tx gov submit-proposal software-upgrade plan_name \
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

6. Check deposits:

```
sifnoded q gov deposits 1
```

7. Vote on proposal:

```
sifnoded tx gov vote 1 yes --from sif --chain-id localnet --keyring-backend test -y --broadcast-mode block
```

The node will have a consensus failure when it reaches the "upgrade-height". Restarting the node will not be enough for the chain to continue a new sifnoded release is required

8. Make a new sifnoded release:

  i. Update "version" file content to "plan_name"
  ii. Update "app/setup_handlers.go" "releaseVersion" constant to "plan_name"

6. Run the new release:

```
make run
```
