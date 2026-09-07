package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

var vpcQuotaTypes = []string{
	"vpc",
	"subnet",
	"securityGroup",
	"securityGroupRule",
	"publicIp",
	"vpn",
	"vpcPeer",
	"loadbalancer",
	"listener",
	"physicalConnect",
	"virtualInterface",
	"firewall",
	"shareBandwidthIP",
	"shareBandwidth",
	"address_group",
	"flow_log",
	"vpcContainRoutetable",
	"routetableContainRoutes",
}

func DataSourceVpcQuotasV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcQuotasV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(vpcQuotaTypes, false),
			},
			"quotas": {
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
						"min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcQuotasV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	quotaType := d.Get("type").(string)
	allQuotas, err := quotas.List(client, quotas.ListOpts{Type: quotaType})
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud VPC v1 quotas: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", client.ProjectID, quotaType))
	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("quotas", flattenVpcQuotas(allQuotas)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVpcQuotas(allQuotas []quotas.Quota) []map[string]interface{} {
	if len(allQuotas) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(allQuotas))
	for _, quota := range allQuotas {
		result = append(result, map[string]interface{}{
			"type":  quota.Type,
			"used":  quota.Used,
			"quota": quota.Quota,
			"min":   quota.Min,
		})
	}
	return result
}
