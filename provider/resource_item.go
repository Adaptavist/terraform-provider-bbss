package provider

import (
	"context"
	"fmt"
	"github.com/adaptavist/terraform-provider-bitbucket-bbss/provider/log"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceItem defined the item resource within the provider
func resourceItem() *schema.Resource {
	return &schema.Resource{
		Description:   "An item is an instance of a product which is a Bitbucket repository constructed in an opinionated way so it can be executed remotely providing self-service infrastructure.",
		CreateContext: resourceItemApply,     // Create runs a terraform apply operation in a remote pipeline
		ReadContext:   resourceItemDoNothing, // We never read from Bitbucket to refresh state, only run pipelines
		UpdateContext: resourceItemApply,     // Update runs a terraform apply operation in a remote pipeline
		DeleteContext: resourceItemDestroy,   // Destroy runs a terraform destroy operation in a remote pipeline
		Schema:        resourceItemSchema(),
	}
}

func resourceItemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"product": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The repository containing the target pipeline",
		},
		"version": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Description:   "The tag containing the target pipeline",
		},
		"variables": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Map of variables for the pipeline",
		},
		"log": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The full pipeline output",
		},
		"outputs": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "Map of outputs, if we can grab them from the pipeline text output",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The UUID of the pipeline, use for looking up its status",
		},
	}
}

func resourceItemApply(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("Applying item")
	config := meta.(*config)
	id := d.Get("id").(string)

	if id == "" {
		id = uuid.New().String()
		log.Debug(fmt.Sprintf("Created item %s of \"%s\"", config.KeyIDKey, id))
		d.SetId(id)
	}

	return run(id, "apply", ctx, d, meta)
}

// resourceItemDoNothing only exists because Terraform requires some operations exist
func resourceItemDoNothing(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Debug(fmt.Sprintf("Doing nothing on item %s", d.Get("id").(string)))
	return nil
}

// resourceItemDestroy runs the destroy pipeline, only removing the resource id on success.
func resourceItemDestroy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug(fmt.Sprintf("Destroying item %s", d.Get("id").(string)))
	id := d.Get("id").(string)

	// only run the destroy pipeline if we have an id
	if id != "" {
		err := run(id, "destroy", ctx, d, meta)
		if err != nil {
			return err
		}
		// only unset the id if the run has finished without error, as we may need to run it again
		d.SetId("")
	}

	return nil
}
