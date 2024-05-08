resource "aws_instance" "this" {
  # checkov:skip=CKV_AWS_8: Encrypted volume not needed for testing
  ami                  = data.aws_ami.this.id
  ebs_optimized        = true
  iam_instance_profile = aws_iam_instance_profile.this.name
  instance_type        = var.instance_type
  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }
  monitoring = true
  root_block_device {
    encrypted = true
  }
  user_data                   = data.template_file.user_data_template.rendered
  user_data_replace_on_change = true
  vpc_security_group_ids = [
    aws_security_group.allow_software_updates.id,
    aws_security_group.allow_web_traffic.id
  ]

  tags = {
    Name = var.name
  }

  lifecycle {
    create_before_destroy = true
  }

  # Object file is required during user-data initiation
  depends_on = [
    aws_ssm_parameter.this
  ]
}

resource "aws_iam_role" "this" {
  name        = var.name
  description = "Role for bastion instance with Systems Manager access"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = "Ec2Access"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })

  inline_policy {
    name = "${var.name}-s3-access"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action = ["s3:*"]
          Effect = "Allow"
          Resource = [
            "arn:aws:s3:::${var.aws_bucket_id}",
            "arn:aws:s3:::${var.aws_bucket_id}/*"
          ]
        }
      ]
    })
  }

  inline_policy {
    name = "${var.name}-ssm-access"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action = [
            "ssm:*Parameter",
            "ssm:*Parameters"
          ]
          Effect = "Allow"
          Resource = [
            "arn:aws:ssm:*:*:parameter/AmazonCloudWatch-web-server-log-config"
          ]
        },
        {
          "Effect" : "Allow",
          "Action" : "ssm:DescribeParameters",
          "Resource" : "*"
        }
      ]
    })
  }

  inline_policy {
    name = "${var.name}-cloudwatch-logs-access"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action = [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents",
            "logs:DescribeLogStreams"
          ]
          Effect   = "Allow"
          Resource = ["*"]
        }
      ]

    })
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
  ]

  tags = {
    Name = var.name
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_iam_instance_profile" "this" {
  name = var.name
  role = aws_iam_role.this.name

  tags = {
    Name = var.name
  }
}

resource "aws_security_group" "allow_software_updates" {
  name        = var.name
  description = "Allow software updates"

  egress {
    description = "HTTPS for updates"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "DNS TCP for updates"
    from_port   = 53
    to_port     = 53
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "DNS UDP for updates"
    from_port   = 53
    to_port     = 53
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = var.name
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "allow_web_traffic" {
  name        = var.web_sg_name
  description = "Allow web traffic"

  ingress {
    description = "Allow HTTP web traffic"
    from_port   = var.ingress_port
    to_port     = var.ingress_port
    protocol    = var.protocol
    cidr_blocks = var.cidr_blocks
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = var.tags
}

resource "aws_ssm_parameter" "this" {
  name  = "AmazonCloudWatch-${var.name}-log-config"
  type  = "String"
  value = file("${path.module}/cloudwatch-logs-config.json")
}
