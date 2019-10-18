package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/webtamperprotection_rules"
)

func resourceWafWebTamperProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafWebTamperProtectionRuleV1Create,
		Read:   resourceWafWebTamperProtectionRuleV1Read,
		Delete: resourceWafWebTamperProtectionRuleV1Delete,
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
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWafWebTamperProtectionRuleV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := webtamperprotection_rules.CreateOpts{
		Hostname: d.Get("hostname").(string),
		Url:      d.Get("url").(string),
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := webtamperprotection_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Web Tamper Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf web tamper protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafWebTamperProtectionRuleV1Read(d, meta)
}

func resourceWafWebTamperProtectionRuleV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	n, err := webtamperprotection_rules.Get(wafClient, policy_id, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf Web Tamper Protection Rule: %s", err)
	}

	d.SetId(n.Id)
	d.Set("hostname", n.Hostname)
	d.Set("url", n.Url)
	d.Set("policy_id", n.PolicyID)

	return nil
}

func resourceWafWebTamperProtectionRuleV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = webtamperprotection_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF Web Tamper Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
