module "s3_bucket" {
    name = var.name
}

variable "name" {
  type        = string
  description = "The name of the project"
  default     = "super-cool-bucket"
}

output "s3_bucket_arn" {
  description = "The ARN of the S3 bucket."
  value       = try(module.s3_bucket.s3_bucket_arn, "")
}
