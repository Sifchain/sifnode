module sifchain {
  source              = "../../deploy/terraform/providers/azure"
  location            = "southeastasia"
  resource_group_name = "sifchain-azure-{{.chainnet}}"
  cluster_name        = "sifchain-azure-{{.chainnet}}"
  dns_prefix          = "sifchain-{{.chainnet}}"
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
