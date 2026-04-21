# Sifnode GCP Validation Checklist

**Document Version:** 1.0  
**Created:** 2026-04-21  
**Author:** Kael (Support)

---

## Pre-Deployment Validation

### Terraform Validation
| Check | Command | Expected Result |
|-------|---------|-----------------|
| Syntax check | `terraform validate` | ✓ Success |
| Format check | `terraform fmt -check` | ✓ All files formatted |
| Plan review | `terraform plan` | Review output, no errors |
| State lock | Check GCS bucket | No stale locks |

### Security Validation
| Check | Method | Requirement |
|-------|--------|-------------|
| No hardcoded secrets | `grep -r "password\|secret\|key" *.tf` | Use vars/Secret Manager only |
| Private cluster | Verify GKE config | `enable_private_nodes = true` |
| SQL private IP | Check Cloud SQL config | `ipv4_enabled = false` |
| Firewall rules | Review rules | Least privilege, specific CIDRs |

---

## Module-by-Module Testing

### 1. GKE Module Tests

| Test | Command/Method | Pass Criteria |
|------|----------------|---------------|
| Cluster reachable | `gcloud container clusters get-credentials` | No errors |
| Nodes ready | `kubectl get nodes` | All nodes Ready |
| Workload Identity | `kubectl auth can-i list pods` | Success |
| Network policies | Deploy test pod | Pod isolation works |
| Autoscaling | Check HPA | Scales to demand |

**Sifnode-Specific Tests:**
| Test | Method | Pass Criteria |
|------|--------|---------------|
| Validator pod running | `kubectl get pods -n sifnode` | Status: Running |
| Sentry pod running | `kubectl get pods -n sifnode` | Status: Running |
| ConfigMaps mounted | `kubectl describe pod` | Config volume mounted |
| Secrets accessible | Pod logs | No secret access errors |

---

### 2. Cloud SQL Module Tests

| Test | Command | Pass Criteria |
|------|---------|---------------|
| Instance created | `gcloud sql instances list` | Status: RUNNABLE |
| Private IP working | `gcloud compute networks peerings list` | Peering active |
| Database accessible | `psql -h [PRIVATE_IP] -U postgres` | Connection succeeds |
| Backup configured | `gcloud sql backups list` | Recent backups shown |
| HA configured | `gcloud sql instances describe` | `availabilityType: REGIONAL` |

**Sifnode Database Tests:**
| Test | Method | Pass Criteria |
|------|--------|---------------|
| Chain DB exists | `\l` in psql | sifnode_chain listed |
| Tables created | `\dt` in sifnode_chain | Expected tables present |
| Indices created | Query performance | Queries use indices |
| Connection from K8s | Test pod connection | Success via Cloud SQL Proxy |

---

### 3. Load Balancer Module Tests

| Test | Command | Pass Criteria |
|------|---------|---------------|
| IP reserved | `gcloud compute addresses list` | Status: RESERVED |
| SSL certificate | `gcloud compute ssl-certificates list` | Status: ACTIVE |
| Backend service | `gcloud compute backend-services list` | Health check passing |
| URL map | `gcloud compute url-maps list` | Routes configured |

**Endpoint Tests:**
| Test | URL | Expected Response |
|------|-----|-------------------|
| Health check | `https://[DOMAIN]/health` | HTTP 200 |
| RPC endpoint | `https://[DOMAIN]:443` | JSON-RPC response |
| Tendermint RPC | `https://[DOMAIN]:26657/status` | Node status JSON |
| REST API | `https://[DOMAIN]:1317/node_info` | Node info JSON |

---

### 4. Monitoring Module Tests

| Test | Method | Pass Criteria |
|------|--------|---------------|
| Dashboards created | Cloud Console → Monitoring → Dashboards | Sifnode dashboards present |
| Alert policies | Cloud Console → Alerting | Policies active |
| Uptime checks | Cloud Console → Uptime | Checks passing |
| Log sinks | `gcloud logging sinks list` | Sinks configured |

**Metric Verification:**
| Metric | Query | Expected |
|--------|-------|----------|
| CPU usage | `compute.googleapis.com/instance/cpu/utilization` | Data flowing |
| Disk I/O | `compute.googleapis.com/instance/disk/write_bytes_count` | Data flowing |
| Custom metrics | Sifnode exporter | Block height metric |

---

## Sifnode Application Tests

### Sync Status
| Test | Command | Pass Criteria |
|------|---------|---------------|
| Node catching up | `curl localhost:26657/status` | `catching_up: false` |
| Block height | Compare with public explorer | Within 5 blocks |
| Sync rate | Monitor over 10 min | > 100 blocks/min |

### API Tests
| Test | Command | Pass Criteria |
|------|---------|---------------|
| RPC available | `curl localhost:26657/net_info` | Returns peer info |
| REST API available | `curl localhost:1317/bank/balances/[ADDR]` | Returns balance |
| gRPC available | `grpcurl -plaintext localhost:9090 list` | Services listed |

### Health Checks
| Test | Method | Pass Criteria |
|------|--------|---------------|
| Process running | `pgrep sifnoded` | Process found |
| Port listening | `netstat -tlnp | grep 26657` | Port open |
| Memory usage | `kubectl top pod` | < 8GB per pod |
| Disk usage | `df -h` | < 85% used |

---

## Failure Scenario Tests

### Node Failure Simulation
| Scenario | Action | Expected Recovery |
|----------|--------|-------------------|
| Validator pod restart | `kubectl delete pod` | Auto-restart, re-sync |
| Node drain | `kubectl drain` | Pods migrate to other nodes |
| Zone outage | Simulate zone failure | HA keeps service up |

### Database Failure
| Scenario | Action | Expected Recovery |
|----------|--------|-------------------|
| Primary failover | Trigger failover | Automatic, < 60s downtime |
| Connection loss | Block SQL port | App reconnects automatically |
| Backup restore | Restore from backup | Data integrity maintained |

### Network Failure
| Scenario | Action | Expected Recovery |
|----------|--------|-------------------|
| LB backend unhealthy | Stop pods | LB routes to healthy |
| SSL cert expiry | Check cert rotation | Auto-renewal |
| DDoS attack | Cloud Armor test | Rate limiting triggers |

---

## Performance Benchmarks

### Load Testing
| Test | Tool | Target |
|------|------|--------|
| RPC throughput | `wrk` | 1000 req/s |
| Tendermint queries | Custom script | 500 req/s |
| Sync speed | Monitor | 100+ blocks/sec |

### Latency Targets
| Metric | Target | Critical |
|--------|--------|----------|
| RPC response | < 100ms | > 500ms |
| Block propagation | < 5s | > 30s |
| DB query | < 50ms | > 200ms |

---

## Sign-Off Checklist

| Validator | Role | Date | Signature |
|-----------|------|------|-----------|
| | William (Lead) | | |
| | Kael (Support) | | |

**Final Approval:**
- [ ] All module tests passing
- [ ] All application tests passing
- [ ] Failure scenarios validated
- [ ] Performance benchmarks met
- [ ] Documentation updated
- [ ] Ready for production

---

**Next:** See `../runbooks/deployment-runbook.md` for deployment procedures.
