terraform {
  required_providers {
    kustomization = {
      source  = "kbst/kustomization"
      version = "~> 0.2.2"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
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

provider "kubectl" {
  host                   = element(concat(data.aws_eks_cluster.cluster[*].endpoint, list("")), 0)
  cluster_ca_certificate = base64decode(element(concat(data.aws_eks_cluster.cluster[*].certificate_authority.0.data, list("")), 0))
  token                  = element(concat(data.aws_eks_cluster_auth.cluster[*].token, list("")), 0)
  load_config_file       = false
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

module "eks" {
  source       = "terraform-aws-modules/eks/aws"
  cluster_name = var.cluster_name
  subnets      = module.vpc.public_subnets
  vpc_id       = module.vpc.vpc_id
  tags         = merge({ "Name" = var.cluster_name }, var.tags)

  node_groups_defaults = {
    ami_type  = var.ami_type
    disk_size = var.disk_size
  }

  node_groups = {
    main = {
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

resource "aws_security_group" "security_group" {
  name        = "NFS security group for Kubernetes"
  description = "EFS Plugin"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "NFS inbound rule"
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = var.cluster_name
  }
}

resource "aws_efs_file_system" "efs_file_system" {
  creation_token = var.cluster_name
  depends_on     = [module.vpc]
  tags           = {
    Name = var.cluster_name
  }

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

resource "aws_efs_mount_target" "mount" {
  count           = length(var.az)
  file_system_id  = aws_efs_file_system.efs_file_system.id
  subnet_id       = element(local.subnet_ids_list, count.index)
  security_groups = [aws_security_group.security_group.id]
  depends_on      = [module.vpc]
}

provider "kustomization" {
  kubeconfig_raw = module.eks.kubeconfig
}

data "kustomization" "efs_manifests" {
  path = var.efs_csi_driver
}

data "kustomization" "ebs_manifests" {
  path = var.ebs_csi_driver
}

resource "kustomization_resource" "efs_csi_driver" {
  for_each = data.kustomization.efs_manifests.ids
  manifest = data.kustomization.efs_manifests.manifests[each.value]
  lifecycle {
      ignore_changes = all
  }
}

resource "kustomization_resource" "ebs_csi_driver" {
  for_each = data.kustomization.ebs_manifests.ids
  manifest = data.kustomization.ebs_manifests.manifests[each.value]
  lifecycle {
      ignore_changes = all
  }
}

resource "kubectl_manifest" "efs_storageclass" {
  yaml_body = <<YAML
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: efs-sc
provisioner: efs.csi.aws.com
YAML
  depends_on = [module.eks]
  lifecycle {
      ignore_changes = all
  }
}

resource "kubectl_manifest" "ebs_storageclass" {
  yaml_body = <<YAML
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ebs-sc
provisioner: ebs.csi.aws.com
volumeBindingMode: WaitForFirstConsumer
YAML
  depends_on = [module.eks]
  lifecycle {
      ignore_changes = all
  }
}

resource "kubectl_manifest" "efs_pv_sifnoded" {
  yaml_body = <<YAML
apiVersion: v1
kind: PersistentVolume
metadata:
  name: "${var.efs_pv_sifnoded_name}"
spec:
  capacity:
    storage: "${var.efs_pv_capacity}"
  volumeMode: Filesystem
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: "${var.efs_pv_storageclass}"
  csi:
    driver: efs.csi.aws.com
    volumeHandle: "${aws_efs_file_system.efs_file_system.id}"
YAML
  depends_on = [module.eks]
}

resource "kubectl_manifest" "efs_pv_sifnodecli" {
  yaml_body = <<YAML
apiVersion: v1
kind: PersistentVolume
metadata:
  name: "${var.efs_pv_sifnodecli_name}"
spec:
  capacity:
    storage: "${var.efs_pv_capacity}"
  volumeMode: Filesystem
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: "${var.efs_pv_storageclass}"
  csi:
    driver: efs.csi.aws.com
    volumeHandle: "${aws_efs_file_system.efs_file_system.id}"
YAML
  depends_on = [module.eks]
}
