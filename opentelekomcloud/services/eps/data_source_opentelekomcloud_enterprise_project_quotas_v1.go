package eps

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceEnterpriseProjectQuotasV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseProjectQuotasV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenResourceQuotas(resources []projects.Quota) []map[string]interface{} {
	if len(resources) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"quota": resource.Quota,
			"type":  resource.Type,
			"used":  resource.Used,
		})
	}

	return result
}

func dataSourceEnterpriseProjectQuotasV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(region)
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	qts, err := projects.ListQuotas(client)
	if err != nil {
		return diag.Errorf("error querying resource quotas: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenResourceQuotas(qts.Resources)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of the EPS resource quotas: %s", err)
	}
	return nil
}
