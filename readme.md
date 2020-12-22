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
rake "genesis:sifnode:scaffold[sandpit,<moniker>,'<mnemonic>','',53b13d8391031f39f846c920762f322e5dde6af1@100.20.113.245:26656, http://100.20.113.245:26657/genesis]"
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
curl http://100.20.113.245:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "53b13d8391031f39f846c920762f322e5dde6af1@100.20.113.245:26656,f9459f02c0591ee25f79d9a1a9430e56a00a23b9@34.250.107.224:26656,dd5bc4a38852bc209278ca24a5ee959480432ec7@3.106.154.25:26656,9f0113530dc0574505e902298c2cf2e160a4f7a0@3.1.98.35:26656"
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
- address: sifvalcons1pulw6ucsaf2dtvuqwlr0qvyvhmql6xj60a45ua
  pubkey: sifvalconspub1zcjduepqe4ddf985za6ffdme07xeu7zp8svkw6lf4pygq6dg856eu76gfxeskzpmqh
  proposerpriority: 1875
  votingpower: 5000
- address: sifvalcons1wncevpl7pj0wdyy7u9mq2u8d906jrjce7nuznd
  pubkey: sifvalconspub1zcjduepqcpmsx3l7axe8wjzxw4j0kd7gr9qv9q0qlktgyyaztgp7826r677qtkckcu
  proposerpriority: -625
  votingpower: 5000
- address: sifvalcons1k0qxjng82sz57w2rstgd5thfpxjwjkq9gjmarm
  pubkey: sifvalconspub1zcjduepqh376rykq7qk9t4erl0g30l6f0gn24lmtph9grrf45mzun269g77qx4p0tt
  proposerpriority: -1875
  votingpower: 5000
- address: sifvalcons16cqp74x5cwat84fu7fpm7y0l5m0t6q3lg0m5tp
  pubkey: sifvalconspub1zcjduepqp3y73kvq7992h7jcjmrv4tvv8pj4k8u7v2tzht94pqxcmfzj85aqlz3804
  proposerpriority: 625
  votingpower: 5000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
53b13d8391031f39f846c920762f322e5dde6af1@100.20.113.245:26656
f9459f02c0591ee25f79d9a1a9430e56a00a23b9@34.250.107.224:26656
dd5bc4a38852bc209278ca24a5ee959480432ec7@3.106.154.25:26656
9f0113530dc0574505e902298c2cf2e160a4f7a0@3.1.98.35:26656
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
