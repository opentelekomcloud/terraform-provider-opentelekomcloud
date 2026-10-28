package v3

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/monitors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLBMonitorV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBMonitorV3Read,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 0xff),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"TCP", "UDP_CONNECT", "HTTP", "HTTPS", "PING"},
					false,
				),
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"max_retries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"max_retries_down": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"monitor_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsPortNumber,
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
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 0xff),
					validation.StringMatch(regexp.MustCompile(`^/.*$`), "value must start with slash"),
				),
			},
			"expected_codes": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMonitorV3ListOpts(d *schema.ResourceData) monitors.ListOpts {
	filter := func(key string) string {
		if !common.IsAttrSet(d, key) {
			return ""
		}
		return d.Get(key).(string)
	}
	opts := monitors.ListOpts{
		Name: filter("name"), Type: filter("type"), URLPath: filter("url_path"),
		ExpectedCodes: filter("expected_codes"),
	}
	if common.IsAttrSet(d, "delay") {
		opts.Delay = d.Get("delay").(int)
	}
	if common.IsAttrSet(d, "timeout") {
		opts.Timeout = d.Get("timeout").(int)
	}
	if common.IsAttrSet(d, "max_retries") {
		opts.MaxRetries = d.Get("max_retries").(int)
	}
	if common.IsAttrSet(d, "max_retries_down") {
		opts.MaxRetriesDown = []int{d.Get("max_retries_down").(int)}
	}
	if common.IsAttrSet(d, "monitor_port") {
		opts.MonitorPort = []int{d.Get("monitor_port").(int)}
	}
	if common.IsAttrSet(d, "domain_name") {
		opts.DomainName = []string{d.Get("domain_name").(string)}
	}
	if common.IsAttrSet(d, "admin_state_up") {
		value := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &value
	}
	return opts
}

func dataSourceLBMonitorV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}
	var monitor *monitors.Monitor
	if common.IsAttrSet(d, "id") {
		monitor, err = monitors.Get(client, d.Get("id").(string))
	} else {
		var matches []monitors.Monitor
		matches, err = monitors.List(client, buildMonitorV3ListOpts(d))
		if err == nil {
			if len(matches) == 0 {
				return common.DataSourceTooFewDiag
			}
			if len(matches) > 1 {
				return common.DataSourceTooManyDiag
			}
			monitor = &matches[0]
		}
	}
	if err != nil {
		return fmterr.Errorf("error retrieving LB monitor v3: %w", err)
	}
	if monitor.ID == "" {
		return diag.Errorf("error retrieving LB monitor v3: response does not contain a monitor ID")
	}
	// Unlike the resource's required pool, a data-source pool is computed and must not retain stale state.
	if err := d.Set("pool_id", ""); err != nil {
		return diag.FromErr(err)
	}
	if diagnostics := setMonitorV3Fields(d, monitor); diagnostics.HasError() {
		return diagnostics
	}
	d.SetId(monitor.ID)
	return nil
}
