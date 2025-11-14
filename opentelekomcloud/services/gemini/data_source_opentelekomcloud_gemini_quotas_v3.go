package gemini

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/quota"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceGeminiQuotasV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiQuotasRead,

		Schema: map[string]*schema.Schema{
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used": {
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

func dataSourceGeminiQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	quotasResp, err := quota.GetQuotas(client)
	if err != nil {
		return diag.Errorf("error getting GeminiDB quotas: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	if err := d.Set("region", config.GetRegion(d)); err != nil {
		return diag.Errorf("error setting region: %s", err)
	}

	if err := d.Set("quotas", flattenQuotas(quotasResp.Quotas.Resources)); err != nil {
		return diag.Errorf("error setting quotas: %s", err)
	}

	return nil
}

func flattenQuotas(quotas []quota.ShowResourcesDetailResponseBody) []map[string]interface{} {
	result := make([]map[string]interface{}, len(quotas))
	for i, q := range quotas {
		result[i] = map[string]interface{}{
			"type":  q.Type,
			"quota": q.Quota,
			"used":  q.Used,
		}
	}
	return result
}
