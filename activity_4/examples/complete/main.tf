module "s3_bucket" {
  source = "../../storage"
  name   = var.name
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
