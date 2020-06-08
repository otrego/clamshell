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

variable "instance_name" {
  description = "The desired name to assign to the deployed instance"
  default     = "hello-world-container-vm"
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
  description = "The forced COS image to use instead of latest"
  default     = "cos-stable-77-12371-89-0"
}
