# Connecting to the Merry-go-Round Testnet with Kubernetes (k8s).

## Demo Videos

1. https://youtu.be/dlPLIivwRGg
2. https://youtu.be/ff9CZkmHo3o
3. https://youtu.be/iJjXGXWMfsk

## Scaffold and deploy a new cluster

1. Switch to the root of the sifchain project.

2. Scaffold a new cluster:

```
rake "cluster:scaffold[<cluster>,<provider>]"
```

e.g.:

```
rake "cluster:scaffold[my-cluster,aws]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|A name for your new cluster.|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|

3. Once complete, you'll notice that several Terraform files/folders have been setup inside of the `.live` directory. We recommend you leave the defaults as-is, but for those that have experience with Terraform, feel free to adjust the configuration as you see fit.

4. Deploy the cluster to AWS:

```
rake "cluster:deploy[<cluster>,<provider>]"
```

e.g.:

```
rake "cluster:deploy[my-cluster,aws]"
```

5. Once complete, you should see your cluster on your AWS account. You can also check using `kubectl`:

```
kubectl get pods --all-namespaces --kubeconfig ./.live/sifchain-aws-my-cluster/kubeconfig_sifchain-aws-my-cluster
```

## Deploy a new node

1. Generate a new mnemonic key for your node. This key is what your node will use to eventually sign transactions/blocks on the network.

```
rake "keys:generate:mnemonic"
```

2. Import your newly generated key:

```
rake "keys:import[<moniker>]"
```

where:

|Param|Description|
|-----|----------|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|

e.g.:

```
rake "keys:import[my-node]"
```

3. Check that it's been imported accordingly:

```
sifnodecli keys show <moniker> --keyring-backend file 
```

4. Deploy a new node to your cluster and connect to an existing network:

```
rake "cluster:sifnode:deploy:peer[<cluster>,<chain_id>,<provider>,<namespace>,<image>,<image_tag>,<moniker>,<mnemonic>,<peer_address>,<genesis_url>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your cluster.|
|`<chain_id>`|The Chain ID of the network (e.g.: merry-go-round).|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|
|`<namespace>`|The Kubernetes namespace to use (e.g.: sifnode).|
|`<image>`|The image to pull down from Docker Hub (e.g.: sifchain/sifnoded).|
|`<image_tag>`|The image tag to use (this must be `testnet-genesis`)..|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<peer_address>`|The address of the peer to connect to.|
|`<genesis_url>`|The URL of genesis file for the network.|

e.g.:

```
rake "cluster:sifnode:deploy:peer[my-cluster,merry-go-round,aws,sifnode,sifchain/sifnoded,testnet-genesis,my-node,'my mnemonic',f214ec6828b85793289fcb0b025bc260747983f0@100.20.201.226:26656,http://100.20.201.226:26657/genesis]"
```

_Please note: the image tag *must* be `testnet-genesis`._

5. Once deployed, check the status of the pods:

```
kubectl get pods -n sifnode --kubeconfig ./.live/sifchain-aws-merry-go-round/kubeconfig_sifchain-aws-merry-go-round
```

and you should see something that resembles the following:

```                            
NAME                           READY   STATUS     RESTARTS   AGE
sifnode-75464fcd4c-dsmzq       0/1     Init:0/2   0          10s
```

_It may take several minutes for your node to become active._

6. Once your node is active (Status of "Running"), you can view it's sync status by looking at the logs. Run:

```
kubectl -n sifnode logs <pod> --kubeconfig ./.live/sifchain-aws-my-cluster/kubeconfig_sifchain-aws-my-cluster
```

e.g.:

```
kubectl -n sifnode logs sifnode-65fbd7798f-6wqhb --kubeconfig ./.live/sifchain-aws-my-cluster/kubeconfig_sifchain-aws-my-cluster
```

## Stake to become a validator

In order to become a validator, that is a node which can participate in consensus on the network, you'll need to stake `rowan`.

1. If using testnet, obtain funds from the faucet.

2. Get the public key of your node:

```
rake "validator:expose:pub_key[<cluster>,<provider>,<namespace>]"
```

e.g.:

```
rake "validator:expose:pub_key[my-cluster,aws,sifnode]"
```

3. Stake:

```
rake "validator:stake[<chain_id>,<moniker>,<amount>,<public_key>,<node_rpc_address>]"
```

where:

|Param|Description|
|-----|----------|
|`<chain_id>`|The Chain ID of the network (e.g.: merry-go-round).|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<amount>`|The amount to stake, including the denomination (e.g.: 100000000rowan). The precision used is 1e18.|
|`<gas>`|The gas price (e.g.: 0.5rowan).|
|`<public_key>`|The public key of your validator (you got this in the previous step).|
|`<node_rpc_address>`|The address to broadcast the transaction to (e.g.: tcp://<node IP address>:26657).|

e.g.:

```
rake "validator:stake[merry-go-round,my-node,10000000rowan,0.5rowan,<public_key>,tcp://100.20.201.226:26657]"
```

4. It may take several blocks before your node appears as a validator on the network, but you can always check by running:

```
sifnodecli q tendermint-validator-set --node <node_rpc_address> --trust-node
```

e.g.:

```
sifnodecli q tendermint-validator-set --node tcp://100.20.201.226:26657 --trust-node
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer-merry-go-round.sifchain.finance|
|RPC|https://rpc-merry-go-round.sifchain.finance|
|API|https://lcd-merry-go-round.sifchain.finance|
