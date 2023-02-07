module "acm" {
  source  = "Selleo/acm/aws//modules/wildcard"
  version = "0.6.0"

  context = {
    namespace = "{{{ .Namespace }}}"
    stage     = "{{{ .Stage }}}"
    name      = "{{{ .Name }}}"
  }

  domain          = "{{{ .Subdomain }}}"
  validation_zone = "{{{ .Domain }}}"
  wildcard        = false
}
