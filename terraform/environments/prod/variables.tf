# Prod environment input variables
variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "notification_email" {
  description = "Email for alert notifications"
  type        = string
}

variable "lb_domain" {
  description = "Domain for SSL certificate"
  type        = string
  default     = ""
}

variable "alert_channels" {
  description = "Additional alert notification channels"
  type        = list(string)
  default     = []
}