output "vm_container_label" {
  description = "The instance label containing container configuration"
  value       = module.api_server_mod.vm_container_label
}

output "container" {
  description = "The container metadata provided to the module"
  value       = module.api_server_mod.container
}

output "volumes" {
  description = "The volume metadata provided to the module"
  value       = module.api_server_mod.volumes
}

output "instance_name" {
  description = "The deployed instance name"
  value       = module.api_server_mod.instance_name
}

output "ipv4" {
  description = "The public IP address of the deployed instance"
  value       = module.api_server_mod.ipv4
}


