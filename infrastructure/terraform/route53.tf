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
  records = [aws_elasticache_replication_group.cache-cluster.primary_endpoint_address]
  ttl     = 60
}

data "aws_route53_zone" "public" {
  name = "staging.udacity.impara8.com"
}

resource "aws_route53_record" "cloudfront" {
  count   = 1
  zone_id = data.aws_route53_zone.public.zone_id
  name    = "tech-search-dist"
  type    = "CNAME"
  records = [data.aws_cloudfront_distribution.frontend.domain_name]
  ttl     = 3600
}