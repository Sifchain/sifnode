variable "project_id" {
  description = "Google Cloud project hosting the Sifnode cluster."
  type        = string
}

variable "region" {
  description = "GCP region for the regional GKE cluster."
  type        = string
  default     = "us-central1"
}

variable "cluster_name" {
  description = "GKE cluster name."
  type        = string
  default     = "sifnode"
}

variable "network_cidr" {
  description = "Primary subnet CIDR."
  type        = string
  default     = "10.20.0.0/20"
}

variable "pods_cidr" {
  description = "Secondary range for GKE pods."
  type        = string
  default     = "10.24.0.0/16"
}

variable "services_cidr" {
  description = "Secondary range for GKE services."
  type        = string
  default     = "10.28.0.0/20"
}

variable "node_machine_type" {
  description = "GCE machine type for Sifnode workers."
  type        = string
  default     = "e2-standard-4"
}

variable "node_count" {
  description = "Initial node count."
  type        = number
  default     = 2
}

variable "disk_size_gb" {
  description = "Boot disk size per worker."
  type        = number
  default     = 100
}

variable "labels" {
  description = "Labels applied to cluster resources."
  type        = map(string)
  default = {
    application = "sifnode"
    managed-by  = "terraform"
  }
}
