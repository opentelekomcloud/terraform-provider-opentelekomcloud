package ces

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/alarms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRuleCreate,
		ReadContext:   resourceAlarmRuleRead,
		UpdateContext: resourceAlarmRuleUpdate,
		DeleteContext: resourceAlarmRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: checkCesAlarmRestrictions,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"alarm_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[\w-]+$`),
						"Only lowercase/uppercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed and must start with a letter.",
					),
				),
			},
			"alarm_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"alarm_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 4),
			},
			"alarm_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EVENT.SYS", "EVENT.CUSTOM",
				}, false),
			},
			"metric": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 64),
								validation.StringMatch(
									regexp.MustCompile(`^[a-zA-Z\][a-z0-9\/_]+[a-z0-9]+$`),
									"Only lowercase/uppercase letters, digits, underscores (_) and slashes (/) are allowed.",
								),
							),
						},
						"dimensions": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 32),
											validation.StringMatch(
												regexp.MustCompile(`^[a-zA-Z].+`),
												"Must start with a letter."),
											validation.StringMatch(
												regexp.MustCompile(`^[\w-]+$`),
												"Only lowercase/uppercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
											),
										),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 256),
											validation.StringMatch(
												regexp.MustCompile(`^[a-zA-Z0-9].+`),
												"Must start with a letter."),
											validation.StringMatch(
												regexp.MustCompile(`^[\w-]+$`),
												"Only lowercase/uppercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
											),
										),
									},
								},
							},
						},
					},
				},
			},
			"condition": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.IntInSlice([]int{
								1, 300, 1200, 3600, 14400, 86400,
							}),
						},
						"filter": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"max", "min", "average", "sum", "variance",
							}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								">", "=", "<", ">=", "<=",
							}, false),
						},
						"value": {
							Type:         schema.TypeFloat,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.FloatAtLeast(0),
						},
						"unit": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(0, 32),
						},
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},
						"alarm_frequency": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400,
							}),
						},
					},
				},
			},
			"alarm_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"notification", "autoscaling",
							}, false),
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ok_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"notification", "autoscaling",
							}, false),
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alarm_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getMetricOpts(d *schema.ResourceData) alarms.MetricForAlarm {
	metricListRaw := d.Get("metric").([]interface{})
	metricElement := metricListRaw[0].(map[string]interface{})

	metricDimensions := metricElement["dimensions"].([]interface{})
	dimensionOpts := make([]alarms.MetricsDimension, len(metricDimensions))
	for i, dimensionElement := range metricDimensions {
		dimension := dimensionElement.(map[string]interface{})
		dimensionOpts[i] = alarms.MetricsDimension{
			Name:  dimension["name"].(string),
			Value: dimension["value"].(string),
		}
	}
	metricName := metricElement["metric_name"].(string)
	metricName = strings.ReplaceAll(metricName, "/", "SlAsH")
	return alarms.MetricForAlarm{
		Namespace:  metricElement["namespace"].(string),
		MetricName: metricName,
		Dimensions: dimensionOpts,
	}
}

