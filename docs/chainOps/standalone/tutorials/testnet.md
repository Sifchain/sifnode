# Connecting to the Sifchain BetaNet.

## Prerequisites / Dependencies:

- [Docker](https://www.docker.com/get-started)
- [Ruby 2.7.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install)
  - Add `export GOPATH=~/go` to your shell
  - Add `export PATH=$PATH:$GOPATH/bin` to your shell

## Scaffold and run your node

1. Clone the repository:

```
git clone https://github.com/Sifchain/sifnode && cd sifnode
```

2. Build:

```
make clean install
```

3. Generate a mnemonic (if you don't already have one):

```
rake "keys:generate:mnemonic"
```

4. Boot your node:

```
rake "genesis:sifnode:boot[testnet,<moniker>,'<mnemonic>',<gas_price>,<bind_ip_address>]"
```

Where:

|Param|Description|
|-----|----------|
|`<moniker>`|A name for your node.|
|`<mnemonic>`|The mnemonic phrase generated in the previous step.|
|`<gas_price>`|The minimum gas price (e.g.: 0.5rowan).|
|`<bind_ip_address>`|The IP Address to bind to (*Important:* this is what your node will advertise to the rest of the network). This should be the public IP of the host.|

and your node will start synchronizing with the network. Please note that this may take several hours or more.

## Verify

You can verify that you're connected by running:

```
sifnodecli q tendermint-validator-set --node tcp://100.20.201.226:26657 --trust-node
```

and you should see the following primary validator node/s for Sifchain:

```
validators:
- address: sifvalcons1g4jce0n7d0zjrckkcysqya4q76yzls7xnhwsx0
  pubkey: sifvalconspub1zcjduepq08swr6w2q2c6utx5j2wm8xfezy3ayluq683dqwj4w358zr3qd95sfaq99k
  proposerpriority: -207656
  votingpower: 1000000
- address: sifvalcons1f9z0jnz96f63enn2c980a6qjwvtxjsmczafk9n
  pubkey: sifvalconspub1zcjduepqquanx4sm7wf07lm76fv5sqwawsl4tnw7ny9m86kmhw7ekkutjrwq8qw8my
  proposerpriority: 1107281
  votingpower: 1000000
- address: sifvalcons1c7rdw0kuqt89s5p4nywczpxraregz3lgsv77g8
  pubkey: sifvalconspub1zcjduepq9nhfkyft4fcvtuac8we7f9hyqxtf70hjyfzwwj6h6g23pdyjg4rsd369d2
  proposerpriority: 232906
  votingpower: 1000000
- address: sifvalcons1c75nmvdge0sxdg4ep3v3l60rrmj5dm9vurpv47
  pubkey: sifvalconspub1zcjduepq03aglq7ux9e5q8sqezhnjrfekslgqhwypzh9myy83wapt7h2ckus5wzvjs
  proposerpriority: -1132531
  votingpower: 1000
```

Congratulations! You are now connected to the network.

## Become a Validator

You won't be able to participate in consensus until you become a validator.

1. Import your mnemonic locally:

```
rake "keys:import[<moniker>]"
```

Where:

|Param|Description|
|-----|----------|
|`<moniker>`|A name for your node.|

*You will need to have tokens (rowan) on your account in order to become a validator.*

2. From within your running container, obtain your node's public key:

```
sifnoded tendermint show-validator
```

3. Run the following command to become a validator: 

```
sifnodecli tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount="<amount>" \
    --pubkey=<pub_key> \
    --moniker=<moniker> \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=<moniker> \
    --keyring-backend=file \
    --node tcp://100.20.201.226:26657
```

Where:

|Param|Description|
|-----|----------|
|`<amount>`|The amount of rowan you wish to stake (the more the better). The precision used is 1e18.|
|`<pub_key>`|The public key of your node, that you got in the previous step.|
|`<moniker>`|The moniker (name) of your node.|

e.g.:

```
sifnodecli tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount="1000000000000000000000rowan" \
    --pubkey=thepublickeyofyournode \
    --moniker=my-node \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=my-node \
    --keyring-backend=file \
    --node tcp://100.20.201.226:26657
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer-merry-go-round.sifchain.finance|
|RPC|https://rpc-merry-go-round.sifchain.finance|
|API|https://lcd-merry-go-round.sifchain.finance|
