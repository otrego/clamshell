terraform {
  backend "gcs" {
    bucket  = "otrego-prod-infrastructure"
    prefix  = "terraform/prod/api"
  }
}

module "api_server_mod" {
    source = "../api_server_mod"
    dns_managed_zone = "otrego-prod"
    dns_name = "www.otrego.com."
    project_id = "otrego-prod"
    subnetwork_project = "otrego-prod"
    api_docker_image = var.api_docker_image
    cos_image_name = "cos-stable-85-13310-1209-17"
}