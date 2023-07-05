package dis

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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	_ "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/dump"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDisDumpV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDisDumpV2Create,
		ReadContext:   resourceDisDumpV2Read,
		DeleteContext: resourceDisDumpV2Delete,
		UpdateContext: resourceDisDumpV2Update,
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"OBS",
				}, false),
			},
			"obs_destination_descriptor": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 64),
								validation.StringMatch(
									regexp.MustCompile(`^[A-Za-z0-9\-_]+$`),
									"Only letters, digits, underscores (_), and hyphens (-) are allowed.",
								),
							),
						},
						"agency_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"deliver_time_interval": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(30, 900),
						},
						"consumer_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"LATEST", "TRIM_HORIZON",
							}, false),
						},
						"file_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 50),
							),
						},
						"partition_format": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"yyyy", "yyyy/MM", "yyyy/MM/dd", "yyyy/MM/dd/HH", "yyyy/MM/dd/HH/mm",
							}, false),
						},
						"obs_bucket_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination_file_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"text",
							}, false),
						},
						"record_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								",", ";", "|", "\\n",
							}, false),
						},
					},
				},
			},
			"obs_processing_schema": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timestamp_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timestamp_format": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"yyyy/MM/dd HH:mm:ss", "MM/dd/yyyy HH:mm:ss", "dd/MM/yyyy HH:mm:ss",
								"yyyy-MM-dd HH:mm:ss", "MM-dd-yyyy HH:mm:ss", "dd-MM-yyyy HH:mm:ss",
							}, false),
						},
					},
				},
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start", "stop",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_transfer_timestamp": {
				Type:     schema.TypeInt,
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

func getOBSDestinationDescriptorOpts(d *schema.ResourceData) dump.OBSDestinationDescriptorOpts {
	var componentOpts dump.OBSDestinationDescriptorOpts

	descriptor := d.Get("obs_destination_descriptor").(*schema.Set).List()
	processing := d.Get("obs_processing_schema").(*schema.Set).List()

	if len(descriptor) > 0 {
		rawDesc := descriptor[0].(map[string]interface{})
		componentOpts = dump.OBSDestinationDescriptorOpts{
			TaskName:            rawDesc["task_name"].(string),
			AgencyName:          rawDesc["agency_name"].(string),
			DeliverTimeInterval: pointerto.Int(rawDesc["deliver_time_interval"].(int)),
			ConsumerStrategy:    rawDesc["consumer_strategy"].(string),
			FilePrefix:          rawDesc["file_prefix"].(string),
			PartitionFormat:     rawDesc["partition_format"].(string),
			OBSBucketPath:       rawDesc["obs_bucket_path"].(string),
			DestinationFileType: rawDesc["destination_file_type"].(string),
			RecordDelimiter:     rawDesc["record_delimiter"].(string),
		}
	}
	if len(processing) > 0 {
		rawProc := processing[0].(map[string]interface{})
		componentOpts.ProcessingSchema = dump.ProcessingSchema{
			TimestampName:   rawProc["timestamp_name"].(string),
			TimestampType:   rawProc["timestamp_type"].(string),
			TimestampFormat: rawProc["timestamp_format"].(string),
		}
	}

	return componentOpts
}

func resourceDisDumpV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := dump.CreateOBSDumpTaskOpts{
		StreamName:               d.Get("stream_name").(string),
		DestinationType:          d.Get("destination").(string),
		OBSDestinationDescriptor: getOBSDestinationDescriptorOpts(d),
	}

	log.Printf("[DEBUG] Creating new dump task: %#v", opts)
	err = dump.CreateOBSDumpTask(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating DIS dump task: %s", err)
	}

	state := d.Get("action").(string)
	if state == "stop" {
		getDump, err := dump.GetTransferTask(client, dump.GetTransferTaskOpts{
			StreamName: d.Get("stream_name").(string),
			TaskName:   opts.OBSDestinationDescriptor.TaskName,
		})
		if err != nil {
			return fmterr.Errorf("error query OpenTelekomCloud DIS dump task partitions: %s", err)
		}

		err = dump.TransferTaskAction(client, dump.TransferTaskActionOpts{
			StreamName: d.Get("stream_name").(string),
			Action:     d.Get("action").(string),
			Tasks: []dump.BatchTransferTask{
				{
					Id: getDump.TaskId,
				},
			},
		})
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud DIS dump task state: %s", err)
		}
	}

	d.SetId(opts.OBSDestinationDescriptor.TaskName)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDisDumpV2Read(clientCtx, d, meta)
}

func resourceDisDumpV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := dump.GetTransferTaskOpts{
		StreamName: d.Get("stream_name").(string),
		TaskName:   d.Id(),
	}
	getDump, err := dump.GetTransferTask(client, opts)
	if err != nil {
		return fmterr.Errorf("error query OpenTelekomCloud DIS dump task partitions: %s", err)
	}
	var partitions []map[string]interface{}
	for _, partition := range getDump.Partitions {
		partitions = append(partitions, map[string]interface{}{
			"id":                    partition.PartitionId,
			"status":                partition.Status,
			"hash_range":            partition.HashRange,
			"sequence_number_range": partition.SequenceNumberRange,
			"parent_partitions":     partition.ParentPartitions,
		})
	}
	obs := getDump.OBSDestinationDescription
	var obsDescriptor = []map[string]interface{}{
		{
			"task_name":             d.Id(),
			"agency_name":           obs.AgencyName,
			"deliver_time_interval": obs.DeliverTimeInterval,
			"consumer_strategy":     obs.ConsumerStrategy,
			"file_prefix":           obs.FilePrefix,
			"partition_format":      obs.PartitionFormat,
			"obs_bucket_path":       obs.OBSBucketPath,
			"destination_file_type": obs.DestinationFileType,
			"record_delimiter":      obs.RecordDelimiter,
		},
	}
	obsProcessingSchema := getDump.OBSDestinationDescription.ProcessingSchema
	var obsSchema []map[string]interface{}
	if obsProcessingSchema.TimestampName != "" {
		obsSchema = []map[string]interface{}{
			{
				"timestamp_name":   obsProcessingSchema.TimestampName,
				"timestamp_type":   obsProcessingSchema.TimestampType,
				"timestamp_format": obsProcessingSchema.TimestampFormat,
			},
		}
	}

	mErr := multierror.Append(
		d.Set("name", getDump.TaskName),
		d.Set("stream_name", getDump.StreamName),
		d.Set("status", getDump.State),
		d.Set("task_id", getDump.TaskId),
		d.Set("destination", getDump.DestinationType),
		d.Set("obs_destination_descriptor", obsDescriptor),
		d.Set("obs_processing_schema", obsSchema),
		d.Set("created_at", getDump.CreatedAt),
		d.Set("partitions", partitions),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDisDumpV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = dump.DeleteTransferTask(client, dump.DeleteTransferTaskOpts{
		StreamName: d.Get("stream_name").(string),
		TaskName:   d.Id(),
	})
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DIS dump transfer task: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDisDumpV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	if d.HasChange("action") {
		err = dump.TransferTaskAction(client, dump.TransferTaskActionOpts{
			StreamName: d.Get("stream_name").(string),
			Action:     d.Get("action").(string),
			Tasks: []dump.BatchTransferTask{
				{
					Id: d.Get("task_id").(string),
				},
			},
		})
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud DIS dump task state: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDisDumpV2Read(clientCtx, d, meta)
}
