terraform {
  required_version = ">= 1.7, < 2.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.56.0"
    }
  }
}

resource "azurerm_resource_group" "sifchain" {
  name     = var.cluster_name
  location = var.location
  tags     = var.tags
}

resource "azurerm_kubernetes_cluster" "sifchain" {
  name                = var.cluster_name
  location            = azurerm_resource_group.sifchain.location
  resource_group_name = azurerm_resource_group.sifchain.name
  dns_prefix          = var.cluster_name
  kubernetes_version  = var.kubernetes_version
  sku_tier            = var.sku_tier
  tags                = var.tags

  role_based_access_control_enabled = true
  automatic_upgrade_channel         = "patch"

  default_node_pool {
    name                        = "system"
    temporary_name_for_rotation = "systemtmp"
    node_count                  = var.node_count
    vm_size                     = var.vm_size
    os_disk_size_gb             = 100
    os_sku                      = "Ubuntu"
    type                        = "VirtualMachineScaleSets"
    upgrade_settings {
      max_surge = "10%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  api_server_access_profile {
    authorized_ip_ranges = var.api_server_authorized_ip_ranges
  }

  network_profile {
    network_plugin      = "azure"
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
  }

  storage_profile {
    disk_driver_enabled = true
  }
}
