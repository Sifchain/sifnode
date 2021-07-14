variable "region" {
  description = "AWS region"
  type        = string
}

variable "az" {
  description = "AWS availability zones"
  default     = ["a", "b", "c"]
}

variable "vpc_cidr" {
  description = "VPC cidr"
  type        = string
  default     = "10.0.0.0/16"
}

variable "cluster_version" {
  description = "EKS cluster version"
  type        = string
  default     = 1.19
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "tags" {
  description = "Tags"
  type        = map(string)
}

variable "desired_capacity_0" {
  description = "Desired nodes per cluster for node group main-0"
  default     = 1
}

variable "desired_capacity_1" {
  description = "Desired nodes per cluster for node group main-1"
  default     = 1
}

variable "desired_capacity_2" {
  description = "Desired nodes per cluster for node group main-2"
  default     = 1
}

variable "max_capacity_0" {
  description = "Max nodes per cluster for node group main-0"
  default     = 3
}

variable "max_capacity_1" {
  description = "Max nodes per cluster for node group main-1"
  default     = 3
}

variable "max_capacity_2" {
  description = "Max nodes per cluster for node group main-2"
  default     = 3
}

variable "min_capacity_0" {
  description = "Min nodes per cluster for node group main-0"
  default     = 1
}

variable "min_capacity_1" {
  description = "Min nodes per cluster for node group main-1"
  default     = 1
}

variable "min_capacity_2" {
  description = "Min nodes per cluster for node group main-2"
  default     = 1
}

variable "ami_type" {
  default = "AL2_x86_64"
}

variable "disk_size" {
  default = 100
}

variable "policy_name" {
  description = "Policy name for the AWS EBS CSI Driver"
  type        = string
  default     = "Amazon_EBS_CSI_Driver"
}

variable "ebs_csi_driver" {
  description = "GitHub path to the CSI plugin for data persistence"
  default     = "github.com/kubernetes-sigs/aws-ebs-csi-driver/deploy/kubernetes/overlays/stable/?ref=master"
}

variable "efs_csi_driver" {
  description = "GitHub path to the CSI plugin for data persistence"
  default     = "github.com/kubernetes-sigs/aws-efs-csi-driver/deploy/kubernetes/overlays/stable/?ref=master"
}

variable "efs_pv_sifnoded_name" {
  description = "The name of the PV object"
  default     = "efs-sifnoded"
}

variable "efs_pv_storageclass" {
  description = "The name of the storageclass for the EFS driver"
  default     = "efs-sc"
}

variable "ebs_pv_storageclass" {
  description = "The name of the storageclass for the EBS driver"
  default     = "ebs-sc"
}

variable "efs_pv_capacity" {
  description = "EFS storage capacity"
  default     = "5Gi"
}

variable "profile" {
  description = "AWS profile settings"
  default     = "default"
}

variable "instance_type" {
  description = "The instance_type of the node_group for the eks cluster"
  type        = string
  default   = "m5.2xlarge"
}
