package nat

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/dnatrules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDnatRulesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnatRulesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_service_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_service_port": {
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
			"global_eip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_eip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_service_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"external_service_port": {
							Type:     schema.TypeInt,
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
						"internal_service_port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_service_port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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

func dataSourceDnatRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	listOpts := dnatrules.ListOpts{
		Id:                  d.Get("rule_id").(string),
		NatGatewayId:        d.Get("gateway_id").(string),
		PortId:              d.Get("port_id").(string),
		PrivateIp:           d.Get("private_ip").(string),
		InternalServicePort: d.Get("internal_service_port").(int),
		FloatingIpId:        d.Get("floating_ip_id").(string),
		FloatingIpAddress:   d.Get("floating_ip_address").(string),
		ExternalServicePort: d.Get("external_service_port").(int),
		Protocol:            d.Get("protocol").(string),
		Description:         d.Get("description").(string),
		Status:              d.Get("status").(string),
	}

	rules, err := dnatrules.List(client, listOpts)
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
		d.Set("rules", flattenListDnatRulesResponseBody(rules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDnatRulesResponseBody(rules []dnatrules.DnatRule) []interface{} {
	if rules == nil {
		return nil
	}

	dnatRules := make([]interface{}, len(rules))
	for _, rule := range rules {
		dnatRules = append(dnatRules, map[string]interface{}{
			"id":                          rule.ID,
			"gateway_id":                  rule.NatGatewayId,
			"protocol":                    rule.Protocol,
			"port_id":                     rule.PortId,
			"private_ip":                  rule.PrivateIp,
			"internal_service_port":       rule.InternalServicePort,
			"external_service_port":       rule.ExternalServicePort,
			"floating_ip_id":              rule.FloatingIpId,
			"floating_ip_address":         rule.FloatingIpAddress,
			"internal_service_port_range": rule.InternalServicePort,
			"external_service_port_range": rule.ExternalServicePort,
			"description":                 rule.Description,
			"status":                      rule.Status,
			"created_at":                  rule.CreatedAt,
		})
	}
	return dnatRules
}
