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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDedicatedAlarmMaskingRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedAlarmMaskingRuleV1Create,
		ReadContext:   resourceWafDedicatedAlarmMaskingRuleV1Read,
		DeleteContext: resourceWafDedicatedAlarmMaskingRuleV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("policy_id", "id"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
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
								"suffix", "not_suffix",
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
						"index": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"advanced_settings": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contents": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"index": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"rule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func getDomains(d *schema.ResourceData) []string {
	domains := d.Get("domains").([]interface{})
	dom := make([]string, len(domains))
	for i, domain := range domains {
		dom[i] = domain.(string)
	}
	return dom
}

func getAmConditions(d *schema.ResourceData) []rules.IgnoreCondition {
	var conditionList []rules.IgnoreCondition
	conditions := d.Get("conditions").([]interface{})
	for _, c := range conditions {
		cond := c.(map[string]interface{})
		contentsRaw := cond["contents"].([]interface{})
		contents := make([]string, len(contentsRaw))

		for i, content := range contentsRaw {
			contents[i] = content.(string)
		}

		condition := rules.IgnoreCondition{
			Category:       cond["category"].(string),
			Index:          cond["index"].(string),
			LogicOperation: cond["logic_operation"].(string),
			Contents:       contents,
		}
		conditionList = append(conditionList, condition)
	}
	return conditionList
}

func getAmAdvancedSettings(d *schema.ResourceData) []rules.AdvancedIgnoreObject {
	var advancedList []rules.AdvancedIgnoreObject
	advanced := d.Get("advanced_settings").([]interface{})
	for _, a := range advanced {
		adv := a.(map[string]interface{})
		contentsRaw := adv["contents"].([]interface{})
		contents := make([]string, len(contentsRaw))

		for i, content := range contentsRaw {
			contents[i] = content.(string)
		}

		advancedObj := rules.AdvancedIgnoreObject{
			Index:    adv["index"].(string),
			Contents: contents,
		}
		advancedList = append(advancedList, advancedObj)
	}
	return advancedList
}

func resourceWafDedicatedAlarmMaskingRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	createOpts := rules.CreateIgnoreOpts{
		Domains:     getDomains(d),
		Conditions:  getAmConditions(d),
		Mode:        1,
		Rule:        d.Get("rule").(string),
		Advanced:    getAmAdvancedSettings(d),
		Description: d.Get("description").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.CreateIgnore(client, policyID, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated Alarms Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf dedicated alarm masking rule created: %#v", rule)
	d.SetId(rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedAlarmMaskingRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedAlarmMaskingRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.GetIgnore(client, policyID, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated Alarm Masking Rule: %s", err)
	}

	mErr := multierror.Append(
		d.Set("policy_id", rule.PolicyId),
		d.Set("rule", rule.Rule),
		d.Set("advanced_settings", rule.Advanced),
		d.Set("domains", rule.Domains),
		d.Set("description", rule.Description),
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
		}
		conditions = append(conditions, condition)
	}

	var advanced []map[string]interface{}
	for _, advancedObj := range rule.Advanced {
		adv := map[string]interface{}{
			"index":    advancedObj.Index,
			"contents": advancedObj.Contents,
		}
		advanced = append(advanced, adv)
	}
	mErr = multierror.Append(mErr,
		d.Set("conditions", conditions),
		d.Set("advanced_settings", advanced),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedAlarmMaskingRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteIgnoreRule(client, policyID, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated Alarm Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
