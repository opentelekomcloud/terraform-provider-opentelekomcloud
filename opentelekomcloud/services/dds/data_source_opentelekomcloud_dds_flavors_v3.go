package dds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/flavors"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDdsFlavorV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"engine_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mongos", "shard", "config", "replica",
				}, true),
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
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

func dataSourceDdsFlavorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ddsClient, err := config.DdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DDS client: %s", err)
	}

	listOpts := flavors.ListFlavorOpts{
		EngineName: d.Get("engine_name").(string),
	}

	extractedFlavors, err := flavors.List(ddsClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to all flavor pages: %s", err)
	}

	matchFlavorList := make([]map[string]interface{}, 0)
	expectedType := d.Get("type").(string)
	expectedVcpus := d.Get("vcpus").(string)
	expectedMemory := d.Get("memory").(string)

	for _, item := range extractedFlavors.Flavors {
		if matchesFilters(item, expectedType, expectedVcpus, expectedMemory) {
			continue
		}

		matchFlavor := map[string]interface{}{
			"spec_code": item.SpecCode,
			"type":      item.Type,
			"vcpus":     item.Vcpus,
			"memory":    item.Ram,
			"az_status": item.AZStatus,
		}
		matchFlavorList = append(matchFlavorList, matchFlavor)
	}

	if len(matchFlavorList) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	d.SetId("flavors")
	mErr := multierror.Append(nil,
		d.Set("flavors", matchFlavorList),
		d.Set("region", config.GetRegion(d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func matchesFilters(item flavors.FlavorResponse, flavorType, flavorVcpus, flavorMemory string) bool {
	if flavorType != "" && flavorType != item.Type {
		return true
	}
	if flavorVcpus != "" && flavorVcpus != item.Vcpus {
		return true
	}
	if flavorMemory != "" && flavorMemory != item.Ram {
		return true
	}

	return false
}
