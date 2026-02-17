package ces

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/oneclickalarms"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceOneClickAlarmV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOneClickAlarmV2Create,
		ReadContext:   resourceOneClickAlarmV2Read,
		UpdateContext: resourceOneClickAlarmV2Update,
		DeleteContext: resourceOneClickAlarmV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"one_click_alarm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimension_names": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Type:         schema.TypeList,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: []string{"dimension_names.0.metric", "dimension_names.0.event"},
						},
						"event": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"notification_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"alarm_notifications": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 20,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ok_notifications": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 20,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"notification_begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"internal_alarm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildOneClickAlarmNotifications(d *schema.ResourceData, key string) []oneclickalarms.Notification {
	rawNotifications := d.Get(key).([]interface{})
	if len(rawNotifications) == 0 {
		return nil
	}

	notifications := make([]oneclickalarms.Notification, len(rawNotifications))
	for i, v := range rawNotifications {
		notification := v.(map[string]interface{})
		notifyListRaw := notification["notification_list"].([]interface{})
		notifyList := make([]string, len(notifyListRaw))
		for j, n := range notifyListRaw {
			notifyList[j] = n.(string)
		}
		notifications[i] = oneclickalarms.Notification{
			Type:             notification["type"].(string),
			NotificationList: notifyList,
		}
	}
	return notifications
}

func buildDimensionNames(d *schema.ResourceData) oneclickalarms.DimensionNames {
	raw := d.Get("dimension_names").([]interface{})
	if len(raw) == 0 {
		return oneclickalarms.DimensionNames{}
	}

	dimMap := raw[0].(map[string]interface{})

	var metric []string
	if v, ok := dimMap["metric"].([]interface{}); ok {
		metric = make([]string, len(v))
		for i, s := range v {
			metric[i] = s.(string)
		}
	}

	var event []string
	if v, ok := dimMap["event"].(bool); ok && v {
		event = []string{""}
	}

	return oneclickalarms.DimensionNames{
		Metric: metric,
		Event:  event,
	}
}

func resourceOneClickAlarmV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	createOpts := oneclickalarms.CreateOpts{
		OneClickAlarmId:       d.Get("one_click_alarm_id").(string),
		DimensionNames:        buildDimensionNames(d),
		NotificationEnabled:   d.Get("notification_enabled").(bool),
		AlarmNotifications:    buildOneClickAlarmNotifications(d, "alarm_notifications"),
		OkNotifications:       buildOneClickAlarmNotifications(d, "ok_notifications"),
		NotificationBeginTime: d.Get("notification_begin_time").(string),
		NotificationEndTime:   d.Get("notification_end_time").(string),
	}

	log.Printf("[DEBUG] CES One-Click Alarm V2 Create Options: %#v", createOpts)

	alarmId, err := oneclickalarms.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CES One-Click Alarm V2: %w", err)
	}

	log.Printf("[DEBUG] CES One-Click Alarm V2 created with ID: %s", alarmId)
	d.SetId(alarmId)

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceOneClickAlarmV2Read(clientCtx, d, meta)
}

func resourceOneClickAlarmV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	alarmList, err := oneclickalarms.List(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CES One-Click Alarm V2")
	}

	var alarm *oneclickalarms.OneClickAlarm
	for _, a := range alarmList {
		if a.OneClickAlarmId == d.Id() {
			alarm = &a
			break
		}
	}

	if alarm == nil {
		log.Printf("[WARN] CES One-Click Alarm V2 (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(nil,
		d.Set("internal_alarm_id", alarm.OneClickAlarmId),
		d.Set("namespace", alarm.Namespace),
		d.Set("description", alarm.Description),
		d.Set("enabled", alarm.Enabled),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOneClickAlarmV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	if d.HasChanges("notification_enabled", "alarm_notifications", "ok_notifications",
		"notification_begin_time", "notification_end_time") {
		updateOpts := oneclickalarms.UpdateNotificationsOpts{
			NotificationEnabled:   d.Get("notification_enabled").(bool),
			AlarmNotifications:    buildOneClickAlarmNotifications(d, "alarm_notifications"),
			OkNotifications:       buildOneClickAlarmNotifications(d, "ok_notifications"),
			NotificationBeginTime: d.Get("notification_begin_time").(string),
			NotificationEndTime:   d.Get("notification_end_time").(string),
		}

		log.Printf("[DEBUG] CES One-Click Alarm V2 Update Notifications: %#v", updateOpts)

		err := oneclickalarms.UpdateNotifications(client, d.Id(), updateOpts)
		if err != nil {
			return fmterr.Errorf("error updating CES One-Click Alarm V2 notifications: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceOneClickAlarmV2Read(clientCtx, d, meta)
}

func resourceOneClickAlarmV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	alarmId := d.Id()
	log.Printf("[DEBUG] Deleting CES One-Click Alarm V2: %s", alarmId)

	_, err = oneclickalarms.BatchDelete(client, oneclickalarms.BatchDeleteOpts{
		OneClickAlarmIds: []string{alarmId},
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES One-Click Alarm V2")
	}

	return nil
}
