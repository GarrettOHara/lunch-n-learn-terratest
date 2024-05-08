variable "name" {
  type        = string
  description = "The name of the project."
}

variable "region" {
  type        = string
  description = "The AWS Region"
  default     = "us-west-1"
}

variable "tags" {
  type        = map(string)
  description = "The tags for the module's resources."
  default     = {}

}
