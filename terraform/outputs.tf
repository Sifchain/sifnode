# -----------------------------------------------------------------------------
# Sifnode GCP - Root Outputs
# -----------------------------------------------------------------------------

# VPC
output "vpc_id" {
  description = "VPC network self-link"
  value       = module.vpc.vpc_id
}

output "subnet_name" {
  description = "Subnetwork name"
  value       = module.vpc.subnet_name
}

# GKE
output "gke_cluster_name" {
  description = "GKE cluster name"
  value       = module.gke.cluster_name
}

output "gke_cluster_id" {
  description = "GKE cluster ID"
  value       = module.gke.cluster_id
}

output "gke_endpoint" {
  description = "GKE API server endpoint"
  value       = module.gke.endpoint
}

output "gke_service_account" {
  description = "GKE node pool service account email"
  value       = module.gke.service_account
}

# Cloud SQL
output "cloud_sql_instance_name" {
  description = "Cloud SQL instance name"
  value       = module.cloud_sql.instance_name
}

output "cloud_sql_connection_name" {
  description = "Cloud SQL connection name (project:region:instance)"
  value       = module.cloud_sql.instance_connection_name
}

output "cloud_sql_private_ip" {
  description = "Cloud SQL private IP address"
  value       = module.cloud_sql.private_ip
}

output "cloud_sql_database_name" {
  description = "Database name"
  value       = module.cloud_sql.database_name
}

# Load Balancer
output "load_balancer_ip" {
  description = "External IP of the load balancer"
  value       = module.load_balancer.load_balancer_ip
}

output "ssl_certificate_name" {
  description = "SSL certificate name (empty if no domain)"
  value       = module.load_balancer.ssl_certificate_name
}

# Monitoring
output "notification_channel_id" {
  description = "Monitoring notification channel ID"
  value       = module.monitoring.notification_channel_id
}

output "dashboard_id" {
  description = "Monitoring dashboard ID"
  value       = module.monitoring.dashboard_id
}

output "alert_policy_ids" {
  description = "Map of alert policy IDs"
  value       = module.monitoring.alert_policy_ids
}