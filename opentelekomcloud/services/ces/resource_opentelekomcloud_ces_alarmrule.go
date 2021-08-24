package ces

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cloudeyeservice/alarmrule"
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
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{0,63}$"),
					"`alarm_name` must be a string of 1 to 64 characters that starts with a letter or digit and consists of uppercase/lowercase letters, digits and hyphens(-)",
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
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{0,63}$"),
								"`alarm_name` must be a string of 1 to 64 characters that starts with a letter or digit and consists of uppercase/lowercase letters, digits and hyphens(-)",
							),
						},
						"dimensions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,31}$"),
											"`name` must be a string of 1 to 32 characters that starts with a letter and consists of uppercase/lowercase letters, digits and underscores(_)",
										),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{0,255}$"),
											"`value` must be a string of 1 to 256 characters that starts with a letter or digit and consists of uppercase/lowercase letters, digits and hyphens(-)",
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
							ValidateFunc: validation.IntInSlice([]int{
								1, 300, 1200, 3600, 14400, 86400,
							}),
						},
						"filter": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"max", "min", "average", "sum", "variance",
							}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								">", "=", "<", ">=", "<=",
							}, false),
						},
						"value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"unit": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
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
							ValidateFunc: validation.StringInSlice([]string{
								"notification", "autoscaling",
							}, false),
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"insufficientdata_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"notification", "autoscaling",
							}, false),
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
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
							ValidateFunc: validation.StringInSlice([]string{
								"notification", "autoscaling",
							}, false),
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
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
				Default:  true,
				ForceNew: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"alarm_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func getMetricOpts(d *schema.ResourceData) (alarmrule.MetricOpts, error) {
	mos, ok := d.Get("metric").([]interface{})
	if !ok {
		return alarmrule.MetricOpts{}, fmt.Errorf("error converting opt of metric:%v", d.Get("metric"))
	}
	mo := mos[0].(map[string]interface{})

	mod := mo["dimensions"].([]interface{})
	dopts := make([]alarmrule.DimensionOpts, len(mod))
	for i, v := range mod {
		v1 := v.(map[string]interface{})
		dopts[i] = alarmrule.DimensionOpts{
			Name:  v1["name"].(string),
			Value: v1["value"].(string),
		}
	}
	return alarmrule.MetricOpts{
		Namespace:  mo["namespace"].(string),
		MetricName: mo["metric_name"].(string),
		Dimensions: dopts,
	}, nil
}

