resource "aws_s3_bucket" "storage" {
  bucket = "{{ .BucketName }}"
}

resource "aws_s3_bucket_acl" "storage" {
  bucket = aws_s3_bucket.storage.id
  acl    = "private"
}

resource "aws_s3_bucket_public_access_block" "storage" {
  bucket = aws_s3_bucket.storage.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
