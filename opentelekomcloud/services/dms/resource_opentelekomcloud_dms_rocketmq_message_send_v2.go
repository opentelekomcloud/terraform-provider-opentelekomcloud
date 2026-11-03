package dms

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/messages"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQMessageSendV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQMessageSendV2Create,
		ReadContext:   resourceDmsRocketMQMessageSendV2Read,
		DeleteContext: resourceDmsRocketMQMessageSendV2Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"body": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"property_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"msg_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"queue_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"queue_offset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"broker_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsRocketMQMessageSendV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	var properties []messages.Property
	for _, raw := range d.Get("property_list").([]interface{}) {
		p := raw.(map[string]interface{})
		properties = append(properties, messages.Property{
			Name:  p["name"].(string),
			Value: p["value"].(string),
		})
	}

	res, err := messages.Send(client, d.Get("instance_id").(string), messages.SendOpts{
		Topic:        d.Get("topic").(string),
		Body:         d.Get("body").(string),
		PropertyList: properties,
	})
	if err != nil {
		return diag.Errorf("error sending DMS RocketMQ message: %s", err)
	}
	d.SetId(res.MsgID)

	// The API returns the message metadata only in the response.
	mErr := multierror.Append(
		d.Set("msg_id", res.MsgID),
		d.Set("queue_id", res.QueueID),
		d.Set("queue_offset", res.QueueOffset),
		d.Set("broker_name", res.BrokerName),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQMessageSendV2Read(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsRocketMQMessageSendV2Delete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The resource is only removed from the state, the message remains in the topic.",
		},
	}
}
