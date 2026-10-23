# GKE provider — mirrors deploy/terraform/providers/aws structure.
module "cluster" {
  source       = "../../modules/cluster"
  name         = var.name
  provider     = "gcp"
  region       = var.region
  node_count   = var.node_count
  node_size    = var.node_size
  kube_version = var.kube_version
  tags         = var.tags
}
resource "google_container_cluster" "gke" {
  name     = module.cluster.cluster_name
  location = var.region
  initial_node_count = var.node_count
  min_master_version = var.kube_version
  node_config { machine_type = var.node_size }
}
