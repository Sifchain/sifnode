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

and within a few seconds, your node should start synchronizing with the network.

You can verify that you're connected by running:

```
sifnodecli q tendermint-validator-set
```

and you should see the following main validator nodes for Sifchain:

```
validators:
- address: sifvalcons1z6jhzs0f7v02ny6k5x5rekf7gyx9400zyxmzve
  pubkey: sifvalconspub1zcjduepq4zyan4mlm8fpku5jd7zu7f59k863x4g2wrzkku0285z6xylppk6q6nkzrk
  proposerpriority: -5000
  votingpower: 5000
- address: sifvalcons192ljdnz3u6d7l7vg9zgstlnqyczqhwz4wj5ltz
  pubkey: sifvalconspub1zcjduepq8zdt2xty2kk87zrzn95crwjkpmhmzxu6w05wtn08dxhq0qnj090sxg634l
  proposerpriority: -5000
  votingpower: 5000
- address: sifvalcons1v38zwh0f9hwq5x6hfna35pr9x5r5wpydqfgyat
  pubkey: sifvalconspub1zcjduepqefzlm5pymv84kfxdrzm627pw9ty6v2zd49dzuc3aan9z2pftk4rqckj2gz
  proposerpriority: -5000
  votingpower: 5000
- address: sifvalcons1kulx53jp3vnhmagsha5ncnsuewqf3s00nwzffv
  pubkey: sifvalconspub1zcjduepqrs8w58a59cu3wtt03rtm0c03gyt84f8pxwvtp7cptly39vhcdyxsyqmf62
  proposerpriority: 15000
  votingpower: 5000
```

you are now connected to the network.

#### Become a Validator

You won't be able to participate in consensus until you become a validator.

1. Reach out to us on [Discord](https://discord.gg/3gQsRvjsRx) to request some tokens.

2. Obtain your node moniker (if you don't already know it):

```
cat ~/.sifnoded/config/config.toml | grep moniker
```

3. Run the following command to become a validator (*replace `<moniker>` with your node's actual moniker*): 

```
sifnodecli tx staking create-validator \
    --commission-max-change-rate 0.1 \
    --commission-max-rate 0.1 \
    --commission-rate 0.1 \
    --amount 1000000000rowan \
    --pubkey $(sifnoded tendermint show-validator) \
    --moniker <moniker> \
    --chain-id monkey-bars \
    --min-self-delegation 1 \
    --gas auto \
    --from <moniker> \
    --keyring-backend file
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
