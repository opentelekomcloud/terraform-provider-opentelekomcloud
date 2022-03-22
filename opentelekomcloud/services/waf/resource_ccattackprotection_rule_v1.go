package waf

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/ccattackprotection_rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafCcAttackProtectionRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafCcAttackProtectionRuleV1Create,
		ReadContext:   resourceWafCcAttackProtectionRuleV1Read,
		DeleteContext: resourceWafCcAttackProtectionRuleV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("policy_id", "id"),
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
			"limit_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"limit_period": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"lock_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tag_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tag_index": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag_contents": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"action_category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"block_content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"block_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func getTagCondition(d *schema.ResourceData) ccattackprotection_rules.TagCondition {
	v := d.Get("tag_contents").([]interface{})
	contents := make([]string, len(v))
	for i, v := range v {
		contents[i] = v.(string)
	}

	condition := ccattackprotection_rules.TagCondition{
		Category: d.Get("tag_category").(string),
		Contents: contents,
	}

	log.Printf("[DEBUG] getTagCondition: %#v", condition)
	return condition
}

func getCcAction(d *schema.ResourceData) ccattackprotection_rules.Action {
	response := ccattackprotection_rules.Response{
		ContentType: d.Get("block_content_type").(string),
		Content:     d.Get("block_content").(string),
	}
	detail := ccattackprotection_rules.Detail{
		Response: response,
	}

	action := ccattackprotection_rules.Action{
		Category: d.Get("action_category").(string),
		Detail:   detail,
	}

	log.Printf("[DEBUG] getAction: %#v", action)
	return action
}

func resourceWafCcAttackProtectionRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	wafClient, err := config.WafV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud WAF Client: %s", err)
	}

	limitNum := d.Get("limit_num").(int)
	limitPeriod := d.Get("limit_period").(int)
	lockTime := d.Get("lock_time").(int)
	createOpts := ccattackprotection_rules.CreateOpts{
		Path:        d.Get("url").(string),
		LimitNum:    &limitNum,
		LimitPeriod: &limitPeriod,
		LockTime:    &lockTime,
		TagType:     d.Get("tag_type").(string),
		TagIndex:    d.Get("tag_index").(string),
		Action:      getCcAction(d),
	}

	_, tagCategoryOk := d.GetOk("tag_category")
	_, tagContentsOk := d.GetOk("tag_contents")
	if tagCategoryOk && tagContentsOk {
		createOpts.TagCondition = getTagCondition(d)
	}

	policyID := d.Get("policy_id").(string)
	rule, err := ccattackprotection_rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF CC Attack Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf cc attack protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafCcAttackProtectionRuleV1Read(ctx, d, meta)
}

func resourceWafCcAttackProtectionRuleV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}
	policyID := d.Get("policy_id").(string)
	n, err := ccattackprotection_rules.Get(wafClient, policyID, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Waf CC Attack Protection Rule: %s", err)
	}

	d.SetId(n.Id)
	mErr := multierror.Append(
		d.Set("policy_id", n.PolicyID),
		d.Set("url", n.Path),
		d.Set("limit_num", n.LimitNum),
		d.Set("limit_period", n.LimitPeriod),
		d.Set("lock_time", n.LockTime),
		d.Set("tag_type", n.TagType),
		d.Set("tag_index", n.TagIndex),
		d.Set("tag_category", n.TagCondition.Category),
		d.Set("tag_contents", n.TagCondition.Contents),
		d.Set("action_category", n.Action.Category),
		d.Set("block_content_type", n.Action.Detail.Response.ContentType),
		d.Set("block_content", n.Action.Detail.Response.Content),
		d.Set("default", n.Default),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafCcAttackProtectionRuleV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = ccattackprotection_rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud WAF CC Attack Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
