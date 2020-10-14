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

module sifnode {
    source                          = "github.com/sifchain/environments/providers/aws/eks"
    chain_id                        = {{.chain_id}}
    aws_cidr                        = {{.aws.cidr}}
    aws_cluster_version             = {{.aws.cluster.version}}
    aws_cluster_name                = "{{.aws.cluster.name}}"
    aws_cluster_ami_type            = "{{.aws.cluster.ami_type}}"
    aws_cluster_desired_capacity    = {{.aws.cluster.desired_capacity}}
    aws_cluster_max_capacity        = {{.aws.cluster.max_capacity}}
    aws_cluster_min_capacity        = {{.aws.cluster.min_capacity}}
    aws_cluster_instance_type       = {{.aws.cluster.instance_type}}
    aws_cluster_disk_size           = {{.aws.cluster.disk_size}}
}
