terraform {
  required_version = ">= 0.13"
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  token                  = data.google_client_config.current.access_token
}

data "google_client_config" "current" {}

resource "google_service_account" "sifnode" {
  account_id   = "${var.cluster_name}-sa"
  display_name = "sifnode GKE Service Account"
  description  = "Service account for sifnode GKE cluster operations"
}

resource "google_project_iam_member" "sifnode_permissions" {
  for_each = toset([
    "roles/container.admin",
    "roles/compute.admin",
    "roles/iam.serviceAccountUser",
    "roles/storage.admin"
  ])
  role    = each.key
  member  = "serviceAccount:${google_service_account.sifnode.email}"
}

resource "google_container_cluster" "primary" {
  provider = google-beta
  name     = var.cluster_name
  location = var.region

  remove_default_node_pool = true
  initial_node_count       = var.initial_node_count

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  private_cluster_config {
    enable_private_nodes    = var.enable_private_nodes
    enable_private_endpoint = false
    master_ipv4_cidr_block  = var.master_ipv4_cidr
  }

  master_authorized_networks_config {
    dynamic "cidr_blocks" {
      for_each = var.authorized_cidrs
      content {
        cidr_block   = cidr_blocks.value
        display_name = "authorized-${cidr_blocks.key}"
      }
    }
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  cluster_autoscaling {
    enabled = var.enable_autoscaling
    resource_limits {
      resource_type = "cpu"
      minimum       = var.min_cpu
      maximum       = var.max_cpu
    }
    resource_limits {
      resource_type = "memory"
      minimum       = var.min_memory_mb
      maximum       = var.max_memory_mb
    }
  }

  node_pool_defaults {
    node_config_defaults {
      image_type = var.image_type
    }
  }

  release_channel {
    channel = var.release_channel
  }

  maintenance_policy {
    daily_maintenance_window {
      start_time = var.maintenance_window
    }
  }

  labels = var.tags

  lifecycle {
    ignore_changes = [
      node_pool,
      initial_node_count
    ]
  }
}

resource "google_compute_network" "vpc" {
  name                    = "${var.cluster_name}-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "${var.cluster_name}-subnet"
  ip_cidr_range = var.subnet_cidr
  region        = var.region
  network       = google_compute_network.vpc.id

  secondary_ip_range = [
    {
      range_name    = "pods"
      ip_cidr_range = var.pods_cidr
    },
    {
      range_name    = "services"
      ip_cidr_range = var.services_cidr
    }
  ]
}

resource "google_compute_router" "router" {
  name    = "${var.cluster_name}-router"
  region  = var.region
  network = google_compute_network.vpc.id
}

resource "google_compute_router_nat" "nat" {
  name                               = "${var.cluster_name}-nat"
  router                             = google_compute_router.router.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.cluster_name}-node-pool"
  cluster    = google_container_cluster.primary.id
  location   = var.region
  node_count = var.desired_capacity

  autoscaling {
    min_node_count = var.min_capacity
    max_node_count = var.max_capacity
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    machine_type = var.machine_type
    disk_size_gb = var.disk_size
    disk_type    = var.disk_type
    image_type   = var.image_type

    service_account = google_service_account.sifnode.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring.write"
    ]

    labels = merge(var.tags, {
      role = "sifnode"
    })

    tags = ["sifnode", var.cluster_name]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }

  lifecycle {
    ignore_changes = [node_count]
  }
}

resource "google_compute_firewall" "sifnode" {
  name    = "${var.cluster_name}-sifnode"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["26656", "26657", "1317", "9090"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["sifnode", var.cluster_name]
}

resource "google_compute_firewall" "ssh" {
  name    = "${var.cluster_name}-ssh"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = var.ssh_source_ranges
  target_tags   = ["sifnode"]
}

resource "google_compute_disk" "sifnoded_data" {
  name  = "${var.cluster_name}-sifnoded-data"
  type  = "pd-ssd"
  zone  = var.default_zone
  size  = var.data_disk_size
  labels = var.tags
}

resource "google_compute_disk" "sifnodecli_data" {
  name  = "${var.cluster_name}-sifnodecli-data"
  type  = "pd-ssd"
  zone  = var.default_zone
  size  = var.data_disk_size
  labels = var.tags
}

resource "kubernetes_storage_class" "ssd" {
  metadata {
    name = "gke-ssd"
  }
  storage_provisioner = "kubernetes.io/gce-pd"
  parameters = {
    type = "pd-ssd"
  }
  allow_volume_expansion = true
  reclaim_policy        = "Retain"
}

resource "kubernetes_persistent_volume_claim" "sifnoded" {
  metadata {
    name      = "sifnoded-pvc"
    namespace = var.namespace
    labels    = var.tags
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.storage_size
      }
    }
    storage_class_name = kubernetes_storage_class.ssd.metadata[0].name
  }
}

resource "kubernetes_persistent_volume_claim" "sifnodecli" {
  metadata {
    name      = "sifnodecli-pvc"
    namespace = var.namespace
    labels    = var.tags
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = var.storage_size
      }
    }
    storage_class_name = kubernetes_storage_class.ssd.metadata[0].name
  }
}
