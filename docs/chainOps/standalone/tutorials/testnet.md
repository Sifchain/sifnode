# Connecting to the Merry-go-Round Testnet. 

## Scaffold and run your node

1. Clone the respository.

```
git clone ssh://git@github.com/Sifchain/sifnode && cd sifnode
```

2. Checkout the testnet release:

```
git checkout tags/testnet-genesis
```

3. Build:

```
make clean install
```

4. Generate a mnemonic:

```
rake "keys:generate:mnemonic"
```

5. Scaffold your node:

```
rake "genesis:sifnode:scaffold[merry-go-round, <moniker>, '<mnemonic>', e99deeec54ca1c477f8826801bc1fd29f5539a45@44.226.150.203:26656, http://44.226.150.203:26657/genesis]"
```

Where:

|Param|Description|
|-----|----------|
|`<moniker>`|A moniker (name) for your node.|
|`<mnemonic>`|The mnemonic phrase generated in the previous step.|

This step will also output the keyring password, so please record this and the moniker somewhere secure.

7. Connect:

```
rake "genesis:sifnode:boot[<gas_price>]"
```

Where:

|Param|Description|
|-----|----------|
|`<gas_price>`|The minimum gas price (e.g.: 0.5rowan).|

e.g.:

```
rake "genesis:sifnode:boot[0.5rowan]"
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

## Verify

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
  votingpower: 1000
- address: sifvalcons1xnn99v7n50m79pxmveeqdj0f6xuefhn4vay6d6
  pubkey: sifvalconspub1zcjduepqpug6d7p5gt69j9jzpltxlheukz86u3emq2qe89d5jlx7cctfyceq4warur
  proposerpriority: -2125000
  votingpower: 1000
- address: sifvalcons1jycgr4z4trfrwk3x04x3nssyv6uyzf2yzen9j3
  pubkey: sifvalconspub1zcjduepq6h5q8vf9z7h5kt9d986vl7x6x0zx9w9x805f9wn5ys5rw4mme0yqxddd2g
  proposerpriority: -875000
  votingpower: 1000
- address: sifvalcons1hdh8xnvg3ckrl0pjnryc225qdg0rr0rgsdmqwz
  pubkey: sifvalconspub1zcjduepq69qwgexfxxfh9jtp4nu0jrvf4hq8m9k4079z995ss0qx0qp5qd7sln7cm2
  proposerpriority: 1375000
  votingpower: 1000
```

Congratulations. You are now connected to the network.

## Additional Peers

The following can be used as additional peers on the network:

```
e99deeec54ca1c477f8826801bc1fd29f5539a45@44.226.150.203:26656
fb4bfcaf9980a2ee3fe8298eadd21c9757d83f6c@52.49.165.39:26656
39172b7f5f8c2394af86f3174e4b8c9f6eb3ad3b@13.210.25.108:26656
e032bdcfad831c7c7131e02453fdc736625547a6@18.136.70.31:26656
```

## Become a Validator

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
    --amount="<amount>" \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=<moniker> \
    --keyring-backend=file
```

Where:

|Param|Description|
|-----|----------|
|`<amount>`|The amount of rowan you wish to stake (e.g.: 10000000000000000000rowan).|
|`<moniker>`|The moniker (name) of your node.|

e.g.:

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
    --from=my-node \
    --keyring-backend=file
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer-merry-go-round.sifchain.finance|
|RPC|https://rpc-merry-go-round.sifchain.finance|
|API|https://lcd-merry-go-round.sifchain.finance|
