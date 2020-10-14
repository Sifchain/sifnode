terraform {
  backend "s3" {
    bucket = "{{.aws.state.bucket_name}}"
    key    = "{{.aws.state.bucket_key}}"
    region = "{{.aws.region}}"
  }
}
