# Sifnode GCP - Prod Environment
# Wires root module with prod-specific values (HA, more nodes)

module "sifnode_prod" {
  source = "../../"

  project_id              = var.project_id
  region                  = "us-central1"
  vpc_name                = "sifnode-prod-vpc"
  subnet_cidr             = "10.1.0.0/20"
  gke_cluster_name        = "sifnode-prod"
  gke_node_count_min      = 5
  gke_node_count_max      = 7
  gke_machine_type        = "n2-standard-4"
  sql_instance_name       = "sifnode-db-prod"
  sql_database_version    = "POSTGRES_14"
  sql_tier                = "db-custom-4-8192"
  sql_database_name       = "sifnodedb"
  sql_deletion_protection = true
  lb_domain               = var.lb_domain
  notification_email      = var.notification_email
  alert_channels          = var.alert_channels
}