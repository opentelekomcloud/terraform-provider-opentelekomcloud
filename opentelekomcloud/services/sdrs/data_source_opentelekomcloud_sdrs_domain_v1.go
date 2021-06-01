package sdrs

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/domains"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceSdrsDomainV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSdrsDomainV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSdrsDomainV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	sdrsV1Client, err := config.SdrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating SDRS client: %s", err)
	}

	v, err := domains.Get(sdrsV1Client).Extract()
	if err != nil {
		return err
	}

	adomains := v.Domains
	var filteredDomains []domains.Domain
	for _, dm := range adomains {
		name := d.Get("name").(string)
		if name != "" && dm.Name != name {
			continue
		}
		filteredDomains = append(filteredDomains, dm)
	}
	if len(filteredDomains) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	dm := filteredDomains[0]
	d.SetId(dm.Id)
	d.Set("name", dm.Name)
	d.Set("description", dm.Description)
	log.Printf("[DEBUG] SDRS Domain : %+v", dm)

	return nil
}
