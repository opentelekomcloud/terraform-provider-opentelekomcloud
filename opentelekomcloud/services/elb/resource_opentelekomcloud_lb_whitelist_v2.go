package elb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/whitelists"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWhitelistV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWhitelistV2Create,
		ReadContext:   resourceWhitelistV2Read,
		UpdateContext: resourceWhitelistV2Update,
		DeleteContext: resourceWhitelistV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enable_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"whitelist": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: common.SuppressLBWhitelistDiffs,
			},
		},
	}
}

func resourceWhitelistV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	enableWhitelist := d.Get("enable_whitelist").(bool)
	createOpts := whitelists.CreateOpts{
		TenantId:        d.Get("tenant_id").(string),
		ListenerId:      d.Get("listener_id").(string),
		EnableWhitelist: &enableWhitelist,
		Whitelist:       d.Get("whitelist").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	wl, err := whitelists.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Whitelist: %s", err)
	}

	d.SetId(wl.ID)
	return resourceWhitelistV2Read(ctx, d, meta)
}

func resourceWhitelistV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	wl, err := whitelists.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "whitelist"))
	}

	log.Printf("[DEBUG] Retrieved whitelist %s: %#v", d.Id(), wl)

	d.SetId(wl.ID)
	d.Set("tenant_id", wl.TenantId)
	d.Set("listener_id", wl.ListenerId)
	d.Set("enable_whitelist", wl.EnableWhitelist)
	d.Set("whitelist", wl.Whitelist)

	return nil
}

func resourceWhitelistV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts whitelists.UpdateOpts
	if d.HasChange("enable_whitelist") {
		ew := d.Get("enable_whitelist").(bool)
		updateOpts.EnableWhitelist = &ew
	}
	if d.HasChange("whitelist") {
		updateOpts.Whitelist = d.Get("whitelist").(string)
	}

	log.Printf("[DEBUG] Updating whitelist %s with options: %#v", d.Id(), updateOpts)
	_, err = whitelists.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("Unable to update whitelist %s: %s", d.Id(), err)
	}

	return resourceWhitelistV2Read(ctx, d, meta)
}

func resourceWhitelistV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete whitelist %s", d.Id())
	err = whitelists.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud whitelist: %s", err)
	}
	d.SetId("")
	return nil
}
