package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
)

func ResourceWafDedicatedPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedPolicyV1Create,
		ReadContext:   resourceWafDedicatedPolicyV1Read,
		UpdateContext: resourceWafDedicatedPolicyV1Update,
		DeleteContext: resourceWafDedicatedPolicyV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protection_mode": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"log", "block",
				}, false),
			},
			"level": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web_attack": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"common": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"anti_crawler": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_engine": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_script": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_scanner": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"crawler_other": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"web_shell": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"cc": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"custom": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"blacklist": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"geolocation_access_control": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"ignore": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"privacy": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"anti_tamper": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"anti_leakage": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"followed_action": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"bot_enable": {
							Type:     schema.TypeBool,
							Computed: true,
							// Optional: true,
						},
						"precise": {
							Type:     schema.TypeBool,
							Computed: true,
							// Optional: true,
						},
					},
				},
			},
			"full_detection": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceWafDedicatedPolicyV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	r, err := policies.Create(client, policies.CreateOpts{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)

	if d.HasChanges("protection_mode", "level", "options", "full_detection") {
		if err := updateWafPolicy(ctx, d, meta); err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud WAF dedicated policy: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedPolicyV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedPolicyV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policy, err := policies.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud WAF dedicated policy.")
	}

	options := []map[string]interface{}{
		{
			"web_attack":                 policy.Options.WebAttack,
			"common":                     policy.Options.Common,
			"crawler":                    policy.Options.Crawler,
			"anti_crawler":               policy.Options.AntiCrawler,
			"crawler_engine":             policy.Options.CrawlerEngine,
			"crawler_scanner":            policy.Options.CrawlerScanner,
			"crawler_script":             policy.Options.CrawlerScript,
			"crawler_other":              policy.Options.CrawlerOther,
			"web_shell":                  policy.Options.WebShell,
			"cc":                         policy.Options.Cc,
			"custom":                     policy.Options.Custom,
			"blacklist":                  policy.Options.WhiteblackIp,
			"geolocation_access_control": policy.Options.GeoIp,
			"ignore":                     policy.Options.Ignore,
			"privacy":                    policy.Options.Privacy,
			"anti_tamper":                policy.Options.AntiTamper,
			"anti_leakage":               policy.Options.AntiLeakage,
			"followed_action":            policy.Options.FollowedAction,
			"bot_enable":                 policy.Options.BotEnable,
			"precise":                    policy.Options.Precise,
		},
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", policy.Name),
		d.Set("level", policy.Level),
		d.Set("protection_mode", policy.Action.Category),
		d.Set("full_detection", policy.FullDetection),
		d.Set("options", options),
		d.Set("domains", policy.Hosts),
		d.Set("created_at", policy.CreatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return fmterr.Errorf("error setting opentelekomcloud WAF dedicated instance fields: %w", err)
	}
	return nil
}

func resourceWafDedicatedPolicyV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	if d.HasChanges("name", "protection_mode", "level", "options", "full_detection") {
		if err := updateWafPolicy(ctx, d, meta); err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud WAF dedicated policy: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedPolicyV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedPolicyV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	err = policies.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting opentelekomcloud WAF dedicated policy : %w", err)
	}

	d.SetId("")
	return nil
}

func updateWafPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	var updateOpts policies.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("level") {
		updateOpts.Level = d.Get("level").(int)
	}

	if d.HasChange("full_detection") {
		updateOpts.FullDetection = pointerto.Bool(d.Get("full_detection").(bool))
	}

	if d.HasChange("options") {
		updateOpts.Options = buildOptions(d)
	}

	if d.HasChange("protection_mode") {
		updateOpts.Action = &policies.PolicyAction{
			Category: d.Get("protection_mode").(string),
		}
	}

	_, err = policies.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF Policy: %s", err)
	}
	return nil
}

func buildOptions(d *schema.ResourceData) *policies.PolicyOption {
	optionsRaw := d.Get("options").([]interface{})
	rawMap := optionsRaw[0].(map[string]interface{})

	options := &policies.PolicyOption{
		WebAttack:      pointerto.Bool(rawMap["web_attack"].(bool)),
		Common:         pointerto.Bool(rawMap["common"].(bool)),
		Crawler:        pointerto.Bool(rawMap["crawler"].(bool)),
		AntiCrawler:    pointerto.Bool(rawMap["anti_crawler"].(bool)),
		CrawlerEngine:  pointerto.Bool(rawMap["crawler_engine"].(bool)),
		CrawlerScanner: pointerto.Bool(rawMap["crawler_scanner"].(bool)),
		CrawlerScript:  pointerto.Bool(rawMap["crawler_script"].(bool)),
		CrawlerOther:   pointerto.Bool(rawMap["crawler_other"].(bool)),
		WebShell:       pointerto.Bool(rawMap["web_shell"].(bool)),
		Cc:             pointerto.Bool(rawMap["cc"].(bool)),
		Custom:         pointerto.Bool(rawMap["custom"].(bool)),
		WhiteblackIp:   pointerto.Bool(rawMap["blacklist"].(bool)),
		GeoIp:          pointerto.Bool(rawMap["geolocation_access_control"].(bool)),
		Ignore:         pointerto.Bool(rawMap["ignore"].(bool)),
		Privacy:        pointerto.Bool(rawMap["privacy"].(bool)),
		AntiTamper:     pointerto.Bool(rawMap["anti_tamper"].(bool)),
		AntiLeakage:    pointerto.Bool(rawMap["anti_leakage"].(bool)),
		FollowedAction: pointerto.Bool(rawMap["followed_action"].(bool)),
		// not works in swiss
		// BotEnable:      pointerto.Bool(rawMap["bot_enable"].(bool)),
		// lead to 500 in eu-de and swiss
		// Precise:        pointerto.Bool(rawMap["precise"].(bool)),
	}

	log.Printf("[DEBUG] Options for OpenTelekomCloud WAF Policy: %#v", options)
	return options
}
