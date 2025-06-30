package ces

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/events"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCesEventsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesEventsRead,

		Schema: map[string]*schema.Schema{
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EVENT.SYS", "EVENT.CUSTOM",
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
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_occur_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_event_source": {
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

func dataSourceCesEventsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	listOpts := events.ListEventsOpts{
		EventName: d.Get("event_name").(string),
		EventType: d.Get("event_type").(string),
		From:      InterfaceToInt64(d.Get("from")),
		To:        InterfaceToInt64(d.Get("to")),
		Limit:     d.Get("limit").(int),
	}

	listEvents, err := events.ListEvents(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching event : %w", err)
	}

	d.SetId(tools.RandomString("events_", 10))

	log.Printf("[DEBUG] Retrieved events list %s: %#v", d.Id(), listEvents)

	metadata := []map[string]interface{}{
		{
			"total": listEvents.MetaData.Total,
		},
	}

	mErr := multierror.Append(
		d.Set("events", setEvents(listEvents.Events)),
		d.Set("meta_data", metadata),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setEvents(eventInfoListInResp []events.EventInfo) []map[string]interface{} {
	var eventInfoList []map[string]interface{}
	for _, eventInfoInResp := range eventInfoListInResp {
		eventInfo := map[string]interface{}{
			"event_name":          eventInfoInResp.EventName,
			"event_type":          eventInfoInResp.EventType,
			"event_count":         eventInfoInResp.EventCount,
			"latest_occur_time":   eventInfoInResp.LatestOccurTime,
			"latest_event_source": eventInfoInResp.LatestEventSource,
		}
		eventInfoList = append(eventInfoList, eventInfo)
	}
	return eventInfoList
}
