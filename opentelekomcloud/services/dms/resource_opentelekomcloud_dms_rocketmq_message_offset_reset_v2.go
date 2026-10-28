package dms

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQMessageOffsetResetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQMessageOffsetResetV2Create,
		ReadContext:   resourceDmsRocketMQMessageOffsetResetV2Read,
		DeleteContext: resourceDmsRocketMQMessageOffsetResetV2Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// String, as a timestamp in milliseconds doesn't fit into TypeInt on 32-bit platforms.
			"timestamp": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+$`), "must be a Unix timestamp in milliseconds"),
			},
			"queues": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timestamp_offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDmsRocketMQMessageOffsetResetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	timestamp, err := strconv.ParseInt(d.Get("timestamp").(string), 10, 64)
	if err != nil {
		return diag.Errorf("error parsing timestamp: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	group := d.Get("group").(string)
	topic := d.Get("topic").(string)
	res, err := groups.ResetOffset(client, instanceID, group, groups.ResetOffsetOpts{
		Topic:     topic,
		Timestamp: timestamp,
	})
	if err != nil {
		return diag.Errorf("error resetting DMS RocketMQ consumer offset: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, group, topic))

	// The API returns the reset queues only in the response.
	queues := make([]map[string]interface{}, len(res.Queues))
	for i, q := range res.Queues {
		queues[i] = map[string]interface{}{
			"broker_name":      q.BrokerName,
			"queue_id":         q.QueueID,
			"timestamp_offset": q.TimestampOffset,
		}
	}
	return diag.FromErr(d.Set("queues", queues))
}

func resourceDmsRocketMQMessageOffsetResetV2Read(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsRocketMQMessageOffsetResetV2Delete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The resource is only removed from the state, the consumer offset remains reset.",
		},
	}
}
