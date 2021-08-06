package bms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/flavors"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBMSFlavorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBMSFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"min_disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"swap": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"rx_tx_factor": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "id",
			},

			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "asc",
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, false),
			},
		},
	}
}

func dataSourceBMSFlavorV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	flavorClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekom bms client: %s", err)
	}

	listOpts := flavors.ListOpts{
		MinDisk: d.Get("min_disk").(int),
		MinRAM:  d.Get("min_ram").(int),
		Name:    d.Get("name").(string),
		ID:      d.Id(),
		SortKey: d.Get("sort_key").(string),
		SortDir: d.Get("sort_dir").(string),
	}
	var flavor flavors.Flavor
	refinedflavors, err := flavors.List(flavorClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve flavors: %s", err)
	}

	if len(refinedflavors) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	} else {
		flavor = refinedflavors[0]
	}

	log.Printf("[DEBUG] Single Flavor found: %s", flavor.ID)
	d.SetId(flavor.ID)
	mErr := multierror.Append(
		d.Set("name", flavor.Name),
		d.Set("disk", flavor.Disk),
		d.Set("min_disk", flavor.MinDisk),
		d.Set("min_ram", flavor.MinRAM),
		d.Set("ram", flavor.RAM),
		d.Set("rx_tx_factor", flavor.RxTxFactor),
		d.Set("swap", flavor.Swap),
		d.Set("vcpus", flavor.VCPUs),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
