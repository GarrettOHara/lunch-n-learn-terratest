variable "cidr_blocks" {
  type        = list(string)
  description = "The CIDR range of the ingress web traffic"
  default     = ["0.0.0.0/0"]
}

variable "ingress_port" {
  type        = number
  description = "The HTTP port"
  default     = 80
}

variable "instance_type" {
  type        = string
  description = "The EC2 instance type"
  default     = "t3.micro"
}

variable "name" {
  type        = string
  description = "The name of the project"
  default     = "web-server"
}

variable "protocol" {
  type        = string
  description = "The layer 4 protocol"
  default     = "tcp"
}

variable "region" {
  type        = string
  description = "The AWS Region"
  default     = "us-west-1"
}

variable "tags" {
  type        = map(string)
  description = "Resource tags"
  default     = {}
}

