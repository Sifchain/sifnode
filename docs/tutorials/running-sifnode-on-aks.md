# Running sifnode on Azure AKS

The deployment CLI selects either Azure AKS or AWS EKS. Each provider has a
Terraform module under `deploy/terraform/providers/`, a scaffold template under
`deploy/terraform/template/`, and uses the same Helm chart. This follows the
historical AWS layout referenced in [issue #503](https://github.com/Sifchain/sifnode/issues/503).
That deployment directory is absent from the current base; these files provide
a focused cluster and node workflow without restoring unrelated old Rake tasks.

This guide starts a single full node. Registering or funding a validator is a
separate operator action. Never run two pods with the same validator signing key.

## Prerequisites

- Python 3.8+, Terraform 1.7+ (below 2.0), Helm 3, and kubectl compatible with
  your chosen Kubernetes version.
- An Azure subscription, Azure CLI authentication, permission to create a resource
  group and AKS cluster, and sufficient regional VM/public IP/disk quota.
- The image version, chain ID, reviewed genesis and peers for the target network.
  Consult the [Sifchain validators repository](https://github.com/Sifchain/sifchain-validators)
  for network-specific instructions. A historical image tag or genesis is not
  necessarily compatible with the current network.
- A cost budget approved by the operator. AKS worker VMs, disks, public IPs and
  load balancers incur charges, including with the `Free` control-plane tier.

The configurations pin AzureRM to `~> 4.56.0`. Follow the
[AzureRM authentication documentation](https://registry.terraform.io/providers/hashicorp/azurerm/4.56.0/docs)
for CLI, workload identity or service-principal authentication. For an existing
Azure CLI session, select the subscription and export its ID:

```sh
az account set --subscription YOUR_SUBSCRIPTION_ID
export ARM_SUBSCRIPTION_ID="$(az account show --query id -o tsv)"
```

Keep authentication outside the repository. Never put passwords, client secrets,
wallet recovery phrases or signing keys in Terraform inputs or Helm values.

## Scaffold and provision AKS

Run these commands from the repository root. Replace the example documentation
CIDR with your real administrator public egress IP/CIDR. It must also cover the
machine running Helm and kubectl; VPN/NAT egress may differ from your local IP.

```sh
python3 deploy/cluster.py scaffold --provider azure --name testnet \
  --location eastus --admin-cidr 192.0.2.1/32
```

`scaffold` only writes files; it does not initialize Terraform or provision
resources. It refuses to overwrite an existing environment. The generated
directory is `.live/sifchain-azure-testnet/`, ignored by Git. Review its `main.tf`
and `terraform.tfvars.json`. Repeat `--admin-cidr` for multiple administrators.
Optionally pass `--kubernetes-version` with a version supported in that region;
otherwise Azure chooses its recommended version when creating the cluster.

The module defaults to one `Standard_D4s_v5` worker with a 100 GB OS disk, managed
identity, Kubernetes RBAC, Azure CNI overlay and the Azure Disk CSI driver. The
API accepts only the configured administrator CIDRs. Adjust the module arguments
in `main.tf` for `node_count`, `vm_size`, `sku_tier` or `tags`. Choose `Standard`
if a control-plane uptime SLA is needed, and size memory/disk for the network.

```sh
python3 deploy/cluster.py init --provider azure --name testnet
python3 deploy/cluster.py validate --provider azure --name testnet
python3 deploy/cluster.py plan --provider azure --name testnet
python3 deploy/cluster.py apply --provider azure --name testnet
python3 deploy/cluster.py kubeconfig --provider azure --name testnet
export KUBECONFIG="$PWD/.live/sifchain-azure-testnet/kubeconfig"
kubectl get nodes
kubectl get storageclass managed-csi-premium
```

`apply` and `destroy` retain Terraform's confirmation prompt. The CLI checks the
recorded provider/name before operating on state and never falls back to your
default Kubernetes context. `kubeconfig` writes the `kubectl_config` output to a
mode-0600 file. AKS uses local cluster credentials; this scaffold does not set up
Entra ID integration. Treat the kubeconfig **and Terraform state** as credentials,
use an access-controlled encrypted remote backend for shared/production state,
and manage cluster access according to your organization policy. Do not commit
the generated directory. Re-scaffold if you move the repository: module paths in
the generated file refer to the checkout that created it.

## Supply genesis and deploy the node

Create the namespace and a ConfigMap from an operator-reviewed `genesis.json`:

```sh
kubectl create namespace sifnode
kubectl -n sifnode create configmap sifnode-genesis --from-file=genesis.json=/path/to/genesis.json
```

ConfigMaps are limited to 1 MiB. For larger genesis files, put `genesis.json` at
the root of a separate, pre-populated PVC in the same namespace, using a storage
class suitable for the node's availability zone. Set `genesis.existingClaim` to
that PVC's name instead of `genesis.existingConfigMap`. The chart mounts it
read-only. Only one of these genesis sources may be configured.

Create `.live/sifchain-azure-testnet/node-values.yaml`, replacing the capitalized
placeholders with the network's actual settings:

```yaml
image:
  repository: YOUR_SIFNODE_IMAGE
  tag: YOUR_NETWORK_VERSION
chain:
  id: YOUR_CHAIN_ID
  moniker: my-aks-node
  persistentPeers: "NODE_ID@PEER_ADDRESS:26656"
genesis:
  existingConfigMap: sifnode-genesis
persistence:
  size: 500Gi
```

The image must contain `sifnoded`, `/bin/sh`, `cp` and `touch`, support the worker
architecture, and run as UID/GID 1000 with a writable `/data` volume. Use a trusted,
versioned image. The chart requires explicit image/genesis/chain settings and
uses `sifnoded init` and `validate-genesis` on first boot. It creates node/signing
keys locally on the PVC; it does not import a mnemonic or stake funds. A reviewed
existing node data PVC can be supplied as `persistence.existingClaim`.

```sh
python3 deploy/cluster.py node --provider azure --name testnet \
  --values .live/sifchain-azure-testnet/node-values.yaml
kubectl -n sifnode get pods,pvc,services
kubectl -n sifnode logs statefulset/sifnode-sifnode -c sifnode --tail=50
```

The node is a single-replica StatefulSet backed by a `ReadWriteOnce`
`managed-csi-premium` disk. Only P2P port 26656 is exposed by the public
LoadBalancer. RPC port 26657 is a separate ClusterIP service. After the P2P service
receives an IP, set `chain.externalAddress: "PUBLIC_IP:26656"` in the values file
and rerun the `node` command to advertise it to peers. Inspect synchronization
through a local port-forward:

```sh
kubectl -n sifnode port-forward service/sifnode-sifnode-rpc 26657:26657
# In another terminal:
curl --fail http://127.0.0.1:26657/status
```

On restarts, the initialization marker preserves existing genesis, keys and
consensus state. Changing the genesis ConfigMap does not replace a running
chain's genesis. Restore/migrate an existing database only after reviewing its
genesis, chain ID and signing state; the bootstrap script refuses an unmarked
existing application database. Do not manually remove the marker on a live node.
Network upgrades and snapshot imports require the network's operator procedure.

## AWS selection

The same commands accept `--provider aws`, for example:

```sh
python3 deploy/cluster.py scaffold --provider aws --name testnet \
  --location us-west-2 --admin-cidr 192.0.2.1/32
python3 deploy/cluster.py init --provider aws --name testnet
python3 deploy/cluster.py plan --provider aws --name testnet
```

Authenticate using your normal AWS CLI/SDK credential chain outside the repo.
The AWS scaffold follows the historical VPC + EKS module pattern, with pinned
VPC/EKS module versions and AWS provider 5.100.x. It provisions a public-subnet
EKS node group, an administrator CIDR-restricted API, and EBS CSI. Its default
`t4g.xlarge` workers are **ARM64**, so use an ARM64 or multi-architecture sifnode
image. The generated kubeconfig runs `aws eks get-token` using your active AWS
credentials. The node role has the EBS CSI policy; operators needing stricter
separation can configure a dedicated CSI service-account role. The chart uses a
gp3 CSI StorageClass on AWS, without applying AWS annotations to Azure services.
Continue with `apply`, `kubeconfig`, genesis setup and `node` using `--provider aws`
throughout. The AWS environment has separate state and kubeconfig from Azure.

## Teardown and troubleshooting

```sh
python3 deploy/cluster.py uninstall --provider azure --name testnet
kubectl -n sifnode get pvc
# Back up/review the data before explicitly deleting any PVC.
python3 deploy/cluster.py destroy --provider azure --name testnet
```

Helm uninstall leaves StatefulSet data claims in place. Azure's built-in storage
class normally deletes a disk when its PVC is explicitly deleted; back it up or
set the PV reclaim policy to `Retain` before deleting the claim if necessary.
`terraform destroy` removes the AKS cluster and dedicated resource group, including
managed cluster resources: retained PVCs are not a backup. Remove remaining
operator-created resources after reviewing their ownership and charges.

- Pending PVC: check the CSI driver, selected storage class, Azure disk quota,
  and availability-zone constraints (`kubectl describe pvc` and pod events).
- Init failure: inspect `kubectl logs POD -c init`. Check genesis validity, volume
  ownership, and whether the image supports the worker architecture and commands.
- API timeout: check your actual egress CIDR against the Terraform API allowlist.
- RPC healthy but not syncing: check network version, genesis, peers and node logs.
  Large state replay may exceed Helm's 15-minute wait; inspect the pod before retrying.

## Local verification

The following checks do not provision cloud resources. Install the test dependency
PyYAML 6.x in your usual development environment; the CLI itself uses only Python's
standard library.

```sh
python3 -m unittest discover -s deploy/tests -v
terraform fmt -check -recursive deploy/terraform
terraform -chdir=deploy/terraform/providers/azure init -backend=false
terraform -chdir=deploy/terraform/providers/azure validate
terraform -chdir=deploy/terraform/providers/azure test
shellcheck deploy/helm/sifnode/scripts/init.sh deploy/tests/smoke-kind.sh
```

Set `SIFNODE_TEST_BINARY=/absolute/path/to/sifnoded` when running unittest to also
run the native bootstrap tests (otherwise those four tests are explicitly skipped).
They verify first boot, invalid-genesis recovery and preservation of genesis,
signing keys and consensus state on restart.

Terraform tests use a mocked Azure provider, checking configuration and input
validation without credentials. To exercise the node on a disposable kind cluster,
load a compatible image containing the real binary and run:

```sh
kind load docker-image YOUR_LOCAL_SIFNODE_IMAGE:TAG --name YOUR_KIND_CLUSTER
KUBECONFIG=/path/to/kind-config SIFNODE_SMOKE_IMAGE=YOUR_LOCAL_SIFNODE_IMAGE:TAG \
  sh deploy/tests/smoke-kind.sh
```

The smoke test refuses non-kind contexts, creates a fresh namespace and an
unfunded local test validator, runs the actual chart, checks block production and
key persistence after pod recreation, and checks that uninstall retains the PVC.
It deletes its namespace on exit. This verifies Kubernetes/node behavior locally;
Azure provisioning, regional capacity and load balancer allocation require an
operator-authorized AKS deployment.
