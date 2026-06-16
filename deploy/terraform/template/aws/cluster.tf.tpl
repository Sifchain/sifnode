module sifchain {
  source       = "../../deploy/terraform/providers/aws"
  region       = "us-west-2"
  cluster_name = "sifchain-aws-{{.chainnet}}"
  tags = {
    Terraform = "true"
    Project   = "sifchain"
    Chainnet  = "{{.chainnet}}"
  }
}

output "kubectl_config" {
  value     = module.sifchain.kubectl_config
  sensitive = true
}
