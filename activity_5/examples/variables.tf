variable "name" {
  type        = string
  description = "The name of the project"
  default     = "super-cool-bucket"
}

variable "region" {
  type        = string
  description = "The AWS Region"
  default     = "us-west-1"
}
