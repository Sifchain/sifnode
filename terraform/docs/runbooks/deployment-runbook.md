# Sifnode GCP Deployment Runbook

**Document Version:** 1.0  
**Created:** 2026-04-21  
**Author:** Kael (Support)  
**Status:** Ready for Execution

---

## Overview

This runbook provides step-by-step procedures for deploying Sifnode validator infrastructure on Google Cloud Platform using Terraform.

**Estimated Time:** 4-6 hours  
**Prerequisites:** GCP project with billing, Terraform >= 1.0

---

## Phase 1: Environment Setup (30 min)

### 1.1 GCP Project Configuration

```bash
# Set project
export PROJECT_ID="sifnode-gcp-mainnet"
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable container.googleapis.com
gcloud services enable sqladmin.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable monitoring.googleapis.com
gcloud services enable logging.googleapis.com
```

### 1.2 Service Account Setup

```bash
# Create Terraform service account
gcloud iam service-accounts create terraform-deploy \
  --display-name="Terraform Deployer"

# Grant IAM roles
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:terraform-deploy@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/editor"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:terraform-deploy@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/container.admin"

# Create and download key
gcloud iam service-accounts keys create terraform-key.json \
  --iam-account=terraform-deploy@$PROJECT_ID.iam.gserviceaccount.com
```

### 1.3 Terraform Backend

```bash
# Create GCS bucket for state
gsutil mb -p $PROJECT_ID gs://$PROJECT_ID-tfstate

gsutil versioning set on gs://$PROJECT_ID-tfstate
```

---

## Phase 2: Terraform Deployment (2-3 hours)

### 2.1 Initialize Terraform

```bash
cd /home/wls/bounty-workspaces/sifnode-gcp/terraform

# Initialize with backend
cat > backend.tf << 'EOF'
terraform {
  backend "gcs" {
    bucket = "sifnode-gcp-mainnet-tfstate"
    prefix = "terraform/state"
  }
}
EOF

terraform init
```

### 2.2 Configure Variables

```bash
# Create terraform.tfvars
cat > terraform.tfvars << 'EOF'
project_id = "sifnode-gcp-mainnet"
region     = "us-central1"

# GKE Configuration
cluster_name     = "sifnode-mainnet"
node_pool_name   = "validator-pool"
node_count       = 3
node_machine_type = "n2-standard-4"

# Cloud SQL Configuration
sql_instance_name = "sifnode-db"
database_tier     = "db-custom-4-16384"
database_version  = "POSTGRES_14"

# Network
vpc_name    = "sifnode-vpc"
subnet_name = "sifnode-subnet"
subnet_cidr = "10.0.0.0/24"

# Monitoring
enable_monitoring = true
alert_email       = "alerts@example.com"
EOF
```

### 2.3 Deploy Network Module

```bash
terraform plan -target=module.network
terraform apply -target=module.network -auto-approve

# Verify
gcloud compute networks list
gcloud compute networks subnets list
```

### 2.4 Deploy Cloud SQL Module

```bash
terraform plan -target=module.cloud-sql
terraform apply -target=module.cloud-sql -auto-approve

# Verify
gcloud sql instances list
```

### 2.5 Deploy GKE Module

```bash
terraform plan -target=module.gke
terraform apply -target=module.gke -auto-approve

# Verify
gcloud container clusters list
kubectl get nodes
```

### 2.6 Deploy Load Balancer Module

```bash
terraform plan -target=module.load-balancer
terraform apply -target=module.load-balancer -auto-approve

# Verify
gcloud compute forwarding-rules list
```

### 2.7 Deploy Monitoring Module

```bash
terraform plan -target=module.monitoring
terraform apply -target=module.monitoring -auto-approve

# Verify
gcloud monitoring dashboards list
```

---

## Phase 3: Database Setup (30 min)

### 3.1 Create Databases

```bash
# Get private IP
SQL_IP=$(gcloud sql instances describe sifnode-db --format='value(ipAddresses[0].ipAddress)')

# Connect and create databases
gcloud sql connect sifnode-db --user=postgres --quiet << 'EOF'
CREATE DATABASE sifnode_chain;
CREATE DATABASE sifnode_indices;
CREATE DATABASE sifnode_app;
CREATE USER sifnode WITH PASSWORD '[SECURE_PASSWORD]';
GRANT ALL PRIVILEGES ON DATABASE sifnode_chain TO sifnode;
GRANT ALL PRIVILEGES ON DATABASE sifnode_indices TO sifnode;
\q
EOF
```

