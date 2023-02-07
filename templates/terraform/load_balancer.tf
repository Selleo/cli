module "lb" {
  source  = "Selleo/lb/aws//modules//lb"
  version = "0.2.0"

  context = {
    namespace = "{{{ .Namespace }}}"
    stage     = "{{{ .Stage }}}"
    name      = "{{{ .Name }}}"
  }

  name        = "{{{ .LBName }}}"
  vpc_id      = module.vpc.id
  subnet_ids  = module.vpc.public_subnets
  force_https = true
}


module "lb_https" {
  source  = "Selleo/lb/aws//modules//lb"
  version = "0.2.0"

  load_balancer_arn = module.lb.id
  certificate_arn   = module.acm.arn
  target_group_arn  = module.service.lb_target_group_id
}
