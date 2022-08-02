package v2

import (
	"context"
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

func ResourceMemberV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberV2Create,
		ReadContext:   resourceMemberV2Read,
		UpdateContext: resourceMemberV2Update,
		DeleteContext: resourceMemberV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("pool_id", "id"),
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
				Default:      1,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMemberV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := pools.CreateMemberOpts{
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
		Name:         d.Get("name").(string),
		TenantID:     d.Get("tenant_id").(string),
		Weight:       golangsdk.IntToPointer(d.Get("weight").(int)),
		SubnetID:     d.Get("subnet_id").(string),
		AdminStateUp: &adminStateUp,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LB to become active before continuing
	poolID := d.Get("pool_id").(string)
	timeout := d.Timeout(schema.TimeoutCreate)
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create member")
	var member *pools.Member
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		member, err = pools.CreateMember(client, poolID, createOpts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault429); ok {
				time.Sleep(1 * time.Minute)
				return resource.RetryableError(err)
			}
			time.Sleep(5 * time.Second)
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error creating member: %w", err)
	}

	// Wait for LB to become ACTIVE again
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(member.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceMemberV2Read(clientCtx, d, meta)
}

func resourceMemberV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	poolID := d.Get("pool_id").(string)
	member, err := pools.GetMember(client, poolID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "member")
	}

	log.Printf("[DEBUG] Retrieved member %s: %#v", d.Id(), member)

	mErr := multierror.Append(
		d.Set("name", member.Name),
		d.Set("weight", member.Weight),
		d.Set("admin_state_up", member.AdminStateUp),
		d.Set("tenant_id", member.TenantID),
		d.Set("subnet_id", member.SubnetID),
		d.Set("address", member.Address),
		d.Set("protocol_port", member.ProtocolPort),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMemberV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	var updateOpts pools.UpdateMemberOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("weight") {
		updateOpts.Weight = golangsdk.IntToPointer(d.Get("weight").(int))
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Wait for LB to become active before continuing
	poolID := d.Get("pool_id").(string)
	timeout := d.Timeout(schema.TimeoutUpdate)
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating member %s with options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = pools.UpdateMember(client, poolID, d.Id(), updateOpts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault429); ok {
				time.Sleep(1 * time.Minute)
				return resource.RetryableError(err)
			}
			return common.CheckForRetryableError(err)
		}
		time.Sleep(5 * time.Second)
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to update member %s: %w", d.Id(), err)
	}

	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceMemberV2Read(clientCtx, d, meta)
}

func resourceMemberV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	// Wait for Pool to become active before continuing
	poolID := d.Get("pool_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to delete member %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = pools.DeleteMember(client, poolID, d.Id()).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault429); ok {
				time.Sleep(1 * time.Minute)
				return resource.RetryableError(err)
			}
			return common.CheckForRetryableError(err)
		}
		time.Sleep(5 * time.Second)
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to delete member %s: %w", d.Id(), err)
	}

	// Wait for LB to become ACTIVE
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
