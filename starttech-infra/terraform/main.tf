terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    bucket = "starttech-terraform-state-uniqueid"
    key    = "prod/terraform.tfstate"
    region = "us-east-1"
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region
}

module "networking" {
  source = "./modules/networking"

  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
  environment        = var.environment
}

module "compute" {
  source = "./modules/compute"

  vpc_id             = module.networking.vpc_id
  public_subnet_ids  = module.networking.public_subnet_ids
  private_subnet_ids = module.networking.private_subnet_ids
  backend_image      = var.backend_image
  instance_type      = var.instance_type
  min_size           = var.min_size
  max_size           = var.max_size
  desired_capacity   = var.desired_capacity
  environment        = var.environment
  redis_endpoint     = module.storage.redis_endpoint
  mongodb_uri        = var.mongodb_uri
}

module "storage" {
  source = "./modules/storage"

  vpc_id            = module.networking.vpc_id
  private_subnet_ids = module.networking.private_subnet_ids
  environment       = var.environment
}

module "monitoring" {
  source = "./modules/monitoring"

  environment = var.environment
  alb_arn     = module.compute.alb_arn
  asg_name    = module.compute.asg_name
}
