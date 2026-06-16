# Running a Sifchain Validator on AKS

This repository now includes Azure AKS deployment scaffolding that matches the older AWS/EKS workflow.

## Prerequisites

- Azure CLI installed
- Terraform installed
- Helm 3 installed
- `kubectl` installed
- Ruby installed
- Go installed

## 1. Sign in to Azure

```bash
az login
az account set --subscription "<your-subscription-id>"
```

The Terraform under `deploy/terraform/providers/azure` uses the standard Azure CLI authentication flow, so `az login` is enough for local testing.

## 2. Scaffold the AKS environment

From the repo root:

```bash
rake cluster:scaffold[testnet,azure]
```

This creates `.live/sifchain-azure-testnet/main.tf`. Before deploying, review it and adjust any of these defaults if needed:

- `location`
- `resource_group_name`
- `cluster_name`
- node sizing or Kubernetes version via the module inputs

The scaffold currently defaults to `southeastasia`, which is a safer starting point for restricted student subscriptions than `eastus`. If your subscription is limited to a different region set, change `location` before applying.

## 3. Create the AKS cluster

```bash
rake cluster:deploy[testnet,azure]
```

After Terraform finishes, the task writes the kubeconfig file to:

```text
.live/sifchain-azure-testnet/kubeconfig_sifchain-azure-testnet
```

## 4. Deploy sifnode

Standalone validator-style deployment:

```bash
rake cluster:sifnode:deploy:standalone[testnet,azure,sifnode,sifchain/sifnoded,testnet,my-validator,"your mnemonic here"]
```

If you want the node to join an existing network with an explicit peer and genesis:

```bash
rake cluster:sifnode:deploy:peer[testnet,azure,sifnode,sifchain/sifnoded,testnet,my-validator,"your mnemonic here",<peer-address>,<genesis-url>]
```

## 5. Verify the deployment

Point `kubectl` at the generated kubeconfig:

```bash
kubectl --kubeconfig .live/sifchain-azure-testnet/kubeconfig_sifchain-azure-testnet get pods -n sifnode
kubectl --kubeconfig .live/sifchain-azure-testnet/kubeconfig_sifchain-azure-testnet get svc -n sifnode
```

The chart waits for the AKS load balancer address and uses that external address when generating the node config. On Azure, the persistent volume claim defaults to the `managed-csi` storage class.

## 6. Clean up

Delete the namespace only:

```bash
rake cluster:sifnode:destroy[testnet,azure,sifnode,skip]
```

Destroy the entire AKS cluster:

```bash
rake cluster:destroy[testnet,azure]
```
