variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "default_zone" {
  description = "Default GCP zone for zonal resources"
  type        = string
  default     = "us-central1-a"
}

variable "cluster_name" {
  description = "GKE cluster name"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
  default     = "sifnode"
}

variable "cluster_version" {
  description = "GKE cluster version (use 'latest' for current default)"
  type        = string
  default     = "latest"
}

variable "release_channel" {
  description = "GKE release channel (REGULAR, STABLE, RAPID)"
  type        = string
  default     = "REGULAR"
}

variable "subnet_cidr" {
  description = "Subnet CIDR range"
  type        = string
  default     = "10.0.0.0/20"
}

variable "pods_cidr" {
  description = "CIDR range for pods"
  type        = string
  default     = "10.4.0.0/14"
}

variable "services_cidr" {
  description = "CIDR range for services"
  type        = string
  default     = "10.8.0.0/20"
}

variable "master_ipv4_cidr" {
  description = "CIDR range for the GKE master"
  type        = string
  default     = "172.16.0.0/28"
}

variable "authorized_cidrs" {
  description = "Map of authorized CIDR blocks for master access"
  type        = map(string)
  default = {
    "default" = "0.0.0.0/0"
  }
}

variable "enable_private_nodes" {
  description = "Enable private nodes"
  type        = bool
  default     = false
}

variable "enable_autoscaling" {
  description = "Enable cluster autoscaling"
  type        = bool
  default     = true
}

variable "initial_node_count" {
  description = "Initial number of nodes"
  default     = 1
}

variable "desired_capacity" {
  description = "Desired node count"
  default     = 1
}

variable "min_capacity" {
  description = "Minimum node count"
  default     = 1
}

variable "max_capacity" {
  description = "Maximum node count"
  default     = 5
}

variable "machine_type" {
  description = "GCE machine type for nodes"
  default     = "e2-standard-2"
}

variable "image_type" {
  description = "Node image type"
  default     = "COS_CONTAINERD"
}

variable "disk_size" {
  description = "Node boot disk size in GB"
  default     = 100
}

variable "disk_type" {
  description = "Node boot disk type"
  default     = "pd-standard"
}

variable "data_disk_size" {
  description = "Persistent data disk size in GB"
  default     = 50
}

variable "storage_size" {
  description = "PVC storage size"
  default     = "50Gi"
}

variable "min_cpu" {
  description = "Minimum CPU for autoscaler"
  default     = 1
}

variable "max_cpu" {
  description = "Maximum CPU for autoscaler"
  default     = 8
}

variable "min_memory_mb" {
  description = "Minimum memory in MB for autoscaler"
  default     = 4096
}

variable "max_memory_mb" {
  description = "Maximum memory in MB for autoscaler"
  default     = 32768
}

variable "maintenance_window" {
  description = "Daily maintenance window start time (UTC)"
  default     = "03:00"
}

variable "ssh_source_ranges" {
  description = "CIDR ranges allowed for SSH access"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
