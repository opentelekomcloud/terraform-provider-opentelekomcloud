package cc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/quota"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCcCentralNetworkQuotasV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworkQuotasV3Read,

		Schema: map[string]*schema.Schema{
			"quota_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCcCentralNetworkQuotasV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CcV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := quota.List(client, quota.ListOpts{
		QuotaType: common.ExpandToStringList(d.Get("quota_type").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud CC central network quotas: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("quotas", flattenCentralNetworkQuotas(resp.Quotas)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCentralNetworkQuotas(quotas []quota.CentralNetworkQuota) []map[string]interface{} {
	if len(quotas) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(quotas))
	for i, q := range quotas {
		result[i] = map[string]interface{}{
			"quota_key":   q.QuotaKey,
			"quota_limit": q.QuotaLimit,
			"used":        q.Used,
			"unit":        q.Unit,
		}
	}
	return result
}
