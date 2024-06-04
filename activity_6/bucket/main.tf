resource "aws_s3_bucket" "this" {
  # checkov:skip=CKV_AWS_18: "Ensure the S3 bucket has access logging enabled"
  # checkov:skip=CKV_AWS_21: "Ensure all data stored in the S3 bucket have versioning enabled"
  # checkov:skip=CKV2_AWS_61: "Ensure that an S3 bucket has a lifecycle configuration"
  # checkov:skip=CKV2_AWS_62: "Ensure S3 buckets should have event notifications enabled"
  # checkov:skip=CKV_AWS_144: "Ensure that S3 bucket has cross-region replication enabled"
  # checkov:skip=CKV_AWS_145: "Ensure that S3 buckets are encrypted with KMS by default"
  # checkov:skip=CKV_AWS_186: No encryption needed for tests
  bucket        = "var.name fix me random_string.random.result"

  force_destroy = true
}

# Generate random string to allow distinct S3 bucket name
resource "random_string" "random" {
  length  = 12
  upper   = false
  special = false
}
