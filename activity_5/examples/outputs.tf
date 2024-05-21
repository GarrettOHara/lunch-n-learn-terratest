output "s3_bucket_arn" {
  description = "The ARN of the S3 bucket."
  value       = try(aws_s3_bucket.private_bucket.arn, "")
}

