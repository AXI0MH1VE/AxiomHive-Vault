terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
}

provider "google" {
  # Configure your credentials
  credentials = file("<path-to-service-account>")
  project     = var.project_id
  region      = var.region
}

# VPC Network
resource "google_compute_network" "ailock_vpc" {
  name                    = "ailock-vpc"
  auto_create_subnetworks = false
}

# Subnet
resource "google_compute_subnetwork" "ailock_subnet" {
  name          = "ailock-subnet"
  ip_cidr_range = "10.0.1.0/24"
  network       = google_compute_network.ailock_vpc.id
  region        = var.region
}

# GKE Cluster
resource "google_container_cluster" "ailock_cluster" {
  name     = "ailock-cluster"
  location = var.region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.ailock_vpc.id
  subnetwork = google_compute_subnetwork.ailock_subnet.id
}

# Node Pool
resource "google_container_node_pool" "ailock_nodes" {
  name       = "ailock-node-pool"
  location   = var.region
  cluster    = google_container_cluster.ailock_cluster.name
  node_count = 3

  node_config {
    machine_type = "e2-medium"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
