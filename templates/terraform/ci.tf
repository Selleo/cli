module "iam_ci" {
  source  = "Selleo/iam/aws//modules/user-with-access-key"
  version = "0.6.0"

  name = "{{ .IAMCI }}"

  groups = [
    module.cluster.deployment_group,
    module.service.deployment_group,
    module.ecr.deployment_group
  ]
}

module "secrets_ci" {
  source  = "Selleo/ssm/aws//modules/parameters"
  version = "0.3.0"

  context = {
    namespace = "{{ .Namespace }}"
    stage     = "{{ .Stage }}"
    name      = "{{ .Name }}"
  }

  path = "/ci/{{ .Namespace }}/{{ .Stage }}"

  secrets = {
    AWS_REGION            = var.region
    AWS_ACCESS_KEY_ID     = module.iam_ci.key_id
    AWS_SECRET_ACCESS_KEY = module.iam_ci.key_secret
  }
}
