package ces

import (
	"context"
	"log"
	"regexp"
	"time"

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

func ResourceCesEventReportV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCesEventReportCreate,
		ReadContext:   resourceCesEventReportRead,
		DeleteContext: resourceCesEventReportDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: checkCesAlarmRestrictions,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"event_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`),
						"Only lowercase/uppercase letters, digits, and underscores (_) are allowed and must start with a letter.",
					),
				),
			},
			"event_source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 32),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*\.[a-zA-Z][a-zA-Z0-9_]*$`),
						"Must be of type service.item. service and item each must only have lowercase/uppercase letters, digits, and underscores (_) and must start with a letter.",
					),
				),
			},
			"time": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"detail": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"event_state": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"normal", "warning", "incident",
							}, false),
						},
						"event_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Critical", "Major", "Minor", "Info",
							}, false),
						},
						"event_user": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"event_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"EVENT.SYS", "EVENT.CUSTOM",
							}, false),
						},
					},
				},
			},
			"event_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getEventDetails(d *schema.ResourceData) events.EventItemDetail {
	eventDetailsRaw := d.Get("detail").([]interface{})
	eventDetail := eventDetailsRaw[0].(map[string]interface{})

	return events.EventItemDetail{
		Content:      eventDetail["content"].(string),
		GroupId:      eventDetail["group_id"].(string),
		ResourceId:   eventDetail["resource_id"].(string),
		ResourceName: eventDetail["resource_name"].(string),
		EventState:   eventDetail["event_state"].(string),
		EventLevel:   eventDetail["event_level"].(string),
		EventUser:    eventDetail["event_user"].(string),
		EventType:    eventDetail["event_type"].(string),
	}
}

func resourceCesEventReportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := events.EventItem{
		EventName:   d.Get("event_name").(string),
		EventSource: d.Get("event_source").(string),
		Time:        InterfaceToInt64(d.Get("time")),
		Detail:      getEventDetails(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	events, err := events.CreateEvents(client, []events.EventItem{
		createOpts,
	})
	if err != nil {
		return fmterr.Errorf("error creating event: %w", err)
	}
	log.Printf("[DEBUG] Created custom event, ID: %#v", events[0].EventId)

	d.SetId(events[0].EventId)

	clientCtx := common.CtxWithClient(ctx, client, cesClientV1)
	return resourceCesEventReportRead(clientCtx, d, meta)
}

func resourceCesEventReportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mErr := multierror.Append(
		d.Set("event_id", d.Id()),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCesEventReportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}
