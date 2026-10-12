# Sifnode on Google Kubernetes Engine (GKE)

Independent implementation for bounty #524. The bundle provides Terraform infrastructure, Kubernetes workload definitions, a guarded deployment CLI, and pull-request validation.

## 1. Configure GKE

```bash
cp deploy/gke/terraform/terraform.tfvars.example deploy/gke/terraform/terraform.tfvars
# edit project_id and any sizing/network values
```

Set the cluster connection variables used by the lifecycle CLI:

```bash
export SIFNODE_GKE_PROJECT=<gcp-project-id>
export SIFNODE_GKE_REGION=<gcp-region>
export SIFNODE_GKE_CLUSTER=<cluster-name>
```

Equivalent `GOOGLE_CLOUD_PROJECT` and `GOOGLE_REGION` variables are also accepted.

## 2. Preflight and plan

```bash
bash deploy/gke/deploy_gke.sh validate
bash deploy/gke/deploy_gke.sh plan
```

`validate` checks the local toolchain, active Google Cloud authentication, Terraform formatting/configuration, and Kubernetes manifest syntax without changing cloud resources. `plan` shows the Terraform changes that would be applied.

## 3. Provision GKE

```bash
bash deploy/gke/deploy_gke.sh provision
```

`provision` runs Terraform with the configured project/region and then establishes `kubectl` credentials for the resulting cluster. It is intentionally separate from `apply` so infrastructure creation is explicit.

## 4. Configure the node

```bash
cp deploy/gke/kubernetes/sifnode.env.example deploy/gke/kubernetes/sifnode.env
# set CHAINNET, MONIKER, EXTERNAL_IP, GENESIS_URL and MNEMONIC
```

`MNEMONIC` is consumed only from the local environment file and is converted into a Kubernetes Secret by the deployment script. Do not commit the local environment file. Placeholder values are rejected before deployment.

## 5. Deploy and verify

```bash
bash deploy/gke/deploy_gke.sh apply
bash deploy/gke/deploy_gke.sh status
```

`apply` establishes cluster credentials, creates/updates the namespace and configuration, applies the workload, waits for rollout, and verifies the Deployment is available. The lifecycle CLI supports `validate`, `plan`, `provision`, `apply`, `status`, and guarded `destroy`/`destroy-infra`. Destructive operations require `CONFIRM_DESTROY=yes`.

`status` is intentionally read-only and only requires `kubectl`; it does not require the local secret file or cloud credentials.

## Architecture

Terraform creates a regional VPC-native GKE cluster, dedicated subnet with pod/service secondary ranges, Cloud NAT for private node egress, Workload Identity, and the GKE persistent-disk CSI integration. Kubernetes then runs Sifnode with a persistent volume and exposes P2P, RPC, and API ports through a LoadBalancer service.

A GitHub Actions workflow runs shell syntax validation and Terraform formatting/validation for every change under `deploy/gke`.

## Validation boundary

Local/CI validation does not claim a live GKE deployment. A real `provision` + `apply` run requires an authorized Google Cloud project, quota, credentials, and a valid Sifnode network configuration. No payment or acceptance is implied until the maintainer verifies the acceptance criteria.

## Bounty/payment contact

For bounty submission/payment administration, the contributor may provide: https://www.paypal.com/paypalme/my/profile
