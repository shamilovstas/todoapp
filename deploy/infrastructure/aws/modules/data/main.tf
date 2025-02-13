data "aws_ssm_parameter" "database_password" {
  name            = "db_pass"
  with_decryption = true
}

resource "aws_db_subnet_group" "db_subnet_group" {
  name        = "database-subnet-group"
  description = "Group of subnets for RDS"
  subnet_ids  = var.db_subnet_ids
}

resource "aws_db_parameter_group" "db_param_group" {
  name   = "rds-postgres-pg"
  family = "postgres17"

  parameter {
    name  = "rds.force_ssl"
    value = 0
  }
}

resource "aws_security_group" "db_sg" {
  name        = "rds-sg"
  description = "Database security group"
  vpc_id      = var.vpc_id
}

resource "aws_db_instance" "database" {
  allocated_storage      = var.db_storage
  db_name                = var.db_name
  engine                 = "postgres"
  engine_version         = "17.2"
  instance_class         = var.db_instance_class
  skip_final_snapshot    = true
  multi_az               = false
  username               = var.db_username
  password               = data.aws_ssm_parameter.database_password.value
  db_subnet_group_name   = aws_db_subnet_group.db_subnet_group.name
  vpc_security_group_ids = [aws_security_group.db_sg.id]
  parameter_group_name   = aws_db_parameter_group.db_param_group.name
}

resource "aws_vpc_security_group_ingress_rule" "allow_postgres_in" {
  description       = "Allow Postgres inbound connection"
  from_port         = aws_db_instance.database.port
  to_port           = aws_db_instance.database.port
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  security_group_id = aws_security_group.db_sg.id
}
