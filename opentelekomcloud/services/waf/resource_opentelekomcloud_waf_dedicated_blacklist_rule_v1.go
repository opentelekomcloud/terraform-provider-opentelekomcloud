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

func ResourceWafDedicatedBlacklistRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedBlacklistRuleV1Create,
		ReadContext:   resourceWafDedicatedBlacklistRuleV1Read,
		DeleteContext: resourceWafDedicatedBlacklistRuleV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("policy_id", "id"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					ProtectionActionBlock, ProtectionActionAllow, ProtectionActionLog,
				}),
			},
			"followed_action_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceWafDedicatedBlacklistRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	createOpts := rules.BlacklistCreateOpts{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Addresses:        d.Get("ip_address").(string),
		Action:           pointerto.Int(d.Get("action").(int)),
		FollowedActionId: d.Get("followed_action_id").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.CreateBlacklist(client, policyID, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF Dedicated Blacklist or Whitelist Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf dedicated blacklist or whitelist rule created: %#v", rule)
	d.SetId(rule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceWafDedicatedBlacklistRuleV1Read(clientCtx, d, meta)
}

func resourceWafDedicatedBlacklistRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.GetBlacklist(client, policyID, d.Id())

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf Dedicated Blacklist or Whitelist Rule: %s", err)
	}

	mErr := multierror.Append(
		d.Set("policy_id", rule.PolicyId),
		d.Set("name", rule.Name),
		d.Set("action", rule.Action),
		d.Set("ip_address", rule.Addresses),
		d.Set("description", rule.Description),
		d.Set("followed_action_id", rule.FollowedActionId),
		d.Set("status", rule.Status),
		d.Set("created_at", rule.CreatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafDedicatedBlacklistRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteBlacklistRule(client, policyID, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF Dedicated Blacklist or Whitelist Rule: %s", err)
	}

	d.SetId("")
	return nil
}
