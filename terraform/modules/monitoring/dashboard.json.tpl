{
  "displayName": "Sifnode Validator Overview",
  "mosaicLayout": {
    "columns": 12,
    "tiles": [
      {
        "width": 6,
        "height": 4,
        "widget": {
          "title": "CPU Utilization",
          "xyChart": {
            "dataSets": [
              {
                "timeSeriesFilter": {
                  "filter": "metric.type=\"kubernetes.io/container/cpu/core_usage_time\" resource.type=\"k8s_container\"%{if cluster_filter != ""} AND ${cluster_filter}%{endif}",
                  "plotType": "LINE"
                }
              }
            ],
            "timeshift": "0s",
            "yAxis": {
              "label": "CPU"
            }
          }
        }
      },
      {
        "width": 6,
        "height": 4,
        "widget": {
          "title": "Memory Utilization",
          "xyChart": {
            "dataSets": [
              {
                "timeSeriesFilter": {
                  "filter": "metric.type=\"kubernetes.io/container/memory/used_bytes\" resource.type=\"k8s_container\"%{if cluster_filter != ""} AND ${cluster_filter}%{endif}",
                  "plotType": "LINE"
                }
              }
            ],
            "timeshift": "0s",
            "yAxis": {
              "label": "Memory (bytes)"
            }
          }
        }
      },
      {
        "width": 6,
        "height": 4,
        "widget": {
          "title": "Disk Read / Write",
          "xyChart": {
            "dataSets": [
              {
                "timeSeriesFilter": {
                  "filter": "metric.type=\"agent.googleapis.com/disk/read_bytes_count\" resource.type=\"gce_instance\"%{if cluster_filter != ""} AND ${cluster_filter}%{endif}",
                  "plotType": "LINE"
                }
              },
              {
                "timeSeriesFilter": {
                  "filter": "metric.type=\"agent.googleapis.com/disk/write_bytes_count\" resource.type=\"gce_instance\"%{if cluster_filter != ""} AND ${cluster_filter}%{endif}",
                  "plotType": "LINE"
                }
              }
            ],
            "timeshift": "0s",
            "yAxis": {
              "label": "Bytes/sec"
            }
          }
        }
      },
      {
        "width": 6,
        "height": 4,
        "widget": {
          "title": "Received Bytes (Network)",
          "xyChart": {
            "dataSets": [
              {
                "timeSeriesFilter": {
                  "filter": "metric.type=\"kubernetes.io/pod/network/received_bytes_count\" resource.type=\"k8s_pod\"%{if cluster_filter != ""} AND ${cluster_filter}%{endif}",
                  "plotType": "LINE"
                }
              }
            ],
            "timeshift": "0s",
            "yAxis": {
              "label": "Bytes/sec"
            }
          }
        }
      },
      {
        "width": 12,
        "height": 4,
        "widget": {
          "title": "Uptime Check Status",
          "xyChart": {
            "dataSets": [
              {
                "timeSeriesFilter": {
                  "filter": "metric.type=\"monitoring.googleapis.com/uptime_check/check_passed\" resource.type=\"uptime_url\"",
                  "plotType": "LINE"
                }
              }
            ],
            "timeshift": "0s",
            "yAxis": {
              "label": "Status"
            }
          }
        }
      }
    ]
  }
}