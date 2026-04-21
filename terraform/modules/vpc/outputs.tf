# -----------------------------------------------------------------------------
# VPC Module - Outputs
# Sifnode GCP Network
# -----------------------------------------------------------------------------

output "vpc_id" {
  description = "The self-link ID of the VPC network"
  value       = google_compute_network.vpc.id
}

output "vpc_name" {
  description = "The name of the VPC network"
  value       = google_compute_network.vpc.name
}

output "subnet_id" {
  description = "The self-link ID of the subnetwork"
  value       = google_compute_subnetwork.subnet.id
}

output "subnet_name" {
  description = "The name of the subnetwork"
  value       = google_compute_subnetwork.subnet.name
}

output "pods_range_name" {
  description = "The secondary range name for GKE pods"
  value       = "pods"
}

output "services_range_name" {
  description = "The secondary range name for GKE services"
  value       = "services"
}

output "router_name" {
  description = "The name of the Cloud Router"
  value       = google_compute_router.router.name
}