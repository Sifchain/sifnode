variable "vpc_cidr" {
  description = "VPC cidr"
  type = string
  default = "10.0.0.0/16"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type = string
}

variable "dns_prefix" {
  description = "DNS prefix"
  type = string
  default = "k8s-sifchain-1"
}

variable "cluster_location" {
  description = "Deploy location"
  type = string
  default = "West US"
}

variable "cluster_resource_group" {
  description = "Cluster resource group"
  type = string
  default = "sifchain-resource-group-1"
}

variable "disk_size"{
  description = "Disk size(GB)"
  default = 100
}

variable "node_count" {
  description = "Node counts"
  default = 2
}

variable "vm_size" {
  description = "VM sizes"
  default = "Standard_D2_v2"
}
variable "node_pool_name"{
  description = "Node pool name"
  default = "default"
}

variable "project_name" {
  description = "sifchain"
  type = string
  default = "sifchain"
}

variable "tags" {
  description = "Tags"
  type = map(string)
}


variable "cluster_version" {
  description = "AKS cluster version"
  type = string
  default = "1.19.11"
}

variable "environment" {
  description = "Environment"
  type = string
  default = "Production"
}

variable "vnet_name" {
  description = "AKS Virtual Network Name"
  type = string
  default = "vnet-1"
}

variable "subnet_name" {
  description = "AKS SubNet Name"
  type = string
  default = "subnet-1"
}

