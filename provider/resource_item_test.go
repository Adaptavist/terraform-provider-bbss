package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testBitbucketPipelinesRun = `
resource "bbss_item" "this" {
	product = "my-product"
	version = "v1.0.0"
	variables  = {
		"var_key" = "var_value"
	}
}
`

func TestBitbucketPipelines_Run(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testBitbucketPipelinesRun,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("bbss_item.this", "product", regexp.MustCompile("^my-product$")),
					resource.TestMatchResourceAttr("bbss_item.this", "version", regexp.MustCompile("^v1.0.0$")),
					resource.TestMatchResourceAttr("bbss_item.this", "outputs.string_output", regexp.MustCompile("value")),
					resource.TestMatchResourceAttr("bbss_item.this", "outputs.map_output.map_key", regexp.MustCompile("value")),
					resource.TestMatchResourceAttr("bbss_item.this", "log", regexp.MustCompile("Some other logging info before outputs")),
				),
			},
		},
	})
}
