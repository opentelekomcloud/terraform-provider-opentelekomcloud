package dns

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/resolverrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDNSResolverRulesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSResolverRulesV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resolver_rules": {
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
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceDNSResolverRulesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	resp, err := resolverrules.List(client, resolverrules.ListOpts{
		Name:       d.Get("name").(string),
		DomainName: d.Get("domain_name").(string),
		EndpointID: d.Get("endpoint_id").(string),
		Limit:      500,
	})
	if err != nil {
		return fmterr.Errorf("error listing DNS resolver rules: %w", err)
	}

	ids := make([]string, 0, len(resp.ResolverRules))
	result := make([]map[string]interface{}, 0, len(resp.ResolverRules))
	for _, rule := range resp.ResolverRules {
		vpcIDs := make([]string, 0, len(rule.Routers))
		for _, r := range rule.Routers {
			vpcIDs = append(vpcIDs, r.RouterID)
		}
		ids = append(ids, rule.ID)
		result = append(result, map[string]interface{}{
			"id":               rule.ID,
			"name":             rule.Name,
			"domain_name":      rule.DomainName,
			"endpoint_id":      rule.EndpointID,
			"status":           rule.Status,
			"ip_address_count": rule.IPAddressCount,
			"vpc_ids":          vpcIDs,
			"created_at":       rule.CreateTime,
			"updated_at":       rule.UpdateTime,
		})
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("resolver_rules", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
