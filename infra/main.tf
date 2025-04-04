// ========================================
// VPC Setup
// ========================================

resource "aws_vpc" "cv_c9_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "${var.project}-vpc"
  }
}

// ========================================
// Private Subnets
// ========================================

resource "aws_subnet" "private_1" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = var.availability_zone
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.project}-private-subnet-1"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}

resource "aws_subnet" "private_2" {
  vpc_id                  = aws_vpc.cv_c9_vpc.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = var.availability_zone
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.project}-private-subnet-2"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}

// ========================================
// Security Group for Redis
// ========================================

resource "aws_security_group" "redis_sg" {
  name        = "${var.project}-redis-sg"
  description = "Allow Redis traffic"
  vpc_id      = aws_vpc.cv_c9_vpc.id

  # 🧠 Conditionally allow Redis access based on environment
  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = var.env == "dev" ? ["0.0.0.0/0"] : ["10.0.0.0/16"]
    description = var.env == "dev" ? "Allow from anywhere (dev)" : "Internal VPC access only (prod)"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-redis-sg"
  }

  depends_on = [aws_vpc.cv_c9_vpc]
}


// ========================================
// Redis Subnet Group
// ========================================

resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "${var.project}-redis-subnet-group"
  subnet_ids = [aws_subnet.private_1.id, aws_subnet.private_2.id]

  tags = {
    Name = "${var.project}-redis-subnet-group"
  }

  depends_on = [aws_subnet.private_1, aws_subnet.private_2]
}

// ========================================
// ElastiCache Redis Cluster
// ========================================

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${var.project}-redis"
  engine               = "redis"
  node_type            = "cache.t2.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379

  subnet_group_name   = aws_elasticache_subnet_group.redis_subnet_group.name
  security_group_ids  = [aws_security_group.redis_sg.id]

  tags = {
    Name = "${var.project}-redis"
  }

  depends_on = [
    aws_elasticache_subnet_group.redis_subnet_group,
    aws_security_group.redis_sg
  ]
}