func getAlarmAction(d *schema.ResourceData, name string) []alarms.AlarmActions {
	alarmListRaw := d.Get(name).([]interface{})
	if len(alarmListRaw) == 0 {
		return nil
	}
	actionOpts := make([]alarms.AlarmActions, len(alarmListRaw))
	for i, alarmElement := range alarmListRaw {
		alarm := alarmElement.(map[string]interface{})

		notifyListRaw := alarm["notification_list"].([]interface{})
		notifyList := make([]string, len(notifyListRaw))
		for j, notifiedObject := range notifyListRaw {
			notifyList[j] = notifiedObject.(string)
		}

		actionOpts[i] = alarms.AlarmActions{
			Type:             alarm["type"].(string),
			NotificationList: notifyList,
		}
	}
	return actionOpts
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	alarmEnabled := d.Get("alarm_enabled").(bool)
	alarmActionEnabled := d.Get("alarm_action_enabled").(bool)

	metric := getMetricOpts(d)
	conditionListRaw := d.Get("condition").([]interface{})
	if len(conditionListRaw) == 0 {
		return fmterr.Errorf("invalid `condition` field provided in configuration")
	}

	conditionElement := conditionListRaw[0].(map[string]interface{})
	createOpts := alarms.CreateAlarmOpts{
		AlarmName:          d.Get("alarm_name").(string),
		AlarmType:          d.Get("alarm_type").(string),
		AlarmDescription:   d.Get("alarm_description").(string),
		AlarmLevel:         d.Get("alarm_level").(int),
		AlarmActions:       getAlarmAction(d, "alarm_actions"),
		AlarmEnabled:       &alarmEnabled,
		AlarmActionEnabled: &alarmActionEnabled,
		OkActions:          getAlarmAction(d, "ok_actions"),
		Metric:             metric,
		Condition: alarms.Condition{
			Period:             conditionElement["period"].(int),
			Filter:             conditionElement["filter"].(string),
			ComparisonOperator: conditionElement["comparison_operator"].(string),
			Value:              conditionElement["value"].(float64),
			Unit:               conditionElement["unit"].(string),
			Count:              conditionElement["count"].(int),
			SuppressDuration:   conditionElement["alarm_frequency"].(int),
		},
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	alarmId, err := alarms.CreateAlarm(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating alarm rule: %w", err)
	}
	log.Printf("[DEBUG] Created alarm rule, ID: %#v", alarmId)

	d.SetId(alarmId)

	clientCtx := common.CtxWithClient(ctx, client, cesClientV1)
	return resourceAlarmRuleRead(clientCtx, d, meta)
}

func resourceAlarmRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	respAlarm, err := alarms.ShowAlarm(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "alarmrule")
	}
	alarm := respAlarm[0]

	log.Printf("[DEBUG] Retrieved alarm rule %s: %#v", d.Id(), alarm)

	dimensionInfoList := make([]map[string]interface{}, len(alarm.Metric.Dimensions))
	for i, v := range alarm.Metric.Dimensions {
		dimensionInfoItem := map[string]interface{}{
			"name":  v.Name,
			"value": v.Value,
		}
		dimensionInfoList[i] = dimensionInfoItem
	}
	metricName := alarm.Metric.MetricName
	metricName = strings.ReplaceAll(metricName, "SlAsH", "/")
	metricInfo := []map[string]interface{}{
		{
			"namespace":   alarm.Metric.Namespace,
			"metric_name": metricName,
			"dimensions":  dimensionInfoList,
		},
	}

	conditionInfo := []map[string]interface{}{
		{
			"period":              alarm.Condition.Period,
			"filter":              alarm.Condition.Filter,
			"comparison_operator": alarm.Condition.ComparisonOperator,
			"value":               alarm.Condition.Value,
			"unit":                alarm.Condition.Unit,
			"count":               alarm.Condition.Count,
			"alarm_frequency":     alarm.Condition.SuppressDuration,
		},
	}

	alarmActionsInfo := make([]map[string]interface{}, len(alarm.AlarmActions))
	for i, alarmActionItem := range alarm.AlarmActions {
		notificationList := make([]string, len(alarmActionItem.NotificationList))
		copy(notificationList, alarmActionItem.NotificationList)
		alarmAction := map[string]interface{}{
			"type":              alarmActionItem.Type,
			"notification_list": notificationList,
		}
		alarmActionsInfo[i] = alarmAction
	}

	okActionsInfo := make([]map[string]interface{}, len(alarm.OkActions))
	for i, okActionItem := range alarm.OkActions {
		notificationList := make([]string, len(okActionItem.NotificationList))
		copy(notificationList, okActionItem.NotificationList)
		insufficientAction := map[string]interface{}{
			"type":              okActionItem.Type,
			"notification_list": notificationList,
		}
		okActionsInfo[i] = insufficientAction
	}

	mErr := multierror.Append(
		d.Set("alarm_name", alarm.AlarmName),
		d.Set("alarm_type", alarm.AlarmType),
		d.Set("alarm_description", alarm.AlarmDescription),
		d.Set("alarm_level", alarm.AlarmLevel),
		d.Set("metric", metricInfo),
		d.Set("condition", conditionInfo),
		d.Set("alarm_actions", alarmActionsInfo),
		d.Set("ok_actions", okActionsInfo),
		d.Set("alarm_enabled", &alarm.AlarmEnabled),
		d.Set("alarm_action_enabled", &alarm.AlarmActionEnabled),
		d.Set("update_time", alarm.UpdateTime),
		d.Set("alarm_state", alarm.AlarmState),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	alarmRuleID := d.Id()

	if !d.HasChange("alarm_enabled") {
		log.Printf("[WARN] Nothing will be updated")
		return nil
	}
	updateOpts := alarms.ModifyAlarmActionRequest{AlarmEnabled: d.Get("alarm_enabled").(bool)}
	log.Printf("[DEBUG] Updating %s with options: %#v", alarmRuleID, updateOpts)

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := alarms.UpdateAlarmAction(client, alarmRuleID, updateOpts)
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error updating alarm rule %s: %w", alarmRuleID, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV1)
	return resourceAlarmRuleRead(clientCtx, d, meta)
}

func resourceAlarmRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	alarmRuleID := d.Id()
	log.Printf("[DEBUG] Deleting alarm rule %s", alarmRuleID)

	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := alarms.DeleteAlarm(client, alarmRuleID)
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if common.IsResourceNotFound(err) {
			log.Printf("[INFO] Deleting an unavailable: %s", alarmRuleID)
			return nil
		}
		return fmterr.Errorf("error deleting alarm rule %s: %s", alarmRuleID, err)
	}

	return nil
}

func checkCesAlarmRestrictions(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if v, ok := d.GetOk("alarm_action_enabled"); ok {
		alarmActionEnabled := v.(bool)
		if alarmActionEnabled {
			_, alarmCheck := d.GetOk("alarm_actions")
			_, okCheck := d.GetOk("ok_actions")
			if !(alarmCheck || okCheck) {
				return fmt.Errorf("either `alarm_actions` or `ok_actions` should be specified when `alarm_action_enabled` set to `true`")
			}
		}
	}
	return nil
}
