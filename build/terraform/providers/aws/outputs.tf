output "kubectl_config" {
  description = "kubectl config"
  value       = module.eks.kubeconfig
}
