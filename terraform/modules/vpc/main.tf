# -----------------------------------------------------------------------------
# VPC Module - Main Resources
# Sifnode GCP Network
# -----------------------------------------------------------------------------

resource "google_compute_network" "vpc" {
  name                    = var.vpc_name
  project                 = var.project_id
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"
}

resource "google_compute_subnetwork" "subnet" {
  name          = "${var.vpc_name}-subnet"
  project       = var.project_id
  region        = var.region
  network       = google_compute_network.vpc.id
  ip_cidr_range = var.subnet_cidr

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = var.pods_cidr
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = var.services_cidr
  }

  private_ip_google_access = true
}

# Firewall rules for Sifnode
resource "google_compute_firewall" "allow_internal" {
  name    = "${var.vpc_name}-allow-internal"
  project = var.project_id
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.subnet_cidr]
  description   = "Allow internal communication within the VPC"
}

resource "google_compute_firewall" "allow_sifnode_ports" {
  name    = "${var.vpc_name}-allow-sifnode"
  project = var.project_id
  network = google_compute_network.vpc.name

  # P2P port (Tendermint)
  allow {
    protocol = "tcp"
    ports    = ["26656"]
  }

  # RPC port
  allow {
    protocol = "tcp"
    ports    = ["26657"]
  }

  # REST API port
  allow {
    protocol = "tcp"
    ports    = ["1317"]
  }

  # Health check port
  allow {
    protocol = "tcp"
    ports    = ["1317"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["sifnode-validator"]
  description   = "Allow Sifnode P2P, RPC, and REST API ports"
}

resource "google_compute_firewall" "allow_health_checks" {
  name    = "${var.vpc_name}-allow-health-checks"
  project = var.project_id
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["1317", "26657"]
  }

  source_ranges = ["35.191.0.0/16", "130.211.0.0/22"]
  target_tags   = ["sifnode-validator"]
  description   = "Allow GCP health check probes"
}

# Cloud NAT for egress (private nodes need outbound access)
resource "google_compute_router" "router" {
  name    = "${var.vpc_name}-router"
  project = var.project_id
  region  = var.region
  network = google_compute_network.vpc.id
}

resource "google_compute_router_nat" "nat" {
  name                               = "${var.vpc_name}-nat"
  project                            = var.project_id
  region                             = var.region
  router                             = google_compute_router.router.name
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}