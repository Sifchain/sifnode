# Connecting to the Sifchain TestNet.

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
rake "genesis:sifnode:boot[testnet,<moniker>,'<mnemonic>',<gas_price>,<bind_ip_address>,'<flags>']"
```

Where:

|Param|Description|
|-----|----------|
|`<moniker>`|A name for your node.|
|`<mnemonic>`|The mnemonic phrase generated in the previous step.|
|`<gas_price>`|The minimum gas price (e.g.: 0.5rowan).|
|`<bind_ip_address>`|The IP Address to bind to (*Important:* this is what your node will advertise to the rest of the network). This should be the public IP of the host.|
|`<flags>`|Optional. Docker compose run flags (see [here](https://docs.docker.com/compose/reference/run/)).|

and your node will start synchronizing with the network. Please note that this may take several hours or more.

## Verify

You can verify that you're connected by running:

```
sifnoded q tendermint-validator-set --node tcp://rpc-testnet.sifchain.finance:80 --trust-node
```

and you should see the following primary validator node/s for Sifchain:

```
validators:
- address: sifvalcons1w8gnu0k86dxs0nsjff7hh77wzhaxmmdcwf5krj
  pubkey: sifvalconspub1zcjduepqvwq9cv5j3kv7x23g3zhd2qhcwp53quq6d72mr4f53xq0fpaqx0psvrmwxz
  proposerpriority: 1125
  votingpower: 1000
- address: sifvalcons1j3z8x9f3x5zxpdzdk9e08kqvlke7y5tsshy9tv
  pubkey: sifvalconspub1zcjduepq54s2l4facx6l3g9rt8jwuhwxfealmt6d2r6v7wn7sthfgwx9n6rq9alcm3
  proposerpriority: -1625
  votingpower: 1000
- address: sifvalcons1hyw42srq4u7766y9zqx4vfurfcasjwwulllll7
  pubkey: sifvalconspub1zcjduepqjdmaurc6v3ueskv7cgrzdw7phhzeg904akj3agz4yw87juq6yknsyw8qkf
  proposerpriority: 1625
  votingpower: 1000
- address: sifvalcons1ufd9g9txtz0vptflty94ey6p97rh6l7q33w9hc
  pubkey: sifvalconspub1zcjduepqcs85pl337wf2da6ucuu6vcs8l9xwtp2cv05swy7l8csr6uzs5vcq37pc07
  proposerpriority: -1125
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
docker exec -ti mainnet_sifnode_1 sh
/root/.sifnoded/cosmovisor/current/bin/sifnoded tendermint show-validator
```

3. Run the following command to become a validator: 

```
sifnoded tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount="<amount>" \
    --pubkey=<pub_key> \
    --moniker=<moniker> \
    --chain-id=sifchain-testnet \
    --min-self-delegation="1" \
    --gas="300000" \
    --gas-prices="0.5rowan" \
    --from=<moniker> \
    --keyring-backend=file \
    --node tcp://rpc-testnet.sifchain.finance:80
```

Where:

|Param|Description|
|-----|----------|
|`<amount>`|The amount of rowan you wish to stake (the more the better). The precision used is 1e18.|
|`<pub_key>`|The public key of your node, that you got in the previous step.|
|`<moniker>`|The moniker (name) of your node.|

e.g.:

```
sifnoded tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount="1000000000000000000000rowan" \
    --pubkey=thepublickeyofyournode \
    --moniker=my-node \
    --chain-id=sifchain-testnet \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=my-node \
    --keyring-backend=file \
    --node tcp://rpc-testnet.sifchain.finance:80
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer-testnet.sifchain.finance|
|RPC|https://rpc-testnet.sifchain.finance|
|API|https://api-testnet.sifchain.finance|
