# Sifnode on Google Kubernetes Engine (GKE)

This directory provides an independent GKE deployment path for bounty #524. It intentionally keeps the deployment provider-neutral: Terraform provisions GKE, while Kubernetes manifests deploy the node.

## Quick start

1. Authenticate with Google Cloud and select a project.
2. Create the cluster with `terraform -chdir=deploy/gke/terraform init && terraform -chdir=deploy/gke/terraform apply`.
3. Fetch credentials with `gcloud container clusters get-credentials` using the Terraform outputs.
4. Copy `deploy/gke/kubernetes/sifnode.env.example` to a local environment file and set the chain/network values.
5. Run `deploy/gke/deploy_gke.sh`.

The script supports `apply`, `status`, and `destroy`, and refuses destructive operations unless explicitly requested.

## Configuration

Required environment variables are documented in `kubernetes/sifnode.env.example`. Secrets are supplied through the environment and are never committed to the repository.

The default storage class is `standard-rwo`, which maps to GKE's CSI-backed persistent disk implementation. The node exposes Tendermint RPC/P2P ports and the Cosmos REST port through a LoadBalancer service.

## Validation

Before applying infrastructure, run:

```bash
deploy/gke/deploy_gke.sh validate
```

This checks required local tools, required environment variables, and Kubernetes manifest rendering prerequisites without changing cloud resources.

## Payment contact

For bounty submission/payment administration, the contributor may provide: https://www.paypal.com/paypalme/my/profile
