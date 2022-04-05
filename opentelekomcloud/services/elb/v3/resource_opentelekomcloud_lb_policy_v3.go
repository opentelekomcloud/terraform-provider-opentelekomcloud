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
					"REDIRECT_TO_POOL", "REDIRECT_TO_LISTENER",
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
								"HOST_NAME", "PATH",
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

func resourceLBPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := policies.CreateOpts{
		Action:             policies.Action(d.Get("action").(string)),
		Description:        d.Get("description").(string),
		ListenerID:         d.Get("listener_id").(string),
		Name:               d.Get("name").(string),
		Position:           d.Get("position").(int),
		ProjectID:          d.Get("project_id").(string),
		RedirectListenerID: d.Get("redirect_listener_id").(string),
		RedirectPoolID:     d.Get("redirect_pool_id").(string),
		Rules:              getRules(d),
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
	if d.HasChange("rules") {
		opts.Rules = getRules(d)
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
