package dns

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/endpoints"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDNSEndpointsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSEndpointsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"inbound", "outbound"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resolver_rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDNSEndpointsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	resp, err := endpoints.List(client, endpoints.ListOpts{
		Direction: d.Get("direction").(string),
		VpcID:     d.Get("vpc_id").(string),
		Name:      d.Get("name").(string),
		Limit:     500,
	})
	if err != nil {
		return fmterr.Errorf("error listing DNS endpoints: %w", err)
	}

	ids := make([]string, 0, len(resp.Endpoints))
	result := make([]map[string]interface{}, 0, len(resp.Endpoints))
	for _, ep := range resp.Endpoints {
		ids = append(ids, ep.ID)
		result = append(result, map[string]interface{}{
			"id":                  ep.ID,
			"name":                ep.Name,
			"direction":           ep.Direction,
			"status":              ep.Status,
			"vpc_id":              ep.VpcID,
			"ip_address_count":    ep.IPAddressCount,
			"resolver_rule_count": ep.ResolverRuleCount,
			"created_at":          ep.CreateTime,
			"updated_at":          ep.UpdateTime,
		})
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("endpoints", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
