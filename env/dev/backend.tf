terraform {
  backend "s3" {
    bucket         = "sportsengine-dev-terraform-state"
    key            = "lunch-n-learn/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-lock"
  }
}
