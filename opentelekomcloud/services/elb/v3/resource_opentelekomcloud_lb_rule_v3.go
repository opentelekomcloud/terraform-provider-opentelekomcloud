package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBRuleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBRuleV3Create,
		ReadContext:   resourceLBRuleV3Read,
		UpdateContext: resourceLBRuleV3Update,
		DeleteContext: resourceLBRuleV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("policy_id", "rule_id"),
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func getConditions(d *schema.ResourceData) []rules.Condition {
	conditionListRaw := d.Get("conditions").(*schema.Set).List()
	var conditionList []rules.Condition

	for _, rule := range conditionListRaw {
		ruleRaw := rule.(map[string]interface{})

		conditionList = append(conditionList, rules.Condition{
			Key:   ruleRaw["key"].(string),
			Value: ruleRaw["value"].(string),
		})
	}

	return conditionList
}

func resourceLBRuleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := rules.CreateOpts{
		Type:        rules.RuleType(d.Get("type").(string)),
		CompareType: rules.CompareType(d.Get("compare_type").(string)),
		Value:       d.Get("value").(string),
		ProjectID:   d.Get("project_id").(string),
		Conditions:  getConditions(d),
	}
	policyID := d.Get("policy_id").(string)

	rule, err := rules.Create(client, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Rule: %w", err)
	}
	_ = d.Set("rule_id", rule.ID) // this can't ever return an error

	if err := common.SetComplexID(d, "policy_id", "rule_id"); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBRuleV3Read(clientCtx, d, meta)
}

func resourceLBRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	rule, err := rules.Get(client, policyID(d), ruleID(d)).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error viewing details of LB Rule v3")
	}

	var conditionsList []interface{}
	for _, v := range rule.Conditions {
		conditionsList = append(conditionsList, map[string]interface{}{
			"key":   v.Key,
			"value": v.Value,
		})
	}
	mErr := multierror.Append(
		d.Set("project_id", rule.ProjectID),
		d.Set("type", rule.Type),
		d.Set("compare_type", rule.CompareType),
		d.Set("value", rule.Value),
		d.Set("conditions", conditionsList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB Rule v3 fields: %w", err)
	}

	return nil
}

func resourceLBRuleV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	updateOpts := rules.UpdateOpts{}
	if d.HasChange("compare_type") {
		compareType := rules.CompareType(d.Get("compare_type").(string))
		updateOpts.CompareType = compareType
	}
	if d.HasChange("value") {
		updateOpts.Value = d.Get("value").(string)
	}
	if d.HasChange("conditions") {
		updateOpts.Conditions = getConditions(d)
	}

	_, err = rules.Update(client, policyID(d), ruleID(d), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating LB Rule v3: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBRuleV3Read(clientCtx, d, meta)
}

func resourceLBRuleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if err := rules.Delete(client, policyID(d), ruleID(d)).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting LB Rule v3: %w", err)
	}

	return nil
}

func policyID(d *schema.ResourceData) string {
	return d.Get("policy_id").(string)
}

func ruleID(d *schema.ResourceData) string {
	return d.Get("rule_id").(string)
}
