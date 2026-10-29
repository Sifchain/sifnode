terraform {
  required_version = ">= 1.7, < 2.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.56.0"
    }
  }
}

# Authenticate outside this directory with Azure CLI or ARM_* environment variables.
# AzureRM 4 requires ARM_SUBSCRIPTION_ID even when using Azure CLI authentication.
provider "azurerm" {
  features {}
}

variable "cluster_name" { type = string }
variable "location" { type = string }
variable "api_server_authorized_ip_ranges" { type = list(string) }
variable "kubernetes_version" {
  type    = string
  default = null
}

module "sifchain" {
  source                          = __MODULE_SOURCE__
  cluster_name                    = var.cluster_name
  location                        = var.location
  api_server_authorized_ip_ranges = var.api_server_authorized_ip_ranges
  kubernetes_version              = var.kubernetes_version
}

output "kubectl_config" {
  value     = module.sifchain.kubectl_config
  sensitive = true
}

output "cluster_name" { value = module.sifchain.cluster_name }
output "resource_group_name" { value = module.sifchain.resource_group_name }
