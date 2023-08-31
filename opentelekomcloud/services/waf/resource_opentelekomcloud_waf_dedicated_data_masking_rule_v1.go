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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDedicatedDataMaskingRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedDataMaskingRuleV1Create,
		ReadContext:   resourceWafDedicatedDataMaskingRuleV1Read,
		UpdateContext: resourceWafDedicatedDataMaskingRuleV1Update,
		DeleteContext: resourceWafDedicatedDataMaskingRuleV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("policy_id", "id"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"params", "cookie", "header", "form"},
					false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceWafDedicatedDataMaskingRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	createOpts := rules.CreatePrivacyOpts{
		Url:         d.Get("url").(string),
		Category:    d.Get("category").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.CreatePrivacy(client, policyID, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated Data Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf dedicated data masking rule created: %#v", rule)
	d.SetId(rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedDataMaskingRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedDataMaskingRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.GetPrivacy(client, policyID, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated Data Masking Rule: %s", err)
	}

	mErr := multierror.Append(
		d.Set("policy_id", rule.PolicyId),
		d.Set("url", rule.Url),
		d.Set("name", rule.Name),
		d.Set("category", rule.Category),
		d.Set("description", rule.Description),
		d.Set("status", rule.Status),
		d.Set("created_at", rule.CreatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedDataMaskingRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}
	policyId := d.Get("policy_id").(string)
	var updateOpts rules.UpdatePrivacyOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("url") {
		updateOpts.Url = d.Get("url").(string)
	}

	if d.HasChange("category") {
		updateOpts.Category = d.Get("category").(string)
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	_, err = rules.UpdatePrivacy(client, policyId, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF Dedicated Data Masking Rule: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedDataMaskingRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedDataMaskingRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeletePrivacyRule(client, policyID, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated Data Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
