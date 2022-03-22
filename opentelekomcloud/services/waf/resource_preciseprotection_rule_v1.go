package waf

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/preciseprotection_rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafPreciseProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafPreciseProtectionRuleV1Create,
		ReadContext:   resourceWafPreciseProtectionRuleV1Read,
		DeleteContext: resourceWafPreciseProtectionRuleV1Delete,
		Importer:      wafRuleImporter(),

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"time": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"start": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"path", "user-agent", "ip", "params", "cookie", "referer", "header"}, false),
						},
						"index": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"logic": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"contain", "not_contain", "equal", "not_equal", "prefix", "not_prefix", "suffix", "not_suffix",
							}, false),
						},
						"contents": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"action_category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func getConditions(d *schema.ResourceData) []preciseprotection_rules.Condition {
	var conditionOpts []preciseprotection_rules.Condition

	conditions := d.Get("conditions").([]interface{})
	for _, v := range conditions {
		cond := v.(map[string]interface{})
		contentsRaw := cond["contents"].([]interface{})
		contents := make([]string, len(contentsRaw))

		for i, v := range contentsRaw {
			contents[i] = v.(string)
		}

		condition := preciseprotection_rules.Condition{
			Category: cond["category"].(string),
			Index:    cond["index"].(string),
			Logic:    cond["logic"].(string),
			Contents: contents,
		}
		conditionOpts = append(conditionOpts, condition)
	}

	log.Printf("[DEBUG] getConditions: %#v", conditionOpts)
	return conditionOpts
}

func getPreciseAction(d *schema.ResourceData) preciseprotection_rules.Action {
	action := preciseprotection_rules.Action{
		Category: d.Get("action_category").(string),
	}

	log.Printf("[DEBUG] getPreciseAction: %#v", action)
	return action
}

func resourceWafPreciseProtectionRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	wafClient, err := config.WafV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}
	priority := d.Get("priority").(int)
	createOpts := preciseprotection_rules.CreateOpts{
		Name:       d.Get("name").(string),
		Time:       d.Get("time").(bool),
		Conditions: getConditions(d),
		Action:     getPreciseAction(d),
		Priority:   &priority,
	}

	if _, ok := d.GetOk("start"); ok {
		start, err := strconv.ParseInt(d.Get("start").(string), 10, 64)
		if err != nil {
			return fmterr.Errorf("error converting start: %s", err)
		}
		createOpts.Start = start
	}
	if _, ok := d.GetOk("cache_control"); ok {
		end, err := strconv.ParseInt(d.Get("end").(string), 10, 64)
		if err != nil {
			return fmterr.Errorf("error converting end: %s", err)
		}
		createOpts.End = end
	}

	policyID := d.Get("policy_id").(string)
	rule, err := preciseprotection_rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Precise Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf precise protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafPreciseProtectionRuleV1Read(ctx, d, meta)
}

func resourceWafPreciseProtectionRuleV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	policyID := d.Get("policy_id").(string)
	n, err := preciseprotection_rules.Get(wafClient, policyID, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Precise Protection Rule: %s", err)
	}

	d.SetId(n.Id)
	mErr := multierror.Append(
		d.Set("policy_id", n.PolicyID),
		d.Set("name", n.Name),
		d.Set("time", n.Time),
		d.Set("start", strconv.FormatInt(n.Start, 10)),
		d.Set("end", strconv.FormatInt(n.End, 10)),
	)
	conditions := make([]map[string]interface{}, len(n.Conditions))
	for i, condition := range n.Conditions {
		conditions[i] = make(map[string]interface{})
		conditions[i]["category"] = condition.Category
		conditions[i]["index"] = condition.Index
		conditions[i]["logic"] = condition.Logic
		conditions[i]["contents"] = condition.Contents
	}
	mErr = multierror.Append(mErr,
		d.Set("conditions", conditions),
		d.Set("action_category", n.Action.Category),
		d.Set("priority", n.Priority),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafPreciseProtectionRuleV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = preciseprotection_rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Precise Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
