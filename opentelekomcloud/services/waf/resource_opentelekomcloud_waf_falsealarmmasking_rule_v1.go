package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/falsealarmmasking_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafFalseAlarmMaskingRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafFalseAlarmMaskingRuleV1Create,
		ReadContext:   resourceWafFalseAlarmMaskingRuleV1Read,
		DeleteContext: resourceWafFalseAlarmMaskingRuleV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceWafFalseAlarmMaskingRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	wafClient, err := config.WafV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := falsealarmmasking_rules.CreateOpts{
		Url:  d.Get("url").(string),
		Rule: d.Get("rule").(string),
	}

	policy_id := d.Get("policy_id").(string)
	rule, err := falsealarmmasking_rules.Create(wafClient, policy_id, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF False Alarm Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf falsealarmmasking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafFalseAlarmMaskingRuleV1Read(ctx, d, meta)
}

func resourceWafFalseAlarmMaskingRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	policy_id := d.Get("policy_id").(string)
	rules, err := falsealarmmasking_rules.List(wafClient, policy_id).Extract()

	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf False Alarm Masking Rule: %s", err)
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

func resourceWafFalseAlarmMaskingRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policy_id := d.Get("policy_id").(string)
	err = falsealarmmasking_rules.Delete(wafClient, policy_id, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF False Alarm Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
