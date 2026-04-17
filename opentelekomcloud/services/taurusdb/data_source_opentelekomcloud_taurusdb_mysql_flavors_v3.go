package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/instance"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBV3MysqlFlavorsRead,

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "gaussdb-mysql",
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "8.0",
			},
			"availability_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "single",
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBV3MysqlFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	opts := instance.ListFlavorsOpts{
		DatabaseName:         d.Get("engine").(string),
		VersionName:          d.Get("version").(string),
		AvailabilityZoneMode: d.Get("availability_zone_mode").(string),
	}

	flavorList, err := instance.ListFlavors(client, opts)
	if err != nil {
		return diag.Errorf("error fetching flavors for TaurusDB MySQL: %s", err)
	}

	flavors := make([]interface{}, 0, len(flavorList))
	for _, flavor := range flavorList {
		flavors = append(flavors, map[string]interface{}{
			"name":      flavor.SpecCode,
			"vcpus":     flavor.Vcpus,
			"memory":    flavor.Ram,
			"type":      flavor.Type,
			"mode":      flavor.InstanceMode,
			"version":   flavor.VersionName,
			"az_status": flavor.AzStatus,
		})
	}

	d.SetId("flavors")

	var mErr *multierror.Error
	mErr = multierror.Append(mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("flavors", flavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
