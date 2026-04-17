package vpn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	cgw "github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/customer-gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceEnterpriseCustomerGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvpnCustomerGatewayRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id_type": {
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvpnCustomerGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	gatewayId := d.Get("id").(string)
	gw, err := cgw.Get(client, gatewayId)
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN customer gateway (%s): %s", gatewayId, err)
	}
	d.SetId(gw.ID)

	log.Printf("[DEBUG] Retrieved Enterprise VPN Gateway %s: %#v", d.Id(), gw)

	mErr := multierror.Append(
		nil,
		d.Set("id", d.Id()),
		d.Set("name", gw.Name),
		d.Set("asn", gw.BgpAsn),
		d.Set("id_value", gw.IdValue),
		d.Set("id_type", gw.IdType),
		d.Set("created_at", gw.CreatedAt),
		d.Set("updated_at", gw.UpdatedAt),
		d.Set("region", config.GetRegion(d)),
		d.Set("ip", gw.Ip),
		d.Set("route_mode", gw.RouteMode),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
