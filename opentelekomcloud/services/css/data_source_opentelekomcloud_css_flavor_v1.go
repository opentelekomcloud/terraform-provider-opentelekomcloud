package css

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/css/v1/flavors"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCSSFlavorV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCSSFlavorV1Read,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ess",
			},

			"min_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"disk_range": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_from": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"min_to": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"from": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"to": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCSSFlavorV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CssV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSS v1 client: %s", err)
	}

	versions, err := flavors.List(client)
	if err != nil {
		return fmterr.Errorf("error extracting versions")
	}

	opts := flavors.FilterOpts{
		Version: d.Get("version").(string),
		Type:    d.Get("type").(string),
	}

	filtered := flavors.FilterVersions(versions, opts)

	if len(filtered) == 0 {
		return fmterr.Errorf("can't find version matching criteria: %+v", opts)
	}

	result := findFlavorInVersions(d, filtered)
	if result == nil {
		return fmterr.Errorf("can't find flavor matching criteria")
	}

	d.SetId(result.FlavorID)

	mErr := multierror.Append(
		d.Set("name", result.Name),
		d.Set("region", result.Region),
		d.Set("ram", result.RAM),
		d.Set("cpu", result.CPU),

		setDiskRange(d, result),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	return nil
}

func setDiskRange(d *schema.ResourceData, flavor *flavors.Flavor) error {
	diskRange := d.Get("disk_range").([]interface{})
	var item map[string]interface{}
	if len(diskRange) == 0 {
		item = make(map[string]interface{})
		diskRange = make([]interface{}, 1)
	} else {
		item = diskRange[0].(map[string]interface{})
	}
	item["from"] = flavor.DiskMin
	item["to"] = flavor.DiskMin
	diskRange[0] = item
	return d.Set("disk_range", diskRange)
}

func findFlavorInVersions(d *schema.ResourceData, versions []flavors.Version) *flavors.Flavor {
	if name := d.Get("name").(string); name != "" {
		return findFlavorByName(versions, name)
	}

	opts := flavors.FilterOpts{}

	if minCPU := d.Get("min_cpu").(int); minCPU != 0 {
		opts.CPU = &flavors.Limit{Min: minCPU}
	}

	if minRAM := d.Get("min_ram").(int); minRAM != 0 {
		opts.RAM = &flavors.Limit{Min: minRAM}
	}

	if d.Get("disk_range.#").(int) != 0 {
		minFrom := d.Get("disk_range.0.min_from").(int)
		if minFrom != 0 {
			opts.DiskMin = &flavors.Limit{Min: minFrom}
		}
		minTo := d.Get("disk_range.0.min_to").(int)
		if minFrom != 0 {
			opts.DiskMax = &flavors.Limit{Min: minTo}
		}
	}

	return flavors.FindFlavor(versions, opts)
}

func findFlavorByName(versions []flavors.Version, name string) *flavors.Flavor {
	for _, version := range versions {
		for _, flavor := range version.Flavors {
			if flavor.Name == name {
				return &flavor
			}
		}
	}
	return nil
}
