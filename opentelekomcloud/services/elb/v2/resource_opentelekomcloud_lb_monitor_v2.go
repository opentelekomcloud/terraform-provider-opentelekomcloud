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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/monitors"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceMonitorV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorV2Create,
		ReadContext:   resourceMonitorV2Read,
		UpdateContext: resourceMonitorV2Update,
		DeleteContext: resourceMonitorV2Delete,

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
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP_CONNECT", "HTTP",
				}, false),
			},
			"delay": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "OPTIONS", "CONNECT", "PATCH",
				}, false),
			},
			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"monitor_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsPortNumber,
			},
		},
	}
}

func resourceMonitorV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := monitors.CreateOpts{
		PoolID:        d.Get("pool_id").(string),
		Type:          d.Get("type").(string),
		Delay:         d.Get("delay").(int),
		DomainName:    d.Get("domain_name").(string),
		Timeout:       d.Get("timeout").(int),
		MaxRetries:    d.Get("max_retries").(int),
		URLPath:       d.Get("url_path").(string),
		HTTPMethod:    d.Get("http_method").(string),
		ExpectedCodes: d.Get("expected_codes").(string),
		TenantID:      d.Get("tenant_id").(string),
		Name:          d.Get("name").(string),
		AdminStateUp:  &adminStateUp,
		MonitorPort:   d.Get("monitor_port").(int),
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	poolID := createOpts.PoolID
	// Wait for parent pool to become active before continuing
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	log.Printf("[DEBUG] Attempting to create monitor")
	var monitor *monitors.Monitor
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		monitor, err = monitors.Create(client, createOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to create monitor: %s", err)
	}

	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(monitor.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceMonitorV2Read(clientCtx, d, meta)
}

func resourceMonitorV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	monitor, err := monitors.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "monitor")
	}

	log.Printf("[DEBUG] Retrieved monitor %s: %#v", d.Id(), monitor)

	mErr := multierror.Append(nil,
		d.Set("tenant_id", monitor.TenantID),
		d.Set("type", monitor.Type),
		d.Set("delay", monitor.Delay),
		d.Set("timeout", monitor.Timeout),
		d.Set("max_retries", monitor.MaxRetries),
		d.Set("url_path", monitor.URLPath),
		d.Set("http_method", monitor.HTTPMethod),
		d.Set("expected_codes", monitor.ExpectedCodes),
		d.Set("admin_state_up", monitor.AdminStateUp),
		d.Set("name", monitor.Name),
		d.Set("monitor_port", monitor.MonitorPort),
		d.Set("region", config.GetRegion(d)),
		d.Set("domain_name", monitor.DomainName),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceMonitorV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	var updateOpts monitors.UpdateOpts
	if d.HasChange("url_path") {
		updateOpts.URLPath = d.Get("url_path").(string)
	}
	if d.HasChange("expected_codes") {
		updateOpts.ExpectedCodes = d.Get("expected_codes").(string)
	}
	if d.HasChange("delay") {
		updateOpts.Delay = d.Get("delay").(int)
	}
	if d.HasChange("domain_name") {
		updateOpts.DomainName = d.Get("domain_name").(string)
	}
	if d.HasChange("timeout") {
		updateOpts.Timeout = d.Get("timeout").(int)
	}
	if d.HasChange("max_retries") {
		updateOpts.MaxRetries = d.Get("max_retries").(int)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("http_method") {
		updateOpts.HTTPMethod = d.Get("http_method").(string)
	}
	if d.HasChange("monitor_port") {
		updateOpts.MonitorPort = d.Get("monitor_port").(int)
	}

	log.Printf("[DEBUG] Updating monitor %s with options: %#v", d.Id(), updateOpts)
	timeout := d.Timeout(schema.TimeoutUpdate)
	poolID := d.Get("pool_id").(string)
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = monitors.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to update monitor %s: %s", d.Id(), err)
	}

	// Wait for LB to become active before continuing
	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceMonitorV2Read(clientCtx, d, meta)
}

func resourceMonitorV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	log.Printf("[DEBUG] Deleting monitor %s", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	poolID := d.Get("pool_id").(string)

	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = monitors.Delete(client, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to delete monitor %s: %s", d.Id(), err)
	}

	if err := waitForLBV2viaPool(ctx, client, poolID, "ACTIVE", timeout); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
