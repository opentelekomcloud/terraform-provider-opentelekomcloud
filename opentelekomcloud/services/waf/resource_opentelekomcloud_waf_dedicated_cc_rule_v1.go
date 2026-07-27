package waf

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDedicatedCcRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedCcRuleV1Create,
		ReadContext:   resourceWafDedicatedCcRuleV1Read,
		DeleteContext: resourceWafDedicatedCcRuleV1Delete,
		CustomizeDiff: validateWafDedicatedCcRuleV1Mode,
		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("policy_id", "id"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"url", "ip", "params", "cookie", "header"},
								false),
						},
						"logic_operation": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"contain", "not_contain", "equal",
								"not_equal", "prefix", "not_prefix",
								"suffix", "not_suffix", "contain_any",
								"not_contain_all", "equal_any", "not_equal_all",
								"prefix_any", "not_prefix_all", "suffix_any",
								"not_suffix_all", "num_greater", "num_less",
								"num_equal", "num_not_equal", "exist",
								"not_exist",
							}, false),
						},
						"contents": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"value_list_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"index": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"action": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"captcha", "block", "log", "dynamic_block"},
								false),
						},
						"content_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"application/json", "text/html", "text/xml"},
								false),
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"tag_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"ip", "cookie", "header", "other"},
					false),
			},
			"tag_index": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag_contents": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"limit_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},
			"limit_period": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
			},
			"unlock_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 2147483647),
			},
			"lock_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func validateWafDedicatedCcRuleV1Mode(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	switch d.Get("mode").(int) {
	case 0:
		if d.Get("url").(string) == "" {
			return fmt.Errorf("url must be configured when mode is 0")
		}
	case 1:
		if len(d.Get("conditions").([]interface{})) == 0 {
			return fmt.Errorf("at least one conditions block must be configured when mode is 1")
		}
	}
	return nil
}

func resourceWafDedicatedCcRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	tagC := d.Get("tag_contents").([]interface{})
	tagContents := make([]string, len(tagC))
	for i, v := range tagC {
		tagContents[i] = v.(string)
	}

	tagCondition := rules.CcTagConditionObject{
		Category: d.Get("tag_category").(string),
		Contents: tagContents,
	}

	rawActionList := d.Get("action").(*schema.Set).List()
	var action rules.CcActionObject
	if len(rawActionList) > 0 {
		rawAction := rawActionList[0].(map[string]interface{})
		action = rules.CcActionObject{
			Category: rawAction["category"].(string),
		}
		if rawAction["content_type"].(string) != "" || rawAction["content"].(string) != "" {
			action.Detail = &rules.CcDetailObject{
				Response: &rules.CcResponseObject{
					ContentType: rawAction["content_type"].(string),
					Content:     rawAction["content"].(string),
				},
			}
		}
	}

	var conditionList []rules.CcConditionsObject
	conditions := d.Get("conditions").([]interface{})
	for _, c := range conditions {
		cond := c.(map[string]interface{})
		contentsRaw := cond["contents"].([]interface{})
		contents := make([]string, len(contentsRaw))

		for i, content := range contentsRaw {
			contents[i] = content.(string)
		}

		condition := rules.CcConditionsObject{
			Category:       cond["category"].(string),
			LogicOperation: cond["logic_operation"].(string),
			ValueListId:    cond["value_list_id"].(string),
			Contents:       contents,
		}
		if index := cond["index"].(string); index != "" {
			condition.Index = pointerto.String(index)
		}
		conditionList = append(conditionList, condition)
	}

	mode := d.Get("mode").(int)
	createOpts := rules.CreateCcOpts{
		Mode:         &mode,
		Url:          d.Get("url").(string),
		Conditions:   conditionList,
		Action:       &action,
		TagType:      d.Get("tag_type").(string),
		TagIndex:     d.Get("tag_index").(string),
		TagCondition: &tagCondition,
		LimitNum:     int64(d.Get("limit_num").(int)),
		LimitPeriod:  int64(d.Get("limit_period").(int)),
		UnlockNum:    int64(d.Get("unlock_num").(int)),
		LockTime:     pointerto.Int(d.Get("lock_time").(int)),
		Description:  d.Get("description").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.CreateCc(client, policyID, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated CC Attack Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf dedicated cc attack protection rule created: %#v", rule)
	d.SetId(rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedCcRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedCcRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.GetCc(client, policyID, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated CC Attack Protection Rule: %s", err)
	}

	mErr := multierror.Append(
		d.Set("policy_id", rule.PolicyId),
		d.Set("mode", rule.Mode),
		d.Set("url", rule.Url),
		d.Set("limit_num", rule.LimitNum),
		d.Set("limit_period", rule.LimitPeriod),
		d.Set("lock_time", rule.LockTime),
		d.Set("tag_type", rule.TagType),
		d.Set("tag_index", rule.TagIndex),
		d.Set("tag_category", rule.TagCondition.Category),
		d.Set("tag_contents", rule.TagCondition.Contents),
		d.Set("description", rule.Description),
		d.Set("unlock_num", rule.UnlockNum),
		d.Set("status", rule.Status),
		d.Set("created_at", rule.CreatedAt),
	)

	var conditions []map[string]interface{}
	for _, conditionObj := range rule.Conditions {
		condition := map[string]interface{}{
			"category":        conditionObj.Category,
			"contents":        conditionObj.Contents,
			"logic_operation": conditionObj.LogicOperation,
			"value_list_id":   conditionObj.ValueListId,
		}
		if conditionObj.Index != nil {
			condition["index"] = *conditionObj.Index
		}
		conditions = append(conditions, condition)
	}

	actionValue := map[string]interface{}{
		"category": rule.Action.Category,
	}
	if rule.Action.Detail != nil && rule.Action.Detail.Response != nil {
		actionValue["content_type"] = rule.Action.Detail.Response.ContentType
		actionValue["content"] = rule.Action.Detail.Response.Content
	}
	action := []map[string]interface{}{actionValue}
	mErr = multierror.Append(mErr,
		d.Set("conditions", conditions),
		d.Set("action", action),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedCcRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteCcRule(client, policyID, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated CC Attack Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
