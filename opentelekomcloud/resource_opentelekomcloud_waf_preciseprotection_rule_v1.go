package opentelekomcloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/preciseprotection_rules"
)

func resourceWafPreciseProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafPreciseProtectionRuleV1Create,
		Read:   resourceWafPreciseProtectionRuleV1Read,
		Delete: resourceWafPreciseProtectionRuleV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"index": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"logic": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
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
		contents_raw := cond["contents"].([]interface{})
		contents := make([]string, len(contents_raw))

		for i, v := range contents_raw {
			contents[i] = v.(string)
		}

		condition := preciseprotection_rules.Condition{
			Category: cond["category"].(string),
			Index:    cond["index"].(string),
			Logic:    cond["logic"].(int),
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

func resourceWafPreciseProtectionRuleV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
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
			return fmt.Errorf("Error converting start: %s", err)
		}
		createOpts.Start = start
	}
	if _, ok := d.GetOk("cache_control"); ok {
		end, err := strconv.ParseInt(d.Get("end").(string), 10, 64)
		if err != nil {
			return fmt.Errorf("Error converting end: %s", err)
		}
		createOpts.End = end
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := preciseprotection_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Precise Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf precise protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafPreciseProtectionRuleV1Read(d, meta)
}

func resourceWafPreciseProtectionRuleV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	n, err := preciseprotection_rules.Get(wafClient, policy_id, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf Precise Protection Rule: %s", err)
	}

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("name", n.Name)
	d.Set("time", n.Time)
	d.Set("start", strconv.FormatInt(n.Start, 10))
	d.Set("end", strconv.FormatInt(n.End, 10))

	conditions := make([]map[string]interface{}, len(n.Conditions))
	for i, condition := range n.Conditions {
		conditions[i] = make(map[string]interface{})
		conditions[i]["category"] = condition.Category
		conditions[i]["index"] = condition.Index
		conditions[i]["logic"] = condition.Logic
		conditions[i]["contents"] = condition.Contents
	}
	d.Set("conditions", conditions)
	d.Set("action_category", n.Action.Category)
	d.Set("priority", n.Priority)

	return nil
}

func resourceWafPreciseProtectionRuleV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = preciseprotection_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF Precise Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
