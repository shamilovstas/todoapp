provider "aws" {
  region = "us-east-1"
}

# VPC
resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"

  tags = {
    Name = "main"
  }
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-1a"
  tags = {
    Name = "PublicSubnet"
  }
}

resource "aws_subnet" "private_subnet_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = false
  availability_zone       = "us-east-1a"
  tags = {
    Name = "PrivateSubnetA"
  }
}

resource "aws_subnet" "private_subnet_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.3.0/24"
  map_public_ip_on_launch = false
  availability_zone       = "us-east-1b"
  tags = {
    Name = "PrivateSubnetB"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "main"
  }
}

resource "aws_route_table" "private_rt" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "private_rt"
  }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
  tags = {
    Name = "public_rt"
  }
}

resource "aws_route_table_association" "private_a_rt_association" {
  subnet_id      = aws_subnet.private_subnet_a.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "private_b_rt_association" {
  subnet_id      = aws_subnet.private_subnet_b.id
  route_table_id = aws_route_table.private_rt.id
}

resource "aws_route_table_association" "public_rt_association" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_eip" "webserver_ip" {
  instance = aws_instance.backend-instance.id
  domain   = "vpc"
}
// Instance
data "aws_key_pair" "ssh-key" {
  key_name = "aws-keypair"
}

resource "aws_security_group" "instance-sg" {
  name        = "instance-sg"
  description = "Security group for the instance"
  vpc_id      = aws_vpc.main.id
  depends_on  = [aws_vpc.main]
}

resource "aws_vpc_security_group_ingress_rule" "allow_http_in" {
  description       = "Allow HTTP"
  from_port         = "80"
  to_port           = "80"
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  security_group_id = aws_security_group.instance-sg.id
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh_in" {
  description       = "Allow SSH"
  from_port         = "22"
  to_port           = "22"
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  security_group_id = aws_security_group.instance-sg.id
}

resource "aws_vpc_security_group_egress_rule" "allow_all_egress" {
  description       = "Allow all egress"
  cidr_ipv4         = "0.0.0.0/0"
  security_group_id = aws_security_group.instance-sg.id
  ip_protocol       = -1
}

data "aws_ami" "backend-ami" {
  most_recent = true
  owners      = ["self"]
  filter {
    name   = "name"
    values = ["ami-todoapp-backend-aws"]
  }
}

resource "aws_instance" "backend-instance" {
  ami           = data.aws_ami.backend-ami.id
  instance_type = "t3a.nano"
  user_data = templatefile("user_data.sh", {
    db_address = aws_db_instance.database.address,
    db_port    = aws_db_instance.database.port,
    db_pass    = data.aws_ssm_parameter.database_password.value,
    db_user    = var.db_username,
    db_name    = var.db_name
  })
  key_name        = data.aws_key_pair.ssh-key.key_name
  security_groups = [aws_security_group.instance-sg.id]
  subnet_id       = aws_subnet.public_subnet.id
  root_block_device {
    delete_on_termination = true
    volume_size           = 20
    volume_type           = "gp2"
  }
  depends_on = [aws_db_instance.database, aws_security_group.instance-sg]

  tags = {
    Name = "server"
    OS   = "Ubuntu"
  }
}

// Database

resource "aws_db_subnet_group" "private_db_subnet" {
  name        = "postgres-private-subnet-group"
  description = "Private subnet for database"
  subnet_ids  = [aws_subnet.private_subnet_a.id, aws_subnet.private_subnet_b.id]
}

data "aws_ssm_parameter" "database_password" {
  name            = "db_pass"
  with_decryption = true
}

resource "aws_db_parameter_group" "main_database_param_group" {
  name   = "rds-postgres-pg"
  family = "postgres17"

  parameter {
    name  = "rds.force_ssl"
    value = 0
  }

}
resource "aws_db_instance" "database" {
  allocated_storage      = 10
  db_name                = var.db_name
  engine                 = "postgres"
  engine_version         = "17.2"
  instance_class         = "db.t3.micro"
  skip_final_snapshot    = true
  multi_az               = false
  db_subnet_group_name   = aws_db_subnet_group.private_db_subnet.name
  vpc_security_group_ids = [aws_security_group.db-sg.id]
  parameter_group_name   = aws_db_parameter_group.main_database_param_group.name
  username               = var.db_username
  password               = data.aws_ssm_parameter.database_password.value
}

resource "aws_security_group" "db-sg" {
  name        = "rds-sg"
  description = "Database security group"
  vpc_id      = aws_vpc.main.id
  depends_on  = [aws_vpc.main]
}

resource "aws_vpc_security_group_ingress_rule" "allow_postgres_in" {
  description                  = "Allow Postgres inbound"
  from_port                    = aws_db_instance.database.port
  to_port                      = aws_db_instance.database.port
  ip_protocol                  = "tcp"
  security_group_id            = aws_security_group.db-sg.id
  referenced_security_group_id = aws_security_group.instance-sg.id
}

output "db_address" {
  value       = aws_db_instance.database.address
  description = "Database endpoint"
}

output "db_port" {
  value       = aws_db_instance.database.port
  description = "Database port"
}

variable "db_username" {
  description = "The username for the database"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
}
