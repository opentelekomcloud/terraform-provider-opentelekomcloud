package fw

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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceFWFirewallGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFWFirewallGroupV2Create,
		ReadContext:   resourceFWFirewallGroupV2Read,
		UpdateContext: resourceFWFirewallGroupV2Update,
		DeleteContext: resourceFWFirewallGroupV2Delete,
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
			"ingress_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"egress_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFWFirewallGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var createOpts firewall_groups.CreateOptsBuilder

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts = FirewallGroupCreateOpts{
		firewall_groups.CreateOpts{
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			IngressPolicyID: d.Get("ingress_policy_id").(string),
			EgressPolicyID:  d.Get("egress_policy_id").(string),
			AdminStateUp:    &adminStateUp,
			TenantID:        d.Get("tenant_id").(string),
		},
		common.MapValueSpecs(d),
	}

	portsRaw := d.Get("ports").(*schema.Set).List()
	if len(portsRaw) > 0 {
		log.Printf("[DEBUG] Will attempt to associate Firewall group with port(s): %+v", portsRaw)

		var portIds []string
		for _, v := range portsRaw {
			portIds = append(portIds, v.(string))
		}

		createOpts = &routerinsertion.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			PortIDs:           portIds,
		}
	}

	log.Printf("[DEBUG] Create firewall group: %#v", createOpts)

	firewall_group, err := firewall_groups.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Firewall group created: %#v", firewall_group)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE", "INACTIVE"},
		Refresh:    waitForFirewallGroupActive(networkingClient, firewall_group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for FW group to be active: %w", err)
	}
	log.Printf("[DEBUG] Firewall group (%s) is active.", firewall_group.ID)

	d.SetId(firewall_group.ID)

	return resourceFWFirewallGroupV2Read(ctx, d, meta)
}

func resourceFWFirewallGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about firewall: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var firewallGroup FirewallGroup
	err = firewall_groups.Get(networkingClient, d.Id()).ExtractInto(&firewallGroup)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "firewall")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud Firewall group %s: %#v", d.Id(), firewallGroup)

	mErr := multierror.Append(
		d.Set("name", firewallGroup.Name),
		d.Set("description", firewallGroup.Description),
		d.Set("ingress_policy_id", firewallGroup.IngressPolicyID),
		d.Set("egress_policy_id", firewallGroup.EgressPolicyID),
		d.Set("admin_state_up", firewallGroup.AdminStateUp),
		d.Set("tenant_id", firewallGroup.TenantID),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ports", firewallGroup.PortIDs); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving ports to state for OpenTelekomCloud firewall group (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceFWFirewallGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	// PolicyID is required
	opts := firewall_groups.UpdateOpts{
		IngressPolicyID: d.Get("ingress_policy_id").(string),
		EgressPolicyID:  d.Get("egress_policy_id").(string),
	}

	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &adminStateUp
	}

	var updateOpts firewall_groups.UpdateOptsBuilder
	var portIds []string
	if d.HasChange("ports") {
		portsRaw := d.Get("ports").(*schema.Set).List()
		log.Printf("[DEBUG] Will attempt to associate Firewall group with port(s): %+v", portsRaw)
		for _, v := range portsRaw {
			portIds = append(portIds, v.(string))
		}

		updateOpts = routerinsertion.UpdateOptsExt{
			UpdateOptsBuilder: opts,
			PortIDs:           portIds,
		}
	} else {
		updateOpts = opts
	}

	log.Printf("[DEBUG] Updating firewall with id %s: %#v", d.Id(), updateOpts)

	err = firewall_groups.Update(networkingClient, d.Id(), updateOpts).Err
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "INACTIVE"},
		Refresh:    waitForFirewallGroupActive(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for firewall group to become active: %w", err)
	}

	return resourceFWFirewallGroupV2Read(ctx, d, meta)
}

func resourceFWFirewallGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy firewall group: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	// Ensure the firewall group was fully created/updated before being deleted.
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "INACTIVE"},
		Refresh:    waitForFirewallGroupActive(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for FW group to be active: %w", err)
	}

	err = firewall_groups.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForFirewallGroupDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	return diag.FromErr(err)
}

func waitForFirewallGroupActive(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var fw FirewallGroup
		err := firewall_groups.Get(networkingClient, id).ExtractInto(&fw)
		if err != nil {
			return nil, "", err
		}
		return fw, fw.Status, nil
	}
}

func waitForFirewallGroupDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		fw, err := firewall_groups.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got firewall group %s => %#v", id, fw)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Firewall group %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] Firewall group %s deletion is pending", id)
		return fw, "DELETING", nil
	}
}
