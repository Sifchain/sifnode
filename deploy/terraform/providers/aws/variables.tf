variable "cluster_name" {
  type        = string
  description = "EKS cluster name."
}

variable "kubernetes_version" {
  type        = string
  default     = null
  description = "EKS-supported version; null selects the AWS default."
}

variable "api_server_authorized_ip_ranges" {
  type        = list(string)
  description = "Administrator public IPv4 CIDRs for the EKS API."
  validation {
    condition = length(var.api_server_authorized_ip_ranges) > 0 && alltrue([
      for cidr in var.api_server_authorized_ip_ranges :
      can(cidrnetmask(cidr)) && try(tonumber(split("/", cidr)[1]) > 0, false)
    ])
    error_message = "Provide valid administrator IPv4 CIDRs narrower than /0."
  }
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "node_count" {
  type    = number
  default = 1
  validation {
    condition     = var.node_count >= 1 && var.node_count <= 3 && floor(var.node_count) == var.node_count
    error_message = "node_count must be an integer from 1 to 3."
  }
}

variable "instance_type" {
  type        = string
  default     = "t4g.xlarge"
  description = "ARM64 worker type; the sifnode image must support ARM64."
}

variable "tags" {
  type    = map(string)
  default = { Project = "sifchain", Terraform = "true" }
}
