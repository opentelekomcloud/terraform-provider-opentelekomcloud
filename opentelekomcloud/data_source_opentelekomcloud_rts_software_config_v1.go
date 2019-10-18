package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
)

func dataSourceRtsSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRtsSoftwareConfigV1Read,

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

func dataSourceRtsSoftwareConfigV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	listOpts := softwareconfig.ListOpts{
		Id:   d.Id(),
		Name: d.Get("name").(string),
	}

	refinedConfigs, err := softwareconfig.List(orchestrationClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve RTS Software Configs: %s", err)
	}

	if len(refinedConfigs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedConfigs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Config := refinedConfigs[0]

	log.Printf("[INFO] Retrieved RTS Software Config using given filter %s: %+v", Config.Id, Config)
	d.SetId(Config.Id)

	d.Set("name", Config.Name)
	d.Set("group", Config.Group)
	d.Set("region", GetRegion(d, config))

	n, err := softwareconfig.Get(orchestrationClient, Config.Id).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve RTS Software Config: %s", err)
	}

	d.Set("config", n.Config)
	d.Set("options", n.Options)
	if err := d.Set("input_values", n.Inputs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving inputs to state for OpenTelekomCloud RTS Software Config (%s): %s", d.Id(), err)
	}
	if err := d.Set("output_values", n.Outputs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving outputs to state for OpenTelekomCloud RTS Software Config (%s): %s", d.Id(), err)
	}

	return nil
}