### 3.2 Configure Connection

```bash
# Create Cloud SQL Proxy secret
kubectl create secret generic cloudsql-credentials \
  --from-literal=username=sifnode \
  --from-literal=password='[SECURE_PASSWORD]'
```

---

## Phase 4: Sifnode Deployment (1-2 hours)

### 4.1 Build Docker Image

```bash
cd /home/wls/bounty-workspaces/sifnode-gcp

# Build
docker build -t gcr.io/$PROJECT_ID/sifnode:latest -f Dockerfile .

# Push
docker push gcr.io/$PROJECT_ID/sifnode:latest
```

### 4.2 Deploy to GKE

```bash
# Create namespace
kubectl create namespace sifnode

# Deploy validator nodes
kubectl apply -f k8s/validator-deployment.yaml
kubectl apply -f k8s/validator-service.yaml

# Deploy sentry nodes
kubectl apply -f k8s/sentry-deployment.yaml
kubectl apply -f k8s/sentry-service.yaml

# Verify
kubectl get pods -n sifnode
kubectl get services -n sifnode
```

### 4.3 Verify Sifnode Status

```bash
# Check pod logs
kubectl logs -f deployment/sifnode-validator -n sifnode

# Check sync status
kubectl exec -it deployment/sifnode-validator -n sifnode -- \
  sifnoded status | jq '.SyncInfo.catching_up'

# Expected: false (when synced)
```

---

## Phase 5: Verification (30 min)

### 5.1 Infrastructure Tests

```bash
# Run validation checklist tests
./scripts/validate-infrastructure.sh

# Expected: All tests pass
```

### 5.2 Endpoint Tests

```bash
# Get load balancer IP
LB_IP=$(gcloud compute forwarding-rules describe sifnode-lb \
  --global --format='value(IPAddress)')

# Test RPC
curl -s http://$LB_IP:26657/status | jq '.result.sync_info'

# Test REST API
curl -s http://$LB_IP:1317/node_info | jq '.node_info'

# Expected: Valid JSON responses
```

### 5.3 Monitoring Verification

```bash
# Check dashboards
gcloud monitoring dashboards list

# Check alerts
gcloud alpha monitoring policies list

# Verify uptime checks
gcloud monitoring uptime list
```

---

## Rollback Procedures

### Scenario 1: Terraform Destroy (Full Teardown)

```bash
# Destroy all resources
terraform destroy -auto-approve

# Verify
gcloud compute instances list  # Should be empty
gcloud container clusters list  # Should be empty
gcloud sql instances list       # Should be empty
```

### Scenario 2: Partial Rollback (Single Module)

```bash
# Example: Rollback GKE only
terraform destroy -target=module.gke -auto-approve

# Re-deploy
terraform apply -target=module.gke -auto-approve
```

### Scenario 3: Application Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/sifnode-validator -n sifnode

# Verify rollback
kubectl rollout status deployment/sifnode-validator -n sifnode
```

---

## Post-Deployment Tasks

### Immediate (Today)
- [ ] Verify all endpoints responding
- [ ] Confirm monitoring dashboards populated
- [ ] Test alert notifications
- [ ] Document any deviations from plan

### Short-term (This Week)
- [ ] Configure backup verification
- [ ] Document operational procedures
- [ ] Train team on incident response
- [ ] Review and optimize costs

### Ongoing
- [ ] Monitor sync status daily
- [ ] Review alerts weekly
- [ ] Update documentation monthly
- [ ] Security patches as needed

---

## Troubleshooting

### Common Issues

**Issue:** Terraform apply fails with API not enabled  
**Fix:** Run API enablement commands in Phase 1.1

**Issue:** GKE nodes not joining cluster  
**Fix:** Check firewall rules allow master-to-node communication

**Issue:** Cloud SQL connection refused  
**Fix:** Verify private IP peering and service account permissions

**Issue:** Sifnode not syncing  **Fix:** Check peer connections (`sifnoded net_info`), verify genesis file

---

## Contact Information

| Role | Responsibility | Contact |
|------|---------------|---------|
| William | Lead Implementation | bridge handoff |
| Kael | Support/Documentation | bridge handoff |

**Emergency:** Escalate to Lisa via bridge if deployment blocked > 30 min

---

**Reference:** See `../architecture/gcp-reference-architecture.md` for architecture details.
