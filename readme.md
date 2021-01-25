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
git checkout tags/merry-go-round-2
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
rake "genesis:sifnode:scaffold[merry-go-round, <moniker>, '<mnemonic>', a75d98a0195596ce7043f7fe14a5498df6279bd3@34.212.71.53:26656, http://34.212.71.53:26657/genesis]"
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
curl http://34.212.71.53:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "a75d98a0195596ce7043f7fe14a5498df6279bd3@34.212.71.53:26656,c4205fd291f3f8d163e0055d859e23ea1b31219a@34.248.75.35:26656,b77e0a8c16462f105669ad5966a10993cbf23205@52.65.129.116:26656,62a91f190ff08861f0602d146c39367b4be7589a@18.140.216.206:26656"
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
- address: sifvalcons1vnys998s8d0xkl9lqvzwn6wp6hyfjlp9k7kyrf
  pubkey: sifvalconspub1zcjduepqw0yxxc28dl2e5ruhzssmpah7kvlxgs0f7m54akcsv9ph529pqxhqmjna3h
  proposerpriority: 125000
  votingpower: 1000000
- address: sifvalcons10qh4zdt9c6htlayjcm96t6umnndhh4049se4cl
  pubkey: sifvalconspub1zcjduepqsu6vj3y2mpzq384q8f3xed85y5upv7gqzscxp0aegax4rg8uh93seezgwk
  proposerpriority: -1125000
  votingpower: 1000000
- address: sifvalcons1sqnsu6zd3tsqah9052xwsxjucauwfcrhx6xn7s
  pubkey: sifvalconspub1zcjduepqu22nqvrq0dntlpa3d8fxx0wq8agc0jr66jkqy7nhc0pr0uujjlmqxam24v
  proposerpriority: -1375000
  votingpower: 1000000
- address: sifvalcons16akkf32nnt44988u5skh0d9llpa2y7tawxumxh
  pubkey: sifvalconspub1zcjduepqk8yrxfvg432qnq489s5pdnq5vg9v4459ewplek6vcmrcl2fp9h3qgha6qw
  proposerpriority: 2375000
  votingpower: 1000000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
a75d98a0195596ce7043f7fe14a5498df6279bd3@34.212.71.53:26656
c4205fd291f3f8d163e0055d859e23ea1b31219a@34.248.75.35:26656
b77e0a8c16462f105669ad5966a10993cbf23205@52.65.129.116:26656
62a91f190ff08861f0602d146c39367b4be7589a@18.140.216.206:26656
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
