terraform {
  backend "gcs" {
    bucket  = "otrego-dev-infrastructure"
    prefix  = "terraform/dev/api"
  }
}

module "api_server_mod" {
    source = "../api_server_mod"
    dns_managed_zone = "otrego-dev"
    dns_name = "dev.otrego.com."
    project_id = "otrego-dev"
    subnetwork_project = "otrego-dev"
    api_docker_image = var.api_docker_image
    cos_image_name = "cos-stable-85-13310-1209-17"
}