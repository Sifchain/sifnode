variable "name" {
  description = "Cluster name"
  type        = string
}

variable "region" {
  description = "AWS region (e.g. us-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "node_count" {
  description = "Desired node count"
  type        = number
  default     = 2
}

variable "node_size" {
  description = "EC2 instance type for nodes"
  type        = string
  default     = "t3.medium"
}

variable "kube_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.29"
}

variable "subnet_ids" {
  description = "VPC subnet IDs for EKS and node groups"
  type        = list(string)
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default     = {}
}
