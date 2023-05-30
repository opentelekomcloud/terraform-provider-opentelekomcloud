package lts

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/transfers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSTransferV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTransferV2Create,
		ReadContext:   resourceTransferV2Read,
		UpdateContext: resourceTransferV2Update,
		DeleteContext: resourceTransferV2Delete,

		CustomizeDiff: validatePeriods,

		Schema: map[string]*schema.Schema{
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"obs_bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_format": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RAW", "JSON",
				}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					1, 2, 3, 5, 6, 12, 30,
				}),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"min", "hour",
				}, false),
			},
			"switch_on": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dir_prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_transfer_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_transfer_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"obs_encryption_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"obs_encryption_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLTSTransferV2Streams(d *schema.ResourceData) []string {
	streamIdsRaw := d.Get("log_stream_ids").(*schema.Set)
	var streamIds []string

	for _, v := range streamIdsRaw.List() {
		streamIds = append(streamIds, v.(string))
	}
	return streamIds
}

func resourceTransferV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("log_group_id").(string)
	switchOn := d.Get("switch_on").(bool)
	period := d.Get("period").(int)

	createOpts := transfers.CreateLogDumpObsOpts{
		LogGroupId:    groupId,
		LogStreamIds:  resourceLTSTransferV2Streams(d),
		ObsBucketName: d.Get("obs_bucket_name").(string),
		Type:          "cycle",
		StorageFormat: d.Get("storage_format").(string),
		SwitchOn:      &switchOn,
		Period:        int32(period),
		PeriodUnit:    d.Get("period_unit").(string),
		PrefixName:    d.Get("prefix_name").(string),
		DirPrefixName: d.Get("dir_prefix_name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	transferCreate, err := transfers.CreateLogDumpObs(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating log transfer: %s", err)
	}

	d.SetId(transferCreate)
	return resourceTransferV2Read(ctx, d, meta)
}

func resourceTransferV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	allTransfers, err := transfers.ListTransfers(client, transfers.ListTransfersOpts{})
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud log transfer %s: %s", d.Id(), err)
	}

	var transfer transfers.Transfer
	for _, transferRaw := range allTransfers {
		if transferRaw.LogTransferId == d.Id() {
			transfer = transferRaw
			break
		}
	}

	if transfer.LogTransferId == "" {
		return fmterr.Errorf("OpenTelekomCloud log transfer %s was not found", d.Id())
	}

	log.Printf("[DEBUG] Retrieved log transfer %s: %#v", d.Id(), transfer)

	transferStatus := transfer.LogTransferInfo.LogTransferStatus
	var switchOn bool

	switch transferStatus {
	case "ENABLE":
		switchOn = true
	case "DISABLE":
		switchOn = false
	default:
		return fmterr.Errorf("'%s' status received from LTS transfer %s", transferStatus, d.Id())
	}

	mErr := multierror.Append(
		d.Set("log_group_id", transfer.LogGroupId),
		d.Set("log_group_name", transfer.LogGroupName),
		d.Set("log_stream_ids", getLogStreamsID(transfer)),
		d.Set("storage_format", transfer.LogTransferInfo.LogStorageFormat),
		d.Set("log_transfer_mode", transfer.LogTransferInfo.LogTransferMode),
		d.Set("status", transferStatus),
		d.Set("switch_on", switchOn),
		d.Set("log_transfer_type", transfer.LogTransferInfo.LogTransferType),
		d.Set("create_time", transfer.LogTransferInfo.LogCreateTime),
		d.Set("period", transfer.LogTransferInfo.LogTransferDetail.ObsPeriod),
		d.Set("period_unit", transfer.LogTransferInfo.LogTransferDetail.ObsPeriodUnit),
		d.Set("obs_encryption_id", transfer.LogTransferInfo.LogTransferDetail.ObsEncryptedId),
		d.Set("obs_encryption_enable", transfer.LogTransferInfo.LogTransferDetail.ObsEncryptedEnable),
		d.Set("prefix_name", transfer.LogTransferInfo.LogTransferDetail.ObsPrefixName),
		d.Set("dir_prefix_name", transfer.LogTransferInfo.LogTransferDetail.ObsDirPreFixName),
		d.Set("obs_bucket_name", transfer.LogTransferInfo.LogTransferDetail.ObsBucketName),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceTransferV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	updateOpts := transfers.UpdateTransferOpts{
		TransferId: d.Id(),
		TransferInfo: transfers.TransferInfo{
			StorageFormat: d.Get("storage_format").(string),
			TransferDetail: transfers.TransferDetail{
				ObsPeriod:        d.Get("period").(int),
				ObsPeriodUnit:    d.Get("period_unit").(string),
				ObsBucketName:    d.Get("obs_bucket_name").(string),
				ObsPrefixName:    d.Get("prefix_name").(string),
				ObsDirPreFixName: d.Get("dir_prefix_name").(string),
			},
		},
	}

	switchOn := d.Get("switch_on").(bool)
	switch switchOn {
	case false:
		updateOpts.TransferInfo.TransferStatus = "DISABLE"
	default:
		updateOpts.TransferInfo.TransferStatus = "ENABLE"
	}

	log.Printf("[DEBUG]  Options: %#v", updateOpts)

	_, err = transfers.UpdateTransfer(client, updateOpts)
	if err != nil {
		return fmterr.Errorf("error creating log transfer: %s", err)
	}

	return resourceTransferV2Read(ctx, d, meta)
}

func resourceTransferV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	err = transfers.DeleteTransfer(client, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			d.SetId("")
			return nil
		} else {
			return common.CheckDeletedDiag(d, err, "Error deleting log transfer")
		}
	}

	d.SetId("")
	return nil
}

func getLogStreamsID(transfer transfers.Transfer) []string {
	var logStreamIds []string

	for _, stream := range transfer.LogStreams {
		logStreamIds = append(logStreamIds, stream.LogStreamId)
	}
	return logStreamIds
}

func validatePeriods(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	periodUnit := d.Get("period_unit").(string)
	period := d.Get("period").(int)
	periodList := []string{
		"2min", "5min", "30min", "1hour", "3hour", "6hour", "12hour",
	}
	if !common.StringInSlice(strconv.Itoa(period)+periodUnit, periodList) {
		return fmt.Errorf("log transfer interval must be set to one of the following: " +
			"`2 min`, `5 min`, `30 min`, `1 hour`, `3 hours`, `6 hours` or `12 hours`")
	}

	return nil
}
