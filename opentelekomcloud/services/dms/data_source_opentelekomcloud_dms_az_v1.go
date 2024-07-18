package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/availablezones"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDmsAZV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsAZV1Read,

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

func dataSourceDmsAZV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms client: %s", err)
	}

	v, err := availablezones.Get(DmsV1Client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Dms az : %+v", v)
	var filteredAZs []availablezones.AvailableZone
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

	if len(filteredAZs) < 1 {
		return fmterr.Errorf("not found any available zones")
	}

	az := filteredAZs[0]
	log.Printf("[DEBUG] Dms az : %+v", az)

	d.SetId(az.ID)
	mErr := multierror.Append(
		d.Set("code", az.Code),
		d.Set("name", az.Name),
		d.Set("port", az.Port),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
