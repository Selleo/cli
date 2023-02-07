module "vpc" {
  source  = "Selleo/vpc/aws//modules/vpc"
  version = "0.5.0"

  context = {
    namespace = "{{{ .Namespace }}}"
    stage     = "{{{ .Stage }}}"
    name      = "{{{ .Name }}}"
  }

  availability_zone_identifiers = ["a", "b", "c"]
  cidr                          = "10.0.0.0/16"
  private_subnets               = ["10.0.1.0/24", "10.0.2.0/24", "10.0.2.0/24"]
  database_subnets              = ["10.0.51.0/24", "10.0.52.0/24", "10.0.53.0/24"]
  public_subnets                = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  single_nat_gateway = true
}
