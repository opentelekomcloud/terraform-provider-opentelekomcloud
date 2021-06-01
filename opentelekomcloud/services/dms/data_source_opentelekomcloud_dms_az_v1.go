package dms

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/availablezones"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceDmsAZV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDmsAZV1Read,

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

func dataSourceDmsAZV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud dms client: %s", err)
	}

	v, err := availablezones.Get(DmsV1Client).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dms az : %+v", v)
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
			code := d.Get("code").(string)
			if code != "" && newAZ.Code != code {
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
	log.Printf("[DEBUG] Dms az : %+v", az)

	d.SetId(az.ID)
	d.Set("code", az.Code)
	d.Set("name", az.Name)
	d.Set("port", az.Port)

	return nil
}
