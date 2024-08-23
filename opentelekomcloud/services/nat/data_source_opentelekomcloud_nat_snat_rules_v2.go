package nat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/snatrules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSnatRulesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSnatRulesRead,
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"floating_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"admin_state_up": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSnatRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	listOpts := snatrules.ListOpts{
		NetworkId:         d.Get("subnet_id").(string),
		Cidr:              d.Get("cidr").(string),
		SourceType:        strconv.Itoa(d.Get("source_type").(int)),
		Id:                d.Get("rule_id").(string),
		NatGatewayId:      d.Get("gateway_id").(string),
		ProjectId:         d.Get("project_id").(string),
		FloatingIpId:      d.Get("floating_ip_id").(string),
		FloatingIpAddress: d.Get("floating_ip_address").(string),
		Description:       d.Get("description").(string),
		Status:            d.Get("status").(string),
	}

	rules, err := snatrules.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve NAT gateway pages: %w", err)
	}

	uID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uID)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("rules", flattenListSnatRulesResponseBody(rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListSnatRulesResponseBody(rules []snatrules.SnatRule) []interface{} {
	if rules == nil {
		return nil
	}

	var snatRules []interface{}
	for _, rule := range rules {
		source, err := getSourceType(rule.SourceType)
		if err != nil {
			return nil
		}
		snatRules = append(snatRules, map[string]interface{}{
			"id":                  rule.ID,
			"gateway_id":          rule.NatGatewayID,
			"subnet_id":           rule.NetworkID,
			"project_id":          rule.TenantID,
			"floating_ip_id":      rule.FloatingIPID,
			"floating_ip_address": rule.FloatingIPAddress,
			"status":              rule.Status,
			"created_at":          rule.CreatedAt,
			"source_type":         source,
			"admin_state_up":      rule.AdminStateUp,
			"cidr":                rule.Cidr,
		})
	}
	return snatRules
}

func getSourceType(s interface{}) (int, error) {
	switch v := s.(type) {
	case float64:
		return int(v), nil
	case string:
		sourceType, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("error converting `source_type` from string: %w", err)
		}
		return sourceType, nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("unsupported type for `source_type`: %T", v)
	}
}
