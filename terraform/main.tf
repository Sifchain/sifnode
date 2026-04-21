# -----------------------------------------------------------------------------
# Sifnode GCP - Root Module
# Wires all sub-modules together
# -----------------------------------------------------------------------------

# VPC Network
module "vpc" {
  source = "./modules/vpc"

  project_id  = var.project_id
  region      = var.region
  vpc_name    = var.vpc_name
  subnet_cidr = var.subnet_cidr
}

# GKE Cluster
module "gke" {
  source = "./modules/gke"

  project_id     = var.project_id
  region         = var.region
  cluster_name   = var.gke_cluster_name
  node_count_min = var.gke_node_count_min
  node_count_max = var.gke_node_count_max
  machine_type   = var.gke_machine_type
  network_vpc_id = module.vpc.vpc_id
  subnetwork_id  = module.vpc.subnet_id

  depends_on = [module.vpc]
}

# Cloud SQL
module "cloud_sql" {
  source = "./modules/cloud-sql"

  project_id          = var.project_id
  region              = var.region
  instance_name       = var.sql_instance_name
  database_version    = var.sql_database_version
  tier                = var.sql_tier
  database_name       = var.sql_database_name
  vpc_id              = module.vpc.vpc_id
  deletion_protection = var.sql_deletion_protection

  depends_on = [module.vpc]
}

# Load Balancer
module "load_balancer" {
  source = "./modules/load-balancer"

  project_id = var.project_id
  name       = "sifnode-lb"
  region     = var.region
  domain     = var.lb_domain

  depends_on = [module.gke]
}

# Monitoring
module "monitoring" {
  source = "./modules/monitoring"

  project_id         = var.project_id
  notification_email = var.notification_email
  alert_channels     = var.alert_channels
  cluster_name       = var.gke_cluster_name

  depends_on = [module.gke]
}