data "aws_route53_zone" "main" {
  name = "{{ .Domain }}"
}

module "dns_load_balancer" {
  source  = "Selleo/route53/aws//modules/load-balancer-record"
  version = "0.4.0"

  lb_arn  = module.lb.id
  name    = "{{ .Subdomain }}"
  zone_id = data.aws_route53_zone.main.id
}
