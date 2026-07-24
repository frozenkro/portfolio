terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# ─── Data: default VPC + subnet ───────────────────────────────────────────────

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

data "aws_subnet" "default" {
  id = data.aws_subnets.default.ids[0]
}

# ─── AMI ──────────────────────────────────────────────────────────────────────

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

# ─── S3 bucket for pre-built binary ───────────────────────────────────────────

resource "aws_s3_bucket" "app_binary" {
  bucket_prefix = "${var.app_name}-binary-"
  force_destroy = true
}

resource "aws_s3_bucket_ownership_controls" "app_binary" {
  bucket = aws_s3_bucket.app_binary.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

# ─── IAM: instance role for S3 read access ────────────────────────────────────

data "aws_iam_policy_document" "instance_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "instance_role" {
  name               = "${var.app_name}-instance-role"
  assume_role_policy = data.aws_iam_policy_document.instance_assume_role.json
}

data "aws_iam_policy_document" "s3_read" {
  statement {
    actions = [
      "s3:GetObject",
      "s3:GetObjectVersion",
    ]
    resources = ["${aws_s3_bucket.app_binary.arn}/*"]
  }

  statement {
    actions   = ["s3:ListBucket"]
    resources = [aws_s3_bucket.app_binary.arn]
  }
}

resource "aws_iam_role_policy" "s3_read" {
  name   = "${var.app_name}-s3-read"
  role   = aws_iam_role.instance_role.id
  policy = data.aws_iam_policy_document.s3_read.json
}

resource "aws_iam_instance_profile" "instance_profile" {
  name = "${var.app_name}-instance-profile"
  role = aws_iam_role.instance_role.id
}

# ─── SSH key pair ─────────────────────────────────────────────────────────────

resource "aws_key_pair" "web" {
  key_name   = "${var.app_name}-ssh-key"
  public_key = file(var.ssh_public_key_path)
}

# ─── Security group ───────────────────────────────────────────────────────────

resource "aws_security_group" "web" {
  name        = "${var.app_name}-web-sg"
  description = "Allow SSH, HTTP, and HTTPS"
  vpc_id      = data.aws_vpc.default.id
}

resource "aws_vpc_security_group_ingress_rule" "ssh" {
  security_group_id = aws_security_group.web.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  to_port           = 22
  ip_protocol       = "tcp"
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.web.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
  ip_protocol       = "tcp"
}

resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.web.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
}

resource "aws_vpc_security_group_egress_rule" "all" {
  security_group_id = aws_security_group.web.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

# ─── Elastic IP ───────────────────────────────────────────────────────────────

resource "aws_eip" "web" {
  domain     = "vpc"
  depends_on = [aws_instance.web]
}

# ─── EC2 instance ─────────────────────────────────────────────────────────────

locals {
  user_data = templatefile("${path.module}/user_data.sh", {
    s3_bucket     = aws_s3_bucket.app_binary.id
    s3_object_key = var.s3_object_key
    app_name      = var.app_name
    app_port      = "8080"
  })
}

resource "aws_instance" "web" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  subnet_id              = data.aws_subnet.default.id
  iam_instance_profile   = aws_iam_instance_profile.instance_profile.name
  vpc_security_group_ids = [aws_security_group.web.id]
  key_name               = aws_key_pair.web.key_name

  user_data = local.user_data
  user_data_replace_on_change = true

  tags = {
    Name    = title(var.app_name)
    Project = var.app_name
  }
}

resource "aws_eip_association" "web" {
  instance_id   = aws_instance.web.id
  allocation_id = aws_eip.web.id
}

# ─── Route53 (optional — only if domain_name is set) ──────────────────────────

resource "aws_route53_zone" "primary" {
  count   = var.domain_name != "" ? 1 : 0
  name    = var.domain_name
  comment = "Hosted zone for ${var.app_name}"
}

resource "aws_route53_record" "root" {
  count   = var.domain_name != "" ? 1 : 0
  zone_id = aws_route53_zone.primary[0].zone_id
  name    = var.domain_name
  type    = "A"
  ttl     = 300
  records = [aws_eip.web.public_ip]
}

resource "aws_route53_record" "www" {
  count   = var.domain_name != "" ? 1 : 0
  zone_id = aws_route53_zone.primary[0].zone_id
  name    = "www.${var.domain_name}"
  type    = "A"
  ttl     = 300
  records = [aws_eip.web.public_ip]
}
