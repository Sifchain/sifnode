# Sifnode GCP Reference Architecture

**Document Version:** 1.0  
**Created:** 2026-04-21  
**Author:** Kael (Support)  
**Status:** Ready for Implementation

---

## Overview

This document defines the Google Cloud Platform (GCP) infrastructure architecture for deploying Sifnode validator nodes using Terraform.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         GCP Project                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                      VPC Network                           │  │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐   │  │
│  │  │   GKE       │◄───│  Cloud SQL  │    │   Cloud     │   │  │
│  │  │  Cluster    │    │ PostgreSQL  │    │  Load Balancer│  │  │
│  │  │             │    │             │    │             │   │  │
│  │  │ • Validator │    │ • Chain data│    │ • External  │   │  │
│  │  │   nodes     │    │ • State     │    │ • HTTPS     │   │  │
│  │  │ • Sentry    │    │ • Indices   │    │ • SSL term  │   │  │
│  │  │   nodes     │    │             │    │             │   │  │
│  │  └──────┬──────┘    └─────────────┘    └──────┬──────┘   │  │
│  │         │                                       │         │  │
│  │  ┌──────▼───────────────────────────────────────▼──────┐  │  │
│  │  │              Cloud Monitoring                         │  │  │
│  │  │  • Metrics  • Alerts  • Dashboards  • Logging       │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Module Specifications

### 1. GKE Module (modules/gke/)

**Purpose:** Kubernetes cluster for Sifnode validator and sentry nodes

**Configuration:**
- **Cluster Type:** Regional (multi-zone for HA)
- **Region:** us-central1 (configurable)
- **Node Pool:**
  - Machine type: n2-standard-4 (4 vCPU, 16 GB RAM)
  - Min nodes: 3
  - Max nodes: 5 (autoscaling enabled)
  - Disk: 100GB SSD per node
  - Preemptible: false (validators need stability)

**Node Types:**
| Role | Count | Resources | Purpose |
|------|-------|-----------|---------|
| Validator | 2 | n2-standard-4 | Core validation, signing |
| Sentry | 2 | n2-standard-2 | Public RPC, peer proxy |

**Security:**
- Private cluster (master authorized networks)
- Workload Identity enabled
- Network policies for pod isolation

---

### 2. Cloud SQL Module (modules/cloud-sql/)

**Purpose:** PostgreSQL database for chain data and application state

**Configuration:**
- **Database Version:** PostgreSQL 14
- **Tier:** db-custom-4-16384 (4 vCPU, 16GB RAM)
- **Storage:** 500GB SSD, auto-expand
- **High Availability:** Zone-redundant (us-central1-a, us-central1-b)
- **Backup:** Automated daily, 30-day retention

**Network:**
- Private IP only (VPC peering)
- No public IP allocation
- Cloud SQL Auth Proxy for K8s access

**Databases:**
| Database | Purpose |
|----------|---------|
| sifnode_chain | Chain state, blocks |
| sifnode_indices | Query indices |
| sifnode_app | Application data |

---

### 3. Load Balancer Module (modules/load-balancer/)

**Purpose:** External HTTPS load balancer for RPC endpoints

**Configuration:**
- **Type:** External HTTP(S) Load Balancer
- **SSL:** Google-managed certificate (Let's Encrypt)
- **Backend:** GKE service (sifnode-rpc)
- **Health Checks:** HTTP /health endpoint

**Ports:**
| Port | Service | Description |
|------|---------|-------------|
| 443 | RPC | JSON-RPC over HTTPS |
| 26657 | Tendermint RPC | Node status, queries |
| 1317 | REST API | LCD endpoints |

**Rate Limiting:**
- Cloud Armor: 1000 req/min per IP
- DDoS protection enabled

---

### 4. Monitoring Module (modules/monitoring/)

**Purpose:** Observability stack for infrastructure and application

**Components:**
- **Cloud Monitoring:** Metrics, dashboards, alerts
- **Cloud Logging:** Log aggregation, export
- **Cloud Trace:** Distributed tracing
- **Uptime Checks:** External health monitoring

**Key Metrics:**
| Metric | Alert Threshold | Severity |
|--------|-----------------|----------|
| Node CPU | > 80% for 5 min | Warning |
| Node Memory | > 90% for 5 min | Critical |
| Block height stuck | No change for 2 min | Critical |
| Peer count | < 5 peers | Warning |
| Disk usage | > 85% | Warning |

**Dashboards:**
- Sifnode Overview (blocks, peers, sync status)
- Infrastructure Health (CPU, memory, disk)
- RPC Performance (latency, error rates)

---

## Security Best Practices

### Network Security
- [ ] VPC with custom subnets (no default network)
- [ ] Private GKE cluster (master auth networks)
- [ ] Cloud SQL private IP only
- [ ] Firewall rules: least privilege
- [ ] Cloud Armor for DDoS protection

### IAM & Access
- [ ] Service accounts per module
- [ ] Workload Identity for GKE
- [ ] No service account keys (use IAM roles)
- [ ] Cloud IAM audit logging enabled

### Data Protection
- [ ] Cloud SQL encrypted at rest
- [ ] Cloud SQL automated backups
- [ ] GKE node disk encryption
- [ ] Secrets in Secret Manager

---

## Cost Estimation

| Resource | Monthly Cost |
|----------|--------------|
| GKE (4 nodes, n2-standard-4) | ~$350 |
| Cloud SQL (db-custom-4-16384, HA) | ~$120 |
| Load Balancer | ~$18 |
| Cloud Monitoring | ~$25 |
| **Total** | **~$513/month** |

*Estimate for production deployment. Dev environments can use smaller instances.*

---

## Implementation Checklist

### Pre-Deployment
- [ ] GCP project created with billing
- [ ] Required APIs enabled (GKE, Cloud SQL, LB, Monitoring)
- [ ] Service accounts created with IAM roles
- [ ] Terraform state backend configured (GCS bucket)

### Deployment Order
1. [ ] Network module (VPC, subnets, firewall)
2. [ ] Cloud SQL module (database provisioning)
3. [ ] GKE module (cluster, node pools)
4. [ ] Load Balancer module (IP, SSL, backend)
5. [ ] Monitoring module (alerts, dashboards)

### Post-Deployment
- [ ] Verify all services healthy
- [ ] Test RPC endpoints
- [ ] Validate monitoring dashboards
- [ ] Document any deviations

---

**Next:** See `../testing/validation-checklist.md` for test procedures.
