output "cluster_name" {
  value = google_container_cluster.sifnode.name
}

output "cluster_location" {
  value = google_container_cluster.sifnode.location
}

output "network_name" {
  value = google_compute_network.sifnode.name
}

output "subnetwork_name" {
  value = google_compute_subnetwork.sifnode.name
}

output "get_credentials_command" {
  value = "gcloud container clusters get-credentials ${google_container_cluster.sifnode.name} --region ${google_container_cluster.sifnode.location} --project ${var.project_id}"
}
