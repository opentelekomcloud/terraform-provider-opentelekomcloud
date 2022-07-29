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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/pools"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBPoolV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBPoolV2Create,
		ReadContext:   resourceLBPoolV2Read,
		UpdateContext: resourceLBPoolV2Update,
		DeleteContext: resourceLBPoolV2Delete,

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
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP",
				}, false),
			},

			// One of loadbalancer_id or listener_id must be provided
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// One of loadbalancer_id or listener_id must be provided
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP",
				}, false),
			},
			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE",
							}, false),
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourcePoolV2Persistence(d *schema.ResourceData) (*pools.SessionPersistence, error) {
	persistenceRaw := d.Get("persistence").([]interface{})

	var persistence pools.SessionPersistence
	if len(persistenceRaw) > 0 && persistenceRaw[0] != nil {
		persistenceInfo := persistenceRaw[0].(map[string]interface{})

		persistence = pools.SessionPersistence{
			Type: persistenceInfo["type"].(string),
		}
		if persistence.Type == "APP_COOKIE" {
			if persistenceInfo["cookie_name"].(string) == "" {
				return nil, fmt.Errorf("persistence cookie_name needs to be " +
					"set if using 'APP_COOKIE' persistence type")
			}
			persistence.CookieName = persistenceInfo["cookie_name"].(string)
		} else if persistenceInfo["cookie_name"].(string) != "" {
			return nil, fmt.Errorf("persistence cookie_name can only be " +
				"set if using 'APP_COOKIE' persistence type")
		}
		return &persistence, nil
	}
	return nil, nil
}

func resourceLBPoolV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)

	persistence, err := resourcePoolV2Persistence(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := pools.CreateOpts{
		TenantID:       d.Get("tenant_id").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Protocol:       pools.Protocol(d.Get("protocol").(string)),
		LoadbalancerID: d.Get("loadbalancer_id").(string),
		ListenerID:     d.Get("listener_id").(string),
		LBMethod:       pools.LBMethod(d.Get("lb_method").(string)),
		AdminStateUp:   &adminStateUp,
	}

	// Must omit if not set
	if persistence != nil {
		createOpts.Persistence = persistence
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	lbID := createOpts.LoadbalancerID
	listenerID := createOpts.ListenerID
	if lbID != "" {
		if err := waitForLBV2LoadBalancer(ctx, client, lbID, "ACTIVE", nil, timeout); err != nil {
			return diag.FromErr(err)
		}
	} else if listenerID != "" {
		// Wait for Listener to become active before continuing
		if err := waitForLBV2Listener(ctx, client, listenerID, "ACTIVE", nil, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Attempting to create pool")
	var pool *pools.Pool
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		pool, err = pools.Create(client, createOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error creating pool: %w", err)
	}

	// Wait for LoadBalancer to become active before continuing
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, client, lbID, "ACTIVE", nil, timeout)
	} else {
		// Pool exists by now so we can ask for lbID
		err = waitForLBV2viaPool(ctx, client, pool.ID, "ACTIVE", timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pool.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPoolV2Read(clientCtx, d, meta)
}

func resourceLBPoolV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	pool, err := pools.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "pool")
	}

	log.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	mErr := multierror.Append(
		d.Set("lb_method", pool.LBMethod),
		d.Set("protocol", pool.Protocol),
		d.Set("description", pool.Description),
		d.Set("tenant_id", pool.TenantID),
		d.Set("admin_state_up", pool.AdminStateUp),
		d.Set("name", pool.Name),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	log.Printf("[DEBUG] pool persistence: %+v", pool.Persistence)
	// UNDONE: This needs to check for failure, and it currently fails...
	// d.Set("persistence", pool.Persistence);

	return nil
}

func resourceLBPoolV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = pools.LBMethod(d.Get("lb_method").(string))
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, client, lbID, "ACTIVE", nil, timeout)
	} else {
		err = waitForLBV2viaPool(ctx, client, d.Id(), "ACTIVE", timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = pools.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to update pool %s: %w", d.Id(), err)
	}

	// Wait for LoadBalancer to become active before continuing
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, client, lbID, "ACTIVE", nil, timeout)
	} else {
		err = waitForLBV2viaPool(ctx, client, d.Id(), "ACTIVE", timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPoolV2Read(clientCtx, d, meta)
}

func resourceLBPoolV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutDelete)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		if err := waitForLBV2LoadBalancer(ctx, client, lbID, "ACTIVE", nil, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Attempting to delete pool %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = pools.Delete(client, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, client, lbID, "ACTIVE", nil, timeout)
	} else {
		// Wait for Pool to delete
		err = waitForLBV2Pool(ctx, client, d.Id(), "DELETED", nil, timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
