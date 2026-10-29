variable "cluster_name" {
  description = "AKS cluster and dedicated resource group name."
  type        = string
  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{2,60}[a-z0-9]$", var.cluster_name))
    error_message = "Use 4-62 lowercase letters, digits or hyphens, starting with a letter and ending with a letter or digit."
  }
}

variable "location" {
  description = "Azure location, for example eastus."
  type        = string
}

variable "kubernetes_version" {
  description = "AKS-supported Kubernetes version; null selects Azure's recommended version."
  type        = string
  default     = null
}

variable "api_server_authorized_ip_ranges" {
  description = "Administrator public IPv4 CIDRs permitted to reach the Kubernetes API."
  type        = list(string)
  validation {
    condition = length(var.api_server_authorized_ip_ranges) > 0 && alltrue([
      for cidr in var.api_server_authorized_ip_ranges :
      can(cidrnetmask(cidr)) && try(tonumber(split("/", cidr)[1]) > 0, false)
    ])
    error_message = "Provide at least one valid administrator IPv4 CIDR; /0 is not permitted."
  }
}

variable "node_count" {
  description = "Number of AKS worker nodes."
  type        = number
  default     = 1
  validation {
    condition     = var.node_count >= 1 && floor(var.node_count) == var.node_count
    error_message = "node_count must be a positive integer."
  }
}

variable "vm_size" {
  description = "Worker VM size; choose sufficient RAM for the target chain."
  type        = string
  default     = "Standard_D4s_v5"
}

variable "sku_tier" {
  description = "AKS control plane tier; Standard is recommended when an uptime SLA is required."
  type        = string
  default     = "Free"
  validation {
    condition     = contains(["Free", "Standard", "Premium"], var.sku_tier)
    error_message = "sku_tier must be Free, Standard or Premium."
  }
}

variable "tags" {
  description = "Resource tags, following the AWS module's convention."
  type        = map(string)
  default     = { Project = "sifchain", Terraform = "true" }
}
