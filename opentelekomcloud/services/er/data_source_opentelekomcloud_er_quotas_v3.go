package er

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/quota"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceErQuotasV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceErQuotasReadV3,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceErQuotasReadV3(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := quota.List(client, quota.ListOpts{})
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud ER v3 quotas: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("quotas", flattenListQuotas(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListQuotas(qts []quota.QuotaResponse) []interface{} {
	if len(qts) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(qts))
	for _, item := range qts {
		result = append(result, map[string]interface{}{
			"used":      int(item.UsedQuota),
			"unit":      item.Unit,
			"type":      item.Type,
			"available": int(item.AvailableQuota),
		})
	}
	return result
}
