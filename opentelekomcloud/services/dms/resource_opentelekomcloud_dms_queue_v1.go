package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/queues"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsQueuesV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsQueuesV1Create,
		ReadContext:   resourceDmsQueuesV1Read,
		DeleteContext: resourceDmsQueuesV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		DeprecationMessage: "Support will be discontinued in favor of DMS Kafka Premium. " +
			"Please use `opentelekomcloud_dms_instance_v1` resource instead",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"redrive_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"max_consume_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"created": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"reservation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_msg_size_byte": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsQueuesV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms queue client: %s", err)
	}

	createOpts := &queues.CreateOpts{
		Name:            d.Get("name").(string),
		QueueMode:       d.Get("queue_mode").(string),
		Description:     d.Get("description").(string),
		RedrivePolicy:   d.Get("redrive_policy").(string),
		MaxConsumeCount: d.Get("max_consume_count").(int),
		RetentionHours:  d.Get("retention_hours").(int),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := queues.Create(DmsV1Client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud queue: %s", err)
	}
	log.Printf("[INFO] Queue ID: %s", v.ID)

	// Store the queue ID now
	d.SetId(v.ID)

	return resourceDmsQueuesV1Read(ctx, d, meta)
}

func resourceDmsQueuesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms queue client: %s", err)
	}
	v, err := queues.Get(DmsV1Client, d.Id(), true).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Dms queue %s: %+v", d.Id(), v)

	d.SetId(v.ID)
	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("created", v.Created),
		d.Set("description", v.Description),
		d.Set("queue_mode", v.QueueMode),
		d.Set("reservation", v.Reservation),
		d.Set("max_msg_size_byte", v.MaxMsgSizeByte),
		d.Set("produced_messages", v.ProducedMessages),
		d.Set("redrive_policy", v.RedrivePolicy),
		d.Set("max_consume_count", v.MaxConsumeCount),
		d.Set("group_count", v.GroupCount),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDmsQueuesV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms queue client: %s", err)
	}

	v, err := queues.Get(DmsV1Client, d.Id(), false).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "queue")
	}

	err = queues.Delete(DmsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud queue: %s", err)
	}

	log.Printf("[DEBUG] Dms queue %s: %+v deactivated.", d.Id(), v)
	d.SetId("")
	return nil
}
