variable "vpc_id" {
    type = string
    description = "VPC id"
}
variable "instance_type" {
    type = string
    description = "Backend instance type"
    default = "t3a.nano"
}

variable "user_data" {
    type = string
    description = "User data for instance"
}

variable "subnet_id" {
    type = string
    description = "Subnet id for the instance"
}

variable "volume_size" {
    type = number
    description = "Root block device size"
    default = 20
}

variable "volume_type" {
    type = string
    description = "Root block device type"
    default = "gp2"
}

variable "ssh_cidr_ipv4" {
    type = string
    description = "Allowed IP CIDR to use for SSH"
    default = "0.0.0.0/0"
}
