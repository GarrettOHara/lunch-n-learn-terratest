module "web_server" {
  source = "git::https://github.com/garrettohara/lunch-n-learn-terratest.git//terraform/web_server?ref=main" # "../../../terraform/web_server/"
  name   = var.name

  cidr_blocks = [
    "165.1.165.11/32",   # My VPN IP
    "50.232.111.124/32", # SE HQ Corp employee wired and wireless
    "67.130.26.171/32",
    "137.83.201.101/32", # SE Global Protect Cloud Service (GPCS) US East
    "137.83.201.167/32",
    "54.67.50.109/32", # SE Global Protect Cloud Service (GPCS) US West
    "13.52.120.179/32",
    "50.228.144.140/32", # NBC Global Protect Contractor VPN US East
    "50.228.144.124/32",
    "50.230.144.156/32", # NBC Global Protect Contractor VPN US West
  ]
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
