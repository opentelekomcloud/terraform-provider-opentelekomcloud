package er

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	az "github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/availability-zones"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceErAvailabilityZonesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceErAvailabilityZonesV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceErAvailabilityZonesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := az.List(client, az.ListOpts{})
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud ER v3 availability zones: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("names", flattenListAvailabilityZone(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAvailabilityZone(azs []az.AZsResponse) []interface{} {
	if len(azs) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(azs))
	for _, item := range azs {
		if item.State == "available" {
			result = append(result, item.Code)
		}
	}
	return result
}
