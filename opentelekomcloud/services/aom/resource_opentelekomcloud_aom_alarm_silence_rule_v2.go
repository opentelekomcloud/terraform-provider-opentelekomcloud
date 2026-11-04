package aom

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/aom/v2/alarm"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAlarmSilenceRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmSilenceRuleV2Create,
		ReadContext:   resourceAlarmSilenceRuleV2Read,
		UpdateContext: resourceAlarmSilenceRuleV2Update,
		DeleteContext: resourceAlarmSilenceRuleV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"match": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operate": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"EXIST", "REGEX", "EQUALS"}, false),
									},
									"value": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"mute_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"FIXED", "DAILY", "WEEKLY", "MONTHLY"}, false),
						},
						"starts_at": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"ends_at": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"scope": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSilenceRuleCreateOpts(d *schema.ResourceData) alarm.CreateSilenceRuleOpts {
	return alarm.CreateSilenceRuleOpts{
		Name:       d.Get("name").(string),
		Desc:       d.Get("description").(string),
		Timezone:   d.Get("timezone").(string),
		UserId:     d.Get("user_id").(string),
		Match:      expandSilenceMatch(d.Get("match").([]interface{})),
		MuteConfig: expandSilenceMuteConfig(d.Get("mute_config").([]interface{})),
	}
}

func resourceAlarmSilenceRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := buildSilenceRuleCreateOpts(d)
	if err = alarm.CreateSilenceRule(client, opts); err != nil {
		return diag.Errorf("error creating AOM alarm silence rule (%s): %s", opts.Name, err)
	}

	d.SetId(opts.Name)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAlarmSilenceRuleV2Read(clientCtx, d, meta)
}

func resourceAlarmSilenceRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	rule, err := GetAlarmSilenceRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AOM alarm silence rule")
	}

	mErr := multierror.Append(nil,
		d.Set("name", rule.Name),
		d.Set("description", rule.Desc),
		d.Set("timezone", rule.Timezone),
		d.Set("user_id", rule.UserId),
		d.Set("create_time", rule.CreateTime),
		d.Set("update_time", rule.UpdateTime),
		d.Set("match", flattenSilenceMatch(rule.Match)),
		d.Set("mute_config", flattenSilenceMuteConfig(rule.MuteConfig)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving AOM alarm silence rule (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceAlarmSilenceRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := alarm.UpdateSilenceRuleOpts{
		Name:       d.Get("name").(string),
		Desc:       d.Get("description").(string),
		Timezone:   d.Get("timezone").(string),
		UserId:     d.Get("user_id").(string),
		Match:      expandSilenceMatch(d.Get("match").([]interface{})),
		MuteConfig: expandSilenceMuteConfig(d.Get("mute_config").([]interface{})),
	}
	if err = alarm.UpdateSilenceRule(client, opts); err != nil {
		return diag.Errorf("error updating AOM alarm silence rule (%s): %s", opts.Name, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAlarmSilenceRuleV2Read(clientCtx, d, meta)
}

func resourceAlarmSilenceRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = alarm.DeleteSilenceRule(client, []alarm.DeleteSilenceRuleName{{Name: d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AOM alarm silence rule")
	}
	return nil
}

func expandSilenceMatch(raw []interface{}) [][]alarm.SilenceMatch {
	result := make([][]alarm.SilenceMatch, 0, len(raw))
	for _, groupRaw := range raw {
		group := groupRaw.(map[string]interface{})
		conditionsRaw := group["conditions"].([]interface{})
		conditions := make([]alarm.SilenceMatch, 0, len(conditionsRaw))
		for _, c := range conditionsRaw {
			item := c.(map[string]interface{})
			conditions = append(conditions, alarm.SilenceMatch{
				Key:     item["key"].(string),
				Operate: item["operate"].(string),
				Value:   common.ExpandToStringList(item["value"].([]interface{})),
			})
		}
		result = append(result, conditions)
	}
	return result
}

func flattenSilenceMatch(match [][]alarm.SilenceMatch) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(match))
	for _, conditions := range match {
		flatConditions := make([]map[string]interface{}, 0, len(conditions))
		for _, c := range conditions {
			flatConditions = append(flatConditions, map[string]interface{}{
				"key":     c.Key,
				"operate": c.Operate,
				"value":   c.Value,
			})
		}
		result = append(result, map[string]interface{}{
			"conditions": flatConditions,
		})
	}
	return result
}

func expandSilenceMuteConfig(raw []interface{}) alarm.SilenceMuteConfig {
	item := raw[0].(map[string]interface{})
	return alarm.SilenceMuteConfig{
		Type:     item["type"].(string),
		StartsAt: int64(item["starts_at"].(int)),
		EndsAt:   int64(item["ends_at"].(int)),
		Scope:    common.ExpandToIntList(item["scope"].([]interface{})),
	}
}

func flattenSilenceMuteConfig(config alarm.SilenceMuteConfig) []map[string]interface{} {
	return []map[string]interface{}{{
		"type":      config.Type,
		"starts_at": config.StartsAt,
		"ends_at":   config.EndsAt,
		"scope":     config.Scope,
	}}
}

// GetAlarmSilenceRule finds the silence rule by name in the rule list.
func GetAlarmSilenceRule(client *golangsdk.ServiceClient, name string) (*alarm.SilenceRule, error) {
	rules, err := alarm.ListSilenceRules(client)
	if err != nil {
		return nil, err
	}
	for i := range rules {
		if rules[i].Name == name {
			return &rules[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}
