output "redis_cluster_id" {
  value = aws_elasticache_cluster.redis.id
}

output "redis_security_group_id" {
  value = aws_security_group.redis_sg.id
}

output "redis_subnet_group_name" {
  value = aws_elasticache_subnet_group.redis_subnet_group.name
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.redis.cache_nodes[0].address
}

output "redis_port" {
  value = aws_elasticache_cluster.redis.cache_nodes[0].port
}


