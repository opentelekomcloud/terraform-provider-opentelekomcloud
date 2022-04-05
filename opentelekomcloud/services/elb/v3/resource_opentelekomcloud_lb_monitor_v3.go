package v3

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/monitors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceMonitorV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorV3Create,
		ReadContext:   resourceMonitorV3Read,
		UpdateContext: resourceMonitorV3Update,
		DeleteContext: resourceMonitorV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"TCP", "UDP_CONNECT", "HTTP", "HTTPS", "PING"},
					false,
				),
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"delay": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"max_retries": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"max_retries_down": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 0xff),
					validation.StringMatch(regexp.MustCompile(`^/.*$`), "value must start with slash"),
				),
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(
						regexp.MustCompile(`^[0-9a-zA-Z][0-9a-zA-Z.-]*$`),
						"The value can contain only digits, letters, hyphens (-), and periods (.) and must start with a digit or letter.",
					),
				),
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "OPTIONS", "CONNECT", "PATCH"},
					false,
				),
			},
			"expected_codes": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 0xff),
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

func resourceMonitorV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := monitors.CreateOpts{
		AdminStateUp:   iBool(d.Get("admin_state_up")),
		PoolID:         d.Get("pool_id").(string),
		Type:           monitors.Type(d.Get("type").(string)),
		Delay:          d.Get("delay").(int),
		Timeout:        d.Get("timeout").(int),
		MaxRetries:     d.Get("max_retries").(int),
		MaxRetriesDown: d.Get("max_retries_down").(int),
		URLPath:        d.Get("url_path").(string),
		DomainName:     d.Get("domain_name").(string),
		HTTPMethod:     d.Get("http_method").(string),
		ExpectedCodes:  d.Get("expected_codes").(string),
		ProjectID:      d.Get("project_id").(string),
		Name:           d.Get("name").(string),
		MonitorPort:    d.Get("monitor_port").(int),
	}
	monitor, err := monitors.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating LB monitor v3: %w", err)
	}
	d.SetId(monitor.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceMonitorV3Read(clientCtx, d, meta)
}

func resourceMonitorV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	monitor, err := monitors.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error reading LB monitor v3")
	}

	mErr := multierror.Append(
		d.Set("admin_state_up", monitor.AdminStateUp),
		d.Set("delay", monitor.Delay),
		d.Set("domain_name", monitor.DomainName),
		d.Set("expected_codes", monitor.ExpectedCodes),
		d.Set("http_method", monitor.HTTPMethod),
		d.Set("max_retries", monitor.MaxRetries),
		d.Set("max_retries_down", monitor.MaxRetriesDown),
		d.Set("monitor_port", monitor.MonitorPort),
		d.Set("name", monitor.Name),
		d.Set("project_id", monitor.ProjectID),
		d.Set("timeout", monitor.Timeout),
		d.Set("type", monitor.Type),
		d.Set("url_path", monitor.URLPath),
	)
	if len(monitor.Pools) > 0 {
		mErr = multierror.Append(mErr, d.Set("pool_id", monitor.Pools[0].ID))
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB monitor v3 fields: %w", err)
	}

	return nil
}

func resourceMonitorV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := monitors.UpdateOpts{
		AdminStateUp:   iBool(d.Get("admin_state_up")),
		Type:           d.Get("type").(string),
		Delay:          d.Get("delay").(int),
		Timeout:        d.Get("timeout").(int),
		MaxRetries:     d.Get("max_retries").(int),
		MaxRetriesDown: d.Get("max_retries_down").(int),
		URLPath:        d.Get("url_path").(string),
		DomainName:     d.Get("domain_name").(string),
		HTTPMethod:     d.Get("http_method").(string),
		ExpectedCodes:  d.Get("expected_codes").(string),
		Name:           d.Get("name").(string),
		MonitorPort:    d.Get("monitor_port").(int),
	}
	_, err = monitors.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating LB monitor v3: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceMonitorV3Read(clientCtx, d, meta)
}

func resourceMonitorV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	err = monitors.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting LB monitor v3")
	}

	return nil
}
