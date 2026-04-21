# ---------------------------------------------------------------------------
# Health Check – probes /status on port 1317 (default) via HTTP
# ---------------------------------------------------------------------------
resource "google_compute_health_check" "this" {
  name        = "${var.name}-hc"
  project     = var.project_id
  description = "Health check for ${var.name} load balancer"

  http_health_check {
    port         = var.health_check_port
    request_path = var.health_check_path
  }

  check_interval_sec  = 10
  timeout_sec         = 5
  healthy_threshold   = 2
  unhealthy_threshold = 3
}

# ---------------------------------------------------------------------------
# Backend Service – no CDN (real-time blockchain data), 30s timeout
# ---------------------------------------------------------------------------
resource "google_compute_backend_service" "this" {
  name    = "${var.name}-backend"
  project = var.project_id

  protocol    = "HTTP"
  port_name   = "http"
  timeout_sec = 30

  health_checks = [google_compute_health_check.this.id]

  # CDN is disabled by default; explicitly set to false
  enable_cdn = false

  # If the caller passed an explicit backend_service reference we do not
  # manage backends here – they are expected to be NEGs attached externally.
  # When a dedicated backend_service var is NOT supplied we leave backends
  # empty so they can be referenced / attached later.
  dynamic "backend" {
    for_each = var.backend_service != "" ? [var.backend_service] : []
    content {
      group = backend.value
    }
  }
}

# ---------------------------------------------------------------------------
# URL Map – routes all traffic to the backend service
# ---------------------------------------------------------------------------
resource "google_compute_url_map" "this" {
  name            = "${var.name}-url-map"
  project         = var.project_id
  default_service = google_compute_backend_service.this.id

  description = "URL map for ${var.name}"
}

# ---------------------------------------------------------------------------
# Managed SSL Certificate (conditional – only when domain is provided)
# ---------------------------------------------------------------------------
resource "google_compute_managed_ssl_certificate" "this" {
  count   = var.domain != "" ? 1 : 0
  name    = "${var.name}-ssl-cert"
  project = var.project_id

  managed {
    domains = [var.domain]
  }
}

# ---------------------------------------------------------------------------
# Target HTTPS Proxy (when domain/SSL cert exists)
# ---------------------------------------------------------------------------
resource "google_compute_target_https_proxy" "this" {
  count            = var.domain != "" ? 1 : 0
  name             = "${var.name}-https-proxy"
  project          = var.project_id
  url_map          = google_compute_url_map.this.id
  ssl_certificates = [google_compute_managed_ssl_certificate.this[0].id]
}

# ---------------------------------------------------------------------------
# Target HTTP Proxy (when NO domain – fallback)
# ---------------------------------------------------------------------------
resource "google_compute_target_http_proxy" "this" {
  count   = var.domain == "" ? 1 : 0
  name    = "${var.name}-http-proxy"
  project = var.project_id
  url_map = google_compute_url_map.this.id
}

# ---------------------------------------------------------------------------
# Global Forwarding Rule – HTTPS (when domain provided)
# ---------------------------------------------------------------------------
resource "google_compute_global_forwarding_rule" "https" {
  count       = var.domain != "" ? 1 : 0
  name        = "${var.name}-https-fwd-rule"
  project     = var.project_id
  target      = google_compute_target_https_proxy.this[0].id
  port_range  = "443"
  ip_protocol = "TCP"

  load_balancing_scheme = "EXTERNAL_MANAGED"
}

# ---------------------------------------------------------------------------
# Global Forwarding Rule – HTTP (when no domain)
# ---------------------------------------------------------------------------
resource "google_compute_global_forwarding_rule" "http" {
  count       = var.domain == "" ? 1 : 0
  name        = "${var.name}-http-fwd-rule"
  project     = var.project_id
  target      = google_compute_target_http_proxy.this[0].id
  port_range  = "80"
  ip_protocol = "TCP"

  load_balancing_scheme = "EXTERNAL_MANAGED"
}