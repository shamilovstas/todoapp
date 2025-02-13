variable "cidr_block" {
  type        = string
  description = "VPC CIDR"
}

variable "public_subnets" {
  type        = list(string)
  description = "Public subnets of the VPC"
}

variable "private_subnets" {
  type        = list(string)
  description = "private subnets of the VPC"
}

variable "availability_zones" {
  type        = list(string)
  description = "List of availability zones"
}
