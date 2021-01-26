# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Prerequisites / Dependencies

- [Ruby 2.6.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install) 
  - Add `export GOPATH=~/go` to your shell
  - Add `export PATH=$PATH:$GOPATH/bin` to your shell
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
git checkout tags/merry-go-round-3
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
rake "genesis:sifnode:scaffold[merry-go-round, <moniker>, '<mnemonic>', e99deeec54ca1c477f8826801bc1fd29f5539a45@44.226.150.203:26656, http://44.226.150.203:26657/genesis]"
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
git fetch && git checkout tags/merry-go-round-3
```

2. Build:

```
make install
```

3. Reset your local state (please take a backup of your keyring first):

```
sifnoded unsafe-reset-all
```

4. Download the new genesis file:

```
curl http://44.226.150.203:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "e99deeec54ca1c477f8826801bc1fd29f5539a45@44.226.150.203:26656,fb4bfcaf9980a2ee3fe8298eadd21c9757d83f6c@52.49.165.39:26656,39172b7f5f8c2394af86f3174e4b8c9f6eb3ad3b@13.210.25.108:26656,e032bdcfad831c7c7131e02453fdc736625547a6@18.136.70.31:26656"
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
- address: sifvalcons1ya2a3w32py5w32lzhxttl04n80vxx5ke8xlt3k
  pubkey: sifvalconspub1zcjduepqw4p8n5s5zpdm9f9r0s0x8lz0e08vmp9j7zrpfmsyg3np0xvmvndqsjkvgw
  proposerpriority: 1625000
  votingpower: 1000000
- address: sifvalcons1xnn99v7n50m79pxmveeqdj0f6xuefhn4vay6d6
  pubkey: sifvalconspub1zcjduepqpug6d7p5gt69j9jzpltxlheukz86u3emq2qe89d5jlx7cctfyceq4warur
  proposerpriority: -2125000
  votingpower: 1000000
- address: sifvalcons1jycgr4z4trfrwk3x04x3nssyv6uyzf2yzen9j3
  pubkey: sifvalconspub1zcjduepq6h5q8vf9z7h5kt9d986vl7x6x0zx9w9x805f9wn5ys5rw4mme0yqxddd2g
  proposerpriority: -875000
  votingpower: 1000000
- address: sifvalcons1hdh8xnvg3ckrl0pjnryc225qdg0rr0rgsdmqwz
  pubkey: sifvalconspub1zcjduepq69qwgexfxxfh9jtp4nu0jrvf4hq8m9k4079z995ss0qx0qp5qd7sln7cm2
  proposerpriority: 1375000
  votingpower: 1000000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
e99deeec54ca1c477f8826801bc1fd29f5539a45@44.226.150.203:26656
fb4bfcaf9980a2ee3fe8298eadd21c9757d83f6c@52.49.165.39:26656
39172b7f5f8c2394af86f3174e4b8c9f6eb3ad3b@13.210.25.108:26656
e032bdcfad831c7c7131e02453fdc736625547a6@18.136.70.31:26656
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
    --amount="10000000000000000000rowan" \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=<moniker> \
    --keyring-backend=file
```

* Replace `<moniker>` with the moniker (name) of your node. 

#### Block Explorer

A block explorer is available at:

* https://blockexplorer-testnet.sifchain.finance

## Additional Resources

- [Additional instructions on standing up Sifnode](https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_).
- [Instructions on using Ethereum <> Sifchain cross-chain functionality](https://www.youtube.com/watch?v=z1EZcetmDMI&t=2s).

You can also ask questions on Discord [here](https://discord.com/invite/zZTYnNG).
