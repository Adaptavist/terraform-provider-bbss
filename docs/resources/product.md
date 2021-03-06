---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bbss_product Resource - terraform-provider-bbss"
subcategory: ""
description: |-
  An item is an instance of a product which is a Bitbucket repository constructed in an opinionated way so it can be executed remotely providing self-service infrastructure.
---

# bbss_product (Resource)

An item is an instance of a product which is a Bitbucket repository constructed in an opinionated way so it can be executed remotely providing self-service infrastructure.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **product** (String) The repository containing the target pipeline
- **variables** (Map of String) Map of variables for the pipeline
- **version** (String) The tag containing the target pipeline

### Read-Only

- **id** (String) The UUID of the pipeline, use for looking up its status
- **log** (String) The full pipeline output
- **outputs** (Map of String) Map of outputs, if we can grab them from the pipeline text output


