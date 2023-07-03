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
	_ "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/checkpoints"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDisCheckpointV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDisCheckpointV2Create,
		ReadContext:   resourceDisCheckpointV2Read,
		DeleteContext: resourceDisCheckpointV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 200),
					validation.StringMatch(
						regexp.MustCompile(`^[A-Za-z0-9\-_]+$`),
						"Only letters, digits, underscores (_), and hyphens (-) are allowed.",
					),
				),
			},
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[A-Za-z0-9\-_]+$`),
						"Only letters, digits, underscores (_), and hyphens (-) are allowed.",
					),
				),
			},
			"checkpoint_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LAST_READ",
				}, false),
				Default: "LAST_READ",
			},
			"partition_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sequence_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metadata": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		},
	}
}

func resourceDisCheckpointV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := checkpoints.CommitCheckpointOpts{
		AppName:        d.Get("app_name").(string),
		CheckpointType: d.Get("checkpoint_type").(string),
		StreamName:     d.Get("stream_name").(string),
		PartitionId:    d.Get("partition_id").(string),
		SequenceNumber: d.Get("sequence_number").(string),
		Metadata:       d.Get("metadata").(string),
	}

	log.Printf("[DEBUG] Creating new Checkpoint: %s", opts.StreamName)
	err = checkpoints.CommitCheckpoint(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating DIS checkpoint: %s", err)
	}

	d.SetId(opts.StreamName)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDisCheckpointV2Read(clientCtx, d, meta)
}

func resourceDisCheckpointV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := checkpoints.GetCheckpointOpts{
		StreamName:     d.Id(),
		AppName:        d.Get("app_name").(string),
		CheckpointType: checkpointType,
		PartitionId:    d.Get("partition_id").(string),
	}
	checkpoint, err := checkpoints.GetCheckpoint(client, opts)
	if err != nil {
		return fmterr.Errorf("error getting checkpoint details: %s", err)
	}

	mErr := multierror.Append(
		d.Set("sequence_number", checkpoint.SequenceNumber),
		d.Set("metadata", checkpoint.Metadata),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDisCheckpointV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := checkpoints.DeleteCheckpointOpts{
		StreamName:     d.Id(),
		AppName:        d.Get("app_name").(string),
		CheckpointType: checkpointType,
		PartitionId:    d.Get("partition_id").(string),
	}

	err = checkpoints.DeleteCheckpoint(client, opts)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DIS checkpoint: %s", err)
	}

	d.SetId("")
	return nil
}
