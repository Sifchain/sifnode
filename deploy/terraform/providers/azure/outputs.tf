output "kubectl_config" {
  description = "Local account kubeconfig. Contains credentials; keep state and this output private."
  value       = azurerm_kubernetes_cluster.sifchain.kube_config_raw
  sensitive   = true
}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.sifchain.name
}

output "resource_group_name" {
  value = azurerm_resource_group.sifchain.name
}
