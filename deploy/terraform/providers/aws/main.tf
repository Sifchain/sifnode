terraform {
  required_version = ">= 0.13"
}

provider "aws" {
  region = var.region
  profile = var.profile
}

provider "kubernetes" {
  host                   = element(concat(data.aws_eks_cluster.cluster[*].endpoint, list("")), 0)
  cluster_ca_certificate = base64decode(element(concat(data.aws_eks_cluster.cluster[*].certificate_authority.0.data, list("")), 0))
  token                  = element(concat(data.aws_eks_cluster_auth.cluster[*].token, list("")), 0)
  load_config_file       = false
  version                = "~> 1.9"
}

data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

data "aws_iam_role" "cluster" {
  name = module.eks.worker_iam_role_name
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
locals {

node_final_name = "${var.cluster_name}-node-group"
}
module "eks" {
  source       = "terraform-aws-modules/eks/aws"
  cluster_name = var.cluster_name
  subnets      = module.vpc.public_subnets
  vpc_id       = module.vpc.vpc_id
  tags         = merge({ "Name" = var.cluster_name }, var.tags)
  version      = "13.2.1"

  node_groups_defaults = {
    ami_type  = var.ami_type
    disk_size = var.disk_size
  }

  node_groups = {
    main = {
      name = var.customize_node_group_name == "yes" ? var.node_group_name : local.node_final_name
      desired_capacity = var.desired_capacity
      max_capacity     = var.max_capacity
      min_capacity     = var.min_capacity
      instance_type    = var.instance_type

      k8s_labels = {
        Environment = "${var.cluster_name}-${var.region}"
      }
      additional_tags = var.tags
    }
  }

  cluster_version  = var.cluster_version
  write_kubeconfig = true
}

resource "aws_iam_policy" "policy" {
  name   = var.cluster_name
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:AttachVolume",
        "ec2:CreateSnapshot",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:DeleteSnapshot",
        "ec2:DeleteTags",
        "ec2:DeleteVolume",
        "ec2:DescribeInstances",
        "ec2:DescribeSnapshots",
        "ec2:DescribeTags",
        "ec2:DescribeVolumes",
        "ec2:DetachVolume"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "attach" {
  name       = var.cluster_name
  roles      = [
    data.aws_iam_role.cluster.id
  ]
  policy_arn = aws_iam_policy.policy.arn
}

data "aws_subnet_ids" "subnets" {
  vpc_id     = module.vpc.vpc_id
  depends_on = [module.vpc]
}

locals {
  subnet_ids_string = join(",", data.aws_subnet_ids.subnets.ids)
  subnet_ids_list   = split(",", local.subnet_ids_string)
  depends_on        = [module.vpc]
}
