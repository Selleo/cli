module "ecr" {
  source  = "Selleo/ecr/aws//modules/repository"
  version = "0.5.0"

  context = {
    namespace = "{{ .Namespace }}"
    stage     = "{{ .Stage }}"
    name      = "{{ .Name }}"
  }
}
