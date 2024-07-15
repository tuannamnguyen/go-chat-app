locals {
  domain_name = "gochatapp.sbs"
}

resource "aws_route53_zone" "chat_app_zone" {
  name          = local.domain_name
  force_destroy = true
}

resource "aws_route53_record" "route_to_alb" {
  zone_id = aws_route53_zone.chat_app_zone.zone_id
  name    = local.domain_name
  type    = "A"

  alias {
    name                   = aws_lb.chat_app_lb.dns_name
    zone_id                = aws_lb.chat_app_lb.zone_id
    evaluate_target_health = false
  }
}
resource "aws_route53_record" "route_to_alb_subdomain" {
  name    = "www"
  zone_id = aws_route53_zone.chat_app_zone.zone_id
  type    = "CNAME"
  ttl     = 300
  records = [aws_lb.chat_app_lb.dns_name]
}
