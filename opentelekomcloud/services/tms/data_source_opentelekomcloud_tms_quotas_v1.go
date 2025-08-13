package tms

import (
	"context"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/quota"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceTmsQuotasV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsQuotasV1Read,

		Schema: map[string]*schema.Schema{
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
		},
	}
}

func dataSourceTmsQuotasV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	allQuotas, err := quota.List(client)
	if err != nil {
		return fmterr.Errorf("Error listing TMS tag quotas: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	var quotaList []map[string]interface{}
	for _, t := range allQuotas {
		quota := map[string]interface{}{
			"quota_key":   t.QuotaKey,
			"quota_limit": t.QuotaLimit,
			"used":        t.Used,
			"unit":        t.Unit,
		}
		quotaList = append(quotaList, quota)
	}

	if err = d.Set("quotas", quotaList); err != nil {
		return fmterr.Errorf("Error setting TMS tag quotas: %s", err)
	}
	return nil

}
