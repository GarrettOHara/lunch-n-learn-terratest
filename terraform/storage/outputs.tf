output "aws_s3_bucket_id" {
  value = try(aws_s3_bucket.this.id, "")
  description = "The S3 bucket ID."
  sensitive = false
}
