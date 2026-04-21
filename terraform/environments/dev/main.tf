# Sifnode GCP - Dev Environment
# Wires root module with dev-specific values

module "sifnode_dev" {
  source = "../../"

  project_id              = var.project_id
  region                  = "us-central1"
  vpc_name                = "sifnode-dev-vpc"
  subnet_cidr             = "10.0.0.0/20"
  gke_cluster_name        = "sifnode-dev"
  gke_node_count_min      = 3
  gke_node_count_max      = 5
  gke_machine_type        = "n2-standard-4"
  sql_instance_name       = "sifnode-db-dev"
  sql_database_version    = "POSTGRES_14"
  sql_tier                = "db-custom-2-4096"
  sql_database_name       = "sifnodedb"
  sql_deletion_protection = false
  lb_domain               = ""
  notification_email      = var.notification_email
  alert_channels          = []
}