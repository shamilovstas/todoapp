output "db_address" {
  value       = aws_db_instance.database.address
  description = "Database endpoint"
}

output "db_port" {
  value       = aws_db_instance.database.port
  description = "Database port"
}

output "db_pass" {
  value = data.aws_ssm_parameter.database_password.value
  sensitive = true
}