package vpn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/connection"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceEnterpriseConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvpnConnectionRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_subnets": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"tunnel_local_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_peer_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_nqa": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy_rules": {
				Type:     schema.TypeList,
				Elem:     DataSourceConnectionPolicyRuleSchema(),
				Computed: true,
			},
			"ikepolicy": {
				Type:     schema.TypeList,
				Elem:     DataSourceConnectionIkePolicySchema(),
				Computed: true,
			},
			"ipsecpolicy": {
				Type:     schema.TypeList,
				Elem:     DataSourceConnectionIpsecPolicySchema(),
				Computed: true,
			},
			"ha_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
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
		},
	}
}

func DataSourceConnectionIkePolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"local_id_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_id_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phase_one_negotiation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dh_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dpd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceConnectionPolicyDPDSchema(),
			},
		},
	}
	return &sc
}

func DataSourceConnectionPolicyDPDSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func DataSourceConnectionIpsecPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pfs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transform_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encapsulation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func DataSourceConnectionPolicyRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceEvpnConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	connectionId := d.Get("id").(string)
	gw, err := connection.Get(client, connectionId)
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN connections (%s): %s", d.Id(), err)
	}
	d.SetId(gw.ID)
	log.Printf("[DEBUG] Retrieved Enterprise VPN Gateway %s: %#v", d.Id(), gw)

	tagsMap := make(map[string]string)
	for _, tag := range gw.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(
		nil,
		d.Set("id", d.Id()),
		d.Set("region", config.GetRegion(d)),
		d.Set("name", gw.Name),
		d.Set("gateway_id", gw.VgwId),
		d.Set("gateway_ip", gw.VgwIp),
		d.Set("vpn_type", gw.Style),
		d.Set("customer_gateway_id", gw.CgwId),
		d.Set("peer_subnets", gw.PeerSubnets),
		d.Set("tunnel_local_address", gw.TunnelLocalAddress),
		d.Set("tunnel_peer_address", gw.TunnelPeerAddress),
		d.Set("enable_nqa", gw.EnableNqa),
		d.Set("ha_role", gw.HaRole),
		d.Set("created_at", gw.CreatedAt),
		d.Set("updated_at", gw.UpdatedAt),
		d.Set("status", gw.Status),
		d.Set("tags", tagsMap),
		d.Set("ikepolicy", flattenConnectionIkePolicy(gw.IkePolicy)),
		d.Set("ipsecpolicy", flattenConnectionIpSecPolicy(gw.IpSecPolicy)),
		d.Set("policy_rules", flattenConnectionPolicyRule(gw.PolicyRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
