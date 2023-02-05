# https://github.com/Selleo/terraform-aws-rds/tree/main/modules/postgres
module "db" {
  source = "Selleo/rds/aws///modules/postgres"
  version = "0.1.0"

  context = {
    namespace = "{{ .Namespace }}"
    stage     = "{{ .Stage }}"
    name      = "{{ .Name }}"
  }

  vpc = {
    id           = module.vpc.id
    cidr         = module.vpc.cidr_block
    subnet_group = module.vpc.database_subnet_group
  }

  identifier = "{{ .DBIdentifier }}"
  db_name    = "{{ .DBName }}"
  db_user    = "{{ .DBUser }}"

  multi_az          = {{ .DBMultiAZ }}
  apply_immediately = {{ .DBApplyImmediately }}

  # instance_class         = "db.t4g.micro"
  # engine_version         = "14.5"
  # parameter_group_family = "postgres14"
  # allocated_storage      = 20
  # max_allocated_storage  = 100
  # maintenance_window     = "Mon:00:00-Mon:02:00"
  # backup_window          = "02:30-03:30"
}
