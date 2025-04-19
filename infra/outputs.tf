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

data "aws_caller_identity" "current" {}

output "aws_account_id" {
  value = coalesce(var.aws_account_id, data.aws_caller_identity.current.account_id)
}
