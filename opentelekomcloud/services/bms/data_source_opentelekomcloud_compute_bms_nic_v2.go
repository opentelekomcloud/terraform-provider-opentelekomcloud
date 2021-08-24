package bms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/nics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBMSNicV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBMSNicV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBMSNicV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	nicClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating compute v2 client: %w", err)
	}

	listOpts := nics.ListOpts{
		ID:     d.Id(),
		Status: d.Get("status").(string),
	}

	refinedNics, err := nics.List(nicClient, d.Get("server_id").(string), listOpts)
	log.Printf("[DEBUG] Nic info: %#v", refinedNics)
	if err != nil {
		return fmterr.Errorf("unable to retrieve nics: %s", err)
	}

	if len(refinedNics) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNics) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Nic := refinedNics[0]

	var s []map[string]interface{}
	for _, fixedips := range Nic.FixedIP {
		mapping := map[string]interface{}{
			"subnet_id":  fixedips.SubnetID,
			"ip_address": fixedips.IPAddress,
		}
		s = append(s, mapping)
	}

	log.Printf("[INFO] Retrieved Nic using given filter %s: %+v", Nic.ID, Nic)
	d.SetId(Nic.ID)

	mErr := multierror.Append(
		d.Set("status", Nic.Status),
		d.Set("network_id", Nic.NetworkID),
		d.Set("mac_address", Nic.MACAddress),
		d.Set("region", config.GetRegion(d)),
		d.Set("fixed_ips", s),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
