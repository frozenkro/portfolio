output "instance_public_ip" {
  description = "Public IP of the EC2 instance (via EIP)"
  value       = aws_eip.web.public_ip
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.web.id
}

output "s3_bucket_name" {
  description = "Name of the S3 bucket holding the pre-built binary"
  value       = aws_s3_bucket.app_binary.id
}

output "s3_object_path" {
  description = "S3 path where the GHA workflow should upload the binary"
  value       = "s3://${aws_s3_bucket.app_binary.id}/${var.s3_object_key}"
}

output "route53_nameservers" {
  description = "Nameservers to point your registrar at (empty if no domain configured)"
  value       = var.domain_name != "" ? aws_route53_zone.primary[0].name_servers : []
}

output "site_url" {
  description = "URL to access the site (https://domain if configured, otherwise http://EIP)"
  value       = var.domain_name != "" ? "https://${var.domain_name}" : "http://${aws_eip.web.public_ip}"
}

output "gha_iam_access_key_id" {
  description = "Access key ID for the GHA S3 upload user — store as a GitHub secret (e.g. AWS_ACCESS_KEY_ID)"
  value       = aws_iam_access_key.gha.id
  sensitive   = true
}

output "gha_iam_secret_access_key" {
  description = "Secret access key for the GHA S3 upload user — store as a GitHub secret (e.g. AWS_SECRET_ACCESS_KEY)"
  value       = aws_iam_access_key.gha.secret
  sensitive   = true
}