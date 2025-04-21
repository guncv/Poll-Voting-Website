module "ecs" {
  source         = "./ecs"                  
  project        = var.project             
  aws_account_id = var.aws_account_id
  aws_region     = var.aws_region
  vpc_id         = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
  redis_endpoint = module.elasticache.redis_endpoint
  redis_port = module.elasticache.redis_port
  private_route_table_id = module.vpc.private_route_table_id
  public_route_table_id = module.vpc.public_route_table_id
}

module "vpc" {
  source = "./vpc"
  project = var.project
  aws_region = var.aws_region
  availability_zones = var.availability_zones
}

module "elasticache" {
  source = "./elasticache"
  project = var.project
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  env = var.env
  availability_zones = var.availability_zones
  vpc_id = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
}

