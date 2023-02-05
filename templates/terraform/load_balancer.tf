module "lb" {
  source  = "Selleo/lb/aws"
  version = "0.1.0"

  context = {
    namespace = "{{ .Namespace }}"
    stage     = "{{ .Stage }}"
    name      = "{{ .Name }}"
  }

  name        = "{{ .LBName }}"
  vpc_id      = module.vpc.id
  subnet_ids  = module.vpc.public_subnets
  force_https = true
}

resource "aws_alb_listener" "https" {
  load_balancer_arn = module.lb.id
  port              = 443
  protocol          = "HTTPS"
  certificate_arn   = module.acm.arn

  # https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html
  ssl_policy = "ELBSecurityPolicy-TLS-1-2-Ext-2018-06"

  default_action {
    target_group_arn = module.service.lb_target_group_id
    type             = "forward"
  }
}
