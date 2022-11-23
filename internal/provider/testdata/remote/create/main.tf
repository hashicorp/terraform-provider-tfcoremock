provider "tfcoremock" {}

resource "null_resource" "test" {
  provider = tfcoremock
}

resource "aws_route53_zone" "primary" {
  provider = tfcoremock

  name = "primary.com"
}

resource "aws_route53_record" "www" {
  provider = tfcoremock

  zone_id = aws_route53_zone.primary.id
  name = "www"
  type = "A"
  ttl = 300
  records = [
    "192.168.0.1"
  ]
}
