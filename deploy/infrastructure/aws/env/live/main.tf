provider "aws" {
  region = "us-east-1"
}

locals {
  database_user = "postgres"
  database_name = "postgres"
}

module "vpc" {
  source             = "../../modules/networking"
  cidr_block         = "10.0.0.0/16"
  public_subnets     = ["10.0.1.0/24"]
  private_subnets    = ["10.0.2.0/24", "10.0.3.0/24"]
  availability_zones = ["us-east-1a", "us-east-1b"]
}

module "database" {
    source = "../../modules/data"
    vpc_id = module.vpc.vpc_id
    db_subnet_ids = module.vpc.private_subnets
    db_name = local.database_name
    db_username = local.database_user
}

module "instance" {
    source = "../../modules/services"
    vpc_id = module.vpc.vpc_id
    subnet_id = module.vpc.public_subnets[0]
    user_data = templatefile("../../../../ami/user_data.sh", {
        db_address = module.database.db_address
        db_port = module.database.db_port
        db_pass = module.database.db_pass
        db_user = local.database_user
        db_name = local.database_name
    })
}
