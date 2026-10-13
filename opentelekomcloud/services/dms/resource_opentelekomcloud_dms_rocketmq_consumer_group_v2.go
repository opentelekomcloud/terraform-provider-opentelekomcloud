package dms

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQConsumerGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQConsumerGroupV2Create,
		ReadContext:   resourceDmsRocketMQConsumerGroupV2Read,
		UpdateContext: resourceDmsRocketMQConsumerGroupV2Update,
		DeleteContext: resourceDmsRocketMQConsumerGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"broadcast": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"retry_max_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16),
			},
			"consume_orderly": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
			},
		},
	}
}

func resourceDmsRocketMQConsumerGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	createOpts := groups.CreateOpts{
		Name:           d.Get("name").(string),
		Enabled:        pointerto.Bool(d.Get("enabled").(bool)),
		Broadcast:      pointerto.Bool(d.Get("broadcast").(bool)),
		RetryMaxTime:   d.Get("retry_max_time").(int),
		ConsumeOrderly: pointerto.Bool(d.Get("consume_orderly").(bool)),
		Description:    d.Get("description").(string),
	}
	log.Printf("[DEBUG] Create DMS RocketMQ consumer group options: %#v", createOpts)

	created, err := groups.Create(client, instanceID, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS RocketMQ consumer group: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceID, created.Name))

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQConsumerGroupV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQConsumerGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	v, err := groups.Get(client, instanceID, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS RocketMQ consumer group")
	}
	log.Printf("[DEBUG] Get DMS RocketMQ consumer group: %+v", v)

	mErr := multierror.Append(
		d.Set("instance_id", instanceID),
		d.Set("name", v.Name),
		d.Set("enabled", v.Enabled),
		d.Set("broadcast", v.Broadcast),
		d.Set("retry_max_time", v.RetryMaxTime),
		d.Set("consume_orderly", v.ConsumeOrderly),
		d.Set("description", v.Description),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQConsumerGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Single group modification API doesn't accept consume_orderly and resets it to false,
	// batch modification API accepts all the fields.
	want := groups.ConsumerGroup{
		Enabled:        d.Get("enabled").(bool),
		Broadcast:      d.Get("broadcast").(bool),
		RetryMaxTime:   d.Get("retry_max_time").(int),
		ConsumeOrderly: d.Get("consume_orderly").(bool),
		Description:    d.Get("description").(string),
	}
	_, err = groups.BatchUpdate(client, instanceID, groups.BatchUpdateOpts{
		Groups: []groups.BatchUpdateGroup{
			{
				Name:           name,
				Enabled:        &want.Enabled,
				Broadcast:      &want.Broadcast,
				RetryMaxTime:   want.RetryMaxTime,
				ConsumeOrderly: &want.ConsumeOrderly,
				Description:    &want.Description,
			},
		},
	})
	if err != nil {
		return diag.Errorf("error updating DMS RocketMQ consumer group (%s): %s", d.Id(), err)
	}

	// Batch modification is asynchronous.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"DONE"},
		Refresh: func() (interface{}, string, error) {
			v, err := groups.Get(client, instanceID, name)
			if err != nil {
				return nil, "", err
			}
			if v.Enabled == want.Enabled && v.Broadcast == want.Broadcast && v.RetryMaxTime == want.RetryMaxTime &&
				v.ConsumeOrderly == want.ConsumeOrderly && v.Description == want.Description {
				return v, "DONE", nil
			}
			return v, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 5 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for DMS RocketMQ consumer group (%s) to be updated: %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQConsumerGroupV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQConsumerGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err = groups.Delete(client, instanceID, name); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DMS RocketMQ consumer group")
	}
	return nil
}
