package rts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/stackresources"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRTSStackResourcesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRTSStackResourcesV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"required_by": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"resource_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"physical_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceRTSStackResourcesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating RTS client: %w", err)
	}

	listOpts := stackresources.ListOpts{
		Name:       d.Get("resource_name").(string),
		PhysicalID: d.Get("physical_resource_id").(string),
		Type:       d.Get("resource_type").(string),
	}

	refinedResources, err := stackresources.List(orchestrationClient, d.Get("stack_name").(string), listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve Stack Resources: %s", err)
	}

	if len(refinedResources) < 1 {
		return fmterr.Errorf("no matching resource found. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedResources) > 1 {
		return fmterr.Errorf("multiple resources matched; use additional constraints to reduce matches to a single resource")
	}

	stackResource := refinedResources[0]
	d.SetId(stackResource.PhysicalID)

	mErr := multierror.Append(
		d.Set("resource_name", stackResource.Name),
		d.Set("resource_status", stackResource.Status),
		d.Set("logical_resource_id", stackResource.LogicalID),
		d.Set("physical_resource_id", stackResource.PhysicalID),
		d.Set("required_by", stackResource.RequiredBy),
		d.Set("resource_status_reason", stackResource.StatusReason),
		d.Set("resource_type", stackResource.Type),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
