package dcs

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/availablezones"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceDcsAZV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsAZV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsAZV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	DcsV1Client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating dcs key client: %s", err)
	}

	v, err := availablezones.Get(DcsV1Client).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dcs az : %+v", v)
	var filteredAZs []availablezones.AvailableZone
	if v.RegionID == config.GetRegion(d) {
		AZs := v.AvailableZones
		for _, newAZ := range AZs {
			if newAZ.ResourceAvailability != "true" {
				continue
			}

			name := d.Get("name").(string)
			if name != "" && newAZ.Name != name {
				continue
			}

			port := d.Get("port").(string)
			if port != "" && newAZ.Port != port {
				continue
			}
			filteredAZs = append(filteredAZs, newAZ)
		}
	}

	if len(filteredAZs) < 1 {
		return fmt.Errorf("Not found any available zones")
	}

	az := filteredAZs[0]
	log.Printf("[DEBUG] Dcs az : %+v", az)

	d.SetId(az.ID)
	d.Set("code", az.Code)
	d.Set("name", az.Name)
	d.Set("port", az.Port)

	return nil
}
