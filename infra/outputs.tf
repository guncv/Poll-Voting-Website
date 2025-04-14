output "redis_endpoint" {
  description = "🔗 Redis Endpoint"
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
}

output "redis_port" {
  description = "🚪 Redis Port"
  value       = aws_elasticache_cluster.redis.cache_nodes[0].port
}

output "redis_sg_id" {
  description = "🛡 Security Group ID"
  value       = aws_security_group.redis_sg.id
}

output "redis_subnet_group" {
  description = "📦 Subnet Group Name"
  value       = aws_elasticache_subnet_group.redis_subnet_group.name
}

output "alb_dns_name" {
  description = "🌐 ALB DNS Name"
  value       = aws_lb.ecs_alb.dns_name
}

output "aws_account_id" {
  description = "AWS Account ID for ECR"
  value       = var.aws_account_id
}
