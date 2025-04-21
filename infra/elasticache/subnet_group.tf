// ========================================
// Redis Subnet Group
// ========================================

resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "${var.project}-redis-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project}-redis-subnet-group"
  }

  depends_on = [var.vpc_id]
}