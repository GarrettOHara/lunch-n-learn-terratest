module "s3_bucket" {
  source              = "../../storage"
  name                = var.name
  region              = var.region
  static_website      = true
  static_website_CIDR = var.static_website_CIDR
}

output "s3_bucket_id" {
  description = "The ID of the S3 bucket."
  value       = try(module.s3_bucket.aws_s3_bucket_id, "")
}

output "website_endpoint" {
  value       = try(module.s3_bucket.website_endpoint, "")
  description = "The static website endpoint."
  sensitive   = false
}

variable "name" {
  type        = string
  description = "The name of the project."
  default     = "super-cool-bucket"
}

variable "region" {
  type        = string
  description = "The AWS Region"
  default     = "us-west-1"
}

variable "static_website_CIDR" {
  type        = list(string)
  description = "A list of valid CIDR ranges to serve content to."
  default = [
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
