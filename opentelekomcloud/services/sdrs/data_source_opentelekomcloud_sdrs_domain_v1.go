package sdrs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/domains"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSdrsDomainV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSdrsDomainV1Read,

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

func dataSourceSdrsDomainV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	sdrsV1Client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	v, err := domains.Get(sdrsV1Client).Extract()
	if err != nil {
		return diag.FromErr(err)
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
		return fmterr.Errorf("your query returned no results. Please change your filters and try again.")
	}
	dm := filteredDomains[0]
	d.SetId(dm.Id)
	_ = d.Set("name", dm.Name)
	_ = d.Set("description", dm.Description)
	log.Printf("[DEBUG] SDRS Domain : %+v", dm)

	return nil
}
