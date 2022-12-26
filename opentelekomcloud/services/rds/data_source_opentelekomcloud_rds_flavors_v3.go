package rds

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/flavors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRdsFlavorV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MySQL", "PostgreSQL", "SQLServer",
				}, false),
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsFlavorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	listOpts := flavors.ListOpts{
		VersionName:  d.Get("db_version").(string),
		SpecCode:     d.Get("instance_mode").(string),
		DatabaseName: d.Get("db_type").(string),
	}
	allFlavorsList, err := flavors.ListFlavors(client, listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	var refinedFlavors []map[string]interface{}
	instanceMode := d.Get("instance_mode").(string)
	for _, flavor := range allFlavorsList {
		if flavor.InstanceMode == instanceMode {
			refinedFlavors = append(
				refinedFlavors,
				map[string]interface{}{
					"vcpus":     flavor.VCPUs,
					"memory":    flavor.RAM,
					"name":      flavor.SpecCode,
					"mode":      flavor.InstanceMode,
					"az_status": flavor.AzStatus,
				},
			)
		}
	}

	if err := d.Set("flavors", refinedFlavors); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("flavors")
	return nil
}
