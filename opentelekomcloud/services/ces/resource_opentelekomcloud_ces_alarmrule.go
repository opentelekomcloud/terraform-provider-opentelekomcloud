package ces

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
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
									regexp.MustCompile(`^[a-zA-Z].+`),
									"Must start with a letter."),
								validation.StringMatch(
									regexp.MustCompile(`^\w+$`),
									"Only lowercase/uppercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
								),
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
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(0),
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
			"insufficientdata_actions": {
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
				Default:  true,
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

func getMetricOpts(d *schema.ResourceData) alarmrule.MetricOpts {
	metricListRaw := d.Get("metric").([]interface{})
	metricElement := metricListRaw[0].(map[string]interface{})

	metricDimensions := metricElement["dimensions"].([]interface{})
	dimensionOpts := make([]alarmrule.DimensionOpts, len(metricDimensions))
	for i, dimensionElement := range metricDimensions {
		dimension := dimensionElement.(map[string]interface{})
		dimensionOpts[i] = alarmrule.DimensionOpts{
			Name:  dimension["name"].(string),
			Value: dimension["value"].(string),
		}
	}

	return alarmrule.MetricOpts{
		Namespace:  metricElement["namespace"].(string),
		MetricName: metricElement["metric_name"].(string),
		Dimensions: dimensionOpts,
	}
}

func getAlarmAction(d *schema.ResourceData, name string) []alarmrule.ActionOpts {
	alarmListRaw := d.Get(name).([]interface{})
	if len(alarmListRaw) == 0 {
		return nil
	}
	actionOpts := make([]alarmrule.ActionOpts, len(alarmListRaw))
	for i, alarmElement := range alarmListRaw {
		alarm := alarmElement.(map[string]interface{})

		notifyListRaw := alarm["notification_list"].([]interface{})
		notifyList := make([]string, len(notifyListRaw))
		for j, notifiedObject := range notifyListRaw {
			notifyList[j] = notifiedObject.(string)
		}

		actionOpts[i] = alarmrule.ActionOpts{
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

	metric := getMetricOpts(d)
	conditionListRaw := d.Get("condition").([]interface{})
	conditionElement := conditionListRaw[0].(map[string]interface{})
	createOpts := alarmrule.CreateOpts{
		AlarmName:        d.Get("alarm_name").(string),
		AlarmDescription: d.Get("alarm_description").(string),
		AlarmLevel:       d.Get("alarm_level").(int),
		Metric:           metric,
		Condition: alarmrule.ConditionOpts{
			Period:             conditionElement["period"].(int),
			Filter:             conditionElement["filter"].(string),
			ComparisonOperator: conditionElement["comparison_operator"].(string),
			Value:              conditionElement["value"].(int),
			Unit:               conditionElement["unit"].(string),
			Count:              conditionElement["count"].(int),
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
		return fmterr.Errorf("error creating alarm rule: %w", err)
	}
	log.Printf("[DEBUG] Created alarm rule: %#v", *r)

	d.SetId(r.AlarmID)

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

	r, err := alarmrule.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "alarmrule")
	}
	log.Printf("[DEBUG] Retrieved alarm rule %s: %#v", d.Id(), r)

	dimensionInfoList := make([]map[string]interface{}, len(r.Metric.Dimensions))
	for i, v := range r.Metric.Dimensions {
		dimensionInfoItem := map[string]interface{}{
			"name":  v.Name,
			"value": v.Value,
		}
		dimensionInfoList[i] = dimensionInfoItem
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
	for i, alarmActionItem := range r.AlarmActions {
		notificationList := make([]string, len(alarmActionItem.NotificationList))
		for j, notification := range alarmActionItem.NotificationList {
			notificationList[j] = notification
		}
		alarmAction := map[string]interface{}{
			"type":              alarmActionItem.Type,
			"notification_list": notificationList,
		}
		alarmActionsInfo[i] = alarmAction
	}

	insufficientActionsInfo := make([]map[string]interface{}, len(r.InsufficientdataActions))
	for i, insufficientActionItem := range r.InsufficientdataActions {
		notificationList := make([]string, len(insufficientActionItem.NotificationList))
		for j, notification := range insufficientActionItem.NotificationList {
			notificationList[j] = notification
		}
		insufficientAction := map[string]interface{}{
			"type":              insufficientActionItem.Type,
			"notification_list": notificationList,
		}
		insufficientActionsInfo[i] = insufficientAction
	}

	okActionsInfo := make([]map[string]interface{}, len(r.OkActions))
	for i, okActionItem := range r.OkActions {
		notificationList := make([]string, len(okActionItem.NotificationList))
		for j, notification := range okActionItem.NotificationList {
			notificationList[j] = notification
		}
		insufficientAction := map[string]interface{}{
			"type":              okActionItem.Type,
			"notification_list": notificationList,
		}
		okActionsInfo[i] = insufficientAction
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
		d.Set("alarm_enabled", r.AlarmEnabled),
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
