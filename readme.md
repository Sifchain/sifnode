# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Requirements

- [Ruby 2.6.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install)

## Getting started

### Setup

#### Connect to Sifchain

1. Ensure you've installed the dependencies listed above.

2. Clone the repository:

```
git clone ssh://git@github.com/Sifchain/sifnode && cd sifnode
```

3. Checkout the latest testnet release:

```
git checkout tags/monkey-bars-testnet-4
```

4. Build:

```
make install
```

5. If you're a new operator (and only if - as otherwise this will reset your node!): 

    5.1 Change to the `build` directory:

    ```
    cd ./build
    ```

    5.2 Scaffold your new node:
    
    ```
    rake 'genesis:sifnode:scaffold[monkey-bars, ec03640d0dcb1160f8cf73c33c63b64a55c93906@35.166.247.98:26656, http://35.166.247.98:26657/genesis]'
    ```

6. If you're an existing node operator:

    6.1 Reset your node state:
    
    ```
    sifnoded unsafe-reset-all
    ```

    6.2 Download the latest genesis file:

    ```
    wget -O ~/.sifnoded/config/genesis.json https://raw.githubusercontent.com/Sifchain/networks/feature/genesis/testnet/monkey-bars-testnet-4/genesis.json
    ```
   
    6.3 Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

    ```
    persistent_peers = "ec03640d0dcb1160f8cf73c33c63b64a55c93906@35.166.247.98:26656,04fad3abcf8d5c6d94d7815f9485830c280a8d73@35.166.247.98:28002,330c1b876d916f7518562b33d2749e3d1fcf7817@35.166.247.98:28004,16d9c23623e42723dfcf3dcbb11d98d989689a7a@35.166.247.98:28006"
    ```

7. Start your node:

```
sifnoded start
```

Run sync progress script and wait for message `Node synced`  

```
curl -s https://raw.githubusercontent.com/c29r3/sifchain-utils/main/sync_progress.sh | bash -s - http://35.166.247.98:26657
```


#### Become a Validator

You won't be able to participate in consensus until you become a validator.

1. Reach out to us on [Discord](https://discord.gg/3gQsRvjsRx) to request some tokens.

2. Run the following command to become a validator: 

```
MONIKER=$(awk -F'[ ="]+' '$1 == "moniker" { print $2 }' $HOME/.sifnoded/config/config.toml); \
sifnodecli tx staking create-validator \
    --amount=1000000rowan \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=$MONIKER \
    --commission-rate="0.10" \
    --commission-max-rate="0.20" \
    --commission-max-change-rate="0.01" \
    --min-self-delegation="1" \
    --gas="auto" \
    --gas-adjustment="1.2" \
    --gas-prices="0.025rowan" \
    --from=$MONIKER \
    --yes
```

## Peers

New node operators may also use the following peer addresses:

```
04fad3abcf8d5c6d94d7815f9485830c280a8d73@35.166.247.98:28002
330c1b876d916f7518562b33d2749e3d1fcf7817@35.166.247.98:28004
16d9c23623e42723dfcf3dcbb11d98d989689a7a@35.166.247.98:28006
```

## Additional Resources

- [Additional instructions on standing up Sifnode](https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_).
- [Instructions on using Ethereum <> Sifchain cross-chain functionality](https://youtu.be/r81NQLxMers).

You can also ask questions on Discord [here](https://discord.com/invite/zZTYnNG).
