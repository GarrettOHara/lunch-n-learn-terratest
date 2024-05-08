variable "aws_region" {
  description = "The AWS region we are deploying to."
  type        = string
  default     = "us-west-1"
}

variable "name" {
  description = "The base name for these resources."
  type        = string
  default     = "lunch-n-learn"
}
