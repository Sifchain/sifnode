# -----------------------------------------------------------------------------
# Sifnode GCP - Root Variables
# -----------------------------------------------------------------------------

variable "project_id" {
  description = "GCP project ID for all resources"
  type        = string
}

variable "region" {
  description = "GCP region for resource deployment"
  type        = string
  default     = "us-central1"
}

# VPC
variable "vpc_name" {
  description = "Name of the VPC network"
  type        = string
  default     = "sifnode-vpc"
}

variable "subnet_cidr" {
  description = "CIDR range for the primary subnetwork"
  type        = string
  default     = "10.0.0.0/20"
}

# GKE
variable "gke_cluster_name" {
  description = "Name of the GKE cluster"
  type        = string
}

variable "gke_node_count_min" {
  description = "Minimum number of GKE nodes (autoscaling lower bound)"
  type        = number
  default     = 3
}

variable "gke_node_count_max" {
  description = "Maximum number of GKE nodes (autoscaling upper bound)"
  type        = number
  default     = 5
}

variable "gke_machine_type" {
  description = "Machine type for GKE nodes"
  type        = string
  default     = "n2-standard-4"
}

# Cloud SQL
variable "sql_instance_name" {
  description = "Name of the Cloud SQL instance"
  type        = string
}

variable "sql_database_version" {
  description = "PostgreSQL version for Cloud SQL"
  type        = string
  default     = "POSTGRES_14"
}

variable "sql_tier" {
  description = "Machine tier for Cloud SQL"
  type        = string
  default     = "db-custom-2-4096"
}

variable "sql_database_name" {
  description = "Name of the application database"
  type        = string
  default     = "sifnodedb"
}

variable "sql_deletion_protection" {
  description = "Enable deletion protection on Cloud SQL"
  type        = bool
  default     = true
}

# Load Balancer
variable "lb_domain" {
  description = "Domain for managed SSL certificate (empty = HTTP only)"
  type        = string
  default     = ""
}

# Monitoring
variable "notification_email" {
  description = "Email address for alert notifications"
  type        = string
}

variable "alert_channels" {
  description = "Additional alert notification channel IDs"
  type        = list(string)
  default     = []
}