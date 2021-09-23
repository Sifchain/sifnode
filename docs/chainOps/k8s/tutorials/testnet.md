# Connecting to the Sifchain TestNet with Kubernetes (k8s).

## Scaffold and deploy a new cluster

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

4. Compile `sifnoded`:

```bash
make install
```

5. Scaffold a new cluster:

```bash
rake "cluster:scaffold[<cluster>,<provider>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|A name for your new cluster.|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|

e.g.:

```bash
rake "cluster:scaffold[my-cluster,aws]"
```

6. Once complete, you'll notice that several Terraform files/folders have been setup inside of a `.live` directory. We recommend you leave the defaults as-is, but for those that have experience with Terraform, feel free to adjust the configuration as you see fit.

7. Deploy the cluster to AWS (Default region is US-WEST-2):

```bash
rake "cluster:deploy[<cluster>,<provider>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your new cluster.|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|

e.g.:

```bash
rake "cluster:deploy[my-cluster,aws]"
```

Once complete, you should see your cluster on your AWS account.

8. Update the cluster kubeconfig (so you may interact with the cluster properly):

```bash
rake "provider:aws:kubeconfig[<cluster>,<region>,<aws_profile>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your cluster.|
|`<region>`|The AWS region where your cluster was deployed.|
|`<aws_profile>`|*Optional. Your AWS profile (used when you have multiple AWS accounts setup locally).|

e.g.:

```bash
rake "provider:aws:kubeconfig[my-cluster,us-west-2]"
```

## Deploy a new node

1. Generate a new mnemonic key for your node. This key is what your node will use to eventually sign transactions/blocks on the network.

```bash
rake "sifnode:keys:generate:mnemonic"
```

2. Import your newly generated key:

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

3. Check that it's been imported accordingly:

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

4. Deploy a new node to your cluster and connect to an existing network:

```bash
rake "sifnode:peer:deploy[<cluster>,<chain_id>,<provider>,<namespace>,<image>,<image_tag>,<moniker>,<mnemonic>,<peer_address>,<genesis_url>,<enable_rpc>,<enable_external_rpc>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your cluster.|
|`<chain_id>`|The Chain ID of the network (e.g.: `sifchain-testnet-1`).|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|
|`<namespace>`|The Kubernetes namespace to use (e.g.: `sifnode`).|
|`<image>`|The image to pull down from Docker Hub (e.g.: `sifchain/sifnoded`).|
|`<image_tag>`|The image tag to use (this must be `testnet-0.9.0`).|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<peer_address>`|The address of the peer to connect to.|
|`<genesis_url>`|The URL of genesis file for the network.|
|`<enable_rpc>`|Enable RPC (if unsure, set to `false`).|
|`<enable_external_rpc>`|Enable external access to the RPC port. If `true`, then `<enable_rpc>` must also be `true`. Recommend this is set to `false`.|

e.g.:

```bash
rake "sifnode:peer:deploy[my-cluster,sifchain-testnet-1,aws,sifnode,sifchain/sifnoded,testnet-0.9.0,my-node,'my mnemonic',a2864737f01d3977211e2ea624dd348595dd4f73@3.222.8.87:26656,https://rpc-testnet.sifchain.finance/genesis,false,false]"
```

5. Once deployed, check the status of the pods:

```bash
rake "cluster:pods[<cluster>,<provider>,<namespace>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your cluster.|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|
|`<namespace>`|The namespace you want to check (e.g.: `sifnode`).|

e.g.:

```bash
rake "cluster:pods[my-cluster,aws,sifnode]"
```

and you should see something that resembles the following:

```bash                    
NAME                           READY   STATUS     RESTARTS   AGE
sifnode-75464fcd4c-dsmzq       0/1     Init:0/2   0          10s
```

_It may take several minutes for your node to become active._

6. Once your node is active (Status of "Running"), you can view its sync status as follows:

```bash
rake "sifnode:status[<cluster>,<namespace>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your cluster.|
|`<namespace>`|The namespace you want to check (e.g.: `sifnode`).|

e.g.:

```bash
rake "sifnode:status[my-cluster,sifnode]"
```

## Stake to become a validator

In order to become a validator, that is a node which can participate in consensus on the network, you'll need to stake `rowan`.

1. Get the public key of your node:

```bash
rake "sifnode:keys:public[<cluster>,<provider>,<namespace>]"
```

where:

|Param|Description|
|-----|----------|
|`<cluster>`|The name of your cluster.|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|
|`<namespace>`|The namespace you want to check (e.g.: `sifnode`).|

e.g.:

```bash
rake "sifnode:keys:public[my-cluster,aws,sifnode]"
```

3. Stake:

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
rake "sifnode:staking:stake[0.1,0.1,0.1,sifchain-testnet-1,my-node,10000000rowan,300000,0.5rowan,<public_key>,tcp://rpc.sifchain.finance:80]"
```

4. It may take several blocks before your node appears as a validator on the network, but you can always check by running:

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
