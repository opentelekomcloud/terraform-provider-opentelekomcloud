package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/datamasking_rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDataMaskingRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDataMaskingRuleV1Create,
		ReadContext:   resourceWafDataMaskingRuleV1Read,
		UpdateContext: resourceWafDataMaskingRuleV1Update,
		DeleteContext: resourceWafDataMaskingRuleV1Delete,
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

func resourceWafDataMaskingRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	wafClient, err := config.WafV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	createOpts := datamasking_rules.CreateOpts{
		Path:     d.Get("url").(string),
		Category: d.Get("category").(string),
		Index:    d.Get("index").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := datamasking_rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF DataMasking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf datamasking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafDataMaskingRuleV1Read(ctx, d, meta)
}

func resourceWafDataMaskingRuleV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	policyID := d.Get("policy_id").(string)
	n, err := datamasking_rules.Get(wafClient, policyID, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf DataMasking Rule: %s", err)
	}

	d.SetId(n.Id)
	mErr := multierror.Append(
		d.Set("url", n.Path),
		d.Set("category", n.Category),
		d.Set("index", n.Index),
		d.Set("policy_id", n.PolicyID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDataMaskingRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Client: %s", err)
	}
	var updateOpts datamasking_rules.UpdateOpts

	if d.HasChange("url") || d.HasChange("category") || d.HasChange("index") {
		updateOpts.Path = d.Get("url").(string)
		updateOpts.Category = d.Get("category").(string)
		updateOpts.Index = d.Get("index").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	if updateOpts != (datamasking_rules.UpdateOpts{}) {
		policyID := d.Get("policy_id").(string)
		_, err = datamasking_rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud WAF DataMasking Rule: %s", err)
		}
	}

	return resourceWafDataMaskingRuleV1Read(ctx, d, meta)
}

func resourceWafDataMaskingRuleV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = datamasking_rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF DataMasking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
