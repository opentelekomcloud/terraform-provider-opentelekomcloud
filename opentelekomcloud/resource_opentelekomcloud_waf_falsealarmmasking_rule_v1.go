package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/falsealarmmasking_rules"
)

func resourceWafFalseAlarmMaskingRuleV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafFalseAlarmMaskingRuleV1Create,
		Read:   resourceWafFalseAlarmMaskingRuleV1Read,
		Delete: resourceWafFalseAlarmMaskingRuleV1Delete,
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
			"rule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWafFalseAlarmMaskingRuleV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := falsealarmmasking_rules.CreateOpts{
		Url:  d.Get("url").(string),
		Rule: d.Get("rule").(string),
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := falsealarmmasking_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF False Alarm Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf falsealarmmasking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafFalseAlarmMaskingRuleV1Read(d, meta)
}

func resourceWafFalseAlarmMaskingRuleV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	rules, err := falsealarmmasking_rules.List(wafClient, policy_id).Extract()

	if err != nil {
		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf False Alarm Masking Rule: %s", err)
	}
	for _, r := range rules {
		if r.Id == d.Id() {
			d.SetId(r.Id)
			d.Set("url", r.Url)
			d.Set("rule", r.Rule)
			d.Set("policy_id", r.PolicyID)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceWafFalseAlarmMaskingRuleV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = falsealarmmasking_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF False Alarm Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
