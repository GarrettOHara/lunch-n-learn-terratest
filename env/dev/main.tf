module "tags" {
  source = "git::https://github.com/sportngin/tflib-tags.git?ref=v0.1.0"

  application       = var.name
  business_vertical = "sportsengine"
  name              = var.name
  env               = "dev"
  managed_by        = "terraform"
  repository        = "garrettohara/lunch-n-learn"
}

module "s3_storage" {
  source = "../../terraform/storage"
  name   = "lunch-n-learn"
}

module "garrett_server" {
  source        = "../../terraform/web-server"
  aws_bucket_id = s3_storage.aws_bucket_id
}
