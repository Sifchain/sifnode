# Deploying sifnode on Google Kubernetes Engine (GKE)

This guide walks through deploying a sifnode validator on **Google Kubernetes Engine (GKE)**.

## Prerequisites

| Requirement | Version | Install |
|-------------|---------|---------|
| gcloud CLI | latest | [Install](https://cloud.google.com/sdk/docs/install) |
| Terraform | >= 0.13 | [Install](https://learn.hashicorp.com/tutorials/terraform/install-cli) |
| kubectl | >= 1.18 | `gcloud components install kubectl` |
| Helm 3 | >= 3.0 | [Install](https://helm.sh/docs/intro/install/) |

## Quick Start

### 1. Authenticate with GCP

```bash
gcloud auth login
gcloud config set project YOUR_PROJECT_ID
```

### 2. Deploy Infrastructure

```bash
# Clone the repository
git clone https://github.com/Sifchain/sifnode.git
cd sifnode

# Deploy with the automation script
./deploy/gke/deploy_gke.sh create -p YOUR_PROJECT_ID -r us-central1
```

### 3. Verify Deployment

```bash
# Check cluster status
./deploy/gke/deploy_gke.sh status

# Check pod logs
kubectl logs -n sifnode -l app=sifnode

# Check sifnoded status
kubectl exec -n sifnode deploy/sifnode -- sifnoded status
```

## Manual Deployment (Step by Step)

### Cluster Setup

```bash
cd deploy/terraform/providers/gke

terraform init
terraform apply \
  -var="project_id=YOUR_PROJECT_ID" \
  -var="region=us-central1" \
  -var="cluster_name=sifnode-gke"
```

### Configure kubectl

```bash
gcloud container clusters get-credentials sifnode-gke \
  --region us-central1 \
  --project YOUR_PROJECT_ID
```

### Install sifnode via Helm

```bash
# Create namespace
kubectl create namespace sifnode

# Install with GKE-specific values
helm install sifnode deploy/helm/sifnode \
  --namespace sifnode \
  --create-namespace \
  --values deploy/helm/sifnode/values_gke.yaml \
  --set provider=gke \
  --set persistence.storageClass=gke-ssd

# Check deployment
kubectl get pods -n sifnode -w
```

## Configuration Reference

### GKE Terraform Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `project_id` | (required) | GCP project ID |
| `region` | `us-central1` | GCP region |
| `cluster_name` | (required) | GKE cluster name |
| `machine_type` | `e2-standard-2` | Node machine type |
| `desired_capacity` | `1` | Initial node count |
| `min_capacity` | `1` | Minimum nodes (autoscaling) |
| `max_capacity` | `5` | Maximum nodes (autoscaling) |
| `disk_size` | `100` | Node disk size (GB) |
| `storage_size` | `50Gi` | PVC storage size |
| `enable_private_nodes` | `false` | Enable private GKE nodes |
| `release_channel` | `REGULAR` | GKE release channel |

### Helm Values (GKE-specific)

Key settings in `values_gke.yaml`:

```yaml
provider: gke
persistence:
  storageClass: gke-ssd  # SSD persistent storage
resources:
  requests:
    cpu: "500m"
    memory: "1Gi"
  limits:
    cpu: "2"
    memory: "4Gi"
autoscaling:
  enabled: true
```

## Comparison: AWS (EKS) vs GKE

| Feature | AWS (EKS) | GKE |
|---------|-----------|-----|
| Provisioning | Terraform (terraform-aws-modules/eks) | Terraform (google_container_cluster) |
| Node Type | t2.medium (x86) | e2-standard-2 (x86) |
| Storage | EBS CSI / EFS CSI | GCE Persistent Disk (pd-ssd) |
| Networking | VPC + Public Subnets | VPC + Subnet + Cloud NAT |
| DNS | Route53 (via ExternalDNS) | Cloud DNS |
| Ingress | ALB Ingress Controller | GCE Ingress |
| IAM | AWS IAM Roles | GCP IAM + Workload Identity |
| CLI Profile | AWS_PROFILE | gcloud config |

## Troubleshooting

### Pods stuck in Pending state
```bash
kubectl describe pod -n sifnode <pod-name>
# Check for resource constraints or PVC binding issues
```

### GKE cluster not reachable
```bash
gcloud container clusters get-credentials sifnode-gke --region us-central1
# Ensure your gcloud is authenticated
```

### Storage issues
```bash
kubectl get pvc -n sifnode
kubectl get storageclass
# Ensure gke-ssd storage class exists
```

## Clean Up

```bash
# Option 1: Using the deployment script
./deploy/gke/deploy_gke.sh destroy

# Option 2: Manual
helm uninstall sifnode --namespace sifnode
cd deploy/terraform/providers/gke
terraform destroy -var="project_id=YOUR_PROJECT_ID"
```

---

*For additional assistance, refer to the [sifnode documentation](https://docs.sifchain.finance/).*
