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
rake "genesis:sifnode:scaffold[merry-go-round, <moniker>, '<mnemonic>', '', a00c8a5f07d87754e9a2d428ad7e1877dbe12ddd@35.165.17.164:26656, http://35.165.17.164:26657/genesis]"
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
git fetch && git checkout tags/merry-go-round-2
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
curl http://35.165.17.164:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "a00c8a5f07d87754e9a2d428ad7e1877dbe12ddd@35.165.17.164:26656,231a0b22837c1ad627d34748ac27e21f540dbb87@54.195.149.238:26656,69b0b8fe353ea3e25b23fad18796056bdf5ce9c1@3.24.153.2:26656,38305e1e64fadb5f1201819d4a1d41f4f83cee60@52.77.101.121:26656"
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
- address: sifvalcons1rcwmttxswcc8sqq5rp4ysj4gl8963a0kkfl7kc
  pubkey: sifvalconspub1zcjduepqdgnxx8hl83tukd7ly2dg3jr5rf625na5p4vx0ccxqemkxr7f2wtsahhpyq
  proposerpriority: -375000
  votingpower: 1000000
- address: sifvalcons1hd3u9fy64saqcndxgmwqcqlg8xt67x8vl9wse2
  pubkey: sifvalconspub1zcjduepqj6fmv5pqdg9xfpuydmzv3lupve6cr4uefmfmvpmc54qwhmeh64qqc26amp
  proposerpriority: 125000
  votingpower: 1000000
- address: sifvalcons16n3fc45snzmhp8evw5pwlevff2t2rc6lk5de2g
  pubkey: sifvalconspub1zcjduepqr3mjukr0hzr8xce22lfxq4fhwp8wwflga8jea3yzerf0rhr6m0uq5zf37k
  proposerpriority: -125000
  votingpower: 1000000
- address: sifvalcons1mgqwmk7le4geddltf8m0y2y5ym7sk823n0jdp0
  pubkey: sifvalconspub1zcjduepqg5gwfe0eymzu8zqkuyhkqvfjqyd4d6rslm5qx3guxmj8jzjhrzzqwnu444
  proposerpriority: 375000
  votingpower: 1000000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
a00c8a5f07d87754e9a2d428ad7e1877dbe12ddd@35.165.17.164:26656
231a0b22837c1ad627d34748ac27e21f540dbb87@54.195.149.238:26656
69b0b8fe353ea3e25b23fad18796056bdf5ce9c1@3.24.153.2:26656
38305e1e64fadb5f1201819d4a1d41f4f83cee60@52.77.101.121:26656
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
    --amount="1000000000000000000rowan" \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    ---gas-prices="500000000000000000.0rowan" \
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
