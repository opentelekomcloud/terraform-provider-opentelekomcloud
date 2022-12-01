package vpn

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/siteconnections"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

var (
	// A PSK can contain 6 to 128 characters.
	// Spaces and the following special characters are not allowed: <>&?*'"
	pskRegex = regexp.MustCompile("^[^<>&?*'\"]{6,128}$")
)

func ResourceVpnSiteConnectionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnSiteConnectionV2Create,
		ReadContext:   resourceVpnSiteConnectionV2Read,
		UpdateContext: resourceVpnSiteConnectionV2Update,
		DeleteContext: resourceVpnSiteConnectionV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ikepolicy_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_ep_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpnservice_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"local_ep_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipsecpolicy_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"psk": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(pskRegex, "Invalid pks value. "+
					"A PSK can contain 6 to 128 characters. Spaces and the following special characters are not allowed: <>&?*'\"."),
			},
			"initiator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"peer_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dpd": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceVpnSiteConnectionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var createOpts siteconnections.CreateOptsBuilder

	dpd := resourceSiteConnectionV2DPDCreateOpts(d.Get("dpd").(*schema.Set))

	peerCidrsRaw := d.Get("peer_cidrs").([]interface{})
	peerCidrs := make([]string, len(peerCidrsRaw))
	for i, peerCids := range peerCidrsRaw {
		peerCidrs[i] = peerCids.(string)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	initiator := resourceSiteConnectionV2Initiator(d.Get("initiator").(string))

	createOpts = VpnSiteConnectionCreateOpts{
		siteconnections.CreateOpts{
			Name:           d.Get("name").(string),
			Description:    d.Get("description").(string),
			AdminStateUp:   &adminStateUp,
			Initiator:      initiator,
			IKEPolicyID:    d.Get("ikepolicy_id").(string),
			TenantID:       d.Get("tenant_id").(string),
			PeerID:         d.Get("peer_id").(string),
			PeerAddress:    d.Get("peer_address").(string),
			PeerEPGroupID:  d.Get("peer_ep_group_id").(string),
			LocalID:        d.Get("local_id").(string),
			VPNServiceID:   d.Get("vpnservice_id").(string),
			LocalEPGroupID: d.Get("local_ep_group_id").(string),
			IPSecPolicyID:  d.Get("ipsecpolicy_id").(string),
			PSK:            d.Get("psk").(string),
			MTU:            d.Get("mtu").(int),
			PeerCIDRs:      peerCidrs,
			DPD:            &dpd,
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create site connection: %#v", createOpts)

	conn, err := siteconnections.Create(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"PENDING_CREATE"},
		Refresh:    waitForSiteConnectionCreation(client, conn.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] SiteConnection created: %#v", conn)

	d.SetId(conn.ID)

	// create tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "ipsec-site-connections", d.Id(), tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of VPN site connection %s: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVpnSiteConnectionV2Read(clientCtx, d, meta)
}

func resourceVpnSiteConnectionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	conn, err := siteconnections.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "site_connection")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud SiteConnection %s: %#v", d.Id(), conn)

	dpd := []map[string]interface{}{{
		"action":   conn.DPD.Action,
		"interval": conn.DPD.Interval,
		"timeout":  conn.DPD.Timeout,
	}}

	mErr := multierror.Append(
		d.Set("name", conn.Name),
		d.Set("description", conn.Description),
		d.Set("admin_state_up", conn.AdminStateUp),
		d.Set("tenant_id", conn.TenantID),
		d.Set("initiator", conn.Initiator),
		d.Set("ikepolicy_id", conn.IKEPolicyID),
		d.Set("peer_id", conn.PeerID),
		d.Set("peer_address", conn.PeerAddress),
		d.Set("local_id", conn.LocalID),
		d.Set("peer_ep_group_id", conn.PeerEPGroupID),
		d.Set("vpnservice_id", conn.VPNServiceID),
		d.Set("local_ep_group_id", conn.LocalEPGroupID),
		d.Set("ipsecpolicy_id", conn.IPSecPolicyID),
		// Do not set psk here as the response value is not same with the requested
		// d.Set("psk", conn.PSK)
		d.Set("mtu", conn.MTU),
		d.Set("peer_cidrs", conn.PeerCIDRs),
		d.Set("dpd", &dpd),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// Set tags
	resourceTags, err := tags.Get(client, "ipsec-site-connections", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching VPN site connection tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for VPN site connection %s: %s", d.Id(), err)
	}

	return nil
}

func resourceVpnSiteConnectionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := siteconnections.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &adminStateUp
		hasChange = true
	}

	if d.HasChange("local_id") {
		opts.LocalID = d.Get("local_id").(string)
		hasChange = true
	}

	if d.HasChange("peer_address") {
		opts.PeerAddress = d.Get("peer_address").(string)
		hasChange = true
	}

	if d.HasChange("peer_id") {
		opts.PeerID = d.Get("peer_id").(string)
		hasChange = true
	}

	if d.HasChange("local_ep_group_id") {
		opts.LocalEPGroupID = d.Get("local_ep_group_id").(string)
		hasChange = true
	}

	if d.HasChange("peer_ep_group_id") {
		opts.PeerEPGroupID = d.Get("peer_ep_group_id").(string)
		hasChange = true
	}

	if d.HasChange("psk") {
		opts.PSK = d.Get("psk").(string)
		hasChange = true
	}

	if d.HasChange("mtu") {
		opts.MTU = d.Get("mtu").(int)
		hasChange = true
	}

	if d.HasChange("initiator") {
		initiator := resourceSiteConnectionV2Initiator(d.Get("initiator").(string))
		opts.Initiator = initiator
		hasChange = true
	}

	if d.HasChange("peer_cidrs") {
		opts.PeerCIDRs = d.Get("peer_cidrs").([]string)
		hasChange = true
	}

	if d.HasChange("dpd") {
		dpdUpdateOpts := resourceSiteConnectionV2DPDUpdateOpts(d.Get("dpd").(*schema.Set))
		opts.DPD = &dpdUpdateOpts
		hasChange = true
	}
	log.Printf("[DEBUG] Updating site connection with id %s: %#v", d.Id(), opts)

	if hasChange {
		conn, err := siteconnections.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForSiteConnectionUpdate(client, conn.ID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated connection with id %s", d.Id())
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "ipsec-site-connections", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of VPN site connection %s: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVpnSiteConnectionV2Read(clientCtx, d, meta)
}

func resourceVpnSiteConnectionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if err := siteconnections.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSiteConnectionDeletion(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSiteConnectionDeletion(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := siteconnections.Get(client, id).Extract()
		log.Printf("[DEBUG] Got site connection %s => %#v", id, conn)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] SiteConnection %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] SiteConnection %s deletion is pending", id)
		return conn, "DELETING", nil
	}
}

