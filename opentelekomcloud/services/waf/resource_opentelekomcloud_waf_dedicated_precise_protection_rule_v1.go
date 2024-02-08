package waf

import (
	"context"
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

func ResourceWafDedicatedPreciseProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedPreciseProtectionRuleV1Create,
		ReadContext:   resourceWafDedicatedPreciseProtectionRuleV1Read,
		DeleteContext: resourceWafDedicatedPreciseProtectionRuleV1Delete,
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
			"time": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"start": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"terminal": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"description": {
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
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"url", "user-agent", "referer", "ip",
									"method", "request_line", "request", "params",
									"cookie", "header",
								},
								false),
						},
						"logic_operation": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
								[]string{"block", "pass", "log"},
								false),
						},
						"followed_action_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
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

func resourceWafDedicatedPreciseProtectionRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	rawActionList := d.Get("action").(*schema.Set).List()
	var action rules.CustomActionObject
	if len(rawActionList) > 0 {
		rawAction := rawActionList[0].(map[string]interface{})
		action = rules.CustomActionObject{
			Category:         rawAction["category"].(string),
			FollowedActionId: rawAction["followed_action_id"].(string),
		}
	}

	var conditionList []rules.CustomConditionsObject
	conditions := d.Get("conditions").([]interface{})
	for _, c := range conditions {
		cond := c.(map[string]interface{})
		contentsRaw := cond["contents"].([]interface{})
		contents := make([]string, len(contentsRaw))

		for i, content := range contentsRaw {
			contents[i] = content.(string)
		}

		condition := rules.CustomConditionsObject{
			Category:       cond["category"].(string),
			Index:          cond["index"].(string),
			LogicOperation: cond["logic_operation"].(string),
			ValueListId:    cond["value_list_id"].(string),
			Contents:       contents,
		}
		conditionList = append(conditionList, condition)
	}

	createOpts := rules.CreateCustomOpts{
		Time:        pointerto.Bool(d.Get("time").(bool)),
		Start:       int64(d.Get("start").(int)),
		Terminal:    int64(d.Get("terminal").(int)),
		Description: d.Get("description").(string),
		Conditions:  conditionList,
		Action:      &action,
		Priority:    pointerto.Int(d.Get("priority").(int)),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.CreateCustom(client, policyID, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated Precise Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf dedicated precise protection rule created: %#v", rule)
	d.SetId(rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedPreciseProtectionRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedPreciseProtectionRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.GetCustom(client, policyID, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated Precise Protection Rule: %s", err)
	}

	mErr := multierror.Append(
		d.Set("policy_id", rule.PolicyId),
		d.Set("description", rule.Description),
		d.Set("priority", rule.Priority),
		d.Set("start", rule.Start),
		d.Set("terminal", rule.Terminal),
		d.Set("time", d.Get("time").(bool)),
		d.Set("status", rule.Status),
		d.Set("created_at", rule.CreatedAt),
	)

	var conditions []map[string]interface{}
	for _, conditionObj := range rule.Conditions {
		condition := map[string]interface{}{
			"category":        conditionObj.Category,
			"index":           conditionObj.Index,
			"contents":        conditionObj.Contents,
			"logic_operation": conditionObj.LogicOperation,
			"value_list_id":   conditionObj.ValueListId,
		}
		conditions = append(conditions, condition)
	}

	action := []map[string]interface{}{
		{
			"category":           rule.Action.Category,
			"followed_action_id": rule.Action.FollowedActionId,
		},
	}
	mErr = multierror.Append(mErr,
		d.Set("conditions", conditions),
		d.Set("action", action),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedPreciseProtectionRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteCustomRule(client, policyID, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated Precise Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
