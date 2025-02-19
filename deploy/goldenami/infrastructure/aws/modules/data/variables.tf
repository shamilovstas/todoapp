variable "db_subnet_ids" {
  type        = list(string)
  description = "List of subnet ids to be used by RDS"
}

variable "db_storage" {
  type        = number
  description = "Storage size for RDS (GB)"
  default     = 20
}

variable "db_name" {
  type        = string
  description = "Database name"
}

variable "db_instance_class" {
  type        = string
  default     = "db.t3.micro"
  description = "Instance class for RDS"
}

variable "db_username" {
  type = string
  description = "Database master user name"
  
}

variable "vpc_id" {
  type        = string
  description = "VPC id for RDS"
}
