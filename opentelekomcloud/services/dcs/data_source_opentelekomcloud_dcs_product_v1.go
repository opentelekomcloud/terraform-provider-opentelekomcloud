package dcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/others"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDcsProductV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsProductV1Read,

		Schema: map[string]*schema.Schema{
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsProductV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DcsV1Client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error get dcs product client: %s", err)
	}

	v, err := others.GetProducts(DcsV1Client)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Dcs get products : %+v", v)
	var FilteredPd []others.Product
	for _, pd := range v {
		specCode := d.Get("spec_code").(string)
		if specCode != "" && pd.SpecCode != specCode {
			continue
		}
		FilteredPd = append(FilteredPd, pd)
	}

	if len(FilteredPd) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your filters and try again.")
	}

	pd := FilteredPd[0]
	d.SetId(pd.ProductID)
	_ = d.Set("spec_code", pd.SpecCode)
	log.Printf("[DEBUG] Dcs product : %+v", pd)

	return nil
}
