package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/webtamperprotection_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafWebTamperProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafWebTamperProtectionRuleV1Create,
		ReadContext:   resourceWafWebTamperProtectionRuleV1Read,
		DeleteContext: resourceWafWebTamperProtectionRuleV1Delete,
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

func resourceWafWebTamperProtectionRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	wafClient, err := config.WafV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := webtamperprotection_rules.CreateOpts{
		Hostname: d.Get("hostname").(string),
		Path:     d.Get("url").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := webtamperprotection_rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Web Tamper Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf web tamper protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafWebTamperProtectionRuleV1Read(ctx, d, meta)
}

func resourceWafWebTamperProtectionRuleV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	policyID := d.Get("policy_id").(string)
	n, err := webtamperprotection_rules.Get(wafClient, policyID, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Web Tamper Protection Rule: %s", err)
	}

	d.SetId(n.Id)
	mErr := multierror.Append(
		d.Set("hostname", n.Hostname),
		d.Set("url", n.Path),
		d.Set("policy_id", n.PolicyID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafWebTamperProtectionRuleV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = webtamperprotection_rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Web Tamper Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
