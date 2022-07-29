package v2

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
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
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	enableWhitelist := d.Get("enable_whitelist").(bool)
	createOpts := whitelists.CreateOpts{
		TenantId:        d.Get("tenant_id").(string),
		ListenerId:      d.Get("listener_id").(string),
		EnableWhitelist: &enableWhitelist,
		Whitelist:       d.Get("whitelist").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	wl, err := whitelists.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Whitelist: %w", err)
	}

	d.SetId(wl.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceWhitelistV2Read(clientCtx, d, meta)
}

func resourceWhitelistV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	wl, err := whitelists.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "whitelist")
	}

	log.Printf("[DEBUG] Retrieved whitelist %s: %#v", d.Id(), wl)

	mErr := multierror.Append(
		d.Set("tenant_id", wl.TenantId),
		d.Set("listener_id", wl.ListenerId),
		d.Set("enable_whitelist", wl.EnableWhitelist),
		d.Set("whitelist", wl.Whitelist),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWhitelistV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
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
	_, err = whitelists.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("unable to update whitelist %s: %w", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceWhitelistV2Read(clientCtx, d, meta)
}

func resourceWhitelistV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	log.Printf("[DEBUG] Attempting to delete whitelist %s", d.Id())

	if err := whitelists.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud whitelist: %w", err)
	}
	d.SetId("")

	return nil
}
