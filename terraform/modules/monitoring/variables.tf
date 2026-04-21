variable "project_id" {
  description = "GCP project ID where monitoring resources will be created"
  type        = string
}

variable "notification_email" {
  description = "Email address for alert notifications"
  type        = string
}

variable "alert_channels" {
  description = "Additional notification channel IDs to attach to alert policies"
  type        = list(string)
  default     = []
}

variable "cluster_name" {
  description = "GKE cluster name used for metric filtering (optional)"
  type        = string
  default     = ""
}