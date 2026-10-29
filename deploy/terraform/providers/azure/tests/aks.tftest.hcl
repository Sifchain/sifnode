# All runs use Terraform's provider mock; no Azure credentials or resources.
mock_provider "azurerm" {}

variables {
  cluster_name                    = "sifchain-azure-testnet"
  location                        = "eastus"
  api_server_authorized_ip_ranges = ["192.0.2.1/32"]
}

run "aks_defaults" {
  command = plan

  assert {
    condition     = azurerm_kubernetes_cluster.sifchain.role_based_access_control_enabled
    error_message = "AKS must enable Kubernetes RBAC."
  }
  assert {
    condition     = azurerm_kubernetes_cluster.sifchain.identity[0].type == "SystemAssigned"
    error_message = "AKS must have a managed identity."
  }
  assert {
    condition     = azurerm_kubernetes_cluster.sifchain.storage_profile[0].disk_driver_enabled
    error_message = "The Azure Disk CSI driver is required for the node's persistent volume."
  }
  assert {
    condition     = azurerm_kubernetes_cluster.sifchain.network_profile[0].network_plugin_mode == "overlay"
    error_message = "Use Azure CNI overlay networking."
  }
  assert {
    condition     = azurerm_kubernetes_cluster.sifchain.api_server_access_profile[0].authorized_ip_ranges == toset(["192.0.2.1/32"])
    error_message = "The API must be limited to the supplied administrator CIDRs."
  }
  assert {
    condition     = azurerm_kubernetes_cluster.sifchain.default_node_pool[0].node_count == 1
    error_message = "The initial scaffold should provision one worker."
  }
}

run "custom_cluster" {
  command = plan
  variables {
    node_count         = 3
    vm_size            = "Standard_D8s_v5"
    sku_tier           = "Standard"
    kubernetes_version = "1.33"
    tags               = { Project = "operator-test" }
  }
  assert {
    condition = (
      azurerm_kubernetes_cluster.sifchain.default_node_pool[0].node_count == 3 &&
      azurerm_kubernetes_cluster.sifchain.default_node_pool[0].vm_size == "Standard_D8s_v5" &&
      azurerm_kubernetes_cluster.sifchain.sku_tier == "Standard" &&
      azurerm_kubernetes_cluster.sifchain.kubernetes_version == "1.33" &&
      azurerm_resource_group.sifchain.tags.Project == "operator-test"
    )
    error_message = "Cluster sizing, version, tier and tags must be configurable."
  }
}

run "reject_open_api" {
  command = plan
  variables {
    api_server_authorized_ip_ranges = ["0.0.0.0/0"]
  }
  expect_failures = [var.api_server_authorized_ip_ranges]
}

run "reject_missing_api_ranges" {
  command = plan
  variables {
    api_server_authorized_ip_ranges = []
  }
  expect_failures = [var.api_server_authorized_ip_ranges]
}

run "reject_invalid_cidr" {
  command = plan
  variables {
    api_server_authorized_ip_ranges = ["not-a-network"]
  }
  expect_failures = [var.api_server_authorized_ip_ranges]
}

run "reject_no_workers" {
  command = plan
  variables {
    node_count = 0
  }
  expect_failures = [var.node_count]
}

run "reject_fractional_workers" {
  command = plan
  variables {
    node_count = 1.5
  }
  expect_failures = [var.node_count]
}
