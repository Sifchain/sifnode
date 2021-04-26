# Connecting to the Sifchain BetaNet with Kubernetes (k8s).

## Scaffold and deploy a new cluster

1. Switch to the root of the sifchain project.

2. import gotpl module

```
go get github.com/belitre/gotpl
```

3. Scaffold a new cluster:

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

4. Once complete, you'll notice that several Terraform files/folders have been setup inside of the `.live` directory. We recommend you leave the defaults as-is, but for those that have experience with Terraform, feel free to adjust the configuration as you see fit.

5. Deploy the cluster to AWS:

```
rake "cluster:deploy[<cluster>,<provider>]"
```

e.g.:

```
rake "cluster:deploy[my-cluster,aws]"
```

6. Once complete, you should see your cluster on your AWS account. You can also check using `kubectl`:

```
kubectl get pods --all-namespaces --kubeconfig ./.live/sifchain-aws-my-cluster/kubeconfig_sifchain-aws-my-cluster
```

Note: if you get 
```
Unable to connect to the server: getting credentials: exec: exec: "aws-iam-authenticator": executable file not found in $PATH
```
Install aws-iam-authenticator from
```
https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
```

## Deploy a new node

1. Generate a new mnemonic key for your node. This key is what your node will use to eventually sign transactions/blocks on the network.

```
rake "keys:generate:mnemonic"
```

Note: if you get _rake abort!_ error

run these commands
```
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
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
|`<chain_id>`|The Chain ID of the network (e.g.: sifchain).|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|
|`<namespace>`|The Kubernetes namespace to use (e.g.: sifnode).|
|`<image>`|The image to pull down from Docker Hub (e.g.: sifchain/sifnoded).|
|`<image_tag>`|The image tag to use (this must be `mainnet-genesis`).|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<peer_address>`|The address of the peer to connect to.|
|`<genesis_url>`|The URL of genesis file for the network.|

e.g.:

```
rake "cluster:sifnode:deploy:peer[my-cluster,sifchain,aws,sifnode,sifchain/sifnoded,mainnet-genesis,my-node,'my mnemonic',0d4981bdaf4d5d73bad00af3b1fa9d699e4d3bc0@44.235.108.41:26656,https://rpc.sifchain.finance/genesis]"
```

_Please note: the image tag *must* be `mainnet-genesis`._

5. Once deployed, check the status of the pods:

```
kubectl get pods -n sifnode --kubeconfig ./.live/sifchain-aws-my-cluster/kubeconfig_sifchain-aws-my-cluster
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
rake "validator:keys:public[<cluster>,<provider>,<namespace>]"
```

e.g.:

```
rake "validator:keys:public[my-cluster,aws,sifnode]"
```
Note: This requires jq JSON processor if not installed install with ```sudo apt-get install jq```

3. Run the following command to become a validator:

sifnodecli tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount="<amount>" \
    --pubkey=<pub_key> \
    --moniker=<moniker> \
    --chain-id=sifchain \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=<moniker> \
    --keyring-backend=file \
    --node tcp://rpc.sifchain.finance:80
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
    --chain-id=sifchain \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=my-node \
    --keyring-backend=file \
    --node tcp://rpc.sifchain.finance:80
```

4. It may take several blocks before your node appears as a validator on the network, but you can always check by running:

```
sifnodecli q tendermint-validator-set --node <node_rpc_address> --trust-node
```

e.g.:

```
sifnodecli q tendermint-validator-set --node tcp://rpc.sifchain.finance:80 --trust-node
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer.sifchain.finance|
|RPC|https://rpc.sifchain.finance|
|API|https://api.sifchain.finance|
