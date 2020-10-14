data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"

  name           = var.cluster_name
  cidr           = var.vpc_cidr
  azs            = [for az in var.az : format("%s%s", var.region, az)]
  public_subnets = [cidrsubnet(var.vpc_cidr, 4, 1), cidrsubnet(var.vpc_cidr, 4, 2), cidrsubnet(var.vpc_cidr, 4, 3)]

  enable_dns_hostnames = true
  enable_dns_support   = true

  map_public_ip_on_launch = true

  tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb" = "1"
  }
}

module "eks" {
  source       = "terraform-aws-modules/eks/aws"
  cluster_name = var.cluster_name
  subnets      = module.vpc.public_subnets
  vpc_id       = module.vpc.vpc_id
  tags         = merge({ "Name" = var.cluster_name }, var.tags)

  node_groups_defaults = {
    ami_type  = var.node_group_settings["ami_type"]
    disk_size = var.node_group_settings["disk_size"]
  }

  node_groups = {
    main = {
      desired_capacity = var.node_group_settings["desired_capacity"]
      max_capacity     = var.node_group_settings["max_capacity"]
      min_capacity     = var.node_group_settings["min_capacity"]
      instance_type    = var.node_group_settings["instance_type"]

      k8s_labels = {
        Environment = "${var.cluster_name}-${var.region}"
      }
      additional_tags = var.tags
    }
  }

  cluster_version  = var.cluster_version
  write_kubeconfig = true
}
