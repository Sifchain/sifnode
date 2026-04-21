# -----------------------------------------------------------------------------
# GKE Module - Outputs
# Sifnode Validator Cluster
# -----------------------------------------------------------------------------

output "cluster_name" {
  description = "The name of the GKE cluster."
  value       = google_container_cluster.validator.name
}

output "cluster_id" {
  description = "The unique identifier for the GKE cluster (project/location/name)."
  value       = google_container_cluster.validator.id
}

output "endpoint" {
  description = "The public endpoint for the GKE cluster API server."
  value       = google_container_cluster.validator.endpoint
}

output "ca_certificate" {
  description = "The base64-encoded public certificate for the cluster's CA."
  value       = google_container_cluster.validator.master_auth[0].cluster_ca_certificate
  sensitive   = true
}

output "node_pool_name" {
  description = "The name of the validator node pool."
  value       = google_container_node_pool.validator_pool.name
}

output "service_account" {
  description = "The email of the service account created for the validator node pool."
  value       = google_service_account.validator_sa.email
}