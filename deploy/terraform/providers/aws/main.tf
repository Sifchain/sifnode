terraform {
  required_version = ">= 1.7, < 2.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.100.0"
    }
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# Same VPC + EKS module layout as the historical deploy/terraform/providers/aws.
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.21.0"

  name                    = var.cluster_name
  cidr                    = var.vpc_cidr
  azs                     = slice(data.aws_availability_zones.available.names, 0, 3)
  public_subnets          = [for n in range(3) : cidrsubnet(var.vpc_cidr, 4, n)]
  enable_dns_hostnames    = true
  enable_dns_support      = true
  map_public_ip_on_launch = true
  public_subnet_tags      = { "kubernetes.io/role/elb" = "1" }
  tags                    = var.tags
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.37.2"

  cluster_name                             = var.cluster_name
  cluster_version                          = var.kubernetes_version
  vpc_id                                   = module.vpc.vpc_id
  subnet_ids                               = module.vpc.public_subnets
  cluster_endpoint_public_access           = true
  cluster_endpoint_private_access          = true
  cluster_endpoint_public_access_cidrs     = var.api_server_authorized_ip_ranges
  enable_cluster_creator_admin_permissions = true
  tags                                     = var.tags

  cluster_addons = {
    coredns            = {}
    kube-proxy         = {}
    vpc-cni            = { before_compute = true }
    aws-ebs-csi-driver = {}
  }

  eks_managed_node_groups = {
    main = {
      ami_type       = "AL2023_ARM_64_STANDARD"
      instance_types = [var.instance_type]
      min_size       = 1
      max_size       = 3
      desired_size   = var.node_count
      disk_size      = 100
      iam_role_additional_policies = {
        ebs_csi = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
      }
    }
  }
}
