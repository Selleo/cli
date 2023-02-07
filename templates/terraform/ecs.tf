module "cluster" {
  source  = "Selleo/ecs/aws//modules/cluster"
  version = "0.11.0"

  context = {
    namespace = "{{{ .Namespace }}}"
    stage     = "{{{ .Stage }}}"
    name      = "{{{ .Name }}}"
  }

  name_prefix          = "{{{ .Namespace }}}"
  vpc_id               = module.vpc.id
  subnet_ids           = module.vpc.public_subnets
  instance_type        = "{{{ .ECSInstanceType }}}"
  lb_security_group_id = module.lb.security_group_id

  autoscaling_group = {
    min_size         = {{{ .ECSMinSize }}}
    max_size         = 5
    desired_capacity = 1
  }
}

module "service" {
  source  = "Selleo/ecs/aws//modules/service"
  version = "0.11.0"

  context = {
    namespace = "{{{ .Namespace }}}"
    stage     = "{{{ .Stage }}}"
    name      = "{{{ .Name }}}"
  }

  name          = "{{{ .Name }}}"
  vpc_id        = module.vpc.id
  subnet_ids    = module.vpc.public_subnets
  cluster_id    = module.cluster.id
  desired_count = 1
  secrets       = ["/{{{ .Namespace }}}/{{{ .Stage }}}/{{{ .Name }}}"]

  container = {
    mem_reservation_units = {{{ .ECSServiceMinMemory }}}
    cpu_units             = {{{ .ECSServiceCpu }}}
    mem_units             = {{{ .ECSServiceMaxMemory }}}

    image = module.ecr.url_tagged_latest
    port  = {{{ .ECSServicePort }}}
  }
  one_off_commands = [
{{{ range .ECSOneOffs -}}}
    "{{{ . }}}",
{{{ end -}}}
  ]

  health_check = {
    path    = "/health"
    matcher = "200"
  }

  # useful for staging and smaller machines (can cause downtime)
  # deployment_minimum_healthy_percent = 0
  # deregistration_delay               = 15

  depends_on = [module.secrets, module.secrets_editable]
}
