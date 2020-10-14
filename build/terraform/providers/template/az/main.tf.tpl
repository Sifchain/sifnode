provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = var.rg_name
  location = var.rg_location
}

module "aks" {
  //source            = "Azure/aks/azurerm"
  source              = "github.com/utx0/terraform-azurerm-aks"
  resource_group_name = azurerm_resource_group.rg.name
  client_id           = var.service_principal_client_id
  client_secret       = var.service_principal_client_secret
  prefix              = var.aks_prefix
  agents_size         = var.agents_size
  agents_count        = var.agents_count
}

