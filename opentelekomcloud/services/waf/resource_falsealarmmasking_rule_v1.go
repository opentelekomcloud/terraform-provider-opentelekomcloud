package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
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
		Importer:      wafRuleImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		DeprecationMessage: "This resource is known to be broken due to the API changes and will be fixed in the upcoming releases",

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
		Path: d.Get("url").(string),
		// Rule: d.Get("rule").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := falsealarmmasking_rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF False Alarm Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf falsealarmmasking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafFalseAlarmMaskingRuleV1Read(ctx, d, meta)
}

func resourceWafFalseAlarmMaskingRuleV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	policyID := d.Get("policy_id").(string)
	rules, err := falsealarmmasking_rules.List(wafClient, policyID).Extract()

	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf False Alarm Masking Rule: %s", err)
	}
	for _, r := range rules {
		if r.Id == d.Id() {
			d.SetId(r.Id)
			mErr := multierror.Append(
				d.Set("url", r.Path),
				d.Set("rule", r.Rule),
				d.Set("policy_id", r.PolicyID),
			)
			if err := mErr.ErrorOrNil(); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceWafFalseAlarmMaskingRuleV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = falsealarmmasking_rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF False Alarm Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
