package dms

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/topics"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsTopicsV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsTopicsV1Create,
		ReadContext:   resourceDmsTopicsV1Read,
		DeleteContext: resourceDmsTopicsV1Delete,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateName,
			},
			"partition": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"replication": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},
			"sync_replication": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 720),
			},
			"sync_message_flush": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"remain_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDmsTopicsV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := &topics.CreateOpts{
		Name:             d.Get("name").(string),
		Partition:        d.Get("partition").(int),
		Replication:      d.Get("replication").(int),
		SyncReplication:  d.Get("sync_replication").(bool),
		RetentionTime:    d.Get("retention_time").(int),
		SyncMessageFlush: d.Get("sync_message_flush").(bool),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := topics.Create(DmsV1Client, createOpts, d.Get("instance_id").(string)).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud queue: %s", err)
	}
	// Store the topic Name/ID now
	d.SetId(v.Name)

	return resourceDmsTopicsV1Read(ctx, d, meta)
}

func resourceDmsTopicsV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}
	instanceId := d.Get("instance_id").(string)

	v, err := topics.Get(DmsV1Client, instanceId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS topic")
	}

	var fTopic topics.Parameters
	found := false

	for _, topic := range v.Topics {
		if topic.Name == d.Id() {
			fTopic = topic
			found = true
			break
		}
	}
	if !found {
		return fmterr.Errorf("Provided topic doesn't exist")
	}

	// conversion is done because API values are returned as strings
	syncReplication, _ := strconv.ParseBool(fTopic.SyncReplication)
	syncMessageFlush, _ := strconv.ParseBool(fTopic.SyncMessageFlush)

	mErr := multierror.Append(
		d.Set("name", fTopic.Name),
		d.Set("partition", fTopic.Partition),
		d.Set("replication", fTopic.Replication),
		d.Set("sync_replication", syncReplication),
		d.Set("retention_time", fTopic.RetentionTime),
		d.Set("sync_message_flush", syncMessageFlush),
		d.Set("size", v.Size),
		d.Set("remain_partitions", v.RemainPartitions),
		d.Set("max_partitions", v.MaxPartitions),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDmsTopicsV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}
	deleteOpts := topics.DeleteOpts{
		Topics: []string{
			d.Id(),
		},
	}
	_, err = topics.Delete(DmsV1Client, deleteOpts, d.Get("instance_id").(string)).Extract()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud topic: %s", err)
	}

	log.Printf("[DEBUG] Dms topic %s deactivated.", d.Id())
	d.SetId("")
	return nil
}
