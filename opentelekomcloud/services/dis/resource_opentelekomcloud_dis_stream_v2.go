package dis

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
	_ "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/streams"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDisStreamV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDisStreamV2Create,
		ReadContext:   resourceDisStreamV2Read,
		DeleteContext: resourceDisStreamV2Delete,
		UpdateContext: resourceDisStreamV2Update,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partition_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(24, 72),
			},
			"stream_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"COMMON", "ADVANCED",
				}, false),
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BLOB",
				}, false),
			},
			"compression_format": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"snappy", "gzip", "zip",
				}, false),
			},
			"auto_scale_min_partition_count": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"auto_scale_max_partition_count"},
			},
			"auto_scale_max_partition_count": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"auto_scale_min_partition_count"},
			},
			"tags": common.TagsSchema(),
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"readable_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"writable_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hash_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sequence_number_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_partitions": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDisStreamV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := streams.CreateStreamOpts{
		StreamName:        d.Get("stream_name").(string),
		PartitionCount:    d.Get("partition_count").(int),
		StreamType:        d.Get("stream_type").(string),
		DataDuration:      pointerto.Int(d.Get("retention_period").(int)),
		DataType:          d.Get("data_type").(string),
		CompressionFormat: d.Get("compression_format").(string),
		Tags:              common.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	opts.AutoScaleEnabled = pointerto.Bool(false)
	autoScaleMinPartitionCount := d.Get("auto_scale_min_partition_count").(int)
	autoScaleMaxPartitionCount := d.Get("auto_scale_max_partition_count").(int)
	if autoScaleMinPartitionCount > 0 && autoScaleMaxPartitionCount > 0 {
		opts.AutoScaleEnabled = pointerto.Bool(true)
		opts.AutoScaleMinPartitionCount = &autoScaleMinPartitionCount
		opts.AutoScaleMaxPartitionCount = &autoScaleMaxPartitionCount
	}

	log.Printf("[DEBUG] Creating new Stream: %#v", opts)
	err = streams.CreateStream(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating DIS streams: %s", err)
	}

	d.SetId(opts.StreamName)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDisStreamV2Read(clientCtx, d, meta)
}

func resourceDisStreamV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	stream, err := streams.GetStream(client, streams.GetStreamOpts{
		StreamName: d.Id(),
	})
	if err != nil {
		return fmterr.Errorf("error getting stream details: %s", err)
	}

	mErr := multierror.Append(
		d.Set("stream_name", stream.StreamName),
		d.Set("auto_scale_max_partition_count", stream.AutoScaleMaxPartitionCount),
		d.Set("auto_scale_min_partition_count", stream.AutoScaleMinPartitionCount),
		d.Set("compression_format", stream.CompressionFormat),
		d.Set("data_type", stream.DataType),
		d.Set("retention_period", stream.RetentionPeriod),
		d.Set("stream_type", stream.StreamType),
		d.Set("tags", common.TagsToMap(stream.Tags)),
		d.Set("created", stream.CreatedAt),
		d.Set("readable_partition_count", stream.ReadablePartitionCount),
		d.Set("writable_partition_count", stream.WritablePartitionCount),
		d.Set("partition_count", stream.WritablePartitionCount),
		d.Set("status", stream.Status),
		d.Set("stream_id", stream.StreamId),
		queryAndSetPartitionsToState(client, d, stream.StreamName),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDisStreamV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	name := d.Id()
	err = streams.DeleteStream(client, name)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DIS stream: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDisStreamV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	name := d.Id()

	// Update partition count
	if d.HasChange("partition_count") {
		newValue := d.Get("partition_count").(int)
		updateOpts := streams.UpdatePartitionCountOpts{
			StreamName:           name,
			TargetPartitionCount: newValue,
		}
		partErr := streams.UpdatePartitionCount(client, updateOpts)
		if partErr != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud DIS stream: %s", err)
		}

		checkErr := checkPartitionUpdateResult(ctx, client, name, newValue, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud DIS stream: %s", checkErr)
		}
	}

	if d.HasChange("tags") {
		streamId := d.Get("stream_id").(string)
		tagErr := common.UpdateResourceTags(client, d, "stream", streamId)
		if tagErr != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud DIS stream tags: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDisStreamV2Read(clientCtx, d, meta)
}

func checkPartitionUpdateResult(ctx context.Context, client *golangsdk.ServiceClient, name string, targetValue int,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			stream, err := streams.GetStream(client, streams.GetStreamOpts{
				StreamName: name,
			})
			if err != nil {
				return nil, "failed", err
			}
			log.Printf("[DEBUG] OpenTelekomCloud DIS stream WritablePartitionCount=%v, target=%v", stream.WritablePartitionCount, targetValue)
			if *stream.WritablePartitionCount == targetValue {
				return stream, "Done", nil
			}
			return stream, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 5 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to waiting for OpenTelekomCloud DIS stream (%s) partition update: %s", name, err)
	}
	return nil
}

func queryAndSetPartitionsToState(client *golangsdk.ServiceClient, d *schema.ResourceData, streamName string) error {
	var result []map[string]interface{}
	opts := streams.GetStreamOpts{
		StreamName: streamName,
	}
	for {
		stream, err := streams.GetStream(client, opts)
		if err != nil {
			return fmt.Errorf("error query OpenTelekomCloud DIS stream partitions: %s", err)
		}

		for _, partition := range stream.Partitions {
			result = append(result, map[string]interface{}{
				"id":                    partition.PartitionId,
				"status":                partition.Status,
				"hash_range":            partition.HashRange,
				"sequence_number_range": partition.SequenceNumberRange,
				"parent_partitions":     partition.ParentPartitions,
			})
		}

		if !*stream.HasMorePartitions {
			break
		}

		opts.StartPartitionId = stream.Partitions[len(stream.Partitions)-1].PartitionId
	}

	return d.Set("partitions", result)
}
