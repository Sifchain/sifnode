# Sifnode on Google Kubernetes Engine (GKE)

Independent implementation for bounty #524. The bundle provides Terraform infrastructure, Kubernetes workload definitions, a guarded deployment CLI, and pull-request validation.

## 1. Provision GKE

```bash
cp deploy/gke/terraform/terraform.tfvars.example deploy/gke/terraform/terraform.tfvars
# edit project_id
terraform -chdir=deploy/gke/terraform init
terraform -chdir=deploy/gke/terraform plan
terraform -chdir=deploy/gke/terraform apply
```

The lifecycle script can now establish cluster credentials itself. Set the cluster connection variables before `apply`:

```bash
export SIFNODE_GKE_PROJECT=<gcp-project-id>
export SIFNODE_GKE_REGION=<gcp-region>
export SIFNODE_GKE_CLUSTER=<cluster-name>
```

Equivalent `GOOGLE_CLOUD_PROJECT` and `GOOGLE_REGION` variables are also accepted. This removes the need to manually copy the Terraform `get_credentials_command` into a separate shell step.

## 2. Configure the node

```bash
cp deploy/gke/kubernetes/sifnode.env.example deploy/gke/kubernetes/sifnode.env
# set CHAINNET, MONIKER, EXTERNAL_IP, GENESIS_URL and MNEMONIC
```

`MNEMONIC` is consumed only from the local environment file and is converted into a Kubernetes Secret by the deployment script. Do not commit the local environment file. Placeholder values are rejected before deployment.

## 3. Deploy

```bash
bash deploy/gke/deploy_gke.sh validate
bash deploy/gke/deploy_gke.sh apply
bash deploy/gke/deploy_gke.sh status
```

`apply` automatically runs `gcloud container clusters get-credentials`, creates/updates the namespace and configuration, applies the workload, waits for rollout, and then verifies the Deployment is available. The lifecycle CLI supports `validate`, `apply`, `status`, and guarded `destroy`. Destruction requires `CONFIRM_DESTROY=yes`.

`status` is intentionally read-only and only requires `kubectl`; it does not require the local secret file or cloud credentials.

## Architecture

Terraform creates a regional VPC-native GKE cluster, dedicated subnet with pod/service secondary ranges, Cloud NAT for private node egress, Workload Identity, and the GKE persistent-disk CSI integration. Kubernetes then runs Sifnode with a persistent volume and exposes P2P, RPC, and API ports through a LoadBalancer service.

A GitHub Actions workflow runs shell syntax validation and Terraform formatting/validation for every change under `deploy/gke`.

## Bounty/payment contact

For bounty submission/payment administration, the contributor may provide: https://www.paypal.com/paypalme/my/profile
