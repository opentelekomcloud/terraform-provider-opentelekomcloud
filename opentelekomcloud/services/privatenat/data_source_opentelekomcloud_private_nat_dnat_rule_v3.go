package privatenat

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/dnatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourcePrivateNatDnatRuleV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateNatDnatRuleV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dnat_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_service_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_service_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePrivateNatDnatRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var listOpts dnatrules.ListDnatRulesQueryParams
	if v, ok := d.GetOk("id"); ok {
		listOpts.Id = []string{v.(string)}
		d.SetId(v.(string))
	}

	getResp, err := dnatrules.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching private NAT DNAT rule : %w", err)
	}

	if d.Id() == "" {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(id)
	}

	natDnatRules := getResp.DnatRules

	mErr := multierror.Append(
		d.Set("dnat_rules", setDnatRules(natDnatRules)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setDnatRules(dnatRulesInResp []dnatrules.PrivateDnat) []map[string]interface{} {
	var dnatRules []map[string]interface{}
	for _, dnatRuleInResp := range dnatRulesInResp {
		dnatRule := map[string]interface{}{
			"id":                    dnatRuleInResp.Id,
			"description":           dnatRuleInResp.Description,
			"transit_ip_id":         dnatRuleInResp.TransitIpId,
			"network_interface_id":  dnatRuleInResp.NetworkInterfaceId,
			"gateway_id":            dnatRuleInResp.GatewayId,
			"private_ip_address":    dnatRuleInResp.PrivateIpAddress,
			"protocol":              dnatRuleInResp.Protocol,
			"internal_service_port": dnatRuleInResp.InternalServicePort,
			"transit_service_port":  dnatRuleInResp.TransitServicePort,
			"project_id":            dnatRuleInResp.ProjectId,
			"type":                  dnatRuleInResp.Type,
			"enterprise_project_id": dnatRuleInResp.EnterpriseProjectId,
			"created_at":            dnatRuleInResp.CreatedAt,
			"updated_at":            dnatRuleInResp.UpdatedAt,
			"status":                dnatRuleInResp.Status,
		}
		dnatRules = append(dnatRules, dnatRule)
	}
	return dnatRules
}
