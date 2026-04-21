# -----------------------------------------------------------------------------
# GKE Module - Input Variables
# Sifnode Validator Cluster
# -----------------------------------------------------------------------------

variable "project_id" {
  description = "The GCP project ID where the GKE cluster will be created."
  type        = string
}

variable "region" {
  description = "The GCP region for the regional GKE cluster. Controls which zones are used for multi-zone HA."
  type        = string
}

variable "cluster_name" {
  description = "The name of the GKE cluster. Must be unique within the project."
  type        = string
}

variable "node_count_min" {
  description = "Minimum number of nodes in the validator node pool (autoscaling lower bound)."
  type        = number
  default     = 3
}

variable "node_count_max" {
  description = "Maximum number of nodes in the validator node pool (autoscaling upper bound)."
  type        = number
  default     = 5
}

variable "machine_type" {
  description = "The machine type for validator node pool instances. Default n2-standard-4 provides 4 vCPU and 16 GB RAM."
  type        = string
  default     = "n2-standard-4"
}

variable "network_vpc_id" {
  description = "The self-link or name of the VPC network into which the GKE cluster will be deployed."
  type        = string
}

variable "subnetwork_id" {
  description = "The self-link or name of the subnetwork within the VPC for the GKE cluster nodes."
  type        = string
}

variable "master_ipv4_cidr_block" {
  description = "The IPv4 CIDR block assigned to the GKE master internal endpoint. Must be /28 and within the VPC's secondary range."
  type        = string
  default     = "172.16.0.0/28"
}