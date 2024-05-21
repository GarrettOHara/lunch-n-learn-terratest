module "s3_bucket" {
  source         = "../../../terraform/storage/"
  name           = var.name
  static_website = true
}

variable "name" {
  type        = string
  description = "The name of the project."
  default     = "super-cool-bucket"
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
