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