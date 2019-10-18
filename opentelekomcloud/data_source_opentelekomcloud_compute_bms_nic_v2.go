package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/nics"
)

func dataSourceBMSNicV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBMSNicV2Read,

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

func dataSourceBMSNicV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nicClient, err := config.computeV2HWClient(GetRegion(d, config))

	listOpts := nics.ListOpts{
		ID:     d.Id(),
		Status: d.Get("status").(string),
	}

	refinedNics, err := nics.List(nicClient, d.Get("server_id").(string), listOpts)
	log.Printf("[DEBUG] Nic info: %#v", refinedNics)
	if err != nil {
		return fmt.Errorf("Unable to retrieve nics: %s", err)
	}

	if len(refinedNics) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNics) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
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

	d.Set("status", Nic.Status)
	d.Set("network_id", Nic.NetworkID)
	d.Set("mac_address", Nic.MACAddress)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("fixed_ips", s); err != nil {
		return err
	}

	return nil
}
