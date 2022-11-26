resource "aws_route53_zone" "service" {
  name = "${var.environment}.${var.region}.${var.service_name}"

  vpc {
    vpc_id = var.vpc_id
  }
}

resource "aws_route53_record" "cache" {
  zone_id = aws_route53_zone.service.zone_id
  name    = "redis"
  type    = "CNAME"
  records = [aws_elasticache_cluster.cache.cache_nodes.0.address]
  ttl     = 60
}