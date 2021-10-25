package provider

import (
	"context"
	"fmt"
	"github.com/adaptavist/bitbucket_pipelines_client/builders"
	"github.com/adaptavist/bitbucket_pipelines_client/model"
	"github.com/adaptavist/terraform-provider-bitbucket-bbss/provider/outputs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPipelineSteps(d *schema.ResourceData, meta interface{}, pipeline *model.Pipeline) (model.PipelineSteps, error) {
	cli := meta.(*config).Client
	return cli.GetPipelineSteps(model.GetPipelineRequest{
		Repository: valueAsPointer("product", d),
		Pipeline:   pipeline,
	})
}

func getPipelineStepLog(d *schema.ResourceData, meta interface{}, pipeline *model.Pipeline, step *model.PipelineStep) (string, error) {
	cli := meta.(*config).Client
	log, err := cli.GetPipelineStepLog(model.GetPipelineStepRequest{
		Repository:   valueAsPointer("product", d),
		Pipeline:     pipeline,
		PipelineStep: step,
	})
	return string(log), err
}

func run(keyID, action string, _ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cli := meta.(*config).Client
	product := d.Get("product").(string)
	version := d.Get("version").(string)
	pipeline := builders.Pipeline().Variable("KEY_ID", keyID, false)
	target := builders.Target().Pattern(action)

	// apply any additional variables
	for key, value := range d.Get("variables").(map[string]interface{}) {
		pipeline.Variable(key, fmt.Sprintf("%v", value), false)
	}

	// get the commit hash of the release
	tagResponse, err := cli.GetTag(model.GetTagRequest{
		Repository: &product,
		Tag:        version,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	target.Tag(version, tagResponse.Target.Hash)
	pipeline.Target(target.Build())

	// complete the request object for the pipeline

	request := model.PostPipelineRequest{
		Repository: &product,
		Pipeline:   pipeline.Build(),
	}

	// run the pipeline
	response, err := cli.RunPipeline(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// get the completed pipeline steps

	steps, err := getPipelineSteps(d, meta, response)

	if err != nil {
		return diag.FromErr(err)
	}

	// extract the pipeline output and its containing outputMap

	output := ""
	outputMap := map[string]interface{}{}

	for _, step := range steps {
		logStr, err := getPipelineStepLog(d, meta, response, &step)

		if err != nil {
			return diag.FromErr(err)
		}

		output = output + "\n" + logStr
		outputMap, err = outputs.ExtractAndAppend(meta.(*config).OutputDelimiter, logStr, outputMap)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = d.Set("outputs", outputMap)

	if err != nil {
		return diag.FromErr(err)
	}

	// set the pipeline output

	err = d.Set("log", output)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
