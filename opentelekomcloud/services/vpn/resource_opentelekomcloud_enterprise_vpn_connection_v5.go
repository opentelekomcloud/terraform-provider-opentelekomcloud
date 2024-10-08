package vpn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/connection"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEnterpriseConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvpnConnectionCreate,
		UpdateContext: resourceEvpnConnectionUpdate,
		ReadContext:   resourceEvpnConnectionRead,
		DeleteContext: resourceEvpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: common.SuppressCaseDiffs,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_subnets": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"tunnel_local_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tunnel_peer_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_nqa": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"psk": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_rules": {
				Type:     schema.TypeList,
				Elem:     ConnectionPolicyRuleSchema(),
				Optional: true,
				Computed: true,
			},
			"ikepolicy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ConnectionIkePolicySchema(),
				Optional: true,
				Computed: true,
			},
			"ipsecpolicy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ConnectionIpsecPolicySchema(),
				Optional: true,
				Computed: true,
			},
			"ha_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func ConnectionIkePolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ike_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"local_id_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peer_id_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"phase_one_negotiation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"authentication_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dh_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dpd": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     ConnectionPolicyDPDSchema(),
			},
		},
	}
	return &sc
}

func ConnectionPolicyDPDSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"msg": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func ConnectionIpsecPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pfs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"transform_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encapsulation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func ConnectionPolicyRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_index": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceEvpnConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	connectionTags := d.Get("tags").(map[string]interface{})
	var tagSlice []tags.ResourceTag
	for k, v := range connectionTags {
		tagSlice = append(tagSlice, tags.ResourceTag{Key: k, Value: v.(string)})
	}
	createOpts := connection.CreateOpts{
		Name:               d.Get("name").(string),
		VgwId:              d.Get("gateway_id").(string),
		VgwIp:              d.Get("gateway_ip").(string),
		Style:              d.Get("vpn_type").(string),
		CgwId:              d.Get("customer_gateway_id").(string),
		PeerSubnets:        buildPeerSubnets(d),
		TunnelLocalAddress: d.Get("tunnel_local_address").(string),
		TunnelPeerAddress:  d.Get("tunnel_peer_address").(string),
		Psk:                d.Get("psk").(string),
		PolicyRules:        buildConnectionPolicyRules(d),
		IkePolicy:          buildConnectionIkePolicy(d),
		IpSecPolicy:        buildConnectionIpSecPolicy(d),
		HaRole:             d.Get("ha_role").(string),
		Tags:               tagSlice,
	}

	if enableNqa, ok := d.GetOk("enable_nqa"); ok {
		createOpts.EnableNqa = pointerto.Bool(enableNqa.(bool))
	}

	n, err := connection.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud EVPN connection: %w", err)
	}

	d.SetId(n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"ACTIVE"},
		Refresh:      waitForConnectionActive(client, n.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
		MinTimeout:   3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for OpenTelekomCloud EVPN connection (%s) to become ACTIVE: %w", n.ID, err)
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnConnectionRead(clientCtx, d, meta)
}

func buildPeerSubnets(d *schema.ResourceData) []string {
	var peerSubnets []string
	subRaw := d.Get("peer_subnets").([]interface{})
	for _, s := range subRaw {
		peerSubnets = append(peerSubnets, s.(string))
	}
	return peerSubnets
}

func buildConnectionPolicyRules(d *schema.ResourceData) []connection.PolicyRules {
	rawRules := d.Get("policy_rules").([]interface{})
	if len(rawRules) == 0 {
		return nil
	}
	rules := make([]connection.PolicyRules, len(rawRules))
	for i, raw := range rawRules {
		if rawMap, ok := raw.(map[string]interface{}); ok {
			var dest []string
			for _, s := range rawMap["destination"].([]interface{}) {
				dest = append(dest, s.(string))
			}
			rules[i] = connection.PolicyRules{
				RuleIndex:   rawMap["rule_index"].(int),
				Source:      rawMap["source"].(string),
				Destination: dest,
			}
		}
	}
	return rules
}

func buildConnectionIpSecPolicy(d *schema.ResourceData) *connection.IpSecPolicy {
	rawPolicy := d.Get("ipsecpolicy").([]interface{})
	if len(rawPolicy) == 0 {
		return nil
	}
	if raw, ok := rawPolicy[0].(map[string]interface{}); ok {
		policy := connection.IpSecPolicy{
			AuthenticationAlgorithm: raw["authentication_algorithm"].(string),
			EncryptionAlgorithm:     raw["encryption_algorithm"].(string),
			Pfs:                     raw["pfs"].(string),
			LifetimeSeconds:         pointerto.Int(raw["lifetime_seconds"].(int)),
			TransformProtocol:       raw["transform_protocol"].(string),
			EncapsulationMode:       raw["encapsulation_mode"].(string),
		}
		return &policy
	}
	return nil
}

func buildConnectionIkePolicy(d *schema.ResourceData) *connection.IkePolicy {
	rawPolicy := d.Get("ikepolicy").([]interface{})
	if len(rawPolicy) == 0 {
		return nil
	}

	if raw, ok := rawPolicy[0].(map[string]interface{}); ok {
		params := connection.IkePolicy{
			AuthenticationAlgorithm: raw["authentication_algorithm"].(string),
			EncryptionAlgorithm:     raw["encryption_algorithm"].(string),
			IkeVersion:              raw["ike_version"].(string),
			LifetimeSeconds:         pointerto.Int(raw["lifetime_seconds"].(int)),
			LocalIdType:             raw["local_id_type"].(string),
			PeerIdType:              raw["peer_id_type"].(string),
			PhaseOneNegotiationMode: raw["phase_one_negotiation_mode"].(string),
			AuthenticationMethod:    raw["authentication_method"].(string),
			DhGroup:                 raw["dh_group"].(string),
			Dpd:                     buildConnectionDPD(raw["dpd"]),
		}
		if raw["local_id_type"].(string) != "ip" {
			params.LocalId = raw["local_id"].(string)
		}
		if raw["peer_id_type"].(string) != "ip" {
			params.PeerId = raw["peer_id"].(string)
		}
		return &params
	}
	return nil
}

