data "aws_region" "current" {}

output "cluster_name" {
  value = module.eks.cluster_name
}

output "kubectl_config" {
  description = "Kubeconfig using the caller's AWS CLI credentials, without storing an EKS token."
  sensitive   = true
  value = yamlencode({
    apiVersion = "v1"
    kind       = "Config"
    clusters = [{
      name = module.eks.cluster_name
      cluster = {
        server                     = module.eks.cluster_endpoint
        certificate-authority-data = module.eks.cluster_certificate_authority_data
      }
    }]
    contexts = [{
      name = module.eks.cluster_name
      context = {
        cluster = module.eks.cluster_name
        user    = module.eks.cluster_name
      }
    }]
    current-context = module.eks.cluster_name
    users = [{
      name = module.eks.cluster_name
      user = {
        exec = {
          apiVersion      = "client.authentication.k8s.io/v1beta1"
          command         = "aws"
          interactiveMode = "Never"
          args            = ["eks", "get-token", "--cluster-name", module.eks.cluster_name, "--region", data.aws_region.current.name]
        }
      }
    }]
  })
}
