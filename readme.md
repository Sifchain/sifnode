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
git checkout tags/monkey-bars-testnet-5
```

4. Build:

```
make install
```

5. If you're a new operator (and only if - as otherwise this will reset your node!): 

    5.1 Scaffold your new node (from the project root directory):
    
    ```
    rake 'genesis:sifnode:scaffold[monkey-bars, 55f250c42b6e7bdcce6fe1a8af65f13e7c33aafb@35.166.247.98:26656, http://35.166.247.98:26657/genesis]'
    ```

6. If you're an existing node operator:

    6.1 Reset your node state:
    
    ```
    sifnoded unsafe-reset-all
    ```

    6.2 Download the latest genesis file:

    ```
    wget -O ~/.sifnoded/config/genesis.json https://raw.githubusercontent.com/Sifchain/networks/feature/genesis/testnet/monkey-bars-testnet-5/genesis.json
    ```
   
    6.3 Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

    ```
    persistent_peers = "55f250c42b6e7bdcce6fe1a8af65f13e7c33aafb@35.166.247.98:26656"
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

and you should see the following main validator node/s for Sifchain:

```
validators:
- address: sifvalcons1lze43sc3dyr9fnftk3hc3q3t8e4h6g5ar8776x
  pubkey: sifvalconspub1zcjduepqxp03gq8py26fqe9fljuppqv0s7859m2pvqrrgnr9v9p58a8j6eus8yx2rm
  proposerpriority: -5000
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

## Additional Resources

- [Additional instructions on standing up Sifnode](https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_).
- [Instructions on using Ethereum <> Sifchain cross-chain functionality](https://youtu.be/r81NQLxMers).

You can also ask questions on Discord [here](https://discord.com/invite/zZTYnNG).
