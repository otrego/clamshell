terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.59.0"
    }
  }
  required_version = ">= 0.14, < 0.15, < 1.0"
}
