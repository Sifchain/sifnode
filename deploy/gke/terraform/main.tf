resource "google_compute_network" "sifnode" {
  name                    = "${var.cluster_name}-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "sifnode" {
  name          = "${var.cluster_name}-subnet"
  region        = var.region
  network       = google_compute_network.sifnode.id
  ip_cidr_range = var.network_cidr

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = var.pods_cidr
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = var.services_cidr
  }
}

resource "google_compute_router" "sifnode" {
  name    = "${var.cluster_name}-router"
  region  = var.region
  network = google_compute_network.sifnode.id
}

resource "google_compute_router_nat" "sifnode" {
  name                               = "${var.cluster_name}-nat"
  router                             = google_compute_router.sifnode.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.sifnode.id
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

resource "google_container_cluster" "sifnode" {
  name     = var.cluster_name
  location = var.region

  network    = google_compute_network.sifnode.name
  subnetwork = google_compute_subnetwork.sifnode.name

  remove_default_node_pool = true
  initial_node_count       = 1

  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  release_channel {
    channel = "REGULAR"
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  addons_config {
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
  }

  resource_labels = var.labels
}

resource "google_container_node_pool" "sifnode" {
  name       = "${var.cluster_name}-pool"
  location   = var.region
  cluster    = google_container_cluster.sifnode.name
  node_count = var.node_count

  node_config {
    machine_type = var.node_machine_type
    disk_size_gb = var.disk_size_gb
    disk_type    = "pd-balanced"

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    labels = var.labels
    tags   = ["${var.cluster_name}-node"]
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}
