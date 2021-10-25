# Terraform provider for running Bitbucket pipelines

Run Bitbucket pipelines straight from Terraform!

## Usage

```terraform
provider "bbss" {
  owner      = "the-owner-account"
  repository = "repo-slug"
}

resource "bbss_item" "this" {
  product   = "repo-slug"
  version   = "v1.0.0"
  variables = {
    VAR_NAME = "VAR_VALUE"
  }
}

output "outputs" {
  value = bpr_run.this.outputs
}

output "log" {
  value = bpr_run.this.log
}
```

### Example of using it to build a platform using self-service repos

```terraform
provider "bbss" {}

resource "bbss_item" "aws_account_prod" {
  product   = "aws-account"
  version   = "v1.0.0"
  variables = {
    name  = "something"
    email = "something@somewhere.com"
    tags  = jsondecode({ "Team" = "my-team" })
  }
}

resource "bbss_item" "aws_sso_group_prod" {
  product   = "aws-sso-group"
  version   = "v1.0.0"
  variables = {
    permission_sets  = jsonencode(["AWSReadOnly"])
    allowed_accounts = jsonencode([bbss_item.aws_account_prod.outputs.account_id])
    users            = jsonencode(["username"])
  }
}

resource "bbss_item" "service_user_prod" {
  product   = "aws-iam-user"
  version   = "v1.0.0"
  variables = {
    name = "service_user"
  }
}

output "service_user_prod" {
  value = bbss_item.service_user_prod.outputs.ssm_param_key
}
```

We strongly recommend you have a backend or commit the state file when using this provider, otherwise Terraform will
keep running your pipelines and using all your units up!

## Import

```bash
terraform import bbss_item.$key $uuid
```

## Supporting outputs

Currently, output support is a quick and dirty hack using substrings and JSON decoding which could and should be
improved, but for now here's the current convention.

Output from your pipelines should render like so.

```text
---OUTPUTS---
YOUR JSON OUTPUT HERE
---OUTPUTS---
```
