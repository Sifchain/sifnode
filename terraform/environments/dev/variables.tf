# Dev environment input variables
variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "notification_email" {
  description = "Email for alert notifications"
  type        = string
}