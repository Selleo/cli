module "iam_app" {
  source  = "Selleo/iam/aws//modules/user-with-access-key"
  version = "0.6.0"

  name = "{{{ .IAMApp }}}"
}

module "iam_app_allow_s3_read_write" {
  source  = "Selleo/iam/aws//modules/s3-read-write"
  version = "0.5.0"

  bucket_arn  = aws_s3_bucket.storage.arn
  name_prefix = "{{{ .Namespace }}}-{{{ .Stage }}}-{{{ .Name }}}"
  users       = [module.iam_app.name]
}
