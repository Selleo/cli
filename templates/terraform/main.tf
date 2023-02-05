module "secrets" {
  source  = "Selleo/ssm/aws//modules/parameters"
  version = "0.3.0"

  context = {
    namespace = "{{ .Namespace }}"
    stage     = "{{ .Stage }}"
    name      = "{{ .Name }}"
  }

  secrets = {
    AWS_ACCESS_KEY_ID     = module.iam_app.key_id
    AWS_SECRET_ACCESS_KEY = module.iam_app.key_secret
    AWS_REGION            = var.region
    AWS_S3_BUCKET         = aws_s3_bucket.storage.bucket
    DATABASE_URL          = module.db.database_url
  }
}

module "secrets_editable" {
  source  = "Selleo/ssm/aws//modules/editable-parameters"
  version = "0.3.0"

  context = {
    namespace = "{{ .Namespace }}"
    stage     = "{{ .Stage }}"
    name      = "{{ .Name }}"
  }

  # ⚠️n new secrets are initialized with value set in terraform,
  # but any further changes are ignored, editable secrets
  # are meant to be edited via UI.
  secrets = {
    # GOOGLE_CLIENT_ID = "Edit in UI"
    # GOOGLE_CLIENT_SECRET = "Edit in UI"
  }
}
