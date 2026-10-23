# Provider-agnostic dispatcher: keep per-cloud resources in providers/<cloud>/
# following the existing ./deploy/terraform/providers/aws pattern.
locals {
  cluster_name   = "${var.name}-${var.provider}"
  kubeconfig_cmd = {
    aws   = "aws eks update-kubeconfig --name ${local.cluster_name} --region ${var.region}"
    azure = "az aks get-credentials -g ${var.name}-rg -n ${local.cluster_name}"
    gcp   = "gcloud container clusters get-credentials ${local.cluster_name} --region ${var.region}"
  }[var.provider]
}
