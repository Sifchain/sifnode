output "notification_channel_id" {
  description = "ID of the created email notification channel"
  value       = google_monitoring_notification_channel.email.id
}

output "notification_channel_name" {
  description = "Full resource name of the email notification channel"
  value       = google_monitoring_notification_channel.email.name
}

output "alert_policy_ids" {
  description = "Map of alert policy IDs keyed by name"
  value = {
    node_down    = google_monitoring_alert_policy.node_down.id
    disk_usage   = google_monitoring_alert_policy.disk_usage.id
    high_cpu     = google_monitoring_alert_policy.high_cpu.id
    memory_usage = google_monitoring_alert_policy.memory_usage.id
    sync_lag     = google_monitoring_alert_policy.sync_lag.id
  }
}

output "dashboard_id" {
  description = "ID of the Sifnode Validator Overview dashboard"
  value       = google_monitoring_dashboard.validator_overview.id
}