# Terraform provider for running Bitbucket pipelines

Run Bitbucket pipelines straight from Terraform!

## Usage

```terraform
provider "bbss" {
  owner      = "the-owner-account"
  repository = "youre-repo-slug"
}

resource "bbss_item" "this" {
  product   = "youre-repo-slug"
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

## Supporting outputs

Currently, output support is a quick and dirty hack using substrings and JSON decoding which could and should be
improved, but for now here's the current convention.

Output from your pipelines should render like so.

```text
--- OUTPUT JSON START ---
YOUR JSON OUTPUT HERE
--- OUTPUT JSON STOP ---
```

This isn't easily achieved in Bitbucket pipelines as it will print your command because printing its output causing
unwanted behavior. So our best advice is to create a wrapper script to print your outputs and call it from your pipeline

- similar to below.

__terraform-output.sh__

```bash
#!/usr/bin/env sh
echo --- OUTPUT JSON START ---
terraform output -json
echo --- OUTPUT JSON STOP ---
```

__bitbucket-pipelines.yml__

```bash
pipelines:
  custom:
    do_a_thing:
    - step:
        script:
        - sh terraform-output.sh
```

### Complex output types

These are supported by the provider already, however they do get flattened before being made available.

#### Example

Output looking like this

```json
{
  "nested_object": {
    "nested_value": "example"
  }
}
```

Would look like this in the provider

```go
map[string]interface{}{
"nested_object.nested_value": "example"
}
```

Fortunately this has no impact on how you access the outputs within Terraform like so

```terraform
output "this" {
  value = bpr_run.this.outputs.nested_object.nested_value
}
```