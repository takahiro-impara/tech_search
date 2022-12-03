resource "aws_elasticache_cluster" "cache" {
  cluster_id           = "${var.environment}-${var.service_name}"
  engine               = "redis"
  engine_version       = var.redis.engine_version
  node_type            = var.redis.node_type
  num_cache_nodes      = var.redis.num_cache_nodes
  parameter_group_name = aws_elasticache_parameter_group.cache.name
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.cache.name
  security_group_ids   = [aws_security_group.cache.id]

  tags = {
    Name = "${var.environment}-${var.service_name}"
  }
}

resource "aws_elasticache_parameter_group" "cache" {
  name   = "${var.environment}-${var.service_name}-redis"
  family = "redis7"
}

resource "aws_elasticache_subnet_group" "cache" {
  name        = "${var.environment}-${var.service_name}-redis"
  description = "${var.service_name} redis subnet group"

  subnet_ids = data.aws_subnets.private.ids
}

resource "aws_security_group" "cache" {
  name        = "${var.environment}-${var.service_name}-redis"
  description = "${var.service_name} Redis security group"
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.environment}-${var.service_name}-redis"
  }
}

resource "aws_security_group_rule" "cache-egress" {
  depends_on = [aws_security_group.cache]
  type       = "egress"
  from_port  = 0
  to_port    = 0
  protocol   = "-1"

  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.cache.id
}

resource "aws_security_group_rule" "cache-ingress" {
  description = "internal access"
  type        = "ingress"
  from_port   = 6379
  to_port     = 6379
  protocol    = "tcp"

  cidr_blocks = ["10.0.0.0/16"]

  security_group_id = aws_security_group.cache.id
}

resource "aws_security_group_rule" "bastion-ingress" {
  description = "internal bastion access"
  type        = "ingress"
  from_port   = 6379
  to_port     = 6379
  protocol    = "tcp"

  source_security_group_id = aws_security_group.bastion.id

  security_group_id = aws_security_group.cache.id
}