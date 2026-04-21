variable "project_id" {
  description = "GCP project ID where Cloud SQL will be created"
  type        = string
}

variable "region" {
  description = "GCP region for the Cloud SQL instance"
  type        = string
}

variable "instance_name" {
  description = "Name of the Cloud SQL instance"
  type        = string
}

variable "database_version" {
  description = "PostgreSQL database version"
  type        = string
  default     = "POSTGRES_14"
}

variable "tier" {
  description = "Machine tier for the Cloud SQL instance"
  type        = string
  default     = "db-custom-2-4096"
}

variable "storage_size_gb" {
  description = "Storage size in GB for the Cloud SQL instance"
  type        = number
  default     = 100
}

variable "database_name" {
  description = "Name of the default database to create"
  type        = string
  default     = "sifnodedb"
}

variable "vpc_id" {
  description = "VPC network ID for private services access"
  type        = string
  default     = ""
}

variable "deletion_protection" {
  description = "Whether to enable deletion protection on the Cloud SQL instance"
  type        = bool
  default     = true
}