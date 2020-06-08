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
    image="gcr.io/google-samples/hello-app:2.0"
    env = [
      {
        name = "TEST_VAR"
        value = "Hello World!"
      }
    ],

    # Declare volumes to be mounted.
    # This is similar to how docker volumes are declared.
    volumeMounts = [
    #   {
    #     mountPath = "/cache"
    #     name      = "tempfs-0"
    #     readOnly  = false
    #   },
    #   {
    #     mountPath = "/persistent-data"
    #     name      = "data-disk-0"
    #     readOnly  = false
    #   },
    ]
  }

  # Declare the Volumes which will be used for mounting.
  volumes = [
    # {
    #   name = "tempfs-0"
    #
    #   emptyDir = {
    #     medium = "Memory"
    #   }
    # },
    # {
    #   name = "data-disk-0"
    #
    #   gcePersistentDisk = {
    #     pdName = "data-disk-0"
    #     fsType = "ext4"
    #   }
    # },
  ]

  restart_policy = "Always"
}

resource "google_compute_instance" "vm" {
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
    google-logging-enabled    = "true"
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
#
# resource "google_compute_network" "default" {
#   name = "default"
# }
