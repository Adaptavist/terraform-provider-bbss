# Make sure you use a remote backend with state locking, this will prevent end-users running your pipeline more than
# once at the same same, as that could cause quite a bit of pain. Most parts of the backend config can be hard-coded
# however, do not set the key, as this will need to be set but the invoking code to mitigate key clash.
terraform {
  backend "s3" {
    bucket         = "some-bucket-name"
    region         = "some-region-id"
    dynamodb_table = "some-table-name"
    kms_key_id     = "some-key-id"
  }
}

# Variables provide the user with an interface into your module, provide the right amount for the service your providing
variable "region" {
  type = string
}

variable "bucket_name" {
  type = string
}

# Configure your provider creds in the environment, think like region could be configurable by the user in many cases
provider "aws" {
  region = var.region
}

# create resources
resource "aws_s3_bucket" "this" {
  bucket = var.bucket_name
}

# provide some feedback
output "bucket_name" {
  value = aws_s3_bucket.this.bucket
}
