package privatenat

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/natgateway"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourcePrivateNatGatewayV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateNatGatewayV3Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"downlink_vpcs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"virsubnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ngport_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
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
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"transit_ip_pool_size_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePrivateNatGatewayV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var listOpts natgateway.ListGatewaysQueryParams
	if v, ok := d.GetOk("id"); ok {
		listOpts.Id = []string{v.(string)}
		d.SetId(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = []string{v.(string)}
	}

	getResp, err := natgateway.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching private nat gateways : %w", err)
	}

	if d.Id() == "" {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(id)
	}

	natGateways := getResp.Gateways

	mErr := multierror.Append(
		d.Set("gateways", setGateways(natGateways)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setGateways(gatewaysInResp []natgateway.PrivateNATGateway) []map[string]interface{} {
	var gateways []map[string]interface{}
	for _, gatewayInResp := range gatewaysInResp {
		gateway := map[string]interface{}{
			"id":                       gatewayInResp.Id,
			"name":                     gatewayInResp.Name,
			"description":              gatewayInResp.Description,
			"spec":                     gatewayInResp.Spec,
			"downlink_vpcs":            setDownlinkVpcs(gatewayInResp.DownlinkVpcs),
			"tags":                     setGatewayTags(gatewayInResp.Tags),
			"enterprise_project_id":    gatewayInResp.EnterpriseProjectID,
			"project_id":               gatewayInResp.ProjectId,
			"status":                   gatewayInResp.Status,
			"created_at":               gatewayInResp.CreatedAt,
			"updated_at":               gatewayInResp.UpdatedAt,
			"rule_max":                 gatewayInResp.RuleMax,
			"transit_ip_pool_size_max": gatewayInResp.TransitIpPoolSizeMax,
		}
		gateways = append(gateways, gateway)
	}
	return gateways
}
