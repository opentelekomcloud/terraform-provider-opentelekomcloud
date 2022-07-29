package v2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLoadBalancerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV2Create,
		ReadContext:   resourceLoadBalancerV2Read,
		UpdateContext: resourceLoadBalancerV2Update,
		DeleteContext: resourceLoadBalancerV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"vip_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state_up": {
				Type:         schema.TypeBool,
				Default:      true,
				Optional:     true,
				ValidateFunc: common.ValidateTrueOnly,
			},
			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceLoadBalancerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	var lbProvider string
	if v, ok := d.GetOk("loadbalancer_provider"); ok {
		lbProvider = v.(string)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := loadbalancers.CreateOpts{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		VipSubnetID:  d.Get("vip_subnet_id").(string),
		TenantID:     d.Get("tenant_id").(string),
		VipAddress:   d.Get("vip_address").(string),
		AdminStateUp: &adminStateUp,
		Provider:     lbProvider,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	lb, err := loadbalancers.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating LoadBalancer: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(ctx, client, lb.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "loadbalancers", lb.ID, tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of LoadBalancer: %s", err)
		}
	}

	// If all has been successful, set the ID on the resource
	d.SetId(lb.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLoadBalancerV2Read(clientCtx, d, meta)
}

func resourceLoadBalancerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	lb, err := loadbalancers.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "loadbalancer")
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

	mErr := multierror.Append(nil,
		d.Set("name", lb.Name),
		d.Set("description", lb.Description),
		d.Set("vip_subnet_id", lb.VipSubnetID),
		d.Set("tenant_id", lb.TenantID),
		d.Set("vip_address", lb.VipAddress),
		d.Set("vip_port_id", lb.VipPortID),
		d.Set("admin_state_up", lb.AdminStateUp),
		d.Set("loadbalancer_provider", lb.Provider),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	// save tags
	resourceTags, err := tags.Get(client, "loadbalancers", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud LoadCalancer tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud LoadCalancer: %s", err)
	}

	return nil
}

func resourceLoadBalancerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	var updateOpts loadbalancers.UpdateOpts
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
	if err := waitForLBV2LoadBalancer(ctx, client, d.Id(), "ACTIVE", nil, timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating loadbalancer %s with options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = loadbalancers.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to update loadbalancer %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active before continuing
	if err := waitForLBV2LoadBalancer(ctx, client, d.Id(), "ACTIVE", nil, timeout); err != nil {
		return diag.FromErr(err)
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "loadbalancers", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of LoadBalancer %s: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLoadBalancerV2Read(clientCtx, d, meta)
}

func resourceLoadBalancerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = loadbalancers.Delete(client, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to delete loadbalancer %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become delete
	pending := []string{"PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}
	if err := waitForLBV2LoadBalancer(ctx, client, d.Id(), "DELETED", pending, timeout); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
