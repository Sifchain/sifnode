output "cluster_name" {
  description = "GKE cluster name"
  value       = google_container_cluster.primary.name
}

output "cluster_location" {
  description = "GKE cluster location"
  value       = google_container_cluster.primary.location
}

output "cluster_endpoint" {
  description = "GKE cluster endpoint"
  value       = google_container_cluster.primary.endpoint
}

output "cluster_ca_certificate" {
  description = "GKE cluster CA certificate"
  value       = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  sensitive   = true
}

output "kubectl_config" {
  description = "kubectl configuration command"
  value       = "gcloud container clusters get-credentials ${google_container_cluster.primary.name} --region ${var.region} --project ${var.project_id}"
}

output "cluster_region" {
  description = "GCP region"
  value       = var.region
}

output "service_account" {
  description = "GKE service account email"
  value       = google_service_account.sifnode.email
}

output "network_name" {
  description = "VPC network name"
  value       = google_compute_network.vpc.name
}

output "subnet_name" {
  description = "Subnet name"
  value       = google_compute_subnetwork.subnet.name
}

output "node_pool_name" {
  description = "Primary node pool name"
  value       = google_container_node_pool.primary_nodes.name
}

output "storage_class" {
  description = "SSD storage class name"
  value       = kubernetes_storage_class.ssd.metadata[0].name
}

output "helm_install_command" {
  description = "Helm install command for sifnode on GKE"
  value = <<-EOT
    # After configuring kubectl context:
    ${google_container_cluster.primary.name}
    
    # Install sifnode via Helm:
    helm install sifnode deploy/helm/sifnode \\
      --namespace sifnode \\
      --create-namespace \\
      --set provider=gke \\
      --set persistence.storageClass=gke-ssd
  EOT
}
