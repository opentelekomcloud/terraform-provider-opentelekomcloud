package v2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/pools"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceL7PolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyV2Create,
		ReadContext:   resourceL7PolicyV2Read,
		UpdateContext: resourceL7PolicyV2Update,
		DeleteContext: resourceL7PolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceL7PolicyV2Import,
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

			"tenant_id": {
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

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"REDIRECT_TO_POOL", "REDIRECT_TO_LISTENER",
				}, true),
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"redirect_pool_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_listener_id"},
				Optional:      true,
			},

			"redirect_listener_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_pool_id"},
				Optional:      true,
			},

			"admin_state_up": {
				Type:         schema.TypeBool,
				Default:      true,
				Optional:     true,
				ValidateFunc: common.ValidateTrueOnly,
			},
		},
	}
}

func resourceL7PolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	// Assign some required variables for use in creation.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectListenerID := d.Get("redirect_listener_id").(string)

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectListenerID, redirectPoolID)
	if err != nil {
		return fmterr.Errorf("unable to create L7 Policy: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := l7policies.CreateOpts{
		TenantID:           d.Get("tenant_id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Action:             l7policies.Action(action),
		ListenerID:         listenerID,
		RedirectPoolID:     redirectPoolID,
		RedirectListenerID: redirectListenerID,
		AdminStateUp:       &adminStateUp,
	}

	if v, ok := d.GetOk("position"); ok {
		createOpts.Position = int32(v.(int))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Make sure the associated pool is active before proceeding.
	if redirectPoolID != "" {
		pool, err := pools.Get(client, redirectPoolID).Extract()
		if err != nil {
			return fmterr.Errorf("unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(ctx, client, pool.ID, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, client, parentListener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create L7 Policy")
	var l7Policy *l7policies.L7Policy
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		l7Policy, err = l7policies.Create(client, createOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error creating L7 Policy: %s", err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, client, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceL7PolicyV2Read(clientCtx, d, meta)
}

func resourceL7PolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	l7Policy, err := l7policies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	mErr := multierror.Append(
		d.Set("action", l7Policy.Action),
		d.Set("description", l7Policy.Description),
		d.Set("tenant_id", l7Policy.TenantID),
		d.Set("name", l7Policy.Name),
		d.Set("position", int(l7Policy.Position)),
		d.Set("redirect_listener_id", l7Policy.RedirectListenerID),
		d.Set("redirect_pool_id", l7Policy.RedirectPoolID),
		d.Set("region", config.GetRegion(d)),
		d.Set("admin_state_up", l7Policy.AdminStateUp),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7PolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	// Assign some required variables for use in updating.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectListenerID := d.Get("redirect_listener_id").(string)

	var updateOpts l7policies.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("redirect_pool_id") {
		redirectPoolID = d.Get("redirect_pool_id").(string)
		updateOpts.RedirectPoolID = &redirectPoolID
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectListenerID, redirectPoolID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Make sure the pool is active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)
	if redirectPoolID != "" {
		pool, err := pools.Get(client, redirectPoolID).Extract()
		if err != nil {
			return fmterr.Errorf("unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(ctx, client, pool.ID, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve parent listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve L7 Policy: %s: %s", d.Id(), err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, client, parentListener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, client, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = l7policies.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to update L7 Policy %s: %s", d.Id(), err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, client, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceL7PolicyV2Read(clientCtx, d, meta)
}

func resourceL7PolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the listener.
	listener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve parent listener (%s) for the L7 Policy: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Unable to retrieve L7 Policy")
	}

	// Wait for Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, client, listener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = l7policies.Delete(client, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting L7 Policy")
	}

	err = waitForLBV2L7Policy(ctx, client, listener, l7Policy, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7PolicyV2Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*cfg.Config)
	client, err := config.ElbV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf(ErrCreationV2Client, err)
	}

	l7Policy, err := l7policies.Get(client, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s during the import: %#v", d.Id(), l7Policy)

	var listenerID string
	if l7Policy.ListenerID != "" {
		listenerID = l7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		var err error
		listenerID, err = getListenerIDForL7Policy(client, d.Id())
		if err != nil {
			return nil, err
		}
	}
	if err := d.Set("listener_id", listenerID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func checkL7PolicyAction(action, redirectListenerID, redirectPoolID string) error {
	if action == "REDIRECT_TO_POOL" && redirectListenerID != "" {
		return fmt.Errorf("redirect_listener_id must be empty when action is set to %s", action)
	}

	if action == "REDIRECT_TO_LISTENER" && redirectPoolID != "" {
		return fmt.Errorf("redirect_pool_id must be empty when action is set to %s", action)
	}

	return nil
}
