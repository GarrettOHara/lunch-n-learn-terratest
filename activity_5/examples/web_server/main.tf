module "web_server" {
  source      = "../../../terraform/web_server/"
  cidr_blocks = ["10.57.0.0/16"]
  name        = var.name
}

variable "name" {
  type        = string
  description = "The name of the project."
  default     = "super-cool-web-server"
}

output "instance_id" {
  value       = module.web_server.instance_id
  description = "The instance id"
  sensitive   = false
}

output "public_ipv4_addr" {
  value       = module.web_server.public_ipv4_addr
  description = "The instance id"
  sensitive   = false
}

output "public_dns" {
  value       = module.web_server.public_dns
  description = "The IPv4 DNS domain"
  sensitive   = false
}
