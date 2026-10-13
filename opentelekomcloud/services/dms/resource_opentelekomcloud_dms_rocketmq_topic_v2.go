package dms

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/topics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQTopicV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQTopicV2Create,
		ReadContext:   resourceDmsRocketMQTopicV2Read,
		UpdateContext: resourceDmsRocketMQTopicV2Update,
		DeleteContext: resourceDmsRocketMQTopicV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// queue_num and description are not returned by the API, so they are kept as configured.
			"queue_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"message_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NORMAL", "FIFO", "DELAY", "TRANSACTION",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_read_queue_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_write_queue_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsRocketMQTopicV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	createOpts := topics.CreateOpts{
		Name:        d.Get("name").(string),
		QueueNum:    d.Get("queue_num").(int),
		MessageType: d.Get("message_type").(string),
	}
	log.Printf("[DEBUG] Create DMS RocketMQ topic options: %#v", createOpts)

	created, err := topics.Create(client, instanceID, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS RocketMQ topic: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceID, created.ID))

	// Description can't be set on creation.
	if desc := d.Get("description").(string); desc != "" {
		if err = topics.Update(client, instanceID, created.ID, topics.UpdateOpts{Description: &desc}); err != nil {
			return diag.Errorf("error setting description of DMS RocketMQ topic (%s): %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQTopicV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQTopicV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, name, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	v, err := topics.Get(client, instanceID, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS RocketMQ topic")
	}
	log.Printf("[DEBUG] Get DMS RocketMQ topic: %+v", v)

	mErr := multierror.Append(
		d.Set("instance_id", instanceID),
		d.Set("name", v.Name),
		d.Set("message_type", v.MessageType),
		d.Set("permission", v.Permission),
		d.Set("total_read_queue_num", v.TotalReadQueueNum),
		d.Set("total_write_queue_num", v.TotalWriteQueueNum),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQTopicV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, name, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	desc := d.Get("description").(string)
	if err = topics.Update(client, instanceID, name, topics.UpdateOpts{Description: &desc}); err != nil {
		return diag.Errorf("error updating DMS RocketMQ topic (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQTopicV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQTopicV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, name, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = topics.Delete(client, instanceID, name); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DMS RocketMQ topic")
	}
	return nil
}
