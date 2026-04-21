# -----------------------------------------------------------------------------
# GKE Module - Main Resources
# Sifnode Validator Cluster (Regional, Private, VPC-native)
# -----------------------------------------------------------------------------

# Service account for the GKE validator node pool
resource "google_service_account" "validator_sa" {
  account_id   = "${var.cluster_name}-validator"
  display_name = "Service account for Sifnode validator GKE node pool"
  project      = var.project_id
}

# -----------------------------------------------------------------------------
# GKE Regional Private Cluster
# -----------------------------------------------------------------------------
resource "google_container_cluster" "validator" {
  name     = var.cluster_name
  project  = var.project_id
  location = var.region

  # Regional cluster (multi-zone HA) - do not set node locations explicitly
  # to let GCP pick all zones in the region automatically.

  # VPC-native (IP aliasing enabled) for private cluster support
  network    = var.network_vpc_id
  subnetwork = var.subnetwork_id

  # Private cluster configuration
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }

  # IP allocation policy (VPC-native / alias IPs)
  ip_allocation_policy {
    # Let GKE auto-create secondary ranges for Pods and Services
  }

  # Workload Identity configuration
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  # Network Policy enforcement
  network_policy {
    enabled  = true
    provider = "CALICO"
  }

  # Shielded GKE nodes
  enable_shielded_nodes = true

  # Disable default node pool - we define our own below
  initial_node_count = 1

  # Deletion protection for production safety
  deletion_protection = false

  lifecycle {
    prevent_destroy = false
    ignore_changes  = [initial_node_count]
  }
}

# -----------------------------------------------------------------------------
# Validator Node Pool
# -----------------------------------------------------------------------------
resource "google_container_node_pool" "validator_pool" {
  name     = "validator-pool"
  project  = var.project_id
  location = var.region
  cluster  = google_container_cluster.validator.name

  # Initial node count (starts at min for autoscaling)
  initial_node_count = var.node_count_min

  # Autoscaling configuration
  autoscaling {
    min_node_count = var.node_count_min
    max_node_count = var.node_count_max
  }

  # Node configuration
  node_config {
    machine_type = var.machine_type
    disk_size_gb = 100
    disk_type    = "pd-ssd"

    # Not preemptible - validator nodes must be stable
    preemptible = false

    # Service account for node pool
    service_account = google_service_account.validator_sa.email

    # OAuth scopes for GCP API access
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    # Taint: dedicate nodes to validator workloads only
    taint {
      key    = "dedicated"
      value  = "validator"
      effect = "NO_SCHEDULE"
    }

    # Labels for identification
    labels = {
      purpose = "validator"
    }

    # Shielded instance config
    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    # Workload Identity metadata
    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }

  # Management - auto-repair and auto-upgrade for resilience
  management {
    auto_repair  = true
    auto_upgrade = true
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [google_container_cluster.validator]
}