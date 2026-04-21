# Private services access: allocate IP range for Cloud SQL private IP
resource "google_compute_global_address" "private_ip_range" {
  count         = var.vpc_id != "" ? 1 : 0
  name          = "${var.instance_name}-private-ip-range"
  project       = var.project_id
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = var.vpc_id
}

# Establish private service networking connection
resource "google_service_networking_connection" "private_vpc_connection" {
  count                   = var.vpc_id != "" ? 1 : 0
  network                 = var.vpc_id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_range[0].name]
}

# Cloud SQL PostgreSQL instance
resource "google_sql_database_instance" "instance" {
  name             = var.instance_name
  project          = var.project_id
  region           = var.region
  database_version = var.database_version

  settings {
    tier = var.tier

    disk_size       = var.storage_size_gb
    disk_type       = "PD_SSD"
    disk_autoresize = true

    availability_type = "REGIONAL"

    backup_configuration {
      enabled                        = true
      start_time                     = "02:00"
      point_in_time_recovery_enabled = true

      backup_retention_settings {
        retained_backups = 7
      }
    }

    ip_configuration {
      ipv4_enabled    = var.vpc_id != "" ? false : true
      private_network = var.vpc_id != "" ? var.vpc_id : null
    }
  }

  deletion_protection = var.deletion_protection

  depends_on = [
    google_service_networking_connection.private_vpc_connection
  ]
}

# Default database
resource "google_sql_database" "database" {
  name     = var.database_name
  project  = var.project_id
  instance = google_sql_database_instance.instance.name
}

# Default user
resource "google_sql_user" "user" {
  name     = "sifnode"
  project  = var.project_id
  instance = google_sql_database_instance.instance.name
  password = ""
}