func getAlarmAction(d *schema.ResourceData, name string) []alarmrule.ActionOpts {
	aos := d.Get(name).([]interface{})
	if len(aos) == 0 {
		return nil
	}
	opts := make([]alarmrule.ActionOpts, len(aos))
	for i, v := range aos {
		v1 := v.(map[string]interface{})

		v2 := v1["notification_list"].([]interface{})
		nl := make([]string, len(v2))
		for j, v3 := range v2 {
			nl[j] = v3.(string)
		}

		opts[i] = alarmrule.ActionOpts{
			Type:             v1["type"].(string),
			NotificationList: nl,
		}
	}
	return opts
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	metric, err := getMetricOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	cos := d.Get("condition").([]interface{})
	co := cos[0].(map[string]interface{})
	createOpts := alarmrule.CreateOpts{
		AlarmName:        d.Get("alarm_name").(string),
		AlarmDescription: d.Get("alarm_description").(string),
		AlarmLevel:       d.Get("alarm_level").(int),
		Metric:           metric,
		Condition: alarmrule.ConditionOpts{
			Period:             co["period"].(int),
			Filter:             co["filter"].(string),
			ComparisonOperator: co["comparison_operator"].(string),
			Value:              co["value"].(int),
			Unit:               co["unit"].(string),
			Count:              co["count"].(int),
		},
		AlarmActions:            getAlarmAction(d, "alarm_actions"),
		InsufficientdataActions: getAlarmAction(d, "insufficientdata_actions"),
		OkActions:               getAlarmAction(d, "ok_actions"),
		AlarmEnabled:            d.Get("alarm_enabled").(bool),
		AlarmActionEnabled:      d.Get("alarm_action_enabled").(bool),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	r, err := alarmrule.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating: %s", err)
	}
	log.Printf("[DEBUG] Create: %#v", *r)

	d.SetId(r.AlarmID)

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	r, err := alarmrule.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "alarmrule"))
	}
	log.Printf("[DEBUG] Retrieved alarm rule %s: %#v", d.Id(), r)

	dimensionInfoList := make([]map[string]interface{}, len(r.Metric.Dimensions))
	for _, v := range r.Metric.Dimensions {
		dimensionInfoItem := map[string]interface{}{
			"name":  v.Name,
			"value": v.Value,
		}
		dimensionInfoList = append(dimensionInfoList, dimensionInfoItem)
	}
	metricInfo := []map[string]interface{}{
		{
			"namespace":   r.Metric.Namespace,
			"metric_name": r.Metric.MetricName,
			"dimensions":  dimensionInfoList,
		},
	}

	conditionInfo := []map[string]interface{}{
		{
			"period":              r.Condition.Period,
			"filter":              r.Condition.Filter,
			"comparison_operator": r.Condition.ComparisonOperator,
			"value":               r.Condition.Value,
			"unit":                r.Condition.Unit,
			"count":               r.Condition.Count,
		},
	}

	alarmActionsInfo := make([]map[string]interface{}, len(r.AlarmActions))
	for _, alarmActionItem := range r.AlarmActions {
		notificationList := make([]string, len(alarmActionItem.NotificationList))
		for _, notification := range alarmActionItem.NotificationList {
			notificationList = append(notificationList, notification)
		}
		alarmAction := map[string]interface{}{
			"type":              alarmActionItem.Type,
			"notification_list": notificationList,
		}
		alarmActionsInfo = append(alarmActionsInfo, alarmAction)
	}

	insufficientActionsInfo := make([]map[string]interface{}, len(r.InsufficientdataActions))
	for _, insufficientActionItem := range r.InsufficientdataActions {
		notificationList := make([]string, len(insufficientActionItem.NotificationList))
		for _, notification := range insufficientActionItem.NotificationList {
			notificationList = append(notificationList, notification)
		}
		insufficientAction := map[string]interface{}{
			"type":              insufficientActionItem.Type,
			"notification_list": notificationList,
		}
		insufficientActionsInfo = append(insufficientActionsInfo, insufficientAction)
	}

	okActionsInfo := make([]map[string]interface{}, len(r.OkActions))
	for _, okActionItem := range r.OkActions {
		notificationList := make([]string, len(okActionItem.NotificationList))
		for _, notification := range okActionItem.NotificationList {
			notificationList = append(notificationList, notification)
		}
		insufficientAction := map[string]interface{}{
			"type":              okActionItem.Type,
			"notification_list": notificationList,
		}
		okActionsInfo = append(okActionsInfo, insufficientAction)
	}

	mErr := multierror.Append(
		d.Set("alarm_name", r.AlarmName),
		d.Set("alarm_description", r.AlarmDescription),
		d.Set("alarm_level", r.AlarmLevel),
		d.Set("metric", metricInfo),
		d.Set("condition", conditionInfo),
		d.Set("alarm_actions", alarmActionsInfo),
		d.Set("insufficientdata_actions", insufficientActionsInfo),
		d.Set("ok_actions", okActionsInfo),
		d.Set("alarm_enabled", r.AlarmActionEnabled),
		d.Set("alarm_action_enabled", r.AlarmActionEnabled),
		d.Set("update_time", r.UpdateTime),
		d.Set("alarm_state", r.AlarmState),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	alarmRuleID := d.Id()

	if !d.HasChange("alarm_enabled") {
		log.Printf("[WARN] Nothing will be updated")
		return nil
	}
	updateOpts := alarmrule.UpdateOpts{AlarmEnabled: d.Get("alarm_enabled").(bool)}
	log.Printf("[DEBUG] Updating %s with options: %#v", alarmRuleID, updateOpts)

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := alarmrule.Update(client, alarmRuleID, updateOpts).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error updating alarm rule %s: %s", alarmRuleID, err)
	}

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	alarmRuleID := d.Id()
	log.Printf("[DEBUG] Deleting alarm rule %s", alarmRuleID)

	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := alarmrule.Delete(client, alarmRuleID).ExtractErr()
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
