# Sifnode GCP Testing & Validation Checklist

**Document:** Testing Framework for GCP Implementation  
**Role:** Kael Support - Validation & QA  
**Date:** 2026-04-21

---

## Pre-Deployment Validation

### Terraform Validation
- [ ] `terraform fmt` passes (no formatting issues)
- [ ] `terraform validate` passes (syntax valid)
- [ ] `terraform plan` shows expected resources
- [ ] No hardcoded credentials in `.tf` files
- [ ] All variables have descriptions
- [ ] All outputs defined and documented

### Security Validation
- [ ] GKE cluster uses private nodes
- [ ] Cloud SQL has no public IP
- [ ] Service accounts use Workload Identity
- [ ] Firewall rules restrict to required ports only
- [ ] Secrets stored in Kubernetes Secrets, not plain text
- [ ] SSL certificate configured for load balancer

### Cost Validation
- [ ] Resource limits set (prevent runaway costs)
- [ ] Billing alerts configured
- [ ] Preemptible nodes NOT used for validators (require stability)

---

## Deployment Testing

### GKE Module Tests
1. **Cluster Creation**
   ```bash
   gcloud container clusters list --filter="name=sifnode-cluster"
   ```
   - [ ] Cluster exists in correct region
   - [ ] Node pool has 3+ nodes
   - [ ] Nodes show "Ready" status

2. **Node Configuration**
   ```bash
   kubectl get nodes -o wide
   ```
   - [ ] Machine type matches spec (n2-standard-4)
   - [ ] Disk size >= 100GB
   - [ ] Shielded nodes enabled

3. **Networking**
   ```bash
   kubectl get svc
   ```
   - [ ] Services accessible within VPC
   - [ ] External IPs assigned correctly

### Cloud SQL Module Tests
1. **Instance Creation**
   ```bash
   gcloud sql instances describe sifnode-db
   ```
   - [ ] Instance status is RUNNABLE
   - [ ] Private IP configured
   - [ ] Public IP disabled
   - [ ] Automated backups enabled

2. **Database Connectivity**
   ```bash
   # From GKE pod
   pg_isready -h <CLOUD_SQL_PRIVATE_IP> -p 5432
   ```
   - [ ] Connection succeeds
   - [ ] Authentication works
   - [ ] Database created

### Load Balancer Tests
1. **Health Checks**
   ```bash
   curl https://<DOMAIN>/status
   ```
   - [ ] Returns HTTP 200
   - [ ] Response includes sync status

2. **SSL Certificate**
   ```bash
   openssl s_client -connect <DOMAIN>:443 -servername <DOMAIN>
   ```
   - [ ] Certificate valid
   - [ ] Chain complete
   - [ ] Not expired

3. **Backend Service**
   - [ ] All backends healthy
   - [ ] Traffic distributed across nodes

---

## Sifnode Application Tests

### Pod Status
```bash
kubectl get pods -n sifnode
```
- [ ] Pods in "Running" state
- [ ] No CrashLoopBackOff
- [ ] Restarts < 3

### Node Synchronization
```bash
kubectl exec -it sifnode-0 -- sifnoded status
```
- [ ] Latest_block_height increasing
- [ ] Catching_up: false
- [ ] Sync complete within 30 minutes of deployment

### API Endpoints
```bash
# RPC endpoint
curl http://<LB_IP>:26657/status

# REST API
curl http://<LB_IP>:1317/cosmos/base/tendermint/v1beta1/blocks/latest
```
- [ ] RPC responds with valid JSON
- [ ] REST API returns block data
- [ ] No timeout errors

---

## Monitoring Validation

### Metrics Collection
- [ ] Node CPU metrics visible in Cloud Monitoring
- [ ] Node memory metrics visible
- [ ] Disk usage metrics visible
- [ ] Sifnode-specific metrics (if custom metrics configured)

### Alert Tests
- [ ] Test alert: Simulate node down → verify notification received
- [ ] Test alert: High CPU usage → verify threshold triggers
- [ ] Test alert: Disk full → verify early warning

### Dashboards
- [ ] Dashboard exists: "Sifnode Validator Overview"
- [ ] All charts populate with data
- [ ] Time range selector works

---

## Failure Scenario Tests

### Node Failure
- [ ] Kill one pod: `kubectl delete pod sifnode-0`
- [ ] Verify pod auto-recreates
- [ ] Verify chain sync resumes
- [ ] Verify no data loss

### Database Connection Loss
- [ ] Simulate Cloud SQL restart
- [ ] Verify Sifnode reconnects automatically
- [ ] Verify no corruption

### Network Partition
- [ ] Block egress from one node
- [ ] Verify load balancer routes traffic away
- [ ] Verify recovery when partition heals

---

## Performance Benchmarks

| Metric | Target | Actual | Pass |
|--------|--------|--------|------|
| Block sync speed | >100 blocks/sec | | |
| API response time | <200ms p95 | | |
| Node startup time | <5 minutes | | |
| Failover time | <30 seconds | | |

---

## Sign-Off

**Tested by:** Kael  
**Date:** ___________  
**Result:** ☐ PASS / ☐ FAIL

**Notes:**
_____________________________________________
_____________________________________________
