resource "aws_elasticache_replication_group" "cache-cluster" {
  replication_group_id       = "${var.environment}-v2-${var.service_name}"
  description                = "${var.environment}-${var.service_name}"
  engine                     = "redis"
  engine_version             = var.redis.engine_version
  node_type                  = var.redis.node_type
  parameter_group_name       = aws_elasticache_parameter_group.cache.name
  port                       = 6379
  automatic_failover_enabled = true
  subnet_group_name          = aws_elasticache_subnet_group.cache.name
  security_group_ids         = [aws_security_group.cache.id]
  multi_az_enabled           = true
  num_cache_clusters         = 2

  tags = {
    Name    = "${var.environment}-${var.service_name}"
    Country = "gl"
    Brand   = "fr"
  }
}