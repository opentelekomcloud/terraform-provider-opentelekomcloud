package taurusdb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/quota"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlProjectQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBMysqlProjectQuotasRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"used": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"quota": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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

func dataSourceTaurusDBMysqlProjectQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	opts := quota.ListProjectQuotasOpts{}
	if v, ok := d.GetOk("type"); ok {
		typeStr := v.(string)
		opts.Type = &typeStr
	}

	result, err := quota.ListProjectQuotas(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB project quotas: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	if err := d.Set("region", config.GetRegion(d)); err != nil {
		return diag.FromErr(err)
	}

	quotasList := make([]map[string]interface{}, 1)
	resourcesList := make([]map[string]interface{}, len(result.Quotas.Resources))

	for i, resource := range result.Quotas.Resources {
		resourcesList[i] = map[string]interface{}{
			"type":  resource.Type,
			"used":  resource.Used,
			"quota": resource.Quota,
		}
	}

	quotasList[0] = map[string]interface{}{
		"resources": resourcesList,
	}

	if err := d.Set("quotas", quotasList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
