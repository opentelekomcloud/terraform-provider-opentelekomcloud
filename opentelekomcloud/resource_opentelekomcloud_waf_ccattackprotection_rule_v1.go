package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/ccattackprotection_rules"
)

func resourceWafCcAttackProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafCcAttackProtectionRuleV1Create,
		Read:   resourceWafCcAttackProtectionRuleV1Read,
		Delete: resourceWafCcAttackProtectionRuleV1Delete,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"limit_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"limit_period": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"lock_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tag_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"action_category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"block_content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"block_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func getTagCondition(d *schema.ResourceData) ccattackprotection_rules.TagCondition {
	v := d.Get("tag_contents").([]interface{})
	contents := make([]string, len(v))
	for i, v := range v {
		contents[i] = v.(string)
	}

	condition := ccattackprotection_rules.TagCondition{
		Category: d.Get("tag_category").(string),
		Contents: contents,
	}

	log.Printf("[DEBUG] getTagCondition: %#v", condition)
	return condition
}

func getCcAction(d *schema.ResourceData) ccattackprotection_rules.Action {
	response := ccattackprotection_rules.Response{
		ContentType: d.Get("block_content_type").(string),
		Content:     d.Get("block_content").(string),
	}
	detail := ccattackprotection_rules.Detail{
		Response: response,
	}

	action := ccattackprotection_rules.Action{
		Category: d.Get("action_category").(string),
		Detail:   detail,
	}

	log.Printf("[DEBUG] getAction: %#v", action)
	return action
}

func resourceWafCcAttackProtectionRuleV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	limit_num := d.Get("limit_num").(int)
	limit_period := d.Get("limit_period").(int)
	lock_time := d.Get("lock_time").(int)
	createOpts := ccattackprotection_rules.CreateOpts{
		Url:         d.Get("url").(string),
		LimitNum:    &limit_num,
		LimitPeriod: &limit_period,
		LockTime:    &lock_time,
		TagType:     d.Get("tag_type").(string),
		TagIndex:    d.Get("tag_index").(string),
		Action:      getCcAction(d),
	}

	_, tag_category_ok := d.GetOk("tag_category")
	_, tag_contents_ok := d.GetOk("tag_contents")
	if tag_category_ok && tag_contents_ok {
		createOpts.TagCondition = getTagCondition(d)
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := ccattackprotection_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF CC Attack Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf cc attack protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafCcAttackProtectionRuleV1Read(d, meta)
}

func resourceWafCcAttackProtectionRuleV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	n, err := ccattackprotection_rules.Get(wafClient, policy_id, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf CC Attack Protection Rule: %s", err)
	}

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("url", n.Url)
	d.Set("limit_num", n.LimitNum)
	d.Set("limit_period", n.LimitPeriod)
	d.Set("lock_time", n.LockTime)
	d.Set("tag_type", n.TagType)
	d.Set("tag_index", n.TagIndex)
	d.Set("tag_category", n.TagCondition.Category)
	d.Set("tag_contents", n.TagCondition.Contents)
	d.Set("action_category", n.Action.Category)
	d.Set("block_content_type", n.Action.Detail.Response.ContentType)
	d.Set("block_content", n.Action.Detail.Response.Content)
	d.Set("default", n.Default)

	return nil
}

func resourceWafCcAttackProtectionRuleV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = ccattackprotection_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF CC Attack Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
