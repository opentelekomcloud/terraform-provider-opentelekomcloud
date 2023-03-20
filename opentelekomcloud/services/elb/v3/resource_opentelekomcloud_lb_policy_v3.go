package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/policies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBPolicyV3Create,
		ReadContext:   resourceLBPolicyV3Read,
		UpdateContext: resourceLBPolicyV3Update,
		DeleteContext: resourceLBPolicyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 0xff),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 0xff),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"REDIRECT_TO_POOL", "REDIRECT_TO_LISTENER", "REDIRECT_TO_URL", "FIXED_RESPONSE",
				}, false),
			},
			"listener_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"redirect_listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"redirect_pool_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"position": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 10000),
			},
			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 10,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HOST_NAME", "PATH", "METHOD",
								"HEADER", "QUERY_STRING", "SOURCE_IP",
							}, false),
						},
						"compare_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"EQUAL_TO", "REGEX", "STARTS_WITH",
							}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_response_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"message_body": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"redirect_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"redirect_url_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "${protocol}",
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "${host}",
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "${port}",
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "${path}",
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "${query}",
						},
					},
				},
			},
			"redirect_pools_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},
		},
	}
}

func getRules(d *schema.ResourceData) []policies.Rule {
	ruleListRaw := d.Get("rules").(*schema.Set).List()
	var ruleList []policies.Rule

	for _, rule := range ruleListRaw {
		ruleRaw := rule.(map[string]interface{})

		ruleList = append(ruleList, policies.Rule{
			Type:        rules.RuleType(ruleRaw["type"].(string)),
			CompareType: rules.CompareType(ruleRaw["compare_type"].(string)),
			Value:       ruleRaw["value"].(string),
		})
	}

	return ruleList
}

func getFixedResponseConfig(d *schema.ResourceData) *policies.FixedResponseOptions {
	responseListRaw := d.Get("fixed_response_config").(*schema.Set).List()
	var fixedResponse *policies.FixedResponseOptions
	if len(responseListRaw) == 1 {
		for _, rule := range responseListRaw {
			fixedResponseRaw := rule.(map[string]interface{})

			fixedResponse = &policies.FixedResponseOptions{
				StatusCode:  fixedResponseRaw["status_code"].(string),
				ContentType: fixedResponseRaw["content_type"].(string),
				MessageBody: fixedResponseRaw["message_body"].(string),
			}
		}
	}

	return fixedResponse
}

func getRedirectUrlConfig(d *schema.ResourceData) *policies.RedirectUrlOptions {
	ruleListRaw := d.Get("redirect_url_config").(*schema.Set).List()
	var redirectUrlConfig *policies.RedirectUrlOptions
	if len(ruleListRaw) == 1 {
		for _, rule := range ruleListRaw {
			redirectUrlConfigRaw := rule.(map[string]interface{})

			redirectUrlConfig = &policies.RedirectUrlOptions{
				StatusCode: redirectUrlConfigRaw["status_code"].(string),
				Protocol:   redirectUrlConfigRaw["protocol"].(string),
				Host:       redirectUrlConfigRaw["host"].(string),
				Path:       redirectUrlConfigRaw["path"].(string),
				Query:      redirectUrlConfigRaw["query"].(string),
				Port:       redirectUrlConfigRaw["port"].(string),
			}
		}
	}

	return redirectUrlConfig
}

func getRedirectPoolsConfig(d *schema.ResourceData) []policies.RedirectPoolOptions {
	ruleListRaw := d.Get("redirect_pools_config").(*schema.Set).List()
	var redirectPoolsList []policies.RedirectPoolOptions

	for _, rule := range ruleListRaw {
		redirectPoolsRaw := rule.(map[string]interface{})

		redirectPoolsList = append(redirectPoolsList, policies.RedirectPoolOptions{
			PoolId: redirectPoolsRaw["pool_id"].(string),
			Weight: redirectPoolsRaw["weight"].(string),
		})
	}

	return redirectPoolsList
}
func resourceLBPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := policies.CreateOpts{
		Action:              policies.Action(d.Get("action").(string)),
		Description:         d.Get("description").(string),
		ListenerID:          d.Get("listener_id").(string),
		Name:                d.Get("name").(string),
		Position:            d.Get("position").(int),
		Priority:            d.Get("priority").(int),
		ProjectID:           d.Get("project_id").(string),
		RedirectListenerID:  d.Get("redirect_listener_id").(string),
		RedirectPoolID:      d.Get("redirect_pool_id").(string),
		RedirectUrl:         d.Get("redirect_url").(string),
		Rules:               getRules(d),
		FixedResponseConfig: getFixedResponseConfig(d),
		RedirectUrlConfig:   getRedirectUrlConfig(d),
		RedirectPoolsConfig: getRedirectPoolsConfig(d),
	}

	policy, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud policy: %w", err)
	}

	d.SetId(policy.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPolicyV3Read(clientCtx, d, meta)
}

func resourceLBPolicyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	policy, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error viewing details of LB Policy v3")
	}

	var ruleList []interface{}

	for _, v := range policy.Rules {
		rule, err := rules.Get(client, d.Id(), v.ID).Extract()
		if err != nil {
			return fmterr.Errorf("error receiving policy rule: %w", err)
		}
		ruleList = append(ruleList, map[string]interface{}{
			"type":         rule.Type,
			"compare_type": rule.CompareType,
			"value":        rule.Value,
		})
	}
	fixedResponseConfig := make([]map[string]interface{}, 0)
	if policy.FixedResponseConfig.StatusCode != "" {
		fixedResponse := map[string]interface{}{
			"status_code":  policy.FixedResponseConfig.StatusCode,
			"content_type": policy.FixedResponseConfig.ContentType,
			"message_body": policy.FixedResponseConfig.MessageBody,
		}
		fixedResponseConfig = append(fixedResponseConfig, fixedResponse)
	}
	redirectUrlConfig := make([]map[string]interface{}, 0)
	if policy.RedirectUrlConfig.StatusCode != "" {
		redirectUrl := map[string]interface{}{
			"status_code": policy.RedirectUrlConfig.StatusCode,
			"path":        policy.RedirectUrlConfig.Path,
			"port":        policy.RedirectUrlConfig.Port,
			"query":       policy.RedirectUrlConfig.Query,
			"host":        policy.RedirectUrlConfig.Host,
			"protocol":    policy.RedirectUrlConfig.Protocol,
		}
		redirectUrlConfig = append(redirectUrlConfig, redirectUrl)
	}

	var redirectPoolsList []interface{}
	for _, v := range policy.RedirectPoolsConfig {
		redirectPoolsList = append(redirectPoolsList, map[string]interface{}{
			"pool_id": v.PoolId,
			"weight":  v.Weight,
		})
	}

	mErr := multierror.Append(
		d.Set("name", policy.Name),
		d.Set("description", policy.Description),
		d.Set("project_id", policy.ProjectID),
		d.Set("action", policy.Action),
		d.Set("listener_id", policy.ListenerID),
		d.Set("redirect_listener_id", policy.RedirectListenerID),
		d.Set("redirect_pool_id", policy.RedirectPoolID),
		d.Set("position", policy.Position),
		d.Set("rules", ruleList),
		d.Set("status", policy.Status),
		d.Set("priority", policy.Priority),
		d.Set("redirect_url", policy.RedirectUrl),
		d.Set("fixed_response_config", fixedResponseConfig),
		d.Set("redirect_url_config", redirectUrlConfig),
		d.Set("redirect_pools_config", redirectPoolsList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB Policy v3 fields: %w", err)
	}

	return nil
}

func resourceLBPolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := policies.UpdateOpts{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
	}
	if d.HasChange("redirect_listener_id") {
		opts.RedirectListenerID = d.Get("redirect_listener_id").(string)
	}
	if d.HasChange("redirect_pool_id") {
		opts.RedirectPoolID = d.Get("redirect_pool_id").(string)
	}
	if d.HasChange("priority") {
		opts.Priority = d.Get("priority").(int)
	}
	if d.HasChange("rules") {
		opts.Rules = getRules(d)
	}
	if d.HasChange("fixed_response_config") {
		opts.FixedResponseConfig = getFixedResponseConfig(d)
	}
	if d.HasChange("redirect_url_config") {
		opts.RedirectUrlConfig = getRedirectUrlConfig(d)
	}
	if d.HasChange("redirect_pools_config") {
		opts.RedirectPoolsConfig = getRedirectPoolsConfig(d)
	}

	_, err = policies.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating LB Policy v3: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPolicyV3Read(clientCtx, d, meta)
}

func resourceLBPolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if err := policies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting LB Policy v3: %w", err)
	}

	return nil
}
