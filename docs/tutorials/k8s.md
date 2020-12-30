# Sifchain - Kubernetes (k8s) Tutorial

#### Demo Videos

1. https://youtu.be/dlPLIivwRGg
2. https://youtu.be/ff9CZkmHo3o
3. https://youtu.be/iJjXGXWMfsk

#### Prerequisites / Dependencies:

- Clone the repository (`git clone git@github.com:Sifchain/sifnode.git`)
- [Ruby 2.7.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install)
- [AWS CLI Tool](https://aws.amazon.com/cli/)
- [kubectl](https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html)
- [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)

_This tutorial assumes that you have at least a basic understanding of setting up AWS and configuring your access keys accordingly, so that you may interact with AWS via the CLI Tool._

#### What is Kubernetes? (k8s)

Kubernetes is an open-source container-orchestration system for automating application deployment, scaling, and management.

## Scaffold and deploy a new cluster

1. Switch to the root of the sifchain project.

2. Scaffold a new cluster:

```
rake "cluster:scaffold[<chainID>,<provider>]"
```

e.g.:

```
rake "cluster:scaffold[merry-go-round,aws]"
```

where:

|Param|Description|
|-----|----------|
|`<chainID>`|The Chain ID of the network (e.g.: merry-go-round).|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|

3. Once complete, you'll notice that several Terraform files/folders have been setup inside of the `.live` directory. We recommend you leave the defaults as-is, but for those that have experience with Terraform, feel free to adjust the configuration as you see fit.

4. Deploy the cluster to AWS:

```
rake "cluster:deploy[<chainID>,<provider>]"
```

e.g.:

```
rake "cluster:deploy[merry-go-round,aws]"
```

5. Once complete, you should see your cluster on your AWS account. You can also check using `kubectl`:

```
kubectl get pods --all-namespaces --kubeconfig ./.live/sifchain-aws-merry-go-round/kubeconfig_sifchain-aws-merry-go-round
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
rake "cluster:sifnode:deploy:peer[<chainID>,<provider>,<namespace>,<image>,<image tag>,<moniker>,<mnemonic>,<peer address>,<genesis URL>]"
```

where:

|Param|Description|
|-----|----------|
|`<chainID>`|The Chain ID of the network (e.g.: merry-go-round).|
|`<provider>`|The cloud provider to use (currently only AWS is supported).|
|`<namespace>`|The Kubernetes namespace to use (e.g.: sifnode).|
|`<image>`|The image to pull down from Docker Hub (e.g.: sifchain/sifnoded).|
|`<image tag>`|The image tag to use (e.g.: merry-go-round).|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<peer address>`|The address of the peer to connect to.|
|`<genesis URL>`|The URL of genesis file for the network.|

e.g.:

```
rake "cluster:sifnode:deploy:peer[merry-go-round,aws,sifnode,sifchain/sifnoded,merry-go-round-1,my-node,'my mnemonic',ff0dd55dffa0e67fe21e2c85c80b0c2894bf2586@52.89.19.109:26656,http://52.89.19.109:26657/genesis]"
```

5. Once deployed, check the status of the pods:

```
kubectl get pods -n sifnode --kubeconfig ./.live/sifchain-aws-merry-go-round/kubeconfig_sifchain-aws-merry-go-round
```

and you should see something that resembles the following:

```                            
NAME                           READY   STATUS     RESTARTS   AGE
sifnode-75464fcd4c-dsmzq       0/1     Init:0/1   0          10s
sifnode-cli-67bcfd4b54-mhdjx   0/1     Running    0          10s
```

_It may take several minutes for your node to become active._

6. Once your node is active (Status of "Running"), you can view it's sync status by looking at the logs. Run:

```
kubectl -n sifnode logs <pod> --kubeconfig ./.live/sifchain-aws-merry-go-round/kubeconfig_sifchain-aws-merry-go-round
```

e.g.:

```
kubectl -n sifnode logs sifnode-65fbd7798f-6wqhb --kubeconfig ./.live/sifchain-aws-merry-go-round/kubeconfig_sifchain-aws-merry-go-round
```

## Stake to become a validator

In order to become a validator, that is a node which can participate in consensus on the network, you'll need to stake `rowan`.

1. If using testnet, obtain funds from the faucet.

2. Get the public key of your node:

```
rake "validator:expose:pub_key[<chainID>,<provider>,<namespace>]"
```

e.g.:

```
rake "validator:expose:pub_key[merry-go-round,aws,sifnode]"
```

3. Stake:

```
rake "validator:stake[<chainID>,<moniker>,<amount>,<public key>,<node RPC address>]"
```

where:

|Param|Description|
|-----|----------|
|`<chainID>`|The Chain ID of the network (e.g.: merry-go-round).|
|`<moniker>`|The moniker or name of your node as you want it to appear on the network.|
|`<amount>`|The amount to stake, including the denomination (e.g.: 100000000rowan).|
|`<public key>`|The public key of your validator (you got this in the previous step).|
|`<node RPC address>`|The address to broadcast the transaction to (e.g.: tcp://<node IP address>:26657).|

e.g.:

```
rake "validator:stake[merry-go-round,my-node,10000000rowan,<public key>,tcp://52.89.19.109:26657]"
```

4. It may take several blocks before your node appears as a validator on the network, but you can always check by running:

```
sifnodecli q tendermint-validator-set --node <node RPC address> --trust-node
```

e.g.:

```
sifnodecli q tendermint-validator-set --node tcp://52.89.19.109:26657 --trust-node
```
