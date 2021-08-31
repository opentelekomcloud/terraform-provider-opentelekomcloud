package ecs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/flavors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceComputeFlavorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name", "min_ram", "min_disk"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"flavor_id"},
			},
			"min_ram": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"flavor_id"},
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_disk": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"flavor_id"},
			},
			"disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"swap": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rx_tx_factor": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extra_specs": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func checkAZExtraOption(es map[string]string, az string) bool {
	zoneNormalState := fmt.Sprintf("%s(normal)", az)

	if azProperty, ok := es["cond:operation:az"]; ok {
		zones := strings.Split(azProperty, ",")
		for _, zone := range zones {
			if zone != zoneNormalState {
				continue
			}
			return true
		}
	}

	return false
}

func dataSourceComputeFlavorV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	var allFlavors []flavors.Flavor
	if flavorID := d.Get("flavor_id").(string); flavorID != "" {
		flavor, err := flavors.Get(client, flavorID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmterr.Errorf("no flavor found")
			}
			return fmterr.Errorf("unable to retrieve OpenTelekomCloud %s flavor: %w", flavorID, err)
		}

		allFlavors = append(allFlavors, *flavor)
	} else {
		listOpts := flavors.ListOpts{
			MinDisk: d.Get("min_disk").(int),
			MinRAM:  d.Get("min_ram").(int),
		}

		log.Printf("[DEBUG] opentelekoncloud_compute_flavor_v2 ListOpts: %#v", listOpts)

		allPages, err := flavors.ListDetail(client, listOpts).AllPages()
		if err != nil {
			return fmterr.Errorf("Unable to query OpenTelekomCloud flavors: %w", err)
		}

		allFlavors, err = flavors.ExtractFlavors(allPages)
		if err != nil {
			return fmterr.Errorf("Unable to retrieve OpenTelekomCloud flavors: %w", err)
		}
	}

	// Loop through all flavors to find a more specific one.
	if len(allFlavors) > 0 {
		var filteredFlavors []flavors.Flavor
		for _, flavor := range allFlavors {
			if v := d.Get("name").(string); v != "" {
				if flavor.Name != v {
					continue
				}
			}

			// d.GetOk is used because 0 might be a valid choice.
			if v, ok := d.GetOk("ram"); ok {
				if flavor.RAM != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("vcpus"); ok {
				if flavor.VCPUs != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("disk"); ok {
				if flavor.Disk != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("swap"); ok {
				if flavor.Swap != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("rx_tx_factor"); ok {
				if flavor.RxTxFactor != v.(float64) {
					continue
				}
			}

			es, err := flavors.ListExtraSpecs(client, flavor.ID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			if v, ok := d.GetOk("resource_type"); ok {
				if resourceType, okType := es["resource_type"]; okType {
					if resourceType != v.(string) {
						continue
					}
				}
			}

			if v, ok := d.GetOk("availability_zone"); ok {
				if !checkAZExtraOption(es, v.(string)) {
					continue
				}
			}

			filteredFlavors = append(filteredFlavors, flavor)
		}

		allFlavors = filteredFlavors
	}

	if len(allFlavors) < 1 {
		return fmterr.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	flavor := allFlavors[0]

	log.Printf("[DEBUG] Retrieved opentelekomcloud_compute_flavor_v2 %s: %#v", flavor.ID, flavor)

	d.SetId(flavor.ID)
	mErr := multierror.Append(
		d.Set("name", flavor.Name),
		d.Set("flavor_id", flavor.ID),
		d.Set("disk", flavor.Disk),
		d.Set("ram", flavor.RAM),
		d.Set("rx_tx_factor", flavor.RxTxFactor),
		d.Set("swap", flavor.Swap),
		d.Set("vcpus", flavor.VCPUs),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	es, err := flavors.ListExtraSpecs(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("extra_specs", es); err != nil {
		log.Printf("[WARN] Unable to set extra_specs for opentelekomcloud_compute_flavor_v2 %s: %s", d.Id(), err)
	}

	return nil
}
