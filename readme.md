# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Prerequisites / Dependencies

- [Ruby 2.6.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)
- [curl](https://curl.haxx.se/download.html)

## Getting started

### Setup

#### Connect to Sifchain

##### New Validators

1. Ensure you've installed the dependencies listed above.

2. Clone the repository:

```
git clone ssh://git@github.com/Sifchain/sifnode && cd sifnode
```

3. Checkout the latest testnet release:

```
git checkout tags/merry-go-round-1
```

4. Build:

```
make clean install
```

5. Generate a mnemonic:

```
rake "keys:generate:mnemonic"
```

6. Scaffold your node:

```
rake "genesis:sifnode:scaffold[merry-go-round, <moniker>, '<mnemonic>', '', ff0dd55dffa0e67fe21e2c85c80b0c2894bf2586@52.89.19.109:26656, http://52.89.19.109:26657/genesis]"
```

* Replace `<moniker>` with the moniker (name) of your node. 
* Replace `<mnemonic>` with the mnemonic phrase generated in the previous step.

This step will also output the keyring password, so please record this and the moniker somewhere secure.

7. Connect:

```
sifnoded start
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

##### Existing Validators

1. Checkout the latest testnet release:

```
git fetch && git checkout tags/merry-go-round-1
```

2. Build:

```
make install
```

3. Reset your local state (please take a backup of your keyring first):

```
sifnodecli unsafe-reset-all
```

4. Download the new genesis file:

```
curl http://52.89.19.109:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "ff0dd55dffa0e67fe21e2c85c80b0c2894bf2586@52.89.19.109:26656,8e58e41e9e47c53f63755d60fe0f35286a96b70f@54.74.58.153:26656,853ae9203af1606ca3497b845f775ece249be5ff@13.54.242.178:26656,910336f8e1342915c3adc40a73a6924d4d974c85@3.0.235.227:26656"
```

6. Start your node:

```
sifnoded start
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

##### Verify

You can verify that you're connected by running:

```
sifnodecli q tendermint-validator-set
```

and you should see the following primary validator node/s for Sifchain:

```
validators:
- address: sifvalcons1z8jyamggawyute8m7a6tfk76whdegz4hhu47kx
  pubkey: sifvalconspub1zcjduepq5geuxq3kyuwayc9ht82y997ncmh3qfe4eqg837kmf6d3tnyspemq6e83zz
  proposerpriority: 5625
  votingpower: 5000
- address: sifvalcons1rya4cf6ejuzsn3qv5c97j3spsr70dwftdygktq
  pubkey: sifvalconspub1zcjduepqewcxpth6dtk82f826gh5te07xyuk04t9y8dg63ndkngsr79dtu0skarrel
  proposerpriority: 4375
  votingpower: 5000
- address: sifvalcons1tzsa80axse3urga7vcck2r638awkgfj6sddm8q
  pubkey: sifvalconspub1zcjduepqdp72sdqtwjujpqlzfg0ku8smark7832ck4440nnqh5yz7yly78fsc0sjqx
  proposerpriority: 3125
  votingpower: 5000
- address: sifvalcons1s00sdg5z5wv89yxjc66ft6uaf0lcqphvef4f9h
  pubkey: sifvalconspub1zcjduepq25y7tsy0c9f0d7u43x7csfpry7t5ur4lcemfdjcrkctjv4hv7taqxvvhx7
  proposerpriority: -13125
  votingpower: 5000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
ff0dd55dffa0e67fe21e2c85c80b0c2894bf2586@52.89.19.109:26656
8e58e41e9e47c53f63755d60fe0f35286a96b70f@54.74.58.153:26656
853ae9203af1606ca3497b845f775ece249be5ff@13.54.242.178:26656
910336f8e1342915c3adc40a73a6924d4d974c85@3.0.235.227:26656
```

#### Become a Validator

You won't be able to participate in consensus until you become a validator.

1. Reach out to us on [Discord](https://discord.gg/3gQsRvjsRx) to request some tokens.

2. Obtain your node moniker (if you don't already know it):

```
cat ~/.sifnoded/config/config.toml | grep moniker
```

3. Run the following command to become a validator: 

```
sifnodecli tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount=1000000000rowan \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    --gas="auto" \
    --from=<moniker> \
    --keyring-backend=file
```

* Replace `<moniker>` with the moniker (name) of your node. 

## Additional Resources

- [Additional instructions on standing up Sifnode](https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_).
- [Instructions on using Ethereum <> Sifchain cross-chain functionality](https://youtu.be/r81NQLxMers).

You can also ask questions on Discord [here](https://discord.com/invite/zZTYnNG).
