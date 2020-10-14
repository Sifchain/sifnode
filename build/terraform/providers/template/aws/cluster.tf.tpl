provider "aws" {
  region = {{.aws.region}}
}

provider "kubernetes" {
  host                   = element(concat(data.aws_eks_cluster.cluster[*].endpoint, list("")), 0)
  cluster_ca_certificate = base64decode(element(concat(data.aws_eks_cluster.cluster[*].certificate_authority.0.data, list("")), 0))
  token                  = element(concat(data.aws_eks_cluster_auth.cluster[*].token, list("")), 0)
  load_config_file       = false
  version                = "~> 1.9"
}

module sifchain {
    source                  = "../../build/terraform/providers/aws"
    chainnet                = "{{.chainnet}}"
    vpc_cidr                = "{{.aws.cidr}}"
    cluster_version         = {{.aws.cluster.version}}
    cluster_name            = "{{.aws.cluster.name}}"
    tags = {
        Terraform           = true
        Sifnode             = true
        ChainNet            = "{{.chainnet}}"
    }
    node_group_settings = {
        ami_type            = "{{.aws.cluster.ami_type}}"
        desired_capacity    = {{.aws.cluster.desired_capacity}}
        max_capacity        = {{.aws.cluster.max_capacity}}
        min_capacity        = {{.aws.cluster.min_capacity}}
        instance_type       = "{{.aws.cluster.instance_type}}"
        disk_size           = {{.aws.cluster.disk_size}}
    }
}
