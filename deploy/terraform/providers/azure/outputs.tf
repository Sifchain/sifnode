output "kubectl_config" {
  description = "kubectl config"
  value       = azurerm_kubernetes_cluster.cluster.kube_config_raw
  sensitive   = true
}

output "resource_group_name" {
  value = azurerm_resource_group.cluster.name
}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.cluster.name
}
