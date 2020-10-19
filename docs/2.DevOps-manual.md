# Sifnode DevOps

## Terraform Cluster Config

Sifnode currently only supports AWS.

### AWS

#### Generate

To generate the necessary cluster config for AWS, simply run (from the project root)

```
make PROVIDER=aws CHAINNET=<chain> create-new-environment
```

Where:

| Variable | Description |
| ---------| ----------------|
| `<chain>` | The name of the chain/network this is for (e.g.: `testnet`, `chaosnet`, `mainnet`) |

Other variables that can be passed to the `create-new-environment` make target, include:

| Variable | Description |
| ---------| ----------------|
| `STATE_BACKEND` | The Terraform state backend (one of `local`, `s3` or `remote`) |
| `AWS_REGION` | The region in which the cluster will be deployed (defaults to `us-west-1`) |
| `AWS_STATE_BUCKET_NAME` | If `STATE_BACKEND=s3`, then you will need to provide the name of the s3 bucket to use. |
| `AWS_CIDR` | The CIDR (classless inter-domain routing) block to use for the cluster. Defaults to `10.0.0.0/16` |
| `AWS_CLUSTER_VERSION` | The cluster version. Defaults to `1.17`. |
| `AWS_CLUSTER_NAME` | The cluster name. Defaults to `sifnode`. |
| `AWS_CLUSTER_AMI_TYPE` | The cluster AMI (Amazon Machine Image) type. Defaults to `AL2_x86_64`. |
| `AWS_CLUSTER_DESIRED_CAPACITY` | The desired number of nodes. Defaults to `1`. |
| `AWS_CLUSTER_MAX_CAPACITY` | The maximum number of nodes. Defaults to `1`. |
| `AWS_CLUSTER_MIN_CAPACITY` | The minimum number of nodes. Defaults to `1`. |
| `AWS_CLUSTER_INSTANCE_TYPE` | The instance type. Defaults to `t2.small`. |
| `AWS_CLUSTER_DISK_SIZE` | The default disk size (in GB). Defaults to `100`. |

#### Apply

The generation step above will output a compiled Terraform config into the following directory:

`.live/sifchain-aws-<chain>`

Simply switch to that directory and run:

`terraform init && terraform apply`

to deploy the cluster.
