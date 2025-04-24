package lts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/alarm"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceKeywordsAlarmRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeywordsAlarmRuleV2Create,
		UpdateContext: resourceKeywordsAlarmRuleV2Update,
		ReadContext:   resourceKeywordsAlarmRuleV2Read,
		DeleteContext: resourceKeywordsAlarmRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"keywords_requests": {
				Type:     schema.TypeList,
				Elem:     keywordsRequestsSchema(),
				Required: true,
			},
			"frequency": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     frequencySchema(),
				Required: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_frequency": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"alarm_action_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"send_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"notification_rule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     notificationRuleSchema(),
				Optional: true,
				Computed: true,
			},
			"trigger_condition_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"trigger_condition_frequency": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"send_recovery_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"recovery_policy": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func keywordsRequestsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_time_range_unit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_time_range": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
	return &sc
}

func frequencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cron_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hour_of_day": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"day_of_week": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"fixed_rate_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func notificationRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topics": {
				Type:     schema.TypeList,
				Elem:     topicSchema(),
				Required: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func topicSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"push_policy": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceKeywordsAlarmRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := alarm.CreateOpts{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		Details:                   buildKeywordsRequests(d.Get("keywords_requests")),
		Frequency:                 buildFrequency(d.Get("frequency")),
		Severity:                  d.Get("severity").(string),
		Send:                      pointerto.Bool(d.Get("send_notifications").(bool)),
		DomainId:                  client.DomainID,
		TriggerConditionCount:     d.Get("trigger_condition_count").(int),
		TriggerConditionFrequency: d.Get("trigger_condition_frequency").(int),
		EnableRecoveryPolicy:      pointerto.Bool(d.Get("send_recovery_notifications").(bool)),
		RecoveryPolicy:            d.Get("recovery_policy").(int),
		NotificationFrequency:     d.Get("notification_frequency").(int),
		AlarmActionRuleName:       d.Get("alarm_action_rule_name").(string),
		NotificationSave:          buildNotificationRule(d.Get("notification_rule")),
	}
	rule, err := alarm.CreateKeywordRule(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS v2 Alarm Keyword Rule: %s", err)
	}
	d.SetId(rule)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceKeywordsAlarmRuleV2Read(clientCtx, d, meta)
}

func buildKeywordsRequests(rawParams interface{}) []alarm.Details {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]alarm.Details, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = alarm.Details{
				Keyword:             raw["keyword"].(string),
				Condition:           raw["condition"].(string),
				Number:              raw["number"].(int),
				LogStreamId:         raw["log_stream_id"].(string),
				LogGroupId:          raw["log_group_id"].(string),
				SearchTimeRangeUnit: raw["search_time_range_unit"].(string),
				SearchTimeRange:     raw["search_time_range"].(int),
			}
		}
		return rst
	}
	return nil
}

func buildFrequency(rawParams interface{}) *alarm.Frequency {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := alarm.Frequency{
			Type:          raw["type"].(string),
			CronExpr:      raw["cron_expression"].(string),
			HourOfDay:     raw["hour_of_day"].(int),
			DayOfWeek:     raw["day_of_week"].(int),
			FixedRateUnit: raw["fixed_rate_unit"].(string),
			FixedRate:     raw["fixed_rate"].(int),
		}
		return &params
	}
	return nil
}

func buildNotificationRule(rawParams interface{}) *alarm.NotificationSave {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := alarm.NotificationSave{
			TemplateName: raw["template_name"].(string),
			UserName:     raw["user_name"].(string),
			Topics:       buildTopic(raw["topics"]),
			Timezone:     raw["timezone"].(string),
			Language:     raw["language"].(string),
		}
		return &params
	}
	return nil
}

func buildTopic(rawParams interface{}) []alarm.TopicsCreate {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]alarm.TopicsCreate, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = alarm.TopicsCreate{
				Name:        raw["name"].(string),
				TopicUrn:    raw["topic_urn"].(string),
				DisplayName: raw["display_name"].(string),
				PushPolicy:  raw["push_policy"].(int),
			}
		}
		return rst
	}
	return nil
}

func resourceKeywordsAlarmRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	requestResp, err := alarm.ListKeywordRules(client)
	if err != nil {
		return diag.FromErr(err)
	}
	var keywordResult alarm.KeywordRule
	for _, kw := range requestResp {
		if kw.ID == d.Id() {
			keywordResult = kw
			break
		}
	}
	if keywordResult.ID == "" {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find OpenTelekomCloud LTS v2 cce access config by its ID (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", keywordResult.Name),
		d.Set("description", keywordResult.Description),
		d.Set("keywords_requests", flattenKeywordsRequests(keywordResult.Details)),
		d.Set("frequency", flattenFrequency(keywordResult.Frequency)),
		d.Set("severity", keywordResult.Severity),
		d.Set("send_notifications", keywordResult.Send),
		d.Set("domain_id", keywordResult.DomainId),
		d.Set("trigger_condition_count", keywordResult.TriggerConditionCount),
		d.Set("trigger_condition_frequency", keywordResult.TriggerConditionFrequency),
		d.Set("send_recovery_notifications", keywordResult.Send),
		d.Set("recovery_policy", keywordResult.RecoveryPolicy),
		d.Set("alarm_action_rule_name", keywordResult.AlarmActionRuleName),
		d.Set("notification_frequency", keywordResult.NotificationFrequency),
		d.Set("status", keywordResult.Status),
		d.Set("created_at", common.FormatTimeStampRFC3339(keywordResult.CreatedAt/1000, false)),
		d.Set("updated_at", common.FormatTimeStampRFC3339(keywordResult.UpdatedAt/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKeywordsRequests(resp []alarm.Details) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"keyword":                v.Keyword,
			"condition":              v.Condition,
			"number":                 v.Number,
			"log_stream_id":          v.LogStreamId,
			"log_group_id":           v.LogGroupId,
			"search_time_range_unit": v.SearchTimeRangeUnit,
			"search_time_range":      v.SearchTimeRange,
		})
	}
	return rst
}

func flattenFrequency(resp *alarm.Frequency) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"type":            resp.Type,
			"cron_expression": resp.CronExpr,
			"hour_of_day":     resp.HourOfDay,
			"day_of_week":     resp.DayOfWeek,
			"fixed_rate_unit": resp.FixedRateUnit,
			"fixed_rate":      resp.FixedRate,
		},
	}
}

func resourceKeywordsAlarmRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	updateKeywordsAlarmRuleChanges := []string{
		"description",
		"keywords_requests",
		"frequency",
		"severity",
		"trigger_condition_count",
		"trigger_condition_frequency",
		"send_recovery_notifications",
		"recovery_frequency",
		"notification_rule",
		"alarm_action_rule_name",
		"notification_frequency",
		"send_notifications",
	}

	if d.HasChanges(updateKeywordsAlarmRuleChanges...) {
		_, err = alarm.UpdateKeywordRule(client, alarm.UpdateOpts{
			ID:                        d.Id(),
			Name:                      d.Get("name").(string),
			Description:               d.Get("description").(string),
			Details:                   buildKeywordsRequests(d.Get("keywords_requests")),
			Frequency:                 buildFrequency(d.Get("frequency")),
			Severity:                  d.Get("severity").(string),
			Send:                      pointerto.Bool(d.Get("send_notifications").(bool)),
			DomainId:                  client.DomainID,
			TriggerConditionCount:     d.Get("trigger_condition_count").(int),
			TriggerConditionFrequency: d.Get("trigger_condition_frequency").(int),
			EnableRecoveryPolicy:      pointerto.Bool(d.Get("send_recovery_notifications").(bool)),
			RecoveryPolicy:            d.Get("recovery_policy").(int),
			NotificationFrequency:     d.Get("notification_frequency").(int),
			AlarmActionRuleName:       d.Get("alarm_action_rule_name").(string),
			NotificationSave:          buildNotificationRule(d.Get("notification_rule")),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceKeywordsAlarmRuleV2Read(clientCtx, d, meta)
}

func resourceKeywordsAlarmRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.LtsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = alarm.DeleteKeywordRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud LTS v2 Alarm Keyword Rule")
	}
	return nil
}
