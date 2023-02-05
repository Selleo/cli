terraform {
  required_version = "~> 1.0"

  cloud {
    organization = "{{ .TerraformCloudOrganization }}"

    workspaces {
      name = "{{ .TerraformCloudWorkspace }}"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.4.3"
    }
  }
}

provider "aws" {
  region = var.region
}

provider "aws" {
  alias  = "global"
  region = "us-east-1"
}
