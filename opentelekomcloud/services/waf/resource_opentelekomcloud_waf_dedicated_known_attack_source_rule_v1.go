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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafDedicatedKnownAttackSourceRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedKnownAttackSourceRuleV1Create,
		ReadContext:   resourceWafDedicatedKnownAttackSourceRuleV1Read,
		UpdateContext: resourceWafDedicatedKnownAttackSourceRuleV1Update,
		DeleteContext: resourceWafDedicatedKnownAttackSourceRuleV1Delete,
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
			"block_time": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1800),
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"long_ip_block", "long_cookie_block", "long_params_block",
						"short_ip_block", "short_cookie_block", "short_params_block"},
					false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceWafDedicatedKnownAttackSourceRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	createOpts := rules.CreateKnownAttackSourceOpts{
		BlockTime:   pointerto.Int(d.Get("block_time").(int)),
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.CreateKnownAttackSource(client, policyID, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated Known Attack Source Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf dedicated known attack source rule created: %#v", rule)
	d.SetId(rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedKnownAttackSourceRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedKnownAttackSourceRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.GetKnownAttackSource(client, policyID, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated Known Attack Source Rule: %s", err)
	}

	mErr := multierror.Append(
		d.Set("policy_id", rule.PolicyId),
		d.Set("category", rule.Category),
		d.Set("block_time", rule.BlockTime),
		d.Set("description", rule.Description),
		d.Set("created_at", rule.CreatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedKnownAttackSourceRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}
	policyId := d.Get("policy_id").(string)
	var updateOpts rules.UpdateKnownAttackSourceOpts

	if d.HasChange("block_time") {
		blockTime := d.Get("block_time").(int)
		updateOpts.BlockTime = &blockTime
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	_, err = rules.UpdateKnownAttackSource(client, policyId, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud WAF Dedicated Known Attack Source Rule: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedKnownAttackSourceRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedKnownAttackSourceRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteKnownAttackSourceRule(client, policyID, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated Known Attack Source Rule: %s", err)
	}

	d.SetId("")
	return nil
}
