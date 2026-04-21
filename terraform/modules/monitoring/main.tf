locals {
  cluster_filter = var.cluster_name != "" ? "resource.labels.cluster_name=\"${var.cluster_name}\"" : ""
}

# -----------------------------------------------------------------------------
# Notification Channel
# -----------------------------------------------------------------------------

resource "google_monitoring_notification_channel" "email" {
  display_name = "Sifnode Alert Email"
  type         = "email"

  labels = {
    email_address = var.notification_email
  }
}

# -----------------------------------------------------------------------------
# Alert Policy: Node Down (uptime check failure > 5 min)
# -----------------------------------------------------------------------------

resource "google_monitoring_alert_policy" "node_down" {
  display_name = "Sifnode Node Down"
  combiner     = "OR"

  conditions {
    display_name = "Uptime check failure for > 5 minutes"

    condition_threshold {
      filter          = "metric.type=\"monitoring.googleapis.com/uptime_check/check_passed\" resource.type=\"uptime_url\""
      duration        = "300s"
      comparison      = "COMPARISON_LT"
      threshold_value = 1

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_NEXT_OLDER"
      }
    }
  }

  notification_channels = concat(
    [google_monitoring_notification_channel.email.id],
    var.alert_channels,
  )

  severity = "CRITICAL"
}

# -----------------------------------------------------------------------------
# Alert Policy: Disk Usage > 80%
# -----------------------------------------------------------------------------

resource "google_monitoring_alert_policy" "disk_usage" {
  display_name = "Sifnode Disk Usage > 80%"
  combiner     = "OR"

  conditions {
    display_name = "Disk usage above 80 percent"

    condition_threshold {
      filter     = "metric.type=\"agent.googleapis.com/disk/percent_used\" resource.type=\"gce_instance\"${local.cluster_filter != "" ? " AND ${local.cluster_filter}" : ""}"
      duration   = "300s"
      comparison = "COMPARISON_GT"

      threshold_value = 80

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_MEAN"
      }
    }
  }

  notification_channels = concat(
    [google_monitoring_notification_channel.email.id],
    var.alert_channels,
  )

  severity = "WARNING"
}

# -----------------------------------------------------------------------------
# Alert Policy: High CPU > 90% for 5 minutes
# -----------------------------------------------------------------------------

resource "google_monitoring_alert_policy" "high_cpu" {
  display_name = "Sifnode High CPU > 90%"
  combiner     = "OR"

  conditions {
    display_name = "CPU utilization above 90 percent for 5 minutes"

    condition_threshold {
      filter     = "metric.type=\"kubernetes.io/container/cpu/core_usage_time\" resource.type=\"k8s_container\"${local.cluster_filter != "" ? " AND ${local.cluster_filter}" : ""}"
      duration   = "300s"
      comparison = "COMPARISON_GT"

      threshold_value = 0.9

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }

  notification_channels = concat(
    [google_monitoring_notification_channel.email.id],
    var.alert_channels,
  )

  severity = "WARNING"
}

# -----------------------------------------------------------------------------
# Alert Policy: Memory Usage > 85%
# -----------------------------------------------------------------------------

resource "google_monitoring_alert_policy" "memory_usage" {
  display_name = "Sifnode Memory Usage > 85%"
  combiner     = "OR"

  conditions {
    display_name = "Memory usage above 85 percent"

    condition_threshold {
      filter     = "metric.type=\"kubernetes.io/container/memory/used_bytes\" resource.type=\"k8s_container\"${local.cluster_filter != "" ? " AND ${local.cluster_filter}" : ""}"
      duration   = "300s"
      comparison = "COMPARISON_GT"

      threshold_value = 0.85

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_MEAN"
      }
    }
  }

  notification_channels = concat(
    [google_monitoring_notification_channel.email.id],
    var.alert_channels,
  )

  severity = "WARNING"
}

# -----------------------------------------------------------------------------
# Alert Policy: Sync Lag > 10 Blocks (placeholder — custom metric)
# -----------------------------------------------------------------------------

resource "google_monitoring_alert_policy" "sync_lag" {
  display_name = "Sifnode Sync Lag > 10 Blocks"
  combiner     = "OR"

  conditions {
    display_name = "Sync lag above 10 blocks"

    condition_threshold {
      filter     = "metric.type=\"custom.googleapis.com/sifnode/sync_lag\" resource.type=\"generic_node\""
      duration   = "300s"
      comparison = "COMPARISON_GT"

      threshold_value = 10

      aggregations {
        alignment_period   = "300s"
        per_series_aligner = "ALIGN_MEAN"
      }
    }
  }

  notification_channels = concat(
    [google_monitoring_notification_channel.email.id],
    var.alert_channels,
  )

  severity = "WARNING"
}

# -----------------------------------------------------------------------------
# Dashboard: Sifnode Validator Overview
# -----------------------------------------------------------------------------

resource "google_monitoring_dashboard" "validator_overview" {
  dashboard_json = templatefile(
    "${path.module}/dashboard.json.tpl",
    {
      cluster_filter = local.cluster_filter
    }
  )
}