func waitForSiteConnectionCreation(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := siteconnections.Get(client, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return service, "PENDING_CREATE", nil
	}
}

func waitForSiteConnectionUpdate(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := siteconnections.Get(client, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return conn, "UPDATED", nil
	}
}

func resourceSiteConnectionV2Initiator(initiatorString string) siteconnections.Initiator {
	var ini siteconnections.Initiator
	switch initiatorString {
	case "bi-directional":
		ini = siteconnections.InitiatorBiDirectional
	case "response-only":
		ini = siteconnections.InitiatorResponseOnly
	}
	return ini
}

func resourceSiteConnectionV2DPDCreateOpts(d *schema.Set) siteconnections.DPDCreateOpts {
	dpd := siteconnections.DPDCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		dpd.Action = resourceSiteConnectionV2Action(rawMap["action"].(string))

		timeout := rawMap["timeout"].(int)
		dpd.Timeout = timeout

		interval := rawMap["interval"].(int)
		dpd.Interval = interval
	}
	return dpd
}
func resourceSiteConnectionV2Action(actionString string) siteconnections.Action {
	var act siteconnections.Action
	switch actionString {
	case "hold":
		act = siteconnections.ActionHold
	case "restart":
		act = siteconnections.ActionRestart
	case "disabled":
		act = siteconnections.ActionDisabled
	case "restart-by-peer":
		act = siteconnections.ActionRestartByPeer
	case "clear":
		act = siteconnections.ActionClear
	}
	return act
}

func resourceSiteConnectionV2DPDUpdateOpts(d *schema.Set) siteconnections.DPDUpdateOpts {
	dpd := siteconnections.DPDUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		dpd.Action = resourceSiteConnectionV2Action(rawMap["action"].(string))

		timeout := rawMap["timeout"].(int)
		dpd.Timeout = timeout

		interval := rawMap["interval"].(int)
		dpd.Interval = interval
	}
	return dpd
}
