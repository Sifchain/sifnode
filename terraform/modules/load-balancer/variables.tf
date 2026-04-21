variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "name" {
  description = "Name prefix for load balancer resources"
  type        = string
}

variable "region" {
  description = "GCP region for the load balancer resources"
  type        = string
}

variable "domain" {
  description = "Domain for managed SSL certificate. If empty, an HTTP proxy is created instead of HTTPS."
  type        = string
  default     = ""
}

variable "backend_service" {
  description = "Self-link of the backend service or NEG to attach to the load balancer"
  type        = string
  default     = ""
}

variable "health_check_path" {
  description = "HTTP path for health check probes"
  type        = string
  default     = "/status"
}

variable "health_check_port" {
  description = "Port for health check probes"
  type        = number
  default     = 1317
}