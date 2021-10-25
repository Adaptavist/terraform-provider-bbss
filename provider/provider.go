package provider

import (
	"context"
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ProviderName = "bbss"

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"base_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BBSS_BASE_URL", "https://api.bitbucket.org"),
					Description: "Customise the base URL for API calls, this is used for testing purposes",
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BBSS_USERNAME", nil),
					Description: "Bitbucket pipelines username. The user must have repo:admin, and pipeline:write scopes",
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BBSS_PASSWORD", nil),
					Description: "Bitbucket pipelines app password",
				},
				"owner": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BBSS_OWNER", nil),
					Description: "Bitbucket workspace slug, where the pipeline repos exist",
				},
				"key_id_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "KEY_ID",
					DefaultFunc: schema.EnvDefaultFunc("BBSS_KEY_ID_KEY", nil),
					Description: "The name of variable to assign a v4 UUID used in each resource for a backend key in the remote pipeline",
				},
				"output_delimiter": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "---OUTPUTS---",
					DefaultFunc: schema.EnvDefaultFunc("BBSS_OUTPUT_DELIMITER", nil),
					Description: "String used to detected the beginning and end outputs which are used to extract from pipeline logs",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"bbss_item": resourceItem(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

type config struct {
	OutputDelimiter string
	KeyIDKey        string
	Client          *client.Client
}

func configure(_ string, _ *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &config{
			OutputDelimiter: d.Get("output_delimiter").(string),
			KeyIDKey:        d.Get("key_id_key").(string),
			Client: &client.Client{
				Config: &client.Config{
					Username:  d.Get("username").(string),
					Password:  d.Get("password").(string),
					BaseURL:   d.Get("base_url").(string),
					Workspace: valueAsPointer("owner", d),
				},
			},
		}, nil
	}
}

func valueAsPointer(k string, d *schema.ResourceData) *string {
	if v := d.Get(k).(string); v != "" {
		return &v
	}
	return nil
}
