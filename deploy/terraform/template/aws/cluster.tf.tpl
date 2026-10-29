terraform {
  required_version = ">= 1.7, < 2.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.100.0"
    }
  }
}

provider "aws" {
  region = var.region
}

variable "cluster_name" { type = string }
variable "region" { type = string }
variable "api_server_authorized_ip_ranges" { type = list(string) }
variable "kubernetes_version" {
  type    = string
  default = null
}

module "sifchain" {
  source                          = __MODULE_SOURCE__
  cluster_name                    = var.cluster_name
  api_server_authorized_ip_ranges = var.api_server_authorized_ip_ranges
  kubernetes_version              = var.kubernetes_version
}

output "kubectl_config" {
  value     = module.sifchain.kubectl_config
  sensitive = true
}

output "cluster_name" { value = module.sifchain.cluster_name }