func buildConnectionDPD(dpd interface{}) *connection.Dpd {
	rawDpd := dpd.([]interface{})
	if len(rawDpd) == 0 {
		return nil
	}
	if raw, ok := rawDpd[0].(map[string]interface{}); ok {
		policy := connection.Dpd{
			Timeout:  pointerto.Int(raw["timeout"].(int)),
			Interval: pointerto.Int(raw["interval"].(int)),
			Msg:      raw["msg"].(string),
		}
		return &policy
	}
	return nil
}

func resourceEvpnConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	gw, err := connection.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN connections (%s): %s", d.Id(), err)
	}

	tagsMap := make(map[string]string)
	for _, tag := range gw.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(
		nil,
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

func flattenConnectionIkePolicy(resp connection.IkePolicy) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"authentication_algorithm":   resp.AuthenticationAlgorithm,
			"encryption_algorithm":       resp.EncryptionAlgorithm,
			"ike_version":                resp.IkeVersion,
			"lifetime_seconds":           resp.LifetimeSeconds,
			"local_id_type":              resp.LocalIdType,
			"local_id":                   resp.LocalId,
			"peer_id_type":               resp.PeerIdType,
			"peer_id":                    resp.PeerId,
			"phase_one_negotiation_mode": resp.PhaseOneNegotiationMode,
			"authentication_method":      resp.AuthenticationMethod,
			"dh_group":                   resp.DhGroup,
			"dpd":                        flattenConnectionDPD(resp.Dpd),
		},
	}
	return rst
}

func flattenConnectionDPD(resp *connection.Dpd) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"timeout":  resp.Timeout,
			"interval": resp.Interval,
			"msg":      resp.Msg,
		},
	}
	return rst
}

func flattenConnectionIpSecPolicy(resp connection.IpSecPolicy) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"authentication_algorithm": resp.AuthenticationAlgorithm,
			"encryption_algorithm":     resp.EncryptionAlgorithm,
			"pfs":                      resp.Pfs,
			"lifetime_seconds":         resp.LifetimeSeconds,
			"transform_protocol":       resp.TransformProtocol,
			"encapsulation_mode":       resp.EncapsulationMode,
		},
	}
	return rst
}

func flattenConnectionPolicyRule(resp []connection.PolicyRules) []interface{} {
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"rule_index":  v.RuleIndex,
			"destination": v.Destination,
			"source":      v.Source,
		})
	}
	return rst
}

func resourceEvpnConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	updateConnectionhasChanges := []string{
		"customer_gateway_id",
		"enable_nqa",
		"ikepolicy",
		"ipsecpolicy",
		"name",
		"peer_subnets",
		"policy_rules",
		"psk",
		"tunnel_local_address",
		"tunnel_peer_address",
	}

	if d.HasChanges(updateConnectionhasChanges...) {
		opts := connection.UpdateOpts{
			ConnectionID:       d.Id(),
			Name:               d.Get("name").(string),
			CgwId:              d.Get("customer_gateway_id").(string),
			PeerSubnets:        buildPeerSubnets(d),
			TunnelLocalAddress: d.Get("tunnel_local_address").(string),
			TunnelPeerAddress:  d.Get("tunnel_peer_address").(string),
			Psk:                d.Get("psk").(string),
			PolicyRules:        buildConnectionPolicyRules(d),
			IkePolicy:          buildConnectionIkePolicy(d),
			IpSecPolicy:        buildConnectionIpSecPolicy(d),
		}
		_, err = connection.Update(client, opts)
		if err != nil {
			return diag.Errorf("error updating OpenTelekomCloud EVPN connection: %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"CREATING"},
			Target:  []string{"ACTIVE"},
			Refresh: waitForConnectionActive(client, d.Id()),
			Timeout: d.Timeout(schema.TimeoutUpdate),
			Delay:   10 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for updating OpenTelekomCloud EVPN connection (%s) to complete: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err = updateTags(client, d, "vpn-connection", d.Id()); err != nil {
			return diag.Errorf("error updating tags of OpenTelekomCloud EVPN connection (%s): %s", d.Id(), err)
		}
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnConnectionRead(clientCtx, d, meta)
}

func resourceEvpnConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	err = connection.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud EVPN connection")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: waitForConnectionDeletion(client, d.Id()),
		Timeout: d.Timeout(schema.TimeoutDelete),
		Delay:   10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for deleting OpenTelekomCloud EVPN connection (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func waitForConnectionActive(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := connection.Get(client, id)
		if err != nil {
			return nil, "", err
		}
		if n.Status == "ACTIVE" || n.Status == "DOWN" {
			return n, "ACTIVE", nil
		}
		return n, "CREATING", nil
	}
}

func waitForConnectionDeletion(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := connection.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The OpenTelekomCloud EVPN connection has been deleted (ID:%s).", id)
				return r, "DELETED", nil
			}
			return nil, "ERROR", err
		}
		switch r.Status {
		case "ACTIVE", "PENDING_DELETE":
			return r, "DELETING", nil
		default:
			err = fmt.Errorf("error deleting OpenTelekomCloud EVPN connection [%s]. "+
				"Unexpected status: %v", r.ID, r.Status)
			return r, "ERROR", err
		}
	}
}
