package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwaredeployment"
)

func dataSourceRtsSoftwareDeploymentV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRTSSoftwareDeploymentV1Read,

		Schema: map[string]*schema.Schema{ //request and response parameters
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

func dataSourceRTSSoftwareDeploymentV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	listOpts := softwaredeployment.ListOpts{
		Id:       d.Id(),
		ServerId: d.Get("server_id").(string),
		ConfigId: d.Get("config_id").(string),
		Action:   d.Get("action").(string),
		Status:   d.Get("status").(string),
	}

	refinedDeployments, err := softwaredeployment.List(orchestrationClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve RTS Software Deployment: %s", err)
	}

	if len(refinedDeployments) < 1 {
		return fmt.Errorf("No matching resource found. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedDeployments) > 1 {
		return fmt.Errorf("Multiple resources matched; use additional constraints to reduce matches to a single resource")
	}

	stackResource := refinedDeployments[0]
	d.SetId(stackResource.Id)

	d.Set("id", stackResource.Id)
	d.Set("status", stackResource.Status)
	d.Set("server_id", stackResource.ServerId)
	d.Set("config_id", stackResource.ConfigId)
	d.Set("status_reason", stackResource.StatusReason)
	d.Set("action", stackResource.Action)
	d.Set("output_values", stackResource.OutputValues)
	d.Set("input_values", stackResource.InputValues)
	d.Set("region", GetRegion(d, config))
	return nil
}
