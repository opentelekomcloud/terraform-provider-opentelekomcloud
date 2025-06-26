package ces

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/events"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCesEventDetailsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesEventDetailsRead,

		Schema: map[string]*schema.Schema{
			"event_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EVENT.SYS", "EVENT.CUSTOM",
				}, false),
			},
			"event_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_level": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Critical", "Major", "Minor", "Info",
				}, false),
			},
			"event_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"normal", "warning", "incident",
				}, false),
			},
			"from": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"event_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"event_sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"event_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"detail": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"event_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"event_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"event_user": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"event_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"event_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"meta_data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCesEventDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	listOpts := events.ListEventDetailOpts{
		EventName:   d.Get("event_name").(string),
		EventType:   d.Get("event_type").(string),
		EventSource: d.Get("event_source").(string),
		EventLevel:  d.Get("event_level").(string),
		EventUser:   d.Get("event_user").(string),
		EventState:  d.Get("event_state").(string),
		From:        InterfaceToInt64(d.Get("from")),
		To:          InterfaceToInt64(d.Get("to")),
		Limit:       d.Get("limit").(int),
	}

	listEventDetails, err := events.ListEventDetail(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching event details: %w", err)
	}

	d.SetId(listEventDetails.EventName)

	log.Printf("[DEBUG] Retrieved events list %s: %#v", d.Id(), listEventDetails)

	metadata := []map[string]interface{}{
		{
			"total": listEventDetails.MetaData.Total,
		},
	}

	mErr := multierror.Append(
		d.Set("event_name", listEventDetails.EventName),
		d.Set("event_type", listEventDetails.EventType),
		d.Set("event_users", listEventDetails.EventUsers),
		d.Set("event_sources", listEventDetails.EventSources),
		d.Set("event_info", setEventInfo(listEventDetails.EventInfo)),
		d.Set("meta_data", metadata),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setEventInfo(eventInfoListInResp []events.EventInfoDetail) []map[string]interface{} {
	var eventInfoList []map[string]interface{}
	for _, eventInfoInResp := range eventInfoListInResp {
		eventInfo := map[string]interface{}{
			"event_name":   eventInfoInResp.EventName,
			"event_source": eventInfoInResp.EventSource,
			"time":         eventInfoInResp.Time,
			"detail":       setEventDetails(eventInfoInResp.Detail),
			"event_id":     eventInfoInResp.EventId,
		}
		eventInfoList = append(eventInfoList, eventInfo)
	}
	return eventInfoList
}

func setEventDetails(eventItemDetailInResp events.EventItemDetail) []map[string]interface{} {
	var eventItemDetailList []map[string]interface{}
	eventItemDetail := map[string]interface{}{
		"content":       eventItemDetailInResp.Content,
		"group_id":      eventItemDetailInResp.GroupId,
		"resource_id":   eventItemDetailInResp.ResourceId,
		"resource_name": eventItemDetailInResp.ResourceName,
		"event_state":   eventItemDetailInResp.EventState,
		"event_level":   eventItemDetailInResp.EventLevel,
		"event_user":    eventItemDetailInResp.EventUser,
		"event_type":    eventItemDetailInResp.EventType,
	}
	eventItemDetailList = append(eventItemDetailList, eventItemDetail)
	return eventItemDetailList
}
