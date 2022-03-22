package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/whiteblackip_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafWhiteBlackIpRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafWhiteBlackIpRuleV1Create,
		ReadContext:   resourceWafWhiteBlackIpRuleV1Read,
		UpdateContext: resourceWafWhiteBlackIpRuleV1Update,
		DeleteContext: resourceWafWhiteBlackIpRuleV1Delete,
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
			"addr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"white": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
		},
	}
}

func resourceWafWhiteBlackIpRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	wafClient, err := config.WafV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := whiteblackip_rules.CreateOpts{
		Addr:  d.Get("addr").(string),
		White: d.Get("white").(int),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := whiteblackip_rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF WhiteBlackIP Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf whiteblackip rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafWhiteBlackIpRuleV1Read(ctx, d, meta)
}

func resourceWafWhiteBlackIpRuleV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := whiteblackip_rules.Get(wafClient, policyID, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf WhiteBlackIP Rule: %s", err)
	}

	d.SetId(n.Id)
	mErr := multierror.Append(
		d.Set("addr", n.Addr),
		d.Set("white", n.White),
		d.Set("policy_id", n.PolicyID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafWhiteBlackIpRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts whiteblackip_rules.UpdateOpts

	if d.HasChange("addr") {
		updateOpts.Addr = d.Get("addr").(string)
	}
	if d.HasChange("white") {
		white := d.Get("white").(int)
		updateOpts.White = &white
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	policyID := d.Get("policy_id").(string)
	_, err = whiteblackip_rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF WhiteBlackIP Rule: %s", err)
	}

	return resourceWafWhiteBlackIpRuleV1Read(ctx, d, meta)
}

func resourceWafWhiteBlackIpRuleV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = whiteblackip_rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF WhiteBlackIP Rule: %s", err)
	}

	d.SetId("")
	return nil
}
