output "instance_id" {
  value       = aws_instance.this.id
  description = "The instance id"
  sensitive   = false
}

output "public_ipv4_addr" {
  value       = aws_instance.this.public_ip
  description = "The instance id"
  sensitive   = false
}

output "public_dns" {
  value       = aws_instance.this.public_dns
  description = "The IPv4 DNS domain"
  sensitive   = false
}
