package dcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/products"

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

	v, err := products.Get(DcsV1Client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Dcs get products : %+v", v)
	var FilteredPd []products.Product
	for _, pd := range v.Products {
		spec_code := d.Get("spec_code").(string)
		if spec_code != "" && pd.SpecCode != spec_code {
			continue
		}
		FilteredPd = append(FilteredPd, pd)
	}

	if len(FilteredPd) < 1 {
		return fmterr.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	pd := FilteredPd[0]
	d.SetId(pd.ProductID)
	d.Set("spec_code", pd.SpecCode)
	log.Printf("[DEBUG] Dcs product : %+v", pd)

	return nil
}
