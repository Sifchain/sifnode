# Sifnode GCP Reference Architecture

**Document:** GCP Architecture for Sifnode Validator Nodes  
**Role:** Kael Support Documentation  
**Date:** 2026-04-21  
**Status:** Draft (pre-implementation)

---

## Overview

This document defines the Google Cloud Platform architecture for deploying Sifnode validator nodes with high availability, security, and observability.

## Architecture Components

```
┌─────────────────────────────────────────────────────────────────┐
│                        GCP Project                               │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Cloud Load Balancer (Layer 7)               │   │
│  │              - SSL termination                           │   │
│  │              - Health checks                             │   │
│  └──────────────────────┬──────────────────────────────────┘   │
│                         │                                       │
│  ┌──────────────────────▼──────────────────────────────────┐   │
│  │              GKE Cluster (Regional)                    │   │
│  │  ┌─────────────────────────────────────────────────┐   │   │
│  │  │         Node Pool: validator-pool               │   │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐        │   │   │
│  │  │  │ Node 1  │  │ Node 2  │  │ Node 3  │        │   │   │
│  │  │  │Sifnode  │  │Sifnode  │  │Sifnode  │        │   │   │
│  │  │  │Pod      │  │Pod      │  │Pod      │        │   │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘        │   │   │
│  │  └─────────────────────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
│                         │                                       │
│  ┌──────────────────────▼──────────────────────────────────┐   │
│  │              Cloud SQL (PostgreSQL)                      │   │
│  │              - Chain data persistence                    │   │
│  │              - Automated backups                         │   │
│  │              - Private IP only                           │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Cloud Monitoring & Logging                │   │
│  │              - Metrics collection                        │   │
│  │              - Alerting policies                         │   │
│  │              - Log aggregation                           │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## Module Specifications

### 1. GKE Module (modules/gke)

**Purpose:** Container orchestration for Sifnode validator nodes

**Configuration:**
- **Cluster Type:** Regional (multi-zone for HA)
- **Node Pool:** validator-pool
  - Machine: n2-standard-4 (4 vCPU, 16GB RAM)
  - Disk: 100GB SSD persistent
  - Preemptible: false (validators need stability)
  - Autoscaling: 3-5 nodes
  - Taints: dedicated=validator:NoSchedule
- **Networking:** Private cluster, VPC-native
- **Security:** Workload Identity, Shielded GKE nodes

**Terraform Variables:**
```hcl
project_id       = string
region           = string
cluster_name     = string
node_count_min   = number (default: 3)
node_count_max   = number (default: 5)
machine_type     = string (default: "n2-standard-4")
```

---

### 2. Cloud SQL Module (modules/cloud-sql)

**Purpose:** Chain data persistence for Sifnode

**Configuration:**
- **Database:** PostgreSQL 14
- **Tier:** db-custom-2-4096 (2 vCPU, 4GB RAM)
- **Storage:** 100GB SSD, auto-expand
- **HA:** Regional availability
- **Backup:** Daily automated backups, 7-day retention
- **Access:** Private IP only (no public IP)

**Terraform Variables:**
```hcl
project_id       = string
region           = string
instance_name    = string
database_version = string (default: "POSTGRES_14")
tier             = string (default: "db-custom-2-4096")
```

---

### 3. Load Balancer Module (modules/load-balancer)

**Purpose:** RPC endpoint distribution and SSL termination

**Configuration:**
- **Type:** External HTTPS Load Balancer
- **Backend:** GKE service backend
- **SSL:** Managed SSL certificate
- **Health Checks:** /status endpoint
- **CDN:** Disabled (real-time blockchain data)

**Terraform Variables:**
```hcl
project_id       = string
name             = string
domain           = string  # Optional, for SSL
backend_service  = string
```

---

### 4. Monitoring Module (modules/monitoring)

**Purpose:** Observability for validator infrastructure

**Configuration:**
- **Metrics:** Node CPU, memory, disk, Sifnode sync status
- **Alerts:**
  - Node down > 5 minutes
  - Disk usage > 80%
  - Sync lag > 10 blocks
- **Dashboards:** Validator health overview

**Terraform Variables:**
```hcl
project_id       = string
notification_email = string
alert_channels   = list(string)
```

---

## Security Best Practices

1. **Network Security**
   - Private GKE cluster (no public control plane)
   - VPC peering for Cloud SQL access
   - Firewall rules: allow only necessary ports (26656, 26657, 1317)

2. **IAM**
   - Workload Identity for pod-to-GCP service
   - Least privilege service accounts
   - No default service account usage

3. **Secrets Management**
   - Kubernetes Secrets for node keys
   - Cloud KMS for encryption at rest
   - No hardcoded credentials

4. **Data Protection**
   - Cloud SQL encrypted with customer-managed keys
   - Automated backups with point-in-time recovery
   - VPC Service Controls for data exfiltration prevention

---

## Cost Estimation (Monthly)

| Resource | Configuration | Cost |
|----------|--------------|------|
| GKE Cluster | 3x n2-standard-4 | ~$290 |
| Cloud SQL | db-custom-2-4096 | ~$145 |
| Load Balancer | External HTTPS | ~$18 |
| Storage | 100GB SSD x 3 | ~$40 |
| Monitoring | Cloud Monitoring | ~$20 |
| **Total** | | **~$513/month** |

---

## Implementation Checklist

- [ ] VPC network created with custom subnets
- [ ] GKE cluster provisioned (regional)
- [ ] Cloud SQL instance with private IP
- [ ] Load balancer configured with health checks
- [ ] Monitoring dashboards and alerts created
- [ ] Sifnode Docker image built and pushed to GCR
- [ ] Kubernetes manifests created (StatefulSet, Service, ConfigMap)
- [ ] End-to-end deployment tested
- [ ] Documentation updated
- [ ] PR submitted to Sifnode repo

---

**Prepared by:** Kael (Support Role)  
**Review:** William (Lead Implementation)
