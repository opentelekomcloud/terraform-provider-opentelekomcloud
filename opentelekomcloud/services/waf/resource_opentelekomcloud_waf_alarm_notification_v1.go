package waf

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf/v1/alarmnotifications"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceWafAlarmNotificationV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafAlarmNotificationV1Create,
		ReadContext:   resourceWafAlarmNotificationV1Read,
		UpdateContext: resourceWafAlarmNotificationV1Update,
		DeleteContext: resourceWafAlarmNotificationV1Delete,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"send_frequency": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					5, 15, 30, 60,
				}),
			},
			"times": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"threat": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"all", "cc", "cmdi", "custom", "illegal", "sqli", "lfi", "robot",
						"antitamper", "rfi", "vuln", "xss", "whiteblackip", "webshell",
					}, false),
				},
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceWAFThreat(d *schema.ResourceData) []string {
	threatRaw := d.Get("threat").(*schema.Set).List()
	var threat []string
	for _, v := range threatRaw {
		threat = append(threat, v.(string))
	}
	return threat
}

func resourceWafAlarmNotificationV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(WafClientError, err)
	}
	alarmNotification, err := alarmnotifications.List(client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	enabled := d.Get("enabled").(bool)
	topicURN := d.Get("topic_urn").(string)
	updateOpts := alarmnotifications.UpdateOpts{
		Enabled:       &enabled,
		TopicURN:      &topicURN,
		SendFrequency: d.Get("send_frequency").(int),
		Times:         d.Get("times").(int),
		Threat:        resourceWAFThreat(d),
		Locale:        d.Get("locale").(string),
	}

	alarmNotification, err = alarmnotifications.Update(client, alarmNotification.ID, updateOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alarmNotification.ID)

	return resourceWafAlarmNotificationV1Read(ctx, d, meta)
}

func resourceWafAlarmNotificationV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(WafClientError, err)
	}

	alarmNotification, err := alarmnotifications.List(client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	var threat []string
	for _, v := range alarmNotification.Threat {
		threat = append(threat, v)
	}

	mErr := multierror.Append(
		d.Set("enabled", alarmNotification.Enabled),
		d.Set("topic_urn", alarmNotification.TopicURN),
		d.Set("send_frequency", alarmNotification.SendFrequency),
		d.Set("times", alarmNotification.Times),
		d.Set("threat", threat),
		d.Set("locale", alarmNotification.Locale),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWafAlarmNotificationV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(WafClientError, err)
	}

	enabled := d.Get("enabled").(bool)
	topicURN := d.Get("topic_urn").(string)
	updateOpts := alarmnotifications.UpdateOpts{
		Enabled:       &enabled,
		TopicURN:      &topicURN,
		SendFrequency: d.Get("send_frequency").(int),
		Times:         d.Get("times").(int),
		Threat:        resourceWAFThreat(d),
		Locale:        d.Get("locale").(string),
	}

	_, err = alarmnotifications.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWafAlarmNotificationV1Read(ctx, d, meta)
}

func resourceWafAlarmNotificationV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(WafClientError, err)
	}

	disabled := false
	emptyTopicURN := ""
	updateOpts := alarmnotifications.UpdateOpts{
		Enabled:       &disabled,
		TopicURN:      &emptyTopicURN,
		SendFrequency: 5,
		Times:         1,
		Threat:        []string{"all"},
	}

	_, err = alarmnotifications.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
