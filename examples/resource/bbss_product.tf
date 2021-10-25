resource "bbss_product" "ecr" {
  repository = "aws-ecr-repo"
  version    = "v1.0.0"
  variables  = {
    VAR_KEY = "value"
  }
}