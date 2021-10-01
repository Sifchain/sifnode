# Connecting to the Sifchain TestNet.

## Scaffold and run your node

1. Switch to the root of the sifchain project.

2. Install bundler/gems:

```bash
make -C ./deploy bundler
```

3. Ensure your `$GOPATH` is setup correctly:

```bash
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
```

4. Compile `sifnoded` and `sifgen`:

```bash
make install
```

5. Generate a new mnemonic key for your node. This key is what your node will use to eventually sign transactions/blocks on the network.

```bash
rake "sifnode:keys:generate:mnemonic"
```

6. Import your newly generated key:

```bash
rake "sifnode:keys:import[<moniker>]"
```

where:

|Param|Description|
|-----|----------|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|

e.g.:

```bash
rake "sifnode:keys:import[my-node]"
```

7. Check that it's been imported accordingly:

```bash
rake "sifnode:keys:show[<moniker>]"
```

where:

|Param|Description|
|-----|----------|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|

e.g.:

```bash
rake "sifnode:keys:show[my-node]"
```

8. Boot your node:

```bash
rake "sifnode:standalone:boot[<chain_id>,<moniker>,'<mnemonic>',<gas_price>,<bind_ip_address>,'<flags>']"
```

Where:

|Param|Description|
|-----|----------|
|`<chain_id>`|The Chain ID of the network (e.g.: `sifchain-testnet-1`).|
|`<moniker>`|A name for your node.|
|`<mnemonic>`|The mnemonic phrase generated in the previous step.|
|`<gas_price>`|The minimum gas price (e.g.: `0.5rowan`).|
|`<bind_ip_address>`|The IP Address to bind to (*Important:* this is what your node will advertise to the rest of the network). This should be the public IP of the host.|
|`<flags>`|Optional. Docker compose run flags (see [here](https://docs.docker.com/compose/reference/run/)).|

e.g.:

```bash
rake "sifnode:standalone:boot[sifchain-testnet-1,my-node,'my mnemonic',0.5rowan,127.0.0.1]"
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

## Stake to become a validator

In order to become a validator, that is a node which can participate in consensus on the network, you'll need to stake `rowan`.

1. Get the public key of your node:

```bash
rake "sifnode:keys:docker:public[<image>,<image_tag>]"
```

where:

|Param|Description|
|-----|----------|
|`<image>`|The docker image your node is running  (e.g.: `sifchain/sifnoded`).|
|`<image_tag>`|The tag of the docker image your node is running  (e.g.: `sifchain-testnet-1`).|

e.g.:

```bash
rake "sifnode:keys:docker:public[sifchain/sifnoded, sifchain-testnet-1]"
```

2. Stake:

```bash
rake "sifnode:staking:stake[<commission_max_change_rate>,<commission_max_rate>,<commission_rate>,<chain_id>,<moniker>,<amount>,<gas>,<gas_prices>,<public_key>,<node>]"
```

where:

|Param|Description|
|-----|----------|
|`<commission_max_change_rate>`|The maximum commission change rate percentage (per day).|
|`<commission_max_rate>`|The maximum commission rate percentage.|
|`<commission_rate>`|The initial commission rate percentage.|
|`<chain_id>`|The Chain ID of the network (e.g.: `sifchain-testnet-1`).|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<amount>`|The amount to stake, including the denomination (e.g.: `100000000rowan`). The precision used is 1e18.|
|`<gas>`| The per-transaction gas limit (e.g.: `300000`).|
|`<gas_prices>`|The minimum gas price to use  (e.g.: `0.5rowan`).|
|`<public_key>`|The public key of your validator (you got this in the previous step).|
|`<node_rpc_address>`|The address to broadcast the transaction to (e.g.: `tcp://rpc-testnet.sifchain.finance:80`).|

e.g.:

```bash
rake "sifnode:staking:stake[0.1,0.1,0.1,sifchain-testnet-1,my-node,10000000rowan,300000,0.5rowan,<public_key>,tcp://rpc-testnet.sifchain.finance:80]"
```

3. It may take several blocks before your node appears as a validator on the network, but you can always check by running:

```bash
rake "sifnode:staking:validators"
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer-testnet.sifchain.finance|
|RPC|https://rpc-testnet.sifchain.finance|
|API|https://api-testnet.sifchain.finance|
