output "instance_name" {
  description = "Name of the Cloud SQL instance"
  value       = google_sql_database_instance.instance.name
}

output "instance_connection_name" {
  description = "Connection name of the Cloud SQL instance (project:region:instance)"
  value       = google_sql_database_instance.instance.connection_name
}

output "private_ip" {
  description = "Private IP address of the Cloud SQL instance"
  value       = google_sql_database_instance.instance.private_ip_address
}

output "database_name" {
  description = "Name of the created database"
  value       = google_sql_database.database.name
}

output "server_ca_cert" {
  description = "Server CA certificate of the Cloud SQL instance"
  value       = google_sql_database_instance.instance.server_ca_cert
}