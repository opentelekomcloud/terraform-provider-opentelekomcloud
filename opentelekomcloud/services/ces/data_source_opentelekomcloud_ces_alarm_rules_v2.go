package ces

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/alarms"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/policies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/resources"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceCesAlarmRulesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesAlarmRulesV2Read,

		Schema: map[string]*schema.Schema{
			"alarm_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"notification_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"notification_begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"period": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"filter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"comparison_operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"suppress_duration": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"level": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"alarm_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"notification_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"ok_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"notification_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCesAlarmRulesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	listOpts := alarms.ListOpts{
		AlarmId:    d.Get("alarm_id").(string),
		Name:       d.Get("name").(string),
		Namespace:  d.Get("namespace").(string),
		ResourceId: d.Get("resource_id").(string),
	}

	listResp, err := alarms.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error listing CES alarm rules v2: %w", err)
	}

	ids := make([]string, 0, len(listResp.Alarms))
	alarmsList := make([]map[string]interface{}, 0, len(listResp.Alarms))

	for _, alarm := range listResp.Alarms {
		ids = append(ids, alarm.AlarmId)

		policiesResp, err := policies.List(client, alarm.AlarmId, policies.ListOpts{})
		if err != nil {
			log.Printf("[WARN] Error retrieving policies for alarm %s: %s", alarm.AlarmId, err)
			continue
		}

		resourcesResp, err := resources.List(client, alarm.AlarmId, resources.ListOpts{})
		if err != nil {
			log.Printf("[WARN] Error retrieving resources for alarm %s: %s", alarm.AlarmId, err)
			continue
		}

		alarmMap := map[string]interface{}{
			"alarm_id":                alarm.AlarmId,
			"name":                    alarm.Name,
			"description":             alarm.Description,
			"namespace":               alarm.Namespace,
			"type":                    alarm.Type,
			"alarm_enabled":           alarm.Enabled,
			"notification_enabled":    alarm.NotificationEnabled,
			"notification_begin_time": alarm.NotificationBeginTime,
			"notification_end_time":   alarm.NotificationEndTime,
			"enterprise_project_id":   alarm.EnterpriseProjectId,
			"alarm_template_id":       alarm.AlarmTemplateId,
			"policies":                flattenPoliciesV2(policiesResp.Policies),
			"resources":               flattenResourcesV2(resourcesResp.Resources),
			"alarm_actions":           flattenNotifications(alarm.AlarmNotifications),
			"ok_actions":              flattenNotifications(alarm.OkNotifications),
		}
		alarmsList = append(alarmsList, alarmMap)
	}

	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("alarms", alarmsList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
