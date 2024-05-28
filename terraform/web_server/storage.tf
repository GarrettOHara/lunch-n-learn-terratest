resource "aws_s3_bucket" "this" {
  # checkov:skip=CKV_AWS_18: "Ensure the S3 bucket has access logging enabled"
  # checkov:skip=CKV_AWS_21: "Ensure all data stored in the S3 bucket have versioning enabled"
  # checkov:skip=CKV2_AWS_61: "Ensure that an S3 bucket has a lifecycle configuration"
  # checkov:skip=CKV2_AWS_62: "Ensure S3 buckets should have event notifications enabled"
  # checkov:skip=CKV_AWS_144: "Ensure that S3 bucket has cross-region replication enabled"
  # checkov:skip=CKV_AWS_145: "Ensure that S3 buckets are encrypted with KMS by default"
  # checkov:skip=CKV_AWS_186: No encryption needed for tests
  bucket        = "${var.name}-${random_string.random.result}"
  force_destroy = true
  tags          = var.tags
}

# Generate random string to allow distinct S3 bucket name
resource "random_string" "random" {
  length  = 12
  upper   = false
  special = false
}

resource "aws_s3_bucket_ownership_controls" "this" {
  bucket = aws_s3_bucket.this.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

# Block all public access to the S3 bucket
resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_object" "app_source_code" {
  # checkov:skip=CKV_AWS_186: No encrpytion required for non sensitive files
  for_each    = toset([for file in fileset("${path.module}/../../src/server", "**") : file])
  bucket      = aws_s3_bucket.this.id
  key         = each.value
  source      = "${path.module}/../../src/server/${each.value}"
  source_hash = filemd5("${path.module}/../../src/server/${each.value}")
}

# resource "aws_s3_object" "app_executable" {
#   bucket      = aws_s3_bucket.this.id
#   key         = "server"
#   source      = "${path.module}/../../src/server/server"
#   source_hash = filemd5("${path.module}/../../src/server/server")
# }

resource "aws_s3_object" "app_html_templates" {
  # checkov:skip=CKV_AWS_186: No encrpytion required for non sensitive files
  for_each = toset([for file in fileset("${path.module}/../../src/server/templates", "**") : file])

  bucket      = aws_s3_bucket.this.id
  key         = "templates/${each.value}"
  source      = "${path.module}/../../src/server/templates/${each.value}"
  source_hash = filemd5("${path.module}/../../src/server/templates/${each.value}")
}
