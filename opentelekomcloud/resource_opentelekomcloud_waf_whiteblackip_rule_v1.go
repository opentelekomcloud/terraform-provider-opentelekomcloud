package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/whiteblackip_rules"
)

func resourceWafWhiteBlackIpRuleV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafWhiteBlackIpRuleV1Create,
		Read:   resourceWafWhiteBlackIpRuleV1Read,
		Update: resourceWafWhiteBlackIpRuleV1Update,
		Delete: resourceWafWhiteBlackIpRuleV1Delete,
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
			"addr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"white": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  0,
			},
		},
	}
}

func resourceWafWhiteBlackIpRuleV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wafClient, err := config.wafV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := whiteblackip_rules.CreateOpts{
		Addr:  d.Get("addr").(string),
		White: d.Get("white").(int),
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := whiteblackip_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomcomCloud WAF WhiteBlackIP Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf whiteblackip rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafWhiteBlackIpRuleV1Read(d, meta)
}

func resourceWafWhiteBlackIpRuleV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	n, err := whiteblackip_rules.Get(wafClient, policy_id, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Waf WhiteBlackIP Rule: %s", err)
	}

	d.SetId(n.Id)
	d.Set("addr", n.Addr)
	d.Set("white", n.White)
	d.Set("policy_id", n.PolicyID)

	return nil
}

func resourceWafWhiteBlackIpRuleV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts whiteblackip_rules.UpdateOpts

	if d.HasChange("addr") || d.HasChange("white") {
		updateOpts.Addr = d.Get("addr").(string)
		white := d.Get("white").(int)
		updateOpts.White = &white
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	if updateOpts != (whiteblackip_rules.UpdateOpts{}) {
		policy_id := d.Get("policy_id").(string)
		_, err = whiteblackip_rules.Update(wafClient, policy_id, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud WAF WhiteBlackIP Rule: %s", err)
		}
	}

	return resourceWafWhiteBlackIpRuleV1Read(d, meta)
}

func resourceWafWhiteBlackIpRuleV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.wafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = whiteblackip_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud WAF WhiteBlackIP Rule: %s", err)
	}

	d.SetId("")
	return nil
}
