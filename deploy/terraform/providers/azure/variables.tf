variable "location" {
  description = "Azure region"
  type        = string
}

variable "resource_group_name" {
  description = "Resource group for the AKS cluster"
  type        = string
}

variable "cluster_name" {
  description = "AKS cluster name"
  type        = string
}

variable "dns_prefix" {
  description = "AKS DNS prefix"
  type        = string
}

variable "kubernetes_version" {
  description = "AKS kubernetes version"
  type        = string
  default     = "1.35.5"
}

variable "node_count" {
  description = "Number of nodes in the default pool"
  type        = number
  default     = 1
}

variable "vm_size" {
  description = "VM size for AKS worker nodes"
  type        = string
  default     = "Standard_B2s_v2"
}

variable "os_disk_size_gb" {
  description = "OS disk size for worker nodes"
  type        = number
  default     = 128
}

variable "vnet_cidr" {
  description = "Virtual network CIDR"
  type        = string
  default     = "10.20.0.0/16"
}

variable "subnet_cidr" {
  description = "Cluster subnet CIDR"
  type        = string
  default     = "10.20.1.0/24"
}

variable "tags" {
  description = "Tags"
  type        = map(string)
  default     = {}
}
