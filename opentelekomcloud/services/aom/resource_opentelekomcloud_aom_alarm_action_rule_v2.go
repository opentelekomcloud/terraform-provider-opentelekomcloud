package aom

import (
	"context"
	"fmt"

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

const (
	keyClientV2         = "aom-v2-client"
	errCreationV2Client = "error creating OpenTelekomCloud AOM V2 client: %w"
)

func ResourceAlarmActionRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmActionRuleV2Create,
		ReadContext:   resourceAlarmActionRuleV2Read,
		UpdateContext: resourceAlarmActionRuleV2Update,
		DeleteContext: resourceAlarmActionRuleV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "3"}, false),
			},
			"notification_template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"smn_topic": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_urn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"push_policy": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1}),
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func buildActionRuleOpts(d *schema.ResourceData, config *cfg.Config) (alarm.CreateActionRuleOpts, error) {
	userName := d.Get("user_name").(string)
	if userName == "" {
		userName = config.Username
	}
	if userName == "" {
		return alarm.CreateActionRuleOpts{},
			fmt.Errorf("user_name must be specified for the AOM alarm action rule (%s), "+
				"and the provider configuration does not include an authenticated username", d.Id())
	}

	return alarm.CreateActionRuleOpts{
		RuleName:             d.Get("name").(string),
		ProjectId:            config.TenantID,
		UserName:             userName,
		Desc:                 d.Get("description").(string),
		Type:                 d.Get("type").(string),
		NotificationTemplate: d.Get("notification_template").(string),
		TimeZone:             d.Get("time_zone").(string),
		SmnTopics:            expandSmnTopics(d.Get("smn_topic").(*schema.Set).List()),
	}, nil
}

func resourceAlarmActionRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts, err := buildActionRuleOpts(d, config)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = alarm.CreateActionRule(client, opts); err != nil {
		return diag.Errorf("error creating AOM alarm action rule (%s): %s", opts.RuleName, err)
	}

	d.SetId(opts.RuleName)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAlarmActionRuleV2Read(clientCtx, d, meta)
}

func resourceAlarmActionRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	rule, err := GetAlarmActionRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AOM alarm action rule")
	}

	mErr := multierror.Append(nil,
		d.Set("name", rule.RuleName),
		d.Set("type", rule.Type),
		d.Set("notification_template", rule.NotificationTemplate),
		d.Set("description", rule.Desc),
		d.Set("user_name", rule.UserName),
		d.Set("time_zone", rule.TimeZone),
		d.Set("create_time", rule.CreateTime),
		d.Set("update_time", rule.UpdateTime),
		d.Set("smn_topic", flattenSmnTopics(rule.SmnTopics)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving AOM alarm action rule (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceAlarmActionRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts, err := buildActionRuleOpts(d, config)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = alarm.UpdateActionRule(client, alarm.UpdateActionRuleOpts(opts)); err != nil {
		return diag.Errorf("error updating AOM alarm action rule (%s): %s", opts.RuleName, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAlarmActionRuleV2Read(clientCtx, d, meta)
}

func resourceAlarmActionRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.AomV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = alarm.DeleteActionRule(client, []string{d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AOM alarm action rule")
	}
	return nil
}

func expandSmnTopics(raw []interface{}) []alarm.SmnTopic {
	topics := make([]alarm.SmnTopic, 0, len(raw))
	for _, v := range raw {
		item := v.(map[string]interface{})
		topics = append(topics, alarm.SmnTopic{
			DisplayName: item["display_name"].(string),
			Name:        item["name"].(string),
			PushPolicy:  item["push_policy"].(int),
			TopicUrn:    item["topic_urn"].(string),
		})
	}
	return topics
}

func flattenSmnTopics(topics []alarm.SmnTopic) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(topics))
	for _, t := range topics {
		result = append(result, map[string]interface{}{
			"topic_urn":    t.TopicUrn,
			"name":         t.Name,
			"push_policy":  t.PushPolicy,
			"display_name": t.DisplayName,
			"status":       t.Status,
		})
	}
	return result
}

// GetAlarmActionRule finds the alarm action rule by name in the rule list.
func GetAlarmActionRule(client *golangsdk.ServiceClient, name string) (*alarm.ActionRule, error) {
	rules, err := alarm.ListActionRules(client)
	if err != nil {
		return nil, err
	}
	for i := range rules {
		if rules[i].RuleName == name {
			return &rules[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}
