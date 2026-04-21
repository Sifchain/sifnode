output "load_balancer_ip" {
  description = "External IP address of the global forwarding rule"
  value       = var.domain != "" ? google_compute_global_forwarding_rule.https[0].ip_address : google_compute_global_forwarding_rule.http[0].ip_address
}

output "forwarding_rule_name" {
  description = "Name of the global forwarding rule"
  value       = var.domain != "" ? google_compute_global_forwarding_rule.https[0].name : google_compute_global_forwarding_rule.http[0].name
}

output "ssl_certificate_name" {
  description = "Name of the managed SSL certificate (empty if no domain provided)"
  value       = var.domain != "" ? google_compute_managed_ssl_certificate.this[0].name : ""
}

output "url_map_name" {
  description = "Name of the URL map"
  value       = google_compute_url_map.this.name
}