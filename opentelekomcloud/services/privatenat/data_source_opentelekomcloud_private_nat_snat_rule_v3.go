package privatenat

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/snatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourcePrivateNatSnatRuleV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateNatSnatRuleV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_rules": {
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
						"virsubnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_ip_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_ip_associations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"transit_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"transit_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
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

func dataSourcePrivateNatSnatRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var listOpts snatrules.ListSnatRulesQueryParams
	if v, ok := d.GetOk("id"); ok {
		listOpts.Id = []string{v.(string)}
		d.SetId(v.(string))
	}

	getResp, err := snatrules.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching private NAT SNAT rule : %w", err)
	}

	if d.Id() == "" {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(id)
	}

	natSnatRules := getResp.SnatRules

	mErr := multierror.Append(
		d.Set("snat_rules", setSnatRules(natSnatRules)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setSnatRules(snatRulesInResp []snatrules.PrivateSnat) []map[string]interface{} {
	var snatRules []map[string]interface{}
	for _, snatRuleInResp := range snatRulesInResp {
		transitIpIds, transitIpAssociations := processTransitIpAssociations(snatRuleInResp.TransitIpAssociations)
		snatRule := map[string]interface{}{
			"id":                      snatRuleInResp.Id,
			"gateway_id":              snatRuleInResp.GatewayId,
			"cidr":                    snatRuleInResp.Cidr,
			"virsubnet_id":            snatRuleInResp.VirSubnetId,
			"description":             snatRuleInResp.Description,
			"transit_ip_ids":          transitIpIds,
			"project_id":              snatRuleInResp.ProjectId,
			"transit_ip_associations": transitIpAssociations,
			"created_at":              snatRuleInResp.CreatedAt,
			"updated_at":              snatRuleInResp.UpdatedAt,
			"enterprise_project_id":   snatRuleInResp.EnterpriseProjectId,
			"status":                  snatRuleInResp.Status,
		}
		snatRules = append(snatRules, snatRule)
	}
	return snatRules
}
