# Sifnode GCP Deployment Runbook

**Document:** Step-by-step deployment procedures  
**Role:** Kael Support - Operations Documentation  
**Date:** 2026-04-21  
**Version:** 1.0 (Draft)

---

## Prerequisites

- [ ] GCP project created with billing enabled
- [ ] gcloud CLI authenticated: `gcloud auth login`
- [ ] Terraform >= 1.5 installed
- [ ] kubectl installed
- [ ] Docker installed (for image building)
- [ ] Access to Sifnode source code

---

## Phase 1: Environment Setup (15 minutes)

### 1.1 Configure GCP Project
```bash
export PROJECT_ID="sifnode-gcp-support"
export REGION="us-central1"
export ZONE="us-central1-a"

gcloud config set project $PROJECT_ID
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE
```

### 1.2 Enable Required APIs
```bash
gcloud services enable container.googleapis.com
gcloud services enable sqladmin.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable monitoring.googleapis.com
gcloud services enable logging.googleapis.com
```

### 1.3 Create Service Account
```bash
gcloud iam service-accounts create terraform-sa \
  --display-name="Terraform Service Account"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:terraform-sa@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/editor"
```

---

## Phase 2: Terraform Deployment (30 minutes)

### 2.1 Initialize Terraform
```bash
cd /home/wls/bounty-workspaces/sifnode-gcp/terraform/environments/dev

terraform init
terraform workspace new dev || terraform workspace select dev
```

### 2.2 Configure Variables
Create `terraform.tfvars`:
```hcl
project_id = "sifnode-gcp-support"
region     = "us-central1"

# GKE
gke_cluster_name = "sifnode-dev"
gke_node_count   = 3

# Cloud SQL
sql_instance_name = "sifnode-db-dev"
sql_database_name = "sifnodedb"

# Notification
notification_email = "ops@example.com"
```

### 2.3 Plan and Apply
```bash
terraform plan -out=tfplan
terraform apply tfplan
```

### 2.4 Configure kubectl
```bash
gcloud container clusters get-credentials sifnode-dev \
  --region=us-central1
```

---

## Phase 3: Database Setup (10 minutes)

### 3.1 Create Database User
```bash
gcloud sql users create sifnode \
  --instance=sifnode-db-dev \
  --password="$(openssl rand -base64 32)"
```

### 3.2 Create Database
```bash
gcloud sql databases create sifnodedb \
  --instance=sifnode-db-dev
```

### 3.3 Store Secret in Kubernetes
```bash
kubectl create secret generic db-credentials \
  --from-literal=username=sifnode \
  --from-literal=password="YOUR_PASSWORD"
```

---

## Phase 4: Sifnode Deployment (20 minutes)

### 4.1 Build Docker Image
```bash
cd /home/wls/bounty-workspaces/sifnode-gcp

docker build -t gcr.io/$PROJECT_ID/sifnode:latest .
docker push gcr.io/$PROJECT_ID/sifnode:latest
```

### 4.2 Deploy to GKE
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/statefulset.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

### 4.3 Verify Deployment
```bash
kubectl get pods -n sifnode -w
kubectl logs -f sifnode-0 -n sifnode
```

---

## Phase 5: Verification (10 minutes)

### 5.1 Check Node Status
```bash
kubectl exec -it sifnode-0 -n sifnode -- sifnoded status
```

### 5.2 Test RPC Endpoint
```bash
export LB_IP=$(kubectl get svc sifnode-lb -n sifnode -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl http://$LB_IP:26657/status | jq
```

### 5.3 Verify Metrics
```bash
gcloud monitoring dashboards list
```

---

## Phase 6: Production Deployment

Repeat Phases 2-5 with:
- Workspace: `prod`
- Region: `us-central1` (or preferred region)
- Node count: 5 (for HA)
- Cloud SQL HA: Enabled

---

## Rollback Procedures

### Scenario: Deployment Failure
```bash
# Destroy resources
terraform destroy

# Verify cleanup
gcloud compute instances list
gcloud container clusters list
gcloud sql instances list
```

### Scenario: Application Issue
```bash
# Rollback to previous version
kubectl rollout undo statefulset/sifnode -n sifnode

# Check rollout status
kubectl rollout status statefulset/sifnode -n sifnode
```

---

## Post-Deployment Tasks

- [ ] Document actual resource names
- [ ] Save terraform.tfvars to secure location
- [ ] Configure backup verification
- [ ] Test alert notifications
- [ ] Update monitoring dashboards
- [ ] Schedule security review

---

## Support Contacts

- **Technical Lead:** William
- **Support/Documentation:** Kael
- **Infrastructure:** Sifnode GCP Team

---

**Last Updated:** 2026-04-21  
**Next Review:** Before production deployment
