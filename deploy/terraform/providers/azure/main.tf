# AKS provider — mirrors deploy/terraform/providers/aws structure.
module "cluster" {
  source       = "../../modules/cluster"
  name         = var.name
  provider     = "azure"
  region       = var.location
  node_count   = var.node_count
  node_size    = var.node_size
  kube_version = var.kube_version
  tags         = var.tags
}
resource "azurerm_resource_group" "rg" { name = "${var.name}-rg" location = var.location }
resource "azurerm_kubernetes_cluster" "aks" {
  name                = module.cluster.cluster_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = var.name
  kubernetes_version  = var.kube_version
  default_node_pool { name = "default" node_count = var.node_count vm_size = var.node_size }
  identity { type = "SystemAssigned" }
  tags = var.tags
}
