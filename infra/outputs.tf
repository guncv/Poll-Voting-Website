output "redis_endpoint" {
  description = "🔗 Redis Endpoint"
  value       = module.elasticache.redis_endpoint
}

output "redis_port" {
  description = "🚪 Redis Port"
  value       = module.elasticache.redis_port
}

output "redis_sg_id" {
  description = "🛡 Security Group ID"
  value       = module.elasticache.redis_security_group_id
}

output "redis_subnet_group" {
  description = "📦 Subnet Group Name"
  value       = module.elasticache.redis_subnet_group_name
}

output "alb_dns_name" {
  description = "🌐 ALB DNS Name"
  value       = module.ecs.alb_dns_name
}

data "aws_caller_identity" "current" {}

output "aws_account_id" {
  value = coalesce(var.aws_account_id, data.aws_caller_identity.current.account_id)
}
