package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/datamasking_rules"
)

func resourceWafDataMaskingRuleV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafDataMaskingRuleV1Create,
		Read:   resourceWafDataMaskingRuleV1Read,
		Update: resourceWafDataMaskingRuleV1Update,
		Delete: resourceWafDataMaskingRuleV1Delete,
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
				ForceNew: false,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"index": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceWafDataMaskingRuleV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := datamasking_rules.CreateOpts{
		Url:      d.Get("url").(string),
		Category: d.Get("category").(string),
		Index:    d.Get("index").(string),
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := datamasking_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF DataMasking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf datamasking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafDataMaskingRuleV1Read(d, meta)
}

func resourceWafDataMaskingRuleV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	n, err := datamasking_rules.Get(wafClient, policy_id, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf DataMasking Rule: %s", err)
	}

	d.SetId(n.Id)
	d.Set("url", n.Url)
	d.Set("category", n.Category)
	d.Set("index", n.Index)
	d.Set("policy_id", n.PolicyID)

	return nil
}

func resourceWafDataMaskingRuleV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts datamasking_rules.UpdateOpts

	if d.HasChange("url") || d.HasChange("category") || d.HasChange("index") {
		updateOpts.Url = d.Get("url").(string)
		updateOpts.Category = d.Get("category").(string)
		updateOpts.Index = d.Get("index").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	if updateOpts != (datamasking_rules.UpdateOpts{}) {
		policy_id := d.Get("policy_id").(string)
		_, err = datamasking_rules.Update(wafClient, policy_id, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF DataMasking Rule: %s", err)
		}
	}

	return resourceWafDataMaskingRuleV1Read(d, meta)
}

func resourceWafDataMaskingRuleV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = datamasking_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF DataMasking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
