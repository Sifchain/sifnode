# EKS provider — canonical layout that azure/ and gcp/ mirror.
module "cluster" {
  source       = "../../modules/cluster"
  name         = var.name
  provider     = "aws"
  region       = var.region
  node_count   = var.node_count
  node_size    = var.node_size
  kube_version = var.kube_version
  tags         = var.tags
}

resource "aws_eks_cluster" "eks" {
  name     = module.cluster.cluster_name
  version  = var.kube_version
  role_arn = aws_iam_role.eks.arn

  vpc_config { subnet_ids = var.subnet_ids }

  tags = var.tags
}

resource "aws_eks_node_group" "default" {
  cluster_name    = aws_eks_cluster.eks.name
  node_group_name = "default"
  node_role_arn   = aws_iam_role.nodes.arn
  subnet_ids      = var.subnet_ids
  instance_types  = [var.node_size]

  scaling_config {
    desired_size = var.node_count
    min_size     = 1
    max_size     = var.node_count + 1
  }
}

resource "aws_iam_role" "eks" {
  name = "${var.name}-eks-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{ Action = "sts:AssumeRole", Effect = "Allow",
      Principal = { Service = "eks.amazonaws.com" } }]
  })
}

resource "aws_iam_role_policy_attachment" "eks_cluster" {
  role       = aws_iam_role.eks.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}

resource "aws_iam_role" "nodes" {
  name = "${var.name}-eks-nodes-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{ Action = "sts:AssumeRole", Effect = "Allow",
      Principal = { Service = "ec2.amazonaws.com" } }]
  })
}

resource "aws_iam_role_policy_attachment" "nodes" {
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
    "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
  ])
  role       = aws_iam_role.nodes.name
  policy_arn = each.value
}
