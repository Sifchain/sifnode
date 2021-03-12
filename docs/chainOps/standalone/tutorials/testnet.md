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
sifnodecli q tendermint-validator-set --node tcp://rpc-merry-go-round.sifchain.finance:80 --trust-node
```

and you should see the following primary validator node/s for Sifchain:

```
validators:
- address: sifvalcons19vn3vtccx7amdgfk07ne4x2jj87fpsgwusrpcm
  pubkey: sifvalconspub1zcjduepqrczw0kfj0rm3u094lvk7le4u6enrdtm5wtjn8k207gnvtwerlcmshk3lvk
  proposerpriority: -1125
  votingpower: 1000
- address: sifvalcons10mmxf7rmtw2u8lr0d4efzjk4s9q4u8pepn9r95
  pubkey: sifvalconspub1zcjduepq2h0eupzq25g7g6s4uwwrwqaex0m37mhfdc9gh3rtxae0aylg7scq3r93et
  proposerpriority: 125
  votingpower: 1000
- address: sifvalcons1sm9s7f0arcmveuqfx6jtjjhhedj58q3syjf68z
  pubkey: sifvalconspub1zcjduepq98kugnwx40hfjy895n40ex6dwaj02zh06pj6zdjash0nj3wy4v0qhq9x2l
  proposerpriority: 625
  votingpower: 1000
- address: sifvalcons15mnhwsqlkyaacjthmax5zn69ukpvg0wr7hhh87
  pubkey: sifvalconspub1zcjduepq2kd405jrgsjtjmjmx038er8t4l0mcw808z8ka4a2t8q5d6dfxuqsrmdhux
  proposerpriority: 375
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
    --node tcp://rpc-merry-go-round.sifchain.finance:80
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
    --node tcp://rpc-merry-go-round.sifchain.finance:80
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer-merry-go-round.sifchain.finance|
|RPC|https://rpc-merry-go-round.sifchain.finance|
|API|https://lcd-merry-go-round.sifchain.finance|
