output "aws_s3_bucket_id" {
  value       = try(aws_s3_bucket.this.id, "")
  description = "The S3 bucket ID."
  sensitive   = false
}

output "website_endpoint" {
  value       = try(aws_s3_bucket_website_configuration.this.website_endpoint, "")
  description = "The static website endpoint."
  sensitive   = false
}
