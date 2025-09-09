package privatenat

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/transitip"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourcePrivateNatTransitIpV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateNatTransitIpV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virsubnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virsubnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
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
						"network_interface_id": {
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
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePrivateNatTransitIpV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var listOpts transitip.ListTransitIpsQueryParams
	if v, ok := d.GetOk("id"); ok {
		listOpts.Id = []string{v.(string)}
		d.SetId(v.(string))
	}
	if v, ok := d.GetOk("ip_address"); ok {
		listOpts.IpAddress = []string{v.(string)}
		d.SetId(v.(string))
	}
	if v, ok := d.GetOk("virsubnet_id"); ok {
		listOpts.VirSubnetID = []string{v.(string)}
	}

	getResp, err := transitip.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching private NAT Transit IPs : %w", err)
	}

	if d.Id() == "" {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(id)
	}

	natTransitIps := getResp.TransitIps

	mErr := multierror.Append(
		d.Set("id", d.Id()),
		d.Set("transit_ips", setTransitIps(natTransitIps)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setTransitIps(transitIpsInResp []transitip.TransitIP) []map[string]interface{} {
	var transitIps []map[string]interface{}
	for _, transitIpInResp := range transitIpsInResp {
		transitIp := map[string]interface{}{
			"id":                    transitIpInResp.Id,
			"virsubnet_id":          transitIpInResp.VirSubnetID,
			"ip_address":            transitIpInResp.IpAddress,
			"tags":                  setTransitIpTags(transitIpInResp.Tags),
			"enterprise_project_id": transitIpInResp.EnterpriseProjectID,
			"project_id":            transitIpInResp.ProjectId,
			"status":                transitIpInResp.Status,
			"created_at":            transitIpInResp.CreatedAt,
			"updated_at":            transitIpInResp.UpdatedAt,
			"network_interface_id":  transitIpInResp.NetworkInterfaceId,
			"gateway_id":            transitIpInResp.GatewayId,
		}
		transitIps = append(transitIps, transitIp)
	}
	return transitIps
}
