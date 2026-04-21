# -----------------------------------------------------------------------------
# VPC Module - Input Variables
# Sifnode GCP Network
# -----------------------------------------------------------------------------

variable "project_id" {
  description = "The GCP project ID where the VPC will be created."
  type        = string
}

variable "region" {
  description = "The GCP region for subnetwork creation."
  type        = string
}

variable "vpc_name" {
  description = "The name of the VPC network."
  type        = string
  default     = "sifnode-vpc"
}

variable "subnet_cidr" {
  description = "CIDR range for the primary subnetwork."
  type        = string
  default     = "10.0.0.0/20"
}

variable "pods_cidr" {
  description = "CIDR range for GKE pods (secondary IP range)."
  type        = string
  default     = "10.4.0.0/14"
}

variable "services_cidr" {
  description = "CIDR range for GKE services (secondary IP range)."
  type        = string
  default     = "10.8.0.0/16"
}