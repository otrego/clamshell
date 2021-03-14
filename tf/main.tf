terraform {
  backend "gcs" {
    bucket  = "otrego-dev-infrastructure"
    prefix  = "terraform/dev/api"
  }
}

provider "google" {
}

locals {
  instance_name = format("%s-%s", var.instance_name, substr(md5(module.gce-container.container.image), 0, 8))
}

module "gce-container" {
  source = "terraform-google-modules/container-vm/google"

  container = {
    image = var.api_docker_image
    env = [
      {
        name = "TEST_VAR"
        value = "Hello World!"
      }
    ],

    # Declare volumes to be mounted.
    # This is similar to how docker volumes are declared.
    volumeMounts = [
    ]
  }

  # Declare the Volumes which will be used for mounting.
  volumes = [
  ]

  restart_policy = "Always"
}

resource "google_compute_instance" "otrego_instance" {
  project      = var.project_id
  name         = local.instance_name
  machine_type = "f1-micro"
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = module.gce-container.source_image
    }
  }

  network_interface {
    network = "default"
    # subnetwork = var.subnetwork
    # subnetwork_project = var.subnetwork_project
    access_config {}
  }

  tags = ["web", "http-server"]

  metadata = {
    gce-container-declaration = module.gce-container.metadata_value
    google-logging-enabled    = "false"
    google-monitoring-enabled = "false"
  }

  labels = {
    container-vm = module.gce-container.vm_container_label
  }

  service_account {
    email = var.client_email
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}

resource "google_compute_firewall" "default" {
  name    = "otrego-dev-firewall"
  network = "default"
  # network = var.subnetwork
  project = var.subnetwork_project

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }

  direction = "INGRESS"
  source_ranges = ["0.0.0.0/0"]

  source_tags = ["web"]
}

resource "google_dns_record_set" "frontend" {
  name = "dev.otrego.com."
  type = "A"
  ttl  = 300

  managed_zone = var.dns_managed_zone
  project = var.project_id

  rrdatas = [google_compute_instance.otrego_instance.network_interface[0].access_config[0].nat_ip]
}
