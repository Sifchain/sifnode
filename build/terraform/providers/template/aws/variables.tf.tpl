variable "region" {
  description = "AWS region"
  type        = string
  default     = "{{.aws.region}}"
}

variable "az" {
  description = "AWS availability zones"
  default     = ["a", "b", "c"]
}

variable "vpc_cidr" {
  description = "VPC cidr"
  type        = string
  default     = "{{.aws.cidr}}"
}

variable "cluster_version" {
  description = "EKS cluster version"
  type        = string
  default     = "{{.aws.cluster.version}}"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "{{.aws.cluster.name}}"
}

variable "tags" {
  description = "Tags"
  type        = map(string)
  default = {
    Terraform = true
    Sifnode   = true
    ChainID   = "{{.chain_id}}"
  }
}

variable "node_group_settings" {
  description = "Cluster node group settings"
  type        = map(string)
  default = {
    ami_type         = "{{.aws.cluster.ami_type}}"
    desired_capacity = {{.aws.cluster.desired_capacity}}
    max_capacity     = {{.aws.cluster.max_capacity}}
    min_capacity     = {{.aws.cluster.min_capacity}}
    instance_type    = "{{.aws.cluster.instance_type}}"
    disk_size        = {{.aws.cluster.disk_size}}
  }
}
