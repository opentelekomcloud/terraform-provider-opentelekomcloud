package vpn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceEnterpriseVpnGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvpnGatewayRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_subnets": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"connect_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"er_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ha_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip1": {
				Type:     schema.TypeList,
				Elem:     DSGatewayEipSchema(),
				Computed: true,
			},
			"eip2": {
				Type:     schema.TypeList,
				Elem:     DSGatewayEipSchema(),
				Computed: true,
			},
			"access_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_private_ip_1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_private_ip_2": {
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
			"used_connection_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_connection_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DSGatewayEipSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceEvpnGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}
	gatewayId := d.Get("id").(string)
	gw, err := gateway.Get(client, gatewayId)
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN gateway (%s): %s", gatewayId, err)
	}
	d.SetId(gw.ID)

	log.Printf("[DEBUG] Retrieved Enterprise VPN Gateway %s: %#v", d.Id(), gw)

	mErr := multierror.Append(
		nil,
		d.Set("id", gw.ID),
		d.Set("region", config.GetRegion(d)),
		d.Set("attachment_type", gw.AttachmentType),
		d.Set("availability_zones", gw.AvailabilityZoneIds),
		d.Set("asn", gw.BgpAsn),
		d.Set("connect_subnet", gw.ConnectSubnet),
		d.Set("created_at", gw.CreatedAt),
		d.Set("flavor", gw.Flavor),
		d.Set("local_subnets", gw.LocalSubnets),
		d.Set("ha_mode", gw.HaMode),
		d.Set("eip1", flattenEvpGatewayResponseEip(gw.Eip1)),
		d.Set("name", gw.Name),
		d.Set("eip2", flattenEvpGatewayResponseEip(gw.Eip2)),
		d.Set("status", gw.Status),
		d.Set("updated_at", gw.UpdatedAt),
		d.Set("used_connection_group", gw.UsedConnectionGroup),
		d.Set("used_connection_number", gw.UsedConnectionNumber),
		d.Set("vpc_id", gw.VpcId),
		d.Set("access_vpc_id", gw.AccessVpcId),
		d.Set("access_subnet_id", gw.AccessSubnetId),
		d.Set("er_id", gw.ErId),
		d.Set("network_type", gw.NetworkType),
		d.Set("access_private_ip_1", gw.AccessPrivateIp1),
		d.Set("access_private_ip_2", gw.AccessPrivateIp2),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
