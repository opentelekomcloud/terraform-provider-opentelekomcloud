package rts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/softwaredeployment"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRtsSoftwareDeploymentV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRTSSoftwareDeploymentV1Read,

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_values": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_values": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceRTSSoftwareDeploymentV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating RTS client: %w", err)
	}

	listOpts := softwaredeployment.ListOpts{
		Id:       d.Id(),
		ServerId: d.Get("server_id").(string),
		ConfigId: d.Get("config_id").(string),
		Action:   d.Get("action").(string),
		Status:   d.Get("status").(string),
	}

	refinedDeployments, err := softwaredeployment.List(orchestrationClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve RTS Software Deployment: %s", err)
	}

	if len(refinedDeployments) < 1 {
		return fmterr.Errorf("no matching resource found. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedDeployments) > 1 {
		return fmterr.Errorf("multiple resources matched; use additional constraints to reduce matches to a single resource")
	}

	stackResource := refinedDeployments[0]
	d.SetId(stackResource.Id)

	mErr := multierror.Append(
		d.Set("id", stackResource.Id),
		d.Set("status", stackResource.Status),
		d.Set("server_id", stackResource.ServerId),
		d.Set("config_id", stackResource.ConfigId),
		d.Set("status_reason", stackResource.StatusReason),
		d.Set("action", stackResource.Action),
		d.Set("output_values", stackResource.OutputValues),
		d.Set("input_values", stackResource.InputValues),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
