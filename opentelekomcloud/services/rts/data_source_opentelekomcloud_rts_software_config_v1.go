package rts

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/softwareconfig"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRtsSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRtsSoftwareConfigV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"input_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"output_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRtsSoftwareConfigV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating RTS client: %w", err)
	}

	listOpts := softwareconfig.ListOpts{
		Id:   d.Id(),
		Name: d.Get("name").(string),
	}

	refinedConfigs, err := softwareconfig.List(orchestrationClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve RTS Software Configs: %s", err)
	}

	if len(refinedConfigs) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedConfigs) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Config := refinedConfigs[0]

	log.Printf("[INFO] Retrieved RTS Software Config using given filter %s: %+v", Config.Id, Config)
	d.SetId(Config.Id)

	n, err := softwareconfig.Get(orchestrationClient, Config.Id).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve RTS Software Config: %s", err)
	}
	mErr := multierror.Append(
		d.Set("name", Config.Name),
		d.Set("group", Config.Group),
		d.Set("region", config.GetRegion(d)),
		d.Set("config", n.Config),
		d.Set("options", n.Options),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("input_values", n.Inputs); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving inputs to state for OpenTelekomCloud RTS Software Config (%s): %s", d.Id(), err)
	}
	if err := d.Set("output_values", n.Outputs); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving outputs to state for OpenTelekomCloud RTS Software Config (%s): %s", d.Id(), err)
	}

	return nil
}
