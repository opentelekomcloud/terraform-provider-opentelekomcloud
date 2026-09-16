package gemini

import (
	"context"
	"sort"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/spec"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

// flavorsPageLimit is the maximum page size accepted by the specification API.
const flavorsPageLimit = 100

func DataSourceGeminiFlavorsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiFlavorsRead,

		Schema: map[string]*schema.Schema{
			"engine_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cassandra",
				ValidateFunc: validation.StringInSlice([]string{
					"cassandra", "influxdb",
				}, false),
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CloudNativeCluster",
				}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"availability_zone": {
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
						"engine_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func defaultFlavorMode(engineName string) string {
	if engineName == "influxdb" {
		return "CloudNativeCluster"
	}
	return ""
}

// listAllFlavors walks over all specification pages, as the API returns at most 100 entries at once.
func listAllFlavors(client *golangsdk.ServiceClient, opts spec.ListFlavorsOpts) ([]spec.Flavors, error) {
	var allFlavors []spec.Flavors
	offset := 0

	for {
		opts.Offset = pointerto.Int(offset)
		opts.Limit = pointerto.Int(flavorsPageLimit)

		page, err := spec.ListFlavors(client, opts)
		if err != nil {
			return nil, err
		}

		allFlavors = append(allFlavors, page.Flavors...)
		offset += len(page.Flavors)

		if len(page.Flavors) < flavorsPageLimit || offset >= page.TotalCount {
			return allFlavors, nil
		}
	}
}

func filterFlavors(d *schema.ResourceData, flavors []spec.Flavors) ([]map[string]interface{}, []string) {
	version := d.Get("engine_version").(string)
	az := d.Get("availability_zone").(string)
	vcpus, vcpusSet := d.GetOk("vcpus")
	ram, ramSet := d.GetOk("ram")

	result := make([]map[string]interface{}, 0, len(flavors))
	specCodes := make([]string, 0, len(flavors))

	for _, flavor := range flavors {
		if version != "" && version != flavor.EngineVersion {
			continue
		}
		if vcpusSet && strconv.Itoa(vcpus.(int)) != flavor.Vcpus {
			continue
		}
		if ramSet && strconv.Itoa(ram.(int)) != flavor.Ram {
			continue
		}

		azList := make([]string, 0, len(flavor.AzStatus))
		for zone, status := range flavor.AzStatus {
			if status != "normal" {
				continue
			}
			if az != "" && az != zone {
				continue
			}
			azList = append(azList, zone)
		}
		if len(azList) == 0 {
			continue
		}
		sort.Strings(azList)

		result = append(result, map[string]interface{}{
			"spec_code":          flavor.SpecCode,
			"engine_name":        flavor.EngineName,
			"engine_version":     flavor.EngineVersion,
			"vcpus":              flavor.Vcpus,
			"ram":                flavor.Ram,
			"availability_zones": azList,
			"az_status":          flavor.AzStatus,
		})
		specCodes = append(specCodes, flavor.SpecCode)
	}

	return result, specCodes
}

func dataSourceGeminiFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV31Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB v3.1 client: %s", err)
	}

	engineName := d.Get("engine_name").(string)
	mode := d.Get("mode").(string)
	if mode == "" {
		mode = defaultFlavorMode(engineName)
	}

	allFlavors, err := listAllFlavors(client, spec.ListFlavorsOpts{
		EngineName: engineName,
		Mode:       mode,
	})
	if err != nil {
		return diag.Errorf("error getting GeminiDB flavors: %s", err)
	}

	flavors, specCodes := filterFlavors(d, allFlavors)
	if len(flavors) == 0 {
		return diag.Errorf("no GeminiDB flavor found matching the given criteria, " +
			"please change your search criteria and try again")
	}

	d.SetId(hashcode.Strings(specCodes))

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("flavors", flavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
