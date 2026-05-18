package ces

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/alarms"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/policies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/resources"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceAlarmRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRuleV2Create,
		ReadContext:   resourceAlarmRuleV2Read,
		UpdateContext: resourceAlarmRuleV2Update,
		DeleteContext: resourceAlarmRuleV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[\w-]+$`),
						"Only letters, digits, underscores (_), and hyphens (-) are allowed.",
					),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alarm_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"alarm_template_id", "policies"},
			},
			"policies": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				Set:          resourcePolicyHash,
				ExactlyOneOf: []string{"alarm_template_id", "policies"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 1, 300, 1200, 3600, 14400, 86400,
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
								">", "=", "<", ">=", "<=", "!=",
							}, false),
						},
						"value": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"unit": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 32),
						},
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},
						"suppress_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400,
							}),
						},
						"level": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2,
							ValidateFunc: validation.IntBetween(1, 4),
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimensions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MULTI_INSTANCE",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALL_INSTANCE", "MULTI_INSTANCE", "RESOURCE_GROUP", "EVENT.SYS", "EVENT.CUSTOM",
				}, false),
			},
			"alarm_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
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
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
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
			"notification_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notification_begin_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"notification_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePolicyHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["metric_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_name"].(string)))
	}
	if m["period"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["period"].(int)))
	}
	if m["filter"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["filter"].(string)))
	}
	if m["comparison_operator"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["comparison_operator"].(string)))
	}
	if m["value"] != nil {
		buf.WriteString(fmt.Sprintf("%f-", m["value"].(float64)))
	}
	if m["count"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["count"].(int)))
	}
	if m["unit"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["unit"].(string)))
	}
	if m["suppress_duration"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["suppress_duration"].(int)))
	}
	if m["level"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["level"].(int)))
	}

	return hashcode.String(buf.String())
}

func buildPoliciesOptsV2(d *schema.ResourceData) []alarms.Policy {
	rawPolicies := d.Get("policies").(*schema.Set).List()
	if len(rawPolicies) == 0 {
		return nil
	}

	policiesOpts := make([]alarms.Policy, len(rawPolicies))
	for i, v := range rawPolicies {
		policy := v.(map[string]interface{})
		policiesOpts[i] = alarms.Policy{
			MetricName:         policy["metric_name"].(string),
			Period:             policy["period"].(int),
			Filter:             policy["filter"].(string),
			ComparisonOperator: policy["comparison_operator"].(string),
			Value:              policy["value"].(float64),
			Unit:               policy["unit"].(string),
			Count:              policy["count"].(int),
			SuppressDuration:   policy["suppress_duration"].(int),
			Level:              policy["level"].(int),
		}
	}
	return policiesOpts
}

func buildResourcesOptsV2(d *schema.ResourceData) [][]alarms.Dimension {
	rawResources := d.Get("resources").(*schema.Set).List()
	if len(rawResources) == 0 {
		resourceType := d.Get("type").(string)
		if resourceType == "EVENT.SYS" || resourceType == "EVENT.CUSTOM" || resourceType == "RESOURCE_GROUP" {
			return [][]alarms.Dimension{}
		}
		return nil
	}

	resourcesOpts := make([][]alarms.Dimension, len(rawResources))
	for i, v := range rawResources {
		resource := v.(map[string]interface{})
		dimensionsRaw, ok := resource["dimensions"].([]interface{})
		if !ok || len(dimensionsRaw) == 0 {
			resourcesOpts[i] = []alarms.Dimension{}
			continue
		}

		dimensions := make([]alarms.Dimension, len(dimensionsRaw))
		for j, dim := range dimensionsRaw {
			dimension := dim.(map[string]interface{})
			dimensions[j] = alarms.Dimension{
				Name:  dimension["name"].(string),
				Value: dimension["value"].(string),
			}
		}
		resourcesOpts[i] = dimensions
	}
	return resourcesOpts
}

func buildNotificationsOpts(d *schema.ResourceData, key string) []alarms.Notification {
	rawNotifications := d.Get(key).([]interface{})
	if len(rawNotifications) == 0 {
		return nil
	}

	notifications := make([]alarms.Notification, len(rawNotifications))
	for i, v := range rawNotifications {
		notification := v.(map[string]interface{})
		notifyListRaw := notification["notification_list"].([]interface{})
		notifyList := make([]string, len(notifyListRaw))
		for j, n := range notifyListRaw {
			notifyList[j] = n.(string)
		}
		notifications[i] = alarms.Notification{
			Type:             notification["type"].(string),
			NotificationList: notifyList,
		}
	}
	return notifications
}

func resourceAlarmRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	alarmEnabled := d.Get("alarm_enabled").(bool)
	notificationEnabled := d.Get("notification_enabled").(bool)

	createOpts := alarms.CreateOpts{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Namespace:             d.Get("namespace").(string),
		ResourceGroupId:       d.Get("resource_group_id").(string),
		AlarmTemplateId:       d.Get("alarm_template_id").(string),
		Policies:              buildPoliciesOptsV2(d),
		Resources:             buildResourcesOptsV2(d),
		Type:                  d.Get("type").(string),
		NotificationEnabled:   &notificationEnabled,
		Enabled:               &alarmEnabled,
		NotificationBeginTime: d.Get("notification_begin_time").(string),
		NotificationEndTime:   d.Get("notification_end_time").(string),
		EnterpriseProjectId:   d.Get("enterprise_project_id").(string),
	}

	if notificationEnabled {
		createOpts.AlarmNotifications = buildNotificationsOpts(d, "alarm_actions")
		createOpts.OkNotifications = buildNotificationsOpts(d, "ok_actions")
	}

	log.Printf("[DEBUG] CES Alarm Rule V2 Create Options: %#v", createOpts)

	alarmId, err := alarms.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CES Alarm Rule V2: %w", err)
	}

	log.Printf("[DEBUG] CES Alarm Rule V2 created with ID: %s", alarmId)
	d.SetId(alarmId)

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceAlarmRuleV2Read(clientCtx, d, meta)
}

func resourceAlarmRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	listOpts := alarms.ListOpts{
		AlarmId: d.Id(),
	}
	listResp, err := alarms.List(client, listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CES Alarm Rule V2")
	}

	if listResp == nil || len(listResp.Alarms) == 0 {
		log.Printf("[WARN] CES Alarm Rule V2 (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	alarm := listResp.Alarms[0]

	// Get policies
	policiesResp, err := policies.List(client, d.Id(), policies.ListOpts{})
	if err != nil {
		return fmterr.Errorf("error retrieving CES Alarm Rule V2 policies: %w", err)
	}

	// Get resources
	resourcesResp, err := resources.List(client, d.Id(), resources.ListOpts{})
	if err != nil {
		return fmterr.Errorf("error retrieving CES Alarm Rule V2 resources: %w", err)
	}

	// Extract resource_group_id from the first resource if available
	var resourceGroupId string
	if len(alarm.Resources) > 0 {
		resourceGroupId = alarm.Resources[0].ResourceGroupId
	}

	mErr := multierror.Append(nil,
		d.Set("name", alarm.Name),
		d.Set("description", alarm.Description),
		d.Set("namespace", alarm.Namespace),
		d.Set("resource_group_id", resourceGroupId),
		d.Set("type", alarm.Type),
		d.Set("alarm_enabled", alarm.Enabled),
		d.Set("notification_enabled", alarm.NotificationEnabled),
		d.Set("notification_begin_time", alarm.NotificationBeginTime),
		d.Set("notification_end_time", alarm.NotificationEndTime),
		d.Set("alarm_id", alarm.AlarmId),
		d.Set("alarm_template_id", alarm.AlarmTemplateId),
		d.Set("enterprise_project_id", alarm.EnterpriseProjectId),
		d.Set("policies", flattenPoliciesV2(policiesResp.Policies)),
		d.Set("resources", flattenResourcesV2(resourcesResp.Resources)),
		d.Set("alarm_actions", flattenNotifications(alarm.AlarmNotifications)),
		d.Set("ok_actions", flattenNotifications(alarm.OkNotifications)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenPoliciesV2(policiesResp []policies.PolicyResp) []map[string]interface{} {
	if len(policiesResp) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(policiesResp))
	for i, policy := range policiesResp {
		result[i] = map[string]interface{}{
			"metric_name":         policy.MetricName,
			"period":              policy.Period,
			"filter":              policy.Filter,
			"comparison_operator": policy.ComparisonOperator,
			"value":               policy.Value,
			"unit":                policy.Unit,
			"count":               policy.Count,
			"suppress_duration":   policy.SuppressDuration,
			"level":               policy.Level,
		}
	}
	return result
}

func flattenResourcesV2(resourcesResp [][]resources.Dimension) []map[string]interface{} {
	if len(resourcesResp) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(resourcesResp))
	for i, dims := range resourcesResp {
		dimensions := make([]map[string]interface{}, len(dims))
		for j, dim := range dims {
			dimensions[j] = map[string]interface{}{
				"name":  dim.Name,
				"value": dim.Value,
			}
		}
		result[i] = map[string]interface{}{
			"dimensions": dimensions,
		}
	}
	return result
}

func flattenNotifications(notifications []alarms.Notification) []map[string]interface{} {
	if len(notifications) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(notifications))
	for i, notification := range notifications {
		notifyList := make([]string, len(notification.NotificationList))
		copy(notifyList, notification.NotificationList)
		result[i] = map[string]interface{}{
			"type":              notification.Type,
			"notification_list": notifyList,
		}
	}
	return result
}

func resourceAlarmRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	alarmId := d.Id()

	// Update policies if changed (only when not using alarm_template_id)
	// Template-associated alarm's policies cannot be modified via API
	if d.HasChange("policies") && d.Get("alarm_template_id").(string) == "" {
		updatePolicies := buildUpdatePoliciesOpts(d)
		_, err := policies.Update(client, alarmId, policies.UpdateOpts{
			Policies: updatePolicies,
		})
		if err != nil {
			return fmterr.Errorf("error updating CES Alarm Rule V2 policies: %w", err)
		}
	}

	// Update resources if changed
	if d.HasChange("resources") {
		oldRaw, newRaw := d.GetChange("resources")
		oldResources := oldRaw.(*schema.Set).List()
		newResources := newRaw.(*schema.Set).List()

		// Delete old resources
		if len(oldResources) > 0 {
			deleteOpts := buildResourcesDimensionDelete(oldResources)
			if len(deleteOpts) > 0 {
				err := resources.Delete(client, alarmId, resources.DeleteOpts{
					Resources: deleteOpts,
				})
				if err != nil {
					return fmterr.Errorf("error deleting old resources from CES Alarm Rule V2: %w", err)
				}
			}
		}

		// Add new resources
		if len(newResources) > 0 {
			addOpts := buildResourcesDimensionAdd(newResources)
			if len(addOpts) > 0 {
				err := resources.Add(client, alarmId, resources.AddOpts{
					Resources: addOpts,
				})
				if err != nil {
					return fmterr.Errorf("error adding new resources to CES Alarm Rule V2: %w", err)
				}
			}
		}
	}

	// Update alarm enabled state if changed
	if d.HasChange("alarm_enabled") {
		enabled := d.Get("alarm_enabled").(bool)
		_, err := alarms.Action(client, alarms.ActionOpts{
			AlarmIds:     []string{alarmId},
			AlarmEnabled: enabled,
		})
		if err != nil {
			return fmterr.Errorf("error updating CES Alarm Rule V2 enabled state: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceAlarmRuleV2Read(clientCtx, d, meta)
}

func buildUpdatePoliciesOpts(d *schema.ResourceData) []policies.Policy {
	rawPolicies := d.Get("policies").(*schema.Set).List()
	if len(rawPolicies) == 0 {
		return nil
	}

	policyOpts := make([]policies.Policy, len(rawPolicies))
	for i, v := range rawPolicies {
		policy := v.(map[string]interface{})
		policyOpts[i] = policies.Policy{
			MetricName:         policy["metric_name"].(string),
			Period:             policy["period"].(int),
			Filter:             policy["filter"].(string),
			ComparisonOperator: policy["comparison_operator"].(string),
			Value:              policy["value"].(float64),
			Unit:               policy["unit"].(string),
			Count:              policy["count"].(int),
			SuppressDuration:   policy["suppress_duration"].(int),
			Level:              policy["level"].(int),
		}
	}
	return policyOpts
}

func buildResourcesDimensionDelete(resourcesList []interface{}) [][]resources.Dimension {
	if len(resourcesList) == 0 {
		return nil
	}

	result := make([][]resources.Dimension, len(resourcesList))
	for i, v := range resourcesList {
		resource := v.(map[string]interface{})
		dimensionsRaw := resource["dimensions"].([]interface{})
		dimensions := make([]resources.Dimension, len(dimensionsRaw))
		for j, dim := range dimensionsRaw {
			dimension := dim.(map[string]interface{})
			dimensions[j] = resources.Dimension{
				Name:  dimension["name"].(string),
				Value: dimension["value"].(string),
			}
		}
		result[i] = dimensions
	}
	return result
}

func buildResourcesDimensionAdd(resourcesList []interface{}) [][]resources.Dimension {
	if len(resourcesList) == 0 {
		return nil
	}

	result := make([][]resources.Dimension, len(resourcesList))
	for i, v := range resourcesList {
		resource := v.(map[string]interface{})
		dimensionsRaw := resource["dimensions"].([]interface{})
		dimensions := make([]resources.Dimension, len(dimensionsRaw))
		for j, dim := range dimensionsRaw {
			dimension := dim.(map[string]interface{})
			dimensions[j] = resources.Dimension{
				Name:  dimension["name"].(string),
				Value: dimension["value"].(string),
			}
		}
		result[i] = dimensions
	}
	return result
}

func resourceAlarmRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	alarmId := d.Id()
	log.Printf("[DEBUG] Deleting CES Alarm Rule V2: %s", alarmId)

	_, err = alarms.Delete(client, alarms.DeleteOpts{
		AlarmIds: []string{alarmId},
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES Alarm Rule V2")
	}

	return nil
}
