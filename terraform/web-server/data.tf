data "aws_ami" "this" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["amazon"]
}

data "template_file" "user_data_template" {
  template = file("${path.module}/user-data.sh")

  vars = {
    s3_bucket     = aws_s3_bucket.this.id
    ssm_parameter = aws_ssm_parameter.this.name
  }
}
