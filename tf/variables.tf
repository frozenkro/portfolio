variable "aws_region" {
  description = "AWS region for all resources"
  type        = string
  default     = "us-east-1"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "domain_name" {
  description = "Domain name for Route53 zone. Leave empty to skip DNS resources and use EIP directly."
  type        = string
  default     = ""
}

variable "app_name" {
  description = "Name used for resource tags and S3 bucket naming"
  type        = string
  default     = "portfolio"
}

variable "s3_object_key" {
  description = "S3 object key for the pre-built Go binary"
  type        = string
  default     = "portfolio"
}

variable "ssh_public_key_path" {
  description = "Path to the local SSH public key to register in AWS"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}