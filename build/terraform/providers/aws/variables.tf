variable "region" {
  description = "AWS region"
  type        = string
}

variable "az" {
  description = "AWS availability zones"
  default     = ["a", "b", "c"]
}

variable "vpc_cidr" {
  description = "VPC cidr"
  type        = string
  default     = "10.0.0.0/16"
}

variable "cluster_version" {
  description = "EKS cluster version"
  type        = string
  default     = 1.18
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "tags" {
  description = "Tags"
  type        = map(string)
}

variable "desired_capacity" {
  description = "desired kubes nodes pre cluster"
  default = 1
}

variable "max_capacity" {
  description = "Max kubes nodes pre cluster"
  default = 3
}

variable "min_capacity" {
  description = "Min capacity of nodes pre kubes cluster"
  default = 1
}

variable "instance_type" {
  default = "t2.medium"
}

variable "ami_type" {
  default = "AL2_x86_64"
}

variable "disk_size" {
  default = 100
}

