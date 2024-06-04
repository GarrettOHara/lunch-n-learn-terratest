terraform {
  backend "s3" {
    bucket = "sportsengine-dev-terraform-state"
    region = "us-east-1"
  }
}
