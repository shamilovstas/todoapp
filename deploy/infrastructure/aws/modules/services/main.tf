data "aws_ami" "backend_ami" {
    most_recent = true
    owners = ["self"]
    filter {
        name = "name"
        values = ["ami-todoapp-backend-aws"]
    }
}

data "aws_key_pair" "ssh_key" {
    key_name = "aws-keypair"
}

resource "aws_security_group" "instance_sg" {
    name = "instance_sg"
    description = "Security group for the instance"
    vpc_id = var.vpc_id
}

resource "aws_vpc_security_group_ingress_rule" "allow_http_in" {
    description = "Allow HTTP in"
    from_port = "80"
    to_port = "80"
    ip_protocol = "tcp"
    cidr_ipv4 = "0.0.0.0/0"
    security_group_id = aws_security_group.instance_sg.id
}

resource "aws_vpc_security_group_ingress_rule" "allow_https_in" {
    description = "Allow HTTPS in"
    from_port = "443"
    to_port = "443"
    ip_protocol = "tcp"
    cidr_ipv4 = "0.0.0.0/0"
    security_group_id = aws_security_group.instance_sg.id
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh" {
    description = "Allow SSH in"
    from_port = "22"
    to_port = "22"
    ip_protocol = "tcp"
    cidr_ipv4 = var.ssh_cidr_ipv4
    security_group_id = aws_security_group.instance_sg.id
}

resource "aws_vpc_security_group_egress_rule" "allow_all_egress" {
    description = "Allow all egress"
    cidr_ipv4 = "0.0.0.0/0"
    security_group_id = aws_security_group.instance_sg.id
    ip_protocol = -1
}

resource "aws_eip" "webserver_ip" {
    instance = aws_instance.backend_instance.id
    domain = "vpc"
}

resource "aws_instance" "backend_instance" {
    ami = data.aws_ami.backend_ami.id
    instance_type = var.instance_type
    user_data = var.user_data
    key_name = data.aws_key_pair.ssh_key.key_name
    subnet_id = var.subnet_id
    vpc_security_group_ids = [aws_security_group.instance_sg.id]
    root_block_device {
      delete_on_termination = true
      volume_size = var.volume_size
      volume_type = var.volume_type
    }
    
}