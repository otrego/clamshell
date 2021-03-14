variable "dns_managed_zone" {
  description = "The dns managed zone in GCP Cloud DNS. This is where DNS entries should be writable. Example: otrego-dev"
}
variable "dns_name" {
  description = "A record entry within the DNS Manged Zone. Eg: dev.otrego.com."
}

variable "project_id" {
  description = "The project ID to deploy resources into"
  default = "otrego-dev"
}

variable "subnetwork_project" {
  description = "The project ID where the desired subnetwork is provisioned"
  default     = "otrego-dev"
}

variable "subnetwork" {
  description = "The name of the subnetwork to deploy instances into"
  default     = "default"
}

variable "zone" {
  description = "The GCP zone to deploy instances into"
  type        = string
  default     = "us-west1-a"
}

variable "client_email" {
  description = "Service account email address"
  type        = string
  default     = ""
}

variable "cos_image_name" {
  description = "The container optimized image. This is the base image for the machine on which docker is run. Should be updated regularly. Example: cos-stable-77-12371-89-0"
}

variable "api_docker_image" {
  description = "Docker Image of the Otrego API"
  type        = string
}
